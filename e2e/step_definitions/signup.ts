import { Given } from "@cucumber/cucumber"
import { FiderWorld } from "../world"
import { delay } from "./fns/delay"

Given("I go to the signup page", async function (this: FiderWorld) {
  await this.page.goto("https://login.dev.fider.io:3000/signup")
})

Given("I create a new account using my email", async function (this: FiderWorld) {
  const now = new Date().getTime()
  this.tenantName = `feedback${now}`

  await this.page.goto("https://login.dev.fider.io:3000/signup")
  await this.page.type("#input-name", "admin")
  await this.page.type("#input-email", `${this.tenantName}@fider.io`)
  await this.page.type("#input-tenantName", this.tenantName)
  await this.page.type("#input-subdomain", this.tenantName)
  await this.page.check("#input-legalAgreement")
  await this.page.click(".c-button--primary")
})

Given("I activate my account", async function (this: FiderWorld) {
  await delay(1000)

  const response = await fetch(`http://localhost:8025/api/v2/search?kind=to&query=${this.tenantName}@fider.io`)
  const responseBody = await response.json()
  const emailHtml = responseBody.items[0].Content.Body
  const reg = /https:\/\/feedback\d+\.dev\.fider\.io:3000\/signup\/verify\?k=.+?(?=')/gim
  const result = reg.exec(emailHtml)
  if (result) {
    await this.page.goto(result[0])
  }
})
