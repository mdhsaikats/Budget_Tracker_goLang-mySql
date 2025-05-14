document.getElementById("registration-form").addEventListener("submit", async function (e) {
    e.preventDefault();
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;
    const name = document.getElementById("userName").value;
    const email = document.getElementById("email").value;

    const res = await fetch("http://localhost:8080/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password, name, email })
    });

    if (res.ok) {
        alert("Registration successful. Please log in.");
        window.location.href = "index.html";
    } else {
        document.getElementById("error-message").textContent = "Error during registration.";
    }
});
