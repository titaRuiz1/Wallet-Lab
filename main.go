package main

import(
	"net/http"
	"github.com/gorilla/mux"
	"github.com/titaruiz1/wallet-lab/controllers"
	"fmt"
	"github.com/titaruiz1/wallet-lab/db"

)

func main(){

	db.UpDb()
	router := mux.NewRouter()

	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	// Definir una ruta
	router.HandleFunc("/wallets", controllers.CreateWallet).Methods("POST")
	router.HandleFunc("/wallets/{dni}", controllers.DeleteWallet).Methods("DELETE")
	router.HandleFunc("/wallets/{dni}", controllers.WalletStatus).Methods("GET")
	router.HandleFunc("/transaction", controllers.CreateTransaction).Methods("POST")
	router.HandleFunc("/wallet/{id}", controllers.GetMovements).Methods("GET")
	
	handler := corsOptions.Handler(router)
	db.Db.PingOrDie()
	// Iniciar el servidor HTTP
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", handler)

}