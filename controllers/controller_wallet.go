package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"log"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/titaruiz1/wallet-lab/models"
	"github.com/titaruiz1/wallet-lab/services"
	"strconv"
)



func CreateWallet(w http.ResponseWriter, r *http.Request) {
	nationalID := r.URL.Query().Get("national_id")
	country := r.URL.Query().Get("country")
	balance, err := strconv.ParseFloat(r.URL.Query().Get("balance"), 64)
	if err != nil {
		fmt.Println("Error al convertir el monto de balance en float:", err)
		return
	}
	exists, err := services.CkeckIfExistWallet(nationalID)
	if err != nil {
		http.Error(w,"Error verifying if wallet exist", http.StatusInternalServerError)
	}

	if exists {
		http.Error(w, "Wallet already exist", http.StatusConflict)
	}
	
	// defer response.Body.Close()
	// Verificar la respuesta deseada
	response, err:= getChecksAPI(nationalID)
	if err != nil {
		http.Error(w,"Error requesting external provider", http.StatusInternalServerError)
	}

	defer response.Body.Close()
	// Decodificar el cuerpo de la solicitud en una estructura de datos de Wallet
	newWallet := models.Wallet{
		DNI:       nationalID,
		CountryID: country,
		Balance:   balance,
	}
	
	var walletData models.ValidationWallet
	json.NewDecoder(response.Body).Decode(&walletData)

	// Acceder a los campos específicos en walletData
	
	if walletData.Check[0].Score == 1 {
		json.NewEncoder(w).Encode("Creando...")

		// Crear la billetera en la base de datos
			
		err := services.Create(newWallet)
		if err != nil {
			http.Error(w, "Failed to create wallet", http.StatusInternalServerError)
			return
		}
		// La billetera se creó correctamente
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Wallet created successfully")

	} else {
		json.NewEncoder(w).Encode("La persona tiene antecedentes.")
	}
}

func getChecksAPI(nationalID string) (*http.Response, error) {
	err := godotenv.Load()
	if err != nil {
			log.Fatal("Error loading .env file")
	}
	// Construir la URL con el parámetro national_id
	url := fmt.Sprintf("https://api.checks.truora.com/v1/checks?national_id=%s&country=PE&type=person&user_authorized=true", nationalID)
	apiKey := os.Getenv("API_KEY")

	// Crear una nueva solicitud GET
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Truora-API-Key", apiKey)
	// Realizar la solicitud HTTP
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func DeleteWallet(w http.ResponseWriter, r *http.Request) {
	dniStr := mux.Vars(r)["dni"]

	// Creas un objeto Wallet con los datos necesarios para eliminarlo
	wallet := models.Wallet{
		DNI: dniStr,
	}

	err := services.Delete(wallet)

	if err != nil {
		http.Error(w, "Error deleting wallet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Wallet with DNI %s has been deleted", dniStr)
}

func WalletStatus(w http.ResponseWriter, r *http.Request) {
	dniStr := mux.Vars(r)["dni"]
	fmt.Println(dniStr)
	// Obtén la billetera utilizando el número de DNI
	wallet, err := services.GetWallet(dniStr)
	if err != nil {
		http.Error(w, "Error retrieving wallet", http.StatusInternalServerError)
		return
	}

	// Si no se encontro billetera, devuelve un mensaje correspondiente
	if wallet.DNI == "" {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	// Devuelve la billetera encontrada con un código de estado HTTP 200 (OK)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Wallet with DNI %s found", dniStr)
}
