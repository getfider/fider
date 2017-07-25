import { WebComponent, Browser, Page, Button, TextInput, findBy, elementIsVisible, pageHasLoaded } from '../lib';
import { GoogleSignInPage, HomePage } from './';

export class SignUpPage extends Page {
  constructor(browser: Browser) {
    super(browser);
    this.setUrl('http://login.dev.fider.io:3000/signup');
  }

  @findBy('#fdr-signup-page')
  private Container: WebComponent;

  @findBy('#fdr-signup-page .button.google')
  public GoogleSignIn: Button;

  @findBy('.form #name')
  public Name: TextInput;

  @findBy('.form #subdomain')
  public Subdomain: TextInput;

  @findBy('.button.positive')
  public Confirm: Button;

  @findBy('.green.basic.label')
  private SubdomainOk: WebComponent;

  public loadCondition() {
    return elementIsVisible(() => this.Container);
  }

  public async signInWithGoogle(): Promise<void> {
    await this.GoogleSignIn.click();
    await this.browser.wait(pageHasLoaded(GoogleSignInPage));
  }

  public async signUpAs(name: string, subdomain: string): Promise<void> {
    await this.Name.type(name);
    await this.Subdomain.type(subdomain);
    await this.browser.wait(elementIsVisible(() => this.SubdomainOk));
    await this.Confirm.click();
    await this.browser.wait(pageHasLoaded(HomePage));
  }
}
