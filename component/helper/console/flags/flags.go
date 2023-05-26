// Author: huaxr
// Time: 2022/7/7 10:05 上午
// Git: huaxr

package flags

import "github.com/urfave/cli/v2"

var (
	dirFlag = &cli.StringFlag{
		Name:    "dir",
		Aliases: []string{"d"},
		Value:   "/tmp",
		Usage:   "config dir path",
	}
)

func GetFlags() []cli.Flag {
	return []cli.Flag{
		dirFlag,
	}
}
