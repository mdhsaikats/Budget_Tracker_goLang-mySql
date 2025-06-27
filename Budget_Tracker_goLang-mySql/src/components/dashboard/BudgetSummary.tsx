import React from 'react';
import { TrendingUp, TrendingDown, DollarSign, PieChart } from 'lucide-react';

interface BudgetSummaryProps {
  income: number;
  expense: number;
}

const BudgetSummary: React.FC<BudgetSummaryProps> = ({ income, expense }) => {
  const balance = income - expense;
  const expensePercentage = income > 0 ? (expense / income) * 100 : 0;

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
        <div className="flex items-center justify-between mb-4">
          <div className="p-2 bg-emerald-100 rounded-lg">
            <TrendingUp className="h-6 w-6 text-emerald-600" />
          </div>
          <span className="text-sm font-medium text-emerald-600">+{expensePercentage > 0 ? (100 - expensePercentage).toFixed(1) : '100'}%</span>
        </div>
        <h3 className="text-2xl font-bold text-gray-900 mb-1">
          ${income.toLocaleString()}
        </h3>
        <p className="text-gray-600 text-sm">Total Income</p>
      </div>

      <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
        <div className="flex items-center justify-between mb-4">
          <div className="p-2 bg-red-100 rounded-lg">
            <TrendingDown className="h-6 w-6 text-red-600" />
          </div>
          <span className="text-sm font-medium text-red-600">{expensePercentage.toFixed(1)}%</span>
        </div>
        <h3 className="text-2xl font-bold text-gray-900 mb-1">
          ${expense.toLocaleString()}
        </h3>
        <p className="text-gray-600 text-sm">Total Expenses</p>
      </div>

      <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
        <div className="flex items-center justify-between mb-4">
          <div className={`p-2 rounded-lg ${balance >= 0 ? 'bg-blue-100' : 'bg-orange-100'}`}>
            <DollarSign className={`h-6 w-6 ${balance >= 0 ? 'text-blue-600' : 'text-orange-600'}`} />
          </div>
          <span className={`text-sm font-medium ${balance >= 0 ? 'text-blue-600' : 'text-orange-600'}`}>
            {balance >= 0 ? 'Surplus' : 'Deficit'}
          </span>
        </div>
        <h3 className={`text-2xl font-bold mb-1 ${balance >= 0 ? 'text-gray-900' : 'text-red-600'}`}>
          ${Math.abs(balance).toLocaleString()}
        </h3>
        <p className="text-gray-600 text-sm">Current Balance</p>
      </div>

      <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
        <div className="flex items-center justify-between mb-4">
          <div className="p-2 bg-purple-100 rounded-lg">
            <PieChart className="h-6 w-6 text-purple-600" />
          </div>
        </div>
        <div className="space-y-2">
          <div className="flex justify-between text-sm">
            <span className="text-gray-600">Spending Rate</span>
            <span className="font-medium">{expensePercentage.toFixed(1)}%</span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div 
              className={`h-2 rounded-full transition-all duration-500 ${
                expensePercentage > 80 ? 'bg-red-500' : 
                expensePercentage > 60 ? 'bg-orange-500' : 'bg-emerald-500'
              }`}
              style={{ width: `${Math.min(expensePercentage, 100)}%` }}
            ></div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default BudgetSummary;