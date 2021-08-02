import { Given, Then } from "@cucumber/cucumber"
import expect from "expect"

let response: Response

Given("I send a {string} request to {string}", async function (method: string, path: string) {
  response = await fetch(`https://login.dev.fider.io:3000${path}`, { method })
})

Then("I should see http status {int}", async function (code: number) {
  expect(response.status).toEqual(code)
})

Then("I should not see a {string} header", async function (headerName: string) {
  expect(response.headers.get(headerName)).toBeNull()
})

Then("I should see a {string} header with value {string}", async function (headerName: string, headerValue: string) {
  expect(response.headers.get(headerName)).toBe(headerValue)
})
