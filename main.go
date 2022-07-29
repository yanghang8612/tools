package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	nowCommand = cli.Command{
		Name:  "now",
		Usage: "CaoJiaJin like command",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				fmt.Println(time.Now().Unix())
				fmt.Println(time.Now().UnixMilli())
				fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
			} else {
				arg := c.Args().Get(0)
				ts, err := strconv.Atoi(arg)
				if err == nil {
					fmt.Println(time.Unix(int64(ts), 0).Format("2006-01-02 15:04:05"))
				} else {
					loc, _ := time.LoadLocation("Asia/Shanghai")
					var dt time.Time
					var err error
					dt, err = time.ParseInLocation("2006-01-02 15:04:05", arg, loc)
					if err != nil {
						dt, err = time.ParseInLocation("2006-01-02", arg, loc)
					}
					if err != nil {
						fmt.Println("Date format err.")
					} else {
						fmt.Println(dt.Unix())
					}
				}
			}
			return nil
		},
	}
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
		&nowCommand,
	}

	for _, cmd := range app.Commands {
		cmd.HideHelpCommand = true
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
