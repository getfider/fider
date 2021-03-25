import puppeteer from "puppeteer"
import { BrowserTab } from "."

export class Browser {
  private browser: puppeteer.Browser

  public constructor(browser: puppeteer.Browser) {
    this.browser = browser
  }

  public static async launch(): Promise<Browser> {
    const browser = await puppeteer.launch({
      headless: true,
      devtools: false,
      ignoreHTTPSErrors: true,
    })
    return new Browser(browser)
  }

  public async newTab(baseURL: string): Promise<BrowserTab> {
    const page = await this.browser.newPage()
    return new BrowserTab(page, baseURL)
  }

  public async close(): Promise<void> {
    await this.browser.close()
  }
}
