import { WebComponent, Browser, Page, Button, TextInput, findBy, elementIsVisible, pageHasLoaded } from "../lib";
import { GoogleSignInPage, FacebookSignInPage, HomePage } from "./";

export class SignUpPage extends Page {
  constructor(browser: Browser) {
    super(browser);
  }

  public getUrl(): string {
    return `http://login.dev.fider.io:3000/signup`;
  }

  @findBy("#p-signup") private Container!: WebComponent;
  @findBy("#p-signup .c-button.google") public GoogleSignIn!: Button;
  @findBy("#p-signup .c-button.facebook") public FacebookSignIn!: Button;
  @findBy("#p-signup .form #name") public UserName!: TextInput;
  @findBy("#p-signup .form #email") public UserEmail!: TextInput;
  @findBy("#p-signup .form #tenantName") public TenantName!: TextInput;
  @findBy("#p-signup .form #subdomain") public Subdomain!: TextInput;
  @findBy("#p-signup .c-button.green") public Confirm!: Button;
  @findBy("#p-signup .green.basic.label") private SubdomainOk!: WebComponent;
  @findBy(".c-modal__window") private ConfirmationModal!: WebComponent;

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
    if (await this.Subdomain.isDisplayed()) {
      await this.Subdomain.type(subdomain);
      await this.browser.wait(elementIsVisible(() => this.SubdomainOk));
    }
    await this.Confirm.click();
    await this.browser.waitAny([pageHasLoaded(HomePage), elementIsVisible(() => this.ConfirmationModal)]);
  }
}
