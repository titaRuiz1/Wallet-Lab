package controllers

import (
	"fmt"
	"log"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/titaruiz1/wallet-lab/models"
	"github.com/titaruiz1/wallet-lab/services"
	// "strconv"
)

func CreateWallet(w http.ResponseWriter,r *http.Request){
	userID  := r.URL.Query().Get("dni")

	// Verificar si el DNI ya existe en la base de datos
	existingWallet, err := services.GetWallet(userID)
	if err != nil {
		http.Error(w, "Error retrieving wallet", http.StatusInternalServerError)
		return
	}

	// Si ya existe una billetera con el mismo DNI, devuelve un mensaje correspondiente
	if existingWallet.DNI != "" {
		http.Error(w, "DNI already exists", http.StatusBadRequest)
		return
	}

	response, err := getChecksAPI(userID)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	// Verificar la respuesta deseada

	// Decodificar el cuerpo de la solicitud en una estructura de datos de Wallet
	newWallet := models.Wallet{ 
		DNI: userID, 
		CountryID: "PE",
	}

	var walletData interface{}
	json.NewDecoder(response.Body).Decode(&walletData)

	// Acceder a los campos específicos en walletData
	dataMap, ok := walletData.(map[string]interface{})["checks"].([]interface{})
	if !ok {
			fmt.Println("El tipo de datos no es el esperado")
			return
	}

	score, ok := dataMap[0].(map[string]interface{})["score"].(float64)
	if !ok {
			fmt.Println("La estructura no coincide con el tipo esperado")
			return
	}

	fmt.Println(score)

	if(score == 1){
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

	}else{
		json.NewEncoder(w).Encode("La persona tiene antecedentes.")
	}
}

func getChecksAPI(nationalID string) (*http.Response, error) {
	// Construir la URL con el parámetro national_id
	url := fmt.Sprintf("https://api.checks.truora.com/v1/checks?national_id=%s&country=PE&type=person&user_authorized=true", nationalID)
	apiKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiIiwiYWRkaXRpb25hbF9kYXRhIjoie30iLCJjbGllbnRfaWQiOiJUQ0k2MmRkOGY2ZTcyMGY2NmVmYTQ2M2Q3ZDQxNzYxYzk0MyIsImV4cCI6MzI2MTU5NDU5MywiZ3JhbnQiOiIiLCJpYXQiOjE2ODQ3OTQ1OTMsImlzcyI6Imh0dHBzOi8vY29nbml0by1pZHAudXMtZWFzdC0xLmFtYXpvbmF3cy5jb20vdXMtZWFzdC0xX09tRlV0bXRMVCIsImp0aSI6ImI1ZmJiODI2LWY0ZjgtNGUxOC05MmExLWE4OTQ1YWI1ZTUxNiIsImtleV9uYW1lIjoicHJ1ZWJhMSIsImtleV90eXBlIjoiYmFja2VuZCIsInVzZXJuYW1lIjoiZ21haWxhc3VudGl0YXJsLXBydWViYTEifQ.kdU1e0oGQn5Dp3B-hcatxyDdfpJH-dhSmX3WityXtwY"
	
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

func DeleteWallet(w http.ResponseWriter,r *http.Request){
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

func WalletStatus(w http.ResponseWriter,r *http.Request){
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




