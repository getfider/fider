import { HomePage } from "./HomePage";
import { SignUpPage } from "./SignUpPage";
import { FacebookSignInPage } from "./FacebookSignInPage";
import { ShowIdeaPage } from "./ShowIdeaPage";
import { BrowserTab } from "../lib";

export { SignUpPage, ShowIdeaPage, HomePage, FacebookSignInPage };

export class AllPages {
  public home: HomePage;
  public signup: SignUpPage;
  public showIdea: ShowIdeaPage;
  public facebook: FacebookSignInPage;

  constructor(public tab: BrowserTab) {
    this.home = new HomePage(tab);
    this.signup = new SignUpPage(tab);
    this.showIdea = new ShowIdeaPage(tab);
    this.facebook = new FacebookSignInPage(tab);
  }

  public async goTo(url: string): Promise<void> {
    return this.tab.navigate(url);
  }
}
