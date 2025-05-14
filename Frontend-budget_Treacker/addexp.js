document.querySelector("form").addEventListener("submit", async function (e) {
    e.preventDefault();
    const expense = document.getElementById("expense").value;
    const expenseName = document.getElementById("expenseName").value;
    const user_id = localStorage.getItem("user_id");

    await fetch("http://localhost:8080/add_expense", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ user_id: parseInt(user_id), amount: parseInt(expense), name: expenseName })
    });

    alert("Expense added.");
});
