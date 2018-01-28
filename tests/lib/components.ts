import { WebElementPromise, Key, By } from 'selenium-webdriver';

export class WebComponent {
  constructor(protected element: WebElementPromise, public selector: string) { }

  public async click() {
    try {
      return await this.element.click();
    } catch (clickErr) {
      try {
        await this.element.getDriver().executeScript('arguments[0].click();', this.element);
      } catch (jsErr) {
        throw clickErr;
      }
    }
  }

  public async getAttribute(attrName: string): Promise<string> {
    return await this.element.getAttribute(attrName);
  }

  public async isDisplayed() {
    try {
      return await this.element.isDisplayed();
    } catch (ex) {
      return false;
    }
  }

  public async getText() {
    if (await this.element.getTagName() === 'input') {
      return await this.element.getAttribute('value');
    }
    return await this.element.getText();
  }
}

export class Button extends WebComponent {
  constructor(element: WebElementPromise, selector: string) {
    super(element, selector);
  }

  public async isDisabled() {
    try {
      return await this.element.getAttribute('disabled') === 'disabled';
    } catch (ex) {
      return false;
    }
  }
}

export class DropDownList extends WebComponent {
  constructor(element: WebElementPromise, selector: string) {
    super(element, selector);
  }

  public async selectByText(text: string) {
    const options = await this.element.findElements(By.tagName('option'));
    for (const option of options) {
      if (await option.getText() === text) {
        return await option.click();
      }
    }
    throw new Error(`No option found for text '${text}'.`);
  }
}

export class TextInput extends WebComponent {
  constructor(element: WebElementPromise, selector: string) {
    super(element, selector);
  }

  public async type(text: string) {
    await this.element.sendKeys(text);
  }

  public async clear() {
    const text = await this.getText();
    for (const char of text) {
      await this.element.sendKeys(Key.ARROW_RIGHT);
      await this.element.sendKeys(Key.BACK_SPACE);
    }
  }
}
