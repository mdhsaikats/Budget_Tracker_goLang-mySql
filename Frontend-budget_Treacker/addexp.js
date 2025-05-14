document.addEventListener("DOMContentLoaded", () => {
    const form = document.querySelector("form");
    const userId = localStorage.getItem("user_id");

    if (!userId) {
        alert("You must be logged in to add expenses.");
        window.location.href = "index.html"; // redirect to login
        return;
    }

    form.addEventListener("submit", async (e) => {
        e.preventDefault();

        const name = document.getElementById("expenseName").value.trim();
        const amount = parseInt(document.getElementById("expense").value);

        if (!name || isNaN(amount) || amount <= 0) {
            alert("Please enter valid expense name and amount.");
            return;
        }

        try {
            const response = await fetch("http://localhost:8080/add_expense", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    user_id: parseInt(userId),
                    name: name,
                    amount: amount
                })
            });

            if (response.ok) {
                alert("Expense added successfully.");
                form.reset();
            } else {
                const errorText = await response.text();
                alert("Error adding expense: " + errorText);
            }
        } catch (err) {
            alert("Request failed: " + err.message);
        }
    });
});
