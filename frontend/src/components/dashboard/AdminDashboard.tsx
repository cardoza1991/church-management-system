'use client';

import React, { useState, useEffect } from 'react';

// Interfaces
interface Room {
  id: number;
  name: string;
  capacity: number;
  location?: string;
  description?: string;
  availability_start?: string;
  availability_end?: string;
  is_available?: boolean;
}

interface User {
  id: number;
  username: string;
  full_name: string;
  email: string;
  role: 'admin' | 'user';
}

interface Lesson {
  id: number;
  title: string;
  description?: string;
  sequence_number: number;
}

interface RoomForm {
  name: string;
  capacity: string;
  location: string;
  description: string;
  availability_start: string;
  availability_end: string;
  is_available: boolean;
}

interface LessonForm {
  title: string;
  description: string;
  sequence_number: string;
}

// Define admin API services
const adminAPI = {
  getRooms: (): Promise<{ rooms: Room[] }> => 
    fetch('http://localhost:8083/rooms').then(res => res.json()),
  createRoom: (data: any): Promise<any> => 
    fetch('http://localhost:8083/rooms', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(data)
    }).then(res => res.json()),
  updateRoom: (id: number, data: any): Promise<any> => 
    fetch(`http://localhost:8083/rooms/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(data)
    }).then(res => res.json()),
  deleteRoom: (id: number): Promise<{ success: boolean }> =>
    fetch(`http://localhost:8083/rooms/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    }).then(res => res.ok ? { success: true } : res.json()),
    
  getUsers: (): Promise<{ users: User[] }> => 
    fetch('http://localhost:8080/users', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    }).then(res => res.json()),
  
  getLessons: (): Promise<{ lessons: Lesson[] }> => 
    fetch('http://localhost:8082/lessons').then(res => res.json()),
  createLesson: (data: any): Promise<any> => 
    fetch('http://localhost:8082/lessons', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(data)
    }).then(res => res.json()),
  updateLesson: (id: number, data: any): Promise<any> => 
    fetch(`http://localhost:8082/lessons/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(data)
    }).then(res => res.json()),
  deleteLesson: (id: number): Promise<{ success: boolean }> =>
    fetch(`http://localhost:8082/lessons/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    }).then(res => res.ok ? { success: true } : res.json()),
};

export default function AdminDashboard() {
  const [activeTab, setActiveTab] = useState<'rooms' | 'users' | 'lessons'>('rooms');
  const [rooms, setRooms] = useState<Room[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [lessons, setLessons] = useState<Lesson[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');
  
  const [showRoomForm, setShowRoomForm] = useState<boolean>(false);
  const [editingRoomId, setEditingRoomId] = useState<number | null>(null);
  const [roomForm, setRoomForm] = useState<RoomForm>({
    name: '',
    capacity: '',
    location: '',
    description: '',
    availability_start: '',
    availability_end: '',
    is_available: true
  });
  
  const [showLessonForm, setShowLessonForm] = useState<boolean>(false);
  const [editingLessonId, setEditingLessonId] = useState<number | null>(null);
  const [lessonForm, setLessonForm] = useState<LessonForm>({
    title: '',
    description: '',
    sequence_number: ''
  });

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError('');
      
      try {
        switch (activeTab) {
          case 'rooms':
            const roomsResult = await adminAPI.getRooms();
            setRooms(roomsResult.rooms || []);
            break;
          case 'users':
            const usersResult = await adminAPI.getUsers();
            setUsers(usersResult.users || []);
            break;
          case 'lessons':
            const lessonsResult = await adminAPI.getLessons();
            setLessons(lessonsResult.lessons || []);
            break;
        }
      } catch (err) {
        console.error(`Error fetching ${activeTab}:`, err);
        setError(`Failed to load ${activeTab}. Please ensure you have admin permissions.`);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [activeTab]);

  const handleRoomFormChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value, type, checked } = e.target as HTMLInputElement;
    setRoomForm(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleRoomSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    
    try {
      const formData = {
        ...roomForm,
        capacity: parseInt(roomForm.capacity, 10)
      };
      
      if (editingRoomId) {
        await adminAPI.updateRoom(editingRoomId, formData);
      } else {
        await adminAPI.createRoom(formData);
      }
      
      setRoomForm({
        name: '',
        capacity: '',
        location: '',
        description: '',
        availability_start: '',
        availability_end: '',
        is_available: true
      });
      setShowRoomForm(false);
      setEditingRoomId(null);
      
      const result = await adminAPI.getRooms();
      setRooms(result.rooms || []);
    } catch (err) {
      console.error('Error saving room:', err);
      setError('Failed to save room. Please check all fields and try again.');
    }
  };

  const handleEditRoom = (room: Room) => {
    setRoomForm({
      name: room.name || '',
      capacity: room.capacity.toString() || '',
      location: room.location || '',
      description: room.description || '',
      availability_start: room.availability_start || '',
      availability_end: room.availability_end || '',
      is_available: room.is_available !== false
    });
    setEditingRoomId(room.id);
    setShowRoomForm(true);
  };

  const handleDeleteRoom = async (id: number) => {
    if (!window.confirm('Are you sure you want to delete this room?')) {
      return;
    }
    
    try {
      await adminAPI.deleteRoom(id);
      const result = await adminAPI.getRooms();
      setRooms(result.rooms || []);
    } catch (err) {
      console.error('Error deleting room:', err);
      setError('Failed to delete room. It may have existing reservations.');
    }
  };

  const handleLessonFormChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setLessonForm(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleLessonSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    
    try {
      const formData = {
        ...lessonForm,
        sequence_number: parseInt(lessonForm.sequence_number, 10)
      };
      
      if (editingLessonId) {
        await adminAPI.updateLesson(editingLessonId, formData);
      } else {
        await adminAPI.createLesson(formData);
      }
      
      setLessonForm({
        title: '',
        description: '',
        sequence_number: ''
      });
      setShowLessonForm(false);
      setEditingLessonId(null);
      
      const result = await adminAPI.getLessons();
      setLessons(result.lessons || []);
    } catch (err) {
      console.error('Error saving lesson:', err);
      setError('Failed to save lesson. This sequence number or title might already exist.');
    }
  };

  const handleEditLesson = (lesson: Lesson) => {
    setLessonForm({
      title: lesson.title || '',
      description: lesson.description || '',
      sequence_number: lesson.sequence_number.toString() || ''
    });
    setEditingLessonId(lesson.id);
    setShowLessonForm(true);
  };

  const handleDeleteLesson = async (id: number) => {
    if (!window.confirm('Are you sure you want to delete this lesson?')) {
      return;
    }
    
    try {
      await adminAPI.deleteLesson(id);
      const result = await adminAPI.getLessons();
      setLessons(result.lessons || []);
    } catch (err) {
      console.error('Error deleting lesson:', err);
      setError('Failed to delete lesson. It may be referenced in study records.');
    }
  };

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold">Admin Dashboard</h1>
        <p className="text-gray-600">Manage rooms, users, and system settings</p>
      </div>

      <div className="mb-6 border-b border-gray-200">
        <ul className="flex flex-wrap -mb-px">
          <li className="mr-2">
            <button
              className={`inline-block p-4 ${activeTab === 'rooms' 
                ? 'text-blue-600 border-b-2 border-blue-600 font-medium' 
                : 'text-gray-500 hover:text-gray-700 hover:border-gray-300'}`}
              onClick={() => setActiveTab('rooms')}
            >
              Rooms
            </button>
          </li>
          <li className="mr-2">
            <button
              className={`inline-block p-4 ${activeTab === 'users' 
                ? 'text-blue-600 border-b-2 border-blue-600 font-medium' 
                : 'text-gray-500 hover:text-gray-700 hover:border-gray-300'}`}
              onClick={() => setActiveTab('users')}
            >
              Users
            </button>
          </li>
          <li className="mr-2">
            <button
              className={`inline-block p-4 ${activeTab === 'lessons' 
                ? 'text-blue-600 border-b-2 border-blue-600 font-medium' 
                : 'text-gray-500 hover:text-gray-700 hover:border-gray-300'}`}
              onClick={() => setActiveTab('lessons')}
            >
              Bible Study Lessons
            </button>
          </li>
        </ul>
      </div>

      {error && (
        <div className="mb-4 p-3 bg-red-100 text-red-700 rounded">
          {error}
        </div>
      )}

      {activeTab === 'rooms' && (
        <div>
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold">Manage Rooms</h2>
            <button 
              onClick={() => {
                setRoomForm({
                  name: '',
                  capacity: '',
                  location: '',
                  description: '',
                  availability_start: '',
                  availability_end: '',
                  is_available: true
                });
                setEditingRoomId(null);
                setShowRoomForm(!showRoomForm);
              }}
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              {showRoomForm ? 'Cancel' : 'Add Room'}
            </button>
          </div>

          {showRoomForm && (
            <div className="mb-6 bg-white shadow-md rounded-lg p-4">
              <h3 className="text-lg font-medium mb-4">
                {editingRoomId ? 'Edit Room' : 'Add New Room'}
              </h3>
              <form onSubmit={handleRoomSubmit}>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Room Name *
                    </label>
                    <input
                      type="text"
                      name="name"
                      value={roomForm.name}
                      onChange={handleRoomFormChange}
                      className="w-full p-2 border border-gray-300 rounded-md"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Capacity *
                    </label>
                    <input
                      type="number"
                      name="capacity"
                      value={roomForm.capacity}
                      onChange={handleRoomFormChange}
                      className="w-full p-2 border border-gray-300 rounded-md"
                      required
                      min="1"
                    />
                  </div>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Location
                    </label>
                    <input
                      type="text"
                      name="location"
                      value={roomForm.location}
                      onChange={handleRoomFormChange}
                      className="w-full p-2 border border-gray-300 rounded-md"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Available From
                    </label>
                    <input
                      type="time"
                      name="availability_start"
                      value={roomForm.availability_start}
                      onChange={handleRoomFormChange}
                      className="w-full p-2 border border-gray-300 rounded-md"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Available Until
                    </label>
                    <input
                      type="time"
                      name="availability_end"
                      value={roomForm.availability_end}
                      onChange={handleRoomFormChange}
                      className="w-full p-2 border border-gray-300 rounded-md"
                    />
                  </div>
                </div>
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Description
                  </label>
                  <textarea
                    name="description"
                    value={roomForm.description}
                    onChange={handleRoomFormChange}
                    className="w-full p-2 border border-gray-300 rounded-md"
                    rows={3}
                  />
                </div>
                <div className="mb-4">
                  <label className="flex items-center">
                    <input
                      type="checkbox"
                      name="is_available"
                      checked={roomForm.is_available}
                      onChange={handleRoomFormChange}
                      className="mr-2"
                    />
                    <span className="text-sm font-medium text-gray-700">
                      Room is available for booking
                    </span>
                  </label>
                </div>
                <div className="flex justify-end">
                  <button
                    type="submit"
                    className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
                  >
                    {editingRoomId ? 'Update Room' : 'Add Room'}
                  </button>
                </div>
              </form>
            </div>
          )}

          {loading ? (
            <div className="text-center py-8">Loading rooms...</div>
          ) : rooms.length === 0 ? (
            <div className="bg-white shadow-md rounded-lg p-6 text-center">
              <p className="text-gray-600">No rooms found. Click "Add Room" to get started.</p>
            </div>
          ) : (
            <div className="bg-white shadow-md rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Name
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Capacity
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Location
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Availability
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {rooms.map((room) => (
                    <tr key={room.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm font-medium text-gray-900">
                          {room.name}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {room.capacity}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {room.location || '—'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {room.availability_start && room.availability_end 
                          ? `${room.availability_start} - ${room.availability_end}`
                          : 'All day'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        {room.is_available !== false ? (
                          <span className="px-2 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800">
                            Available
                          </span>
                        ) : (
                          <span className="px-2 py-1 text-xs font-semibold rounded-full bg-red-100 text-red-800">
                            Unavailable
                          </span>
                        )}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                        <button
                          onClick={() => handleEditRoom(room)}
                          className="text-blue-600 hover:text-blue-900 mr-3"
                        >
                          Edit
                        </button>
                        <button
                          onClick={() => handleDeleteRoom(room.id)}
                          className="text-red-600 hover:text-red-900"
                        >
                          Delete
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}

      {activeTab === 'users' && (
        <div>
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold">Manage Users</h2>
          </div>

          {loading ? (
            <div className="text-center py-8">Loading users...</div>
          ) : users.length === 0 ? (
            <div className="bg-white shadow-md rounded-lg p-6 text-center">
              <p className="text-gray-600">No users found or you don't have permission to view users.</p>
            </div>
          ) : (
            <div className="bg-white shadow-md rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Username
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Full Name
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Email
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Role
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {users.map((user) => (
                    <tr key={user.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm font-medium text-gray-900">
                          {user.username}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {user.full_name}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {user.email}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`px-2 py-1 text-xs font-semibold rounded-full ${
                          user.role === 'admin' 
                            ? 'bg-purple-100 text-purple-800' 
                            : 'bg-blue-100 text-blue-800'
                        }`}>
                          {user.role}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}

      {activeTab === 'lessons' && (
        <div>
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold">Manage Bible Study Lessons</h2>
            <button 
              onClick={() => {
                setLessonForm({
                  title: '',
                  description: '',
                  sequence_number: ''
                });
                setEditingLessonId(null);
                setShowLessonForm(!showLessonForm);
              }}
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              {showLessonForm ? 'Cancel' : 'Add Lesson'}
            </button>
          </div>

          {showLessonForm && (
            <div className="mb-6 bg-white shadow-md rounded-lg p-4">
              <h3 className="text-lg font-medium mb-4">
                {editingLessonId ? 'Edit Lesson' : 'Add New Lesson'}
              </h3>
              <form onSubmit={handleLessonSubmit}>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Lesson Title *
                    </label>
                    <input
                      type="text"
                      name="title"
                      value={lessonForm.title}
                      onChange={handleLessonFormChange}
                      className="w-full p-2 border border-gray-300 rounded-md"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Sequence Number *
                    </label>
                    <input
                      type="number"
                      name="sequence_number"
                      value={lessonForm.sequence_number}
                      onChange={handleLessonFormChange}
                      className="w-full p-2 border border-gray-300 rounded-md"
                      required
                      min="1"
                    />
                  </div>
                </div>
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Description
                  </label>
                  <textarea
                    name="description"
                    value={lessonForm.description}
                    onChange={handleLessonFormChange}
                    className="w-full p-2 border border-gray-300 rounded-md"
                    rows={3}
                  />
                </div>
                <div className="flex justify-end">
                  <button
                    type="submit"
                    className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
                  >
                    {editingLessonId ? 'Update Lesson' : 'Add Lesson'}
                  </button>
                </div>
              </form>
            </div>
          )}

          {loading ? (
            <div className="text-center py-8">Loading lessons...</div>
          ) : lessons.length === 0 ? (
            <div className="bg-white shadow-md rounded-lg p-6 text-center">
              <p className="text-gray-600">No lessons found. Click "Add Lesson" to get started.</p>
            </div>
          ) : (
            <div className="bg-white shadow-md rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Sequence
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
                  {lessons.sort((a, b) => a.sequence_number - b.sequence_number).map((lesson) => (
                    <tr key={lesson.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {lesson.sequence_number}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {lesson.title}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-500">
                        {lesson.description || '—'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                        <button
                          onClick={() => handleEditLesson(lesson)}
                          className="text-blue-600 hover:text-blue-900 mr-3"
                        >
                          Edit
                        </button>
                        <button
                          onClick={() => handleDeleteLesson(lesson.id)}
                          className="text-red-600 hover:text-red-900"
                        >
                          Delete
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}
    </div>
  );
}