import { http, Result } from "@fider/services/http";
import { Tenant, UserRole } from "@fider/models";

export interface CheckAvailabilityResponse {
  message: string;
}

export interface CreateTenantRequest {
  tenantName: string;
  subdomain?: string;
  name?: string;
  token?: string;
  email?: string;
}

export interface CreateTenantResponse {
  token?: string;
}

export const createTenant = async (request: CreateTenantRequest): Promise<Result<CreateTenantResponse>> => {
  return await http.post<CreateTenantResponse>("/api/tenants", request);
};

export const updateTenantSettings = async (
  title: string,
  invitation: string,
  welcomeMessage: string,
  cname: string
): Promise<Result> => {
  return await http.post("/api/admin/settings/general", {
    title,
    invitation,
    welcomeMessage,
    cname
  });
};

export const updateTenantPrivacy = async (isPrivate: boolean): Promise<Result> => {
  return await http.post("/api/admin/settings/privacy", {
    isPrivate
  });
};

export const checkAvailability = async (subdomain: string): Promise<Result<CheckAvailabilityResponse>> => {
  return await http.get<CheckAvailabilityResponse>(`/api/tenants/${subdomain}/availability`);
};

export const signIn = async (email: string): Promise<Result> => {
  return await http.post("/api/signin", {
    email
  });
};

export const completeProfile = async (key: string, name: string): Promise<Result> => {
  return await http.post("/api/signin/complete", {
    key,
    name
  });
};

export const changeUserRole = async (userId: number, role: UserRole): Promise<Result> => {
  return await http.post(`/api/admin/users/${userId}/role`, {
    role
  });
};
