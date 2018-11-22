import { http, Result } from "@fider/services";
import { User } from "@fider/models"

export const get = async (prefix: string): Promise<Result<User[]>> => {
    return http.get<User[]>(`/_api/mentions?query=${prefix}`);
};
  