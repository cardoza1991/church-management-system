import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { contactsAPI, statusesAPI } from '@/lib/api';

export default function ContactForm({ contactId = null }) {
  const router = useRouter();
  const [contact, setContact] = useState({
    name: '',
    location: '',
    phone: '',
    email: '',
    notes: '',
    current_status_id: 1, // Default to "New Contact"
  });
  const [statuses, setStatuses] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [isEdit, setIsEdit] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        // Fetch statuses
        const statusesRes = await statusesAPI.getStatuses();
        setStatuses(statusesRes.data.statuses);
        
        // If we have a contactId, fetch the contact data
        if (contactId) {
          setIsEdit(true);
          const contactRes = await contactsAPI.getContact(contactId);
          setContact(contactRes.data);
        }
      } catch (error) {
        setError('Failed to load data. Please try again.');
        console.error(error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [contactId]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setContact(prev => ({
      ...prev,
      [name]: name === 'current_status_id' ? parseInt(value, 10) : value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      if (isEdit) {
        await contactsAPI.updateContact(contactId, contact);
      } else {
        await contactsAPI.createContact(contact);
      }
      router.push('/contacts');
    } catch (error) {
      setError('Failed to save contact. Please try again.');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  if (loading && !contact.name) {
    return <div className="text-center py-4">Loading...</div>;
  }

  return (
    <form onSubmit={handleSubmit} className="bg-white shadow-md rounded-lg p-6">
      {error && (
        <div className="mb-4 p-3 bg-red-100 text-red-700 rounded">
          {error}
        </div>
      )}
      
      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="name">
          Name *
        </label>
        <input
          id="name"
          name="name"
          type="text"
          value={contact.name}
          onChange={handleChange}
          className="w-full px-3 py-2 border rounded shadow appearance-none"
          required
        />
      </div>

      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="location">
          Location
        </label>
        <input
          id="location"
          name="location"
          type="text"
          value={contact.location}
          onChange={handleChange}
          className="w-full px-3 py-2 border rounded shadow appearance-none"
        />
      </div>

      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="phone">
          Phone
        </label>
        <input
          id="phone"
          name="phone"
          type="tel"
          value={contact.phone}
          onChange={handleChange}
          className="w-full px-3 py-2 border rounded shadow appearance-none"
        />
      </div>

      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="email">
          Email
        </label>
        <input
          id="email"
          name="email"
          type="email"
          value={contact.email}
          onChange={handleChange}
          className="w-full px-3 py-2 border rounded shadow appearance-none"
        />
      </div>

      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="current_status_id">
          Status
        </label>
        <select
          id="current_status_id"
          name="current_status_id"
          value={contact.current_status_id}
          onChange={handleChange}
          className="w-full px-3 py-2 border rounded shadow appearance-none"
          required
        >
          {statuses.map(status => (
            <option key={status.id} value={status.id}>
              {status.name}
            </option>
          ))}
        </select>
      </div>

      <div className="mb-6">
        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="notes">
          Notes
        </label>
        <textarea
          id="notes"
          name="notes"
          value={contact.notes}
          onChange={handleChange}
          className="w-full px-3 py-2 border rounded shadow appearance-none"
          rows={4}
        />
      </div>

      <div className="flex items-center justify-between">
        <button
          type="submit"
          className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
          disabled={loading}
        >
          {loading ? 'Saving...' : isEdit ? 'Update Contact' : 'Create Contact'}
        </button>
        <button
          type="button"
          className="bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
          onClick={() => router.push('/contacts')}
        >
          Cancel
        </button>
      </div>
    </form>
  );
}
