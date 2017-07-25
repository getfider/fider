import { TextInput, Button, Browser, WebComponent, Page, findBy, elementIsVisible, elementIsPresent, pageHasLoaded, delay } from '../lib';
import { HomePage, SignUpPage } from './';
import config from '../config';

export class FacebookSignInPage extends Page {
  constructor(browser: Browser) {
    super(browser);
  }

  @findBy('#email')
  public Email: TextInput;

  @findBy('input[type="password"]')
  public Password: TextInput;

  @findBy('#loginbutton')
  public Confirm: Button;

  public loadCondition() {
    return elementIsPresent(() => this.Email);
  }

  public async signInAsJonSnow() {
    await this.signInAs('jon_qvfhyum_snow@tfbnw.net', 'jon_qvfhyum_snow');
  }

  public async signInAsAryaStark() {
    await this.signInAs('arya_mdmwgdg_stark@tfbnw.net', 'arya_mdmwgdg_stark');
  }

  public async signInAs(email: string, password: string) {
    await this.Email.type(email);
    await this.Password.type(password);
    await this.Confirm.click();
    await this.browser.waitAny([
      pageHasLoaded(HomePage),
      pageHasLoaded(SignUpPage)
    ]);
  }
}
