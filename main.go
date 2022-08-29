package main

import (
    "tools/log"
    "tools/util"

    "fmt"
    "go/token"
    "go/types"
    "math/big"
    "os"
    "path/filepath"
    "regexp"
    "strings"

    "github.com/urfave/cli/v2"
)

func main() {
    app := cli.NewApp()
    app.Name = filepath.Base(os.Args[0])
    app.HideHelp = true
    app.Copyright = "Copyright 2021-2022 Asuka, jeancky"
    app.Usage = "very useful tool kits for TRON and ethereum"
    app.CustomAppHelpTemplate = ""
    app.Action = func(c *cli.Context) error {
        if c.NArg() != 1 {
            return cli.ShowAppHelp(c)
        }
        arg0 := c.Args().Get(0)
        // input can be address
        _ = hexAddrCommand.Action(c)
        // input is in hex
        if data, ok := utils.FromHex(arg0); ok {
            // 0 ~ 31 bytes, pad it to 32 bytes
            if len(data) < 32 {
                _ = abiPadCommand.Action(c)
            } else if len(data) > 32 {
                // > 32 bytes, can be call data
                _ = abiSplitCommand.Action(c)
            }
            // output its value in decimal and try to display its readable string
            _ = hexIntCommand.Action(c)
            _ = hexStrCommand.Action(c)
        } else {
            // input is not hex
            // check if input is in decimal
            if num, ok := new(big.Int).SetString(arg0, 10); ok {
                if num.Cmp(big.NewInt(256)) <= 0 {
                    _ = hexMaxCommand.Action(c)
                }
                // arg is decimal num
                _ = nowCommand.Action(c)
                _ = hexIntCommand.Action(c)
                // try to pad it to 32bytes
                _ = abiPadCommand.Action(c)
            } else {
                // it may be a string
                // check if arg is date
                if strings.ContainsAny(arg0, "-,:") {
                    _ = nowCommand.Action(c)
                }
                // check if arg is function or event
                if matched, _ := regexp.MatchString(`^\w.*\(.*\)$`, arg0); matched {
                    _ = abi4bytesCommand.Action(c)
                    return nil
                }
                // try to eval arg
                ee := regexp.MustCompile(`1e\d+`)
                for _, e := range ee.FindAllString(arg0, -1) {
                    bigfloat := new(big.Float)
                    bigfloat.SetString(e)
                    bigint := new(big.Int)
                    bigfloat.Int(bigint)
                    arg0 = strings.ReplaceAll(arg0, e, bigint.String())
                }
                if res, err := types.Eval(token.NewFileSet(), nil, token.NoPos, arg0); err == nil {
                    log.NewLog("eval result", res.Value.String())
                    return nil
                }
                // otherwise, convert str to hex
                return hexStrCommand.Action(c)
            }
        }
        return nil
    }
    app.Commands = []*cli.Command{
        &callCommand,
        &nowCommand,
        {
            Name:  "abi",
            Usage: "ABI related commands",
            Subcommands: []*cli.Command{
                &abiPadCommand,
                &abiSplitCommand,
                &abiUnpackCommand,
                &abi4bytesCommand,
            },
        },
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
            Name:  "eth",
            Usage: "ETH JSON-RPC related commands",
            Subcommands: []*cli.Command{
                &logsCommand,
            },
        },
        {
            Name:  "hex",
            Usage: "Hex related commands",
            Subcommands: []*cli.Command{
                &hexAddrCommand,
                &hexIntCommand,
                &hexMaxCommand,
                &hexStrCommand,
            },
        },
        {
            Name:  "scan",
            Usage: "TronScan related commands",
            Subcommands: []*cli.Command{
                &txsCommand,
                &txCommand,
            },
        },
    }
    app.After = func(c *cli.Context) error {
        log.FlushLogsToConsole()
        return nil
    }

    for _, cmd := range app.Commands {
        cmd.HideHelpCommand = true
    }

    if err := app.Run(os.Args); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
