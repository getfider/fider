import { elementIsVisible, BrowserTab } from ".";

export class WebComponent {
  constructor(protected tab: BrowserTab, public selector: string) {}

  public async click() {
    await this.tab.click(this.selector);
  }

  public async getText(): Promise<string> {
    return await this.tab.evaluate<string>(
      (selector: string) => {
        return (document.querySelector(selector) as HTMLElement).innerText;
      },
      [this.selector]
    );
  }

  public async isVisible(): Promise<boolean> {
    const condition = elementIsVisible(this.selector);
    const instance = condition(this.tab);
    return this.tab.evaluate<boolean>(instance.function, instance.args);
  }
}

export class Button extends WebComponent {
  constructor(protected tab: BrowserTab, selector: string) {
    super(tab, selector);
  }
}

export class TextInput extends WebComponent {
  constructor(protected tab: BrowserTab, selector: string) {
    super(tab, selector);
  }

  public async type(text: string) {
    await this.tab.type(this.selector, text);
    const current = await this.getText();
    if (current !== text) {
      await this.clear();
      await this.type(text);
    }
  }

  public async getText(): Promise<string> {
    return await this.tab.evaluate<string>(
      (selector: string) => (document.querySelector(selector) as HTMLInputElement).value,
      [this.selector]
    );
  }

  public async clear() {
    await this.tab.click(this.selector, { clickCount: 3 });
    await this.tab.press("Backspace");
  }
}
