type LoginResponse = {
  user_id: number;
};

type BalanceResponse = {
  user_id: number;
  balance: number;
};

type TransferResponse = {
  transfer_id: number;
};

const state = {
  userId: 0,
};

const navButtons = document.querySelectorAll<HTMLButtonElement>(".nav-button");
const views = {
  login: document.querySelector<HTMLElement>("#view-login"),
  balance: document.querySelector<HTMLElement>("#view-balance"),
  transfer: document.querySelector<HTMLElement>("#view-transfer"),
};

const loginForm = document.querySelector<HTMLFormElement>("#login-form");
const loginMessage = document.querySelector<HTMLParagraphElement>("#login-message");
const balanceValue = document.querySelector<HTMLParagraphElement>("#balance-value");
const refreshBalance = document.querySelector<HTMLButtonElement>("#refresh-balance");
const transferForm = document.querySelector<HTMLFormElement>("#transfer-form");
const transferMessage = document.querySelector<HTMLParagraphElement>("#transfer-message");

function setView(view: keyof typeof views) {
  navButtons.forEach((button) => {
    button.classList.toggle("is-active", button.dataset.view === view);
  });

  (Object.keys(views) as Array<keyof typeof views>).forEach((key) => {
    const target = views[key];
    if (!target) return;
    target.classList.toggle("is-hidden", key !== view);
  });
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    headers: {
      "Content-Type": "application/json",
    },
    ...options,
  });

  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || "Request failed");
  }

  return response.json() as Promise<T>;
}

function formatCurrency(amount: number): string {
  return new Intl.NumberFormat("ja-JP", {
    style: "currency",
    currency: "JPY",
  }).format(amount);
}

navButtons.forEach((button) => {
  button.addEventListener("click", () => {
    const view = button.dataset.view as keyof typeof views;
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
  loginMessage && (loginMessage.textContent = "ログイン中...");

  try {
    const data = await request<LoginResponse>("/api/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });

    state.userId = data.user_id;
    loginMessage && (loginMessage.textContent = "ログインに成功しました。");
    loginMessage?.classList.add("is-success");
    await refreshBalanceData();
    setView("balance");
  } catch (error) {
    loginMessage && (loginMessage.textContent = (error as Error).message);
    loginMessage?.classList.add("is-error");
  }
});

async function refreshBalanceData() {
  if (!state.userId) {
    balanceValue && (balanceValue.textContent = "¥0");
    return;
  }

  const data = await request<BalanceResponse>(`/api/balance?user_id=${state.userId}`);
  balanceValue && (balanceValue.textContent = formatCurrency(data.balance));
}

refreshBalance?.addEventListener("click", async () => {
  try {
    await refreshBalanceData();
  } catch (error) {
    alert((error as Error).message);
  }
});

transferForm?.addEventListener("submit", async (event) => {
  event.preventDefault();
  if (!transferForm) return;

  const formData = new FormData(transferForm);
  const toAccountNumber = String(formData.get("toAccountNumber") ?? "");
  const amount = Number(formData.get("amount"));

  transferMessage?.classList.remove("is-error", "is-success");
  transferMessage && (transferMessage.textContent = "振り込み処理中...");

  try {
    const data = await request<TransferResponse>("/api/transfer", {
      method: "POST",
      body: JSON.stringify({
        from_user_id: state.userId,
        to_account_number: toAccountNumber,
        amount,
      }),
    });

    transferMessage &&
      (transferMessage.textContent = `振り込みが完了しました (ID: ${data.transfer_id})`);
    transferMessage?.classList.add("is-success");
    transferForm.reset();
    await refreshBalanceData();
  } catch (error) {
    transferMessage && (transferMessage.textContent = (error as Error).message);
    transferMessage?.classList.add("is-error");
  }
});

setView("login");
