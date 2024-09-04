package qbitorrent

import (
	"anime-go/pkg/logger"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Auth struct {
	Sid    string
	Expiry time.Time
}

var AuthInfo Auth

const User = "jlak"
const Pass = "A13eb3fbb."
const Url = "http://192.168.2.9:8282"

func (auth *Auth) updateSid() error {
	sid, err := login(User, Pass)
	if err != nil {
		return errors.New("登陆错误")
	}
	auth.Sid = sid
	auth.Expiry = time.Now().Add(1 * time.Hour)
	return nil
}

func (auth *Auth) getSid() string {
	if time.Now().After(auth.Expiry) {
		err := auth.updateSid()
		if err != nil {
			logger.Log(err.Error())
		}
	}
	return auth.Sid
}

func login(username, password string) (string, error) {
	baseURL := Url // Replace with your base URL

	// Encode form data
	formData := url.Values{}
	formData.Set("username", username)
	formData.Set("password", password)
	urlEncodedData := formData.Encode()

	// Create HTTP request
	req, err := http.NewRequest("POST", baseURL+"/api/v2/auth/login", strings.NewReader(urlEncodedData))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response cookies
	cookies := resp.Cookies()

	// Function to get cookie value
	getCookieValue := func(cookies []*http.Cookie, name string) string {
		for _, cookie := range cookies {
			if cookie.Name == name {
				return cookie.Value
			}
		}
		return ""
	}

	// Get SID from cookies
	sid := getCookieValue(cookies, "SID")

	if sid == "" {
		return "", fmt.Errorf("SID not found in the response cookies")
	}

	// Log SID for debugging
	// fmt.Println("SID:", sid)

	return sid, nil
}
