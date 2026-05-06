import { http, Result } from "@fider/services/http"

export interface EmailDomainRule {
  id: number
  domain: string
  ruleType: "deny" | "allow"
  createdAt: string
}

export interface EmailDomainRulesResponse {
  deny: EmailDomainRule[]
  allow: EmailDomainRule[]
}

export const listEmailDomainRules = async (): Promise<Result<EmailDomainRulesResponse>> => {
  return await http.get<EmailDomainRulesResponse>("/_api/admin/email-domain-rules")
}

export const addEmailDomainRule = async (domain: string, ruleType: "deny" | "allow"): Promise<Result<EmailDomainRule>> => {
  return await http.post<EmailDomainRule>("/_api/admin/email-domain-rules", { domain, ruleType })
}

export const deleteEmailDomainRule = async (id: number): Promise<Result> => {
  return await http.delete(`/_api/admin/email-domain-rules/${id}`)
}

export interface DisposableUserRow {
  userID: number
  name: string
  email: string
  voteCount: number
  postCount: number
  commentCount: number
}

export interface DisposableUsersResponse {
  total: number
  users: DisposableUserRow[]
}

export const listDisposableUsers = async (): Promise<Result<DisposableUsersResponse>> => {
  return await http.get<DisposableUsersResponse>("/_api/admin/disposable-users")
}

export const bulkDeleteDisposableUsers = async (userIds: number[]): Promise<Result<{ deleted: number }>> => {
  return await http.post<{ deleted: number }>("/_api/admin/disposable-users/delete", { userIds })
}

export const updateBlockDisposableEmails = async (block: boolean): Promise<Result> => {
  return await http.post("/_api/admin/settings/disposable-emails", { blockDisposableEmails: block })
}
