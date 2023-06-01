package models

// Estructura para representar una billetera creada
type Wallet struct {
	ID          	int       `json:"id"`
	DNI         	string    `json:"dni"`
	CountryID   	string    `json:"country_id"`
	CreationDate	string   	`json:"creation_date"`
	Balance       float64   `json:"balance"`
}

type Data struct{
	CheckID string `json:"check_id"`
	Score int `json:"score"`
}

type ValidationWallet struct{
	Check []Data `json:"checks"`
}

type Transaction struct{
	ID          	int       `json:"id"`
	SenderId int `json:"sender_id"`
	ReceiverId int `json:"receiver_id"`
	Amount float64 `json:"amount"`
	Type string `json:"type"`
	Date string	`json:"creation_date"`
}

type WalletMovements struct{
	ID int `json:"id"`
	Balance float64 `json:"balance"`
	Movements []Transaction `json:"movements"`
}