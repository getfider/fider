import { WebComponent, Browser, Page, Button, TextInput, findBy, elementIsVisible, pageHasLoaded, delay } from '../lib';
import { GoogleSignInPage, FacebookSignInPage, HomePage } from './';

export class SignUpPage extends Page {
  constructor(browser: Browser) {
    super(browser);
    this.setUrl('http://login.dev.fider.io:3000/signup');
  }

  @findBy('#fdr-signup-page')
  private Container: WebComponent;

  @findBy('#fdr-signup-page .button.google')
  public GoogleSignIn: Button;

  @findBy('#fdr-signup-page .button.facebook')
  public FacebookSignIn: Button;

  @findBy('#fdr-signup-page .form #name')
  public UserName: TextInput;

  @findBy('#fdr-signup-page .form #email')
  public UserEmail: TextInput;

  @findBy('#fdr-signup-page .form #tenantName')
  public TenantName: TextInput;

  @findBy('#fdr-signup-page .form #subdomain')
  public Subdomain: TextInput;

  @findBy('#fdr-signup-page .button.positive')
  public Confirm: Button;

  @findBy('#fdr-signup-page .page .green.basic.label')
  private SubdomainOk: WebComponent;

  @findBy('#submitted-modal')
  private SubmitConfirmation: WebComponent;

  public loadCondition() {
    return elementIsVisible(() => this.Container);
  }

  public async signInWithGoogle(): Promise<void> {
    await this.GoogleSignIn.click();
    await this.browser.wait(pageHasLoaded(GoogleSignInPage));
  }

  public async signInWithFacebook(): Promise<void> {
    await this.FacebookSignIn.click();
    await this.browser.wait(pageHasLoaded(FacebookSignInPage));
  }

  public async signInWithEmail(name: string, email: string): Promise<void> {
    await this.UserName.type(name);
    await this.UserEmail.type(email);
  }

  public async signUpAs(name: string, subdomain: string): Promise<void> {
    await this.TenantName.type(name);
    await this.Subdomain.type(subdomain);
    await this.browser.wait(elementIsVisible(() => this.SubdomainOk));
    await this.Confirm.click();
    await this.browser.waitAny([
      pageHasLoaded(HomePage),
      elementIsVisible(() => this.SubmitConfirmation)
    ]);
  }
}
