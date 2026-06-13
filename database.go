package main

import (
	"database/sql"
	"fmt"
	"strings"
)

// DB wraps a sql.DB connection with application-specific queries.
type DB struct {
	conn *sql.DB
}

// NewDB creates a new DB connection using the given Config.
func NewDB(cfg *Config) (*DB, error) {
	conn, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(5)
	conn.SetMaxIdleConns(2)
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}
	return &DB{conn: conn}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}

// CreateTables creates the component and transaction tables if they don't exist.
func (db *DB) CreateTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS component (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(200) NOT NULL,
			model VARCHAR(200) NOT NULL,
			package VARCHAR(100) DEFAULT '',
			quantity INT NOT NULL DEFAULT 0,
			location VARCHAR(200) DEFAULT '',
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_name (name),
			INDEX idx_model (model),
			UNIQUE KEY uk_name_model_package (name, model, package)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS transaction (
			id INT AUTO_INCREMENT PRIMARY KEY,
			component_id INT NOT NULL,
			type ENUM('inbound','outbound') NOT NULL,
			quantity INT NOT NULL,
			operator VARCHAR(100) DEFAULT '',
			reason VARCHAR(500) DEFAULT '',
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_component_id (component_id),
			INDEX idx_type (type),
			INDEX idx_created_at (created_at),
			FOREIGN KEY (component_id) REFERENCES component(id) ON DELETE RESTRICT
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
	}
	for _, q := range queries {
		if _, err := db.conn.Exec(q); err != nil {
			return fmt.Errorf("create table: %w", err)
		}
	}
	return nil
}

// GetStats returns dashboard statistics.
func (db *DB) GetStats() (*Stats, error) {
	s := &Stats{}
	if err := db.conn.QueryRow("SELECT COUNT(*) FROM component").Scan(&s.TotalTypes); err != nil {
		return nil, err
	}
	var total sql.NullInt64
	if err := db.conn.QueryRow("SELECT SUM(quantity) FROM component").Scan(&total); err != nil {
		return nil, err
	}
	if total.Valid {
		s.TotalQuantity = int(total.Int64)
	}
	if err := db.conn.QueryRow("SELECT COUNT(*) FROM component WHERE quantity > 0 AND quantity <= 5").Scan(&s.LowStockCount); err != nil {
		return nil, err
	}
	if err := db.conn.QueryRow("SELECT COUNT(*) FROM component WHERE quantity = 0").Scan(&s.OutStockCount); err != nil {
		return nil, err
	}
	return s, nil
}

// ListComponents returns all components, optionally filtered by search keyword.
func (db *DB) ListComponents(search string) ([]Component, error) {
	q := "SELECT id, name, model, COALESCE(package,''), quantity, COALESCE(location,''), COALESCE(notes,''), created_at, updated_at FROM component"
	args := []interface{}{}
	if search != "" {
		q += " WHERE name LIKE ? OR model LIKE ? OR location LIKE ? OR package LIKE ?"
		pattern := "%" + search + "%"
		args = append(args, pattern, pattern, pattern, pattern)
	}
	q += " ORDER BY name"

	rows, err := db.conn.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []Component
	for rows.Next() {
		var c Component
		if err := rows.Scan(&c.ID, &c.Name, &c.Model, &c.Package, &c.Quantity, &c.Location, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		components = append(components, c)
	}
	return components, rows.Err()
}

// GetComponent returns a single component by ID.
func (db *DB) GetComponent(id int) (*Component, error) {
	var c Component
	err := db.conn.QueryRow(
		"SELECT id, name, model, COALESCE(package,''), quantity, COALESCE(location,''), COALESCE(notes,''), created_at, updated_at FROM component WHERE id = ?",
		id,
	).Scan(&c.ID, &c.Name, &c.Model, &c.Package, &c.Quantity, &c.Location, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// CreateComponent inserts a new component and returns its ID.
func (db *DB) CreateComponent(c *Component) (int, error) {
	result, err := db.conn.Exec(
		"INSERT INTO component (name, model, package, quantity, location, notes) VALUES (?, ?, ?, ?, ?, ?)",
		c.Name, c.Model, c.Package, c.Quantity, c.Location, c.Notes,
	)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}

// UpdateComponent updates a component's info (not quantity).
func (db *DB) UpdateComponent(c *Component) error {
	_, err := db.conn.Exec(
		"UPDATE component SET name=?, model=?, package=?, location=?, notes=? WHERE id=?",
		c.Name, c.Model, c.Package, c.Location, c.Notes, c.ID,
	)
	return err
}

// DeleteComponent deletes a component by ID.
func (db *DB) DeleteComponent(id int) error {
	_, err := db.conn.Exec("DELETE FROM component WHERE id = ?", id)
	return err
}

// HasTransactions checks if a component has any transactions.
func (db *DB) HasTransactions(componentID int) (bool, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM transaction WHERE component_id = ?", componentID).Scan(&count)
	return count > 0, err
}

// SearchComponents returns components matching a query (for JSON API).
func (db *DB) SearchComponents(q string) ([]Component, error) {
	pattern := "%" + q + "%"
	rows, err := db.conn.Query(
		"SELECT id, name, model, COALESCE(package,''), quantity, COALESCE(location,''), COALESCE(notes,''), created_at, updated_at FROM component WHERE name LIKE ? OR model LIKE ? ORDER BY name LIMIT 20",
		pattern, pattern,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []Component
	for rows.Next() {
		var c Component
		if err := rows.Scan(&c.ID, &c.Name, &c.Model, &c.Package, &c.Quantity, &c.Location, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		components = append(components, c)
	}
	return components, rows.Err()
}

// UpdateComponentQuantity adjusts a component's stock by delta.
func (db *DB) UpdateComponentQuantity(componentID int, delta int) error {
	_, err := db.conn.Exec(
		"UPDATE component SET quantity = quantity + ? WHERE id = ?",
		delta, componentID,
	)
	return err
}

// GetComponentQuantity returns the current stock for a component.
func (db *DB) GetComponentQuantity(componentID int) (int, error) {
	var qty int
	err := db.conn.QueryRow("SELECT quantity FROM component WHERE id = ?", componentID).Scan(&qty)
	return qty, err
}

// CreateTransaction inserts a new transaction record.
func (db *DB) CreateTransaction(t *Transaction) error {
	_, err := db.conn.Exec(
		"INSERT INTO transaction (component_id, type, quantity, operator, reason, notes) VALUES (?, ?, ?, ?, ?, ?)",
		t.ComponentID, t.Type, t.Quantity, t.Operator, t.Reason, t.Notes,
	)
	return err
}

// ListTransactions returns transactions with component info, filtered by the given criteria.
func (db *DB) ListTransactions(f TransactionFilter) ([]TransactionWithComponent, error) {
	q := `SELECT t.id, t.component_id, t.type, t.quantity, COALESCE(t.operator,''), COALESCE(t.reason,''), COALESCE(t.notes,''), t.created_at,
	      COALESCE(c.name,''), COALESCE(c.model,'')
	      FROM transaction t JOIN component c ON t.component_id = c.id WHERE 1=1`
	args := []interface{}{}

	if f.Type != "" {
		q += " AND t.type = ?"
		args = append(args, f.Type)
	}
	if f.ComponentID > 0 {
		q += " AND t.component_id = ?"
		args = append(args, f.ComponentID)
	}
	if f.DateFrom != "" {
		q += " AND t.created_at >= ?"
		args = append(args, f.DateFrom+" 00:00:00")
	}
	if f.DateTo != "" {
		q += " AND t.created_at <= ?"
		args = append(args, f.DateTo+" 23:59:59")
	}
	q += " ORDER BY t.created_at DESC"

	rows, err := db.conn.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []TransactionWithComponent
	for rows.Next() {
		var t TransactionWithComponent
		if err := rows.Scan(&t.ID, &t.ComponentID, &t.Type, &t.Quantity, &t.Operator, &t.Reason, &t.Notes, &t.CreatedAt, &t.CompName, &t.CompModel); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, rows.Err()
}

// GetComponentTransactions returns all transactions for a specific component.
func (db *DB) GetComponentTransactions(componentID int) ([]Transaction, error) {
	rows, err := db.conn.Query(
		"SELECT id, component_id, type, quantity, COALESCE(operator,''), COALESCE(reason,''), COALESCE(notes,''), created_at FROM transaction WHERE component_id = ? ORDER BY created_at DESC",
		componentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.ComponentID, &t.Type, &t.Quantity, &t.Operator, &t.Reason, &t.Notes, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, rows.Err()
}

// GetRecentTransactions returns the last N transactions with component info.
func (db *DB) GetRecentTransactions(n int) ([]TransactionWithComponent, error) {
	rows, err := db.conn.Query(
		`SELECT t.id, t.component_id, t.type, t.quantity, COALESCE(t.operator,''), COALESCE(t.reason,''), COALESCE(t.notes,''), t.created_at,
		 COALESCE(c.name,''), COALESCE(c.model,'')
		 FROM transaction t JOIN component c ON t.component_id = c.id
		 ORDER BY t.created_at DESC LIMIT ?`, n,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []TransactionWithComponent
	for rows.Next() {
		var t TransactionWithComponent
		if err := rows.Scan(&t.ID, &t.ComponentID, &t.Type, &t.Quantity, &t.Operator, &t.Reason, &t.Notes, &t.CreatedAt, &t.CompName, &t.CompModel); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, rows.Err()
}

// GetLowStockComponents returns components with quantity between 1 and 5.
func (db *DB) GetLowStockComponents() ([]Component, error) {
	rows, err := db.conn.Query(
		"SELECT id, name, model, COALESCE(package,''), quantity, COALESCE(location,''), COALESCE(notes,''), created_at, updated_at FROM component WHERE quantity > 0 AND quantity <= 5 ORDER BY quantity",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []Component
	for rows.Next() {
		var c Component
		if err := rows.Scan(&c.ID, &c.Name, &c.Model, &c.Package, &c.Quantity, &c.Location, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		components = append(components, c)
	}
	return components, rows.Err()
}

// GetOutStockComponents returns components with zero quantity.
func (db *DB) GetOutStockComponents() ([]Component, error) {
	rows, err := db.conn.Query(
		"SELECT id, name, model, COALESCE(package,''), quantity, COALESCE(location,''), COALESCE(notes,''), created_at, updated_at FROM component WHERE quantity = 0 ORDER BY name",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []Component
	for rows.Next() {
		var c Component
		if err := rows.Scan(&c.ID, &c.Name, &c.Model, &c.Package, &c.Quantity, &c.Location, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		components = append(components, c)
	}
	return components, rows.Err()
}

// FindComponentByModel finds a component by exact model match.
func (db *DB) FindComponentByModel(model string) (*Component, error) {
	var c Component
	err := db.conn.QueryRow(
		"SELECT id, name, model, COALESCE(package,''), quantity, COALESCE(location,''), COALESCE(notes,''), created_at, updated_at FROM component WHERE model = ?",
		model,
	).Scan(&c.ID, &c.Name, &c.Model, &c.Package, &c.Quantity, &c.Location, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// FindComponentByNamePackage finds a component by name + package match.
func (db *DB) FindComponentByNamePackage(name, pkg string) (*Component, error) {
	var c Component
	err := db.conn.QueryRow(
		"SELECT id, name, model, COALESCE(package,''), quantity, COALESCE(location,''), COALESCE(notes,''), created_at, updated_at FROM component WHERE name = ? AND package = ?",
		name, pkg,
	).Scan(&c.ID, &c.Name, &c.Model, &c.Package, &c.Quantity, &c.Location, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// IsDuplicateComponent checks if a component with the same name+model+package exists.
func (db *DB) IsDuplicateComponent(name, model, pkg string, excludeID int) (bool, error) {
	var count int
	q := "SELECT COUNT(*) FROM component WHERE name = ? AND model = ? AND package = ?"
	args := []interface{}{name, model, pkg}
	if excludeID > 0 {
		q += " AND id != ?"
		args = append(args, excludeID)
	}
	err := db.conn.QueryRow(q, args...).Scan(&count)
	return count > 0, err
}

// EscapeLike escapes special characters in LIKE pattern.
func EscapeLike(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}
