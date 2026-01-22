import React from "react";
import { createBrowserRouter, Navigate } from "react-router";
import type { RouteObject } from "react-router";
import BasicLayout from "@/layouts/BasicLayout";
import AuthGuard from "./authGuard";
import Login from "@/pages/login/index";

const routes: RouteObject[] = [
  {
    path: "/login",
    element: <Login />,
  },
  {
    path: "/",
    element: <BasicLayout />,
    children: [
      {
        path: "/",
        element: <Navigate to="/dashboard" replace />,
      },
    ],
  },
  {
    path: "*",
    element: <Navigate to="/404" replace />,
  },
];
