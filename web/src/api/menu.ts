import request from "@/utils/request";

export interface MenuItem {
  id: number;
  pid: number;
  name: string;
  path: string;
  component: string;
  icon: string;
  sort: number;
  type: 1 | 2 | 3; // 1: 目录, 2: 菜单, 3: 按钮
  hidden: boolean;
  keepAlive: boolean;
  perms: string;
  redirect: string;
  children?: MenuItem[];
  createdAt?: string;
  updatedAt?: string;
}

export const getMenuList = () => {
  return request.get<MenuItem[]>("/menus");
};

export const getUserMenuTree = () => {
  return request.get<MenuItem[]>("/menus/user");
};

export const createMenu = (data: Partial<MenuItem>) => {
  return request.post("/menus", data);
};

export const updateMenu = (id: number, data: Partial<MenuItem>) => {
  return request.put(`/menus/${id}`, data);
};

export const deleteMenu = (id: number) => {
  return request.delete(`/menus/${id}`);
};
