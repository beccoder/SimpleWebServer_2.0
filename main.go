package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/beccoder/SimpleWebServer_2.0/models"
	"github.com/gin-gonic/gin"
)

func main() {
	err := models.ConnectDatabase()
	checkErr(err)
	defer models.DB.Close()

	server := gin.Default()

	server.POST("/user/register", registerUser)
	server.POST("/user/auth", authenticateUser)
	server.GET("/user/:name", authMiddleware(), getUserByName)
	server.POST("/user/phone", authMiddleware(), addPhoneNumber)
	server.GET("/user/phone", authMiddleware(), getPhoneNumber)
	server.PUT("/user/phone", authMiddleware(), updatePhoneNumber)
	server.DELETE("/user/phone/:phone_id", authMiddleware(), deletePhoneNumber)

	server.Run(":8080")
}

func registerUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBind(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Print(newUser)
	if newUser.Login == "" || newUser.Password == "" || newUser.Name == "" || newUser.Age == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	success, err := models.RegisterUser(newUser)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func authenticateUser(c *gin.Context) {
	var loginData models.LoginData
	if err := c.ShouldBind(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if loginData.Login == "" || loginData.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	success, err := models.AuthenticateUser(loginData, c)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Authentication successful"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func getUserByName(c *gin.Context) {
	name := c.Param("name")
	user, err := models.GetUserHandler(name)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func addPhoneNumber(c *gin.Context) {
	userId := getUserIdFromToken(c)
	var newPhoneData models.Phone
	if err := c.ShouldBind(&newPhoneData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newPhoneData.PhoneNumber == "" || newPhoneData.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}
	count, err := models.CountPhoneNumber(newPhoneData.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already exists"})
		return
	}
	newPhoneData.UserID = userId
	success, err := models.AddPhoneNumber(newPhoneData)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Phone added successfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func getPhoneNumber(c *gin.Context) {
	userId := getUserIdFromToken(c)
	query := c.Query("q")

	phones, err := models.SearchPhoneHandler(userId, query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, phones)
	}
}

func updatePhoneNumber(c *gin.Context) {
	var newPhoneData models.Phone
	userId := getUserIdFromToken(c)

	if err := c.ShouldBind(&newPhoneData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newPhoneData.ID == 0 || newPhoneData.PhoneNumber == "" || newPhoneData.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	isOwner, err := models.CheckPhoneOwnership(userId, newPhoneData.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "Phone does not belong to user"})
		return
	}

	success, err := models.UpdatePhoneData(newPhoneData)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Phone updated successfully!!!"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func deletePhoneNumber(c *gin.Context) {
	userId := getUserIdFromToken(c)
	phoneId, err := strconv.Atoi(c.Param("phone_id"))
	checkErr(err)

	isOwner, err := models.CheckPhoneOwnership(userId, phoneId)
	checkErr(err)

	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "Phone does not belong to user"})
		return
	}

	success, err := models.DeletePhoneData(phoneId)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Phone number deleted successfully"})
		return
	} else {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("SESSTOKEN")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		tokenString := cookie.Value

		user_id, err := models.CheckJWTTokens(tokenString)

		if err != nil || user_id < 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		} else {
			c.Set("user_id", user_id)
			c.Next()
		}
	}
}

func getUserIdFromToken(c *gin.Context) int {
	userId, _ := c.Get("user_id")
	return userId.(int)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
