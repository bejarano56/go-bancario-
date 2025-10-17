package main

import (
	"go-bancario/Controllers"
	"go-bancario/Models"
	"log"
	"net/http"
)

func main() {
	if err := Models.InitDB(); err != nil {
		log.Fatal("Error inicializando la base de datos:", err)
	}
	http.HandleFunc("/cuentas", Controllers.AccountHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "Views/index.html")
	})
	log.Println("Servidor iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
