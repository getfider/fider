import { HomePage } from "./HomePage";
import { SignUpPage } from "./SignUpPage";
import { FacebookSignInPage } from "./FacebookSignInPage";
import { ShowIdeaPage } from "./ShowIdeaPage";
import { Browser } from "../lib";

export { SignUpPage, ShowIdeaPage, HomePage, FacebookSignInPage };

export class AllPages {
  public home: HomePage;
  public signup: SignUpPage;
  public showIdea: ShowIdeaPage;
  public facebook: FacebookSignInPage;

  constructor(public browser: Browser) {
    this.home = new HomePage(browser);
    this.signup = new SignUpPage(browser);
    this.showIdea = new ShowIdeaPage(browser);
    this.facebook = new FacebookSignInPage(browser);
  }

  public async goTo(url: string): Promise<void> {
    return this.browser.navigate(url);
  }
}
