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
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-User-ID")
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
	w.Header().Set("Content-Type", "application/json")
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

	var userID int
	var username string
	err := db.QueryRow(
		"SELECT id, username FROM users WHERE username = ? AND password = ?",
		auth.Username, auth.Password,
	).Scan(&userID, &username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":  userID,
		"username": username,
	})
}

func addIncome(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var income Income
	if err := json.NewDecoder(r.Body).Decode(&income); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if income.UserID == 0 || income.Amount <= 0 {
		http.Error(w, "Invalid input: user_id and amount must be positive", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO income (user_id, amount) VALUES (?, ?)", income.UserID, income.Amount)
	if err != nil {
		http.Error(w, "Error adding income: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Optional: Use a per-user budget map if tracking is needed
	budget["income"] += income.Amount

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Income added successfully"})
}

func addExpense(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var expense Expense
	err := json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch total income
	var totalIncome int
	err = db.QueryRow("SELECT IFNULL(SUM(amount), 0) FROM income WHERE user_id = ?", expense.UserID).Scan(&totalIncome)
	if err != nil {
		http.Error(w, "Error fetching income", http.StatusInternalServerError)
		return
	}

	// Fetch total expense
	var totalExpense int
	err = db.QueryRow("SELECT IFNULL(SUM(amount), 0) FROM expenses WHERE user_id = ?", expense.UserID).Scan(&totalExpense)
	if err != nil {
		http.Error(w, "Error fetching expenses", http.StatusInternalServerError)
		return
	}

	// Check if adding this expense exceeds available income
	if totalExpense+expense.Amount > totalIncome {
		http.Error(w, "Insufficient balance to add this expense", http.StatusBadRequest)
		return
	}

	// Insert expense
	_, err = db.Exec("INSERT INTO expenses (user_id, amount, name) VALUES (?, ?, ?)", expense.UserID, expense.Amount, expense.Name)
	if err != nil {
		http.Error(w, "Error adding expense", http.StatusInternalServerError)
		return
	}

	budget["expense"] += expense.Amount

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Expense added successfully")
}

func allExpenses(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("X-User-ID")
	fmt.Println("Received X-User-ID:", userID)

	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	rows, err := db.Query("SELECT id, user_id, amount, name FROM expenses WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Error fetching expenses", http.StatusInternalServerError)
		fmt.Println("DB query error:", err)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

func getBudget(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	fmt.Println("Budget endpoint hit") // <-- Add this

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	UserID := r.Header.Get("X-User-ID")
	fmt.Println("Received X-User-ID header:", UserID) // <-- Add this

	if UserID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	// Optional: Try converting to int early and log
	var totalIncome int
	err := db.QueryRow("SELECT IFNULL(SUM(amount), 0) FROM income WHERE user_id = ?", UserID).Scan(&totalIncome)
	if err != nil {
		http.Error(w, "Error fetching income: "+err.Error(), http.StatusInternalServerError)
		fmt.Println("Income query error:", err)
		return
	}

	var totalExpense int
	err = db.QueryRow("SELECT IFNULL(SUM(amount), 0) FROM expenses WHERE user_id = ?", UserID).Scan(&totalExpense)
	if err != nil {
		http.Error(w, "Error fetching expenses: "+err.Error(), http.StatusInternalServerError)
		fmt.Println("Expense query error:", err)
		return
	}

	fmt.Println("Fetched income/expense:", totalIncome, totalExpense)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"income":  totalIncome,
		"expense": totalExpense,
	})
}

func resetBudget(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid or missing user_id", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM income WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Error resetting income: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("DELETE FROM expenses WHERE user_id = ?", userID) // Correct table name here
	if err != nil {
		http.Error(w, "Error resetting expenses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// If you want to reset budget totals for the user, you may want to track them differently.
	// For example, if you keep per-user budget summary in a map:
	delete(budget, strconv.Itoa(userID)) // or reset fields if budget is more complex

	fmt.Fprintln(w, "Budget reset successful.")
}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/budget_tacker_db")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	http.HandleFunc("/register", registerUser)
	http.HandleFunc("/signin", signIn)
	http.HandleFunc("/add_income", addIncome)
	http.HandleFunc("/add_expense", addExpense)
	http.HandleFunc("/all_expenses", allExpenses)
	http.HandleFunc("/budget", getBudget)
	http.HandleFunc("/reset_budget", resetBudget)

	fmt.Println("Server started at :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
