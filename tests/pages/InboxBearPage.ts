import {
  WebComponent,
  Browser,
  Page,
  findBy,
  Button,
  elementIsVisible,
  pageHasLoaded,
  delay
} from '../lib';

export class InboxBearPage extends Page {
  constructor(browser: Browser) {
    super(browser);
    this.setUrl(`https://inboxbear.com/q/ktbknuw/83exd39`);
  }

  @findBy('.inbox-toolbar button.btn-danger')
  private Delete: Button;

  @findBy('.inbox-item.unread')
  private UnreadMessage: WebComponent;

  @findBy('#idIframe')
  private MessageBody: WebComponent;

  @findBy('a')
  private MessageBodyLink: WebComponent;

  public loadCondition() {
    return elementIsVisible(() => this.Delete);
  }

  public async clearInbox() {
    await this.Delete.click();
  }

  public async getLinkFromEmail(idx: number): Promise<string> {
    await this.browser.wait(elementIsVisible(() => this.UnreadMessage));
    await this.UnreadMessage.click();

    await this.browser.wait(elementIsVisible(() => this.MessageBody));
    await this.browser.switchTo(this.MessageBody.selector);

    await this.browser.wait(elementIsVisible(() => this.MessageBodyLink));
    return await this.MessageBodyLink.getText();
  }
}
