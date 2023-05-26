package normal

import (
	"net/http"
	"strings"

	"github.com/huaxr/magicflow/component/api/auth"
	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/huaxr/magicflow/pkg/jwtutil"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

var (
	AllowedHosts    = []string{""}
	AllowedLanguage = []string{"en", "zh"}
)

func DebugCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// dynamic using the given origin . when using "*" which will disable cookie by chrome save reasons
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", strings.Join([]string{"content-type", "JWT"}, ","))
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, gin.H{"error_code": 0, "err_msg": nil, "data": "Options Request Success!"})
			c.Abort()
			return
		}
	}
}

func GetHeaderToken(c *gin.Context, head string) (string, bool) {
	h := c.Request.Header.Get(head)
	if len(h) > 0 {
		return strings.TrimPrefix(h, "AppToken "), true
	}
	return "", false
}

var bypasswhitemap = map[string]bool{
	"/callback":                 true,
	"/auth":                     true,
	"/hosts":                    true,
	"/tmp":                      true,
	"/notifycation":             true,
	"/sql_shell":                true,
	"/config/lookups":           true,
	"/config/brokers":           false,
	"/config/auth":              true,
	"/trigger/execute":          true,
	"/trigger/hook":             true,
	"/trigger/worker_exception": true,
	"/trigger/worker_response":  true,
}

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// local env debug bypass all authentication
		if confutil.GetConf().Configuration.Env == "debug" {
			fakeauth := &auth.UserInfo{
				Uid:          1,
				Account:      "guest",
				Name:         "guest",
				Avatar:       "",
				DeptName:     "4s",
				DeptFullName: "4s",
			}
			c.Set("auth", fakeauth)
			c.Next()
			return
		}

		if b, ok := bypasswhitemap[c.Request.URL.Path]; ok && b {
			c.Next()
			return
		}

		// token validate
		//h := c.Request.Header.get("Authorization")
		//t := strings.TrimPrefix(h, "AppToken ")

		coo, _ := c.Cookie("magic")
		token, err := jwtutil.CheckTokenString(coo)
		if coo == "" || err != nil {
			c.JSON(401, "MagicFlow Unauthorized")
			c.Abort()
			return
		}

		infom := token.Claims.(jwt.MapClaims)

		auth := &auth.UserInfo{
			Uid:          cast.ToInt(infom["uid"]),
			Account:      cast.ToString(infom["account"]),
			Name:         cast.ToString(infom["name"]),
			Avatar:       cast.ToString(infom["avatar"]),
			DeptName:     cast.ToString(infom["deptname"]),
			DeptFullName: cast.ToString(infom["deptfullname"]),
			IsAdmin:      cast.ToInt(infom["is_admin"]),
		}

		c.Set("auth", auth)
		c.Next()
	}
}
