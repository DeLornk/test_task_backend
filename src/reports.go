package src

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func OrderToReport(userID, productID, orderID, cost int) (httpStatus int, err error) {
	tx, err := DB.Begin()
	if err != nil {
		log.Println("OrderToReport error: ", err)
		return 500, err
	}

	_, err = tx.Exec(`INSERT INTO reports
							SELECT id, user_id, product_id, cost, started_at
							FROM orders
							WHERE id=$1 and user_id=$2 and product_id=$3 and cost=$4;`, orderID, userID, productID, cost)
	if err != nil {
		_ = tx.Rollback()
		log.Println("OrderToReport error: ", err)
		return 500, err
	}

	_, err = tx.Exec(`UPDATE orders SET cost = 0 WHERE id=$1 and user_id=$2 and product_id=$3 and cost=$4;`,
		orderID, userID, productID, cost)

	if err != nil {
		_ = tx.Rollback()
		log.Println("OrderToReport error: ", err)
		return 500, err
	}

	if err := tx.Commit(); err != nil {
		log.Println("OrderToReport error: ", err)
		return 500, err
	}

	log.Println("OrderToReport done")
	return 200, nil
}

func MonthReport(month, year int) (rprts []MonthOrder, err error) {
	q := `SELECT product_id, SUM(cost) FROM reports WHERE (extract(year from started_at) =$1 and extract(month from started_at)=$2) GROUP BY product_id;`

	rows, err := DB.Query(q, year, month)
	if err != nil {
		log.Println("MonthReport error: ", err)
		return nil, err
	}

	rprts = make([]MonthOrder, 0)

	for rows.Next() {
		var order MonthOrder
		err := rows.Scan(&order.ID, &order.Revenue)
		if err != nil {
			log.Println("MonthReport error: ", err)
			return nil, err
		}
		rprts = append(rprts, order)
	}

	return rprts, err
}

func CreateCSV(records []MonthOrder) (name string, err error) {

	name = "month_report.csv"

	file, err := os.Create(name)
	defer file.Close()
	if err != nil {
		log.Fatalln("CreateCSV error with open file: ", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()

	// Using WriteAll
	var data [][]string
	for _, record := range records {
		row := []string{strconv.Itoa(record.ID), strconv.Itoa(record.Revenue)}
		data = append(data, row)
	}

	err = w.WriteAll(data)
	if err != nil {
		return "", err
	}

	return name, nil
}

func ReportUserOperations(id int, sortByDate, sortBySum, descending bool) (oprtns []Operation, err error) {
	var q string
	if sortByDate {
		if descending {
			q = `SELECT product_id, cost, started_at FROM reports WHERE user_id=$1 ORDER BY started_at DESC;`
		} else {
			q = `SELECT product_id, cost, started_at FROM reports WHERE user_id=$1 ORDER BY started_at ASC;`
		}
	} else if sortBySum {
		if descending {
			q = `SELECT product_id, cost, started_at FROM reports WHERE user_id=$1 ORDER BY cost DESC;`
		} else {
			q = `SELECT product_id, cost, started_at FROM reports WHERE user_id=$1 ORDER BY cost ASC;`
		}
	} else {
		q = `SELECT product_id, cost, started_at FROM reports WHERE user_id=$1`
	}

	rows, err := DB.Query(q, id)
	if err != nil {
		log.Println("ReportUserOperations error: ", err)
	}
	log.Println(id)
	log.Println(q, sortByDate, sortBySum, descending)

	oprtns = make([]Operation, 0)
	i := 1
	for rows.Next() {
		var operation Operation
		operation.Page = i
		operation.Comment = "Оплата услуг"
		err := rows.Scan(&operation.ProductID, &operation.Cost, &operation.StartedAt)
		if err != nil {
			log.Println("ReportUserOperationserror: ", err)
			return nil, err
		}
		oprtns = append(oprtns, operation)
		i++
	}

	return oprtns, nil
}
