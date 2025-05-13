package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var budget int
var expenses []int

func AddIncome(db *sql.DB, userID, income int) {
	fmt.Println("Adding an income...")
	_, err := db.Exec("INSERT INTO income (user_id, amount) VALUES (?,?)", userID, income)
	if err != nil {
		fmt.Println("Error adding income:", err)
		return
	}
	// Code to add an income
	budget += income
	fmt.Printf("New budget: %d\n", budget)
}

func addExp(db *sql.DB, userID, expense int, expenseName string) {
	fmt.Println("Adding an expense...")
	// Code to add an expense
	_, err := db.Exec("INSERT INTO expense (user_id, amount , name) VALUES (?,?,?)", userID, expense, expenseName)
	if err != nil {
		fmt.Println("Error adding expense:", err)
		return
	}
	budget -= expense
	expenses = append(expenses, expense)
	if budget < 0 {
		fmt.Println("Warning: Your budget is negative!")
	} else {
		fmt.Println("Expense added successfully.")
	}
	fmt.Printf("New budget: %d\n", budget)
}

func allExp(db *sql.DB, userID int) {
	fmt.Println("Viewing expenses...")
	// Code to view expenses
	rows, err := db.Query("SELECT id,amount,name FROM expense where user_id = ?", userID)
	if err != nil {
		fmt.Println("Error retrieving expenses:", err)
		return
	}
	defer rows.Close()

	fmt.Println("Expenses:")
	for rows.Next() {
		var id, amount int
		var name string
		err := rows.Scan(&id, &amount, &name)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		fmt.Printf("ID: %d, Amount: %d, Name: %s\n", id, amount, name)
	}
}

func viewExp() {
	// Code to view expenses
	fmt.Printf("Current budget: %d\n", budget)
	if budget == 0 {
		fmt.Println("Insufficient budget.")
	} else if budget <= 100 {
		fmt.Println("Warning: Your budget is low!")
	}
}

func resetBudget(db *sql.DB, userID int) {
	fmt.Println("Resetting budget...")
	_, err := db.Exec("DELETE FROM income WHERE user_id = ?", userID)
	if err != nil {
		fmt.Println("Error clearing user income:", err)
	}
	_, err = db.Exec("DELETE FROM expense WHERE user_id = ?", userID)
	if err != nil {
		fmt.Println("Error clearing user expenses:", err)
	}
	budget = 0
	expenses = []int{}
	fmt.Println("Budget reset successfully.")
}

func signIn(db *sql.DB, username, password string) int {
	var id int
	err := db.QueryRow("SELECT id FROM auth WHERE username = ? AND password = ?", username, password).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Invalid username or password.")
		} else {
			fmt.Println("Error signing in:", err)
		}
		return -1
	}
	fmt.Println("Sign in successful!")
	return id
}

func register(db *sql.DB, username, password, user, mail string) {
	_, err := db.Exec("INSERT INTO auth (username, password, user,mail) VALUES (?, ?,?,?)", username, password, user, mail)
	if err != nil {
		fmt.Println("Error registering user:", err)
		return
	}
	fmt.Println("Registration successful! You can now sign in.")
}

func mainSec(db *sql.DB, userID int) {
	for {
		var income int
		fmt.Printf("Please select an option:\n1. Add Income \n2. Add Expense \n3. View Expenses \n4. View All Exp\n5. Exit program\n6. Reset Budget\n")
		option := 0
		fmt.Scan(&option)
		switch option {
		case 1:
			fmt.Println("Enter your income:")
			fmt.Scan(&income)
			AddIncome(db, userID, income)
		case 2:
			fmt.Println("Adding an expense...")
			// Code to add an expense
			var expense int
			var expenseName string
			fmt.Println("Enter the name of the expense:")
			fmt.Scan(&expenseName)
			fmt.Println("Enter your expense:")
			fmt.Scan(&expense)
			addExp(db, userID, expense, expenseName)
		case 3:
			fmt.Println("Viewing budget...")
			viewExp()
		case 4:
			fmt.Println("Viewing all expenses...")
			allExp(db, userID)
		case 5:
			fmt.Println("Exiting the program...")
			return
		case 6:
			resetBudget(db, userID)
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/budget_tracker_db" // Data Source Name for MySQL
	db, err := sql.Open("mysql", dsn)                    // Open a connection to the database
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close() // Ensure the database connection is closed when the program ends
	fmt.Println("Welcome to the Budget Tracker!")
	for {
		var choose string
		fmt.Println("Signing in or register you can write exit to exit the program")
		fmt.Println("Choose s for Sign Up r for Registration:")
		fmt.Scanln(&choose)
		switch choose {
		case "s":
			var username, password string
			fmt.Print("Enter username: ")
			fmt.Scanln(&username)
			fmt.Print("Enter password: ")
			fmt.Scanln(&password)
			userID := signIn(db, username, password)
			if userID != -1 {
				mainSec(db, userID)
			}
			signIn(db, username, password)
		case "r":
			var username, password, user, mail string
			fmt.Println("Enter new username:")
			fmt.Scanln(&username)
			fmt.Println("Enter new password:")
			fmt.Scanln(&password)
			fmt.Println("Enter your name:")
			fmt.Scanln(&user)
			fmt.Println("Enter your mail:")
			fmt.Scanln(&mail)
			register(db, username, password, user, mail)
		case "exit":
			fmt.Println("Exiting the program...")
			return
		}
	}
}
