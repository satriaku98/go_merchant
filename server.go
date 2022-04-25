package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_merchant/auth"
	"go_merchant/config"
	"go_merchant/delivery/api"
	"go_merchant/delivery/middleware"
	"go_merchant/manager"
)

type AppServer interface {
	Run()
}

type serverConfig struct {
	gin            *gin.Engine
	Name           string
	Port           string
	InfraManager   manager.InfraManager
	RepoManager    manager.RepoManager
	UseCaseManager manager.UseCaseManager
	Config         *config.Config
	Middleware     *middleware.AuthTokenMiddleware
	ConfigToken    auth.Token
}

func (s *serverConfig) initHeader() {
	s.gin.Use(s.Middleware.TokenAuthMiddleware())
	s.routeGroupApi()
}

func (s *serverConfig) routeGroupApi() {

	apiCashier := s.gin.Group("")
	api.NewLoginApi(apiCashier, s.UseCaseManager.LoginUseCase(), s.ConfigToken)

}

func (s *serverConfig) Run() {
	s.initHeader()
	s.gin.Run(fmt.Sprintf("%s:%s", s.Name, s.Port))
}

func Server() AppServer {
	ginStart := gin.Default()
	config := config.NewConfig()
	infra := manager.NewInfraManager(config.ConfigDatabase)
	repo := manager.NewRepoManager(infra.MysqlConn())
	usecase := manager.NewUseCaseManager(repo)
	configToken := infra.ConfigToken(config.ConfigToken)
	middleware := middleware.NewAuthTokenMiddleware(configToken)
	return &serverConfig{
		ginStart,
		config.ConfigServer.Url,
		config.ConfigServer.Port,
		infra,
		repo,
		usecase,
		config,
		middleware,
		configToken,
	}
}
