package Controllers

import (
	"encoding/json"
	"fmt"
	"go-bancario/Models"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Devolver cuentas como JSON para el listado
		cuentas, err := Models.ListarCuentas()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cuentas)
		return
	}
	if r.Method == "POST" {
		log.Println("POST /cuentas - inicio")
		if err := r.ParseForm(); err != nil {
			log.Println("Error parsing form:", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "Error parsing form"}`))
			return
		}
		idTipo, _ := strconv.Atoi(r.FormValue("id_tipo_cuenta"))

		      // datos del cliente desde el formulario
		cliente := Models.Cliente{
			NumeroDocumento: r.FormValue("numero_documento"),
			TipoDocumento:   r.FormValue("tipo_documento"),
			PrimerNombre:    r.FormValue("primer_nombre"),
			PrimerApellido:  r.FormValue("primer_apellido"),
			Email:           r.FormValue("email"),
		}
		// opcionales
		cliente.SegundoNombre.String = r.FormValue("segundo_nombre")
		cliente.SegundoNombre.Valid = cliente.SegundoNombre.String != ""
		cliente.SegundoApellido.String = r.FormValue("segundo_apellido")
		cliente.SegundoApellido.Valid = cliente.SegundoApellido.String != ""
		cliente.Telefono.String = r.FormValue("telefono")
		cliente.Telefono.Valid = cliente.Telefono.String != ""

        // fecha de nacimiento requerida por esquema
        fnacStr := r.FormValue("fecha_nacimiento")
        if fnac, err := time.Parse("2006-01-02", fnacStr); err == nil {
            cliente.FechaNacimiento = fnac
        }

        // validaciones mínimas servidor
        if cliente.NumeroDocumento == "" || cliente.TipoDocumento == "" || cliente.PrimerNombre == "" || cliente.PrimerApellido == "" || cliente.Email == "" || cliente.FechaNacimiento.IsZero() {
        	w.Header().Set("Content-Type", "application/json")
        	w.WriteHeader(http.StatusBadRequest)
        	w.Write([]byte(`{"message": "Campos requeridos faltantes"}`))
        	return
        }
        if idTipo == 0 {
        	log.Println("idTipo es 0, valores del form:", r.FormValue("id_tipo_cuenta"))
        	w.Header().Set("Content-Type", "application/json")
        	w.WriteHeader(http.StatusBadRequest)
        	w.Write([]byte(`{"message": "Seleccione un tipo de cuenta"}`))
        	return
        }

		// si ya existe el cliente, reutilizar ID; si no, crearlo
		log.Println("POST /cuentas - crear/obtener cliente")
		idCliente, err := Models.CrearOObtenerCliente(cliente)
		if err != nil {
			log.Println("error crear/obtener cliente:", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Error creando cliente: %s"}`, err.Error())))
			return
		}

		// validar tipo de cuenta
		log.Println("POST /cuentas - listar tipos")
		tipos, err := Models.ListarTiposCuenta()
		if err != nil {
			log.Println("error listar tipos:", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "No se pudieron cargar los tipos de cuenta"}`))
			return
		}
		tipoValido := false
		for _, t := range tipos {
			if t.IDTipoCuenta == idTipo {
				tipoValido = true
				break
			}
		}
		if !tipoValido {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "Tipo de cuenta inválido"}`))
			return
		}

		// generar número de cuenta y valores por defecto
		log.Println("POST /cuentas - generar numero")
		numero, err := Models.GenerarNumeroCuenta(nil)
		if err != nil {
			log.Println("error generar numero:", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "No se pudo generar el número de cuenta"}`))
			return
		}
        saldo := 500000.00
		fecha := time.Now()
        // obtener id estado 'Activa'
        var idEstado int
        if err := Models.DB.QueryRow(`SELECT id_estado FROM estado_cuenta WHERE nombre_estado='Activa'`).Scan(&idEstado); err != nil {
            log.Println("error estado activa:", err)
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte(`{"message": "No se pudo obtener el estado 'Activa'"}`))
            return
        }

		log.Println("POST /cuentas - crear cuenta")
        cuentaID, err := Models.CrearCuenta(Models.Cuenta{
			NumeroCuenta:  numero,
			IDCliente:     idCliente,
			IDTipoCuenta:  idTipo,
            IDEstado:      idEstado,
            SaldoActual:   saldo,
            SaldoDisp:     saldo,
			FechaApertura: fecha,
		})
		if err != nil {
			log.Println("error crear cuenta:", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Error creando cuenta: %s"}`, err.Error())))
			return
		}
		// devolver JSON para AJAX
		log.Println("POST /cuentas - exito, devolviendo JSON")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"message": "Cuenta creada exitosamente", "numero_cuenta": "%s", "id_cuenta": %d}`, numero, cuentaID)))
		return
	}
	if r.Method == "PUT" {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error parsing form: " + err.Error()))
			return
		}
		fromID, _ := strconv.Atoi(r.FormValue("id_origen"))
		toID, _ := strconv.Atoi(r.FormValue("id_destino"))
		tipoID, _ := strconv.Atoi(r.FormValue("id_tipo_transaccion"))
		amount, _ := strconv.ParseFloat(r.FormValue("monto"), 64)
		referencia := r.FormValue("referencia")
		descripcion := r.FormValue("descripcion")

		// Generar referencia única si no se proporciona
		if referencia == "" {
			referencia = Models.GenerarReferenciaUnica()
		}

		// Validaciones básicas
		if fromID == 0 || toID == 0 || tipoID == 0 || amount <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Campos requeridos faltantes o inválidos: fromID="+strconv.Itoa(fromID)+", toID="+strconv.Itoa(toID)+", tipoID="+strconv.Itoa(tipoID)+", amount="+strconv.FormatFloat(amount, 'f', -1, 64)))
			return
		}

		// Verificar que las cuentas existen
		var count int
		if err := Models.DB.QueryRow(`SELECT COUNT(*) FROM cuenta WHERE id_cuenta = ?`, fromID).Scan(&count); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error verificando cuenta origen: " + err.Error()))
			return
		}
		if count == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Cuenta origen no existe (ID: " + strconv.Itoa(fromID) + ")"))
			return
		}
		if err := Models.DB.QueryRow(`SELECT COUNT(*) FROM cuenta WHERE id_cuenta = ?`, toID).Scan(&count); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error verificando cuenta destino: " + err.Error()))
			return
		}
		if count == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Cuenta destino no existe (ID: " + strconv.Itoa(toID) + ")"))
			return
		}

		if err := Models.Transferir(fromID, toID, tipoID, amount, referencia, descripcion); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// Endpoints básicos de clientes
func ClientesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		clientes, err := Models.ListarClientes()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		json.NewEncoder(w).Encode(clientes)
		return
	}
	if r.Method == "POST" {
		c := Models.Cliente{
			NumeroDocumento: r.FormValue("numero_documento"),
			TipoDocumento:   r.FormValue("tipo_documento"),
			PrimerNombre:    r.FormValue("primer_nombre"),
			PrimerApellido:  r.FormValue("primer_apellido"),
			Email:           r.FormValue("email"),
		}
		if _, err := Models.CrearCliente(c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func TiposCuentaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	tipos, err := Models.ListarTiposCuenta()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(tipos)
}

func urlQuery(s string) string { return url.QueryEscape(s) }
