package Controllers

import (
	"encoding/json"
	"go-bancario/Models"
	"net/http"
	"strconv"
)

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		accounts, err := Models.GetAccounts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		json.NewEncoder(w).Encode(accounts)
		return
	}
	if r.Method == "POST" {
		owner := r.FormValue("owner")
		balance, _ := strconv.ParseFloat(r.FormValue("balance"), 64)
		if err := Models.CreateAccount(owner, balance); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
	if r.Method == "PUT" {
		fromID, _ := strconv.Atoi(r.FormValue("from"))
		toID, _ := strconv.Atoi(r.FormValue("to"))
		amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
		if err := Models.Transfer(fromID, toID, amount); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
