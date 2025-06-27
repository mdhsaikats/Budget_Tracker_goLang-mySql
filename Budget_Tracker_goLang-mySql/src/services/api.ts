const API_BASE_URL = 'http://localhost:8080';

export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'ApiError';
  }
}

async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new ApiError(response.status, errorText || 'Request failed');
  }

  const contentType = response.headers.get('content-type');
  if (contentType && contentType.includes('application/json')) {
    return response.json();
  }
  
  return response.text() as unknown as T;
}

export const authApi = {
  register: async (userData: {
    username: string;
    password: string;
    name: string;
    email: string;
  }) => {
    return apiRequest<string>('/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  },

  signIn: async (credentials: { username: string; password: string }) => {
    return apiRequest<{
      username: string; user_id: number 
}>('/signin', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
  },
};

export const budgetApi = {
  addIncome: async (userId: number, amount: number) => {
    return apiRequest<{ message: string }>('/add_income', {
      method: 'POST',
      body: JSON.stringify({
        user_id: userId,
        amount: amount,
      }),
    });
  },

  addExpense: async (userId: number, amount: number, name: string) => {
    return apiRequest<string>('/add_expense', {
      method: 'POST',
      body: JSON.stringify({
        user_id: userId,
        amount: amount,
        name: name,
      }),
    });
  },

  getAllExpenses: async (userId: number | undefined) => {
    if (typeof userId !== 'number') {
      throw new Error('User ID is undefined. Please sign in again.');
    }
    return apiRequest<Array<{
      id: number;
      user_id: number;
      amount: number;
      name: string;
    }>>('/all_expenses', {
      method: 'GET',
      headers: {
        'X-User-ID': userId.toString(),
      },
    });
  },

  getBudget: async (userId: number | undefined) => {
    if (typeof userId !== 'number') {
      throw new Error('User ID is undefined. Please sign in again.');
    }
    return apiRequest<{ income: number; expense: number }>('/budget', {
      method: 'GET',
      headers: {
        'X-User-ID': userId.toString(),
      },
    });
  },

  resetBudget: async (userId: number) => {
    return apiRequest<string>(`/reset_budget?user_id=${userId}`, {
      method: 'POST',
    });
  },
};