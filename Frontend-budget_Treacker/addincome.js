document.querySelector("form").addEventListener("submit", async function (e) {
    e.preventDefault();
    const income = document.getElementById("income").value;
    const user_id = localStorage.getItem("user_id");

    await fetch("http://localhost:8080/add_income", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ user_id: parseInt(user_id), amount: parseInt(income) })
    });

    alert("Income added.");
});
