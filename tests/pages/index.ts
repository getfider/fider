import { GoogleSignInPage } from "./GoogleSignInPage";
import { FacebookSignInPage } from "./FacebookSignInPage";
import { HomePage } from "./HomePage";
import { AdminSettingsPage } from "./AdminSettingsPage";
import { ShowIdeaPage } from "./ShowIdeaPage";
import { SignUpPage } from "./SignUpPage";
import { Browser } from "../lib";

export { GoogleSignInPage, HomePage, SignUpPage, ShowIdeaPage, FacebookSignInPage, AdminSettingsPage };

export class AllPages {
  public google: GoogleSignInPage;
  public facebook: FacebookSignInPage;
  public home: HomePage;
  public adminSettings: AdminSettingsPage;
  public signup: SignUpPage;
  public showIdea: ShowIdeaPage;

  constructor(public browser: Browser) {
    this.google = new GoogleSignInPage(browser);
    this.facebook = new FacebookSignInPage(browser);
    this.home = new HomePage(browser);
    this.signup = new SignUpPage(browser);
    this.showIdea = new ShowIdeaPage(browser);
    this.adminSettings = new AdminSettingsPage(browser);
  }

  public async goTo(url: string): Promise<void> {
    return this.browser.navigate(url);
  }

  public async dispose(): Promise<void> {
    await this.browser.close();
  }
}
