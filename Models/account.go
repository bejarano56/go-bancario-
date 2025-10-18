package Models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Catálogos
type TipoCuenta struct {
	IDTipoCuenta int
	NombreTipo   string
    SaldoMinimo  float64
}

type TipoTransaccion struct {
	IDTipoTransaccion int
	NombreTipo        string
}

// Principales
type Cliente struct {
	IDCliente       int
	NumeroDocumento string
	TipoDocumento   string
	PrimerNombre    string
	SegundoNombre   sql.NullString
	PrimerApellido  string
	SegundoApellido sql.NullString
	Email           string
	Telefono        sql.NullString
    FechaNacimiento time.Time
	FechaRegistro   time.Time
}

type Cuenta struct {
	IDCuenta      int
	NumeroCuenta  string
	IDCliente     int
	IDTipoCuenta  int
    IDEstado      int
    SaldoActual   float64
    SaldoDisp     float64
    FechaApertura time.Time
    EstadoNombre  string  // Campo temporal para el nombre del estado
    TipoNombre    string  // Campo temporal para el nombre del tipo
}

type Transaccion struct {
	IDTransaccion     int
	IDCuentaOrigen    int
	IDCuentaDestino   sql.NullInt64
	IDTipoTransaccion int
    IDEstadoTrans     int
	Monto             float64
	Referencia        string
	Descripcion       sql.NullString
    FechaTransaccion  time.Time
    FechaProcesamiento sql.NullTime
}

type Auditoria struct {
	IDAuditoria          int
	IDTransaccion        int
	SaldoAnteriorOrigen  sql.NullFloat64
	SaldoNuevoOrigen     sql.NullFloat64
	SaldoAnteriorDestino sql.NullFloat64
	SaldoNuevoDestino    sql.NullFloat64
    Fecha                time.Time
}

// Catálogos
func ListarTiposCuenta() ([]TipoCuenta, error) {
    rows, err := DB.Query(`SELECT id_tipo_cuenta, nombre_tipo, saldo_minimo FROM tipo_cuenta ORDER BY nombre_tipo`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TipoCuenta
	for rows.Next() {
		var t TipoCuenta
		if err := rows.Scan(&t.IDTipoCuenta, &t.NombreTipo, &t.SaldoMinimo); err != nil {
			return nil, err
		}
		items = append(items, t)
	}
	if len(items) == 0 {
		if _, err := DB.Exec(`INSERT INTO tipo_cuenta (nombre_tipo, saldo_minimo) VALUES
            ('Ahorros', 0.00),
            ('Corriente', 0.00),
            ('Nómina', 0.00)`); err != nil {
			return nil, err
		}
		rows2, err := DB.Query(`SELECT id_tipo_cuenta, nombre_tipo, saldo_minimo FROM tipo_cuenta ORDER BY nombre_tipo`)
		if err != nil {
			return nil, err
		}
		defer rows2.Close()
		for rows2.Next() {
			var t TipoCuenta
			if err := rows2.Scan(&t.IDTipoCuenta, &t.NombreTipo, &t.SaldoMinimo); err != nil {
				return nil, err
			}
			items = append(items, t)
		}
	}
	return items, nil
}

// Generación de número de cuenta único (intenta hasta obtener uno libre)
func GenerarNumeroCuenta(tx *sql.Tx) (string, error) {
	for i := 0; i < 5; i++ {
		candidate := time.Now().Format("20060102") + RandDigits(8)
		var exists int
		q := `SELECT 1 FROM cuenta WHERE numero_cuenta=? LIMIT 1`
		var err error
		if tx != nil {
			err = tx.QueryRow(q, candidate).Scan(&exists)
		} else {
			err = DB.QueryRow(q, candidate).Scan(&exists)
		}
		if err == sql.ErrNoRows {
			return candidate, nil
		}
		if err == nil {
			continue
		}
		if err != nil && err != sql.ErrNoRows {
			return "", err
		}
	}
	return "", errors.New("no se pudo generar número de cuenta único")
}

// RandDigits retorna una cadena de n dígitos aleatorios
func RandDigits(n int) string {
	const digits = "0123456789"
	b := make([]byte, n)
	// fuente simple basada en time, suficiente para demo
	seed := time.Now().UnixNano()
	for i := 0; i < n; i++ {
		seed = (seed*1664525 + 1013904223) % 2147483647
		// Usar valor absoluto para asegurar índice positivo
		idx := int(seed)
		if idx < 0 {
			idx = -idx
		}
		b[i] = digits[idx%10]
	}
	return string(b)
}

// GenerarReferenciaUnica genera una referencia única para transacciones
func GenerarReferenciaUnica() string {
	for i := 0; i < 5; i++ {
		candidate := "REF-" + time.Now().Format("20060102150405") + RandDigits(6)
		var exists int
		q := `SELECT 1 FROM transaccion WHERE referencia=? LIMIT 1`
		err := DB.QueryRow(q, candidate).Scan(&exists)
		if err == sql.ErrNoRows {
			return candidate
		}
		if err != nil {
			continue
		}
	}
	return "REF-" + time.Now().Format("20060102150405") + RandDigits(6) + "-FALLBACK"
}

// Clientes
func CrearCliente(c Cliente) (int64, error) {
    res, err := DB.Exec(`INSERT INTO cliente (numero_documento, tipo_documento, primer_nombre, segundo_nombre, primer_apellido, segundo_apellido, email, telefono, fecha_nacimiento) VALUES (?,?,?,?,?,?,?,?,?)`,
        c.NumeroDocumento, c.TipoDocumento, c.PrimerNombre, c.SegundoNombre, c.PrimerApellido, c.SegundoApellido, c.Email, c.Telefono, c.FechaNacimiento)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func ListarClientes() ([]Cliente, error) {
    rows, err := DB.Query(`SELECT id_cliente, numero_documento, tipo_documento, primer_nombre, segundo_nombre, primer_apellido, segundo_apellido, email, telefono, fecha_nacimiento, fecha_registro FROM cliente`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Cliente
	for rows.Next() {
		var c Cliente
        if err := rows.Scan(&c.IDCliente, &c.NumeroDocumento, &c.TipoDocumento, &c.PrimerNombre, &c.SegundoNombre, &c.PrimerApellido, &c.SegundoApellido, &c.Email, &c.Telefono, &c.FechaNacimiento, &c.FechaRegistro); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, nil
}

// ObtenerClienteIDPorDocumento busca un cliente por numero_documento y retorna su ID si existe
func ObtenerClienteIDPorDocumento(numeroDocumento string) (int, error) {
	if numeroDocumento == "" {
		return 0, errors.New("numero_documento requerido")
	}
	var cnt int
	if err := DB.QueryRow(`SELECT COUNT(*) FROM cliente WHERE numero_documento = ?`, numeroDocumento).Scan(&cnt); err != nil {
		return 0, err
	}
	if cnt == 0 {
		return 0, sql.ErrNoRows
	}
	var id int
	if err := DB.QueryRow(`SELECT id_cliente FROM cliente WHERE numero_documento = ? LIMIT 1`, numeroDocumento).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// Crea el cliente o retorna el ID si ya existe por clave única (numero_documento o email)
func CrearOObtenerCliente(c Cliente) (int, error) {
    res, err := DB.Exec(`INSERT INTO cliente (numero_documento, tipo_documento, primer_nombre, segundo_nombre, primer_apellido, segundo_apellido, email, telefono, fecha_nacimiento)
        VALUES (?,?,?,?,?,?,?,?,?)
        ON DUPLICATE KEY UPDATE fecha_nacimiento = VALUES(fecha_nacimiento), id_cliente = LAST_INSERT_ID(id_cliente)`,
        c.NumeroDocumento, c.TipoDocumento, c.PrimerNombre, c.SegundoNombre, c.PrimerApellido, c.SegundoApellido, c.Email, c.Telefono, c.FechaNacimiento)
	if err != nil {
		return 0, err
	}
	id64, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id64), nil
}

// Cuentas
func CrearCuenta(c Cuenta) (int64, error) {
    res, err := DB.Exec(`INSERT INTO cuenta (numero_cuenta, id_cliente, id_tipo_cuenta, id_estado, saldo_actual, saldo_disponible, fecha_apertura) VALUES (?,?,?,?,?,?,?)`,
        c.NumeroCuenta, c.IDCliente, c.IDTipoCuenta, c.IDEstado, c.SaldoActual, c.SaldoDisp, c.FechaApertura)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func ListarCuentas() ([]Cuenta, error) {
    rows, err := DB.Query(`SELECT c.id_cuenta, c.numero_cuenta, c.id_cliente, c.id_tipo_cuenta, c.id_estado, c.saldo_actual, c.saldo_disponible, c.fecha_apertura, ec.nombre_estado, tc.nombre_tipo 
        FROM cuenta c 
        LEFT JOIN estado_cuenta ec ON c.id_estado = ec.id_estado
        LEFT JOIN tipo_cuenta tc ON c.id_tipo_cuenta = tc.id_tipo_cuenta`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Cuenta
	for rows.Next() {
		var c Cuenta
		var nombreEstado, nombreTipo sql.NullString
        if err := rows.Scan(&c.IDCuenta, &c.NumeroCuenta, &c.IDCliente, &c.IDTipoCuenta, &c.IDEstado, &c.SaldoActual, &c.SaldoDisp, &c.FechaApertura, &nombreEstado, &nombreTipo); err != nil {
			return nil, err
		}
		// Debug: imprimir valores
		fmt.Printf("DEBUG - ID: %d, Estado: %s (Valid: %t), Tipo: %s (Valid: %t), Saldo: %f\n", 
			c.IDCuenta, nombreEstado.String, nombreEstado.Valid, nombreTipo.String, nombreTipo.Valid, c.SaldoActual)
		
		// Manejar valores NULL
		if nombreEstado.Valid {
			c.EstadoNombre = nombreEstado.String
		} else {
			c.EstadoNombre = fmt.Sprintf("Estado %d", c.IDEstado)
		}
		if nombreTipo.Valid {
			c.TipoNombre = nombreTipo.String
		} else {
			c.TipoNombre = fmt.Sprintf("Tipo %d", c.IDTipoCuenta)
		}
		items = append(items, c)
	}
	return items, nil
}

// Transferencia ACID con auditoría
func Transferir(idOrigen, idDestino, idTipoTransaccion int, monto float64, referencia, descripcion string) error {
	if monto <= 0 {
		return errors.New("monto debe ser > 0")
	}
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Saldos actuales
	   var saldoOrigen float64
	   if err := tx.QueryRow(`SELECT saldo_disponible FROM cuenta WHERE id_cuenta=? FOR UPDATE`, idOrigen).Scan(&saldoOrigen); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("cuenta origen no existe")
		}
		return err
	}
	if saldoOrigen < monto {
		return errors.New("fondos insuficientes")
	}

	   var saldoDestino sql.NullFloat64
	   if err := tx.QueryRow(`SELECT saldo_disponible FROM cuenta WHERE id_cuenta=? FOR UPDATE`, idDestino).Scan(&saldoDestino); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("cuenta destino no existe")
		}
		return err
	}

	// Registrar transacción pendiente
	   // estado inicial: Pendiente (buscar id)
	   var idEstadoPendiente int
	   if err := tx.QueryRow(`SELECT id_estado_transaccion FROM estado_transaccion WHERE nombre_estado = 'Pendiente'`).Scan(&idEstadoPendiente); err != nil {
	       if err == sql.ErrNoRows {
	           return errors.New("estado de transacción 'Pendiente' no encontrado")
	       }
	       return err
	   }
	   res, err := tx.Exec(`INSERT INTO transaccion (id_cuenta_origen, id_cuenta_destino, id_tipo_transaccion, id_estado_transaccion, monto, referencia, descripcion) VALUES (?,?,?,?,?,?,?)`,
	       idOrigen, idDestino, idTipoTransaccion, idEstadoPendiente, monto, referencia, descripcion)
	if err != nil {
		return err
	}
	transID, _ := res.LastInsertId()

	// Debitar y acreditar
	   if _, err := tx.Exec(`UPDATE cuenta SET saldo_actual = saldo_actual - ?, saldo_disponible = saldo_disponible - ? WHERE id_cuenta = ?`, monto, monto, idOrigen); err != nil {
		return err
	}
	   if _, err := tx.Exec(`UPDATE cuenta SET saldo_actual = saldo_actual + ?, saldo_disponible = saldo_disponible + ? WHERE id_cuenta = ?`, monto, monto, idDestino); err != nil {
		return err
	}

	// Nuevos saldos
	var nuevoOrigen float64
	var nuevoDestino float64
	   if err := tx.QueryRow(`SELECT saldo_disponible FROM cuenta WHERE id_cuenta=?`, idOrigen).Scan(&nuevoOrigen); err != nil {
		return err
	}
	   if err := tx.QueryRow(`SELECT saldo_disponible FROM cuenta WHERE id_cuenta=?`, idDestino).Scan(&nuevoDestino); err != nil {
		return err
	}

	// Auditoría
	   if _, err := tx.Exec(`INSERT INTO auditoria_transacciones (id_transaccion, accion, saldo_anterior_origen, saldo_nuevo_origen, saldo_anterior_destino, saldo_nuevo_destino) VALUES (?,?,?,?,?,?)`,
	       transID, "Movimiento", saldoOrigen, nuevoOrigen, saldoDestino, nuevoDestino); err != nil {
		return err
	}

	// Completar transacción
	   var idEstadoCompletada int
	   if err := tx.QueryRow(`SELECT id_estado_transaccion FROM estado_transaccion WHERE nombre_estado = 'Completada'`).Scan(&idEstadoCompletada); err != nil {
	       if err == sql.ErrNoRows {
	           return errors.New("estado de transacción 'Completada' no encontrado")
	       }
	       return err
	   }
	   if _, err := tx.Exec(`UPDATE transaccion SET id_estado_transaccion=?, fecha_procesamiento = NOW() WHERE id_transaccion=?`, idEstadoCompletada, transID); err != nil {
		return err
	}

	return tx.Commit()
}

// Crear cliente ad-hoc desde un nombre completo (por restricciones NOT NULL en CLIENTE)
// Genera numero_documento y email temporales únicos.
func CrearClienteDesdeNombreCompleto(nombreCompleto string) (int64, error) {
	if nombreCompleto == "" {
		return 0, errors.New("nombre requerido")
	}
	partes := splitNombre(nombreCompleto)
	tmpDoc := "TMP-" + time.Now().Format("20060102150405") + RandDigits(3)
	tmpEmail := "tmp" + RandDigits(6) + "@example.com"
	c := Cliente{
		NumeroDocumento: tmpDoc,
		TipoDocumento:   "CC",
		PrimerNombre:    partes[0],
		PrimerApellido:  partes[1],
		Email:           tmpEmail,
	}
	return CrearCliente(c)
}

func splitNombre(full string) [2]string {
	// extracción muy simple: primer token y último token
	var first, last string
	start := -1
	for i := 0; i <= len(full); i++ {
		if i == len(full) || full[i] == ' ' {
			if start >= 0 {
				token := full[start:i]
				if first == "" {
					first = token
				} else {
					last = token
				}
				start = -1
			}
		} else if start < 0 {
			start = i
		}
	}
	if last == "" {
		last = "N/A"
	}
	if first == "" {
		first = "N/A"
	}
	return [2]string{first, last}
}
