package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/titaruiz1/wallet-lab/db"
	"github.com/titaruiz1/wallet-lab/models"
	"github.com/titaruiz1/wallet-lab/services"
)

func GetMovements(w http.ResponseWriter, r *http.Request){
	idStr, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Error con el id", http.StatusBadRequest)
	}

	var wallet models.Wallet
	row, err := db.Db.Query("SELECT * FROM wallets WHERE id = $1", idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer row.Close()

	if row.Next() {
		err = row.Scan(&wallet.ID, &wallet.DNI, &wallet.CountryID, &wallet.CreationDate, &wallet.Balance )
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// No se encontró una billetera con el ID especificado
		http.Error(w, "Billetera de envío no encontrada", http.StatusNotFound)
		return
	}

	var movements models.WalletMovements
	movements.ID = wallet.ID
	movements.Balance = wallet.Balance
	var arrayTransactions []models.Transaction

	query := `
	SELECT * 
	FROM transactions
	WHERE senderid = $1 OR receiverid = $1
`

	rows, err := db.Db.Query(query, wallet.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.ID, &transaction.SenderId, &transaction.ReceiverId, &transaction.Amount, &transaction.Type, &transaction.Date)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		arrayTransactions = append(arrayTransactions, transaction)
	}
	
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	movements.Movements = arrayTransactions
	json.NewEncoder(w).Encode(&movements)
	fmt.Println(movements)
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var wallet_sender, wallet_receiver models.Wallet
	row, err := db.Db.Query("SELECT * FROM wallets WHERE id = $1", transaction.SenderId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer row.Close()
	if row.Next() {
		err = row.Scan(&wallet_sender.ID, &wallet_sender.DNI, &wallet_sender.CountryID, &wallet_sender.CreationDate, &wallet_sender.Balance )
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// No se encontró una billetera con el ID especificado
		http.Error(w, "Billetera de envío no encontrada", http.StatusNotFound)
		return
	}

	// Encontrar billetera que va a recibir la transaccion.
	row2, err := db.Db.Query("SELECT * FROM wallets WHERE id = $1", transaction.ReceiverId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer row2.Close()

	if row2.Next() {
		err = row2.Scan(&wallet_receiver.ID, &wallet_receiver.DNI, &wallet_receiver.CountryID, &wallet_receiver.CreationDate, &wallet_receiver.Balance )
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// No se encontró una billetera con el ID especificado
		http.Error(w, "Billetera de envío no encontrada", http.StatusNotFound)
		return
	}

	if(wallet_sender.Balance < transaction.Amount){
		fmt.Println(err)
			http.Error(w, "Saldo insuficiente", http.StatusInternalServerError)
			return
	}else{
		err = services.CreateServiceTransaction(transaction, wallet_sender, wallet_receiver)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	json.NewEncoder(w).Encode("Operacion exitosa")
}
