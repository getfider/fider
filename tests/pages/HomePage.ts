import { Browser, Page, findBy, TextInput, elementIsVisible, WebComponent, pageHasLoaded, Button } from "../lib";
import { getTenant } from "../context";
import { ShowIdeaPage, FacebookSignInPage } from ".";

export class HomePage extends Page {
  constructor(browser: Browser) {
    super(browser);
  }

  public getUrl(): string {
    return `http://${getTenant()}.dev.fider.io:3000/`;
  }

  // @findBy(".c-menu-item-title") public MenuTitle!: WebComponent;
  // @findBy(".welcome-message") public WelcomeMessage!: WebComponent;
  @findBy("#input-title") public IdeaTitle!: TextInput;
  @findBy("#input-description") public IdeaDescription!: TextInput;
  @findBy(".c-button.m-positive") public SubmitIdea!: Button;
  @findBy(".c-menu-item-signin") public UserMenu!: WebComponent;
  @findBy(".c-menu-user-heading") public UserName!: WebComponent;
  @findBy(".c-modal-window") public SignInModal!: WebComponent;
  @findBy(".c-modal-window .c-button.m-facebook") public FacebookSignIn!: Button;
  // @findBy(".c-modal-window .c-signin-control #input-email") private EmailSignInInput!: TextInput;
  // @findBy(".c-modal-window .c-signin-control .c-button.m-positive") private EmailSignInButton!: TextInput;
  @findBy(".signout") private SignOut!: Button;
  // @findMultipleBy(".c-idea-list > .c-list-item") public IdeaList!: IdeaList;
  // @findBy(".c-modal-window input") private CompleteEmailSignInInput!: TextInput;
  // @findBy(".c-modal-window button") private CompleteEmailSignInButton!: Button;

  public loadCondition() {
    return elementIsVisible(this.IdeaTitle);
  }

  public async submitNewIdea(title: string, description: string): Promise<void> {
    await this.IdeaTitle.type(title);
    await this.IdeaDescription.type(description);
    await this.SubmitIdea.click();
    await this.browser.wait(pageHasLoaded(ShowIdeaPage));
  }

  public async signOut(): Promise<void> {
    try {
      await this.SignOut.click();
    } catch (ex) {
      return;
    }
    await this.browser.wait(pageHasLoaded(HomePage));
  }

  public async signInWithFacebook(): Promise<void> {
    await this.browser.clearCookies("https://facebook.com");
    await this.signOut();

    await this.signIn(this.FacebookSignIn);
    await this.browser.wait(pageHasLoaded(FacebookSignInPage));
  }

  // public async signInWithEmail(email: string): Promise<void> {
  //   await this.signOut();
  //   await this.UserMenu.click();
  //   await this.browser.wait(elementIsVisible(() => this.EmailSignInInput));
  //   await this.EmailSignInInput.type(email);
  //   await this.EmailSignInButton.click();
  //   await this.browser.wait(elementIsNotVisible(() => this.EmailSignInButton));
  // }

  // public async completeSignIn(name: string): Promise<void> {
  //   await this.browser.wait(elementIsVisible(() => this.CompleteEmailSignInInput));
  //   await this.CompleteEmailSignInInput.type(name);
  //   await this.CompleteEmailSignInButton.click();
  //   await this.browser.wait(elementIsNotVisible(() => this.CompleteEmailSignInInput));
  //   await this.browser.wait(this.loadCondition());
  // }

  private async signIn(locator: WebComponent): Promise<void> {
    await this.UserMenu.click();
    await this.browser.wait(elementIsVisible(locator));
    await locator.click();
  }
}
