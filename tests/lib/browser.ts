import "chromedriver";
import * as puppeteer from "puppeteer";
import { WaitCondition } from "./conditions";

export class Browser {
  private browser: puppeteer.Browser;
  public page: puppeteer.Page; // TODO: make it private

  public constructor(browser: puppeteer.Browser, page: puppeteer.Page) {
    this.browser = browser;
    this.page = page;
  }

  public static async launch(): Promise<Browser> {
    const browser = await puppeteer.launch({
      headless: process.env.CI ? true : false
    });
    const page = await browser.newPage();
    return new Browser(browser, page);
  }

  public async navigate(url: string): Promise<void> {
    await this.page.goto(url);
  }

  public async clearCookies(): Promise<void> {
    const cookies = (await this.page.cookies()).map(x => ({ name: x.name }));
    await this.page.deleteCookie(...cookies);
  }

  public async wait(condition: WaitCondition, timeout = 30000): Promise<void> {
    const inst = condition(this);
    await this.page.waitForFunction(inst.function, { timeout }, inst.args);
  }

  public async waitAny(conditions: WaitCondition | WaitCondition[]): Promise<void> {
    const all = !(conditions instanceof Array) ? [conditions] : conditions;

    let retry = 20;
    do {
      for (const condition of all) {
        retry--;
        try {
          await this.wait(condition, 200);
          return;
        } catch (ex) {
          continue;
        }
      }
    } while (retry >= 0);
  }

  public async findElement(selector: string): Promise<puppeteer.ElementHandle | null> {
    return await this.page.$(selector);
  }

  public async close(): Promise<void> {
    await this.browser.close();
  }
}
