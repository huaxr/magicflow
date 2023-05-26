// Author: huaxr
// Time:   2022/1/14 上午10:37
// Git:    huaxr

package triggersdk

import (
	"encoding/json"
	"fmt"
	"github.com/huaxr/magicflow/pkg/request"
)

func (cli *HttpClient) GetLookUps() (string, error) {
	var result request.LookupAddsRes
	resp, err := cli.GetRestClient().R().
		Get(fmt.Sprintf("%s%s", cli.GetBasePath(), "/config/lookups"))

	if err != nil {
		return "", err
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}

	return result.Data, nil
}
