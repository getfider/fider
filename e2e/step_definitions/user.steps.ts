import { Given } from "@cucumber/cucumber"
import { FiderWorld } from "e2e/world"
import { getLatestLinkSentTo, isAuthenticated, isAuthenticatedAsUser } from "./fns"

Given("I sign in as {string}", async function (this: FiderWorld, userName: string) {
  if (await isAuthenticatedAsUser(this.page, userName)) {
    return
  }

  if (await isAuthenticated(this.page)) {
    await this.page.click(".c-menu-user .c-dropdown__handle")
    await this.page.click("a[href='/signout?redirect=/']")
  }

  const userEmail = `${userName}-${this.tenantName}@fider.io`
  await this.page.click(".c-menu .uppercase.text-sm")
  await this.page.type(".c-signin-control #input-email", userEmail)
  await this.page.click(".c-signin-control .c-button--primary")

  const activationLink = await getLatestLinkSentTo(userEmail)
  await this.page.goto(activationLink)
})
