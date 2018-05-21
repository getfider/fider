import {
  WebComponent,
  Browser,
  Page,
  findBy,
  findMultipleBy,
  Button,
  TextInput,
  elementIsVisible,
  elementIsNotVisible,
  pageHasLoaded
} from "../lib";
import { tenant } from "../context";

export class AdminSettingsPage extends Page {
  constructor(browser: Browser) {
    super(browser);
  }

  public getUrl(): string {
    return `http://${tenant}.dev.fider.io:3000/admin`;
  }

  @findBy("#p-admin-general #input-title") public TitleInput!: TextInput;
  @findBy("#p-admin-general #input-welcomeMessage") public WelcomeMessageInput!: TextInput;
  @findBy("#p-admin-general #input-invitation") public InvitationInput!: TextInput;
  @findBy("#p-admin-general .c-button.m-positive") public ConfirmButton!: Button;

  public loadCondition() {
    return elementIsVisible(() => this.TitleInput);
  }
}
