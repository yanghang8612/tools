package main

import (
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"tools/log"
	utils "tools/util"
)

var (
	nowCommand = cli.Command{
		Name:  "now",
		Usage: "Convert time between datetime and timestamp",
		Action: func(c *cli.Context) error {
			// display now
			if c.NArg() == 0 {
				log.NewLog("in sec", time.Now().Unix())
				log.NewLog("in milli", time.Now().UnixMilli())
				log.NewLog("in datetime", time.Now())
			} else {
				arg := c.Args().Get(0)

				var (
					err error
					ts  uint64
				)

				num, err := strconv.Atoi(arg)
				if err == nil {
					ts = uint64(num)
				} else {
					data, ok := utils.FromHex(arg)
					if ok && len(data) <= 6 {
						ts = big.NewInt(0).SetBytes(data).Uint64()
					}
				}
				// input str is a valid timestamp
				if ts > 0 {
					var dt time.Time
					dt = time.Unix(int64(ts), 0)
					// only display date in 21st Century
					if dt.Year() >= 2000 && dt.Year() <= 2100 {
						log.NewLog("in sec", dt.Format("2006-01-02 15:04:05"))
					}
					dt = time.Unix(int64(ts/1000), 0)
					// only display date in 21st Century
					if dt.Year() >= 2000 && dt.Year() <= 2100 {
						log.NewLog("in milli", dt.Format("2006-01-02 15:04:05"))
					}
				} else {
					// input str is date or time
					loc, _ := time.LoadLocation("Asia/Shanghai")
					formatsWithDateButNoYear := []string{"1-2 15:4:5", "1-2 15:4", "1-2 15", "1-2"}
					formatsWithoutDate := []string{"15:4:5", "15:4"}
					formats := formatsWithDateButNoYear
					// append two kinds of year to formats
					for _, f := range formatsWithDateButNoYear {
						formats = append(formats, "2006-"+f)
						formats = append(formats, "06-"+f)
					}
					formats = append(formats, formatsWithoutDate...)
					for _, format := range formats {
						if dt, err := time.ParseInLocation(format, arg, loc); err == nil {
							if !strings.ContainsAny(format, "-") {
								dt = dt.AddDate(time.Now().Year(), int(time.Now().Month())-1, time.Now().Day()-1)
							}
							if dt.Year() == 0 {
								dt = dt.AddDate(time.Now().Year(), 0, 0)
							}
							log.NewLog("in sec", dt.Unix())
							log.NewLog("in milli", dt.UnixMilli())
							log.NewLog("in datetime", dt)
						}
					}
				}
			}
			return nil
		},
	}
)
