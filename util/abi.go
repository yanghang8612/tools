package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// EncodeCallData build calldata for eth_call / eth_sendTransaction.
//
// inputSignature can be:
//   - "transfer(address,uint256)"
//   - "function transfer(address to, uint256 amount) external returns (bool)"
//
// params: each argument as a string, e.g.
//
//	["0xabc...", "1e18"]
//	["[0xabc...,0xdef...]", "1e18"]
//	["(0xabc...,1e18)"]  // tuple
func EncodeCallData(inputSignature string, params []string) (string, error) {
	name, typeList, err := parseFunctionLikeInput(inputSignature)
	if err != nil {
		return "", err
	}
	if len(typeList) != len(params) {
		return "", fmt.Errorf("param count mismatch: signature expects %d args, got %d", len(typeList), len(params))
	}

	abiJSON := buildSingleFunctionABIJSON(name, typeList)
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return "", fmt.Errorf("abi.JSON parse failed: %w", err)
	}

	method := parsedABI.Methods[name]

	args := make([]any, len(params))
	for i := range params {
		v, err := convertStringToABIValue(method.Inputs[i].Type, params[i])
		if err != nil {
			return "", fmt.Errorf("convert arg[%d] (%s) failed: %w", i, method.Inputs[i].Type.String(), err)
		}
		args[i] = v
	}

	data, err := parsedABI.Pack(name, args...)
	if err != nil {
		return "", fmt.Errorf("abi.Pack failed: %w", err)
	}
	return hexutil.Encode(data), nil
}

// Optional helper: compute 4-byte selector from signature/function-def input
func Selector(inputSignature string) (string, error) {
	name, typeList, err := parseFunctionLikeInput(inputSignature)
	if err != nil {
		return "", err
	}
	abiJSON := buildSingleFunctionABIJSON(name, typeList)
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(parsedABI.Methods[name].ID), nil
}

/* ------------------------- Signature parsing ------------------------- */

var sigRe = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\((.*)\)$`)
var solFnRe = regexp.MustCompile(`(?s)\bfunction\s+([A-Za-z_][A-Za-z0-9_]*)\s*\((.*)\)`)

func parseFunctionLikeInput(input string) (name string, types []string, err error) {
	in := strings.TrimSpace(input)
	if strings.Contains(in, "function") {
		return parseSolidityFunctionDef(in)
	}
	return parseSignature(in)
}

func parseSignature(signature string) (name string, types []string, err error) {
	m := sigRe.FindStringSubmatch(strings.TrimSpace(signature))
	if m == nil {
		return "", nil, fmt.Errorf("invalid signature: %q", signature)
	}
	name = m[1]
	inside := strings.TrimSpace(m[2])
	if inside == "" {
		return name, []string{}, nil
	}
	types, err = splitCommaRespectNesting(inside)
	if err != nil {
		return "", nil, err
	}
	for i := range types {
		types[i] = strings.TrimSpace(types[i])
	}
	return name, types, nil
}

func parseSolidityFunctionDef(def string) (name string, types []string, err error) {
	def = stripSolidityComments(def)

	m := solFnRe.FindStringSubmatch(def)
	if m == nil {
		return "", nil, fmt.Errorf("invalid solidity function definition: %q", def)
	}
	name = m[1]

	// m[2] is "params) ...", we need only params until the matching ')'
	afterLParen := m[2]
	paramsPart, err := takeUntilMatchingParen(afterLParen)
	if err != nil {
		return "", nil, err
	}
	paramsPart = strings.TrimSpace(paramsPart)
	if paramsPart == "" {
		return name, []string{}, nil
	}

	params, err := splitCommaRespectNesting(paramsPart)
	if err != nil {
		return "", nil, err
	}

	types = make([]string, 0, len(params))
	for _, p := range params {
		ty, err := solidityParamToABIType(p)
		if err != nil {
			return "", nil, fmt.Errorf("param %q: %w", p, err)
		}
		types = append(types, ty)
	}
	return name, types, nil
}

func buildSingleFunctionABIJSON(name string, types []string) string {
	// minimal ABI JSON with inputs only; outputs not required for calldata encoding
	// NOTE: abi.JSON requires valid JSON format
	var b strings.Builder
	b.WriteString(`[{"type":"function","name":"`)
	b.WriteString(name)
	b.WriteString(`","stateMutability":"nonpayable","inputs":[`)
	for i, t := range types {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(`{"name":"a`)
		b.WriteString(strconv.Itoa(i))
		b.WriteString(`","type":"`)
		b.WriteString(t)
		b.WriteString(`"}`)
	}
	b.WriteString(`],"outputs":[]}]`)
	return b.String()
}

func stripSolidityComments(s string) string {
	// remove /* ... */ best-effort
	for {
		start := strings.Index(s, "/*")
		if start < 0 {
			break
		}
		end := strings.Index(s[start+2:], "*/")
		if end < 0 {
			s = s[:start]
			break
		}
		s = s[:start] + s[start+2+end+2:]
	}
	// remove //... to endline
	lines := strings.Split(s, "\n")
	for i := range lines {
		if idx := strings.Index(lines[i], "//"); idx >= 0 {
			lines[i] = lines[i][:idx]
		}
	}
	return strings.Join(lines, "\n")
}

// after "(" already consumed: s = "a,b) external ..."
// depth starts at 1, return content until the matching ")"
func takeUntilMatchingParen(s string) (string, error) {
	depth := 1
	var buf bytes.Buffer
	inStr := byte(0)

	for i := 0; i < len(s); i++ {
		c := s[i]

		if inStr != 0 {
			buf.WriteByte(c)
			if c == inStr && (i == 0 || s[i-1] != '\\') {
				inStr = 0
			}
			continue
		}

		switch c {
		case '"', '\'':
			inStr = c
			buf.WriteByte(c)
		case '(':
			depth++
			buf.WriteByte(c)
		case ')':
			depth--
			if depth == 0 {
				return buf.String(), nil
			}
			buf.WriteByte(c)
		default:
			buf.WriteByte(c)
		}
	}
	return "", fmt.Errorf("unbalanced parentheses in solidity function definition")
}

func solidityParamToABIType(param string) (string, error) {
	p := strings.TrimSpace(param)
	if p == "" {
		return "", fmt.Errorf("empty param")
	}
	p = strings.Join(strings.Fields(p), " ")

	// strip common location qualifiers
	for _, q := range []string{" memory ", " calldata ", " storage "} {
		p = strings.ReplaceAll(p, q, " ")
	}
	for _, q := range []string{" memory", " calldata", " storage"} {
		p = strings.ReplaceAll(p, q, "")
	}

	// strip payable keyword (address payable)
	p = strings.ReplaceAll(p, " payable ", " ")
	p = strings.ReplaceAll(p, " payable", "")

	// tuple type starts with '('
	if strings.HasPrefix(strings.TrimSpace(p), "(") {
		tupleType, _, err := takeTupleType(p)
		if err != nil {
			return "", err
		}
		return tupleType, nil
	}

	// otherwise first token is type (may include []/[k])
	parts := strings.SplitN(p, " ", 2)
	ty := strings.TrimSpace(parts[0])
	if ty == "" {
		return "", fmt.Errorf("cannot parse type from %q", param)
	}
	return ty, nil
}

func takeTupleType(s string) (tupleType string, rest string, err error) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "(") {
		return "", "", fmt.Errorf("tuple type must start with '('")
	}
	depth := 0
	i := 0
	inStr := byte(0)

	for i < len(s) {
		c := s[i]
		if inStr != 0 {
			if c == inStr && (i == 0 || s[i-1] != '\\') {
				inStr = 0
			}
			i++
			continue
		}

		switch c {
		case '"', '\'':
			inStr = c
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				i++ // include ')'
				j := i
				// include any immediate array suffix [..][..]...
				for j < len(s) && s[j] == '[' {
					k := strings.IndexByte(s[j:], ']')
					if k < 0 {
						return "", "", fmt.Errorf("invalid array suffix in tuple type: %q", s)
					}
					j = j + k + 1
				}
				return strings.TrimSpace(s[:j]), s[j:], nil
			}
		}
		i++
	}
	return "", "", fmt.Errorf("unbalanced tuple type: %q", s)
}

/* ------------------------- Value parsing & conversion ------------------------- */

func convertStringToABIValue(t abi.Type, s string) (any, error) {
	s = strings.TrimSpace(s)

	// Arrays/slices: prefer [...] syntax
	if t.T == abi.SliceTy || t.T == abi.ArrayTy {
		elems, err := parseBracketListValue(s, '[', ']')
		if err != nil {
			// fallback: allow bare "a,b" for convenience (you说你习惯 a,b)
			elems, err = splitCommaRespectNesting(s)
			if err != nil {
				return nil, err
			}
		}

		elemType := *t.Elem
		goTy := t.GetType() // reflect.Type

		if t.T == abi.SliceTy {
			slice := reflect.MakeSlice(goTy, len(elems), len(elems))
			for i := range elems {
				v, err := convertStringToABIValue(elemType, elems[i])
				if err != nil {
					return nil, fmt.Errorf("array elem[%d]: %w", i, err)
				}
				sv := reflect.ValueOf(v)
				slice.Index(i).Set(coerceValue(sv, goTy.Elem()))
			}
			return slice.Interface(), nil
		}

		// fixed array
		if len(elems) != t.Size {
			return nil, fmt.Errorf("fixed array expects %d elems, got %d", t.Size, len(elems))
		}
		arr := reflect.New(goTy).Elem()
		for i := range elems {
			v, err := convertStringToABIValue(elemType, elems[i])
			if err != nil {
				return nil, fmt.Errorf("array elem[%d]: %w", i, err)
			}
			av := reflect.ValueOf(v)
			arr.Index(i).Set(coerceValue(av, goTy.Elem()))
		}
		return arr.Interface(), nil
	}

	switch t.T {
	case abi.AddressTy:
		if !common.IsHexAddress(s) {
			return nil, fmt.Errorf("invalid address: %s", s)
		}
		return common.HexToAddress(s), nil

	case abi.BoolTy:
		b, err := strconv.ParseBool(strings.ToLower(s))
		if err != nil {
			return nil, fmt.Errorf("invalid bool: %s", s)
		}
		return b, nil

	case abi.StringTy:
		// allow quoted strings
		return trimOptionalQuotes(s), nil

	case abi.BytesTy:
		// dynamic bytes: accept 0x.. hex; otherwise treat as UTF-8 bytes
		if strings.HasPrefix(s, "0x") {
			b, err := hexutil.Decode(s)
			if err != nil {
				return nil, err
			}
			return b, nil
		}
		return []byte(trimOptionalQuotes(s)), nil

	case abi.FixedBytesTy:
		n := t.Size
		var b []byte
		var err error
		if strings.HasPrefix(s, "0x") {
			b, err = hexutil.Decode(s)
		} else {
			b = []byte(trimOptionalQuotes(s))
		}
		if err != nil {
			return nil, err
		}
		if len(b) != n {
			return nil, fmt.Errorf("bytes%d expects %d bytes, got %d", n, n, len(b))
		}
		arr := reflect.New(t.GetType()).Elem()
		reflect.Copy(arr, reflect.ValueOf(b))
		return arr.Interface(), nil

	case abi.IntTy, abi.UintTy:
		bi, err := parseBigIntExtended(s, t.T == abi.UintTy)
		if err != nil {
			return nil, err
		}
		return bi, nil

	case abi.TupleTy:
		// tuple value syntax: (a,b, ...)
		fields, err := parseParenListValue(s, '(', ')')
		if err != nil {
			return nil, fmt.Errorf("tuple expects '(...)' syntax: %w", err)
		}
		if len(fields) != len(t.TupleElems) {
			return nil, fmt.Errorf("tuple expects %d fields, got %d", len(t.TupleElems), len(fields))
		}

		goTy := t.GetType() // struct type
		v := reflect.New(goTy).Elem()

		for i := range fields {
			cv, err := convertStringToABIValue(*t.TupleElems[i], fields[i])
			if err != nil {
				return nil, fmt.Errorf("tuple field[%d]: %w", i, err)
			}
			f := v.Field(i)
			fv := reflect.ValueOf(cv)
			f.Set(coerceValue(fv, f.Type()))
		}
		return v.Interface(), nil

	default:
		return nil, fmt.Errorf("unsupported abi type: %s", t.String())
	}
}

func coerceValue(v reflect.Value, target reflect.Type) reflect.Value {
	// handle assignable / convertible
	if !v.IsValid() {
		return reflect.Zero(target)
	}
	if v.Type().AssignableTo(target) {
		return v
	}
	if v.Type().ConvertibleTo(target) {
		return v.Convert(target)
	}

	// common case: big.Int vs *big.Int
	if target == reflect.TypeOf(big.Int{}) {
		if v.Type() == reflect.TypeOf(&big.Int{}) {
			return v.Elem()
		}
	}
	if target == reflect.TypeOf(&big.Int{}) {
		if v.Type() == reflect.TypeOf(big.Int{}) {
			b := new(big.Int)
			b.Set(v.Addr().Interface().(*big.Int))
			return reflect.ValueOf(b)
		}
	}

	// last resort: try to take address if needed
	if target.Kind() == reflect.Ptr && v.CanAddr() && v.Addr().Type().AssignableTo(target) {
		return v.Addr()
	}

	// will panic later if we return wrong type; better to be explicit here
	panic(fmt.Sprintf("cannot coerce %s to %s", v.Type().String(), target.String()))
}

/* ------------------------- Your preferred list syntax ------------------------- */

// Parse "[a,b,(c,d),[e,f]]" -> []string{"a","b","(c,d)","[e,f]"}
func parseBracketListValue(s string, open, close byte) ([]string, error) {
	s = strings.TrimSpace(s)
	if len(s) < 2 || s[0] != open || s[len(s)-1] != close {
		return nil, fmt.Errorf("not a bracketed list")
	}
	inner := strings.TrimSpace(s[1 : len(s)-1])
	if inner == "" {
		return []string{}, nil
	}
	return splitCommaRespectNesting(inner)
}

func parseParenListValue(s string, open, close byte) ([]string, error) {
	s = strings.TrimSpace(s)
	if len(s) < 2 || s[0] != open || s[len(s)-1] != close {
		return nil, fmt.Errorf("not a parenthesized list")
	}
	inner := strings.TrimSpace(s[1 : len(s)-1])
	if inner == "" {
		return []string{}, nil
	}
	return splitCommaRespectNesting(inner)
}

// Split by commas while respecting nested (), [], and quotes.
func splitCommaRespectNesting(s string) ([]string, error) {
	var out []string
	var buf bytes.Buffer

	depthParen := 0
	depthBrack := 0
	inStr := byte(0)

	for i := 0; i < len(s); i++ {
		c := s[i]

		if inStr != 0 {
			buf.WriteByte(c)
			if c == inStr && (i == 0 || s[i-1] != '\\') {
				inStr = 0
			}
			continue
		}

		switch c {
		case '"', '\'':
			inStr = c
			buf.WriteByte(c)

		case '(':
			depthParen++
			buf.WriteByte(c)
		case ')':
			depthParen--
			if depthParen < 0 {
				return nil, fmt.Errorf("unbalanced ')': %q", s)
			}
			buf.WriteByte(c)

		case '[':
			depthBrack++
			buf.WriteByte(c)
		case ']':
			depthBrack--
			if depthBrack < 0 {
				return nil, fmt.Errorf("unbalanced ']': %q", s)
			}
			buf.WriteByte(c)

		case ',':
			if depthParen == 0 && depthBrack == 0 {
				part := strings.TrimSpace(buf.String())
				out = append(out, part)
				buf.Reset()
			} else {
				buf.WriteByte(c)
			}

		default:
			buf.WriteByte(c)
		}
	}

	if inStr != 0 {
		return nil, fmt.Errorf("unclosed string in: %q", s)
	}
	if depthParen != 0 || depthBrack != 0 {
		return nil, fmt.Errorf("unbalanced nesting in: %q", s)
	}

	last := strings.TrimSpace(buf.String())
	if last != "" {
		out = append(out, last)
	} else if len(out) == 0 {
		// s might be empty or spaces
	}
	return out, nil
}

/* ------------------------- Numbers & helpers ------------------------- */

// Supports:
//   - decimal: 123
//   - hex: 0xFF
//   - scientific: 1e18, 1.5e18, 1_000e18
func parseBigIntExtended(s string, unsigned bool) (*big.Int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty int")
	}
	s = strings.ReplaceAll(s, "_", "")

	// hex
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		hs := s[2:]
		if hs == "" {
			return nil, fmt.Errorf("invalid hex int")
		}
		bi := new(big.Int)
		_, ok := bi.SetString(hs, 16)
		if !ok {
			return nil, fmt.Errorf("invalid hex int: %s", s)
		}
		if unsigned && bi.Sign() < 0 {
			return nil, fmt.Errorf("uint cannot be negative: %s", s)
		}
		return bi, nil
	}

	// scientific notation base10
	if strings.ContainsAny(s, "eE") {
		parts := strings.SplitN(strings.ToLower(s), "e", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid scientific notation: %s", s)
		}
		mantStr := strings.TrimSpace(parts[0])
		expStr := strings.TrimSpace(parts[1])
		if mantStr == "" || expStr == "" {
			return nil, fmt.Errorf("invalid scientific notation: %s", s)
		}
		exp, err := strconv.ParseInt(expStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid exponent: %s", expStr)
		}

		mant := new(big.Rat)
		if _, ok := mant.SetString(mantStr); !ok {
			return nil, fmt.Errorf("invalid mantissa: %s", mantStr)
		}

		scaleInt := new(big.Int).Exp(big.NewInt(10), big.NewInt(absInt64(exp)), nil)
		scale := new(big.Rat).SetInt(scaleInt)

		if exp >= 0 {
			mant.Mul(mant, scale)
		} else {
			mant.Quo(mant, scale)
		}

		if mant.Denom().Cmp(big.NewInt(1)) != 0 {
			return nil, fmt.Errorf("scientific value is not integer after scaling: %s", s)
		}
		bi := new(big.Int).Set(mant.Num())
		if unsigned && bi.Sign() < 0 {
			return nil, fmt.Errorf("uint cannot be negative: %s", s)
		}
		return bi, nil
	}

	// plain decimal
	bi := new(big.Int)
	_, ok := bi.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid int: %s", s)
	}
	if unsigned && bi.Sign() < 0 {
		return nil, fmt.Errorf("uint cannot be negative: %s", s)
	}
	return bi, nil
}

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func trimOptionalQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
