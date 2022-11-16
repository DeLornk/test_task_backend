package main

import (
	"avito_task_2022/src"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

// @title Balance Manipulate API
// @version 1.0
// @description API Server for manipulating with user's balance

// @host localhost 8080
func main() {
	err := src.InitDB()
	if err != nil {
		log.Fatal(fmt.Println("Error: ", err))
	}

	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/balance/find", src.GetBalance)
	router.POST("/balance/change", src.ChangeBalance)
	router.POST("/balance/reserve", src.ReserveMoney)
	router.POST("/balance/unreserve", src.UnReserveMoney)
	router.POST("/transaction", src.TransactionBetween)
	router.POST("/accept_revenue", src.AcceptRevenue)
	router.GET("/month_report", src.GetMonthReport)
	router.GET("/operation_report", src.ReviewOperations)

	router.Run("localhost:8080")

	fmt.Println("Успешно запустились")
}
