document.addEventListener("DOMContentLoaded", () => {
    const loginForm = document.getElementById("login-form");
    const errorMessage = document.getElementById("error-message");

    loginForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const username = document.getElementById("username").value;
        const password = document.getElementById("password").value;

        try {
            const response = await fetch("http://localhost:8080/signin", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ username, password }),
            });

            if (!response.ok) {
                const errorText = await response.text();
                errorMessage.textContent = errorText || "Login failed.";
                return;
            }

            const data = await response.json();
            console.log("Login successful. User ID:", data.user_id);

            // Optionally store user_id in localStorage or sessionStorage
            localStorage.setItem("user_id", data.user_id);

            // ✅ Redirect to main page (e.g., dashboard.html)
            window.location.href = "main.html";

        } catch (err) {
            errorMessage.textContent = "Error: " + err.message;
        }
    });
});


