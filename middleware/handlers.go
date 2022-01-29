package middleware

import (
		"errors"
    "database/sql"
    "encoding/json"
    "fmt"
    "acad-be/models"
    "net/http"
		"time"
		"os"
		"strconv"
    "github.com/gorilla/mux" 
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

type response struct {
	ID      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
			fmt.Print("error loading env")
			return nil
	}

	var dbUser string = os.Getenv("DB_USER")
	var dbHost string = os.Getenv("DB_HOST")
	var dbPort int
	dbPort, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		fmt.Print("error converting string into int")
		return nil
	}
	
	var dbPassword = os.Getenv("DB_PASS")
	var dbName = os.Getenv("DB_NAME")
	db, err := sql.Open("postgres", 
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName))

	if err != nil {
		fmt.Print("error opening postgres db")
		return nil
	}

	err = db.Ping()
	if err != nil {
		fmt.Print("error pinging db")
		return nil
	}
	return db
}

// ROUTES
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	insertID, err := createUser()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := response {
			ID:      insertID,
			Message: "user created successfully",
	}
	json.NewEncoder(w).Encode(res)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	user, err := getUser(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func SwitchBasis(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	var istraditionalStr string
	istraditionalStr = r.URL.Query().Get("is-traditional")
	istraditional, err := strconv.ParseBool(istraditionalStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = switchUserCharacterBasis(params["id"], istraditional)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return	
	}

	res := response {
		Message: "Basis successfully switched",
	}
	json.NewEncoder(w).Encode(res)
}

// HANDLERS
func createUser() (string, error) {
	db := createConnection()
	defer db.Close()

	sqlStatement := 
		`INSERT INTO users (streak, higheststreak, totallearned, reviewpoints, lastlearned, usestraditional) VALUES ($1, $2, $3, $4, $5, $6) RETURNING userid`
	var uuid string
	err := db.QueryRow(sqlStatement, 0, 0, 0, 0, time.Now().Unix(), false).Scan(&uuid)
	if err != nil {
		return "", errors.New("error when executing INSERT query")
	}
	
	return uuid, nil
}

func getUser(id string) (models.User, error) {
	db := createConnection()
	defer db.Close()

	var user models.User
	sqlStatement := `SELECT * FROM users WHERE userid=$1`

	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&user.ID, 
		&user.Streak, 
		&user.HighestStreak, 
		&user.TotalLearned, 
		&user.ReviewPoints, 
		&user.LastLearned,
		&user.UsesTraditional)

	switch err {
		case nil:
			return user, nil
		default:
			return user, err
	}

	return user, err
}

func switchUserCharacterBasis(id string, isTraditional bool) error {
	db := createConnection()
	defer db.Close()
	sqlStatement := `UPDATE users SET usestraditional=$2 WHERE userid=$1`

	_, err := db.Exec(sqlStatement, id, isTraditional)
	if err != nil {
			return errors.New("error while executing UPDATE query")
	}

	return err
}

