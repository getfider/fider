import { Browser, BrowserTab, pageHasLoaded } from "../lib"
import { HomePage } from "../pages"

export interface TestContext {
  browser1: Browser
  browser2: Browser
  tab1: BrowserTab
  tab2: BrowserTab
  tenantName: string
  tenantSubdomain: string
}

export let ctx: TestContext

describe("E2E", () => {
  beforeAll(async () => {
    const now = new Date().getTime()
    const tenantName = `Selenium ${now}`
    const tenantSubdomain = process.env.HOST_MODE === "single" ? "login" : `selenium${now}`
    const baseURL = `https://${tenantSubdomain}.dev.fider.io:3000`
    const browser1 = await Browser.launch()
    const tab1 = await browser1.newTab(baseURL)
    const browser2 = await Browser.launch()
    const tab2 = await browser2.newTab(baseURL)

    ctx = {
      tenantName,
      tenantSubdomain,
      browser1,
      tab1,
      browser2,
      tab2,
    }
  })

  afterAll(async () => {
    await ctx.browser1.close()
    await ctx.browser2.close()
  })

  it("Tab1: User can sign up with facebook", async () => {
    await ctx.tab1.pages.signup.navigate()
    await ctx.tab1.pages.signup.signInWithFacebook()
    await ctx.tab1.pages.facebook.signInAsJonSnow()

    await ctx.tab1.pages.signup.signUpAs(ctx.tenantName, ctx.tenantSubdomain)

    await ctx.tab1.wait(pageHasLoaded(HomePage))
  })

  describe("E2E: Feedback Loop", () => {
    require("./feedback-loop")
  })

  describe("E2E: Admin Settings", () => {
    require("./admin-settings")
  })
})
