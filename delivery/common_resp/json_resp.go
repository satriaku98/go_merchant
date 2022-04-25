package common_resp

import (
	"github.com/gin-gonic/gin"
)

type commonResp struct {
	g *gin.Context
}

func (cr *commonResp) SuccessResp(httpCode int, succesMessage *SuccessResponse) {
	cr.g.JSON(httpCode, succesMessage)
	cr.g.Abort()
}

func (cr *commonResp) FailedResp(httpCode int, failedMessage *ErrorResponse) {
	cr.g.JSON(httpCode, failedMessage)
	cr.g.Abort()
}

func NewCommonResp(g *gin.Context) *commonResp {
	return &commonResp{
		g,
	}
}
