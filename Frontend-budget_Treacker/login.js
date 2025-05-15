document.getElementById("login-form").addEventListener("submit", async function (e) {
    e.preventDefault();
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

        try {
            const response = await fetch("http://localhost:8080/signin", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    username,
                    password
                })
            });

            if (response.ok) {
                const data = await response.json();
                const userId = data.user_id;

                // Save user_id in localStorage for later use
                localStorage.setItem("user_id", userId);

                // Redirect to dashboard or home page
                window.location.href = "dashboard.html";
            } else {
                const errorText = await response.text();
                errorMessage.textContent = "Login failed: " + errorText;
                errorMessage.style.color = "red";
            }
        } catch (error) {
            console.error("Error during login:", error);
            errorMessage.textContent = "Something went wrong. Please try again.";
            errorMessage.style.color = "red";
        }
    });
});


