import { BrowserTab, Page, WebComponent, TextInput, Button, findBy, elementIsVisible, DropDownList } from "../lib"
import { CommentList } from "./components"

export class ShowPostPage extends Page {
  constructor(tab: BrowserTab) {
    super(tab)
  }

  @findBy(".post-header h1")
  public Title!: WebComponent
  @findBy(".description")
  public Description!: WebComponent
  @findBy(".c-segment.l-response .status-label")
  public Status!: WebComponent
  @findBy(".c-segment.l-response .content")
  public ResponseText!: WebComponent
  @findBy(".c-vote-counter button")
  public VoteCounter!: WebComponent
  @findBy(".action-col .c-button.respond")
  public RespondButton!: Button
  @findBy(".c-modal-window .c-response-form")
  public ResponseModal!: WebComponent
  @findBy(".c-comment-list")
  public Comments!: CommentList

  @findBy(".c-modal-window .c-response-form #input-status")
  private ResponseModalStatus!: DropDownList
  @findBy(".c-modal-window .c-response-form #input-text")
  private ResponseModalText!: TextInput
  @findBy(".c-modal-window .c-modal-footer .c-button.m-positive")
  private ResponseModalSubmitButton!: Button

  @findBy(".c-comment-input #input-content")
  private CommentInput!: TextInput
  @findBy(".c-comment-input .c-button.m-positive")
  private SubmitCommentButton!: Button

  @findBy(".action-col .c-button.edit")
  private Edit!: Button
  @findBy("#input-title")
  private EditTitle!: TextInput
  @findBy("#input-description")
  private EditDescription!: TextInput
  @findBy(".action-col .c-button.save")
  private SaveEdit!: Button

  public loadCondition() {
    return elementIsVisible(this.Title)
  }

  public async changeStatus(status: string, text: string): Promise<void> {
    await this.RespondButton.click()
    await this.tab.wait(elementIsVisible(this.ResponseModal))
    await this.ResponseModalStatus.selectByText(status)
    await this.ResponseModalText.clear()
    await this.ResponseModalText.type(text)
    await this.ResponseModalSubmitButton.click()
  }

  public async edit(newTitle: string, newDescription: string): Promise<void> {
    await this.Edit.click()
    await this.EditTitle.clear()
    await this.EditTitle.type(newTitle)
    await this.EditDescription.clear()
    await this.EditDescription.type(newDescription)
    await this.SaveEdit.click()
  }

  public async comment(text: string): Promise<void> {
    await this.CommentInput.type(text)
    await this.SubmitCommentButton.click()
  }
}
