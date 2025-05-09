package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var budget int     // Global variable to track the budget
var expenses []int // Slice to store expenses

func AddIncome(db *sql.DB, income int) {
	fmt.Println("Adding an income...")
	_, err := db.Exec("INSERT INTO income (amount) VALUES (?)", income)
	if err != nil {
		fmt.Println("Error adding income:", err)
		return
	}
	// Code to add an income
	budget += income
	fmt.Printf("New budget: %d\n", budget)
}

func addExp(db *sql.DB, expense int, expenseName string) {
	fmt.Println("Adding an expense...")
	// Code to add an expense
	_, err := db.Exec("INSERT INTO expense (amount , name) VALUES (?,?)", expense, expenseName)
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

func allExp() {
	fmt.Println("Viewing expenses...")
	// Code to view expenses
	if len(expenses) == 0 {
		fmt.Println("No expenses recorded.")
		return
	}
	fmt.Println("Expenses:")
	for i, expense := range expenses {
		fmt.Printf("%d: %d\n", i+1, expense)
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

func resetBudget() {
	fmt.Println("Resetting budget...")
	// Code to reset the budget
	budget = 0
	expenses = []int{}
	fmt.Println("Budget reset successfully.")
}

func main() {
	var income int
	dsn := "root:@tcp(127.0.0.1:3306)/budget_tracker_db" // Data Source Name for MySQL
	db, err := sql.Open("mysql", dsn)                    // Open a connection to the database
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close() // Ensure the database connection is closed when the program ends
	fmt.Println("Welcome to the Budget Tracker!")
	for {
		fmt.Printf("Please select an option:\n1. Add Income \n2. Add Expense \n3. View Expenses \n4. View All Exp\n5. Exit program\n")
		option := 0
		fmt.Scan(&option)
		switch option {
		case 1:
			fmt.Println("Enter your income:")
			fmt.Scan(&income)
			AddIncome(db, income)
		case 2:
			fmt.Println("Adding an expense...")
			// Code to add an expense
			var expense int
			var expenseName string
			fmt.Println("Enter the name of the expense:")
			fmt.Scan(&expenseName)
			fmt.Println("Enter your expense:")
			fmt.Scan(&expense)
			addExp(db, expense, expenseName)

		case 3:
			fmt.Println("Viewing budget...")
			viewExp()
		case 4:
			fmt.Println("Viewing all expenses...")
			allExp()
		case 5:
			fmt.Println("Exiting the program...")
			return
		case 6:
			resetBudget()
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
