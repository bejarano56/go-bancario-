-- Crear base de datos
CREATE DATABASE IF NOT EXISTS sistema_bancario;
USE sistema_bancario;

-- =============================================
-- TABLAS DE CATÁLOGO (Tipos y Estados)
-- =============================================

-- Tabla: TIPO_CUENTA
CREATE TABLE TIPO_CUENTA (
    id_tipo_cuenta INT AUTO_INCREMENT PRIMARY KEY,
    nombre_tipo VARCHAR(50) NOT NULL UNIQUE,
    descripcion VARCHAR(255),
    saldo_minimo DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    interes_anual DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT chk_interes CHECK (interes_anual >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla: ESTADO_CUENTA
CREATE TABLE ESTADO_CUENTA (
    id_estado INT AUTO_INCREMENT PRIMARY KEY,
    nombre_estado VARCHAR(50) NOT NULL UNIQUE,
    descripcion VARCHAR(255)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla: TIPO_TRANSACCION
CREATE TABLE TIPO_TRANSACCION (
    id_tipo_transaccion INT AUTO_INCREMENT PRIMARY KEY,
    nombre_tipo VARCHAR(50) NOT NULL UNIQUE,
    descripcion VARCHAR(255),
    requiere_destino BOOLEAN NOT NULL DEFAULT TRUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla: ESTADO_TRANSACCION
CREATE TABLE ESTADO_TRANSACCION (
    id_estado_transaccion INT AUTO_INCREMENT PRIMARY KEY,
    nombre_estado VARCHAR(50) NOT NULL UNIQUE,
    descripcion VARCHAR(255)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla: TIPO_TARJETA
CREATE TABLE TIPO_TARJETA (
    id_tipo_tarjeta INT AUTO_INCREMENT PRIMARY KEY,
    nombre_tipo VARCHAR(50) NOT NULL UNIQUE,
    descripcion VARCHAR(255),
    limite_diario DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT chk_limite CHECK (limite_diario >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =============================================
-- TABLAS PRINCIPALES
-- =============================================

-- Tabla: CLIENTE
CREATE TABLE CLIENTE (
    id_cliente INT AUTO_INCREMENT PRIMARY KEY,
    numero_documento VARCHAR(20) NOT NULL UNIQUE,
    tipo_documento VARCHAR(20) NOT NULL,
    primer_nombre VARCHAR(50) NOT NULL,
    segundo_nombre VARCHAR(50),
    primer_apellido VARCHAR(50) NOT NULL,
    segundo_apellido VARCHAR(50),
    email VARCHAR(100) NOT NULL UNIQUE,
    telefono VARCHAR(20),
    fecha_nacimiento DATE NOT NULL,
    fecha_registro DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_documento (numero_documento),
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla: CUENTA
CREATE TABLE CUENTA (
    id_cuenta INT AUTO_INCREMENT PRIMARY KEY,
    numero_cuenta VARCHAR(20) NOT NULL UNIQUE,
    id_cliente INT NOT NULL,
    id_tipo_cuenta INT NOT NULL,
    id_estado INT NOT NULL,
    saldo_actual DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    saldo_disponible DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    fecha_apertura DATE NOT NULL,
    fecha_ultima_actualizacion DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_cuenta_cliente FOREIGN KEY (id_cliente) REFERENCES CLIENTE(id_cliente),
    CONSTRAINT fk_cuenta_tipo FOREIGN KEY (id_tipo_cuenta) REFERENCES TIPO_CUENTA(id_tipo_cuenta),
    CONSTRAINT fk_cuenta_estado FOREIGN KEY (id_estado) REFERENCES ESTADO_CUENTA(id_estado),
    CONSTRAINT chk_saldo_actual CHECK (saldo_actual >= 0),
    CONSTRAINT chk_saldo_disponible CHECK (saldo_disponible >= 0),
    INDEX idx_numero_cuenta (numero_cuenta),
    INDEX idx_cliente (id_cliente)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla: TRANSACCION
CREATE TABLE TRANSACCION (
    id_transaccion INT AUTO_INCREMENT PRIMARY KEY,
    id_cuenta_origen INT NOT NULL,
    id_cuenta_destino INT,
    id_tipo_transaccion INT NOT NULL,
    id_estado_transaccion INT NOT NULL,
    monto DECIMAL(15,2) NOT NULL,
    referencia VARCHAR(100) NOT NULL UNIQUE,
    descripcion VARCHAR(255),
    fecha_transaccion DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    fecha_procesamiento DATETIME,
    CONSTRAINT fk_transaccion_origen FOREIGN KEY (id_cuenta_origen) REFERENCES CUENTA(id_cuenta),
    CONSTRAINT fk_transaccion_destino FOREIGN KEY (id_cuenta_destino) REFERENCES CUENTA(id_cuenta),
    CONSTRAINT fk_transaccion_tipo FOREIGN KEY (id_tipo_transaccion) REFERENCES TIPO_TRANSACCION(id_tipo_transaccion),
    CONSTRAINT fk_transaccion_estado FOREIGN KEY (id_estado_transaccion) REFERENCES ESTADO_TRANSACCION(id_estado_transaccion),
    CONSTRAINT chk_monto CHECK (monto > 0),
    INDEX idx_cuenta_origen (id_cuenta_origen),
    INDEX idx_cuenta_destino (id_cuenta_destino),
    INDEX idx_fecha (fecha_transaccion),
    INDEX idx_referencia (referencia)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla: TARJETA
CREATE TABLE TARJETA (
    id_tarjeta INT AUTO_INCREMENT PRIMARY KEY,
    numero_tarjeta VARCHAR(16) NOT NULL UNIQUE,
    id_cuenta INT NOT NULL,
    id_tipo_tarjeta INT NOT NULL,
    cvv VARCHAR(4) NOT NULL,
    fecha_emision DATE NOT NULL,
    fecha_vencimiento DATE NOT NULL,
    activa BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT fk_tarjeta_cuenta FOREIGN KEY (id_cuenta) REFERENCES CUENTA(id_cuenta),
    CONSTRAINT fk_tarjeta_tipo FOREIGN KEY (id_tipo_tarjeta) REFERENCES TIPO_TARJETA(id_tipo_tarjeta),
    INDEX idx_numero_tarjeta (numero_tarjeta),
    INDEX idx_cuenta (id_cuenta)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla: AUDITORIA_TRANSACCIONES
CREATE TABLE AUDITORIA_TRANSACCIONES (
    id_auditoria INT AUTO_INCREMENT PRIMARY KEY,
    id_transaccion INT NOT NULL,
    accion VARCHAR(50) NOT NULL,
    saldo_anterior_origen DECIMAL(15,2),
    saldo_nuevo_origen DECIMAL(15,2),
    saldo_anterior_destino DECIMAL(15,2),
    saldo_nuevo_destino DECIMAL(15,2),
    fecha_auditoria DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    usuario_sistema VARCHAR(100),
    CONSTRAINT fk_auditoria_transaccion FOREIGN KEY (id_transaccion) REFERENCES TRANSACCION(id_transaccion),
    INDEX idx_transaccion (id_transaccion),
    INDEX idx_fecha_auditoria (fecha_auditoria)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =============================================
-- DATOS INICIALES (CATÁLOGOS)
-- =============================================

-- Insertar tipos de cuenta
INSERT INTO TIPO_CUENTA (nombre_tipo, descripcion, saldo_minimo, interes_anual) VALUES
('Cuenta de Ahorros', 'Cuenta básica de ahorros', 0.00, 2.50),
('Cuenta Corriente', 'Cuenta corriente sin intereses', 0.00, 0.00),
('Cuenta Nómina', 'Cuenta para recepción de nómina', 0.00, 1.50),
('Cuenta Premium', 'Cuenta con beneficios especiales', 5000.00, 3.50);

-- Insertar estados de cuenta
INSERT INTO ESTADO_CUENTA (nombre_estado, descripcion) VALUES
('Activa', 'Cuenta operando normalmente'),
('Bloqueada', 'Cuenta temporalmente bloqueada'),
('Cerrada', 'Cuenta cerrada definitivamente'),
('En Revisión', 'Cuenta bajo revisión administrativa');

-- Insertar tipos de transacción
INSERT INTO TIPO_TRANSACCION (nombre_tipo, descripcion, requiere_destino) VALUES
('Transferencia', 'Transferencia entre cuentas', TRUE),
('Depósito', 'Depósito en efectivo o cheque', FALSE),
('Retiro', 'Retiro de efectivo', FALSE),
('Pago de Servicios', 'Pago de servicios públicos', FALSE),
('Compra con Tarjeta', 'Compra realizada con tarjeta', FALSE);

-- Insertar estados de transacción
INSERT INTO ESTADO_TRANSACCION (nombre_estado, descripcion) VALUES
('Pendiente', 'Transacción pendiente de procesar'),
('Procesando', 'Transacción en proceso'),
('Completada', 'Transacción completada exitosamente'),
('Fallida', 'Transacción fallida'),
('Revertida', 'Transacción revertida');

-- Insertar tipos de tarjeta
INSERT INTO TIPO_TARJETA (nombre_tipo, descripcion, limite_diario) VALUES
('Débito', 'Tarjeta débito estándar', 5000.00),
('Crédito', 'Tarjeta de crédito', 10000.00),
('Prepago', 'Tarjeta prepago recargable', 2000.00);
