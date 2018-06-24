import { WebComponent, BrowserTab } from "../../lib";

export class IdeaItem {
  public Vote: WebComponent;

  constructor(protected tab: BrowserTab, public selector: string) {
    this.Vote = new WebComponent(tab, `${this.selector} .c-support-counter button`);
  }
}

export class IdeaList {
  constructor(protected tab: BrowserTab, public selector: string) {}

  private async findIdeaIndex(title: string): Promise<number> {
    return this.tab.evaluate(
      (text: string, selector: string) => {
        const elements = document.querySelectorAll(`${selector} .c-list-item-title`);
        for (let i = 0; i <= elements.length; i++) {
          if (elements[i].textContent === text) {
            return i;
          }
        }
        return -1;
      },
      [title, this.selector]
    );
  }

  public async get(title: string): Promise<IdeaItem> {
    const idx = await this.findIdeaIndex(title);
    return new IdeaItem(this.tab, `${this.selector} > .c-list-item:nth-child(${idx + 1})`);
  }
}
