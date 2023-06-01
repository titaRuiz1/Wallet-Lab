package services

import (
	"github.com/titaruiz1/wallet-lab/models"
	"github.com/titaruiz1/wallet-lab/db"
)

func CreateServiceTransaction(transaction models.Transaction, wallet_sender models.Wallet, wallet_receiver models.Wallet) error{
	// Iniciamos una transacción
	tx, err := db.Db.Begin()
	
	if err != nil {
		return err
	}
	
	const SQLInsertNewWallet = `INSERT INTO transactions (senderid, receiverid, amount, transactiontype)
														VALUES ($1, $2, $3, $4)`
	// Ejecutamos la consulta SQL para insertar un nuevo registro en la tabla 'wallet' dentro de la transacción.
	_, err = tx.Exec(SQLInsertNewWallet,transaction.SenderId, transaction.ReceiverId, transaction.Amount, transaction.Type)
	if err != nil {
		// Si algo salió mal, hacemos un rollback de la transacción
		tx.Rollback()
		return err
	}
	
	// Actualizamos el saldo de la billetera del remitente
	const SQLUpdateSenderWallet = `UPDATE wallets SET balance = balance - $1 WHERE id = $2`
	_, err = tx.Exec(SQLUpdateSenderWallet, transaction.Amount, wallet_sender.ID)
	if err != nil {
		// Si algo salió mal, hacemos un rollback de la transacción
		tx.Rollback()
		return err
	}

	// Actualizamos el saldo de la billetera del destinatario
	const SQLUpdateReceiverWallet = `UPDATE wallets SET balance = balance + $1 WHERE id = $2`
	_, err = tx.Exec(SQLUpdateReceiverWallet, transaction.Amount, wallet_receiver.ID)
	if err != nil {
		// Si algo salió mal, hacemos un rollback de la transacción
		tx.Rollback()
		return err
	}
	// Si todo salió bien, hacemos commit de la transacción
	err = tx.Commit()
	if err != nil {
		return err
	}
	
	return nil
}