-- Migraciones SQL normalizadas
CREATE TABLE customers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

CREATE TABLE accounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner TEXT NOT NULL,
    balance REAL NOT NULL,
    customer_id INTEGER,
    FOREIGN KEY(customer_id) REFERENCES customers(id)
);

CREATE TABLE transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_account INTEGER,
    to_account INTEGER,
    amount REAL NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(from_account) REFERENCES accounts(id),
    FOREIGN KEY(to_account) REFERENCES accounts(id)
);

-- Consultas separadas
-- Insertar cuenta: INSERT INTO accounts (owner, balance, customer_id) VALUES (?, ?, ?);
-- Transferencia: UPDATE accounts SET balance = balance - ? WHERE id = ?;
-- Registro de transacción: INSERT INTO transactions (from_account, to_account, amount) VALUES (?, ?, ?);
