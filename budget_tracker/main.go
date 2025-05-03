package main

import (
	"fmt"
)

var budget int // Global variable to track the budget

func AddIncome(income int) {
	fmt.Println("Adding an income...")
	// Code to add an income
	budget += income
	fmt.Printf("New budget: %d\n", budget)
}

func addExp(expense int) {
	fmt.Println("Adding an expense...")
	// Code to add an expense
	budget -= expense
	fmt.Printf("New budget: %d\n", budget)
}
func viewExp() {
	fmt.Println("Viewing expenses...")
	// Code to view expenses
	fmt.Printf("Current budget: %d\n", budget)

}

func main() {
	var income int
	fmt.Println("Welcome to the Budget Tracker!")
	for {
		fmt.Println("Please select an option: 1. Add Income, 2. Add Expense, 3. View Expenses, 4. Exit program")
		option := 0
		fmt.Scan(&option)
		switch option {
		case 1:
			fmt.Println("Enter your income:")
			fmt.Scan(&income)
			AddIncome(income)
		case 2:
			fmt.Println("Adding an expense...")
			// Code to add an expense
			var expense int
			fmt.Println("Enter your expense:")
			fmt.Scan(&expense)
			addExp(expense)
		case 3:
			fmt.Println("Viewing budget...")
			viewExp()
		case 4:
			fmt.Println("Exiting the program...")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
