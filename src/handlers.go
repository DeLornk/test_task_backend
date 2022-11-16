package src

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// @Summary GetBalance
// @Tags users
// @Description get user's balance
// @ID get-balance
// @Accept json
// @Produce json
// @Param input body User true "Необходимо ввести параметр User.ID"
// @Router /balance/find [get]
func GetBalance(c *gin.Context) {
	var user User

	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"done": false, "message": "некорректно поданы данные"})
		return
	}

	user, err := FindUser(user.ID)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"done": false, "message": "баланс пользователя не найден"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, user)
		return
	}

}

// @Summary ChangeBalance
// @Tags users
// @Change user balance
// @ID change-balance
// @Accept json
// @Produce json
// @Param input body User true "Необходимо указать пользователя и изменение его баланса"
// @Router /balance/change [post]
func ChangeBalance(c *gin.Context) {
	var userChanges User
	if err := c.BindJSON(&userChanges); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"done": false, "message": "некорретный тип данных"})
		return
	}

	user, _ := FindUser(userChanges.ID)

	newBalance := user.Balance + userChanges.Balance

	log.Println("Новый баланс: ", newBalance)

	if newBalance < 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"done": false, "message": "итоговый баланс отрицательный"})
		return
	} else {
		done, _ := UpdateBalance(newBalance, userChanges.ID)

		c.IndentedJSON(http.StatusOK, gin.H{"done": done})
		return
	}
}

// @Summary ReserveMoney
// @Tags orders
// @Description Reserve money of user
// @ID reserve-money
// @Accept json
// @Produce json
// @Param input body Order true "Необходимо ввести параметры заказа"
// @Router /balance/reserve [post]
func ReserveMoney(c *gin.Context) {
	var order Order
	if err := c.BindJSON(&order); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "некорретный тип данных"})
		return
	}

	httpStatus, _ := CreateOrder(order.ID, order.UserId, order.ProductId, order.Cost)
	switch httpStatus {
	case 201:
		c.IndentedJSON(201, gin.H{"done": true, "message": "успешная резервация денег"})
		return
	case 400:
		c.IndentedJSON(400, gin.H{"done": false, "message": "резервация денег не прошла: указан несуществующий товар или указан id уже существующего заказа"})
		return
	case 404:
		c.IndentedJSON(404, gin.H{"done": false, "message": "резервация денег не прошла: пользователь с данным id не найден"})
		return
	case 406:
		c.IndentedJSON(406, gin.H{"done": false, "message": "резервация денег не прошла: недопустимо"})
		return
	case 500:
		c.IndentedJSON(500, gin.H{"done": false, "message": "резервация денег не прошла: ошибка на сервере"})
		return
	}
}

// @Summary UnReserveMoney
// @Tags orders
// @Description Unreserve money of user
// @ID get-balance
// @Accept json
// @Produce json
// @Param input body Order true "Необходимо ввести параметр order.ID"
// @Router /balance/unreserve [post]
func UnReserveMoney(c *gin.Context) {
	var order Order
	if err := c.BindJSON(&order); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "некорретный тип данных"})
		return
	}

	httpStatus, _ := UnReserve(order.ID)
	switch httpStatus {
	case 200:
		c.IndentedJSON(200, gin.H{"done": true, "message": "успешная разрезервация денег"})
		return
	case 404:
		c.IndentedJSON(404, gin.H{"done": false, "message": "разрезервация денег не прошла: указан id несуществующего заказа"})
		return
	case 500:
		c.IndentedJSON(500, gin.H{"done": false, "message": "разрезервация денег не прошла: ошибка на сервере"})
		return
	}
}

// @Summary AcceptRevenue
// @Tags reports
// @Description Accept revenue and add information to repost
// @ID accept-revenue
// @Accept json
// @Produce json
// @Param input body Order true "Необходимо ввести параметры заказа"
// @Router /accept_revenue [post]
func AcceptRevenue(c *gin.Context) {
	var order Order
	if err := c.BindJSON(&order); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "некорретный тип данных"})
		return
	}

	httpStatus, _ := OrderToReport(order.UserId, order.ProductId, order.ID, order.Cost)
	switch httpStatus {
	case 200:
		c.IndentedJSON(200, gin.H{"done": true, "message": "успешное признание выручки"})
		return
	case 404:
		c.IndentedJSON(404, gin.H{"done": false, "message": "признание выручки не прошло: один (или оба) пользователя не найдены"})
		return
	case 406:
		c.IndentedJSON(406, gin.H{"done": false, "message": "признание выручки не прошло: недостаточно средств"})
		return
	case 500:
		c.IndentedJSON(500, gin.H{"done": false, "message": "признание выручки не прошло: ошибка на сервере"})
		return
	}
}

// @Summary TransactionBetween
// @Tags reports
// @Description transaction between two users
// @ID transaction-between
// @Accept json
// @Produce json
// @Param input body Transaction true "Необходимо ввести параметры транзакции"
// @Router /transaction [post]
func TransactionBetween(c *gin.Context) {
	var transaction Transaction
	if err := c.BindJSON(&transaction); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "некорретный тип данных"})
		return
	}

	httpStatus, _ := FromUserToUser(transaction.Cost, transaction.FromID, transaction.TOID)
	switch httpStatus {
	case 200:
		c.IndentedJSON(200, gin.H{"done": true, "message": "успешная транзакция"})
		return
	case 404:
		c.IndentedJSON(404, gin.H{"done": false, "message": "транзакция не прошла: один (или оба) пользователя не найдены"})
		return
	case 406:
		c.IndentedJSON(406, gin.H{"done": false, "message": "транзакция не прошла: недостаточно средств"})
		return
	case 500:
		c.IndentedJSON(500, gin.H{"done": false, "message": "транзакция не прошла: ошибка на сервере"})
		return
	}
}

// @Summary GetMonthReport
// @Tags reports
// @Description get month report
// @ID get-mr
// @Accept json
// @Produce json
// @Param input body MonthYear true "Необходимо ввести год и месяц"
// @Router /month_report [get]
func GetMonthReport(c *gin.Context) {
	var monthYear MonthYear

	if err := c.BindJSON(&monthYear); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "некорретный тип данных"})
		return
	}

	rprts, err := MonthReport(monthYear.Month, monthYear.Year)
	if err != nil {
		log.Println("MonthReport error: ", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "некорретный тип данных"})
		return
	}
	name, _ := CreateCSV(rprts)

	filepath := fmt.Sprint("./", name)
	c.FileAttachment(filepath, name)

}

// @Summary ReviewOperations
// @Tags reports
// @Description get review user operations
// @ID review-operations
// @Accept json
// @Produce json
// @Param input body ConfigureOperationReview true "Необходимо ввести параметр User.ID"
// @Router /operation_report [get]
func ReviewOperations(c *gin.Context) {
	var configOper ConfigureOperationReview

	if err := c.BindJSON(&configOper); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "некорретный тип данных"})
		return
	}
	log.Println(configOper)

	oprtns, err := ReportUserOperations(configOper.Id, configOper.SortByDate, configOper.SortBySum, configOper.Descending)
	if err != nil {
		log.Println("ReportUserOperations error: ", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "некорретный тип данных"})
		return
	}

	log.Println("ReportUserOperations done")
	c.IndentedJSON(http.StatusOK, oprtns)
	return
}
