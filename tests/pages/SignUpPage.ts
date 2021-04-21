import { BrowserTab, WaitCondition, Page, elementIsVisible, findBy, TextInput, WebComponent, Button, pageHasLoaded } from "../lib"
import { HomePage, FacebookSignInPage } from "."

export class SignUpPage extends Page {
  constructor(tab: BrowserTab) {
    super(tab)
  }

  public getURL(): string {
    return `${this.tab.baseURL}/signup`
  }

  @findBy("#p-signup")
  private Container!: WebComponent
  @findBy("#p-signup .c-button.m-facebook")
  public FacebookSignIn!: Button
  @findBy("#p-signup #input-name")
  public UserName!: TextInput
  @findBy("#p-signup #input-email")
  public UserEmail!: TextInput
  @findBy("#p-signup #input-tenantName")
  public TenantName!: TextInput
  @findBy("#p-signup #input-subdomain")
  public Subdomain!: TextInput
  @findBy("#p-signup #input-legalAgreement")
  private LegalAgreement!: WebComponent
  @findBy("#p-signup .c-button.m-positive")
  public Confirm!: Button
  @findBy("#p-signup .c-message.m-success")
  private SubdomainOk!: WebComponent
  @findBy(".c-modal-window")
  private ConfirmationModal!: WebComponent

  public loadCondition(): WaitCondition {
    return elementIsVisible(this.Container)
  }

  public async signInWithFacebook(): Promise<void> {
    await this.FacebookSignIn.click()
    await this.tab.wait(pageHasLoaded(FacebookSignInPage))
  }

  public async signInWithEmail(name: string, email: string): Promise<void> {
    await this.UserName.type(name)
    await this.UserEmail.type(email)
  }

  public async signUpAs(name: string, subdomain: string): Promise<void> {
    await this.TenantName.type(name)
    if (await this.Subdomain.isVisible()) {
      await this.Subdomain.type(subdomain)
      await this.tab.wait(elementIsVisible(this.SubdomainOk))
    }
    if (await this.LegalAgreement.isVisible()) {
      await this.LegalAgreement.click()
    }
    await this.Confirm.click()
    await this.tab.waitAny([pageHasLoaded(HomePage), elementIsVisible(this.ConfirmationModal)])
  }
}
