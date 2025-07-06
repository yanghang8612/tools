package main

import (
	"errors"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli/v2"
	"tools/log"
	utils "tools/util"
)

var (
	signCommand = cli.Command{
		Name:  "sign",
		Usage: "Sign a message with private key",
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				return errors.New("sign subcommand needs msg-hash and private-key args")
			}
			msgHash, ok := utils.FromHex(c.Args().Get(0))
			if ok {
				if len(msgHash) != 32 {
					return errors.New("msg-hash must be 32 bytes")
				}
			} else {
				return errors.New("msg-hash must be in hex format")
			}

			privateKey, ok := utils.FromHex(c.Args().Get(1))
			if ok {
				if len(privateKey) != 32 {
					return errors.New("private-key must be 32 bytes")
				}
			} else {
				return errors.New("private-key must be in hex format")
			}

			pub, err := crypto.ToECDSA(privateKey)

			if err != nil {
				return err
			}

			sig, err := crypto.Sign(msgHash, pub)
			if err != nil {
				return err
			}

			addr := crypto.PubkeyToAddress(pub.PublicKey)
			log.NewLog("eth addr", addr.String())
			log.NewLog("tron addr", base58.CheckEncode(addr.Bytes(), 0x41))
			log.NewLog("signature", common.Bytes2Hex(sig))

			return nil
		},
	}
	recoverCommand = cli.Command{
		Name:  "recover",
		Usage: "Recover address from signature",
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				return errors.New("recover subcommand needs msg-hash and sig args")
			}
			msgHash, ok := utils.FromHex(c.Args().Get(0))
			if ok {
				if len(msgHash) != 32 {
					return errors.New("msg-hash must be 32 bytes")
				}
			} else {
				return errors.New("msg-hash must be in hex format")
			}

			sig, ok := utils.FromHex(c.Args().Get(1))
			if ok {
				if len(sig) != 65 {
					return errors.New("signature must be 65 bytes")
				}
			} else {
				return errors.New("signature must be in hex format")
			}

			pub, err := crypto.SigToPub(msgHash, sig)

			if err != nil {
				return err
			}

			addr := crypto.PubkeyToAddress(*pub)
			log.NewLog("eth addr", addr.String())
			log.NewLog("tron addr", base58.CheckEncode(addr.Bytes(), 0x41))

			return nil
		},
	}
)
