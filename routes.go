package main

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.POST("/user/register", RegisterHandler)
	router.POST("/user/auth", AuthHandler)
	router.GET("/user/:name", AuthMiddleware(), GetUserHandler)
	router.POST("/user/phone", AuthMiddleware(), AddPhoneHandler)
	router.GET("/user/phone", AuthMiddleware(), SearchPhoneHandler)
	router.PUT("/user/phone", AuthMiddleware(), UpdatePhoneHandler)
	router.DELETE("/user/phone/:phone_id", AuthMiddleware(), DeletePhoneHandler)
}
