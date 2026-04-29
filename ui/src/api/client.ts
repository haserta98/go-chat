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
  },

  async addContact(contactId: string): Promise<BaseResponse<null>> {
    const res = await fetch(`${API_BASE}/users/contacts`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ contact_id: contactId }),
      credentials: 'include'
    });
    return res.json();
  },

  async removeContact(contactId: string): Promise<BaseResponse<null>> {
    const res = await fetch(`${API_BASE}/users/contacts`, {
      method: 'DELETE',
      headers: this.getHeaders(),
      body: JSON.stringify({ contact_id: contactId }),
      credentials: 'include'
    });
    return res.json();
  },

  async fetchContactsOnlineStatus(): Promise<BaseResponse<Record<string, boolean>>> {
    const res = await fetch(`${API_BASE}/users/contacts/online`, {
      method: 'GET',
      headers: this.getHeaders(),
      credentials: 'include'
    });
    return res.json();
  },

  async createGroup(name: string): Promise<BaseResponse<any>> {
    const res = await fetch(`${API_BASE}/groups`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ name }),
      credentials: 'include'
    });
    return res.json();
  },

  async fetchAllGroups(): Promise<BaseResponse<any[]>> {
    const res = await fetch(`${API_BASE}/groups`, {
      method: 'GET',
      headers: this.getHeaders(),
      credentials: 'include'
    });
    return res.json();
  },

  async joinGroup(groupId: string): Promise<BaseResponse<null>> {
    const res = await fetch(`${API_BASE}/groups/${groupId}/join`, {
      method: 'POST',
      headers: this.getHeaders(),
      credentials: 'include'
    });
    return res.json();
  },

  async leaveGroup(groupId: string): Promise<BaseResponse<null>> {
    const res = await fetch(`${API_BASE}/groups/${groupId}/leave`, {
      method: 'POST',
      headers: this.getHeaders(),
      credentials: 'include'
    });
    return res.json();
  }
};
