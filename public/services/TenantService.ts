import { get, post, Result } from '@fider/services/http';
import { injectable } from '@fider/di';
import { Tenant } from '@fider/models';

export interface TenantService {
    create(token?: string, name?: string, subdomain?: string): Promise<Result>;
    updateSettings(title: string, invitation: string, welcomeMessage: string): Promise<Result>;
    checkAvailability(subdomain: string): Promise<Result<CheckAvailabilityResponse>>;
    signIn(email: string): Promise<Result>;
    completeProfile(key: string, name: string): Promise<Result>;
}

export interface CheckAvailabilityResponse {
    message: string;
}

@injectable()
export class HttpTenantService implements TenantService {
    public async create(token?: string, name?: string, subdomain?: string): Promise<Result> {
        return await post('/api/tenants', {
            token, name, subdomain,
        });
    }

    public async updateSettings(title: string, invitation: string, welcomeMessage: string): Promise<Result> {
        return await post('/api/settings', {
            title, invitation, welcomeMessage,
        });
    }

    public async checkAvailability(subdomain: string): Promise<Result<CheckAvailabilityResponse>> {
        return await get<CheckAvailabilityResponse>(`/api/tenants/${subdomain}/availability`);
    }

    public async signIn(email: string): Promise<Result> {
        return await post('/api/signin', {
            email,
        });
    }

    public async completeProfile(key: string, name: string): Promise<Result> {
        return await post('/api/signin/complete', {
            key,
            name,
        });
    }
}
