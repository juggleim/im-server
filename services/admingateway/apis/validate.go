package apis

import (
	"fmt"
	"im-server/commons/tools"
	"im-server/services/admingateway/services"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	Header_RequestId     string = "request-id"
	Header_Authorization string = "Authorization"
)

func Validate(ctx *gin.Context) {
	session := fmt.Sprintf("admin_%s", tools.GenerateUUIDShort11())
	ctx.Header(Header_RequestId, session)
	ctx.Set(services.CtxKey_Session, session)

	urlPath := ctx.Request.URL.Path
	if strings.HasSuffix(urlPath, "/login") || strings.HasSuffix(urlPath, "/apps/create") {
		return
	}
	authStr := ctx.Request.Header.Get(Header_Authorization)
	account, err := validateAuthorization(authStr)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_AuthFail,
			Msg:  "auth failed",
		})
		ctx.Abort()
		return
	}
	//check account
	code := services.CheckAccountState(account)
	if code != services.AdminErrorCode_Success {
		ctx.JSON(http.StatusUnauthorized, &services.ApiErrorMsg{
			Code: code,
			Msg:  "auth failed",
		})
		ctx.Abort()
		return
	}
	ctx.Set(services.CtxKey_Account, account)
}

func GetLoginedAccount(ctx *gin.Context) string {
	if account, ok := ctx.Value(services.CtxKey_Account).(string); ok {
		return account
	}
	return ""
}

var jwtkey = []byte("jug9le1m")

type Claims struct {
	Account string
	jwt.RegisteredClaims
}

func generateAuthorization(account string) (string, error) {
	expireTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Account: account,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: expireTime,
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
			Issuer:  "aabbcc",
			Subject: "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthTest(auth string) (string, error) {
	return validateAuthorization(auth)
}

func validateAuthorization(authorization string) (string, error) {
	token, claims, err := parseToken(authorization)
	if err != nil || !token.Valid {
		return "", fmt.Errorf("auth fail")
	}
	return claims.Account, nil
}

func parseToken(tokenString string) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtkey, nil
	})
	return token, Claims, err
}
