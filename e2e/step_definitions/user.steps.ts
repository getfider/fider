import { Given, Then } from "@cucumber/cucumber"
import { expect } from "@playwright/test"
import { FiderWorld } from "../world"

Given("I expand the user menu", async function (this: FiderWorld) {
  await this.page.click(".c-menu-user .c-dropdown__handle")
})

Given("I click on sign out", async function (this: FiderWorld) {
  await this.page.click("a[href='/signout?redirect=/']")
})

Then("I should be logged in", async function (this: FiderWorld) {
  expect(await this.page.isVisible(".c-notification-indicator")).toBe(true)
})

Then("I should not be logged in", async function (this: FiderWorld) {
  expect(await this.page.isVisible(".c-notification-indicator")).toBe(false)
})
