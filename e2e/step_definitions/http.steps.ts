import { Given, Then, When } from "@cucumber/cucumber"
import { FiderWorld } from "e2e/world"
import expect from "expect"

let requestUrl: string
let requestOpts: RequestInit
let requestHeaders: Headers

let response: Response
let responseBody: string

const getFullUrl = (world: FiderWorld, url: string): string => {
  return url.startsWith("/") ? `https://${world.tenantName}.dev.fider.io:3000${url}` : url
}

Given("I prepare a {string} request to {string}", async function (this: FiderWorld, method: string, url: string) {
  requestUrl = getFullUrl(this, url)
  requestHeaders = new Headers()
  requestOpts = {
    method,
    headers: {},
  }
})

Given("I set the {string} header to {string}", async function (headerName: string, headerValue: string) {
  requestHeaders.set(headerName, headerValue)
})

Given("I send a {string} request to {string}", async function (this: FiderWorld, method: string, url: string) {
  response = await fetch(getFullUrl(this, url), { method })
  responseBody = await response.text()
})

When("I send the request", async function () {
  response = await fetch(requestUrl, { ...requestOpts, headers: requestHeaders })
  responseBody = await response.text()
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

Then("I should see {string} on the response body", async function (substring: string) {
  expect(responseBody).toContain(substring)
})

Then("I should not see {string} on the response body", async function (substring: string) {
  expect(responseBody).not.toContain(substring)
})
