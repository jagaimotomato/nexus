import { create } from "zustand";

interface UserState {
  token: string;
  userInfo: {
    name: string;
    avatar?: string;
    id: number;
    username: string;
    roles?: any[];
  } | null;
  setUserInfo: (user: any) => void;
  login: (token: string) => void;
  logout: () => void;
}

const useUserStore = create<UserState>((set) => ({
  token: localStorage.getItem("token") || "",
  userInfo: localStorage.getItem("userInfo")
    ? JSON.parse(localStorage.getItem("userInfo")!)
    : null,
  login: (token: string) => {
    localStorage.setItem("token", token);
    set({ token });
  },
  setUserInfo: (user: any) => {
    localStorage.setItem("userInfo", JSON.stringify(user));
    set({ userInfo: user });
  },
  logout: () => {
    localStorage.removeItem("token");
    set({ token: "", userInfo: null });
  },
}));

export default useUserStore;
