import { Before, BeforeAll, AfterAll, After } from "@cucumber/cucumber"
import * as playwright from "playwright"
import { FiderWorld } from "./world"

let browser: playwright.Browser
type BrowserName = "chromium" | "firefox" | "webkit"

BeforeAll(async function () {
  const name = (process.env.BROWSER || "chromium") as BrowserName
  browser = await playwright[name].launch({
    headless: true,
    slowMo: 50,
  })
})

AfterAll(async function () {
  await browser.close()
})

Before(async function (this: FiderWorld) {
  this.context = await browser.newContext({
    viewport: { width: 1280, height: 720 },
    ignoreHTTPSErrors: true,
  })
  this.page = await this.context.newPage()
})

After(async function (this: FiderWorld) {
  await this.page.close()
  await this.context.close()
})
