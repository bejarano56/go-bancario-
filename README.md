Quick start

1) Importa `migrations.sql` en MySQL y asegúrate de que el usuario `root` no tiene contraseña o ajusta `Models/db.go`.
2) Ejecuta: `go run main.go`.
3) Abre http://localhost:8080/.

Endpoints
- GET /clientes
- POST /clientes
- GET /cuentas
- POST /cuentas
- PUT /cuentas (transferir) body: id_origen, id_destino, id_tipo_transaccion, monto, referencia, descripcion
# Sistema Bancario Go (MVC)

## Ejecución rápida

1. Instala Go y ejecuta:

```powershell
cd "c:\Users\juan david\Documents\Trabajos ADSO\Programación\go-bancario-"
go run main.go
```

2. Inicializa la base de datos:

```powershell
sqlite3 banco.db < migrations.sql
```

3. Prueba los endpoints con Postman o curl:
- GET     /cuentas
- POST    /cuentas (owner, balance)
- PUT     /cuentas (from, to, amount)

## Estructura
- `Models/` lógica y acceso a datos
- `Controllers/` handlers HTTP
- `Views/` vista HTML simple
- `migrations.sql` esquema normalizado y consultas

## ACID
- Las transferencias usan transacciones SQL (`BEGIN`, `COMMIT`, `ROLLBACK`) para garantizar atomicidad, consistencia, aislamiento y durabilidad.

## Consultas SQL
- Separadas en `migrations.sql` para facilitar mantenimiento y normalización.

---

Este sistema es solo un ejemplo educativo, no para producción.
