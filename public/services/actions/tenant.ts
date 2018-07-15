import { http, Result } from "@fider/services/http";
import { Tenant, UserRole, OAuthConfig } from "@fider/models";

export interface CheckAvailabilityResponse {
  message: string;
}

export interface CreateTenantRequest {
  legalAgreement: boolean;
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

export interface UpdateTenantSettingsRequest {
  logo?: {
    upload?: {
      content?: string;
      contentType?: string;
    };
    remove: boolean;
  };
  title: string;
  invitation: string;
  welcomeMessage: string;
  cname: string;
}

export const updateTenantSettings = async (request: UpdateTenantSettingsRequest): Promise<Result> => {
  return await http.post("/api/admin/settings/general", request);
};

export const updateTenantAdvancedSettings = async (customCSS: string): Promise<Result> => {
  return await http.post("/api/admin/settings/advanced", { customCSS });
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

export const getOAuthConfig = async (provider: string): Promise<Result<OAuthConfig>> => {
  return await http.get<OAuthConfig>(`/api/admin/oauth/${provider}`);
};

export interface CreateEditOAuthConfigRequest {
  provider: string;
  displayName: string;
  clientId: string;
  clientSecret: string;
  authorizeUrl: string;
  tokenUrl: string;
  scope: string;
  profileUrl: string;
  jsonUserIdPath: string;
  jsonUserNamePath: string;
  jsonUserEmailPath: string;
}

export const saveOAuthConfig = async (request: CreateEditOAuthConfigRequest): Promise<Result> => {
  return await http.post("/api/admin/oauth", request);
};
