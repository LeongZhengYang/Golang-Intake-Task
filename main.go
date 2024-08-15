package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Request struct {
	InputNumber       string `json:"inputNumber"`
	AccumulatedNumber string `json:"accumulatedNumber"`
}

type Response struct {
	Matched     int `json:"matched"`
	Continuous  int `json:"continuous"`
	Permutation int `json:"permutation"`
}

func main() {
	log.Println("Starting server on http://localhost:8080")

	db, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/transactions_db")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to the database.")

	// Serve files from the current directory
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Handle the /store endpoint
	http.HandleFunc("/store", storeTransaction)

	// Start the HTTP server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func storeTransaction(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		log.Printf("Error decoding request: %v", err)
		return
	}

	db, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/transactions_db")
	if err != nil {
		http.Error(w, "Failed to connect to the database: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error connecting to database: %v", err)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO transactions (accumulated_number) VALUES (?)", req.AccumulatedNumber)
	if err != nil {
		http.Error(w, "Failed to insert data: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error inserting data into database: %v", err)
		return
	}

	matched, continuous, permutation := performMatching(req.InputNumber, req.AccumulatedNumber)

	resp := Response{
		Matched:     matched,
		Continuous:  continuous,
		Permutation: permutation,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func performMatching(inputNumber, accumulatedNumber string) (int, int, int) {
	matched := 0
	continuous := 0
	permutation := 1

	// Frequency map to count occurrences of each digit
	inputCount := make(map[rune]int)
	accumulatedCount := make(map[rune]int)

	// Count frequencies in input number and accumulated number
	for _, digit := range inputNumber {
		inputCount[digit]++
	}
	for _, digit := range accumulatedNumber {
		accumulatedCount[digit]++
	}

	// Calculate matched digits
	for digit, count := range inputCount {
		if accumulatedCount[digit] > 0 {
			matched += min(count, accumulatedCount[digit])
		}
	}

	// Calculate continuous digits
	maxContinuous := 0

	for i := 0; i < len(inputNumber); i++ {
		for j := 0; j < len(accumulatedNumber); j++ {
			k := 0
			for i+k < len(inputNumber) && j+k < len(accumulatedNumber) && inputNumber[i+k] == accumulatedNumber[j+k] {
				k++
			}
			if k >= 2 && k > maxContinuous {
				maxContinuous = k
			}
		}
	}

	if maxContinuous >= 2 {
		continuous = maxContinuous
	} else {
		continuous = 0
	}

	// Check permutation
	if len(inputNumber) != len(accumulatedNumber) {
		permutation = 0
	} else {
		for digit, count := range inputCount {
			if accumulatedCount[digit] != count {
				permutation = 0
				break
			}
		}
	}

	return matched, continuous, permutation
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
