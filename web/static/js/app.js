const $message = document.getElementById("message");

async function fetchMessage() {
  $message.textContent = "Loading...";

  try {
    const response = await fetch();
    if (!response.ok) {
      throw new Error("Failed to fetch message");
    }

    const data = await response.json();
    $message.textContent = data.content;
  } catch (error) {
    $message.textContent = "You are perfect as you are.";
  }
}

fetchMessage();
