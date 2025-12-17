const DEFAULT_API_BASE = 'http://localhost:8080';
const API_BASE_KEY = 'pulseBudgetApiBase';
const STORAGE_KEY = 'pulseBudgetUser';

let apiBase = localStorage.getItem(API_BASE_KEY) || DEFAULT_API_BASE;

const elements = {
	registerForm: document.getElementById('register-form'),
	signInForm: document.getElementById('signin-form'),
	incomeForm: document.getElementById('income-form'),
	expenseForm: document.getElementById('expense-form'),
	authSection: document.getElementById('auth-section'),
	dashboard: document.getElementById('dashboard'),
	usernameChip: document.getElementById('username-chip'),
	expenseList: document.getElementById('expense-list'),
	expenseEmpty: document.getElementById('expense-empty'),
	incomeValue: document.getElementById('budget-income'),
	expenseValue: document.getElementById('budget-expense'),
	balanceValue: document.getElementById('budget-balance'),
	toast: document.getElementById('toast'),
	refreshAllBtn: document.getElementById('refresh-all'),
	resetBtn: document.getElementById('reset-btn'),
	signOutBtn: document.getElementById('signout-btn'),
	apiBaseLabel: document.getElementById('api-base-label'),
	apiBaseInput: document.getElementById('api-base-input'),
	apiBaseForm: document.getElementById('api-base-form'),
	aiForm: document.getElementById('ai-form'),
	aiQuestionInput: document.getElementById('ai-question'),
	aiResponse: document.getElementById('ai-response'),
};

const state = {
	userId: null,
	username: null,
};

const currencyFormatter = new Intl.NumberFormat('en-US', {
	style: 'currency',
	currency: 'USD',
	maximumFractionDigits: 0,
});

function updateAIResponse(message, isError = false) {
	if (!elements.aiResponse) return;
	elements.aiResponse.textContent = message;
	elements.aiResponse.classList.toggle('text-rose-300', isError);
	elements.aiResponse.classList.toggle('text-slate-200', !isError);
}

function resetAIResponse() {
	const defaultMessage = state.userId
		? 'Ask a question above to get a tailored forecast.'
		: 'Sign in to ask the assistant about your budget.';
	updateAIResponse(defaultMessage, false);
}

function toggleAIFormLoading(isLoading) {
	if (!elements.aiForm) return;
	const button = elements.aiForm.querySelector('button[type="submit"]');
	if (!button) return;
	if (!button.dataset.defaultText) {
		button.dataset.defaultText = button.textContent.trim();
	}
	button.disabled = isLoading;
	button.classList.toggle('opacity-60', isLoading);
	button.textContent = isLoading ? 'Analyzing...' : button.dataset.defaultText;
}


function sanitizeApiBase(url) {
	const trimmed = (url || '').trim().replace(/\/$/, '');
	if (!trimmed) throw new Error('Base URL is required');
	try {
		const parsed = new URL(trimmed);
		return parsed.toString().replace(/\/$/, '');
	} catch (error) {
		throw new Error('Enter a valid URL (e.g., http://localhost:8080)');
	}
}

function updateApiBaseUI() {
	if (elements.apiBaseLabel) elements.apiBaseLabel.textContent = apiBase;
	if (elements.apiBaseInput && document.activeElement !== elements.apiBaseInput) {
		elements.apiBaseInput.value = apiBase;
	}
}

function setApiBase(url) {
	apiBase = sanitizeApiBase(url);
	localStorage.setItem(API_BASE_KEY, apiBase);
	updateApiBaseUI();
}

function showToast(message, type = 'success') {
	if (!elements.toast) return;
	elements.toast.textContent = message;
	elements.toast.className = `toast show ${type}`;
	clearTimeout(elements.toast.timeoutId);
	elements.toast.timeoutId = setTimeout(() => {
		elements.toast.classList.remove('show');
	}, 3600);
}

function persistSession() {
	if (state.userId) {
		localStorage.setItem(
			STORAGE_KEY,
			JSON.stringify({ userId: state.userId, username: state.username })
		);
	} else {
		localStorage.removeItem(STORAGE_KEY);
	}
}

function applySessionToUI() {
	const isLoggedIn = Boolean(state.userId);
	if (elements.authSection) {
		elements.authSection.classList.toggle('hidden', isLoggedIn);
	}
	if (elements.dashboard) {
		elements.dashboard.classList.toggle('hidden', !isLoggedIn);
	}
	if (elements.usernameChip) {
		elements.usernameChip.textContent = state.username || '—';
	}
}

function restoreSession() {
	try {
		const stored = localStorage.getItem(STORAGE_KEY);
		if (!stored) return;
		const parsed = JSON.parse(stored);
		if (parsed.userId) {
			state.userId = parsed.userId;
			state.username = parsed.username;
			applySessionToUI();
			resetAIResponse();
			refreshInsights();
		}
	} catch (error) {
		console.error('Unable to restore session', error);
	}
}

function handleApiBaseSubmit(event) {
	event.preventDefault();
	if (!elements.apiBaseInput) return;
	try {
		setApiBase(elements.apiBaseInput.value);
		showToast('API base saved');
		refreshInsights();
	} catch (error) {
		showToast(error.message, 'error');
	}
}

async function handleRegister(event) {
	event.preventDefault();
	const form = event.currentTarget;
	const payload = Object.fromEntries(new FormData(form).entries());
	try {
		const response = await fetch(`${apiBase}/register`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(payload),
		});
		const message = await response.text();
		if (!response.ok) throw new Error(message || 'Registration failed');
		form.reset();
		showToast('User registered successfully');
	} catch (error) {
		showToast(error.message, 'error');
	}
}

async function handleSignIn(event) {
	event.preventDefault();
	const form = event.currentTarget;
	const payload = Object.fromEntries(new FormData(form).entries());
	try {
		const response = await fetch(`${apiBase}/signin`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(payload),
		});
		if (!response.ok) {
			const message = await response.text();
			throw new Error(message || 'Invalid credentials');
		}
		const data = await response.json();
		state.userId = data.user_id;
		state.username = data.username;
		persistSession();
		applySessionToUI();
		resetAIResponse();
		form.reset();
		showToast('Signed in');
		refreshInsights();
	} catch (error) {
		showToast(error.message, 'error');
	}
}

async function handleIncome(event) {
	event.preventDefault();
	if (!state.userId) return showToast('Sign in first', 'error');
	const form = event.currentTarget;
	const amount = Number(form.amount.value);
	try {
		const response = await fetch(`${apiBase}/add_income`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ user_id: state.userId, amount }),
		});
		const message = await response.text();
		if (!response.ok) throw new Error(message || 'Income failed');
		form.reset();
		showToast('Income added');
		refreshInsights();
	} catch (error) {
		showToast(error.message, 'error');
	}
}

async function handleExpense(event) {
	event.preventDefault();
	if (!state.userId) return showToast('Sign in first', 'error');
	const form = event.currentTarget;
	const payload = {
		name: form.name.value,
		amount: Number(form.amount.value),
		user_id: state.userId,
	};
	try {
		const response = await fetch(`${apiBase}/add_expense`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(payload),
		});
		const message = await response.text();
		if (!response.ok) throw new Error(message || 'Expense failed');
		form.reset();
		showToast('Expense added');
		refreshInsights();
	} catch (error) {
		showToast(error.message, 'error');
	}
}

async function handleAIQuestion(event) {
	event.preventDefault();
	if (!state.userId) return showToast('Sign in first', 'error');
	if (!elements.aiQuestionInput) return;

	const question = elements.aiQuestionInput.value.trim();
	if (!question) {
		showToast('Enter a question for the assistant', 'error');
		return;
	}

	toggleAIFormLoading(true);
	updateAIResponse('Analyzing your ledger...', false);

	try {
		const response = await fetch(`${apiBase}/ai_predict`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ user_id: state.userId, question }),
		});

		const data = await response.json().catch(() => ({}));
		if (!response.ok) {
			throw new Error(data.error || data.message || 'Unable to fetch AI insight');
		}

		updateAIResponse(data.answer || 'No insight generated. Try rephrasing.', false);
	} catch (error) {
		updateAIResponse(error.message, true);
		showToast(error.message, 'error');
	} finally {
		toggleAIFormLoading(false);
	}
}

async function fetchBudget() {
	if (!state.userId) return;
	try {
		const response = await fetch(`${apiBase}/budget`, {
			headers: { 'X-User-ID': state.userId },
		});
		if (!response.ok) {
			const message = await response.text();
			throw new Error(message || 'Unable to fetch budget');
		}
		const { income = 0, expense = 0 } = await response.json();
		if (elements.incomeValue) elements.incomeValue.textContent = currencyFormatter.format(income);
		if (elements.expenseValue) elements.expenseValue.textContent = currencyFormatter.format(expense);
		if (elements.balanceValue) elements.balanceValue.textContent = currencyFormatter.format(income - expense);
	} catch (error) {
		showToast(error.message, 'error');
	}
}

async function fetchExpenses() {
	if (!state.userId) return;
	try {
		const response = await fetch(`${apiBase}/all_expenses`, {
			headers: { 'X-User-ID': state.userId },
		});
		if (!response.ok) {
			const message = await response.text();
			throw new Error(message || 'Unable to fetch expenses');
		}
		const expenses = await response.json();
		renderExpenses(Array.isArray(expenses) ? expenses : []);
	} catch (error) {
		showToast(error.message, 'error');
	}
}

function renderExpenses(expenses) {
	if (!elements.expenseList || !elements.expenseEmpty) return;
	elements.expenseList.innerHTML = '';
	if (!expenses.length) {
		elements.expenseEmpty.classList.remove('hidden');
		return;
	}
	elements.expenseEmpty.classList.add('hidden');
	expenses.forEach((expense) => {
		const item = document.createElement('li');
		item.innerHTML = `
			<div>
				<span class="label">${expense.name}</span>
				<span class="meta">#${expense.id} · User ${expense.user_id}</span>
			</div>
			<span class="font-semibold">${currencyFormatter.format(expense.amount)}</span>
		`;
		elements.expenseList.appendChild(item);
	});
}

async function handleReset() {
	if (!state.userId) return showToast('Sign in first', 'error');
	try {
		const response = await fetch(`${apiBase}/reset_budget?user_id=${state.userId}`, {
			method: 'POST',
		});
		const message = await response.text();
		if (!response.ok) throw new Error(message || 'Reset failed');
		showToast('Budget reset');
		refreshInsights();
	} catch (error) {
		showToast(error.message, 'error');
	}
}

function handleSignOut() {
	state.userId = null;
	state.username = null;
	persistSession();
	applySessionToUI();
	if (elements.expenseList) elements.expenseList.innerHTML = '';
	if (elements.expenseEmpty) elements.expenseEmpty.classList.remove('hidden');
	['incomeValue', 'expenseValue', 'balanceValue'].forEach((key) => {
		if (elements[key]) elements[key].textContent = '$0';
	});
	if (elements.aiQuestionInput) elements.aiQuestionInput.value = '';
	resetAIResponse();
	showToast('Signed out');
}

function refreshInsights() {
	fetchBudget();
	fetchExpenses();
}

function attachEvents() {
	updateApiBaseUI();
	elements.apiBaseForm && elements.apiBaseForm.addEventListener('submit', handleApiBaseSubmit);
	elements.registerForm && elements.registerForm.addEventListener('submit', handleRegister);
	elements.signInForm && elements.signInForm.addEventListener('submit', handleSignIn);
	elements.incomeForm && elements.incomeForm.addEventListener('submit', handleIncome);
	elements.expenseForm && elements.expenseForm.addEventListener('submit', handleExpense);
	elements.refreshAllBtn && elements.refreshAllBtn.addEventListener('click', refreshInsights);
	elements.resetBtn && elements.resetBtn.addEventListener('click', handleReset);
	elements.signOutBtn && elements.signOutBtn.addEventListener('click', handleSignOut);
	elements.aiForm && elements.aiForm.addEventListener('submit', handleAIQuestion);
}

function init() {
	attachEvents();
	applySessionToUI();
	resetAIResponse();
	restoreSession();
}

document.addEventListener('DOMContentLoaded', init);
