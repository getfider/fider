import { Then } from "@cucumber/cucumber"
import { FiderWorld } from "../world"
import { expect } from "@playwright/test"

Then("I should be on the show post page", async function (this: FiderWorld) {
  const container = await this.page.$$(".p-show-post")
  expect(container).toBeDefined()
})

Then("I should see {string} as the post title", async function (this: FiderWorld, title: string) {
  const postTitle = await this.page.innerText(".p-show-post__title")
  expect(postTitle).toBe(title)
})

Then("I should see {int} vote\\(s)", async function (this: FiderWorld, voteCount: number) {
  // Look for the vote count number within the post detail view
  await expect(this.page.locator(".p-show-post .text-2xl").filter({ hasText: voteCount.toString() })).toBeVisible()
})
