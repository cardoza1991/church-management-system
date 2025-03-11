'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { contactsAPI, statusesAPI } from '@/lib/api';

export default function ViewContactPage({ params }) {
  const router = useRouter();
  const contactId = parseInt(params.id, 10);
  const [contact, setContact] = useState(null);
  const [statuses, setStatuses] = useState([]);
  const [statusHistory, setStatusHistory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const [contactRes, statusesRes, historyRes] = await Promise.all([
          contactsAPI.getContact(contactId),
          statusesAPI.getStatuses(),
          contactsAPI.getContactStatusHistory(contactId),
        ]);
        
        setContact(contactRes.data);
        setStatuses(statusesRes.data.statuses);
        setStatusHistory(historyRes.data.history || []);
      } catch (error) {
        console.error('Failed to fetch contact data', error);
        setError('Failed to load contact. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [contactId]);

  const getStatusName = (statusId) => {
    const status = statuses.find(s => s.id === statusId);
    return status ? status.name : 'Unknown';
  };

  if (loading) {
    return <div className="text-center py-4">Loading contact information...</div>;
  }

  if (error || !contact) {
    return (
      <div className="bg-red-100 text-red-700 p-4 rounded">
        {error || 'Contact not found'}
        <button 
          onClick={() => router.push('/contacts')}
          className="mt-2 bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700"
        >
          Back to Contacts
        </button>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">{contact.name}</h1>
        <div className="space-x-3">
          <Link 
            href={'/contacts/' + contactId + '/edit'} 
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Edit Contact
          </Link>
          <Link 
            href="/contacts" 
            className="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
          >
            Back to List
          </Link>
        </div>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="md:col-span-2">
          <div className="bg-white shadow-md rounded-lg overflow-hidden mb-6">
            <div className="p-6">
              <h2 className="text-xl font-semibold mb-4">Contact Information</h2>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-gray-500">Status</p>
                  <p className="font-medium">{getStatusName(contact.current_status_id)}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-500">Added On</p>
                  <p className="font-medium">{new Date(contact.date_added).toLocaleDateString()}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-500">Email</p>
                  <p className="font-medium">{contact.email || '—'}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-500">Phone</p>
                  <p className="font-medium">{contact.phone || '—'}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-500">Location</p>
                  <p className="font-medium">{contact.location || '—'}</p>
                </div>
              </div>
              
              {contact.notes && (
                <div className="mt-4">
                  <p className="text-sm text-gray-500">Notes</p>
                  <p className="mt-1 whitespace-pre-line">{contact.notes}</p>
                </div>
              )}
            </div>
          </div>
          
          <div className="bg-white shadow-md rounded-lg overflow-hidden">
            <div className="p-6">
              <h2 className="text-xl font-semibold mb-4">Bible Studies</h2>
              <p className="text-gray-500">No Bible studies recorded yet.</p>
              <button className="mt-4 bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700">
                Add Bible Study Session
              </button>
            </div>
          </div>
        </div>
        
        <div>
          <div className="bg-white shadow-md rounded-lg overflow-hidden">
            <div className="p-6">
              <h2 className="text-xl font-semibold mb-4">Status History</h2>
              {statusHistory.length > 0 ? (
                <div className="space-y-4">
                  {statusHistory.map((change) => (
                    <div key={change.id} className="border-l-4 border-blue-500 pl-4 py-1">
                      <p className="font-medium">{change.status_name}</p>
                      <p className="text-sm text-gray-500">
                        {new Date(change.date_changed).toLocaleString()}
                      </p>
                      {change.notes && (
                        <p className="text-sm mt-1">{change.notes}</p>
                      )}
                    </div>
                  ))}
                </div>
              ) : (
                <p className="text-gray-500">No status changes recorded yet.</p>
              )}
              
              <div className="mt-6">
                <h3 className="font-semibold mb-2">Update Status</h3>
                <select 
                  className="w-full p-2 border rounded mb-2"
                  defaultValue={contact.current_status_id}
                >
                  {statuses.map(status => (
                    <option key={status.id} value={status.id}>
                      {status.name}
                    </option>
                  ))}
                </select>
                <textarea
                  className="w-full p-2 border rounded mb-2"
                  placeholder="Notes about this status change"
                  rows={2}
                />
                <button className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
                  Update Status
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}