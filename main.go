package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MiddlewaresCors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	})
}

func HttpPostForm(targetUrl string, parms map[string][]string) (string, error) {
	if targetUrl == "" {
		return "", errors.New("url is empty ")
	}
	if len(parms) == 0 {
		return "", errors.New("parms is empty ")
	}

	resPost, err := http.PostForm(targetUrl, parms)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	defer resPost.Body.Close()

	resBody, _ := ioutil.ReadAll(resPost.Body)
	//fmt.Println(string(resBody))
	return string(resBody), nil
}

type ResValue struct {
	Status  int      `json:"status"`
	Info    string   `json:"info"`
	ResData *ResData `json:"data"`
}

type ResData struct {
	UserID  string `json:"user_id"`
	Account string `json:"account"`
}

func RegUser(c *gin.Context) {
	account := strings.TrimSpace(c.Query("account"))
	password := strings.TrimSpace(c.Query("password"))

	message := "success"
	status := 200

	if account == "" || password == "" {
		message = "Account or password cannot be empty"
		status = 400
	} else if len(account) < 6 || len(account) > 60 || len(password) < 6 || len(password) > 60 {
		message = "Account or password needs 6~60 characters"
		status = 400
	} else {
		parms := url.Values{"account": {account}, "password": {password}}
		targetUrl := "http://localhost:8080/user/reg"

		resPost, err := HttpPostForm(targetUrl, parms)
		if err != nil {
			message = "Service error" //err.Error()
			status = 400
		} else {
			resValue := ResValue{}
			err := json.Unmarshal([]byte(resPost), &resValue)
			if err != nil {
				message = "Service error" //err.Error()
				status = 400
			} else {
				if resValue.Status == 0 {
					message = resValue.Info
					status = 400
				} else {
					message = "success"
					status = 200
				}
			}
		}
	}
	c.JSON(200, gin.H{
		"status":  status,
		"message": message,
	})
}

func Login(c *gin.Context) {
	account := strings.TrimSpace(c.Query("account"))
	password := strings.TrimSpace(c.Query("password"))

	message := ""
	status := 200

	if account == "" || password == "" {
		message = "Account or password cannot be empty"
		status = 400
	} else if len(account) < 6 || len(account) > 60 || len(password) < 6 || len(password) > 60 {
		message = "Account or password needs 6~60 characters"
		status = 400
	} else {
		parms := url.Values{"account": {account}, "password": {password}}
		targetUrl := "http://localhost:8080/user/login"

		resPost, err := HttpPostForm(targetUrl, parms)
		if err != nil {
			message = "Service error" //err.Error()
			status = 400
		} else {
			resValue := ResValue{}
			err := json.Unmarshal([]byte(resPost), &resValue)
			if err != nil {
				message = "Service error" //err.Error()
				status = 400
			} else {
				if resValue.Status == 0 {
					message = resValue.Info
					status = 400
				} else {
					message = "Welcome " + account
					status = 200
				}
			}
		}
	}

	c.JSON(200, gin.H{
		"status":  status,
		"message": message,
	})
}

func Index(c *gin.Context) {
	c.String(200, "Service running")
}

func main() {
	router := gin.Default()
	router.Use(MiddlewaresCors())
	router.GET("/reguser", RegUser)
	router.GET("/login", Login)
	router.GET("/", Index)
	router.Run(":8080")
}
