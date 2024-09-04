package api

import (
	"anime-go/internal/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func loginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	u := models.User{}
	models.DB.Where("username = ?", username).First(&u)
	if u.ID == 0 || u.Password != password {
		c.JSON(401, gin.H{"status": "error"})
		return
	}

	var mySigningKey = []byte("MySecretKey")

	claims := MyCustomClaims{
		u.Username, // 自定义字段
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // 过期时间1小时后
			Issuer:    "anime-go",                                        // 签发者
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名token
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		Expires:  time.Now().Add(180 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // 在生产环境中应设置为 true，并使用 HTTPS
	}
	http.SetCookie(c.Writer, cookie)
	c.JSON(200, gin.H{"token": tokenString})
}
