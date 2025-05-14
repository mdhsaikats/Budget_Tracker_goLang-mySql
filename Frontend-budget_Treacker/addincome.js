document.addEventListener("DOMContentLoaded", () => {
    const form = document.querySelector("form");
    const userId = localStorage.getItem("user_id");

    if (!userId) {
        alert("You must be logged in to add income.");
        window.location.href = "index.html"; // redirect to login
        return;
    }

    form.addEventListener("submit", async (e) => {
        e.preventDefault();

        const incomeAmount = parseInt(document.getElementById("income").value);
        if (isNaN(incomeAmount) || incomeAmount <= 0) {
            alert("Please enter a valid income amount.");
            return;
        }

        try {
            const response = await fetch("http://localhost:8080/add_income", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    user_id: parseInt(userId),
                    amount: incomeAmount
                })
            });

            if (response.ok) {
                alert("Income added successfully.");
                document.getElementById("income").value = ""; // reset form
            } else {
                const errorText = await response.text();
                alert("Error adding income: " + errorText);
            }
        } catch (err) {
            alert("Request failed: " + err.message);
        }
    });
});
