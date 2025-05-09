package main

import (
	"fmt"
)

var budget int     // Global variable to track the budget
var expenses []int // Slice to store expenses

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
	if budget <= 100 {
		fmt.Println("Warning: Your budget is running low!")
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
	fmt.Println("Welcome to the Budget Tracker!")
	for {
		fmt.Printf("Please select an option:\n1. Add Income \n2. Add Expense \n3. View Expenses \n4. View All Exp\n5. Exit program\n")
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
