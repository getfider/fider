import { GoogleSignInPage } from './GoogleSignInPage';
import { FacebookSignInPage } from './FacebookSignInPage';
import { HomePage } from './HomePage';
import { ShowIdeaPage } from './ShowIdeaPage';
import { SignUpPage } from './SignUpPage';
import { Browser } from '../lib';

export {
  GoogleSignInPage,
  HomePage,
  SignUpPage,
  ShowIdeaPage,
  FacebookSignInPage,
};

export class AllPages {
    public google: GoogleSignInPage;
    public facebook: FacebookSignInPage;
    public home: HomePage;
    public signup: SignUpPage;
    public showIdea: ShowIdeaPage;

    constructor(public browser: Browser) {
      this.google = new GoogleSignInPage(browser);
      this.facebook = new FacebookSignInPage(browser);
      this.home = new HomePage(browser);
      this.signup = new SignUpPage(browser);
      this.showIdea = new ShowIdeaPage(browser);
    }

    public async dispose(): Promise<void> {
      await this.browser.close();
    }
}
