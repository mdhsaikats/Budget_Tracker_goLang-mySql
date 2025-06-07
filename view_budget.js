async function fetchBudgetData() {
    try {
        const userId = localStorage.getItem("user_id");
        console.log("User ID from localStorage:", userId);

        if (!userId) {
            document.getElementById("budgetDetails").textContent = "User not logged in.";
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

        document.getElementById("budgetDetails").innerHTML = `
            <p>Total Income: Tk${income}</p>
            <p>Total Expenses: Tk${expense}</p>
            <p>Remaining Balance: Tk${balance}</p>
        `;
    } catch (err) {
        document.getElementById("budgetDetails").textContent = "Failed to load budget data.";
        console.error("Fetch error:", err);
    }
}
window.addEventListener("DOMContentLoaded", fetchBudgetData);
