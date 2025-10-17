package Models

import (
	"database/sql"
	"errors"
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
	FechaRegistro   time.Time
}

type Cuenta struct {
	IDCuenta      int
	NumeroCuenta  string
	IDCliente     int
	IDTipoCuenta  int
	Saldo         float64
	Activa        bool
	FechaApertura time.Time
}

type Transaccion struct {
	IDTransaccion     int
	IDCuentaOrigen    int
	IDCuentaDestino   sql.NullInt64
	IDTipoTransaccion int
	Monto             float64
	Referencia        string
	Descripcion       sql.NullString
	Fecha             time.Time
	Completada        bool
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
	rows, err := DB.Query(`SELECT id_tipo_cuenta, nombre_tipo, saldo_minimo FROM TIPO_CUENTA ORDER BY nombre_tipo`)
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
		if _, err := DB.Exec(`INSERT INTO TIPO_CUENTA (nombre_tipo, saldo_minimo) VALUES
            ('Ahorros', 0.00),
            ('Corriente', 0.00),
            ('Nómina', 0.00)`); err != nil {
			return nil, err
		}
		rows2, err := DB.Query(`SELECT id_tipo_cuenta, nombre_tipo, saldo_minimo FROM TIPO_CUENTA ORDER BY nombre_tipo`)
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
		q := `SELECT 1 FROM CUENTA WHERE numero_cuenta=? LIMIT 1`
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
		b[i] = digits[int(seed)%10]
	}
	return string(b)
}

// Clientes
func CrearCliente(c Cliente) (int64, error) {
	res, err := DB.Exec(`INSERT INTO CLIENTE (numero_documento, tipo_documento, primer_nombre, segundo_nombre, primer_apellido, segundo_apellido, email, telefono) VALUES (?,?,?,?,?,?,?,?)`,
		c.NumeroDocumento, c.TipoDocumento, c.PrimerNombre, c.SegundoNombre, c.PrimerApellido, c.SegundoApellido, c.Email, c.Telefono)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func ListarClientes() ([]Cliente, error) {
	rows, err := DB.Query(`SELECT id_cliente, numero_documento, tipo_documento, primer_nombre, segundo_nombre, primer_apellido, segundo_apellido, email, telefono, fecha_registro FROM CLIENTE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Cliente
	for rows.Next() {
		var c Cliente
		if err := rows.Scan(&c.IDCliente, &c.NumeroDocumento, &c.TipoDocumento, &c.PrimerNombre, &c.SegundoNombre, &c.PrimerApellido, &c.SegundoApellido, &c.Email, &c.Telefono, &c.FechaRegistro); err != nil {
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
	if err := DB.QueryRow(`SELECT COUNT(*) FROM CLIENTE WHERE numero_documento = ?`, numeroDocumento).Scan(&cnt); err != nil {
		return 0, err
	}
	if cnt == 0 {
		return 0, sql.ErrNoRows
	}
	var id int
	if err := DB.QueryRow(`SELECT id_cliente FROM CLIENTE WHERE numero_documento = ? LIMIT 1`, numeroDocumento).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// Crea el cliente o retorna el ID si ya existe por clave única (numero_documento o email)
func CrearOObtenerCliente(c Cliente) (int, error) {
	res, err := DB.Exec(`INSERT INTO CLIENTE (numero_documento, tipo_documento, primer_nombre, segundo_nombre, primer_apellido, segundo_apellido, email, telefono)
        VALUES (?,?,?,?,?,?,?,?)
        ON DUPLICATE KEY UPDATE id_cliente = LAST_INSERT_ID(id_cliente)`,
		c.NumeroDocumento, c.TipoDocumento, c.PrimerNombre, c.SegundoNombre, c.PrimerApellido, c.SegundoApellido, c.Email, c.Telefono)
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
	res, err := DB.Exec(`INSERT INTO CUENTA (numero_cuenta, id_cliente, id_tipo_cuenta, saldo, activa, fecha_apertura) VALUES (?,?,?,?,?,?)`,
		c.NumeroCuenta, c.IDCliente, c.IDTipoCuenta, c.Saldo, c.Activa, c.FechaApertura)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func ListarCuentas() ([]Cuenta, error) {
	rows, err := DB.Query(`SELECT id_cuenta, numero_cuenta, id_cliente, id_tipo_cuenta, saldo, activa, fecha_apertura FROM CUENTA`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Cuenta
	for rows.Next() {
		var c Cuenta
		if err := rows.Scan(&c.IDCuenta, &c.NumeroCuenta, &c.IDCliente, &c.IDTipoCuenta, &c.Saldo, &c.Activa, &c.FechaApertura); err != nil {
			return nil, err
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
	if err := tx.QueryRow(`SELECT saldo FROM CUENTA WHERE id_cuenta=? FOR UPDATE`, idOrigen).Scan(&saldoOrigen); err != nil {
		return err
	}
	if saldoOrigen < monto {
		return errors.New("fondos insuficientes")
	}

	var saldoDestino sql.NullFloat64
	if err := tx.QueryRow(`SELECT saldo FROM CUENTA WHERE id_cuenta=? FOR UPDATE`, idDestino).Scan(&saldoDestino); err != nil {
		return err
	}

	// Registrar transacción pendiente
	res, err := tx.Exec(`INSERT INTO TRANSACCION (id_cuenta_origen, id_cuenta_destino, id_tipo_transaccion, monto, referencia, descripcion, completada) VALUES (?,?,?,?,?,?,false)`,
		idOrigen, idDestino, idTipoTransaccion, monto, referencia, descripcion)
	if err != nil {
		return err
	}
	transID, _ := res.LastInsertId()

	// Debitar y acreditar
	if _, err := tx.Exec(`UPDATE CUENTA SET saldo = saldo - ? WHERE id_cuenta = ?`, monto, idOrigen); err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE CUENTA SET saldo = saldo + ? WHERE id_cuenta = ?`, monto, idDestino); err != nil {
		return err
	}

	// Nuevos saldos
	var nuevoOrigen float64
	var nuevoDestino float64
	if err := tx.QueryRow(`SELECT saldo FROM CUENTA WHERE id_cuenta=?`, idOrigen).Scan(&nuevoOrigen); err != nil {
		return err
	}
	if err := tx.QueryRow(`SELECT saldo FROM CUENTA WHERE id_cuenta=?`, idDestino).Scan(&nuevoDestino); err != nil {
		return err
	}

	// Auditoría
	if _, err := tx.Exec(`INSERT INTO AUDITORIA (id_transaccion, saldo_anterior_origen, saldo_nuevo_origen, saldo_anterior_destino, saldo_nuevo_destino) VALUES (?,?,?,?,?)`,
		transID, saldoOrigen, nuevoOrigen, saldoDestino, nuevoDestino); err != nil {
		return err
	}

	// Completar transacción
	if _, err := tx.Exec(`UPDATE TRANSACCION SET completada=true WHERE id_transaccion=?`, transID); err != nil {
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
