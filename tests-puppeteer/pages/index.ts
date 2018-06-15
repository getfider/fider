import { HomePage } from "./HomePage";
import { SignUpPage } from "./SignUpPage";
import { FacebookSignInPage } from "./FacebookSignInPage";
import { Browser } from "../lib";

export { SignUpPage, HomePage, FacebookSignInPage };

export class AllPages {
  public home: HomePage;
  public signup: SignUpPage;
  public facebook: FacebookSignInPage;

  constructor(public browser: Browser) {
    this.home = new HomePage(browser);
    this.signup = new SignUpPage(browser);
    this.facebook = new FacebookSignInPage(browser);
  }

  public async goTo(url: string): Promise<void> {
    return this.browser.navigate(url);
  }
}
