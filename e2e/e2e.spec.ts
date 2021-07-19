import { test, expect } from "@playwright/test"
import debug from "debug"
import "isomorphic-fetch"
const log = debug("e2e")

//TODO: rewrite to steps when available https://github.com/microsoft/playwright/issues/7254

test("create a new site using email", async ({ page }) => {
  const now = new Date().getTime()
  const tenantName = `feedback${now}`
  log(`New Tenant › ${tenantName}`)

  // Create a new site
  await page.goto("https://login.dev.fider.io:3000/signup")
  await page.type("#input-name", "admin")
  await page.type("#input-email", `${tenantName}@fider.io`)
  await page.type("#input-tenantName", tenantName)
  await page.type("#input-subdomain", tenantName)
  await page.check("#input-legalAgreement")
  await page.click(".c-button--primary")

  // Verify pending state
  await page.goto(`https://${tenantName}.dev.fider.io:3000`)
  expect(await page.innerText("h1")).toBe("YOUR ACCOUNT IS PENDING ACTIVATION")

  // Find activation link
  const response = await fetch(`http://localhost:8026/api/v2/search?kind=to&query=${tenantName}@fider.io`)
  const responseBody = await response.json()
  const emailHtml = responseBody.items[0].Content.Body
  const reg = /https:\/\/feedback\d+\.dev\.fider\.io:3000\/signup\/verify\?k=.+?(?=')/gim
  const result = reg.exec(emailHtml)
  log(`Found activation link › ${result[0]}`)

  // Activate site
  await page.goto(result[0])
  expect(await page.innerText("h1.text-title")).toBe(tenantName)

  // Log out
  await page.click(".c-menu-user .c-dropdown__handle")
  await page.click("a[href='/signout?redirect=/']")

  expect(await page.isVisible(".c-notification-indicator")).toBe(false)
})
