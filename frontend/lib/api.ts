import { getCookie } from 'cookies-next';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

async function request(endpoint: string, options: RequestInit = {}) {
  const token = getCookie('auth_token');
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
    ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
  };

  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Request failed with status ' + response.status }));
    throw new Error(errorData.error || 'An unknown error occurred');
  }

  return response.json();
}

// Authentication
export const login = (data: any) => request('/api/v1/auth/login', { method: 'POST', body: JSON.stringify(data) });
export const register = (data: any) => request('/api/v1/auth/register', { method: 'POST', body: JSON.stringify(data) });

// Dashboard & Analytics
export const getDashboardData = () => request('/api/v1/analytics/dashboard');

// Videos
export const getVideos = (channelId: string) => request(`/api/v1/videos?channel_id=${channelId}`);
export const createVideo = (data: any) => request('/api/v1/videos', { method: 'POST', body: JSON.stringify(data) });
export const processVideo = (videoId: string) => request(`/api/v1/videos/${videoId}/process`, { method: 'POST' });

// Channels
export const getChannels = () => request('/api/v1/channels');
export const createChannel = (data: any) => request('/api/v1/channels', { method: 'POST', body: JSON.stringify(data) });

// AI Agents
export const getAgents = () => request('/api/v1/ai/agents');
export const createTask = (agentId: string, data: any) => request(`/api/v1/ai/agents/${agentId}/task`, { method: 'POST', body: JSON.stringify(data) });
export const getTasks = () => request('/api/v1/ai/tasks');