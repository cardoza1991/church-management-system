import axios, { AxiosInstance, AxiosResponse } from 'axios';

// Type definitions
export interface User {
  id: number;
  username: string;
  email: string;
  role: string;
  full_name: string;
  phone?: string;
  created_at: string;
  updated_at: string;
}

export interface Contact {
  id: number;
  name: string;
  email?: string;
  phone?: string;
  location?: string;
  notes?: string;
  date_added: string;
  last_updated: string;
  current_status_id: number;
}

export interface Status {
  id: number;
  name: string;
  description?: string;
  display_order: number;
}

export interface StatusHistory {
  id: number;
  contact_id: number;
  status_id: number;
  status_name: string;
  notes?: string;
  date_changed: string;
}

export interface Lesson {
  id: number;
  title: string;
  description?: string;
  sequence_number: number;
  created_at: string;
  updated_at: string;
}

export interface Study {
  id: number;
  contact_id: number;
  lesson_id: number;
  lesson_title?: string;
  date_completed: string;
  location?: string;
  duration_minutes?: number;
  notes?: string;
  taught_by_user_id?: number;
  created_at: string;
  updated_at: string;
}

export interface Room {
  id: number;
  name: string;
  capacity: number;
  location?: string;
  description?: string;
  availability_start?: string;
  availability_end?: string;
  is_available: boolean;
  created_at: string;
  updated_at: string;
}

export interface Reservation {
  id: number;
  room_id: number;
  room_name?: string;
  user_id: number;
  contact_id?: number;
  title: string;
  description?: string;
  start_time: string;
  end_time: string;
  recurring_type: string;
  recurring_end_date?: string;
  created_at: string;
  updated_at: string;
}

// API response types
export interface ContactsResponse {
  contacts: Contact[];
  limit: number;
  offset: number;
}

export interface StatusesResponse {
  statuses: Status[];
}

export interface AuthResponse {
  token: string;
  user: User;
}

// Create base API instance
const baseURL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const api: AxiosInstance = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add a request interceptor to include the auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Add a response interceptor to handle token expiration
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      // Redirect to login page if unauthorized
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;

// Auth API
export const authAPI = {
  login: (username: string, password: string): Promise<AxiosResponse<AuthResponse>> =>
    api.post('/login', { username, password }),
  register: (userData: any): Promise<AxiosResponse<AuthResponse>> =>
    api.post('/register', userData),
  getCurrentUser: (): Promise<AxiosResponse<User>> =>
    api.get('/users/me'),
};

// Contacts API
export const contactsAPI = {
  getContacts: (limit = 20, offset = 0): Promise<AxiosResponse<ContactsResponse>> =>
    api.get(`/contacts?limit=${limit}&offset=${offset}`),
  getContact: (id: number): Promise<AxiosResponse<Contact>> =>
    api.get(`/contacts/${id}`),
  createContact: (contactData: Partial<Contact>): Promise<AxiosResponse<Contact>> =>
    api.post('/contacts', contactData),
  updateContact: (id: number, contactData: Partial<Contact>): Promise<AxiosResponse<Contact>> =>
    api.put(`/contacts/${id}`, contactData),
  updateContactStatus: (id: number, statusId: number, notes: string = ''): Promise<AxiosResponse<Contact>> =>
    api.put(`/contacts/${id}/status`, { status_id: statusId, notes }),
  getContactStatusHistory: (id: number): Promise<AxiosResponse<{contact_id: number, history: StatusHistory[]}>> =>
    api.get(`/contacts/${id}/status-history`),
};

// Statuses API
export const statusesAPI = {
  getStatuses: (): Promise<AxiosResponse<StatusesResponse>> =>
    api.get('/statuses'),
  getStatus: (id: number): Promise<AxiosResponse<Status>> =>
    api.get(`/statuses/${id}`),
};

// Custom API for the study service
export const studiesAPI = {
  getLessons: (): Promise<AxiosResponse<{lessons: Lesson[]}>> => 
    axios.get('http://localhost:8082/lessons'),
  getStudiesByContact: (contactId: number): Promise<AxiosResponse<{contact_id: number, studies: Study[]}>> => 
    axios.get(`http://localhost:8082/contacts/${contactId}/studies`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    }),
  getContactStudyStats: (contactId: number): Promise<AxiosResponse<any>> => 
    axios.get(`http://localhost:8082/contacts/${contactId}/study-stats`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    }),
  createStudy: (studyData: Partial<Study>): Promise<AxiosResponse<Study>> =>
    axios.post('http://localhost:8082/studies', studyData, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    }),
  updateStudy: (id: number, studyData: Partial<Study>): Promise<AxiosResponse<Study>> =>
    axios.put(`http://localhost:8082/studies/${id}`, studyData, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    }),
  deleteStudy: (id: number): Promise<AxiosResponse<any>> =>
    axios.delete(`http://localhost:8082/studies/${id}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
};

// Custom API for the reservation service
export const reservationsAPI = {
  getRooms: (): Promise<AxiosResponse<{rooms: Room[]}>> => 
    axios.get('http://localhost:8083/rooms'),
  getReservationsByDate: (start: string, end: string, roomId: number | null = null): Promise<AxiosResponse<{reservations: Reservation[]}>> => {
    let url = `http://localhost:8083/reservations/by-date?start=${start}&end=${end}`;
    if (roomId) url += `&room_id=${roomId}`;
    
    return axios.get(url, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
  },
  createReservation: (data: Partial<Reservation>): Promise<AxiosResponse<Reservation>> => 
    axios.post('http://localhost:8083/reservations', data, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    }),
  updateReservation: (id: number, data: Partial<Reservation>): Promise<AxiosResponse<Reservation>> => 
    axios.put(`http://localhost:8083/reservations/${id}`, data, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    }),
  deleteReservation: (id: number): Promise<AxiosResponse<any>> =>
    axios.delete(`http://localhost:8083/reservations/${id}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
};