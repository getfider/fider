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
    const current = await this.getText();
    if (current !== text) {
      await this.clear();
      await this.type(text);
    }
  }

  public async getText(): Promise<string> {
    return await (this.browser.page.evaluate(
      selector => (document.querySelector(selector) as HTMLInputElement).value,
      this.selector
    ) as Promise<string>);
  }

  public async clear() {
    await this.browser.page.click(this.selector, { clickCount: 3 });
    await this.browser.page.keyboard.press("Backspace");
  }
}
