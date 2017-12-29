import { http, Result } from '@fider/services/http';
import { injectable } from '@fider/di';
import { Tenant } from '@fider/models';

export interface UserService {
    updateSettings(name: string): Promise<Result>;
}

@injectable()
export class HttpUserService implements UserService {
    public async updateSettings(name: string): Promise<Result> {
        return await http.post('/api/user/settings', {
          name,
        });
    }
}
