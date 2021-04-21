import { BrowserTab, List } from "../../lib"

export class CommentList extends List {
  constructor(tab: BrowserTab, selector: string) {
    super(tab, selector)
  }

  public async count(): Promise<number> {
    return await this.tab.evaluate<number>(
      (selector: string) => {
        return document.querySelectorAll(`${selector} .c-comment`).length
      },
      [this.selector]
    )
  }
}
