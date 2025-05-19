package middelware

import (
	"fmt"
	"gin_gorm_oj/utils"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserClaims struct {
	Identity             string `json:"identity"`
	Name                 string `json:"name"`
	IsAdmin              int    `json:"is_admin"`
	jwt.RegisteredClaims        // 替换 StandardClaims
}

var myKey = []byte(utils.JwtKey)

func GenerateToken(identity, name string, isAdmin int) (string, error) {
	UserClaim := &UserClaims{
		Identity: identity,
		Name:     name,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gin-gorm-oj", // 添加必要字段
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	return token.SignedString(myKey)
}
func AnalyseToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return myKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
