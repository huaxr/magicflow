// Author: XinRui Hua
// Time:   2022/4/18 上午11:38
// Git:    huaxr

package tag

type TagKey int

const (
	Trigger TagKey = iota
	WorkerException
	ServerException
	Complied

	PlaybookPut
	AppPut
	AppSwitch
)

func (m TagKey) String() string {
	switch m {
	case Trigger:
		return "Trigger"
	case Complied:
		return "Complied"
	case WorkerException:
		return "WorkerException"
	case ServerException:
		return "ServerException"
	case PlaybookPut:
		return "PlaybookPut"
	case AppPut:
		return "AppPut"

	case AppSwitch:
		return "AppSwitch"
	default:
		return "unknown"
	}
}
