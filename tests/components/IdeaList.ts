import { WebElement, By } from "selenium-webdriver";
import { Browser, WebComponent, WaitCondition } from "../lib";

export class IdeaList {
  constructor(private elements: Promise<WebElement[]>, private selector: string, private browser: Browser) {}

  public async want(index: number) {
    await (await this.elements)[index].findElement(By.css(".c-support-counter button")).click();
  }

  public async navigateAndWait(index: number, condition: WaitCondition) {
    await (await this.elements)[index].findElement(By.css("a.title")).click();
    await this.browser.wait(condition);
  }

  public async at(index: number): Promise<WebComponent> {
    const selector = ".c-support-counter button";
    const counter = (await this.elements)[index].findElement(By.css(selector));
    return new WebComponent(counter, selector);
  }
}
