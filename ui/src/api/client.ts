import { AuthResponse, BaseResponse, User } from '../types/api';

const API_BASE = 'http://localhost:8080';

export const apiClient = {
  getHeaders() {
    return {
      'Content-Type': 'application/json',
    };
  },

  async login(name: string, password: string): Promise<AuthResponse> {
    const res = await fetch(`${API_BASE}/login`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ name, password }),
      credentials: 'include'
    });
    return res.json();
  },

  async logout(): Promise<BaseResponse<null>> {
    const res = await fetch(`${API_BASE}/logout`, {
      method: 'POST',
      headers: this.getHeaders(),
      credentials: 'include'
    });
    return res.json();
  },

  async register(name: string, password: string): Promise<AuthResponse> {
    const res = await fetch(`${API_BASE}/register`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ name, email: `${name}@example.com`, password }),
      credentials: 'include'
    });
    return res.json();
  },

  async fetchMe(): Promise<BaseResponse<User>> {
    const res = await fetch(`${API_BASE}/users/me`, {
      method: 'GET',
      headers: this.getHeaders(),
      credentials: 'include'
    });
    return res.json();
  },

  async fetchUsers(): Promise<BaseResponse<User[]>> {
    const res = await fetch(`${API_BASE}/users`, {
      method: 'GET',
      headers: this.getHeaders(),
      credentials: 'include'
    });
    return res.json();
  },

  async fetchContacts(): Promise<BaseResponse<User[]>> {
    const res = await fetch(`${API_BASE}/users/contacts`, {
      method: 'GET',
      headers: this.getHeaders(),
      credentials: 'include'
    });
    return res.json();
  },

  async fetchMessagesBetween(otherUserId: string): Promise<BaseResponse<any[]>> {
    const res = await fetch(`${API_BASE}/messages/${otherUserId}`, {
      method: 'GET',
      headers: this.getHeaders(),
      credentials: 'include'
    });
    return res.json();
  },

  async fetchGroupMessages(groupId: string): Promise<BaseResponse<any[]>> {
    const res = await fetch(`${API_BASE}/messages/group/${groupId}`, {
      method: 'GET',
      headers: this.getHeaders(),
      credentials: 'include'
    });
    return res.json();
  },

  async fetchMyGroups(): Promise<BaseResponse<any[]>> {
    const res = await fetch(`${API_BASE}/groups/my`, {
      method: 'GET',
      headers: this.getHeaders(),
      credentials: 'include'
    });
    return res.json();
  }
};
