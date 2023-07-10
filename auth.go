package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JWTClaims struct {
	Login  string `json:"login"`
	UserID int    `json:"user_id"`
	jwt.StandardClaims
}

var signingKey = []byte("secret")

func AuthHandler(c *gin.Context) {
	var loginData LoginData
	if err := c.ShouldBind(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if loginData.Login == "" || loginData.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	user, err := GetUserByLogin(db, loginData.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	err = ComparePasswords(user.Password, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
		return
	}

	tokenString, err := GenerateToken(loginData.Login, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	cookie := &http.Cookie{
		Name:  "SESSTOKEN",
		Value: tokenString,
		Path:  "/",
	}
	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{
	"message": "Authentication successful"})
}
