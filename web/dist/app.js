const state = {
  userId: 0
};
const navButtons = document.querySelectorAll(".nav-button");
const views = {
  login: document.querySelector("#view-login"),
  balance: document.querySelector("#view-balance"),
  transfer: document.querySelector("#view-transfer")
};
const loginForm = document.querySelector("#login-form");
const loginMessage = document.querySelector("#login-message");
const balanceValue = document.querySelector("#balance-value");
const refreshBalance = document.querySelector("#refresh-balance");
const transferForm = document.querySelector("#transfer-form");
const transferMessage = document.querySelector("#transfer-message");
function setView(view) {
  navButtons.forEach((button) => {
    button.classList.toggle("is-active", button.dataset.view === view);
  });
  Object.keys(views).forEach((key) => {
    const target = views[key];
    if (!target) return;
    target.classList.toggle("is-hidden", key !== view);
  });
}
async function request(path, options) {
  const response = await fetch(path, {
    headers: {
      "Content-Type": "application/json"
    },
    ...options
  });
  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || "Request failed");
  }
  return response.json();
}
function formatCurrency(amount) {
  return new Intl.NumberFormat("ja-JP", {
    style: "currency",
    currency: "JPY"
  }).format(amount);
}
navButtons.forEach((button) => {
  button.addEventListener("click", () => {
    const view = button.dataset.view;
    if (view) {
      setView(view);
    }
  });
});
loginForm?.addEventListener("submit", async (event) => {
  event.preventDefault();
  if (!loginForm) return;
  const formData = new FormData(loginForm);
  const email = String(formData.get("email") ?? "");
  const password = String(formData.get("password") ?? "");
  loginMessage?.classList.remove("is-error", "is-success");
  if (loginMessage) {
    loginMessage.textContent = "ログイン中...";
  }
  try {
    const data = await request("/api/login", {
      method: "POST",
      body: JSON.stringify({ email, password })
    });
    state.userId = data.user_id;
    if (loginMessage) {
      loginMessage.textContent = "ログインに成功しました。";
    }
    loginMessage?.classList.add("is-success");
    await refreshBalanceData();
    setView("balance");
  } catch (error) {
    if (loginMessage) {
      loginMessage.textContent = error.message;
    }
    loginMessage?.classList.add("is-error");
  }
});
async function refreshBalanceData() {
  if (!state.userId) {
    if (balanceValue) {
      balanceValue.textContent = "¥0";
    }
    return;
  }
  const data = await request(`/api/balance?user_id=${state.userId}`);
  if (balanceValue) {
    balanceValue.textContent = formatCurrency(data.balance);
  }
}
refreshBalance?.addEventListener("click", async () => {
  try {
    await refreshBalanceData();
  } catch (error) {
    alert(error.message);
  }
});
transferForm?.addEventListener("submit", async (event) => {
  event.preventDefault();
  if (!transferForm) return;
  const formData = new FormData(transferForm);
  const toAccountNumber = String(formData.get("toAccountNumber") ?? "");
  const amount = Number(formData.get("amount"));
  transferMessage?.classList.remove("is-error", "is-success");
  if (transferMessage) {
    transferMessage.textContent = "振り込み処理中...";
  }
  try {
    const data = await request("/api/transfer", {
      method: "POST",
      body: JSON.stringify({
        from_user_id: state.userId,
        to_account_number: toAccountNumber,
        amount
      })
    });
    if (transferMessage) {
      transferMessage.textContent = `振り込みが完了しました (ID: ${data.transfer_id})`;
    }
    transferMessage?.classList.add("is-success");
    transferForm.reset();
    await refreshBalanceData();
  } catch (error) {
    if (transferMessage) {
      transferMessage.textContent = error.message;
    }
    transferMessage?.classList.add("is-error");
  }
});
setView("login");
