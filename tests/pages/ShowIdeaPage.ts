import { WebComponent, Browser, Page, findBy, elementIsVisible } from '../lib';

export class ShowIdeaPage extends Page {
  constructor(browser: Browser) {
    super(browser);
  }

  @findBy('.idea-header .header')
  public Title: WebComponent;

  @findBy('div.description')
  public Description: WebComponent;

  @findBy('.support-counter .button')
  public SupportCounter: WebComponent;

  public loadCondition() {
    return elementIsVisible(() => this.Title);
  }
}
