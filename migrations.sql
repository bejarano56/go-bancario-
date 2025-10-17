CREATE DATABASE IF NOT EXISTS sistema_bancario;
USE sistema_bancario;

-- =============================================
-- TABLAS DE CATÁLOGO
-- =============================================

CREATE TABLE IF NOT EXISTS TIPO_CUENTA (
    id_tipo_cuenta INT AUTO_INCREMENT PRIMARY KEY,
    nombre_tipo VARCHAR(50) NOT NULL UNIQUE,
    saldo_minimo DECIMAL(15,2) DEFAULT 0.00
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS TIPO_TRANSACCION (
    id_tipo_transaccion INT AUTO_INCREMENT PRIMARY KEY,
    nombre_tipo VARCHAR(50) NOT NULL UNIQUE
) ENGINE=InnoDB;

-- =============================================
-- TABLAS PRINCIPALES
-- =============================================

CREATE TABLE IF NOT EXISTS CLIENTE (
    id_cliente INT AUTO_INCREMENT PRIMARY KEY,
    numero_documento VARCHAR(20) NOT NULL UNIQUE,
    tipo_documento VARCHAR(20) NOT NULL,
    primer_nombre VARCHAR(50) NOT NULL,
    segundo_nombre VARCHAR(50),
    primer_apellido VARCHAR(50) NOT NULL,
    segundo_apellido VARCHAR(50),
    email VARCHAR(100) NOT NULL UNIQUE,
    telefono VARCHAR(20),
    fecha_registro DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS CUENTA (
    id_cuenta INT AUTO_INCREMENT PRIMARY KEY,
    numero_cuenta VARCHAR(20) NOT NULL UNIQUE,
    id_cliente INT NOT NULL,
    id_tipo_cuenta INT NOT NULL,
    saldo DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    activa BOOLEAN DEFAULT TRUE,
    fecha_apertura DATE NOT NULL,
    FOREIGN KEY (id_cliente) REFERENCES CLIENTE(id_cliente),
    FOREIGN KEY (id_tipo_cuenta) REFERENCES TIPO_CUENTA(id_tipo_cuenta),
    CHECK (saldo >= 0)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS TRANSACCION (
    id_transaccion INT AUTO_INCREMENT PRIMARY KEY,
    id_cuenta_origen INT NOT NULL,
    id_cuenta_destino INT,
    id_tipo_transaccion INT NOT NULL,
    monto DECIMAL(15,2) NOT NULL,
    referencia VARCHAR(100) NOT NULL UNIQUE,
    descripcion VARCHAR(255),
    fecha DATETIME DEFAULT CURRENT_TIMESTAMP,
    completada BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_cuenta_origen) REFERENCES CUENTA(id_cuenta),
    FOREIGN KEY (id_cuenta_destino) REFERENCES CUENTA(id_cuenta),
    FOREIGN KEY (id_tipo_transaccion) REFERENCES TIPO_TRANSACCION(id_tipo_transaccion),
    CHECK (monto > 0)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS AUDITORIA (
    id_auditoria INT AUTO_INCREMENT PRIMARY KEY,
    id_transaccion INT NOT NULL,
    saldo_anterior_origen DECIMAL(15,2),
    saldo_nuevo_origen DECIMAL(15,2),
    saldo_anterior_destino DECIMAL(15,2),
    saldo_nuevo_destino DECIMAL(15,2),
    fecha DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_transaccion) REFERENCES TRANSACCION(id_transaccion)
) ENGINE=InnoDB;

-- =============================================
-- DATOS INICIALES
-- =============================================

INSERT INTO TIPO_CUENTA (nombre_tipo, saldo_minimo) VALUES
('Ahorros', 0.00),
('Corriente', 0.00),
('Nómina', 0.00)
ON DUPLICATE KEY UPDATE nombre_tipo = VALUES(nombre_tipo);

INSERT INTO TIPO_TRANSACCION (nombre_tipo) VALUES
('Transferencia'),
('Depósito'),
('Retiro'),
('Pago')
ON DUPLICATE KEY UPDATE nombre_tipo = VALUES(nombre_tipo);
