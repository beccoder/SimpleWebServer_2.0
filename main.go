package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	err := ConnectDatabase()
	if err != nil {
		log.Fatal("Cannot connect database:", err)
	}
	defer db.Close()

	router := gin.Default()

	RegisterRoutes(router)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
