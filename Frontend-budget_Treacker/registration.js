document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("registration-form");

    form.addEventListener("submit", async (e) => {
        e.preventDefault();

        const username = document.getElementById("username").value;
        const password = document.getElementById("password").value;
        const name = document.getElementById("userName").value;
        const email = document.getElementById("email").value;
        const errorMessage = document.getElementById("error-message");

        try {
            const response = await fetch("http://localhost:8080/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({ username, password, name, email })
            });

            if (response.ok) {
                errorMessage.style.color = "green";
                errorMessage.textContent = "Registration successful!";
                // Optionally redirect after a second
                setTimeout(() => {
                    window.location.href = "index.html";
                }, 1000);
            } else {
                const errText = await response.text();
                errorMessage.style.color = "red";
                errorMessage.textContent = "Registration failed: " + errText;
            }
        } catch (error) {
            errorMessage.style.color = "red";
            errorMessage.textContent = "Error: " + error.message;
        }
    });
});

