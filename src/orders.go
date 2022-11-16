package src

import "log"

func CreateOrder(orderID, userID, productID, cost int) (httpStatus int, err error) {
	user, err := FindUser(userID)
	if err != nil {
		log.Println("CreateOrder: ", err)
		return 404, err
	}

	newBalance := user.Balance - cost

	if newBalance < 0 {
		log.Println("CreateOrder error: недостаточно средств")
		return 406, nil // можно создать свою кастомную функцию ошибку о недостатке средств
	}

	//log.Println("NewBalance: ", newBalance)

	// Начало транзакции
	tx, err := DB.Begin()
	if err != nil {
		log.Println("CreateOrder error: ", err)
		return 500, err
	}

	_, err = tx.Exec(`UPDATE users SET balance = $1 WHERE id = $2;`, newBalance, userID)
	if err != nil {
		_ = tx.Rollback()
		log.Println("CreateOrder error: ", err)
		return 500, err // ?
	}

	_, err = tx.Exec(
		`INSERT INTO orders (id, user_id, product_id, cost)
		VALUES ($1, $2, $3, $4)`,
		orderID, userID, productID, cost)
	if err != nil {
		_ = tx.Rollback()
		log.Println("CreateOrder error: ", err)
		return 400, err
	}

	if err := tx.Commit(); err != nil {
		log.Println("CreateReserve error: ", err)
		return 500, err
	}
	log.Println("CreateReserve done")

	return 201, nil
}

func UnReserve(orderID int) (httpStatus int, err error) {

	var userID, cost, oldBalance int

	tx, err := DB.Begin()
	if err != nil {
		log.Println("CreateOrder error: ", err)
		return 500, err
	}

	q := `SELECT user_id, cost FROM orders WHERE id=$1;`
	res := tx.QueryRow(q, orderID)
	if err := res.Scan(&userID, &cost); err != nil {
		log.Println("UnReserve error: ", err)
		return 404, err
	}

	q = `SELECT balance FROM users WHERE id=$1;`
	res = tx.QueryRow(q, userID)
	if err := res.Scan(&oldBalance); err != nil {
		log.Println("UnReserve error: ", err)
		return 500, err
	}

	q = `UPDATE users SET balance=$1 WHERE id=$2;`
	_, err = tx.Exec(q, oldBalance+cost, userID)
	if err != nil {
		log.Println("UnReserve error: ", err)
		return 404, err
	}

	q = `UPDATE orders SET cost=0 WHERE id=$1;`
	_, err = tx.Exec(q, orderID)
	if err != nil {
		log.Println("UnReserve error: ", err)
		return 500, err
	}

	if err := tx.Commit(); err != nil {
		log.Println("UnReserve error: ", err)
		return 500, err
	}
	log.Println("CreateReserve done")

	return 200, nil
}
