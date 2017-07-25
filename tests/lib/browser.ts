import 'chromedriver';
import { Builder, ThenableWebDriver, WebElement, By, WebElementPromise } from 'selenium-webdriver';
import { Page, NewablePage, WebComponent, WaitCondition, timeout } from './';

export class Browser {
  private driver: ThenableWebDriver;
  public constructor(private browserName: string) {
    this.driver = new Builder().forBrowser(browserName).build();
  }

  public async navigate(url: string): Promise<void> {
    await this.driver.navigate().to(url);
  }

  public findElement(selector: string): WebElementPromise {
    return this.driver.findElement(By.css(selector));
  }

  public async wait(condition: WaitCondition) {
    await this.waitAny(condition);
  }

  public async waitAny(conditions: WaitCondition | WaitCondition[]): Promise<void> {
    const all = (!(conditions instanceof Array)) ? [ conditions ] : conditions;

    await this.driver.wait(async () => {
        do {
          for (const condition of all) {
            try {
              return await Promise.race([
                condition(this),
                timeout(500)
              ]);
            } catch (ex) {
              continue;
            }
          }
        } while (true);
    });
  }

  public async close(): Promise<void> {
    await this.driver.quit();
  }
}
