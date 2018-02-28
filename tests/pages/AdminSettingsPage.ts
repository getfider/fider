import { WebComponent, Browser, Page, findBy, findMultipleBy, Button, TextInput, elementIsVisible, elementIsNotVisible, pageHasLoaded } from '../lib';
import { tenant } from '../context';

export class AdminSettingsPage extends Page {
  constructor(browser: Browser) {
    super(browser);
  }

  public getUrl(): string {
    return `http://${tenant}.dev.fider.io:3000/admin`;
  }

  @findBy('#title')
  public TitleInput!: TextInput;

  @findBy('#welcome-message')
  public WelcomeMessageInput!: TextInput;

  @findBy('#invitation')
  public InvitationInput!: TextInput;

  @findBy('button.positive')
  public ConfirmButton!: Button;

  public loadCondition() {
    return elementIsVisible(() => this.TitleInput);
  }
}
