package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")

		if err != nil || !verificationHandler(cookie) {
			c.JSON(401, gin.H{"status": "unauthorized"})
			c.Abort() // 阻止请求继续被处理
			return
		}

		c.Next() // 继续处理
	}
}

type JWT struct {
	Token string `json:"token"`
}

func verificationHandler(t string) bool {

	var mySigningKey = []byte("MySecretKey")

	// 解析并验证token
	token, err := jwt.ParseWithClaims(t, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证使用的签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})

	if err != nil {
		// fmt.Println("Error parsing token:", err)
		return false
	}

	if _, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		// c.JSON(200, gin.H{"status": claims.Username})
		return true
	} else {
		// fmt.Println("Invalid token")
		return false
	}
}
