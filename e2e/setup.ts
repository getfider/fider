import { Before, BeforeAll, AfterAll, After } from "@cucumber/cucumber"
import debug from "debug"
import * as playwright from "playwright"
import { getLatestLinkSentTo } from "./step_definitions/fns"
import { FiderWorld } from "./world"

let browser: playwright.Browser
let tenantName: string
type BrowserName = "chromium" | "firefox" | "webkit"

BeforeAll({ timeout: 30 * 1000 }, async function () {
  const name = (process.env.BROWSER || "chromium") as BrowserName
  browser = await playwright[name].launch({
    headless: true,
    slowMo: 10,
  })

  if (!tenantName) {
    const now = new Date().getTime()
    tenantName = `feedback${now}`
    await createNewSite()
  }
})

AfterAll(async function () {
  await browser.close()
})

Before(async function (this: FiderWorld) {
  const context = await browser.newContext({
    viewport: { width: 1280, height: 720 },
    ignoreHTTPSErrors: true,
  })

  this.page = await context.newPage()
  this.tenantName = tenantName
  this.log = debug("e2e")
})

After(async function (this: FiderWorld) {
  await this.page.close()
})

async function createNewSite() {
  const context = await browser.newContext({
    viewport: { width: 1280, height: 720 },
    ignoreHTTPSErrors: true,
  })
  const page = await context.newPage()

  const adminEmail = `admin-${tenantName}@fider.io`
  //Create site
  await page.goto("https://login.dev.fider.io:3000/signup")
  await page.type("#input-name", "admin")
  await page.type("#input-email", adminEmail)
  await page.type("#input-tenantName", tenantName)
  await page.type("#input-subdomain", tenantName)
  await page.check("#input-legalAgreement")
  await page.click(".c-button--primary")

  //Activate site
  const activationLink = await getLatestLinkSentTo(adminEmail)
  await page.goto(activationLink)
  await page.close()
}
