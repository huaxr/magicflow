// Author: huaxr
// Time: 2022/7/1 5:58 下午
// Git: huaxr

package console

import (
	"sort"

	"github.com/huaxr/magicflow/component/helper/console/flags"
	"github.com/urfave/cli/v2"
)

var app = NewApp("the mAgIcfOw command line interface")

// NewApp creates an app with sane defaults.
func NewApp(usage string) *cli.App {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Version = "vv-vv"
	app.Usage = usage
	app.Copyright = "Copyright 2021-2022 The Tal MagicFlow Authors"
	app.Before = func(ctx *cli.Context) error {
		return nil
	}
	return app
}

func LaunchApp(flow func(*cli.Context) error) *cli.App {
	app.Action = nil
	app.HideVersion = true
	app.Commands = append(flags.GetCommand(), &cli.Command{
		Action:      flow,
		Name:        "start",
		Usage:       "start the service command",
		ArgsUsage:   "",
		Flags:       nil,
		Description: "setup your app in your env",
	})

	sort.Sort(cli.CommandsByName(app.Commands))
	// without those flags, the flag parse will not work.
	app.Flags = flags.GetFlags()
	sort.Sort(cli.FlagsByName(app.Flags))

	app.Before = func(ctx *cli.Context) error {
		return nil
	}
	app.After = func(ctx *cli.Context) error {
		return nil
	}

	return app
}
