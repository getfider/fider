import { Given, Then } from "@cucumber/cucumber"
import { FiderWorld } from "../world"
import expect from "expect"

Given("I go to the home page", async function (this: FiderWorld) {
  await this.page.goto(`https://${this.tenantName}.dev.fider.io:3000/`)
})

Then("I should be on the home page", async function (this: FiderWorld) {
  const container = await this.page.$$("#p-home")
  expect(container).toBeDefined()
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
  await this.page.type("#input-description", description)
})

Given("I click submit new post", async function () {
  await this.page.click(".p-home__welcome-col .c-button--primary")
})
