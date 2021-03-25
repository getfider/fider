import { http, Result } from "@fider/services"

export const sendInvites = async (subject: string, message: string, recipients: string[]): Promise<Result> => {
  return http.post("/api/v1/invitations/send", { subject, message, recipients }).then(http.event("invite", "send"))
}

export const sendSampleInvite = async (subject: string, message: string): Promise<Result> => {
  return http.post("/api/v1/invitations/sample", { subject, message }).then(http.event("invite", "sample"))
}
