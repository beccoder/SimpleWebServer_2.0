package models

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
}

type Phone struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	PhoneNumber string `json:"phone"`
	Description string `json:"description"`
	IsFax       bool   `json:"is_fax"`
}

type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JWTClaims struct {
	Login  string `json:"login"`
	UserID int    `json:"user_id"`
	jwt.StandardClaims
}

var DB *sql.DB
var signingKey = []byte("secret")

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	err = CreateTables(db)
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func CreateTables(db *sql.DB) error {
	userTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			login TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			name TEXT NOT NULL,
			age INTEGER NOT NULL
		);
	`
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

	_, err := db.Exec(userTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(phoneTable)
	if err != nil {
		return err
	}

	return nil
}

func RegisterUser(newUser User) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO users (login, password, name, age) VALUES (?, ?, ?, ?)")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newUser.Login, string(hashedPassword), newUser.Name, newUser.Age)
	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func AuthenticateUser(loginData LoginData, c *gin.Context) (bool, error) {
	var user User
	err := DB.QueryRow("SELECT id, password FROM users WHERE login = ?", loginData.Login).Scan(&user.ID, &user.Password)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		return false, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		Login:  loginData.Login,
		UserID: user.ID,
	})
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return false, err
	}

	cookie := &http.Cookie{
		Name:  "SESSTOKEN",
		Value: tokenString,
		Path:  "/",
	}
	http.SetCookie(c.Writer, cookie)

	return true, nil
}

func GetUserHandler(name string) (User, error) {
	var user User
	err := DB.QueryRow("SELECT id, name, age FROM users WHERE name = ?", name).Scan(&user.ID, &user.Name, &user.Age)
	return user, err
}

func CountPhoneNumber(phoneNumber string) (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM phones WHERE phone = ?", phoneNumber).Scan(&count)
	return count, err
}

func AddPhoneNumber(phone Phone) (bool, error) {
	_, err := DB.Exec("INSERT INTO phones (user_id, phone, description, is_fax) VALUES (?, ?, ?, ?)", phone.UserID, phone.PhoneNumber, phone.Description, phone.IsFax)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func SearchPhoneHandler(userId int, query string) ([]Phone, error) {
	var phones []Phone
	rows, err := DB.Query("SELECT id, phone, description, is_fax FROM phones WHERE user_id = ? AND phone LIKE ?", userId, "%"+query+"%")
	if err != nil {
		return phones, err
	}
	defer rows.Close()

	for rows.Next() {
		var phone Phone
		err := rows.Scan(&phone.ID, &phone.PhoneNumber, &phone.Description, &phone.IsFax)
		if err != nil {
			return phones, err
		}
		phones = append(phones, phone)
	}
	return phones, nil
}

func CheckPhoneOwnership(userId int, phoneDataId int) (bool, error) {
	var ownerId int
	err := DB.QueryRow("SELECT user_id FROM phones WHERE id = ?", phoneDataId).Scan(&ownerId)
	if err != nil {
		return false, err
	}
	if userId != ownerId {
		return false, err
	}
	return true, nil
}

func UpdatePhoneData(newPhoneData Phone) (bool, error) {
	_, err := DB.Exec("UPDATE phones SET phone = ?, description = ?, is_fax = ? WHERE id = ?", newPhoneData.PhoneNumber, newPhoneData.Description, newPhoneData.IsFax, newPhoneData.ID)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func DeletePhoneData(phoneId int) (bool, error) {
	_, err := DB.Exec("DELETE FROM phones WHERE id = ?", phoneId)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func CheckJWTTokens(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return 0, err
	}
	var userId int
	err = DB.QueryRow("SELECT id FROM users WHERE login = ?", claims.Login).Scan(&userId)
	if err != nil {
		return 0, err
	} else if userId != claims.UserID {
		return -1, err
	} else {
		return userId, err
	}
}
