import {
  TextInput,
  Button,
  Browser,
  WebComponent,
  Page,
  findBy,
  elementIsVisible,
  elementIsPresent,
  pageHasLoaded
} from "../lib";
import { HomePage, SignUpPage } from "./";
import config from "../config";

export class GoogleSignInPage extends Page {
  constructor(browser: Browser) {
    super(browser);
  }

  @findBy("#view_container") private Container!: WebComponent;
  @findBy("#identifierId") public Email!: TextInput;
  @findBy("#identifierNext") public ConfirmEmail!: Button;
  @findBy('input[type="password"]') public Password!: TextInput;
  @findBy("#passwordNext") public ConfirmPassword!: Button;

  public loadCondition() {
    return elementIsPresent(() => this.Container);
  }

  public async signInAsDarthVader() {
    return this.signInAs("darthvader.fider@gmail.com", process.env.DARTHVADER_PASSWORD!);
  }

  public async signInAs(email: string, password: string) {
    try {
      const element = await this.browser.findElement(`p[data-email='${email}']`);
      await element.click();
    } catch (ex) {
      await this.browser.wait(elementIsVisible(() => this.Email));
      await this.Email.type(email);
      await this.ConfirmEmail.click();
    }

    await this.browser.wait(elementIsVisible(() => this.Password));
    await this.Password.type(password);
    await this.ConfirmPassword.click();
    await this.browser.waitAny([pageHasLoaded(HomePage), pageHasLoaded(SignUpPage)]);
  }
}
