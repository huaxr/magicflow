// Author: XinRui Hua
// Time:   2022/4/8 下午2:46
// Git:    huaxr

package middleware

import (
	"github.com/huaxr/magicflow/component/logx"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoveryInterceptor panic时返回Unknown错误吗
func RecoveryInterceptor() grpc_recovery.Option {
	return grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		logx.L().Errorf("grpc err:%v", p)
		return status.Errorf(codes.Unknown, "grpc panic: %v", p)
	})
}
