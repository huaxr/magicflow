package jwtutil

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
)

func TestX(t *testing.T) {
	x, err := GenTokenString(map[string]interface{}{"user_name": "huaxinrui"})
	t.Log(x, err)

	res, err := CheckTokenString(x)
	t.Log(err, "erris:", res.Claims.(jwt.MapClaims)["user_name"])

}
