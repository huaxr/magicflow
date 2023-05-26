// Author: XinRui Hua
// Time:   2022/3/24 下午2:47
// Git:    huaxr

package transport

import (
	"context"
	"encoding/json"

	"github.com/huaxr/magicflow/component/api/auth"
	"github.com/huaxr/magicflow/pkg/toolutil"
	"google.golang.org/grpc/metadata"
)

func CtxToGRpcCtx(ctx context.Context) context.Context {
	md := metadata.Pairs()
	if ctx.Value("auth") != nil {
		auth := ctx.Value("auth").(*auth.UserInfo)
		b, _ := json.Marshal(auth)
		md.Set("auth", toolutil.Bytes2string(b))
	}
	return metadata.NewOutgoingContext(ctx, md)
}
