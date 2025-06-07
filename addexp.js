document.querySelector("form").addEventListener("submit", async function (e) {
    e.preventDefault();

    const expense = document.getElementById("expense").value.trim();
    const expenseName = document.getElementById("expenseName").value.trim();
    const user_id = localStorage.getItem("user_id");

    if (!expense || !expenseName || !user_id) {
        console.log("Please fill in all fields and ensure you are logged in.");
        return;
    }

    try {
        const res = await fetch("http://localhost:8080/add_expense", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                user_id: parseInt(user_id),
                amount: parseInt(expense),
                name: expenseName
            })
        });

        if (res.ok) {
            alert("Expense added.");
            document.getElementById("expense").value = "";
            document.getElementById("expenseName").value = "";
        } else {
            const errorText = await res.text();
            alert("Failed to add expense: " + errorText);
        }
    } catch (err) {
        alert("Error: " + err.message);
    }
});
