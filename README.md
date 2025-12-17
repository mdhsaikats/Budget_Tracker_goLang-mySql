# Budget Tracker
Track income and expenses, view budget summaries, and ask an AI assistant questions about your ledger.

**Tech:** HTML, CSS, JavaScript (frontend) • Go (backend) • MySQL

**Repository structure**
- `index.html` — frontend UI
- `frontend/` — static assets (`css/`, `js/`)
- `go-backend/` — Go server and DB schema
	- `budget_tacker_db.sql` — SQL to create tables used by the app
	- `main.go` — server implementation

**Quickstart**
1. Install MySQL and create a database. Import schema:

	```sh
	mysql -u root -p < go-backend/budget_tacker_db.sql
	```

2. Configure backend environment. Create `go-backend/.env` (example):

	```env
	AI_API_KEY=YOUR_API_KEY
	AI_API_URL=https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent
	AI_MODEL=gemini-1.5-flash
	```

	- The server will also try OpenAI-compatible endpoints when `AI_API_URL` is not a Google endpoint.

3. Run the backend from `go-backend`:

	```sh
	cd go-backend
	go run .
	```

4. Open `index.html` in your browser (or serve it) and use the UI.

**API Endpoints (backend)**
- `POST /register` — register user (body: `username`, `password`, `name`, `email`)
- `POST /signin` — sign in (body: `username`, `password`) => returns `user_id`
- `POST /add_income` — add income (body: `user_id`, `amount`)
- `POST /add_expense` — add expense (body: `id`, `user_id`, `amount`, `name`)
- `GET  /all_expenses` — list expenses (set header `X-User-ID`)
- `GET  /budget` — get totals (set header `X-User-ID`)
- `POST /reset_budget` — reset (query `?user_id=<id>`)
- `POST /ai_predict` — ask the AI (body: `user_id`, `question`)

Notes:
- Some endpoints expect the `X-User-ID` header for authentication in this demo.
- CORS is enabled for simplicity (`Access-Control-Allow-Origin: *`).

**AI integration**
- The backend supports both OpenAI-style chat completions and Google's Gemini (Generative Language) endpoints.
- Set `AI_API_URL` and `AI_API_KEY` in `go-backend/.env`. The server will detect Google endpoints and try compatible request shapes; otherwise it uses OpenAI-style payloads.
- Example curl to test AI (OpenAI-style):

	```sh
	curl -X POST http://localhost:8080/ai_predict \
	-H "Content-Type: application/json" \
	-d '{"user_id":1,"question":"How can I save more this month?"}'
	```

**Development notes**
- The AI key in this repo (if present) should be replaced with your own key. Keep keys out of public repos.
- The `go-backend` uses `github.com/joho/godotenv` to load `.env` at startup.

If you want, I can add runnable scripts, a Dockerfile, or expand the README with deployment steps.
<h1 align="center">Budget Tracker</h1>
<h3 align="center">Track your budget regularly</h3>


<h3 align="left">Languages I used for this project:</h3>
<p align="left"> <a href="https://www.w3schools.com/css/" target="_blank" rel="noreferrer"> <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/css3/css3-original-wordmark.svg" alt="css3" width="40" height="40"/> </a> <a href="https://golang.org" target="_blank" rel="noreferrer"> <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/go/go-original.svg" alt="go" width="40" height="40"/> </a> <a href="https://www.w3.org/html/" target="_blank" rel="noreferrer"> <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/html5/html5-original-wordmark.svg" alt="html5" width="40" height="40"/> </a> <a href="https://developer.mozilla.org/en-US/docs/Web/JavaScript" target="_blank" rel="noreferrer"> <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/javascript/javascript-original.svg" alt="javascript" width="40" height="40"/> </a> <a href="https://www.mysql.com/" target="_blank" rel="noreferrer"> <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/mysql/mysql-original-wordmark.svg" alt="mysql" width="40" height="40"/> </a> <a href="https://spring.io/" target="_blank" rel="noreferrer"> <img src="https://www.vectorlogo.zone/logos/springio/springio-icon.svg" alt="spring" width="40" height="40"/> </a> </p>
