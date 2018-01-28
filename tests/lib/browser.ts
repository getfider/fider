import 'chromedriver';
import { Builder, ThenableWebDriver, WebElement, Capabilities, By, WebElementPromise } from 'selenium-webdriver';
import { Page, NewablePage, WebComponent, WaitCondition, timeout } from './';

export class Browser {
  private driver: ThenableWebDriver;
  public constructor(private browserName: string) {
    this.driver = new Builder().forBrowser('chrome').build();
  }

  public async navigate(url: string): Promise<void> {
    await this.driver.navigate().to(url);
  }

  public async clearCookies(url?: string): Promise<void> {
    if (url) {
      const currentUrl = await this.driver.getCurrentUrl();
      await this.navigate(url);
      await this.driver.manage().deleteAllCookies();
      await this.navigate(currentUrl);
    }
    await this.driver.manage().deleteAllCookies();
  }

  public findElement(selector: string): WebElementPromise {
    return this.driver.findElement(By.css(selector));
  }

  public async findElements(selector: string): Promise<WebElement[]> {
    return await this.driver.findElements(By.css(selector));
  }

  public async switchTo(element: string | WebElement): Promise<void> {
    if (typeof element === 'string') {
      await this.driver.switchTo().frame(this.findElement(element));
    } else {
      return await this.driver.switchTo().frame(element);
    }
  }

  public async wait(condition: WaitCondition) {
    await this.waitAny(condition);
  }

  public async waitAny(conditions: WaitCondition | WaitCondition[]): Promise<void> {
    const all = (!(conditions instanceof Array)) ? [ conditions ] : conditions;

    await this.driver.wait(async () => {
      for (const condition of all) {
        try {
          if (await condition(this) === true) {
            return true;
          }
          continue;
        } catch (ex) {
          continue;
        }
      }
    });
  }

  public async close(): Promise<void> {
    await this.driver.quit();
  }
}
