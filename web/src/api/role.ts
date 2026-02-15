import request from "@/utils/request";

export interface RoleItem {
  id: number;
  name: string;
  key: string;
  sort: number;
  status: number;
}

export interface RoleDetail extends RoleItem {
  menuIds: number[];
}

export interface RolePayload {
  name: string;
  key: string;
  sort: number;
  status: number;
  menuIds: number[];
}

export const getRoleList = () => {
  return request.get<RoleItem[]>("/roles");
};

export const getRoleDetail = (id: number) => {
  return request.get<RoleDetail>(`/roles/${id}`);
};

export const createRole = (data: RolePayload) => {
  return request.post("/roles", data);
};

export const updateRole = (id: number, data: RolePayload) => {
  return request.put(`/roles/${id}`, data);
};

export const deleteRole = (id: number) => {
  return request.delete(`/roles/${id}`);
};
