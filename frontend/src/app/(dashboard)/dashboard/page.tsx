'use client';

import React, { useState } from 'react';
import { useAuth } from '@/components/auth/AuthProvider';
import ContactsDashboard from '@/components/dashboard/ContactsDashboard';
import StudiesDashboard from '@/components/dashboard/StudiesDashboard';
import ReservationsDashboard from '@/components/dashboard/ReservationsDashboard';
import AdminDashboard from '@/components/dashboard/AdminDashboard';

// Assuming this is the shape of your user object from AuthProvider
interface User {
  full_name: string;
  role: 'admin' | 'user';
}

interface AuthContext {
  user: User | null;
  loading: boolean;
}

// Define props for StatCard
interface StatCardProps {
  title: string;
  value: string;
  icon: string;
  color: string;
}

export default function DashboardPage() {
  const { user, loading } = useAuth() as AuthContext;
  const [activeSection, setActiveSection] = useState<'overview' | 'contacts' | 'studies' | 'reservations' | 'admin'>('overview');

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-xl text-gray-600">Loading...</div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="bg-red-100 text-red-700 p-4 rounded">
        <p>You must be logged in to view this page.</p>
      </div>
    );
  }

  const StatCard: React.FC<StatCardProps> = ({ title, value, icon, color }) => (
    <div className={`bg-white rounded-lg shadow-md p-6 ${color}`}>
      <div className="flex items-center">
        <div className="p-3 rounded-full bg-opacity-20 mr-4">
          <span className="text-2xl">{icon}</span>
        </div>
        <div>
          <p className="text-sm text-gray-500">{title}</p>
          <p className="text-2xl font-bold">{value}</p>
        </div>
      </div>
    </div>
  );

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <div className="text-gray-600">
          Welcome, <span className="font-semibold">{user.full_name}</span>
        </div>
      </div>

      <div className="mb-6 border-b border-gray-200">
        <ul className="flex flex-wrap -mb-px">
          <li className="mr-2">
            <button
              className={`inline-block p-4 ${activeSection === 'overview' 
                ? 'text-blue-600 border-b-2 border-blue-600 font-medium' 
                : 'text-gray-500 hover:text-gray-700 hover:border-gray-300'}`}
              onClick={() => setActiveSection('overview')}
            >
              Overview
            </button>
          </li>
          <li className="mr-2">
            <button
              className={`inline-block p-4 ${activeSection === 'contacts' 
                ? 'text-blue-600 border-b-2 border-blue-600 font-medium' 
                : 'text-gray-500 hover:text-gray-700 hover:border-gray-300'}`}
              onClick={() => setActiveSection('contacts')}
            >
              Contacts
            </button>
          </li>
          <li className="mr-2">
            <button
              className={`inline-block p-4 ${activeSection === 'studies' 
                ? 'text-blue-600 border-b-2 border-blue-600 font-medium' 
                : 'text-gray-500 hover:text-gray-700 hover:border-gray-300'}`}
              onClick={() => setActiveSection('studies')}
            >
              Bible Studies
            </button>
          </li>
          <li className="mr-2">
            <button
              className={`inline-block p-4 ${activeSection === 'reservations' 
                ? 'text-blue-600 border-b-2 border-blue-600 font-medium' 
                : 'text-gray-500 hover:text-gray-700 hover:border-gray-300'}`}
              onClick={() => setActiveSection('reservations')}
            >
              Reservations
            </button>
          </li>
          {user.role === 'admin' && (
            <li className="mr-2">
              <button
                className={`inline-block p-4 ${activeSection === 'admin' 
                  ? 'text-blue-600 border-b-2 border-blue-600 font-medium' 
                  : 'text-gray-500 hover:text-gray-700 hover:border-gray-300'}`}
                onClick={() => setActiveSection('admin')}
              >
                Admin
              </button>
            </li>
          )}
        </ul>
      </div>

      {activeSection === 'overview' && (
        <div>
          <div className="mb-6">
            <h2 className="text-xl font-semibold mb-4">System Overview</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              <StatCard 
                title="Total Contacts" 
                value="42" 
                icon="👥" 
                color="border-t-4 border-blue-500" 
              />
              <StatCard 
                title="Bible Studies" 
                value="18" 
                icon="📚" 
                color="border-t-4 border-yellow-500" 
              />
              <StatCard 
                title="Room Bookings" 
                value="7" 
                icon="🏢" 
                color="border-t-4 border-green-500" 
              />
              <StatCard 
                title="Gospel Workers" 
                value="5" 
                icon="🙏" 
                color="border-t-4 border-purple-500" 
              />
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="bg-white shadow-md rounded-lg p-6">
              <h3 className="text-lg font-semibold mb-4">Recent Contacts</h3>
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <div>
                    <p className="font-medium">John Smith</p>
                    <p className="text-sm text-gray-500">Added 2 days ago</p>
                  </div>
                  <span className="px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800">
                    New Contact
                  </span>
                </div>
                <div className="flex justify-between items-center">
                  <div>
                    <p className="font-medium">Maria Garcia</p>
                    <p className="text-sm text-gray-500">Added 5 days ago</p>
                  </div>
                  <span className="px-2 py-1 text-xs font-semibold rounded-full bg-yellow-100 text-yellow-800">
                    In Studies
                  </span>
                </div>
                <div className="flex justify-between items-center">
                  <div>
                    <p className="font-medium">David Johnson</p>
                    <p className="text-sm text-gray-500">Added 1 week ago</p>
                  </div>
                  <span className="px-2 py-1 text-xs font-semibold rounded-full bg-purple-100 text-purple-800">
                    Baptized
                  </span>
                </div>
              </div>
              <button
                onClick={() => setActiveSection('contacts')}
                className="mt-4 text-blue-600 hover:text-blue-800"
              >
                View all contacts →
              </button>
            </div>

            <div className="bg-white shadow-md rounded-lg p-6">
              <h3 className="text-lg font-semibold mb-4">Today's Reservations</h3>
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <div>
                    <p className="font-medium">Bible Study - Fellowship Hall</p>
                    <p className="text-sm text-gray-500">10:00 AM - 11:30 AM</p>
                  </div>
                </div>
                <div className="flex justify-between items-center">
                  <div>
                    <p className="font-medium">Prayer Meeting - Room A</p>
                    <p className="text-sm text-gray-500">2:00 PM - 3:00 PM</p>
                  </div>
                </div>
                <div className="flex justify-between items-center">
                  <div>
                    <p className="font-medium">Youth Group - Main Hall</p>
                    <p className="text-sm text-gray-500">6:30 PM - 8:30 PM</p>
                  </div>
                </div>
              </div>
              <button
                onClick={() => setActiveSection('reservations')}
                className="mt-4 text-blue-600 hover:text-blue-800"
              >
                Manage reservations →
              </button>
            </div>
          </div>
        </div>
      )}

      {activeSection === 'contacts' && <ContactsDashboard />}
      {activeSection === 'studies' && <StudiesDashboard />}
      {activeSection === 'reservations' && <ReservationsDashboard />}
      {activeSection === 'admin' && user.role === 'admin' && <AdminDashboard />}
    </div>
  );
}