import { Browser, WaitCondition, Page, elementIsVisible, findBy, TextInput, delay, WebComponent, Button } from "../lib";

export class SignUpPage extends Page {
  constructor(browser: Browser) {
    super(browser);
  }

  public getUrl(): string {
    return `http://login.dev.fider.io:3000/signup`;
  }

  // @findBy("#p-signup") private Container!: WebComponent;
  // @findBy("#p-signup .c-button.m-google") public GoogleSignIn!: Button;
  // @findBy("#p-signup .c-button.m-facebook") public FacebookSignIn!: Button;
  @findBy("#p-signup #input-name") public UserName!: TextInput;
  @findBy("#p-signup #input-email") public UserEmail!: TextInput;
  @findBy("#p-signup #input-tenantName") public TenantName!: TextInput;
  @findBy("#p-signup #input-subdomain") public Subdomain!: TextInput;
  @findBy("#p-signup .c-button.m-positive") public Confirm!: Button;
  @findBy("#p-signup .c-message.m-success") private SubdomainOk!: WebComponent;
  @findBy(".c-modal-window") private ConfirmationModal!: WebComponent;

  public loadCondition(): WaitCondition {
    return elementIsVisible("#p-signup");
  }

  // public async signInWithGoogle(): Promise<void> {
  //   await this.GoogleSignIn.click();
  //   await this.browser.wait(pageHasLoaded(GoogleSignInPage));
  // }

  // public async signInWithFacebook(): Promise<void> {
  //   await this.FacebookSignIn.click();
  //   await this.browser.wait(pageHasLoaded(FacebookSignInPage));
  // }

  public async signInWithEmail(name: string, email: string): Promise<void> {
    await this.UserName.type(name);
    await this.UserEmail.type(email);
  }

  public async signUpAs(name: string, subdomain: string): Promise<void> {
    await this.TenantName.type(name);
    if (await this.Subdomain.isVisible()) {
      await this.Subdomain.type(subdomain);
      await this.browser.wait(elementIsVisible(this.SubdomainOk));
    }
    await this.Confirm.click();
    await this.browser.waitAny([pageHasLoaded(HomePage), elementIsVisible(this.ConfirmationModal)]);
  }
}
