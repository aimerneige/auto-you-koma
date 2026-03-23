import { create } from 'zustand'

interface AuthState {
  token: string | null;
  userId: string | null;
  email: string | null;
  setAuth: (token: string, userId: string, email: string) => void;
  clearAuth: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem('token'),
  userId: localStorage.getItem('user_id'),
  email: localStorage.getItem('email'),
  setAuth: (token, userId, email) => {
    localStorage.setItem('token', token);
    localStorage.setItem('user_id', userId);
    localStorage.setItem('email', email);
    set({ token, userId, email });
  },
  clearAuth: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user_id');
    localStorage.removeItem('email');
    set({ token: null, userId: null, email: null });
  },
}));
