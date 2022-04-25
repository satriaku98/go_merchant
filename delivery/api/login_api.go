package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_merchant/auth"
	"go_merchant/delivery/apprequest"
	"go_merchant/delivery/common_resp"
	"go_merchant/usecase"
	"net/http"
	"strconv"
)

type loginApi struct {
	usecase         usecase.LoginUseCase
	usecaseTransfer usecase.TransferUseCase
	configToken     auth.Token
}

func (api *loginApi) LoginCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataLogin apprequest.LoginRequest
		err := c.ShouldBindJSON(&dataLogin)
		if err != nil {
			common_resp.NewCommonResp(c).FailedResp(http.StatusInternalServerError, common_resp.FailedMessage(err.Error()))
			return
		}
		data, isAvailable, err := api.usecase.LoginCustomer(dataLogin.Username, dataLogin.Password)
		if err != nil {
			common_resp.NewCommonResp(c).FailedResp(http.StatusInternalServerError, common_resp.FailedMessage(err.Error()))
			return
		}
		if isAvailable == 0 {
			common_resp.NewCommonResp(c).FailedResp(http.StatusUnauthorized, common_resp.FailedMessage("not register"))
			return
		}
		tokenString, errToken := api.configToken.CreateToken(data)
		fmt.Println(data, tokenString, errToken)
		if errToken != nil {
			common_resp.NewCommonResp(c).FailedResp(http.StatusInternalServerError, common_resp.FailedMessage("Token Failed"))
			return
		}
		common_resp.NewCommonResp(c).SuccessResp(http.StatusOK, common_resp.SuccessMessage("login admin", gin.H{
			"token": tokenString,
		}))
	}
}

func (api *loginApi) LogoutCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataLogin apprequest.LoginRequest
		err := c.ShouldBindJSON(&dataLogin)
		fmt.Println(dataLogin)
		if errBind := c.ShouldBindJSON(&dataLogin); errBind != nil {
			common_resp.NewCommonResp(c).FailedResp(http.StatusInternalServerError, common_resp.FailedMessage(errBind.Error()))
			return
		}
		isAvailable, err := api.usecase.LogoutCustomer(dataLogin.Username, dataLogin.Password)
		if err != nil {
			common_resp.NewCommonResp(c).FailedResp(http.StatusInternalServerError, common_resp.FailedMessage(err.Error()))
			return
		}
		if isAvailable == 0 {
			common_resp.NewCommonResp(c).FailedResp(http.StatusUnauthorized, common_resp.FailedMessage("not register"))
			return
		}
		common_resp.NewCommonResp(c).SuccessResp(http.StatusOK, common_resp.SuccessMessage("Success", ""))
	}
}
func (a *loginApi) Transfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var authHeader apprequest.CustomerRequest
		var custReq apprequest.CustomerUpdateRequest
		senderAccountNumber, _ := strconv.Atoi(c.Param("user"))
		err := a.ParseRequestBody(c, &custReq)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = c.ShouldBindHeader(&authHeader)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = a.usecaseTransfer.Transfer(senderAccountNumber, custReq.AmountTransfer)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":       "TRANSFER SUCCES",
			"Transfer From": senderAccountNumber,
			"Transfer to":   custReq.ReceiverAccountNumber,
			"Amount":        custReq.AmountTransfer,
		})
	}
}
func NewLoginApi(routeGroup *gin.RouterGroup, usecase usecase.LoginUseCase, configToken auth.Token) {
	api := &loginApi{
		usecase,
		configToken,
	}
	routeGroup.POST("/login", api.LoginCustomer())
	routeGroup.POST("/logout", api.LogoutCustomer())
}
