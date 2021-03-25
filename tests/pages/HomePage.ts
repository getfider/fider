import { BrowserTab, Page, findBy, TextInput, elementIsVisible, WebComponent, pageHasLoaded, Button, elementIsNotVisible } from "../lib"
import { ShowPostPage, FacebookSignInPage } from "."
import { PostList } from "./components"

export class HomePage extends Page {
  constructor(tab: BrowserTab) {
    super(tab)
  }

  public getURL(): string {
    return `${this.tab.baseURL}/`
  }

  @findBy(".c-menu-item-title")
  public MenuTitle!: WebComponent
  @findBy(".welcome-message")
  public WelcomeMessage!: WebComponent
  @findBy("#input-title")
  public PostTitle!: TextInput
  @findBy("#input-description")
  public PostDescription!: TextInput
  @findBy(".c-button.m-positive")
  public SubmitPost!: Button
  @findBy(".c-menu-item-signin")
  public UserMenu!: WebComponent
  @findBy(".c-unread-count")
  public UnreadCounter!: WebComponent
  @findBy(".c-menu-user-heading")
  public UserName!: WebComponent
  @findBy(".c-modal-window")
  public SignInModal!: WebComponent
  @findBy(".c-modal-window .c-button.m-facebook")
  public FacebookSignIn!: Button
  @findBy(".c-modal-window .c-signin-control #input-email")
  private EmailSignInInput!: TextInput
  @findBy(".c-modal-window .c-signin-control .c-button.m-positive")
  private EmailSignInButton!: TextInput
  @findBy(".signout")
  private SignOut!: Button
  @findBy(".c-post-list")
  public PostList!: PostList
  @findBy(".c-modal-window input")
  private CompleteEmailSignInInput!: TextInput
  @findBy(".c-modal-window button")
  private CompleteEmailSignInButton!: Button

  public loadCondition() {
    return elementIsVisible("#p-home")
  }

  public async submitNewPost(title: string, description: string): Promise<void> {
    await this.PostTitle.type(title)
    await this.PostDescription.type(description)
    await this.SubmitPost.click()
    await this.tab.wait(pageHasLoaded(ShowPostPage))
  }

  public async signOut(): Promise<void> {
    if (await this.SignOut.isVisible()) {
      await this.UserMenu.click()
      await this.SignOut.click()
      await this.tab.wait(pageHasLoaded(HomePage))
    }
  }

  public async signInWithFacebook(): Promise<void> {
    await this.signOut()
    await this.signIn(this.FacebookSignIn)
    await this.tab.wait(pageHasLoaded(FacebookSignInPage))
  }

  public async signInWithEmail(email: string): Promise<void> {
    await this.signOut()
    await this.UserMenu.click()
    await this.tab.wait(elementIsVisible(this.EmailSignInInput))
    await this.EmailSignInInput.type(email)
    await this.EmailSignInButton.click()
    await this.tab.wait(elementIsNotVisible(this.EmailSignInButton))
  }

  public async completeSignIn(name: string): Promise<void> {
    await this.tab.wait(elementIsVisible(this.CompleteEmailSignInInput))
    await this.CompleteEmailSignInInput.type(name)
    await this.CompleteEmailSignInButton.click()
    await this.tab.wait(elementIsNotVisible(this.CompleteEmailSignInInput))
    await this.tab.wait(this.loadCondition())
  }

  private async signIn(locator: WebComponent): Promise<void> {
    await this.UserMenu.click()
    await this.tab.wait(elementIsVisible(locator))
    await locator.click()
  }
}
