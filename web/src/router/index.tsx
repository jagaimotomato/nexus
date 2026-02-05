import {
  createBrowserRouter,
  Navigate,
  type RouteObject,
} from "react-router-dom";

import BasicLayout from "@/layouts/BasicLayout";
import AuthGuard from "./authGuard";
import Login from "@/pages/login";
import Dashboard from "@/pages/dashboard";
import NotFoundPage from "@/pages/404";

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
        element: <Navigate to="/dashboard" replace />,
      },
      {
        path: "dashboard",
        element: <Dashboard />,
      },
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

const router = createBrowserRouter(routes);

export default router;
