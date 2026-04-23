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
  const emailInput = this.page.locator(".c-signin-control #input-email")
  const signInButton = this.page.getByRole("button", { name: /sign in/i })

  if (!(await emailInput.isVisible({ timeout: 2000 }).catch(() => false))) {
    if (await signInButton.isVisible({ timeout: 5000 }).catch(() => false)) {
      await signInButton.click()
    } else {
      await this.page.locator("#c-header .c-button--primary").first().click()
    }
  }

  await this.page.fill(".c-signin-control #input-email", userEmail)
  await this.page.click(".c-signin-control .c-button--primary")

  // Get the code from email and enter it
  const code = await getLatestCodeSentTo(userEmail)
  await this.page.fill("#input-code", code)
  await this.page.getByRole("button", { name: "submit" }).click()

  // Wait for signed-in home UI instead of waiting for network idle.
  await this.page.getByRole("button", { name: /enter your suggestion/i }).waitFor({ state: "visible" })
})
