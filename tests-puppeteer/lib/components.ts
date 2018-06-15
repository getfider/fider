import { elementIsVisible, Browser } from ".";

export class WebComponent {
  constructor(protected browser: Browser, public selector: string) {}

  public async click() {
    await this.browser.page.click(this.selector);
  }

  public async getText(): Promise<string> {
    return await this.browser.page.evaluate((selector: string) => {
      return (document.querySelector(selector) as HTMLElement).innerText;
    }, this.selector);
  }

  public async isVisible(): Promise<boolean> {
    const condition = elementIsVisible(this.selector);
    const instance = condition(this.browser);
    return this.browser.page.evaluate(instance.function, instance.args);
  }
}

export class Button extends WebComponent {
  constructor(protected browser: Browser, selector: string) {
    super(browser, selector);
  }
}

export class TextInput extends WebComponent {
  constructor(protected browser: Browser, selector: string) {
    super(browser, selector);
  }

  public async type(text: string) {
    await this.browser.page.type(this.selector, text);
  }
}
