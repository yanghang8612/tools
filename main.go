package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.HideHelp = true
	app.Copyright = "Copyright 2021-2021 Asuka"
	app.Usage = "very useful tool kits for Asuka"
	app.CustomAppHelpTemplate = ""
	app.Commands = []*cli.Command{
		{
			Name:  "db",
			Usage: "Database related commands",
			Subcommands: []*cli.Command{
				&dbCountCommand,
				&dbGetCommand,
				&dbRootCommand,
				&dbPrintCommand,
				&dbDiffCommand,
			},
		},
		{
			Name:  "addr",
			Usage: "Address related commands",
			Subcommands: []*cli.Command{
				&addrNumCommand,
				&addrDecodeCommand,
				&addrEncodeCommand,
			},
		},
		{
			Name:  "vm",
			Usage: "EVM related commands",
			Subcommands: []*cli.Command{
				&vmPadCommand,
				&vmSplitCommand,
				&vmUnpackCommand,
				&vm4bytesCommand,
			},
		},
		{
			Name:  "scan",
			Usage: "TronScan related commands",
			Subcommands: []*cli.Command{
				&txsCommand,
				&txCommand,
				&speedCommand,
				&transferCommand,
			},
		},
		{
			Name:  "eth",
			Usage: "ETH JSON-RPC related commands",
			Subcommands: []*cli.Command{
				&logsCommand,
			},
		},
	}

	for _, cmd := range app.Commands {
		cmd.HideHelpCommand = true
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
