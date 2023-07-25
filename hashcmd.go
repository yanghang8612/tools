package main

import (
	"crypto/sha256"
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli/v2"
	"tools/log"
	utils "tools/util"
)

var (
	hashCommand = cli.Command{
		Name:  "hash",
		Usage: "Hash data by keccak256, sha256",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("hash command only needs one arg")
			}
			arg0 := c.Args().Get(0)
			data, ok := utils.FromHex(arg0)
			if !ok {
				data = []byte(arg0)
			}
			log.NewLog("sha256", sha256.Sum256(data))
			log.NewLog("keccak256", crypto.Keccak256(data))
			return nil
		},
	}
)
