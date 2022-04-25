package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	//"go_merchant/delivery/logger"
	"time"
)

type authHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     string
	JwtSigningMethod    *jwt.SigningMethodHMAC
	AccessTokenLifeTime time.Duration
}

type token struct {
	Config TokenConfig
}

func AuthUser(requiredToken string, sqlxdb *sqlx.DB) bool {
	var countAuthUser int
	err := sqlxdb.Get(&countAuthUser, "select count(*) from customer where token = $1", requiredToken)
	if err != nil {
		//logger.SendLogToDiscord("Auth User", err)
		panic("database not connected")
	}
	if countAuthUser == 0 {
		//logger.SendLogToDiscord("Token Auth Not Found", err)
		return true
	}
	return false
}
