import { ElementHandle } from "puppeteer";
import { elementIsVisible } from ".";

export class WebComponent {
  constructor(protected getter: Promise<ElementHandle | null>, public selector: string) {}

  protected async getHandle(): Promise<ElementHandle> {
    const handle = await this.getter;
    if (handle) {
      return handle;
    } else {
      throw new Error(`Element not found for selector ${this.selector}.`);
    }
  }

  public async click() {
    (await this.getHandle()).click();
  }

  public async isVisible(): Promise<boolean> {
    const condition = elementIsVisible(this.selector);
    return (await this.getHandle()).executionContext().evaluate(condition.function, condition.args);
  }
}

export class Button extends WebComponent {
  constructor(getter: Promise<ElementHandle | null>, selector: string) {
    super(getter, selector);
  }
}

export class TextInput extends WebComponent {
  constructor(getter: Promise<ElementHandle | null>, selector: string) {
    super(getter, selector);
  }

  public async type(text: string) {
    await (await this.getHandle()).type(text);
  }
}
