document.addEventListener("DOMContentLoaded", () => {
  const userId = localStorage.getItem("user_id");

  if (!userId) {
    alert("You must be logged in.");
    window.location.href = "index.html"; // Redirect to login page
    return;
  }

  // Logout Button functionality
  const logoutBtn = document.getElementById("logoutBtn");
  if (logoutBtn) {
    logoutBtn.addEventListener("click", () => {
      localStorage.removeItem("user_id");
      alert("Logged out successfully!");
      window.location.href = "index.html"; // Redirect to login page
    });
  }

  // Dummy statistics update (optional — can remove if using real data only)
  const incomeStat = document.getElementById("stat-income");
  const expenseStat = document.getElementById("stat-expense");
  const balanceStat = document.getElementById("stat-balance");

  const income = 1250.50;
  const expense = 780.25;
  const balance = income - expense;

  if (incomeStat) incomeStat.textContent = `$${income.toFixed(2)}`;
  if (expenseStat) expenseStat.textContent = `$${expense.toFixed(2)}`;
  if (balanceStat) balanceStat.textContent = `$${balance.toFixed(2)}`;
});

// Fetch and update budget data from backend
async function fetchBudgetData() {
  try {
    const userId = localStorage.getItem("user_id");
    console.log("User ID from localStorage:", userId);

    if (!userId) {
      document.getElementById("incomeAmt").textContent = "N/A";
      document.getElementById("expenseAmt").textContent = "N/A";
      document.getElementById("balanceAmt").textContent = "N/A";
      return;
    }

    const res = await fetch("http://localhost:8080/budget", {
      method: "GET",
      headers: {
        "X-User-ID": userId
      }
    });

    if (!res.ok) throw new Error(`HTTP error: ${res.status}`);
    const data = await res.json();
    console.log("Fetch result:", data);

    const income = data.income || 0;
    const expense = data.expense || 0;
    const balance = income - expense;

    // Update cards
    document.getElementById("incomeAmt").textContent = `${income} tk`;
    document.getElementById("expenseAmt").textContent = `${expense} tk`;
    document.getElementById("balanceAmt").textContent = `${balance} tk`;

    // Draw chart with fetched values
    renderChart(income, expense, balance);

  } catch (err) {
    document.getElementById("incomeAmt").textContent = "Error";
    document.getElementById("expenseAmt").textContent = "Error";
    document.getElementById("balanceAmt").textContent = "Error";
    console.error("Fetch error:", err);
  }
}

window.addEventListener("DOMContentLoaded", fetchBudgetData);

let budgetChartInstance = null; // holds current Chart instance

function renderChart(income, expense, balance) {
  const ctx = document.getElementById("budgetBarChart").getContext("2d");

  // Destroy old chart if it exists
  if (budgetChartInstance) {
    budgetChartInstance.destroy();
  }

  // Create and store new chart
  budgetChartInstance = new Chart(ctx, {
    type: "bar",
    data: {
      labels: ["Income", "Expenses", "Balance"],
      datasets: [{
        label: "Amount (tk)",
        data: [income, expense, balance],
        backgroundColor: [
          "rgba(46, 204, 113, 0.7)", 
          "rgba(231, 76, 60, 0.7)",
          "rgba(52, 152, 219, 0.7)"
        ],
        borderColor: [
          "rgba(46, 204, 113, 1)",
          "rgba(231, 76, 60, 1)",
          "rgba(52, 152, 219, 1)"
        ],
        borderWidth: 2,
        borderRadius: 10
      }]
    },
    options: {
      responsive: true,
      plugins: { legend: { display: false } },
      scales: { y: { beginAtZero: true, ticks: { stepSize: 200 } } }
    }
  });
}
const calendarDays = document.getElementById('calendarDays');
const monthYear = document.getElementById('monthYear');
const prevBtn = document.getElementById('prevBtn');
const nextBtn = document.getElementById('nextBtn');

let currentDate = new Date();

function renderCalendar(date) {
  calendarDays.innerHTML = '';

  const year = date.getFullYear();
  const month = date.getMonth();

  // Display month and year
  const options = { month: 'long', year: 'numeric' };
  monthYear.textContent = date.toLocaleDateString(undefined, options);

  // First day of the month
  const firstDay = new Date(year, month, 1);
  // Last day of the month
  const lastDay = new Date(year, month + 1, 0);

  // Weekday of first day (0=Sun,...)
  const firstDayIndex = firstDay.getDay();
  // Number of days in month
  const daysInMonth = lastDay.getDate();

  // Add empty divs for days before first day
  for (let i = 0; i < firstDayIndex; i++) {
    const emptyDiv = document.createElement('div');
    emptyDiv.classList.add('empty');
    calendarDays.appendChild(emptyDiv);
  }

  // Add day numbers
  for (let day = 1; day <= daysInMonth; day++) {
    const dayDiv = document.createElement('div');
    dayDiv.textContent = day;

    // Highlight today
    const today = new Date();
    if (
      day === today.getDate() &&
      month === today.getMonth() &&
      year === today.getFullYear()
    ) {
      dayDiv.classList.add('today');
    }

    calendarDays.appendChild(dayDiv);
  }
}

// Navigation handlers
prevBtn.addEventListener('click', () => {
  currentDate.setMonth(currentDate.getMonth() - 1);
  renderCalendar(currentDate);
});

nextBtn.addEventListener('click', () => {
  currentDate.setMonth(currentDate.getMonth() + 1);
  renderCalendar(currentDate);
});

// Initial render
renderCalendar(currentDate);
