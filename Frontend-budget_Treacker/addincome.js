document.querySelector("form").addEventListener("submit", async function (e) {
    e.preventDefault();

    const income = document.getElementById("income").value.trim();
    const user_id = localStorage.getItem("user_id");

    if (!user_id) {
        console.log("User not logged in. Redirecting to login page.");
        window.location.href = "index.html"; // redirect to login page
        return;
    }

    if (!income || isNaN(income) || parseInt(income) <= 0) {
        console.log("Please enter a valid income amount.");
        return;
    }

    try {
        const res = await fetch("http://localhost:8080/add_income", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ user_id: parseInt(user_id), amount: parseInt(income) })
        });

        if (res.ok) {
            console.log("Income added successfully.");
            document.getElementById("income").value = ""; // reset the input
        } else {
            const err = await res.text();
            console.error("Failed to add income:", err);
        }
    } catch (error) {
        console.error("Error adding income:", error);
    }
});
