package main

import (
	"go-bancario/Controllers"
	"go-bancario/Models"
	"log"
	"net/http"
	"path/filepath"
	"runtime/debug"
)

func safeHandle(pattern string, h http.HandlerFunc) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic en %s: %v\n%s", pattern, rec, string(debug.Stack()))
				http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			}
		}()
		h(w, r)
	})
}

func main() {
	if err := Models.InitDB(); err != nil {
		log.Fatal("Error inicializando la base de datos:", err)
	}
	safeHandle("/cuentas", Controllers.AccountHandler)
	safeHandle("/clientes", Controllers.ClientesHandler)
	safeHandle("/tipos-cuenta", Controllers.TiposCuentaHandler)

	// Servir archivos estáticos desde ./public
	fs := http.FileServer(http.Dir(filepath.FromSlash("public")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Páginas (archivos estáticos)
	safeHandle("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.FromSlash("Views/index.html"))
	})
	safeHandle("/crear", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.FromSlash("Views/crear_cuenta.html"))
	})
	safeHandle("/transferir", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.FromSlash("Views/transferir.html"))
	})
	safeHandle("/listado", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.FromSlash("Views/listado.html"))
	})
	log.Println("Servidor iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
