import { create } from 'zustand';
import { apiClient } from '../api/client';

interface AuthState {
  username: string;
  userId: string;
  isLoggedIn: boolean;
  login: (username: string, userId: string) => void;
  logout: () => Promise<void>;
  hydrate: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  username: '',
  userId: '',
  isLoggedIn: false,

  login: (username, userId) => set({ username, userId, isLoggedIn: true }),

  logout: async () => {
    try { await apiClient.logout(); } catch {}
    set({ username: '', userId: '', isLoggedIn: false });
  },

  hydrate: async () => {
    try {
      const res = await apiClient.fetchMe();
      if (res.status === 'success' && res.data) {
        set({ username: res.data.name, userId: res.data.id, isLoggedIn: true });
      }
    } catch {}
  },
}));
