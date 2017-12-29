import { http, Result } from '@fider/services/http';
import { injectable } from '@fider/di';
import { Tenant, UserRole } from '@fider/models';

export interface TenantService {
    create(request: CreateTenantRequest): Promise<Result<CreateTenantResponse>>;
    updateSettings(title: string, invitation: string, welcomeMessage: string, cname: string): Promise<Result>;
    checkAvailability(subdomain: string): Promise<Result<CheckAvailabilityResponse>>;
    signIn(email: string): Promise<Result>;
    completeProfile(key: string, name: string): Promise<Result>;
    changeRole(userId: number, role: UserRole): Promise<Result>;
}

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

@injectable()
export class HttpTenantService implements TenantService {
    public async create(request: CreateTenantRequest): Promise<Result<CreateTenantResponse>> {
        return await http.post<CreateTenantResponse>('/api/tenants', request);
    }

    public async updateSettings(title: string, invitation: string, welcomeMessage: string, cname: string): Promise<Result> {
        return await http.post('/api/admin/settings', {
            title, invitation, welcomeMessage, cname
        });
    }

    public async checkAvailability(subdomain: string): Promise<Result<CheckAvailabilityResponse>> {
        return await http.get<CheckAvailabilityResponse>(`/api/tenants/${subdomain}/availability`);
    }

    public async signIn(email: string): Promise<Result> {
        return await http.post('/api/signin', {
            email,
        });
    }

    public async completeProfile(key: string, name: string): Promise<Result> {
        return await http.post('/api/signin/complete', {
            key,
            name,
        });
    }

    public async changeRole(userId: number, role: UserRole): Promise<Result> {
        return await http.post(`/api/admin/users/${userId}/role`, {
            role,
        });
    }
}
