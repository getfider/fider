import { Page } from "playwright"

export function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

export async function isAuthenticated(page: Page): Promise<boolean> {
  const serverData = JSON.parse(await page.innerText("#server-data"))
  return serverData.user !== undefined
}

// On E2E test, every user is created as {userName}-{tenantName}
export async function isAuthenticatedAsUser(page: Page, userName: string): Promise<boolean> {
  const serverData = JSON.parse(await page.innerText("#server-data"))
  return serverData.user ? serverData.email.startsWith(userName) : false
}

export async function getLatestLinkSentTo(address: string): Promise<string> {
  await delay(1000)

  const response = await fetch(`http://localhost:8025/api/v2/search?kind=to&query=${address}`)
  const responseBody = await response.json()
  const emailHtml = responseBody.items[0].Content.Body
  const reg = /https:\/\/feedback\d+\.dev\.fider\.io:3000\/(.*)verify\?k=.+?(?=')/gim
  const result = reg.exec(emailHtml)
  if (!result) {
    throw new Error("Could not find a link in email content.")
  }

  return result[0]
}
