import { Given } from "@cucumber/cucumber"
import { FiderWorld } from "e2e/world"
import { getLatestCodeSentTo, isAuthenticated, isAuthenticatedAsUser } from "./fns"

Given("I sign in as {string}", async function (this: FiderWorld, userName: string) {
  if (await isAuthenticatedAsUser(this.page, userName)) {
    return
  }

  if (await isAuthenticated(this.page)) {
    await this.page.click(".c-menu-user .c-dropdown__handle")
    await this.page.click("a[href='/signout']")
  }

  const userEmail = `${userName}-${this.tenantName}@fider.io`
  await this.page.getByRole("button", { name: "Sign in" }).click()
  await this.page.fill(".c-signin-control #input-email", userEmail)
  await this.page.click(".c-signin-control .c-button--primary")

  // Get the code from email and enter it
  const code = await getLatestCodeSentTo(userEmail)
  await this.page.fill("#input-code", code)
  await this.page.getByRole("button", { name: "submit" }).click()

  // Verifying the code closes the modal and reloads the page. Wait for the signed-in page itself:
  // waitForLoadState resolves immediately on the old page, letting the next step click it just before it's replaced
  await this.page.locator(".c-menu-user").waitFor()
})
