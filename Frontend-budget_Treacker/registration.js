document.addEventListener("DOMContentLoaded", () => {
    const loginForm = document.getElementById("login-form");
    const registrationForm = document.getElementById("registration-form");

    // Handle Login
    if (loginForm) {
        const errorMessage = document.getElementById("error-message");
        loginForm.addEventListener("submit", async (e) => {
            e.preventDefault();

            const username = document.getElementById("username").value;
            const password = document.getElementById("password").value;

            try {
                const response = await fetch("http://localhost:8080/signin", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({ username, password })
                });

                if (response.ok) {
                    const data = await response.json();
                    localStorage.setItem("user_id", data.user_id);
                    window.location.href = "dashboard.html";
                } else {
                    errorMessage.textContent = "Invalid username or password.";
                    errorMessage.style.color = "red";
                }
            } catch (error) {
                errorMessage.textContent = "Login error: " + error.message;
                errorMessage.style.color = "red";
            }
        });
    }

    // Handle Registration
    if (registrationForm) {
        const errorMessage = document.getElementById("error-message");
        registrationForm.addEventListener("submit", async (e) => {
            e.preventDefault();

            const username = document.getElementById("username").value;
            const password = document.getElementById("password").value;
            const name = document.getElementById("userName").value;
            const email = document.getElementById("email").value;

            try {
                const response = await fetch("http://localhost:8080/register", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({ username, password, name, email })
                });

                if (response.ok) {
                    alert("Registration successful! Please login.");
                    window.location.href = "index.html"; // Redirect to login
                } else {
                    const msg = await response.text();
                    errorMessage.textContent = "Registration failed: " + msg;
                    errorMessage.style.color = "red";
                }
            } catch (error) {
                errorMessage.textContent = "Registration error: " + error.message;
                errorMessage.style.color = "red";
            }
        });
    }
});
