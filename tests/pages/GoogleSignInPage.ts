import { TextInput, Button, Browser, WebComponent, Page, findBy, elementIsVisible, elementIsPresent, pageHasLoaded, delay } from '../lib';
import { HomePage, SignUpPage } from './';
import config from '../config';

export class GoogleSignInPage extends Page {
  constructor(browser: Browser) {
    super(browser);
  }

  @findBy('base[href="https://accounts.google.com/"]')
  private Head: WebComponent;

  @findBy('#identifierId')
  public Email: TextInput;

  @findBy('#identifierNext')
  public ConfirmEmail: Button;

  @findBy('input[type="password"]')
  public Password: TextInput;

  @findBy('#passwordNext')
  public ConfirmPassword: Button;

  public loadCondition() {
    return elementIsPresent(() => this.Head);
  }

  public async signInAsDarthVader() {
    return this.signInAs(config.users.darthvader.email, config.users.darthvader.password);
  }

  public async signInAs(email: string, password: string) {

    try {
      const element = await this.browser.findElement(`p[data-email='${config.users.darthvader.email}']`);
      await element.click();
    } catch (ex) {
      await this.Email.type(email);
      await this.ConfirmEmail.click();
    }

    await this.browser.wait(elementIsVisible(() => this.Password));
    await this.Password.type(password);
    await this.ConfirmPassword.click();
    await this.browser.waitAny([
      pageHasLoaded(HomePage),
      pageHasLoaded(SignUpPage)
    ]);
  }
}
