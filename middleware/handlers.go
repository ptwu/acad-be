package middleware

import (
	"acad-be/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type response struct {
	ID      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

const NumChengyu = 258

func createConnection() *sql.DB {
	var dbUser string = os.Getenv("DB_USER")
	var dbHost string = os.Getenv("DB_HOST")
	var dbPort int
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
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

// =======================================
// ROUTES
// =======================================
// /api/create-user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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

// /api/user/{id}?offset={int[-12, 12]}
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	timezoneOffsetStr := r.URL.Query().Get("offset")
	timezoneOffset, err := strconv.Atoi(timezoneOffsetStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := getUser(params["id"], timezoneOffset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// /api/switch-basis/{id}?is-traditional={true, false}
func SwitchBasis(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	istraditionalStr := r.URL.Query().Get("is-traditional")
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

// /api/review-card/{id}
func ReviewCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	err := markCurrentCardReviewed(params["id"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return	
	}

	res := response {
		Message: "Card successfully reviewed",
	}
	json.NewEncoder(w).Encode(res)
}

// =======================================
// HANDLERS
// =======================================
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

func getUser(id string, timezoneOffset int) (models.User, error) {
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
	
	if (err != nil) {
		return user, err
	}
	// calculate times from unix timestamps and set to UTC time
	lastTimeObj := time.Unix(user.LastLearned, 0).UTC()
	currentTimeObj := time.Now().UTC()

	// adjust time objects for the user's respective timezone offset
	lastTimeObj = lastTimeObj.Add(time.Hour * time.Duration(timezoneOffset))
	currentTimeObj = currentTimeObj.Add(time.Hour * time.Duration(timezoneOffset))

	lastDay := lastTimeObj.Day()
	currentDay := currentTimeObj.Day()
	var streak int64 = user.Streak
	var highestStreak int64 = user.HighestStreak
	var totalLearned int64 = user.TotalLearned
	if (lastDay != currentDay) {
		// we are on a new day, set new streak
		nextDayAfterLastDay := lastTimeObj.Add(time.Hour * 24).Day()
		if (currentDay != nextDayAfterLastDay) {
			// reset streak if the user didn't come back the previous day
			streak = 0
		}
		streak++;
		if (streak > highestStreak) {
			highestStreak = streak
		}
		if (totalLearned < NumChengyu - 1) {
			totalLearned++
		} else {
			// The user is done with all the chengyu in the dataset.
			return user, nil
		}
		var lastLearned int64 = time.Now().Unix()
		sqlStatement := `UPDATE users SET streak=$2, higheststreak=$3, totallearned=$4, lastlearned=$5 WHERE userid=$1`
		_, updateErr := db.Exec(sqlStatement, id, streak, highestStreak, totalLearned, lastLearned)

		if updateErr != nil {
				return user, errors.New("error while executing UPDATE query")
		}
		user.Streak = streak
		user.HighestStreak = highestStreak
		user.LastLearned = lastLearned
		user.TotalLearned = totalLearned
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

func markCurrentCardReviewed(id string) error {
	db := createConnection()
	defer db.Close()

	var reviewPtsStr string
	var totalLearnedStr string
	selectSqlStatement := `SELECT totallearned, reviewpoints FROM users WHERE userid=$1`

	row := db.QueryRow(selectSqlStatement, id)
	err := row.Scan(&totalLearnedStr, &reviewPtsStr)
	if (err != nil) {
		return errors.New("error while querying for this current user")
	}

	reviewPts, err := strconv.Atoi(reviewPtsStr)
	if err != nil {
		return errors.New("error while converting reviewPointsStr to int")
	}

	totalLearned, err := strconv.Atoi(totalLearnedStr)
	if err != nil {
		return errors.New("error while converting totalLearnedStr to int")
	}

	if (reviewPts < NumChengyu - 1 && reviewPts < totalLearned) {
		reviewPts = reviewPts + 1
	} else {
		return errors.New("cannot review card past max")
	}

	sqlStatement := `UPDATE users SET reviewpoints=$2 WHERE userid=$1`

	_, updateErr := db.Exec(sqlStatement, id, reviewPts)
	if updateErr != nil {
			return updateErr
	}

	return err
}


