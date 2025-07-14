import { Given, Then } from "@cucumber/cucumber"
import { FiderWorld } from "../world"
import expect from "expect"
import { getLatestLinkSentTo } from "./fns"

Given("I go to the home page", async function (this: FiderWorld) {
  await this.page.goto(`https://${this.tenantName}.dev.fider.io:3000/`)
})

Then("I should be on the home page", async function (this: FiderWorld) {
  const container = await this.page.locator("#p-home")
  await expect(container).toBeVisible()
})

Then("I click on the first post", async function (this: FiderWorld) {
  await this.page.click(".c-posts-container__post:first-child a")
})

Then("I search for {string}", async function (this: FiderWorld, searchTerm: string) {
  await this.page.type("#input-query", searchTerm)
})

Given("I type {string} as the title", async function (this: FiderWorld, title: string) {
  await this.page.type("#input-title", title)
})

Given("I type {string} as the description", async function (this: FiderWorld, description: string) {
  const editor = this.page.getByTestId("tiptap-editor")
  // Click to focus, then type
  await editor.click()
  await this.page.keyboard.type(description)
})

Given("I click enter your suggestion", async function () {
  await this.page.click(".p-home__welcome-col .c-button--default")
})

Given("I type my email address", async function (this: FiderWorld) {
  const userEmail = `$user-${this.tenantName}@fider.io`
  await this.page.type("#input-email", userEmail)
})

Given("I click continue with email", async function () {
  await this.page.click(".c-signin-control button[type='submit']")
})

Given("I click submit your feedback", async function () {
  await this.page.click(".c-share-feedback__content .c-button--primary")
})

Given("I click on the confirmation link", async function (this: FiderWorld) {
  const userEmail = `$user-${this.tenantName}@fider.io`
  const activationLink = await getLatestLinkSentTo(userEmail)
  await this.page.goto(activationLink)
})

Then("I should be on the complete profile page", async function (this: FiderWorld) {
  const container = await this.page.$$("#p-complete-profile")
  await expect(container).toBeDefined()
})

Then("I should see the new post modal", async function (this: FiderWorld) {
  const container = await this.page.getByTestId("modal")
  await expect(container).toBeVisible()
})

Given("I enter my name as {string}", async function (this: FiderWorld, name: string) {
  await this.page.type("#input-name", name)
})

Given("I click submit", async function () {
  await this.page.click("button[type='submit']")
})

Then("I should be on the confirmation link page", async function (this: FiderWorld) {
  const userEmail = `$user-${this.tenantName}@fider.io`
  await expect(this.page.getByText(`We have just sent a confirmation link to ${userEmail}`)).toBeVisible()
})

Then("I should see {string} as the draft post title", async function (this: FiderWorld, title: string) {
  const postTitle = await this.page.locator("#input-title").inputValue()
  await expect(postTitle).toBe(title)
})
