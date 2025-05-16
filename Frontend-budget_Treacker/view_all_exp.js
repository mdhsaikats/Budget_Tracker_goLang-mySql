async function fetchAllExpenses() {
    const list = document.getElementById("expenseList");
    const userId = localStorage.getItem("user_id");

    if (!userId) {
        list.innerHTML = "<li>User not logged in.</li>";
        return;
    }

    try {
        const res = await fetch("http://localhost:8080/all_expenses", {
            method: "GET",
            headers: {
                "X-User-ID": userId
            }
        });

        if (!res.ok) {
            throw new Error(`HTTP error: ${res.status}`);
        }

        const expenses = await res.json();

        list.innerHTML = "";

        if (!expenses.length) {
            list.innerHTML = "<li>No expenses found.</li>";
            return;
        }

        expenses.forEach(exp => {
            const li = document.createElement("li");
            li.textContent = `${exp.name}: Tk${exp.amount}`;
            list.appendChild(li);
        });

    } catch (err) {
        list.innerHTML = "<li>Failed to load expenses.</li>";
        console.error("Fetch error:", err);
    }
}

window.addEventListener("DOMContentLoaded", fetchAllExpenses);
