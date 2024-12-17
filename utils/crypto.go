package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/bilibili/HCP/app/interface/v1/configs"
	"github.com/dgrijalva/jwt-go"
)

// CustomClaims 自定义声明结构体
type CustomClaims struct {
	UserId int
	jwt.StandardClaims
}

// SHA256 加密
func SHA256(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	res := hex.EncodeToString(hash.Sum(nil))
	return res
}

// GenerateJWT 生成jwt token
func GenerateJWT(userID int, userName string) (string, error) {
	customClaims := &CustomClaims{
		UserId: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(configs.Conf.LoginConf.ExpireSecond) * time.Second).Unix(),
			Issuer:    userName,
		},
	}
	// 采用HMAC SHA256加密算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	tokenString, err := token.SignedString([]byte(SHA256(configs.Conf.LoginConf.SecretKey)))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析jwt token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SHA256(configs.Conf.LoginConf.SecretKey)), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

// GenerateSalt 生成指定字节长度的盐值
func GenerateSalt(bytes int) (string, error) {
	// 创建一个字节切片用于存储盐值
	b := make([]byte, bytes)
	_, err := io.ReadFull(rand.Reader, b) // 使用加密安全的随机数据填充切片
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil // 返回盐值的十六进制字符串形式
}
