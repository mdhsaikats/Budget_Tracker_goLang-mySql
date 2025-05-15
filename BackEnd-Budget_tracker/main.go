package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

var db *sql.DB
var budget = make(map[string]int)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Income struct {
	UserID int `json:"user_id"`
	Amount int `json:"amount"`
}

type Expense struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Amount int    `json:"amount"`
	Name   string `json:"name"`
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	json.NewDecoder(r.Body).Decode(&user)
	_, err := db.Exec("INSERT INTO users (username, password, name, email) VALUES (?, ?, ?, ?)", user.Username, user.Password, user.Name, user.Email)
	if err != nil {
		http.Error(w, "Error Registering User", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User Registered Successfully")
}

func signIn(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var auth Auth
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Received login:", auth.Username, auth.Password) // Debug

	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", auth.Username, auth.Password).Scan(&userID)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"user_id": userID})
}

func addIncome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed) //error
		return
	}
	var income Income
	err := json.NewDecoder(r.Body).Decode(&income)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO income (user_id, amount) VALUES (?, ?)", income.UserID, income.Amount)
	if err != nil {
		http.Error(w, "Error adding income", http.StatusInternalServerError)
		return
	}

	budget["income"] += income.Amount
	fmt.Fprintln(w, "Income added successfully")
}

func addExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed) //error
		return
	}
	var expense Expense
	json.NewDecoder(r.Body).Decode(&expense)

	_, err := db.Exec("INSERT INTO expenses (user_id, amount, name) VALUES (?, ?, ?)", expense.UserID, expense.Amount, expense.Name)
	if err != nil {
		http.Error(w, "Error adding expense", http.StatusInternalServerError)
		return
	}

	budget["expense"] += expense.Amount
	fmt.Fprintln(w, "Expense added successfully")
}

func allExpences(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed) //error
		return
	}
	rows, err := db.Query("SELECT id, user_id, amount, name FROM expenses")
	if err != nil {
		http.Error(w, "Error fetching expenses", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	expenses := []Expense{}
	for rows.Next() {
		var expense Expense
		err := rows.Scan(&expense.ID, &expense.UserID, &expense.Amount, &expense.Name)
		if err != nil {
			http.Error(w, "Error scanning expense", http.StatusInternalServerError)
			return
		}
		expenses = append(expenses, expense)
	}

	json.NewEncoder(w).Encode(expenses)

}

func getBudget(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed) //error
		return
	}
	json.NewEncoder(w).Encode(budget)
}

func resetBudget(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	userIDStr := r.URL.Query().Get("user_id")
	userID, _ := strconv.Atoi(userIDStr)

	_, err := db.Exec("DELETE FROM income WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Error resetting income: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("DELETE FROM expense WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Error resetting expenses: "+err.Error(), http.StatusInternalServerError)
		return
	}
	budget[strconv.Itoa(userID)] = 0
	fmt.Fprintln(w, "Budget reset successful.")
}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/budget_tracker_db")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	http.HandleFunc("/register", registerUser)
	http.HandleFunc("/signin", signIn)
	http.HandleFunc("/add_income", addIncome)
	http.HandleFunc("/add_expense", addExpense)
	http.HandleFunc("/all_expenses", allExpences)
	http.HandleFunc("/budget", getBudget)
	http.HandleFunc("/reset_budget", resetBudget)

	fmt.Println("Server started at :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
