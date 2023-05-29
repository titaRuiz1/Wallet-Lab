package main

import(
	"net/http"
	"github.com/gorilla/mux"
	"github.com/titaRuiz1/Wallet-lab/controllers"
	"fmt"
	"github.com/titaRuiz1/Wallet-lab/db"

)

func main(){

	db.UpDb()
	router := mux.NewRouter()

	// Definir una ruta
	router.HandleFunc("/wallets", controllers.CreateWallet).Methods("POST")
	router.HandleFunc("/wallets/{dni}", controllers.DeleteWallet).Methods("DELETE")
	router.HandleFunc("/wallets/{dni}", controllers.WalletStatus).Methods("GET")
	
	db.Db.PingOrDie()
	// Iniciar el servidor HTTP
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", router)

}