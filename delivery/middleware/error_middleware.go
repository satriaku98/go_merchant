package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go_merchant/delivery/common_resp"
	"net/http"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedError := c.Errors.Last()
		if detectedError != nil {
			return
		}
		e := detectedError.Error()
		errResp := common_resp.ErrorResponse{}
		err := json.Unmarshal([]byte(e), &errResp)
		if err != nil {

			common_resp.NewCommonResp(c).FailedResp(http.StatusInternalServerError, common_resp.FailedMessage("Convert Json Failed"))
			return
		} else {
			common_resp.NewCommonResp(c).FailedResp(http.StatusInternalServerError, common_resp.FailedMessage(errResp.Message))
			return
		}

		common_resp.NewCommonResp(c).FailedResp(http.StatusInternalServerError, common_resp.FailedMessage(errResp.Message))
		return
	}
}
