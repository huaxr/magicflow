// Author: XinRui Hua
// Time:   2022/1/24 下午4:36
// Git:    huaxr

package ssrf

import (
	"encoding/json"
	"fmt"

	"github.com/huaxr/magicflow/component/logx"

	"time"

	"github.com/huaxr/magicflow/pkg/request"
)

// dynamic expansion and contraction capacity
func (cli *HttpClient) DeleteTopic(broker, topic string) (string, error) {
	var result request.TicketResult
	resp, err := cli.GetRestClient().R().SetBody(map[string]string{
		"topic": topic,
	}).Delete(fmt.Sprintf("http://%v/api/nodes/%v", cli.basicPath, broker))

	if err != nil {
		return "", err
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}

	return result.Ticket, nil
}

type ClientStats struct {
	Node              string        `json:"node"`
	RemoteAddress     string        `json:"remote_address"`
	Version           string        `json:"version"`
	ClientID          string        `json:"client_id"`
	Hostname          string        `json:"hostname"`
	UserAgent         string        `json:"user_agent"`
	ConnectTs         int64         `json:"connect_ts"`
	ConnectedDuration time.Duration `json:"connected"`
	InFlightCount     int           `json:"in_flight_count"`
	ReadyCount        int           `json:"ready_count"`
	FinishCount       int64         `json:"finish_count"`
	RequeueCount      int64         `json:"requeue_count"`
	MessageCount      int64         `json:"message_count"`
	SampleRate        int32         `json:"sample_rate"`
	Deflate           bool          `json:"deflate"`
	Snappy            bool          `json:"snappy"`
	Authed            bool          `json:"authed"`
	AuthIdentity      string        `json:"auth_identity"`
	AuthIdentityURL   string        `json:"auth_identity_url"`

	TLS                           bool   `json:"tls"`
	CipherSuite                   string `json:"tls_cipher_suite"`
	TLSVersion                    string `json:"tls_version"`
	TLSNegotiatedProtocol         string `json:"tls_negotiated_protocol"`
	TLSNegotiatedProtocolIsMutual bool   `json:"tls_negotiated_protocol_is_mutual"`
}

type ChannelStats struct {
	//Node string `json:"node"`
	//Hostname      string          `json:"hostname"`
	TopicName string `json:"topic_name"`
	//ChannelName   string          `json:"channel_name"`
	Depth         int64 `json:"depth"`
	MemoryDepth   int64 `json:"memory_depth"`
	BackendDepth  int64 `json:"backend_depth"`
	InFlightCount int64 `json:"in_flight_count"`
	DeferredCount int64 `json:"deferred_count"`
	RequeueCount  int64 `json:"requeue_count"`
	TimeoutCount  int64 `json:"timeout_count"`
	MessageCount  int64 `json:"message_count"`
	ClientCount   int   `json:"client_count"`
	//Selected      bool            `json:"-"`
	//NodeStats     []*ChannelStats `json:"nodes"`
	//Clients []*ClientStats `json:"clients"`
	Paused bool `json:"paused"`
}

// In-Flight	Current count of messages delivered but not yet finished (FIN), requeued (REQ) or timed out
// Deferred	Current count of messages that were requeued and explicitly deferred which are not yet available for delivery.
// Requeued	Total count of messages that have been added back to the queue due to time outs or explicit requeues.
// Timed Out	Total count of messages that were requeued after not receiving a response from the client before the configured timeout.
func (cli *HttpClient) ReportChannelStats(app string) *ChannelStats {
	var result ChannelStats
	_, err := cli.GetRestClient().R().
		SetResult(&result).
		Get(fmt.Sprintf("http://%v/api/topics/%s/nil", cli.basicPath, app))

	if err != nil {
		logx.L().Errorf("heartbeat err %v", err)
		return nil
	}

	return &result
}
