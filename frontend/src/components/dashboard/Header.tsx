'use client';

import { useAuth } from '@/components/auth/AuthProvider';

export default function Header() {
  const { user, logout } = useAuth();

  return (
    <header className="bg-white shadow">
      <div className="container mx-auto py-4 px-6 flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-800">Church Management System</h1>
        
        <div className="flex items-center space-x-4">
          {user && (
            <>
              <span className="text-gray-600">
                Welcome, {user.full_name}
              </span>
              <button
                onClick={logout}
                className="px-3 py-1 border border-gray-300 rounded hover:bg-gray-100"
              >
                Logout
              </button>
            </>
          )}
        </div>
      </div>
    </header>
  );
}