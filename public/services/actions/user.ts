import { http, Result } from '@fider/services/http';

export const updateUserSettings = async (name: string): Promise<Result> => {
  return await http.post('/api/user/settings', {
    name,
  });
};

export const changeUserEmail = async (email: string): Promise<Result> => {
  return await http.post('/api/user/change-email', {
    email,
  });
};
