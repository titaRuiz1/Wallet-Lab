package services

import (
	"github.com/titaruiz1/wallet-lab/db"
	"github.com/titaruiz1/wallet-lab/models"
	
	"log"
	"fmt"
)


func Create(wallet models.Wallet) error {
	// Iniciamos una transacción
	tx, err := db.Db.Begin()
	
	if err != nil {
		log.Fatal(err)
		
	}
	
	const SQLInsertNewWallet = `INSERT INTO wallets (dni, country_id, balance)
														VALUES ($1, $2, $3)`
	// Ejecutamos la consulta SQL para insertar un nuevo registro en la tabla 'wallet' dentro de la transacción.
	_, err = tx.Exec(SQLInsertNewWallet,wallet.DNI, wallet.CountryID, wallet.Balance)
	if err != nil {
		// Si algo salió mal, hacemos un rollback de la transacción
		tx.Rollback()
		log.Fatal(err)
	}

	// Si todo salió bien, hacemos commit de la transacción
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	
	return nil
}


func CkeckIfExistWallet(dniStr string)(bool, error){
	// Verificar si el DNI con el country  ya existe en la base de datos
	existingWallet, err := GetWallet(dniStr)
	if err != nil {
		return false,err
	}

	return existingWallet.DNI != "",nil
	}



func GetWallet(dniStr string) (models.Wallet, error) {
	var wallet models.Wallet
	rows, err := db.Db.Query("SELECT * FROM wallets WHERE dni = $1",dniStr)
	if err != nil {
		
		fmt.Println(err)
		return wallet, err
	}
	defer rows.Close()
	fmt.Println(rows)
	// Itera sobre cada fila en 'rows' y crea una instancia de 'models.Wallet' con los valores de cada columna.
	for rows.Next() {
		err := rows.Scan(&wallet.ID, &wallet.DNI, &wallet.CountryID, &wallet.CreationDate)
		if err != nil {
			fmt.Println(err)
			return wallet, err
		}

	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return wallet, err
	}
	return wallet, nil
}


func Delete(wallet models.Wallet) error {
	// Iniciamos una transacción
	tx, err := db.Db.Begin()
	
	if err != nil {
		log.Fatal(err)	
		return err
	}
	defer tx.Rollback()

	const SQLDeleteWallet = `DELETE FROM wallets WHERE dni = $1`
	// Ejecutamos la consulta SQL para eliminar el registro correspondiente al 'dni' en la transacción
	_, err = tx.Exec(SQLDeleteWallet, wallet.DNI)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
	

