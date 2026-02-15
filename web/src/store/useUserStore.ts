import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

interface UserInfo {
  id: number;
  username: string;
  name: string;
  avatar: string;
  roles: Array<{
    id?: number;
    name?: string;
    key?: string;
  }>;
}

interface UserState {
  token: string;
  userInfo: UserInfo | null;
  setUserInfo: (user: UserInfo) => void;
  login: (token: string) => void;
  logout: () => void;
}

const useUserStore = create<UserState>()(
  persist(
    (set) => ({
      token: "",
      userInfo: null,
      login: (token: string) => set({ token }),
      logout: () => set({ token: "", userInfo: null }),
      setUserInfo: (user: UserInfo) => set({ userInfo: user }),
    }),
    {
      name: "nexus-storage",
      storage: createJSONStorage(() => localStorage),
    }
  )
);

export default useUserStore;
