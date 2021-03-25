import * as puppeteer from "puppeteer"
import { AllPages } from "../pages"
import { WaitCondition, Page, NewablePage } from "."
import { pageHasLoaded } from "./conditions"

export class BrowserTab {
  private page: puppeteer.Page
  public pages: AllPages
  public baseURL: string

  public constructor(page: puppeteer.Page, baseURL: string) {
    this.page = page
    this.baseURL = baseURL
    this.pages = new AllPages(this)
  }

  public async navigate(url: string): Promise<void> {
    await this.page.goto(url)
  }

  public async reload<T extends Page>(page: NewablePage<T>): Promise<void> {
    await this.page.goto(this.page.url())
    await this.wait(pageHasLoaded(page))
  }

  public async clearCookies(): Promise<void> {
    const cookies = (await this.page.cookies()).map((x) => ({ name: x.name }))
    await this.page.deleteCookie(...cookies)
  }

  public async wait(condition: WaitCondition, timeout = 30000): Promise<void> {
    const inst = condition(this)
    await this.page.waitForFunction(inst.function, { timeout }, ...inst.args)
  }

  public async press(key: puppeteer.KeyInput): Promise<void> {
    await this.page.keyboard.press(key)
  }

  public async click(selector: string, options?: puppeteer.ClickOptions): Promise<void> {
    await this.page.click(selector, options)
  }

  public async type(selector: string, text: string): Promise<void> {
    await this.page.type(selector, text)
  }

  public async select(selector: string, value: string): Promise<void> {
    await this.page.select(selector, value)
  }

  public async evaluate<T>(fn: puppeteer.EvaluateFn, args: any[]): Promise<T> {
    return await this.page.evaluate(fn, ...args)
  }

  public async waitAny(conditions: WaitCondition | WaitCondition[]): Promise<void> {
    const all = !(conditions instanceof Array) ? [conditions] : conditions

    let retry = 20
    do {
      for (const condition of all) {
        retry--
        try {
          await this.wait(condition, 200)
          return
        } catch (ex) {
          continue
        }
      }
    } while (retry >= 0)
  }
}
