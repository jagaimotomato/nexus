import { createBrowserRouter, Navigate, type RouteObject } from 'react-router-dom';

import BasicLayout from '@/layouts/BasicLayout';
import AuthGuard from './authGuard';
import Login from '@/pages/login';
import DashboardPage from '@/pages/dashboard';
import NotFoundPage from '@/pages/404';

const routes: RouteObject[] = [
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/',
    element: (
      <AuthGuard>
        <BasicLayout />
      </AuthGuard>
    ),
    children: [
      {
        path: '/',
        element: <Navigate to="/dashboard" replace />,
      },
      {
        path: '/dashboard',
        element: <DashboardPage />,
      },
    ],
  },
  {
    path: '/404',
    element: <NotFoundPage />,
  },
  {
    path: '*',
    element: <Navigate to="/404" replace />,
  },
];

const router = createBrowserRouter(routes);

export default router;
