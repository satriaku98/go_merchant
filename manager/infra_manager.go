package manager

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"go_merchant/auth"
	"go_merchant/config"
	"os"
)

type InfraManager interface {
	MysqlConn() *sqlx.DB
	RedisConn() (context.Context, *redis.Client)
	ConfigToken(tokenConfig auth.TokenConfig) auth.Token
}

type infraManager struct {
	mysqlConn *sqlx.DB
	redisConn *redis.Client
	ctx       context.Context
}

func (i *infraManager) MysqlConn() *sqlx.DB {
	return i.mysqlConn
}

func (i *infraManager) RedisConn() (context.Context, *redis.Client) {
	return i.ctx, i.redisConn
}

func (i *infraManager) ConfigToken(tokenConfig auth.TokenConfig) auth.Token {
	return auth.NewToken(tokenConfig, i.ctx, i.mysqlConn)
}

func NewInfraManager(configDatabase *config.ConfigDatabase) InfraManager {
	urlMysql := configDatabase.PostgreConn()
	redisConfig := configDatabase.RedisConfig()
	conn, err := sqlx.Connect("pgx", urlMysql)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
	})
	ctx := context.Background()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return &infraManager{
		mysqlConn: conn,
		redisConn: rdb,
		ctx:       ctx,
	}
}
