package models

// Estructura para representar una solicitud de creación de billetera
type Log struct {
	ID     			int      		`json:"id"`
	DNI    			string   		`json:"dni"`
	Date 	 			string   		`json:"date"`
	Stage 			string   	 	`json:"stage"`
}


