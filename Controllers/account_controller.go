package Controllers

import (
	"encoding/json"
	"fmt"
	"go-bancario/Models"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
			http.Redirect(w, r, "/crear?error="+urlQuery("Campos requeridos faltantes"), http.StatusSeeOther)
			return
		}
		if idTipo == 0 {
			http.Redirect(w, r, "/crear?error="+urlQuery("Seleccione un tipo de cuenta"), http.StatusSeeOther)
			return
		}

		// si ya existe el cliente, reutilizar ID; si no, crearlo
		log.Println("POST /cuentas - crear/obtener cliente")
		idCliente, err := Models.CrearOObtenerCliente(cliente)
		if err != nil {
			log.Println("error crear/obtener cliente:", err)
			http.Redirect(w, r, "/crear?error="+urlQuery(err.Error()), http.StatusSeeOther)
			return
		}

		// validar tipo de cuenta
		log.Println("POST /cuentas - listar tipos")
		tipos, err := Models.ListarTiposCuenta()
		if err != nil {
			log.Println("error listar tipos:", err)
			http.Redirect(w, r, "/crear?error="+urlQuery("No se pudieron cargar los tipos de cuenta"), http.StatusSeeOther)
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
			http.Error(w, "Tipo de cuenta inválido", http.StatusBadRequest)
			return
		}

		// generar número de cuenta y valores por defecto
		log.Println("POST /cuentas - generar numero")
		numero, err := Models.GenerarNumeroCuenta(nil)
		if err != nil {
			log.Println("error generar numero:", err)
			http.Redirect(w, r, "/crear?error="+urlQuery("No se pudo generar el número de cuenta"), http.StatusSeeOther)
			return
		}
        saldo := 500000.00
		fecha := time.Now()
        // obtener id estado 'Activa'
        var idEstado int
        if err := Models.DB.QueryRow(`SELECT id_estado FROM ESTADO_CUENTA WHERE nombre_estado='Activa'`).Scan(&idEstado); err != nil {
            log.Println("error estado activa:", err)
            http.Redirect(w, r, "/crear?error="+urlQuery("No se pudo obtener el estado 'Activa'"), http.StatusSeeOther)
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
			http.Redirect(w, r, "/crear?error="+urlQuery(err.Error()), http.StatusSeeOther)
			return
		}
		// siempre redirigir con éxito
		log.Println("POST /cuentas - exito, redirigiendo")
		http.Redirect(w, r, "/crear?success=1&numero="+urlQuery(numero)+"&id="+strconv.FormatInt(cuentaID, 10), http.StatusSeeOther)
		return
	}
	if r.Method == "PUT" {
		fromID, _ := strconv.Atoi(r.FormValue("id_origen"))
		toID, _ := strconv.Atoi(r.FormValue("id_destino"))
		tipoID, _ := strconv.Atoi(r.FormValue("id_tipo_transaccion"))
		amount, _ := strconv.ParseFloat(r.FormValue("monto"), 64)
		referencia := r.FormValue("referencia")
		descripcion := r.FormValue("descripcion")
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

// Página de creación: renderiza opciones de tipo de cuenta sin JS
func CrearPageHandler(w http.ResponseWriter, r *http.Request) {
	tipos, err := Models.ListarTiposCuenta()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var opts strings.Builder
	opts.WriteString(`<option value="" disabled selected>Selecciona un tipo</option>`)
	for _, t := range tipos {
		fmt.Fprintf(&opts, `<option value="%d">%s</option>`, t.IDTipoCuenta, t.NombreTipo)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <title>Crear cuenta - Sistema Bancario</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="/static/css/styles.css" rel="stylesheet">
    <script src="https://cdn.jsdelivr.net/npm/sweetalert2@11"></script>
</head>
<body>
<nav class="navbar navbar-expand-lg navbar-dark bg-primary">
  <div class="container-fluid">
    <a class="navbar-brand" href="/">Sistema Bancario</a>
    <div>
      <a class="nav-link d-inline text-white" href="/crear">Crear cuenta</a>
      <a class="nav-link d-inline text-white" href="/transferir">Transferir</a>
      <a class="nav-link d-inline text-white" href="/listado">Cuentas</a>
    </div>
  </div>
  </nav>

<div class="container my-4">
  <div class="col-12 col-lg-8 mx-auto">
    <div class="card shadow-sm">
      <div class="card-header">Crear cuenta</div>
      <div class="card-body">
        <form id="form-cuenta" class="row g-3" method="POST" action="/cuentas">
          <div class="col-6">
            <label class="form-label">Tipo documento</label>
            <select id="tipo_documento" name="tipo_documento" class="form-select" required>
              <option value="CC">Cédula</option>
              <option value="TI">Tarjeta Identidad</option>
              <option value="CE">Cédula Extranjería</option>
            </select>
          </div>
          <div class="col-6">
            <label class="form-label">Número documento</label>
            <input type="text" id="numero_documento" name="numero_documento" class="form-control" required>
          </div>
          <div class="col-6">
            <label class="form-label">Primer nombre</label>
            <input type="text" id="primer_nombre" name="primer_nombre" class="form-control" required>
          </div>
          <div class="col-6">
            <label class="form-label">Segundo nombre</label>
            <input type="text" id="segundo_nombre" name="segundo_nombre" class="form-control">
          </div>
          <div class="col-6">
            <label class="form-label">Primer apellido</label>
            <input type="text" id="primer_apellido" name="primer_apellido" class="form-control" required>
          </div>
          <div class="col-6">
            <label class="form-label">Segundo apellido</label>
            <input type="text" id="segundo_apellido" name="segundo_apellido" class="form-control">
          </div>
          <div class="col-6">
            <label class="form-label">Email</label>
            <input type="email" id="email" name="email" class="form-control" required>
          </div>
          <div class="col-6">
            <label class="form-label">Teléfono</label>
            <input type="text" id="telefono" name="telefono" class="form-control">
          </div>
          <div class="col-6">
            <label class="form-label">Fecha de nacimiento</label>
            <input type="date" id="fecha_nacimiento" name="fecha_nacimiento" class="form-control" required>
          </div>
          <div class="col-12">
            <label class="form-label">Tipo de cuenta</label>
            <select id="id_tipo_cuenta" name="id_tipo_cuenta" class="form-select" required>%s</select>
          </div>
          <div class="col-12">
            <label class="form-label">Saldo inicial</label>
            <input type="text" id="saldo" class="form-control" value="$500.000" disabled>
            <div class="form-text">Moneda: pesos colombianos (COP)</div>
          </div>
          <div class="col-12">
            <button type="submit" class="btn btn-primary w-100">Crear cuenta</button>
          </div>
          <div class="col-12">
            <small class="text-muted">El número de cuenta y el ID del cliente se generan automáticamente.</small>
          </div>
        </form>
      </div>
    </div>
  </div>
</div>

<script>
  (function(){
    const params = new URLSearchParams(window.location.search);
    if (params.get('success') === '1') {
      const numero = params.get('numero') || '';
      const id = params.get('id') || '';
      Swal.fire({ icon: 'success', title: 'Cuenta creada', text: 'Número: ' + numero + '  |  ID: ' + id, timer: 4000, showConfirmButton: false });
      history.replaceState(null, '', window.location.pathname);
    } else if (params.get('error')) {
      Swal.fire({ icon: 'error', title: 'Error', text: params.get('error') });
      history.replaceState(null, '', window.location.pathname);
    }
  })();
</script>

</body>
</html>`, opts.String())
}
