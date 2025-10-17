package Models

import (
	"errors"
)

type Account struct {
	ID      int
	Owner   string
	Balance float64
}

func CreateAccount(owner string, balance float64) error {
	_, err := DB.Exec("INSERT INTO accounts (owner, balance) VALUES (?, ?)", owner, balance)
	return err
}

func GetAccounts() ([]Account, error) {
	rows, err := DB.Query("SELECT id, owner, balance FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accounts []Account
	for rows.Next() {
		var a Account
		if err := rows.Scan(&a.ID, &a.Owner, &a.Balance); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

func Transfer(fromID, toID int, amount float64) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var fromBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ?", fromID).Scan(&fromBalance)
	if err != nil {
		return err
	}
	if fromBalance < amount {
		return errors.New("fondos insuficientes")
	}
	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromID)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toID)
	if err != nil {
		return err
	}
	return tx.Commit()
}
