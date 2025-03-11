'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { contactsAPI, statusesAPI } from '@/lib/api';

export default function ContactsPage() {
  const [contacts, setContacts] = useState([]);
  const [statuses, setStatuses] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [contactsRes, statusesRes] = await Promise.all([
          contactsAPI.getContacts(),
          statusesAPI.getStatuses(),
        ]);
        
        setContacts(contactsRes.data.contacts);
        setStatuses(statusesRes.data.statuses);
      } catch (error) {
        console.error('Failed to fetch contacts', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const getStatusName = (statusId) => {
    const status = statuses.find(s => s.id === statusId);
    return status ? status.name : 'Unknown';
  };
  
  const getStatusColor = (statusId) => {
    const status = getStatusName(statusId);
    switch (status) {
      case 'New Contact':
        return 'bg-blue-100 text-blue-800';
      case 'In Studies':
        return 'bg-yellow-100 text-yellow-800';
      case 'Baptized':
        return 'bg-purple-100 text-purple-800';
      case 'Gospel Worker':
        return 'bg-green-100 text-green-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Contacts</h1>
        <Link 
          href="/contacts/new" 
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          Add New Contact
        </Link>
      </div>
      
      {loading ? (
        <div className="text-center">Loading contacts...</div>
      ) : (
        <div className="bg-white shadow-md rounded-lg overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Name
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Phone
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Added
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {contacts.map((contact) => (
                <tr key={contact.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900">
                      {contact.name}
                    </div>
                    <div className="text-sm text-gray-500">
                      {contact.email}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={\`px-2 py-1 text-xs font-semibold rounded-full \${getStatusColor(contact.current_status_id)}\`}>
                      {getStatusName(contact.current_status_id)}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {contact.phone || '—'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {new Date(contact.date_added).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <Link 
                      href={\`/contacts/\${contact.id}\`}
                      className="text-blue-600 hover:text-blue-900 mr-3"
                    >
                      View
                    </Link>
                    <Link 
                      href={\`/contacts/\${contact.id}/edit\`}
                      className="text-green-600 hover:text-green-900"
                    >
                      Edit
                    </Link>
                  </td>
                </tr>
              ))}
              
              {contacts.length === 0 && (
                <tr>
                  <td colSpan={5} className="px-6 py-4 text-center text-sm text-gray-500">
                    No contacts found. Click "Add New Contact" to get started.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}