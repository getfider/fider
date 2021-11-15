import { Then } from "@cucumber/cucumber"
import { FiderWorld } from "../world"
import expect from "expect"

Then("I should be on the show post page", async function (this: FiderWorld) {
  const container = await this.page.$$("#p-show-post")
  expect(container).toBeDefined()
})

Then("I should see {string} as the post title", async function (this: FiderWorld, title: string) {
  const postTitle = await this.page.innerText("#p-show-post h1")
  expect(postTitle).toBe(title)
})

Then("I should see {int} vote\\(s)", async function (this: FiderWorld, voteCount: number) {
  const postVoteCount = await this.page.innerText("#p-show-post .c-vote-counter__button")
  expect(postVoteCount).toBe(voteCount.toString())
})
