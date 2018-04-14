import { http, Result } from "@fider/services";

export const sendSampleInvite = async (subject: string, message: string): Promise<Result> => {
  return http.post("/api/admin/invitations/sample", { subject, message }).then(http.event("invite", "sample"));
};
