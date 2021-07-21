import { Given, Then } from "@cucumber/cucumber"
import { expect } from "@playwright/test"
import { FiderWorld } from "../world"

Given("I go to the home page", async function (this: FiderWorld) {
  await this.page.goto(`https://${this.tenantName}.dev.fider.io:3000/`)
})

Then("I should be on the home page", async function (this: FiderWorld) {
  const container = await this.page.$$("#p-home")
  expect(container).toBeDefined()
})
