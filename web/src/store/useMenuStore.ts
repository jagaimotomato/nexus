import { create } from "zustand";
import { getUserMenuTree, type MenuItem } from "@/api/menu";

interface MenuState {
  menus: MenuItem[];
  loading: boolean;
  setMenus: (menus: MenuItem[]) => void;
  fetchMenus: () => Promise<void>;
  resetMenus: () => void;
}

const useMenuStore = create<MenuState>((set) => ({
  menus: [],
  loading: false,
  setMenus: (menus) => set({ menus }),
  resetMenus: () => set({ menus: [] }),
  fetchMenus: async () => {
    set({ loading: true });
    try {
      const menus = await getUserMenuTree();
      set({ menus });
    } finally {
      set({ loading: false });
    }
  },
}));

export default useMenuStore;
