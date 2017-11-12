import { WebComponent, Browser, Page, findBy, findMultipleBy, Button, TextInput, elementIsVisible, elementIsNotVisible, pageHasLoaded } from '../lib';
import { ShowIdeaPage, GoogleSignInPage, FacebookSignInPage } from './';
import config from '../config';
import { tenant } from '../context';

import { IdeaList } from '../components/IdeaList';

export class HomePage extends Page {
  constructor(browser: Browser) {
    super(browser);
  }

  public getUrl(): string {
    return `http://${tenant}.dev.fider.io:3000/`;
  }

  @findBy('#new-idea-input')
  public IdeaTitle: TextInput;

  @findBy('.ui.form textarea')
  public IdeaDescription: TextInput;

  @findBy('.ui.button.primary')
  public SubmitIdea: Button;

  @findBy('.signin')
  public UserMenu: WebComponent;

  @findBy('.signin .name')
  public UserName: WebComponent;

  @findBy('#signin-modal')
  public SignInModal: WebComponent;

  @findBy('#signin-modal .button.google')
  public GoogleSignIn: Button;

  @findBy('#signin-modal .button.facebook')
  public FacebookSignIn: Button;

  @findBy('#email-signin input')
  private EmailSignInInput: TextInput;

  @findBy('#email-signin button')
  private EmailSignInButton: TextInput;

  @findBy('.ui.form .ui.negative.message')
  public ErrorBox: WebComponent;

  @findBy('.signout')
  private SignOut: Button;

  @findMultipleBy('.fdr-idea-list > .item')
  public IdeaList: IdeaList;

  @findBy('#signin-complete-modal input')
  private CompleteEmailSignInInput: TextInput;

  @findBy('#signin-complete-modal button')
  private CompleteEmailSignInButton: Button;

  public loadCondition() {
    return elementIsVisible(() => this.IdeaTitle);
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

  public async signInWithGoogle(): Promise<void> {
    await this.browser.clearCookies('https://accounts.google.com');
    await this.signOut();

    await this.signIn(() => this.GoogleSignIn);
    await this.browser.waitAny([
      pageHasLoaded(GoogleSignInPage),
      pageHasLoaded(HomePage)
    ]);
  }

  public async signInWithFacebook(): Promise<void> {
    await this.browser.clearCookies('https://facebook.com');
    await this.signOut();

    await this.signIn(() => this.FacebookSignIn);
    await this.browser.waitAny([
      pageHasLoaded(FacebookSignInPage),
      pageHasLoaded(HomePage)
    ]);
  }

  public async signInWithEmail(email: string): Promise<void> {
    await this.signOut();
    await this.UserMenu.click();
    await this.browser.wait(elementIsVisible(() => this.EmailSignInInput));
    await this.EmailSignInInput.type(email);
    await this.EmailSignInButton.click();
    await this.browser.wait(
      elementIsNotVisible(() => this.EmailSignInButton)
    );
  }

  public async completeSignIn(name: string): Promise<void> {
    await this.browser.wait(elementIsVisible(() => this.CompleteEmailSignInInput));
    await this.CompleteEmailSignInInput.type(name);
    await this.CompleteEmailSignInButton.click();
    await this.browser.wait(this.loadCondition());
  }

  private async signIn(locator: () => WebComponent): Promise<void> {
    await this.UserMenu.click();
    await this.browser.wait(elementIsVisible(locator));
    await locator().click();
  }
}
