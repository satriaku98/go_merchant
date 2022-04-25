package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"go_merchant/model"
	"time"
)

type Token interface {
	CreateToken(dataLogin model.Customer) (string, error)
	VerifAccessToken(tokenString string) (jwt.MapClaims, error)
	GetAppName() string
	CheckTokenAvailable(tokenString string, cashierId interface{}) (bool, error)
	UpdateToken(tokenString string)
}

type TokenConfig struct {
	AplicationName      string
	JwtSignatureKey     string
	JwtSignatureMethod  *jwt.SigningMethodHMAC
	AccessTokenDuration time.Duration
}

type token struct {
	config TokenConfig
	rdb    *sqlx.DB
	ctx    context.Context
}

func (t *token) GetAppName() string {
	return t.config.AplicationName
}

func (t *token) UpdateToken(tokenString string) {
	//t.rdb.Set(t.ctx, "token", tokenString, t.config.AccessTokenDuration)
}

func (t *token) CreateToken(dataLogin model.Customer) (string, error) {
	claims := MyClaims{ // Menyiapkan struct dengan isi yg dibutuhkan
		StandardClaims: jwt.StandardClaims{
			Issuer: t.config.AplicationName,
		},
		CustomerId: dataLogin.CustomerId,
		Name:       dataLogin.Name,
	}
	token := jwt.NewWithClaims(t.config.JwtSignatureMethod, claims) // membuat jwt dengan format method dan claim berupa struct
	tokenString, err := token.SignedString([]byte(t.config.JwtSignatureKey))
	//tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NDkzMTU5NzksInN1YiI6MX0.Eb-zFl9pVL7lmVjJCf74SqUIhfe3VXJ3_uJhTvm7iYc"
	t.UpdateToken(tokenString)
	_, errCreate := t.rdb.Query("update customer set token = $1 where customer_id = $2", tokenString, dataLogin.CustomerId)
	if errCreate != nil {
		return tokenString, errCreate
	}
	return tokenString, err // mendapatkan token
}

func (t *token) CheckTokenAvailable(tokenString string, cashierId interface{}) (bool, error) {
	var data model.Customer
	err := t.rdb.Get(&data, "select token from customer where customer_id = ?", cashierId)
	if err != nil {
		return false, err
	}
	if tokenString != data.Token {
		return false, nil
	}
	return true, nil
}

func (t *token) VerifAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signid Method invalid")
		} else if method != t.config.JwtSignatureMethod {
			return nil, fmt.Errorf("signid Method invalid")
		}
		return []byte(t.config.JwtSignatureKey), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func NewToken(config TokenConfig, ctx context.Context, rdb *sqlx.DB) Token {
	return &token{
		config,
		rdb,
		ctx,
	}
}
