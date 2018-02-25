import { http, Result } from '@fider/services';

export const getTotalUnreadNotifications = async (): Promise<Result<number>> => {
  return http.get<{ total: number}>('/api/notifications/unread/total').then((result) => {
    return {
      ok: result.ok,
      error: result.error,
      data: result.data ? result.data.total : 0,
    };
  });
};
