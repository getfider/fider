import { elementIsVisible, BrowserTab } from "."

export class WebComponent {
  constructor(protected tab: BrowserTab, public selector: string) {}

  public async click() {
    await this.tab.click(this.selector)
  }

  public async getText(): Promise<string> {
    return await this.tab.evaluate<string>(
      (selector: string) => {
        const el = document.querySelector(selector) as HTMLElement | undefined
        return el ? el.textContent : ""
      },
      [this.selector]
    )
  }

  public async getAttribute(attributeName: string): Promise<string> {
    return await this.tab.evaluate<string>(
      (selector: string, attrName: string) => {
        const el = document.querySelector(selector) as HTMLElement | undefined
        return el ? el.getAttribute(attrName) || "" : ""
      },
      [this.selector, attributeName]
    )
  }

  public async isVisible(): Promise<boolean> {
    const condition = elementIsVisible(this.selector)
    const instance = condition(this.tab)
    return this.tab.evaluate<boolean>(instance.function, instance.args)
  }
}

export class Button extends WebComponent {
  constructor(protected tab: BrowserTab, selector: string) {
    super(tab, selector)
  }
}

export class DropDownList extends WebComponent {
  constructor(protected tab: BrowserTab, selector: string) {
    super(tab, selector)
  }

  public async selectByText(text: string): Promise<void> {
    const value = await this.tab.evaluate<string>(
      (selector: string, textToSelect: string) => {
        const options = document.querySelectorAll(`${selector} option`)
        for (let i = 0; i <= options.length; i++) {
          if (options[i] && options[i].textContent === textToSelect) {
            return (options[i] as HTMLOptionElement).value
          }
        }
        return ""
      },
      [this.selector, text]
    )
    await this.tab.select(this.selector, value)
  }
}

export abstract class List {
  constructor(protected tab: BrowserTab, public selector: string) {}
  public abstract count(): Promise<number>
}

export class TextInput extends WebComponent {
  constructor(protected tab: BrowserTab, selector: string) {
    super(tab, selector)
  }

  public async type(text: string) {
    await this.tab.type(this.selector, text)
    const current = await this.getText()
    if (current !== text) {
      await this.clear()
      await this.type(text)
    }
  }

  public async getText(): Promise<string> {
    return await this.tab.evaluate<string>((selector: string) => (document.querySelector(selector) as HTMLInputElement).value, [this.selector])
  }

  public async clear() {
    await this.tab.click(this.selector, { clickCount: 3 })
    await this.tab.press("Backspace")
  }
}
