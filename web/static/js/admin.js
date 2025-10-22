let credentials = null;

async function authenticatedFetch(url, options = {}) {
  const response = await fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      Authorization: "Basic " + credentials,
    },
  });

  // Handle unauthorized
  if (response.status === 401) {
    throw new Error("Unauthorized");
  }

  if (!response.ok) {
    throw new Error(`Request failed: ${response.status}`);
  }

  return response;
}

document
  .getElementById("refresh")
  .addEventListener("click", loadPendingMessages);

document.getElementById("login-form").addEventListener("submit", async (e) => {
  e.preventDefault();

  const username = document.getElementById("username").value;
  const password = document.getElementById("password").value;
  const loginButton = document.getElementById("login-btn");
  const loginError = document.getElementById("login-error");

  loginButton.disabled = true;
  loginButton.textContent = "Logging in...";
  loginError.classList.add("hidden");

  try {
    credentials = btoa(username + ":" + password);
    const response = await authenticatedFetch("/api/admin/pending");

    document.getElementById("login-form").classList.add("hidden");
    document.getElementById("admin-app").classList.remove("hidden");

    const messages = await response.json();
    displayMessages(messages);
  } catch (error) {
    credentials = null;
    loginError.textContent =
      error.message === "Unauthorized"
        ? "Invalid username or password"
        : "Login failed. Please try again.";
    loginError.classList.remove("hidden");
    loginButton.disabled = false;
    loginButton.textContent = "Login";
  }
});

async function loadPendingMessages() {
  try {
    const response = await authenticatedFetch("/api/admin/pending");
    const messages = await response.json();
    displayMessages(messages);
  } catch (error) {
    console.error("Error loading messages:", error);
    document.getElementById("pending-list").innerHTML =
      '<div class="error">Failed to load messages</div>';
  }
}

function displayMessages(messages) {
  const container = document.getElementById("pending-list");

  if (messages.length === 0) {
    container.innerHTML = '<div class="empty">No pending messages</div>';
    return;
  }

  container.innerHTML = messages
    .map(
      (msg) => `
                <div class="message-card" id="message-${msg.id}">
                    <div class="message-content">${escapeHtml(
                      msg.content
                    )}</div>
                    <div class="message-date">
                        ${new Date(msg.created_at).toLocaleString()}
                    </div>
                    <div class="message-actions">
                        <button 
                            class="btn-approve"
                            onclick="approveMessage(${msg.id})">
                            ✓ Approve
                        </button>
                        <button 
                            class="btn-reject"
                            onclick="rejectMessage(${msg.id})">
                            ✗ Reject
                        </button>
                    </div>
                </div>
            `
    )
    .join("");
}

function escapeHtml(unsafe) {
  return unsafe
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}