package main

import (
	"database/sql"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func ConnectDatabase() error {
	DB, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	err = CreateUserTable(DB)
	if err != nil {
		return err
	}

	err = CreatePhoneTable(DB)
	if err != nil {
		return err
	}

	db = DB

	return nil
}

func CreateUserTable(db *sql.DB) error {
	userTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			login TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			name TEXT NOT NULL,
			age INTEGER NOT NULL
		);
	`
	_, err := db.Exec(userTable)
	return err
}

func CreatePhoneTable(db *sql.DB) error {
	phoneTable := `
		CREATE TABLE IF NOT EXISTS phones (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			phone TEXT NOT NULL,
			description TEXT,
			is_fax INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users (id)
		);
	`
	_, err := db.Exec(phoneTable)
	return err
}

func InsertUser(db *sql.DB, login, password, name, age string) (sql.Result, error) {
	return db.Exec("INSERT INTO users (login, password, name, age) VALUES (?, ?, ?, ?)", login, password, name, age)
}

func GetUserByName(db *sql.DB, name string) (User, error) {
	var user User
	row := db.QueryRow("SELECT id, name, age FROM users WHERE name = ?", name)
	err := row.Scan(&user.ID, &user.Name, &user.Age)
	return user, err
}

func InsertPhone(db *sql.DB, userID int, phone, description string, isFax bool) (sql.Result, error) {
	return db.Exec("INSERT INTO phones (user_id, phone, description, is_fax) VALUES (?, ?, ?, ?)", userID, phone, description, isFax)
}

func CountPhonesByPhoneNumber(db *sql.DB, phone string) (int, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM phones WHERE phone = ?", phone)
	err := row.Scan(&count)
	return count, err
}

func SearchPhonesByQuery(db *sql.DB, userID int, query string) ([]Phone, error) {
	rows, err := db.Query("SELECT id, phone, description, is_fax FROM phones WHERE user_id = ? AND phone LIKE ?", userID, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phones []Phone
	for rows.Next() {
		var phone Phone
		err := rows.Scan(&phone.ID, &phone.PhoneNumber, &phone.Description, &phone.IsFax)
		if err != nil {
			return nil, err
		}
		phones = append(phones, phone)
	}

	return phones, nil
}

func GetUserByLogin(db *sql.DB, login string) (User, error) {
	var user User
	row := db.QueryRow("SELECT id, password FROM users WHERE login = ?", login)
	err := row.Scan(&user.ID, &user.Password)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUserByID(db *sql.DB, userID int) (User, error) {
	var user User

	row := db.QueryRow("SELECT id, login, password, name, age FROM users WHERE id = ?", userID)
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Name, &user.Age)
	if err != nil {
		return user, err
	}

	return user, nil
}

func ComparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateToken(login string, userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		Login:  login,
		UserID: userID,
	})
	return token.SignedString(signingKey)
}

func CheckPhoneOwnership(db *sql.DB, phoneID, userID int) (bool, error) {
	row := db.QueryRow("SELECT user_id FROM phones WHERE id = ?", phoneID)
	var ownerID int
	err := row.Scan(&ownerID)
	if err != nil {
		return false, err
	}
	return ownerID == userID, nil
}

func UpdatePhoneData(db *sql.DB, newPhoneData Phone) error {
	_, err := db.Exec("UPDATE phones SET phone = ?, description = ?, is_fax = ? WHERE id = ?",
		newPhoneData.PhoneNumber, newPhoneData.Description, newPhoneData.IsFax, newPhoneData.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeletePhoneData(db *sql.DB, phoneID int) error {
	_, err := db.Exec("DELETE FROM phones WHERE id = ?", phoneID)
	if err != nil {
		return err
	}
	return nil
}
