import { WebComponent, Button, TextInput, delay, List } from "."

const retry = async (fn: () => Promise<void>) => {
  let count = 0
  let lastErr: Error
  do {
    try {
      return await fn()
    } catch (err) {
      lastErr = err
      await delay(500)
      count++
    }
  } while (count !== 3)
  throw lastErr
}

class WebComponentEnsurer {
  constructor(protected component: WebComponent) {}

  public async textIs(expected: string) {
    await retry(async () => {
      const text = await this.component.getText()

      if (expected.trim() !== text.trim()) {
        throw new Error(`Element ${this.component.selector} text is '${text}'. Expected value is '${expected}'`)
      }
    })
  }

  public async attributeIs(attrName: string, expected: string) {
    await retry(async () => {
      const attrValue = await this.component.getAttribute(attrName)

      if (expected.trim() !== attrValue.trim()) {
        throw new Error(`Element ${this.component.selector} ${attrName} is '${attrValue}'. Expected value is '${expected}'`)
      }
    })
  }

  // public async attributeIs(attrName: string, expected: string) {
  //   await this.retry(async () => {
  //     const value = await this.component.getAttribute(attrName);

  //     if (value.trim() !== expected.trim()) {
  //       throw new Error(
  //         `Element ${this.component.selector} attribute '${attrName}' is '${value}'. Expected value is '${expected}'`
  //       );
  //     }
  //   });
  // }

  // public async isVisible() {
  //   await this.retry(async () => {
  //     if (!(await this.component.isDisplayed())) {
  //       throw new Error(`Element ${this.component.selector} is not visible`);
  //     }
  //   });
  // }

  // public async isNotVisible() {
  //   await this.retry(async () => {
  //     if (await this.component.isDisplayed()) {
  //       throw new Error(`Element ${this.component.selector} is visible`);
  //     }
  //   });
  // }
}

class ButtonEnsurer extends WebComponentEnsurer {
  constructor(protected button: Button) {
    super(button)
  }

  // public async isNotDisabled() {
  //   await this.retry(async () => {
  //     if (await this.button.isDisabled()) {
  //       throw new Error(`Button ${this.button.selector} is disabled`);
  //     }
  //   });
  // }
}

class ListEnsurer {
  constructor(protected list: List) {}

  public async countIs(expected: number) {
    await retry(async () => {
      const count = await this.list.count()
      if (count !== expected) {
        throw new Error(`Count of ${this.list.selector} is ${count}, expected ${expected}`)
      }
    })
  }
}

class TextInputEnsurer extends WebComponentEnsurer {
  constructor(component: TextInput) {
    super(component)
  }
}

export function ensure(component: Button): ButtonEnsurer
export function ensure(component: TextInput): TextInputEnsurer
export function ensure(component: WebComponent): WebComponentEnsurer
export function ensure(component: List): ListEnsurer
export function ensure(component: WebComponent): WebComponentEnsurer
export function ensure(component: WebComponent | List | Button): any {
  if (component instanceof Button) {
    return new ButtonEnsurer(component)
  } else if (component instanceof List) {
    return new ListEnsurer(component)
  } else if (component instanceof WebComponent) {
    return new WebComponentEnsurer(component)
  }
}
