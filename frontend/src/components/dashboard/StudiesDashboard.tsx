'use client';

import React, { useState, useEffect } from 'react';
import { contactsAPI } from '@/lib/api';
import Link from 'next/link';

// Interfaces for data structures
interface Contact {
  id: number;
  name: string;
}

interface Lesson {
  id: number;
  title: string;
}

interface Study {
  id: number;
  lesson_id: number;
  lesson_title?: string;
  date_completed: string;
  duration_minutes?: number;
  location?: string;
}

interface StudyStats {
  completed_lessons: number;
  total_lessons: number;
  progress_percentage: number;
  last_study_date?: string;
  total_study_time_minutes: number;
}

// Define a custom API for the study service
const studiesAPI = {
  getLessons: (): Promise<{ lessons: Lesson[] }> => 
    fetch('http://localhost:8082/lessons').then(res => res.json()),
  getStudiesByContact: (contactId: number): Promise<{ studies: Study[] }> => 
    fetch(`http://localhost:8082/contacts/${contactId}/studies`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    }).then(res => res.json()),
  getContactStudyStats: (contactId: number): Promise<StudyStats> => 
    fetch(`http://localhost:8082/contacts/${contactId}/study-stats`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    }).then(res => res.json())
};

export default function StudiesDashboard() {
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [selectedContactId, setSelectedContactId] = useState<number | null>(null);
  const [studies, setStudies] = useState<Study[]>([]);
  const [studyStats, setStudyStats] = useState<StudyStats | null>(null);
  const [lessons, setLessons] = useState<Lesson[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    const fetchInitialData = async () => {
      setLoading(true);
      setError('');
      try {
        // Fetch contacts and lessons in parallel
        const [contactsRes, lessonsRes] = await Promise.all([
          contactsAPI.getContacts(),
          studiesAPI.getLessons()
        ]);
        
        setContacts(contactsRes.data.contacts);
        setLessons(lessonsRes.lessons);

        // Select first contact by default if available
        if (contactsRes.data.contacts.length > 0) {
          setSelectedContactId(contactsRes.data.contacts[0].id);
        }
      } catch (err) {
        console.error('Error fetching initial data:', err);
        setError('Failed to load data. Please try again later.');
      } finally {
        setLoading(false);
      }
    };

    fetchInitialData();
  }, []);

  // Fetch studies when selected contact changes
  useEffect(() => {
    const fetchStudyData = async () => {
      if (!selectedContactId) return;
      
      setLoading(true);
      setError('');
      try {
        // Fetch studies and stats for the selected contact
        const [studiesRes, statsRes] = await Promise.all([
          studiesAPI.getStudiesByContact(selectedContactId),
          studiesAPI.getContactStudyStats(selectedContactId)
        ]);
        
        setStudies(studiesRes.studies || []);
        setStudyStats(statsRes);
      } catch (err) {
        console.error('Error fetching study data:', err);
        setError('Failed to load studies. Please try again later.');
      } finally {
        setLoading(false);
      }
    };

    fetchStudyData();
  }, [selectedContactId]);

  // Find contact name by ID
  const getContactName = (contactId: number): string => {
    const contact = contacts.find(c => c.id === contactId);
    return contact ? contact.name : 'Unknown Contact';
  };

  // Find lesson title by ID
  const getLessonTitle = (lessonId: number): string => {
    const lesson = lessons.find(l => l.id === lessonId);
    return lesson ? lesson.title : 'Unknown Lesson';
  };

  const handleContactChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setSelectedContactId(parseInt(e.target.value, 10));
  };

  if (loading && contacts.length === 0) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-xl text-gray-600">Loading...</div>
      </div>
    );
  }

  if (error && contacts.length === 0) {
    return (
      <div className="bg-red-100 text-red-700 p-4 rounded">
        <p>{error}</p>
        <button 
          onClick={() => window.location.reload()}
          className="mt-2 px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
        >
          Try Again
        </button>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Bible Studies</h1>
        {selectedContactId && (
          <Link 
            href={`/studies/new?contactId=${selectedContactId}`}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Add New Study
          </Link>
        )}
      </div>

      {/* Contact selector */}
      <div className="mb-6 bg-white shadow-md rounded-lg p-4">
        <label className="block text-sm font-medium text-gray-700 mb-2">Select Contact</label>
        <select 
          className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          value={selectedContactId || ''}
          onChange={handleContactChange}
        >
          <option value="">-- Select a contact --</option>
          {contacts.map(contact => (
            <option key={contact.id} value={contact.id}>
              {contact.name}
            </option>
          ))}
        </select>
      </div>

      {/* Stats summary */}
      {selectedContactId && studyStats && (
        <div className="mb-6 bg-white shadow-md rounded-lg p-4">
          <h2 className="text-lg font-semibold mb-4">Study Progress for {getContactName(selectedContactId)}</h2>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="bg-blue-50 p-4 rounded-lg">
              <div className="text-sm text-blue-600 font-medium">Completed Lessons</div>
              <div className="text-2xl font-bold">{studyStats.completed_lessons} / {studyStats.total_lessons}</div>
            </div>
            <div className="bg-green-50 p-4 rounded-lg">
              <div className="text-sm text-green-600 font-medium">Progress</div>
              <div className="text-2xl font-bold">{studyStats.progress_percentage.toFixed(1)}%</div>
            </div>
            <div className="bg-purple-50 p-4 rounded-lg">
              <div className="text-sm text-purple-600 font-medium">Last Study</div>
              <div className="text-2xl font-bold">
                {studyStats.last_study_date 
                  ? new Date(studyStats.last_study_date).toLocaleDateString() 
                  : 'None'}
              </div>
            </div>
            <div className="bg-yellow-50 p-4 rounded-lg">
              <div className="text-sm text-yellow-600 font-medium">Total Study Time</div>
              <div className="text-2xl font-bold">{Math.floor(studyStats.total_study_time_minutes / 60)}h {studyStats.total_study_time_minutes % 60}m</div>
            </div>
          </div>
        </div>
      )}

      {/* Studies list */}
      {selectedContactId && (
        <div className="bg-white shadow-md rounded-lg overflow-hidden">
          <h2 className="text-lg font-semibold p-4 border-b">Completed Studies</h2>
          
          {studies.length === 0 ? (
            <div className="p-6 text-center text-gray-500">
              No studies recorded yet. Click "Add New Study" to get started.
            </div>
          ) : (
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Lesson
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Date Completed
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Duration
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Location
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {studies.map((study) => (
                  <tr key={study.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm font-medium text-gray-900">
                        {study.lesson_title || getLessonTitle(study.lesson_id)}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {new Date(study.date_completed).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {study.duration_minutes 
                        ? `${Math.floor(study.duration_minutes / 60)}h ${study.duration_minutes % 60}m` 
                        : '—'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {study.location || '—'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <Link 
                        href={`/studies/${study.id}/edit`}
                        className="text-blue-600 hover:text-blue-900"
                      >
                        Edit
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}
    </div>
  );
}