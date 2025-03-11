'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';

// Interfaces
interface Room {
  id: number;
  name: string;
  capacity: number;
}

interface Reservation {
  id: number;
  room_id: number;
  room_name?: string;
  title: string;
  description?: string;
  start_time: string;
  end_time: string;
}

interface FormData {
  room_id: string;
  title: string;
  description: string;
  start_time: string;
  end_time: string;
}

// Define a custom API for the reservation service
const reservationsAPI = {
  getRooms: (): Promise<{ rooms: Room[] }> => 
    fetch('http://localhost:8083/rooms').then(res => res.json()),
  getReservationsByDate: (start: string, end: string, roomId: number | null = null): Promise<{ reservations: Reservation[] }> => {
    let url = `http://localhost:8083/reservations/by-date?start=${start}&end=${end}`;
    if (roomId) url += `&room_id=${roomId}`;
    
    return fetch(url, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    }).then(res => res.json());
  },
  createReservation: (data: any): Promise<any> => 
    fetch('http://localhost:8083/reservations', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(data)
    }).then(res => res.json()),
  deleteReservation: (id: number): Promise<{ success: boolean }> =>
    fetch(`http://localhost:8083/reservations/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    }).then(res => res.ok ? { success: true } : res.json())
};

export default function ReservationsDashboard() {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [selectedDate, setSelectedDate] = useState<string>(new Date().toISOString().split('T')[0]);
  const [selectedRoomId, setSelectedRoomId] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');
  const [showNewReservationForm, setShowNewReservationForm] = useState<boolean>(false);
  const [formData, setFormData] = useState<FormData>({
    room_id: '',
    title: '',
    description: '',
    start_time: '',
    end_time: '',
  });

  // Fetch rooms on component mount
  useEffect(() => {
    const fetchRooms = async () => {
      try {
        const result = await reservationsAPI.getRooms();
        setRooms(result.rooms || []);
      } catch (err) {
        console.error('Error fetching rooms:', err);
        setError('Failed to load rooms. Please try again later.');
      }
    };

    fetchRooms();
  }, []);

  // Fetch reservations when selected date or room changes
  useEffect(() => {
    const fetchReservations = async () => {
      if (!selectedDate) return;
      
      setLoading(true);
      setError('');
      try {
        const result = await reservationsAPI.getReservationsByDate(
          selectedDate, 
          selectedDate, 
          selectedRoomId ? parseInt(selectedRoomId, 10) : null
        );
        
        setReservations(result.reservations || []);
      } catch (err) {
        console.error('Error fetching reservations:', err);
        setError('Failed to load reservations. Please try again later.');
      } finally {
        setLoading(false);
      }
    };

    fetchReservations();
  }, [selectedDate, selectedRoomId]);

  const handleDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedDate(e.target.value);
  };

  const handleRoomChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setSelectedRoomId(e.target.value);
  };

  const handleFormChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    
    if (!formData.room_id || !formData.title || !formData.start_time || !formData.end_time) {
      setError('Please fill in all required fields');
      return;
    }

    try {
      const reservationData = {
        ...formData,
        room_id: parseInt(formData.room_id, 10),
        start_time: new Date(`${selectedDate}T${formData.start_time}`).toISOString(),
        end_time: new Date(`${selectedDate}T${formData.end_time}`).toISOString(),
      };
      
      await reservationsAPI.createReservation(reservationData);
      
      setFormData({
        room_id: '',
        title: '',
        description: '',
        start_time: '',
        end_time: '',
      });
      setShowNewReservationForm(false);
      
      const updated = await reservationsAPI.getReservationsByDate(
        selectedDate, 
        selectedDate, 
        selectedRoomId ? parseInt(selectedRoomId, 10) : null
      );
      setReservations(updated.reservations || []);
      
    } catch (err) {
      console.error('Error creating reservation:', err);
      setError('Failed to create reservation. The room might not be available for the selected time.');
    }
  };
  
  const handleDeleteReservation = async (id: number) => {
    if (!window.confirm('Are you sure you want to cancel this reservation?')) {
      return;
    }
    
    try {
      await reservationsAPI.deleteReservation(id);
      
      const updated = await reservationsAPI.getReservationsByDate(
        selectedDate, 
        selectedDate, 
        selectedRoomId ? parseInt(selectedRoomId, 10) : null
      );
      setReservations(updated.reservations || []);
    } catch (err) {
      console.error('Error deleting reservation:', err);
      setError('Failed to cancel reservation.');
    }
  };

  const formatTime = (isoString: string): string => {
    const date = new Date(isoString);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };
  
  const getRoomName = (roomId: number): string => {
    const room = rooms.find(r => r.id === roomId);
    return room ? room.name : 'Unknown Room';
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Room Reservations</h1>
        <button 
          onClick={() => setShowNewReservationForm(!showNewReservationForm)}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          {showNewReservationForm ? 'Cancel' : 'Book a Room'}
        </button>
      </div>

      <div className="mb-6 bg-white shadow-md rounded-lg p-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Date</label>
            <input
              type="date"
              value={selectedDate}
              onChange={handleDateChange}
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Room (Optional)</label>
            <select 
              value={selectedRoomId}
              onChange={handleRoomChange}
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">All Rooms</option>
              {rooms.map(room => (
                <option key={room.id} value={room.id}>
                  {room.name} (Capacity: {room.capacity})
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {showNewReservationForm && (
        <div className="mb-6 bg-white shadow-md rounded-lg p-4">
          <h2 className="text-lg font-semibold mb-4">Book a Room</h2>
          {error && (
            <div className="mb-4 p-3 bg-red-100 text-red-700 rounded">
              {error}
            </div>
          )}
          <form onSubmit={handleSubmit}>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Room *</label>
                <select 
                  name="room_id"
                  value={formData.room_id}
                  onChange={handleFormChange}
                  className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                >
                  <option value="">Select a Room</option>
                  {rooms.map(room => (
                    <option key={room.id} value={room.id}>
                      {room.name} (Capacity: {room.capacity})
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Title *</label>
                <input
                  type="text"
                  name="title"
                  value={formData.title}
                  onChange={handleFormChange}
                  className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Meeting, Bible Study, etc."
                  required
                />
              </div>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Start Time *</label>
                <input
                  type="time"
                  name="start_time"
                  value={formData.start_time}
                  onChange={handleFormChange}
                  className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">End Time *</label>
                <input
                  type="time"
                  name="end_time"
                  value={formData.end_time}
                  onChange={handleFormChange}
                  className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-2">Description</label>
              <textarea
                name="description"
                value={formData.description}
                onChange={handleFormChange}
                className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Optional description of the reservation"
                rows={3}
              />
            </div>
            <div>
              <button 
                type="submit" 
                className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
              >
                Book Now
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="bg-white shadow-md rounded-lg overflow-hidden">
        <h2 className="text-lg font-semibold p-4 border-b">
          Reservations for {new Date(selectedDate).toLocaleDateString()} 
          {selectedRoomId && ` - ${getRoomName(parseInt(selectedRoomId, 10))}`}
        </h2>
        
        {loading ? (
          <div className="p-6 text-center">
            <div className="text-gray-600">Loading reservations...</div>
          </div>
        ) : reservations.length === 0 ? (
          <div className="p-6 text-center text-gray-500">
            No reservations found for this date.
          </div>
        ) : (
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Room
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Time
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Title
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Description
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {reservations.map((reservation) => (
                <tr key={reservation.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900">
                      {reservation.room_name || getRoomName(reservation.room_id)}
                    </div>
                    <div className="text-xs text-gray-500">
                      Capacity: {rooms.find(r => r.id === reservation.room_id)?.capacity || 'N/A'}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {formatTime(reservation.start_time)} - {formatTime(reservation.end_time)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {reservation.title}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-500">
                    {reservation.description || '—'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <Link 
                      href={`/reservations/${reservation.id}/edit`}
                      className="text-blue-600 hover:text-blue-900 mr-3"
                    >
                      Edit
                    </Link>
                    <button
                      onClick={() => handleDeleteReservation(reservation.id)}
                      className="text-red-600 hover:text-red-900"
                    >
                      Cancel
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}