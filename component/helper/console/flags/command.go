// Author: huaxr
// Time: 2022/7/1 6:16 下午
// Git: huaxr

package flags

import "github.com/urfave/cli/v2"

var (
	// 调试-查看内存数据
	getCacheCommand = &cli.Command{
		Action:    GetCache,
		Name:      "GetCache",
		Usage:     "GetCache",
		ArgsUsage: "GetCache",
		Flags:     nil,
		Description: `
GetCache`,
	}

	// 调试-查看内存数据
	getAckCommand = &cli.Command{
		Action:    GetAck,
		Name:      "GetAck",
		Usage:     "GetAck",
		ArgsUsage: "GetAck",
		Flags:     nil,
		Description: `
GetAck`,
	}

	// 调试-查看定时任务数据
	getTickerCommand = &cli.Command{
		Action:    GetTicker,
		Name:      "GetTicker",
		Usage:     "GetTicker",
		ArgsUsage: "GetTicker",
		Flags:     nil,
		Description: `
GetTicker`,
	}

	cleanTableCommand = &cli.Command{
		Action:    CleanTable,
		Name:      "CleanTable",
		Usage:     "CleanTable",
		ArgsUsage: "CleanTable",
		Flags:     nil,
		Description: `
GetTicker`,
	}
)

func GetCommand() []*cli.Command {
	return []*cli.Command{
		getCacheCommand,
		getAckCommand,
		getTickerCommand,
		cleanTableCommand,
	}
}
