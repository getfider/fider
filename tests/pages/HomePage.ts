import { WebComponent, Browser, Page, findBy, Button, TextInput, elementIsVisible, pageHasLoaded } from '../lib';
import { ShowIdeaPage, GoogleSignInPage, FacebookSignInPage } from './';
import config from '../config';
import { tenant } from '../context';

export class HomePage extends Page {
  constructor(browser: Browser) {
    super(browser);
    this.setUrl(`http://${tenant}.dev.fider.io:3000/`);
  }

  @findBy('#new-idea-input')
  public IdeaTitle: TextInput;

  @findBy('.ui.form textarea')
  public IdeaDescription: TextInput;

  @findBy('.ui.button.primary')
  public SubmitIdea: Button;

  @findBy('.signin')
  public UserMenu: WebComponent;

  @findBy('.fdr-profile-popup .button.google')
  public GoogleSignIn: Button;

  @findBy('.fdr-profile-popup .button.facebook')
  public FacebookSignIn: Button;

  @findBy('.ui.form .ui.negative.message')
  public ErrorBox: WebComponent;

  @findBy('.signout')
  private SignOut: Button;

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
      await this.browser.wait(pageHasLoaded(HomePage));
    } catch (ex) {
      return;
    }
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

  public async signIn(locator: () => WebComponent): Promise<void> {
    await this.UserMenu.click();
    await this.browser.wait(elementIsVisible(locator));
    await locator().click();
  }
}
