import { useEffect, useMemo } from "react";
import { RouterProvider } from "react-router-dom";
import { ConfigProvider } from "antd";
import { createAppRouter } from "./router";
import useUserStore from "@/store/useUserStore";
import useMenuStore from "@/store/useMenuStore";

function App() {
  const token = useUserStore((state) => state.token);
  const menus = useMenuStore((state) => state.menus);
  const fetchMenus = useMenuStore((state) => state.fetchMenus);
  const resetMenus = useMenuStore((state) => state.resetMenus);

  useEffect(() => {
    if (token) {
      fetchMenus();
      return;
    }
    resetMenus();
  }, [token, fetchMenus, resetMenus]);

  const router = useMemo(() => createAppRouter(menus), [menus]);

  return (
    <ConfigProvider>
      <RouterProvider router={router} />
    </ConfigProvider>
  );
}

export default App;

