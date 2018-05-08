import {
  WebComponent,
  TextInput,
  Button,
  DropDownList,
  Browser,
  Page,
  findBy,
  findMultipleBy,
  elementIsVisible
} from "../lib";
import { CommentList } from "../components/CommentList";

export class ShowIdeaPage extends Page {
  constructor(browser: Browser) {
    super(browser);
  }

  @findBy(".idea-header .header") public Title!: WebComponent;
  @findBy(".description") public Description!: WebComponent;
  @findBy(".c-response .ui.label") public Status!: WebComponent;
  @findBy(".c-response .content") public ResponseText!: WebComponent;
  @findBy(".c-support-counter button") public SupportCounter!: WebComponent;
  @findBy(".comment-input textarea") public CommentInput!: TextInput;
  @findBy(".comment-input button") public SubmitCommentButton!: Button;
  @findMultipleBy(".ui.comments > .comment") public CommentList!: CommentList;
  @findBy(".action-col button.respond") public RespondButton!: Button;
  @findBy(".c-modal__window .c-response-form") public ResponseModal!: WebComponent;
  @findBy(".c-modal__window .c-response-form select") private ResponseModalStatus!: DropDownList;
  @findBy(".c-modal__window .c-response-form textarea") private ResponseModalText!: TextInput;
  @findBy(".c-modal__window .c-modal__footer button.green") private ResponseModalSubmitButton!: Button;

  public loadCondition() {
    return elementIsVisible(() => this.Title);
  }

  public async changeStatus(status: string, text: string): Promise<void> {
    await this.ResponseModalStatus.selectByText(status);
    await this.ResponseModalText.clear();
    await this.ResponseModalText.type(text);
    await this.ResponseModalSubmitButton.click();
  }
}
