import React from 'react';
import { Receipt, DollarSign } from 'lucide-react';

interface Expense {
  id: number;
  user_id: number;
  amount: number;
  name: string;
}

interface ExpenseListProps {
  expenses: Expense[];
}

const ExpenseList: React.FC<ExpenseListProps> = ({ expenses }) => {
  const totalExpenses = expenses.reduce((sum, expense) => sum + expense.amount, 0);

  return (
    <div className="bg-white rounded-xl shadow-sm border border-gray-100">
      <div className="p-6 border-b border-gray-100">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-blue-100 rounded-lg">
              <Receipt className="h-5 w-5 text-blue-600" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-gray-900">Recent Expenses</h3>
              <p className="text-sm text-gray-600">{expenses.length} transactions</p>
            </div>
          </div>
          <div className="text-right">
            <p className="text-sm text-gray-600">Total</p>
            <p className="text-lg font-bold text-gray-900">${totalExpenses.toLocaleString()}</p>
          </div>
        </div>
      </div>

      <div className="divide-y divide-gray-100 max-h-96 overflow-y-auto">
        {expenses.length === 0 ? (
          <div className="p-8 text-center">
            <Receipt className="h-12 w-12 text-gray-300 mx-auto mb-4" />
            <p className="text-gray-500">No expenses recorded yet</p>
            <p className="text-sm text-gray-400">Add your first expense to get started</p>
          </div>
        ) : (
          expenses.map((expense) => (
            <div key={expense.id} className="p-4 hover:bg-gray-50 transition-colors">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <div className="p-2 bg-gray-100 rounded-lg">
                    <DollarSign className="h-4 w-4 text-gray-600" />
                  </div>
                  <div>
                    <p className="font-medium text-gray-900">{expense.name}</p>
                    <p className="text-sm text-gray-500">ID: {expense.id}</p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="font-semibold text-red-600">-${expense.amount.toLocaleString()}</p>
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default ExpenseList;