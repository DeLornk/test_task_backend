package src

import (
	"log"
)

func FindUser(id int) (User, error) {
	var user User

	sqlQuery := `SELECT id, balance FROM users where id = $1;`

	row := DB.QueryRow(sqlQuery, id)

	log.Println("FindUser: сделали запрос")

	err := row.Scan(&user.ID, &user.Balance)

	if err != nil {
		log.Println("FindUser error: ", err)
		return user, err
	}

	return user, nil
}

func UpdateBalance(newBalance int, id int) (done bool, err error) {
	//q := `UPDATE users set balance = $1 WHERE id = $2`

	q := `INSERT INTO users (balance, id) 
		VALUES ($1, $2)
		ON CONFLICT (id)
		DO UPDATE SET balance=$1;`

	_, err = DB.Exec(q, newBalance, id)

	if err != nil {
		log.Println("UpdateBalance: ", err)
		return false, err
	}

	log.Println("UpdateBalance: успешно")
	return true, err
}

func FromUserToUser(cost, idFrom, idTo int) (httpStatus int, err error) {

	userFrom, err := FindUser(idFrom)
	if err != nil {
		log.Println("FromUserToUser error: ", err)
		return 404, err
	}

	userTo, err := FindUser(idTo)
	if err != nil {
		log.Println("FromUserToUser error: ", err)
		return 404, err
	}

	newBalanceFrom := userFrom.Balance - cost
	newBalanceTo := userTo.Balance + cost

	if newBalanceFrom < 0 {
		log.Println("FromUserToUser error: недостаточно средств")
		return 406, err
	}

	tx, err := DB.Begin()
	if err != nil {
		log.Println("FromUserToUser error: ", err)
		return 500, err
	}

	_, err = tx.Exec(`UPDATE users SET balance = $1 WHERE id = $2;`, newBalanceFrom, idFrom)
	if err != nil {
		_ = tx.Rollback()
		log.Println("FromUserToUser error: ", err)
		return 500, err // ?
	}

	_, err = tx.Exec(`UPDATE users SET balance = $1 WHERE id = $2;`, newBalanceTo, idTo)
	if err != nil {
		_ = tx.Rollback()
		log.Println("FromUserToUser error: ", err)
		return 500, err
	}

	if err := tx.Commit(); err != nil {
		log.Println("FromUserToUser error: ", err)
		return 500, err
	}

	log.Println("FromUserToUser done")
	return 200, nil
}
