import { create } from 'zustand'

interface UserState {
  token: string;
  userInfo: { name: string; avatar?: string } | null;
  login: (token: string) => void;
  logout: () => void
}

const useUserStore = create<UserState>((set) => ({
  token: localStorage.getItem('token') || '',
  userInfo: null,
  login: (token: string) => {
    localStorage.setItem('token', token);
    set({ token })
  },
  logout: () => {
    localStorage.removeItem('token')
    set({ token: '', userInfo: null })
  }
}))

export default useUserStore