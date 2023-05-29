package models

// Estructura para representar una billetera creada
type Wallet struct {
	ID          	int       `json:"id"`
	DNI         	string    `json:"dni"`
	CountryID   	string    `json:"country_id"`
	CreationDate	string   	`json:"creation_date"`
}

