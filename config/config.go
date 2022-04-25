package config

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/viper"
	"go_merchant/auth"
	"time"
)

type Config struct {
	ConfigToken  auth.TokenConfig
	ConfigServer *ConfigServer
	*ConfigDatabase
}

type ConfigDatabase struct {
	dbConn string
	//mysqlConn   string
	configRedis *ConfigRedis
}

type ConfigServer struct {
	Url  string
	Port string
}

type ConfigRedis struct {
	Address  string
	Password string
	Db       int
}

func newTokenConfig() auth.TokenConfig {
	//duration, _ := strconv.Atoi(GetConfigValue("JWTDURATION"))
	return auth.TokenConfig{
		AplicationName:      GetConfigValue("APLICATIONNAME"),
		JwtSignatureKey:     GetConfigValue("JWTKEY"),
		JwtSignatureMethod:  jwt.SigningMethodHS256,
		AccessTokenDuration: 60 * time.Minute,
	}
}

func newServerConfig() *ConfigServer {
	return &ConfigServer{
		GetConfigValue("SERVERURL"),
		GetConfigValue("SERVERPORT"),
	}
}

//func (c *ConfigDatabase) MysqlConn() string {
//	return c.mysqlConn
//}

func (c *ConfigDatabase) PostgreConn() string {
	return c.dbConn
}

func (c *ConfigDatabase) RedisConfig() *ConfigRedis {
	return c.configRedis
}

func ReadConfigFile(configName string) {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the Config file
	if err != nil {             // Handle errors reading the Config file
		panic(fmt.Errorf("Fatal error Config file: %w \n", err))
	}
}

func GetConfigValue(configName string) string {
	ReadConfigFile("Config")
	return viper.GetString(configName)
}

func newPostgreConn() string {
	dbName := GetConfigValue("DBNAME")
	dbHost := GetConfigValue("DBHOST")
	dbUsername := GetConfigValue("DBUSERNAME")
	dbPassword := GetConfigValue("DBPASSWORD")
	dbPort := GetConfigValue("DBPORT")
	urlDb := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUsername, dbPassword, dbHost, dbPort, dbName)
	fmt.Println(urlDb)
	// urlExample := "postgres://username:password@localhost:5432/database_name"

	return urlDb
}

//func newMysqlConn() string {
//	dbUser := GetConfigValue("MYSQLUSER")
//	dbPass := GetConfigValue("MYSQLPASS")
//	dbUrl := GetConfigValue("MYSQLURL")
//	dbPort := GetConfigValue("MYSQLPORT")
//	dbName := GetConfigValue("MYSQLDBNAME")
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbUrl, dbPort, dbName)
//	return dsn
//}

func newRedisConn() *ConfigRedis {
	return &ConfigRedis{
		Address:  fmt.Sprintf("%s:%s", GetConfigValue("REDISURL"), GetConfigValue("REDISPORT")),
		Password: GetConfigValue("REDISPASSWORD"),
		Db:       0,
	}
}

func NewConfig() *Config {
	return &Config{
		ConfigToken:  newTokenConfig(),
		ConfigServer: newServerConfig(),
		ConfigDatabase: &ConfigDatabase{
			newPostgreConn(),
			//newMysqlConn(),
			newRedisConn(),
		},
	}
}
