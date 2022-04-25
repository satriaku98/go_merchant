package middleware

import (
	"github.com/gin-gonic/gin"
	authenticator2 "go_merchant/auth"
	"go_merchant/delivery/common_resp"
	"net/http"
	"strings"
)

type AuthTokenMiddleware struct {
	acctToken authenticator2.Token
}

func NewAuthTokenMiddleware(configToken authenticator2.Token) *AuthTokenMiddleware {
	return &AuthTokenMiddleware{
		acctToken: configToken,
	}
}

func (a *AuthTokenMiddleware) TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "/passcode") || strings.Contains(c.Request.URL.Path, "/login") {
			c.Next()
		} else {
			h := authHeader{}
			err := c.ShouldBindHeader(&h)
			if err != nil {
				common_resp.NewCommonResp(c).FailedResp(http.StatusUnauthorized, common_resp.FailedMessage(err.Error()))
				return
			}
			tokenString := strings.Replace(h.AuthorizationHeader, "Bearer ", "", -1)
			if tokenString == "" {
				common_resp.NewCommonResp(c).FailedResp(http.StatusUnauthorized, common_resp.FailedMessage("Unautherized"))
				return
			}
			token, errToken := a.acctToken.VerifAccessToken(tokenString)
			if errToken != nil {
				common_resp.NewCommonResp(c).FailedResp(http.StatusUnauthorized, common_resp.FailedMessage(errToken.Error()))
				return
			}
			condition, err := a.acctToken.CheckTokenAvailable(tokenString, token["cashierId"])
			if err != nil {
				common_resp.NewCommonResp(c).FailedResp(http.StatusUnauthorized, common_resp.FailedMessage("Unautherized"))
				return
			}
			if condition {
				a.acctToken.UpdateToken(tokenString)
				c.Next()
			} else {
				common_resp.NewCommonResp(c).FailedResp(http.StatusUnauthorized, common_resp.FailedMessage("Unautherized"))
				return
			}

			c.Next()
		}
	}
}
