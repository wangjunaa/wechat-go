package token

import (
	"demo/config"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

// CreateToken 创建token
func CreateToken(id string) (string, error) {
	c := jwt.StandardClaims{
		Issuer:    config.Issuer,
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(config.ExpiresTime)).Unix(),
		Audience:  id,
	}
	log.Println("CreateToken:", c)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		log.Println("tools.token.CreateToken:", err)
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解密token
func ParseToken(tokenString string) (*jwt.Token, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		log.Println("tools.token.ParseToken:", err)
		return nil, false
	}
	return token, true
}

// CheckToken 检查token是否正确
func CheckToken(tokenString string, id string) bool {
	token, ok := ParseToken(tokenString)
	if !ok {
		return false
	}
	if c := token.Claims.(jwt.MapClaims); c["iss"] != config.Issuer || c["aud"] != id {
		return false
	}
	return true
}
