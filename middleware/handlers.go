package middleware

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "acad-be/models"
    "log"
    "net/http"
    "os"
		"time"
    "strconv"
		"github.com/google/uuid"
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
			log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
			panic(err)
	}

	err = db.Ping()
	if err != nil {
			panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

// CreateUser create a user in the postgres db
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	insertID := createUser()

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
			log.Fatalf("Unable to get user. %v", err)
	}

	json.NewEncoder(w).Encode(user)
}

// HANDLERS
func createUser() string {
	db := createConnection()

	defer db.Close()

	newUuid := uuid.New().String()
	sqlStatement := `INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6, $7)`


	err := db.QueryRow(sqlStatement, newUuid, 0, 0, 0, 0, time.Now().Unix(), false)

	if err != nil {
			log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", newUuid)

	return newUuid
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
	case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return user, nil
	case nil:
			return user, nil
	default:
			log.Fatalf("Unable to scan the row. %v", err)
	}

	return user, err
}

func updateUser(id int64, user models.User) int64 {
	db := createConnection()

	defer db.Close()

	sqlStatement := `UPDATE users SET name=$2, location=$3, age=$4 WHERE userid=$1`

	res, err := db.Exec(sqlStatement, id, user.Name, user.Location, user.Age)

	if err != nil {
			log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
			log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}