package src

type User struct {
	ID      int `json:"id"`
	Balance int `json:"balance"`
}

type Product struct {
	ID   int `json:"id"`
	Cost int `json:"cost"`
}

type Order struct {
	ID        int `json:"id"`
	UserId    int `json:"user_id"`
	ProductId int `json:"product_id"`
	Cost      int `json:"cost"`
	StartedAt int `json:"started_at"`
}

type Report struct {
	ID int `json:"id"`
}

type Transaction struct {
	FromID int `json:"id_from"`
	TOID   int `json:"id_to"`
	Cost   int `json:"cost"`
}

type MonthOrder struct {
	ID      int `json:"product_id"`
	Revenue int `json:"revenue"`
}

type MonthYear struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type Operation struct {
	Page      int    `json:"page"`
	Cost      int    `json:"cost"`
	ProductID int    `json:"product_id"`
	Comment   string `json:"comment"`
	StartedAt string `json:"started_at"`
}

type ConfigureOperationReview struct {
	Id         int  `json:"id"`
	SortByDate bool `json:"sortByDate"`
	SortBySum  bool `json:"sortBySum"`
	Descending bool `json:"descending"`
}
