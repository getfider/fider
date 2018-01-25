import { WebElementPromise } from 'selenium-webdriver';
import { WebComponent, Button, TextInput, delay } from './';

class WebComponentEnsurer {
  public async retry(fn: () => Promise<void>) {
    let count = 0;
    let lastErr: Error;
    do {
      try {
        return await fn();
      } catch (err) {
        lastErr = err;
        await delay(500);
        count++;
      }
    } while (count !== 3);
    throw lastErr;
  }

  constructor(private component: WebComponent) {
  }

  public async textIs(expected: string) {
    await this.retry(async () => {
      const text = await this.component.getText();

      if (expected.trim() !== text.trim()) {
        throw new Error(`Element ${this.component.selector} text is '${text}'. Expected value is '${expected}'`);
      }
    });
  }

  public async isVisible() {
    await this.retry(async () => {
      if (!await this.component.isDisplayed()) {
        throw new Error(`Element ${this.component.selector} is not visible`);
      }
    });
  }

  public async isNotVisible() {
    await this.retry(async () => {
      if (await this.component.isDisplayed()) {
        throw new Error(`Element ${this.component.selector} is visible`);
      }
    });
  }
}

class ButtonEnsurer extends WebComponentEnsurer {
  protected button: Button;
  constructor(button: Button) {
    super(button);
    this.button = button;
  }

  public async isNotDisabled() {
    await this.retry(async () => {
      if (await this.button.isDisabled()) {
        throw new Error(`Button ${this.button.selector} is disabled`);
      }
    });
  }
}

class TextInputEnsurer extends WebComponentEnsurer {
  constructor(element: TextInput) {
    super(element);
  }
}

export function ensure(component: Button): ButtonEnsurer;
export function ensure(component: TextInput): TextInputEnsurer;
export function ensure(component: WebComponent): WebComponentEnsurer;
export function ensure(component: WebComponent | Button): any {
    if (component instanceof Button) {
        return new ButtonEnsurer(component);
    } else if (component instanceof WebComponent) {
        return new WebComponentEnsurer(component);
    }
}
