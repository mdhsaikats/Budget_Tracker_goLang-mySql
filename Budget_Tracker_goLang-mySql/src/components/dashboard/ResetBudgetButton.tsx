import React, { useState } from 'react';
import { RotateCcw, AlertTriangle } from 'lucide-react';
import { budgetApi } from '../../services/api';
import { useAuth } from '../../context/AuthContext';

interface ResetBudgetButtonProps {
  onReset: () => void;
}

const ResetBudgetButton: React.FC<ResetBudgetButtonProps> = ({ onReset }) => {
  const { user } = useAuth();
  const [showConfirm, setShowConfirm] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleReset = async () => {
    if (!user) return;

    setLoading(true);
    setError('');

    try {
      await budgetApi.resetBudget(user.user_id);
      onReset();
      setShowConfirm(false);
    } catch (err: any) {
      setError(err.message || 'Failed to reset budget');
    } finally {
      setLoading(false);
    }
  };

  if (showConfirm) {
    return (
      <div className="bg-white rounded-xl p-6 shadow-sm border border-red-200">
        <div className="flex items-center space-x-3 mb-4">
          <div className="p-2 bg-red-100 rounded-lg">
            <AlertTriangle className="h-5 w-5 text-red-600" />
          </div>
          <h3 className="text-lg font-semibold text-red-900">Confirm Reset</h3>
        </div>

        <div className="mb-6">
          <p className="text-gray-700 mb-2">
            This will permanently delete all your income and expense records.
          </p>
          <p className="text-sm text-red-600 font-medium">
            This action cannot be undone.
          </p>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-red-600 text-sm">{error}</p>
          </div>
        )}

        <div className="flex space-x-3">
          <button
            onClick={handleReset}
            disabled={loading}
            className="flex-1 bg-red-500 text-white py-2 px-4 rounded-lg font-semibold hover:bg-red-600 focus:ring-2 focus:ring-red-500 focus:ring-offset-2 disabled:opacity-50 transition-colors"
          >
            {loading ? (
              <div className="flex items-center justify-center">
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                Resetting...
              </div>
            ) : (
              'Yes, Reset Budget'
            )}
          </button>
          <button
            onClick={() => setShowConfirm(false)}
            disabled={loading}
            className="flex-1 bg-gray-200 text-gray-800 py-2 px-4 rounded-lg font-semibold hover:bg-gray-300 transition-colors"
          >
            Cancel
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
      <div className="flex items-center space-x-3 mb-4">
        <div className="p-2 bg-orange-100 rounded-lg">
          <RotateCcw className="h-5 w-5 text-orange-600" />
        </div>
        <h3 className="text-lg font-semibold text-gray-900">Reset Budget</h3>
      </div>

      <p className="text-gray-600 mb-4 text-sm">
        Clear all income and expense records to start fresh.
      </p>

      <button
        onClick={() => setShowConfirm(true)}
        className="w-full bg-orange-500 text-white py-3 px-4 rounded-lg font-semibold hover:bg-orange-600 focus:ring-2 focus:ring-orange-500 focus:ring-offset-2 transition-all duration-200 transform hover:scale-[1.02]"
      >
        Reset All Data
      </button>
    </div>
  );
};

export default ResetBudgetButton;