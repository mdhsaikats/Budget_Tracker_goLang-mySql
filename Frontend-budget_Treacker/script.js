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

  // Optional: Simulate statistics dynamically (You can later fetch real data via API or backend)
  const incomeStat = document.getElementById("stat-income");
  const expenseStat = document.getElementById("stat-expense");
  const balanceStat = document.getElementById("stat-balance");

  // Example dummy data — replace with real calculations later
  const income = 1250.50;
  const expense = 780.25;
  const balance = income - expense;

  if (incomeStat) incomeStat.textContent = `$${income.toFixed(2)}`;
  if (expenseStat) expenseStat.textContent = `$${expense.toFixed(2)}`;
  if (balanceStat) balanceStat.textContent = `$${balance.toFixed(2)}`;
});

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

        // Update respective card elements
        document.getElementById("incomeAmt").textContent = `${income} tk`;
        document.getElementById("expenseAmt").textContent = `${expense} tk`;
        document.getElementById("balanceAmt").textContent = `${balance} tk`;
    } catch (err) {
        document.getElementById("incomeAmt").textContent = "Error";
        document.getElementById("expenseAmt").textContent = "Error";
        document.getElementById("balanceAmt").textContent = "Error";
        console.error("Fetch error:", err);
    }
}

window.addEventListener("DOMContentLoaded", fetchBudgetData);
