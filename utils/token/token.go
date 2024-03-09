package token

import (
	"demo/dao"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

// CreateToken 创建token
func CreateToken(id string) (string, error) {
	c := jwt.StandardClaims{
		Issuer:    dao.Issuer,
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(dao.ExpiresTime)).Unix(),
		Audience:  id,
	}
	log.Println("CreateToken:", c)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString([]byte(dao.SecretKey))
	if err != nil {
		log.Println("utils.token.CreateToken:", err)
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解密token
func ParseToken(tokenString string) (*jwt.Token, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(dao.SecretKey), nil
	})
	if err != nil {
		log.Println("utils.token.ParseToken:", err)
		return nil, false
	}
	return token, true
}

// CheckToken 检查token是否正确并返回用户id
func CheckToken(tokenString string) string {
	token, ok := ParseToken(tokenString)
	if !ok {
		return ""
	}
	c := token.Claims.(jwt.MapClaims)
	if c["iss"] != dao.Issuer {
		return ""
	}
	return c["aud"].(string)
}
