package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Phone struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	Description string `json:"description"`
	IsFax       bool   `json:"is_fax"`
}

func AddPhoneHandler(c *gin.Context) {
	userID := GetUserIDFromToken(c)

	var newPhoneData Phone
	if err := c.ShouldBind(&newPhoneData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newPhoneData.PhoneNumber == "" || newPhoneData.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	count, err := CountPhonesByPhoneNumber(db, newPhoneData.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check phone duplicate"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already exists"})
		return
	}

	_, err = InsertPhone(db, userID, newPhoneData.PhoneNumber, newPhoneData.Description, newPhoneData.IsFax)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert phone"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Phone added successfully"})
}

func SearchPhoneHandler(c *gin.Context) {
	userID := GetUserIDFromToken(c)
	query := c.Query("q")

	// Search for phones by query
	phones, err := SearchPhonesByQuery(db, userID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search phone"})
		return
	}

	c.JSON(http.StatusOK, phones)
}

func UpdatePhoneHandler(c *gin.Context) {
	userID := GetUserIDFromToken(c)
	var newPhoneData Phone

	if err := c.ShouldBind(&newPhoneData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Print(newPhoneData)
	if newPhoneData.ID == 0 || newPhoneData.PhoneNumber == "" || newPhoneData.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	isOwner, err := CheckPhoneOwnership(db, newPhoneData.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "Phone does not belong to user"})
		return
	}

	err = UpdatePhoneData(db, newPhoneData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Phone updated successfully!!!"})
	}
}

func DeletePhoneHandler(c *gin.Context) {
	userID := GetUserIDFromToken(c)
	phoneID, err := strconv.Atoi(c.Param("phone_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	isOwner, err := CheckPhoneOwnership(db, phoneID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "Phone does not belong to user"})
		return
	}

	err = DeletePhoneData(db, phoneID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Phone deleted successfully!!!"})
	}
}
