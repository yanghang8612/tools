package main

import (
    tron "tools/proto"

    "bytes"
    "encoding/binary"
    "encoding/hex"
    "errors"
    "fmt"
    "reflect"
    "strconv"
    "strings"
    "time"

    "github.com/btcsuite/btcutil/base58"
    "github.com/ethereum/go-ethereum/common"
    "github.com/golang/protobuf/proto"
    "github.com/syndtr/goleveldb/leveldb"
    "github.com/syndtr/goleveldb/leveldb/opt"
    "github.com/urfave/cli/v2"
    "golang.org/x/crypto/sha3"
)

const (
    LatestBlockNumKey  = "latest_block_header_number"
    LatestBlockHashKey = "latest_block_header_hash"
)

var (
    dbValueTypeFlag = &cli.StringFlag{Name: "type, t"}
    dbCountCommand  = cli.Command{
        Name:  "count",
        Usage: "Count the total items for given name db",
        Action: func(c *cli.Context) error {
            if c.NArg() != 1 {
                return errors.New("count command needs db path arg")
            }
            return countDb(c.Args().Get(0))
        },
    }
    dbGetCommand = cli.Command{
        Name:  "get",
        Usage: "Get value of the given key in db",
        Flags: []cli.Flag{
            dbValueTypeFlag,
        },
        Subcommands: []*cli.Command{
            {
                Name:  "blocknum",
                Usage: "Get current block number",
                Action: func(c *cli.Context) error {
                    if value, err := queryValue("properties", []byte(LatestBlockNumKey)); err == nil {
                        fmt.Printf("Key `%s` int value is %d\n",
                            LatestBlockNumKey,
                            int64(binary.BigEndian.Uint64(value)))
                        return nil
                    } else {
                        return err
                    }
                },
            },
            {
                Name:  "blockhash",
                Usage: "Get current block hash",
                Action: func(c *cli.Context) error {
                    if value, err := queryValue("properties", []byte(LatestBlockHashKey)); err == nil {
                        fmt.Printf("Key `%s` hex value is %x\n", LatestBlockHashKey, value)
                        return nil
                    } else {
                        return err
                    }
                },
            },
        },
        Action: func(c *cli.Context) error {
            if c.NArg() < 2 {
                return errors.New("get subcommand needs db-path and db-key args")
            }
            dbPath := c.Args().Get(0)
            key := c.Args().Get(1)
            var dbKey []byte
            if strings.ContainsAny(key, "0x") {
                if hexKey, err := hex.DecodeString(key[2:]); err == nil {
                    dbKey = hexKey
                } else {
                    return err
                }
            } else if strings.HasPrefix(key, "T") && len(key) == 34 {
                if decoded, _, err := base58.CheckDecode(key); err == nil {
                    dbKey = decoded
                } else {
                    return err
                }
            } else {
                dbKey = []byte(key)
            }
            if value, err := queryValue(dbPath, dbKey); err != nil {
                return err
            } else {
                outputType := c.String("type")
                switch outputType {
                case "num", "number", "int", "int32", "int64":
                    fmt.Printf("Key `%s` int value is %d\n", key, int64(binary.BigEndian.Uint64(value)))
                case "account", "acc":
                    account := &tron.Account{}
                    if proto.UnmarshalMerge(value, account) == nil {
                        fmt.Printf("%s", proto.MarshalTextString(account))
                    }
                case "hex":
                default:
                    fmt.Printf("Key `%s` hex value is %s\n", key, hex.EncodeToString(value))
                }
                return nil
            }
        },
    }
    dbRootCommand = cli.Command{
        Name:  "hash",
        Usage: "Calculate the hash for given name db",
        Action: func(c *cli.Context) error {
            if c.NArg() != 1 {
                return errors.New("hash subcommand needs db path arg")
            }
            if root, err := calcHash(c.Args().Get(0)); err == nil {
                fmt.Printf("Root is %s\n", hex.EncodeToString(root))
                return nil
            } else {
                return err
            }
        },
    }
    dbPrintCommand = cli.Command{
        Name:  "print",
        Usage: "Print all key-value for given name db",
        Action: func(c *cli.Context) error {
            if c.NArg() != 1 {
                return errors.New("print subcommand needs db path arg")
            }
            return printDb(c.Args().Get(0))
        },
    }
    dbDiffCommand = cli.Command{
        Name:  "diff",
        Usage: "Diff for given db-A and db-B",
        Action: func(c *cli.Context) error {
            if c.NArg() != 2 {
                return errors.New("diff subcommand needs db-A and db-B path args")
            }
            return diffDb(c.Args().Get(0), c.Args().Get(1))
        },
    }
)

func dbOptions() *opt.Options {
    return &opt.Options{
        ErrorIfMissing: true,
        ReadOnly:       true,
    }
}

func countDb(dbPath string) error {
    db, err := leveldb.OpenFile(dbPath, dbOptions())
    if err != nil {
        return err
    }
    defer db.Close()

    var done = make(chan int)
    go func() {
        dot := 0
        for {
            select {
            case count := <-done:
                fmt.Printf("\rDB items count: %d, cost: %ds\n", count, dot)
                return
            default:
                note := "\rCounting"
                for i := 0; i < 5; i++ {
                    if i <= dot%5 {
                        note += "."
                    } else {
                        note += " "
                    }
                }
                dot += 1
                note += "   " + strconv.Itoa(dot) + "s"
                fmt.Print(note)
                time.Sleep(time.Second)
            }
        }
    }()

    count := 0
    zero := 0
    itr := db.NewIterator(nil, nil)
    defer itr.Release()
    for itr.Next() {
        count += 1
        if allZero(itr.Value()) {
            zero += 1
        }
    }
    done <- count
    fmt.Printf("Zero count: %d\n", zero)
    return itr.Error()
}

func allZero(s []byte) bool {
    for _, v := range s {
        if v != 0 {
            return false
        }
    }
    return true
}

func queryValue(dbPath string, key []byte) ([]byte, error) {
    db, err := leveldb.OpenFile(dbPath, dbOptions())
    if err != nil {
        return nil, err
    }
    defer db.Close()

    if value, err := db.Get(key, nil); err == nil {
        return value, nil
    } else {
        return nil, err
    }
}

func calcHash(dbPath string) ([]byte, error) {
    db, err := leveldb.OpenFile(dbPath, dbOptions())
    if err != nil {
        return nil, err
    }
    defer db.Close()

    blackhole, _ := hex.DecodeString("4177944D19C052B73EE2286823AA83F8138CB7032F")
    hash := common.Hash{}.Bytes()
    itr := db.NewIterator(nil, nil)
    defer itr.Release()
    for itr.Next() {
        if !bytes.Equal(itr.Key(), blackhole) {
            //fmt.Printf("%s\n", hex.EncodeToString(itr.Key()))
            hasher := sha3.NewLegacyKeccak256()
            hasher.Write(hash)
            hasher.Write(itr.Key())
            hasher.Write(itr.Value())
            hash = hasher.Sum(nil)
        }
    }
    return hash, itr.Error()
}

func printDb(dbPath string) error {
    db, err := leveldb.OpenFile(dbPath, dbOptions())
    if err != nil {
        return err
    }
    defer db.Close()

    itr := db.NewIterator(nil, nil)
    defer itr.Release()
    for itr.Next() {
        fmt.Printf("%x\n", itr.Key())
    }
    return itr.Error()
}

func diffDb(dbAPath, dbBPath string) error {
    dbA, err := leveldb.OpenFile(dbAPath, dbOptions())
    if err != nil {
        return err
    }
    defer dbA.Close()
    dbB, err := leveldb.OpenFile(dbBPath, dbOptions())
    if err != nil {
        return err
    }
    defer dbB.Close()

    itr := dbA.NewIterator(nil, nil)
    defer itr.Release()
    totalCount, notFoundCount := 0, 0
    for itr.Next() {
        totalCount += 1
        key, aValue := itr.Key(), itr.Value()
        if bValue, err := dbB.Get(key, nil); err == nil {
            if !reflect.DeepEqual(aValue, bValue) {
                fmt.Printf("Different: %x\n", key)
            }
        } else {
            notFoundCount += 1
            fmt.Printf("%s: %x\n", err.Error(), itr.Key())
        }
    }
    fmt.Printf("Total: %d, Not Found: %d\n", totalCount, notFoundCount)
    return itr.Error()
}
