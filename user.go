package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id" form:"id"`
	Login    string `json:"login" form:"login"`
	Password string `json:"password" form:"password"`
	Name     string `json:"name" form:"name"`
	Age      int    `json:"age" form:"age"`
}

func RegisterHandler(c *gin.Context) {
	// Our struct can be assigned to form data as well.
	// Getting data by binding is the best practise, but we can use PostForm also

	var newUser User
	// Apply PostForm method

	// newUser.Login = c.PostForm("login")
	// newUser.Password = c.PostForm("password")
	// newUser.Name = c.PostForm("name")
	// age, err := strconv.Atoi(c.PostForm("age"))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// newUser.Age = age

	if err := c.ShouldBind(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newUser.Login == "" || newUser.Password == "" || newUser.Name == "" || newUser.Age == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	in, err := CheckLoginExistance(db, newUser.Login)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if in {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this login already exists"})
		return
	}

	_, err = InsertUser(db, newUser.Login, string(hashedPassword), newUser.Name, strconv.Itoa(newUser.Age))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func GetUserHandler(c *gin.Context) {
	name := c.Param("name")

	user, err := GetUserByName(db, name)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
