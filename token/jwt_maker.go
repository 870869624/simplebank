package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secret string = "1111111" //不能被捕获----------------------

const minSecretKeySize = 32

// 使用对成密钥算法
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) { //通过实现Maker接口确保JWTMaker必须实现了token maker 接口
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("密钥长度%d小于最短密钥%d要求", len(secretKey), minSecretKeySize)
	}
	return &JWTMaker{secretKey: secretKey}, nil
}
func (JWT *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	Payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Payload)
	return jwtToken.SignedString([]byte(JWT.secretKey)) //这里为什么不能直接用secret----------------
}

func (JWT *JWTMaker) VerifyToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) { //解析令牌会出现两种错误，过期或者无效
		_, ok := token.Method.(*jwt.SigningMethodHMAC) //类型断言，把token.Method转换为HMAC类型
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(JWT.secretKey), nil //这里为什么不能直接用secret--------------------
	})
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payLoad, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payLoad, nil
}
