package main

import (
	"errors"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli/v2"
	"tools/log"
)

var (
	recoverCommand = cli.Command{
		Name:  "recover",
		Usage: "Recover address from signature",
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				return errors.New("recover subcommand needs msg-hash sig args")
			}
			msgHash := common.FromHex(c.Args().Get(0))
			sig := common.FromHex(c.Args().Get(1))

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
