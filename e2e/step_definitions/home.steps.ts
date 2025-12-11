import { Given, Then } from "@cucumber/cucumber"
import { FiderWorld } from "../world"
import expect from "expect"
import { getLatestCodeSentTo } from "./fns"

Given("I go to the home page", async function (this: FiderWorld) {
  await this.page.goto(`https://${this.tenantName}.dev.fider.io:3000/`)
})

Then("I should be on the home page", async function (this: FiderWorld) {
  const container = await this.page.locator("#p-home")
  await expect(container).toBeVisible()
})

Then("I click on the first post", async function (this: FiderWorld) {
  await this.page.click(".c-posts-container__post-link:first-of-type")
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
  await this.page.click(".p-home__add-idea-btn")
})

Given("I type my email address", async function (this: FiderWorld) {
  const userEmail = `$user-${this.tenantName}@fider.io`
  await this.page.type("#input-email", userEmail)
})

Given("I click continue with email", async function () {
  await this.page.click(".c-signin-control button[type='submit']")
})

Then("I should see the name field", async function (this: FiderWorld) {
  // Wait for the name field to appear
  await this.page.waitForSelector("#input-name", { timeout: 5000 })
  const nameField = await this.page.locator("#input-name")
  await expect(nameField).toBeVisible()
})

Given("I click continue", async function () {
  await this.page.getByRole("button", { name: "Sign up" }).click()
})

Given("I click submit your feedback", async function () {
  await this.page.click(".c-share-feedback__content .c-button--primary")
})

Then("I should be on the confirmation code page", async function (this: FiderWorld) {
  const userEmail = `$user-${this.tenantName}@fider.io`
  // Wait for the code entry field to appear
  await this.page.waitForSelector("#input-code", { timeout: 5000 })
  // Check for code entry instruction message
  await expect(this.page.getByText(`Please type in the code we just sent to ${userEmail}`)).toBeVisible()
})

Given("I enter the confirmation code", async function (this: FiderWorld) {
  const userEmail = `$user-${this.tenantName}@fider.io`
  const code = await getLatestCodeSentTo(userEmail)

  // Enter the code in the UI
  await this.page.fill("#input-code", code)
  await this.page.getByRole("button", { name: "submit" }).click()

  // Wait for navigation after successful code verification
  await this.page.waitForLoadState("networkidle")
})

Then("I should see the new post modal", async function (this: FiderWorld) {
  const container = await this.page.getByTestId("modal")
  await expect(container).toBeVisible()
})

Given("I enter my name as {string}", async function (this: FiderWorld, name: string) {
  await this.page.fill("#input-name", name)
})

Given("I click submit", async function () {
  await this.page.getByRole("button", { name: "submit" }).click()
})

Then("I should see {string} as the draft post title", async function (this: FiderWorld, title: string) {
  const postTitle = await this.page.locator("#input-title").inputValue()
  await expect(postTitle).toBe(title)
})
