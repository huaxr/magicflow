// Author: huaxr
// Time:   2022/1/14 上午10:21
// Git:    huaxr

package ssrf

import "github.com/go-resty/resty/v2"

type HttpClient struct {
	basicPath  string
	restClient *resty.Client
}

func NewHttpClient(domain string) *HttpClient {
	hc := new(HttpClient)
	hc.basicPath = domain
	hc.restClient = resty.New()
	return hc
}

func (cli *HttpClient) GetRestClient() *resty.Client {
	return cli.restClient
}

func (cli *HttpClient) GetBasePath() string {
	return cli.basicPath
}
