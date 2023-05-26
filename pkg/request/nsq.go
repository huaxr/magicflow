// Author: huaxr
// Time:   2022/1/14 上午10:39
// Git:    huaxr

package request

// nsq auth
type AuthAccount struct {
	Channels    []string `json:"channels"`
	Topic       string   `json:"topic"`
	Permissions []string `json:"permissions"`
}

type NsqAuth struct {
	Identity       string `json:"identity"`
	TTL            int    `json:"ttl"`
	Authorizations []*AuthAccount
}
