package database

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// User represents the database user record
type User struct {
	JID          string
	Name         string
	Registered   bool
	Limit        int
	Money        int
	Premium      bool
	RegisteredAt string
}

// DB is the global SQLite connection
var DB *sql.DB

// InitDB initializes SQLite database connection and sets up required tables
func InitDB(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Create users table
	query := `
	CREATE TABLE IF NOT EXISTS users (
		jid TEXT PRIMARY KEY,
		name TEXT,
		registered INTEGER DEFAULT 0,
		limit_count INTEGER DEFAULT 20,
		money INTEGER DEFAULT 0,
		premium INTEGER DEFAULT 0,
		registered_at TEXT
	);`
	_, err = DB.Exec(query)
	return err
}

// GetUser retrieves a user from database. Returns nil, nil if user is not found.
func GetUser(jid string) (*User, error) {
	row := DB.QueryRow("SELECT jid, name, registered, limit_count, money, premium, registered_at FROM users WHERE jid = ?", jid)
	var u User
	var reg, prem int
	var regAt sql.NullString
	err := row.Scan(&u.JID, &u.Name, &reg, &u.Limit, &u.Money, &prem, &regAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	u.Registered = reg == 1
	u.Premium = prem == 1
	if regAt.Valid {
		u.RegisteredAt = regAt.String
	}
	return &u, nil
}

// CreateUser inserts a new user record into database
func CreateUser(jid string, name string, limitDefault int) (*User, error) {
	u := &User{
		JID:          jid,
		Name:         name,
		Registered:   false,
		Limit:        limitDefault,
		Money:        0,
		Premium:      false,
		RegisteredAt: "",
	}
	_, err := DB.Exec("INSERT INTO users (jid, name, registered, limit_count, money, premium, registered_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		u.JID, u.Name, 0, u.Limit, u.Money, 0, u.RegisteredAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// UpdateUser saves changes to a user record
func UpdateUser(u *User) error {
	reg := 0
	if u.Registered {
		reg = 1
	}
	prem := 0
	if u.Premium {
		prem = 1
	}
	_, err := DB.Exec("UPDATE users SET name = ?, registered = ?, limit_count = ?, money = ?, premium = ?, registered_at = ? WHERE jid = ?",
		u.Name, reg, u.Limit, u.Money, prem, u.RegisteredAt, u.JID)
	return err
}

// GetTotalUsers returns total number of users in database
func GetTotalUsers() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

// GetRegisteredUsers returns total number of registered users
func GetRegisteredUsers() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE registered = 1").Scan(&count)
	return count, err
}
