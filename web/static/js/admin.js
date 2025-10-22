let credentials = null;

document.getElementById("login-form").addEventListener("submit", async (e) => {
  e.preventDefault();

  const username = document.getElementById("username").value;
  const password = document.getElementById("password").value;
  const loginButton = document.getElementById("login-button");
  const loginError = document.getElementById("login-error");

  loginButton.disabled = true;
  loginButton.textContent = "Logging in...";
  loginError.classList.add("hidden");

  try {
    credentials = btoa(username + ":" + password);

    // Test credentials by fetching pending messages
    const response = await fetch("/api/admin/pending", {
      headers: {
        Authorization: `Basic ${credentials}`,
      },
    });

    if (response.status === 401) {
      throw new Error("Invalid credentials");
    }

    if (!response.ok) {
      throw new Error("Login failed");
    }

    document.getElementById("login-form").classList.add("hidden");
    document.getElementById("admin-app").classList.remove("hidden");

    // Load messages
    const messages = await response.json();
    console.log(messages);
    // displayMessages(messages);
  } catch (error) {
    credentials = null;
    loginError.textContent =
      error.message === "Invalid credentials"
        ? "Invalid username or password"
        : "Login failed. Please try again.";
    loginError.classList.remove("hidden");
    loginButton.disabled = false;
    loginButton.textContent = "Login";
  }
});
