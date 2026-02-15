import React, { Suspense } from "react";
import type { RouteObject } from "react-router-dom";
import type { MenuItem } from "@/api/menu";
import NotFoundPage from "@/pages/404";

const pageModules = import.meta.glob("../pages/**/*.tsx");

const resolveComponent = (componentPath?: string) => {
  if (!componentPath) {
    return null;
  }
  const key = `../pages${componentPath}.tsx`;
  const loader = pageModules[key];
  if (!loader) {
    return <NotFoundPage />;
  }
  const LazyComponent = React.lazy(loader as any);
  return (
    <Suspense fallback={null}>
      <LazyComponent />
    </Suspense>
  );
};

export const buildRoutesFromMenus = (menus: MenuItem[]): RouteObject[] => {
  const routes: RouteObject[] = [];
  const walk = (items: MenuItem[]) => {
    items.forEach((item) => {
      if (item.type === 2 && item.path) {
        const element = resolveComponent(item.component);
        if (element) {
          routes.push({
            path: item.path.replace(/^\//, ""),
            element,
          });
        }
      }
      if (item.children && item.children.length > 0) {
        walk(item.children);
      }
    });
  };
  walk(menus);
  return routes;
};

export const buildLayoutMenus = (menus: MenuItem[]) => {
  const walk = (items: MenuItem[]): any[] => {
    return items
      .filter((item) => item.type !== 3)
      .map((item) => ({
        path: item.path,
        name: item.title,
        icon: item.icon,
        hideInMenu: item.hidden,
        routes: item.children ? walk(item.children) : undefined,
      }));
  };
  return walk(menus);
};

export const getFirstMenuPath = (menus: MenuItem[]): string | null => {
  for (const item of menus) {
    if (item.type === 2 && item.path) {
      return item.path;
    }
    if (item.children && item.children.length > 0) {
      const child = getFirstMenuPath(item.children);
      if (child) {
        return child;
      }
    }
  }
  return null;
};
