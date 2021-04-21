import { findBy, Page, BrowserTab, TextInput, elementIsVisible, Button, pageHasLoaded } from "../lib"
import { HomePage, SignUpPage } from "."

export class FacebookSignInPage extends Page {
  constructor(tab: BrowserTab) {
    super(tab)
  }

  @findBy("#email")
  public Email!: TextInput
  @findBy('input[type="password"]')
  public Password!: TextInput
  @findBy("#loginbutton")
  public Confirm!: Button

  public loadCondition() {
    return elementIsVisible(this.Email)
  }

  public async signInAsJonSnow() {
    await this.signInAs("jon_jdrtrsd_snow@tfbnw.net", "jon_jdrtrsd_snow")
  }

  public async signInAsAryaStark() {
    await this.signInAs("arya_xittsrj_stark@tfbnw.net", "arya_xittsrj_stark")
  }

  public async signInAs(email: string, password: string) {
    await this.Email.type(email)
    await this.Password.type(password)
    await this.Confirm.click()
    await this.tab.waitAny([pageHasLoaded(HomePage), pageHasLoaded(SignUpPage)])
  }
}
