export interface User {
  username: string;
  password: string;
  name: string;
  email: string;
}

export interface AuthCredentials {
  username: string;
  password: string;
}

export interface Income {
  user_id: number;
  amount: number;
}

export interface Expense {
  id: number;
  user_id: number;
  amount: number;
  name: string;
}

export interface BudgetSummary {
  income: number;
  expense: number;
}

export interface AuthUser {
  user_id: number;
  username: string;
}