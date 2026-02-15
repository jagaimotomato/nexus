import {
  createBrowserRouter,
  Navigate,
  type RouteObject,
} from "react-router-dom";

import BasicLayout from "@/layouts/BasicLayout";
import AuthGuard from "./authGuard";
import Login from "@/pages/login";
import NotFoundPage from "@/pages/404";
import type { MenuItem } from "@/api/menu";
import {
  buildRoutesFromMenus,
  getFirstMenuPath,
} from "@/router/routeUtils";

export const createAppRouter = (menus: MenuItem[]) => {
  const firstPath = getFirstMenuPath(menus) || "/404";
  const dynamicRoutes = buildRoutesFromMenus(menus);

  const routes: RouteObject[] = [
    {
      path: "/login",
      element: <Login />,
    },
    {
      path: "/",
      element: (
        <AuthGuard>
          <BasicLayout />
        </AuthGuard>
      ),
      children: [
        {
          index: true,
          element: <Navigate to={firstPath} replace />,
        },
        ...dynamicRoutes,
      ],
    },
    {
      path: "/404",
      element: <NotFoundPage />,
    },
    {
      path: "*",
      element: <Navigate to="/404" replace />,
    },
  ];

  return createBrowserRouter(routes);
};
