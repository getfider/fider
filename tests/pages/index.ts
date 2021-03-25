import { HomePage } from "./HomePage"
import { SignUpPage } from "./SignUpPage"
import { FacebookSignInPage } from "./FacebookSignInPage"
import { GeneralSettingsPage } from "./GeneralSettingsPage"
import { ShowPostPage } from "./ShowPostPage"
import { BrowserTab } from "../lib"

export { SignUpPage, ShowPostPage, HomePage, FacebookSignInPage, GeneralSettingsPage }

export class AllPages {
  public home: HomePage
  public signup: SignUpPage
  public showPost: ShowPostPage
  public facebook: FacebookSignInPage
  public generalSettings: GeneralSettingsPage

  constructor(public tab: BrowserTab) {
    this.home = new HomePage(tab)
    this.signup = new SignUpPage(tab)
    this.showPost = new ShowPostPage(tab)
    this.facebook = new FacebookSignInPage(tab)
    this.generalSettings = new GeneralSettingsPage(tab)
  }

  public async goTo(url: string): Promise<void> {
    return this.tab.navigate(url)
  }
}
