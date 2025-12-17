package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-User-ID")
}

var db *sql.DB
var budget = make(map[string]int)

const defaultAIEndpoint = "https://api.openai.com/v1/chat/completions"
const defaultAIModel = "gpt-4.1-mini"

var (
	aiAPIKey   string
	aiEndpoint string
	aiModel    string
)

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

type AIQuestion struct {
	UserID   int    `json:"user_id"`
	Question string `json:"question"`
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

func aiPrediction(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload AIQuestion
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.UserID <= 0 || strings.TrimSpace(payload.Question) == "" {
		http.Error(w, "User ID and question are required", http.StatusBadRequest)
		return
	}

	snapshot, err := buildLedgerSnapshot(payload.UserID)
	if err != nil {
		http.Error(w, "Error building ledger snapshot", http.StatusInternalServerError)
		return
	}

	answer, aiErr := answerWithAI(payload.Question, snapshot)
	if aiErr != nil {
		fmt.Println("AI fallback triggered:", aiErr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"answer": answer,
	})
}

type LedgerSnapshot struct {
	TotalIncome  int
	TotalExpense int
	Expenses     []Expense
}

func buildLedgerSnapshot(userID int) (*LedgerSnapshot, error) {
	var totalIncome int
	var totalExpense int

	if err := db.QueryRow("SELECT IFNULL(SUM(amount), 0) FROM income WHERE user_id = ?", userID).Scan(&totalIncome); err != nil {
		return nil, err
	}

	if err := db.QueryRow("SELECT IFNULL(SUM(amount), 0) FROM expenses WHERE user_id = ?", userID).Scan(&totalExpense); err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT id, user_id, amount, name FROM expenses WHERE user_id = ? ORDER BY amount DESC LIMIT 20", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := []Expense{}
	for rows.Next() {
		var e Expense
		if err := rows.Scan(&e.ID, &e.UserID, &e.Amount, &e.Name); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}

	return &LedgerSnapshot{TotalIncome: totalIncome, TotalExpense: totalExpense, Expenses: expenses}, nil
}

func ledgerNarrative(snapshot *LedgerSnapshot) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Total income: %d. Total expense: %d. Balance: %d.\n", snapshot.TotalIncome, snapshot.TotalExpense, snapshot.TotalIncome-snapshot.TotalExpense))
	if len(snapshot.Expenses) == 0 {
		b.WriteString("No expenses recorded yet.")
		return b.String()
	}

	b.WriteString("Key expenses: ")
	limit := len(snapshot.Expenses)
	if limit > 5 {
		limit = 5
	}
	for i := 0; i < limit; i++ {
		expense := snapshot.Expenses[i]
		b.WriteString(fmt.Sprintf("%s (%d)", expense.Name, expense.Amount))
		if i < limit-1 {
			b.WriteString(", ")
		}
	}
	return b.String()
}

func answerWithAI(question string, snapshot *LedgerSnapshot) (string, error) {
	context := ledgerNarrative(snapshot)
	answer, err := callAIAPI(question, context)
	if err != nil {
		return fallbackAIResponse(question, snapshot), err
	}
	return answer, nil
}

func fallbackAIResponse(question string, snapshot *LedgerSnapshot) string {
	net := snapshot.TotalIncome - snapshot.TotalExpense
	savingsGuidance := "Expenses exceed income. Consider trimming high-impact categories or boosting income to reach equilibrium."
	if net > 0 {
		savingsGuidance = fmt.Sprintf("You currently clear %d more than you spend. Tucking away 70%% would save roughly %d each cycle.", net, int(float64(net)*0.7))
	}

	topExpense := "n/a"
	if len(snapshot.Expenses) > 0 {
		topExpense = fmt.Sprintf("%s at %d", snapshot.Expenses[0].Name, snapshot.Expenses[0].Amount)
	}

	return fmt.Sprintf("Based on your ledger (income %d vs expense %d), %s Biggest expense: %s. Question reviewed: %s", snapshot.TotalIncome, snapshot.TotalExpense, savingsGuidance, topExpense, question)
}

func callAIAPI(question, context string) (string, error) {
	if aiAPIKey == "" {
		return "", fmt.Errorf("AI API key not configured")
	}

	// If configured to use Google's Generative Language API (Gemini), use its
	// generate endpoint shape and parsing. Otherwise keep OpenAI-compatible flow.
	if strings.Contains(aiEndpoint, "generativelanguage.googleapis.com") {
		promptText := fmt.Sprintf("You are a concise financial planning assistant. Use the provided ledger context to answer savings questions with actionable steps.\n\nLedger data:\n%s\n\nQuestion:\n%s", context, question)

		// Try multiple request shapes accepted by different GenAI endpoint variants.
		tryPayloads := []map[string]interface{}{
			// messages-style (generateMessage / chat-like)
			{
				"messages": []map[string]interface{}{
					{"author": "user", "content": []map[string]string{{"type": "text", "text": promptText}}},
				},
				"temperature":     0.2,
				"maxOutputTokens": 256,
			},
			// snake_case variant
			{
				"messages": []map[string]interface{}{
					{"author": "user", "content": []map[string]string{{"type": "text", "text": promptText}}},
				},
				"temperature":       0.2,
				"max_output_tokens": 256,
			},
			// instances-style
			{
				"instances": []map[string]interface{}{
					{"content": []map[string]string{{"type": "text", "text": promptText}}},
				},
				"temperature":       0.2,
				"max_output_tokens": 256,
			},
		}

		urlBase := aiEndpoint
		// prefer API key as query param unless an OAuth token is detected
		addKeyToURL := func(u string) string {
			if aiAPIKey == "" || strings.Contains(u, "key=") {
				return u
			}
			if strings.Contains(u, "?") {
				return u + "&key=" + aiAPIKey
			}
			return u + "?key=" + aiAPIKey
		}

		client := &http.Client{Timeout: 20 * time.Second}

		for _, p := range tryPayloads {
			body, err := json.Marshal(p)
			if err != nil {
				continue
			}

			url := addKeyToURL(urlBase)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			if err != nil {
				continue
			}
			req.Header.Set("Content-Type", "application/json")

			// If OAuth-like token, set Authorization header
			if strings.HasPrefix(aiAPIKey, "ya29.") || strings.HasPrefix(aiAPIKey, "Bearer ") {
				token := aiAPIKey
				if strings.HasPrefix(token, "Bearer ") {
					token = strings.TrimPrefix(token, "Bearer ")
				}
				req.Header.Set("Authorization", "Bearer "+token)
			}

			resp, err := client.Do(req)
			if err != nil {
				continue
			}
			data, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				continue
			}

			if resp.StatusCode >= 300 {
				// If Google's error mentions unknown fields, try next payload; otherwise return error
				if strings.Contains(string(data), "Unknown name") || strings.Contains(string(data), "Invalid JSON payload") {
					continue
				}
				return "", fmt.Errorf("AI provider error: %s", string(data))
			}

			// parse response flexibly
			var parsed map[string]interface{}
			if err := json.Unmarshal(data, &parsed); err != nil {
				continue
			}

			// candidates -> content
			if cands, ok := parsed["candidates"].([]interface{}); ok && len(cands) > 0 {
				if first, ok := cands[0].(map[string]interface{}); ok {
					if content, ok := first["content"].(string); ok && strings.TrimSpace(content) != "" {
						return strings.TrimSpace(content), nil
					}
				}
			}

			// output or outputs
			if out, ok := parsed["output"].(string); ok && strings.TrimSpace(out) != "" {
				return strings.TrimSpace(out), nil
			}
			if outs, ok := parsed["outputs"].([]interface{}); ok && len(outs) > 0 {
				if o0, ok := outs[0].(map[string]interface{}); ok {
					if cont, ok := o0["content"].(string); ok && strings.TrimSpace(cont) != "" {
						return strings.TrimSpace(cont), nil
					}
				}
			}

			// messages -> content list
			if msgs, ok := parsed["messages"].([]interface{}); ok && len(msgs) > 0 {
				if m0, ok := msgs[0].(map[string]interface{}); ok {
					if contentArr, ok := m0["content"].([]interface{}); ok && len(contentArr) > 0 {
						if c0, ok := contentArr[0].(map[string]interface{}); ok {
							if text, ok := c0["text"].(string); ok && strings.TrimSpace(text) != "" {
								return strings.TrimSpace(text), nil
							}
						}
					}
				}
			}

			// If nothing found, continue to next payload
		}

		return "", fmt.Errorf("no supported response shape matched Gemini responses")
	}

	// Default: OpenAI-compatible chat completions
	payload := map[string]interface{}{
		"model": aiModel,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a concise financial planning assistant. Use the provided ledger context to answer savings questions with actionable steps.",
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("Ledger data:\n%s\n\nQuestion:\n%s", context, question),
			},
		},
		"temperature": 0.2,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, aiEndpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+aiAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("AI provider error: %s", string(data))
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", err
	}

	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("AI provider returned no choices")
	}

	return strings.TrimSpace(parsed.Choices[0].Message.Content), nil
}

func initAIConfig() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warning: unable to load .env -", err)
	}
	aiAPIKey = os.Getenv("AI_API_KEY")
	aiEndpoint = os.Getenv("AI_API_URL")
	if aiEndpoint == "" {
		aiEndpoint = defaultAIEndpoint
	}
	aiModel = os.Getenv("AI_MODEL")
	if aiModel == "" {
		aiModel = defaultAIModel
	}
}

func main() {
	initAIConfig()
	var err error
	db, err = sql.Open("mysql", "root:29112003@tcp(127.0.0.1:3306)/budget_tacker_db")
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
	http.HandleFunc("/ai_predict", aiPrediction)

	fmt.Println("Server started at :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
