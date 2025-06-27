import React, { useState } from 'react';
import AuthForm from '../components/auth/AuthForm';

const AuthPage: React.FC = () => {
  const [isSignIn, setIsSignIn] = useState(true);

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50 flex items-center justify-center p-4">
      <div className="w-full max-w-6xl flex items-center justify-center lg:justify-between">
        <div className="hidden lg:block lg:w-1/2 xl:w-3/5">
          <div className="p-8">
            <h1 className="text-4xl xl:text-5xl font-bold text-gray-900 mb-6">
              Take Control of Your
              <span className="bg-gradient-to-r from-emerald-600 to-blue-600 bg-clip-text text-transparent"> Finances</span>
            </h1>
            <p className="text-xl text-gray-600 mb-8 leading-relaxed">
              Track your income and expenses with ease. Get insights into your spending patterns and make better financial decisions.
            </p>
            <div className="space-y-4">
              <div className="flex items-center space-x-3">
                <div className="w-2 h-2 bg-emerald-500 rounded-full"></div>
                <span className="text-gray-700">Real-time budget tracking</span>
              </div>
              <div className="flex items-center space-x-3">
                <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
                <span className="text-gray-700">Detailed expense categorization</span>
              </div>
              <div className="flex items-center space-x-3">
                <div className="w-2 h-2 bg-purple-500 rounded-full"></div>
                <span className="text-gray-700">Secure data management</span>
              </div>
            </div>
          </div>
        </div>
        
        <div className="w-full lg:w-1/2 xl:w-2/5">
          <AuthForm isSignIn={isSignIn} onToggle={() => setIsSignIn(!isSignIn)} />
        </div>
      </div>
    </div>
  );
};

export default AuthPage;