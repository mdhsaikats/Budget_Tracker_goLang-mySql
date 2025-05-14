document.getElementById("login-form").addEventListener("submit", async function (e) {
    e.preventDefault();
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    const res = await fetch("http://localhost:8080/signin", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password })
    });

    if (res.ok) {
        const data = await res.json();
        localStorage.setItem("user_id", data.user_id); // Save user ID for later
        window.location.href = "main_menu.html";
    } else {
        document.getElementById("error-message").textContent = "Invalid username or password.";
    }
});
