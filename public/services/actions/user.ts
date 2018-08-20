import { http, Result } from "@fider/services/http";
import { UserSettings } from "@fider/models";

export const updateUserSettings = async (name: string, settings: UserSettings): Promise<Result> => {
  return await http.post("/_api/user/settings", {
    name,
    settings
  });
};

export const changeUserEmail = async (email: string): Promise<Result> => {
  return await http.post("/_api/user/change-email", {
    email
  });
};

export const deleteCurrentAccount = async (): Promise<Result> => {
  return await http.delete("/_api/user");
};

export const regenerateAPIKey = async (): Promise<Result<{ apiKey: string }>> => {
  return await http.post<{ apiKey: string }>("/_api/user/regenerate-apikey");
};
