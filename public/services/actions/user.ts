import { http, Result } from '@fider/services/http';
import { UserSettings } from '@fider/models';

export const updateUserSettings = async (name: string, settings: UserSettings): Promise<Result> => {
  return await http.post('/api/user/settings', {
    name,
    settings,
  });
};

export const changeUserEmail = async (email: string): Promise<Result> => {
  return await http.post('/api/user/change-email', {
    email,
  });
};
