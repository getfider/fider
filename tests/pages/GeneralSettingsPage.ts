import { BrowserTab, WaitCondition, Page, elementIsVisible, findBy, TextInput, Button, pageHasLoaded } from "../lib"
import { HomePage } from "."

export class GeneralSettingsPage extends Page {
  constructor(tab: BrowserTab) {
    super(tab)
  }

  public getURL(): string {
    return `${this.tab.baseURL}/admin`
  }

  @findBy("#p-admin-general #input-title")
  private TitleInput!: TextInput
  @findBy("#p-admin-general #input-welcomeMessage")
  private WelcomeMessageInput!: TextInput
  @findBy("#p-admin-general #input-invitation")
  private InvitationInput!: TextInput
  @findBy("#p-admin-general .c-button.m-positive")
  private SaveChangesButton!: Button

  public loadCondition(): WaitCondition {
    return elementIsVisible(this.TitleInput)
  }

  public async changeSettings(title: string, welcomeMessage: string, invitation: string): Promise<void> {
    await this.TitleInput.clear()
    await this.TitleInput.type(title)
    await this.WelcomeMessageInput.clear()
    await this.WelcomeMessageInput.type(welcomeMessage)
    await this.InvitationInput.clear()
    await this.InvitationInput.type(invitation)
    await this.SaveChangesButton.click()
    await this.tab.wait(pageHasLoaded(HomePage))
  }
}
