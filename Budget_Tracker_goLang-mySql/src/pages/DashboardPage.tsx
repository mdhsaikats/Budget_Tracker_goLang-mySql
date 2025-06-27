import React, { useState, useEffect, useCallback } from 'react';
import BudgetSummary from '../components/dashboard/BudgetSummary';
import AddIncomeForm from '../components/dashboard/AddIncomeForm';
import AddExpenseForm from '../components/dashboard/AddExpenseForm';
import ExpenseList from '../components/dashboard/ExpenseList';
import ResetBudgetButton from '../components/dashboard/ResetBudgetButton';
import { budgetApi } from '../services/api';
import { useAuth } from '../context/AuthContext';

interface BudgetData {
  income: number;
  expense: number;
}

interface Expense {
  id: number;
  user_id: number;
  amount: number;
  name: string;
}

const DashboardPage: React.FC = () => {
  const { user } = useAuth();
  const [budgetData, setBudgetData] = useState<BudgetData>({ income: 0, expense: 0 });
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const isUserValid = user && typeof user.user_id === 'number' && !isNaN(user.user_id);

  const fetchData = useCallback(async () => {
    if (!isUserValid) return;
    try {
      setLoading(true);
      const [budgetResponse, expensesResponse] = await Promise.all([
        budgetApi.getBudget(user.user_id),
        budgetApi.getAllExpenses(user.user_id)
      ]);

      setBudgetData(budgetResponse);
      setExpenses(expensesResponse);
      setError('');
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Failed to load data');
      }
      console.error('Error fetching data:', err);
    } finally {
      setLoading(false);
    }
  }, [isUserValid, user]);

  useEffect(() => {
    if (isUserValid) {
      fetchData();
    }
  }, [fetchData, isUserValid]);

  const handleDataUpdate = () => {
    fetchData();
  };


  if (!isUserValid) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <h3 className="text-lg font-semibold text-red-900 mb-2">Please Sign In</h3>
          <p className="text-gray-600 mb-4">You must be signed in to view your dashboard.</p>
          <a
            href="/"
            className="bg-emerald-500 text-white px-4 py-2 rounded-lg hover:bg-emerald-600 transition-colors"
          >
            Go to Sign In
          </a>
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-500 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading your budget...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-xl p-6">
        <h3 className="text-lg font-semibold text-red-900 mb-2">Error Loading Data</h3>
        <p className="text-red-700 mb-4">{error}</p>
        <button
          onClick={fetchData}
          className="bg-red-500 text-white px-4 py-2 rounded-lg hover:bg-red-600 transition-colors"
        >
          Try Again
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Budget Dashboard</h1>
        <p className="text-gray-600">Monitor your financial health and track your spending</p>
      </div>

      <BudgetSummary income={budgetData.income} expense={budgetData.expense} />

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-1 space-y-6">
          <AddIncomeForm onIncomeAdded={handleDataUpdate} />
          <AddExpenseForm onExpenseAdded={handleDataUpdate} />
          <ResetBudgetButton onReset={handleDataUpdate} />
        </div>
        
        <div className="lg:col-span-2">
          <ExpenseList expenses={expenses} />
        </div>
      </div>
    </div>
  );
};

export default DashboardPage;