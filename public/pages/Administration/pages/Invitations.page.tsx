import * as React from "react";

import { SystemSettings, CurrentUser, Tenant } from "@fider/models";
import { Button, ButtonClickEvent, Textarea, DisplayError } from "@fider/components/common";
import { actions, page, Failure } from "@fider/services";
import { AdminBasePage } from "../components";

interface InvitationsPageProps {
  user: CurrentUser;
  tenant: Tenant;
}

interface InvitationsPageState {
  subject: string;
  message: string;
  recipients: string[];
  numOfRecipients: number;
  rawRecipients: string;
  error?: Failure;
}

export class InvitationsPage extends AdminBasePage<InvitationsPageProps, InvitationsPageState> {
  public name = "invitations";
  public icon = "envelope";
  public title = "Invitations";
  public subtitle = "Invite people to share their feedback";

  constructor(props: InvitationsPageProps) {
    super(props);

    this.state = {
      subject: "",
      message: "",
      recipients: [],
      numOfRecipients: 0,
      rawRecipients: ""
    };

    page.setTitle(`Invitations · Site Settings · ${document.title}`);
  }

  private setRecipients = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const rawRecipients = e.currentTarget.value;
    const recipients = rawRecipients.split(/\n|;|,|\s/gm).filter(x => !!x);
    this.setState({ rawRecipients, recipients, numOfRecipients: recipients.length });
  };

  private sendSample = async (e: ButtonClickEvent) => {
    return;
  };

  private sendInvites = async (e: ButtonClickEvent) => {
    return;
  };

  public content() {
    return (
      <div className="ui form">
        <DisplayError fields={["recipients"]} error={this.state.error} />
        <div className="field">
          <label htmlFor="recipients">Send invitations to</label>
          <Textarea
            id="recipients"
            placeholder="william@example.com; michael@company.com"
            rows={1}
            minRows={1}
            onChange={this.setRecipients}
          />
          <div className="info">
            <p>
              Input the list of all email addresses you wish to invite. Separate each address with either{" "}
              <strong>semicolon</strong>, <strong>comma</strong>, <strong>whitespace</strong> or{" "}
              <strong>line break</strong>.
            </p>
            <p>You can send this invite to a maximum of 30 recipients each time.</p>
          </div>
        </div>
        <DisplayError fields={["subject"]} error={this.state.error} />
        <div className="field">
          <label htmlFor="subject">Subject</label>
          <input
            id="subject"
            defaultValue={`Share your ideas and thoughts about ${this.props.tenant.name}`}
            type="text"
            maxLength={70}
            onChange={e => this.setState({ subject: e.currentTarget.value })}
          />
          <p className="info">
            This is the subject that will be used on the invitation email. Keep it short and sweet.
          </p>
        </div>
        <DisplayError fields={["message"]} error={this.state.error} />
        <div className="field">
          <label htmlFor="message">Message</label>
          <Textarea
            id="message"
            defaultValue={`Hi,

At ${
              this.props.tenant.name
            } we take feedback very seriously, which is why we've launched a space where you can vote, discuss and share your ideas and thoughts about our products and services.

We'd like to extend an invite for you to join this community and raise awareness on topics you care about!

To join, click on the link below.

%invite%

Regards,
${this.props.user.name} (${this.props.tenant.name})`}
            rows={8}
            minRows={8}
            onChange={e => this.setState({ message: e.currentTarget.value })}
          />
          <div className="info">
            <p>
              This is the content of the invite. Be polite and explain what this invite is for, otherwise there's a high
              change people will simply ignore your message.
            </p>
            <p>
              You're allowed to write whatever you want as long as you include the accept link placeholder named{" "}
              <strong>%invite%</strong>.
            </p>
          </div>
        </div>

        <div className="ui tiny header">Sample Invite</div>

        <div className="field">
          <p className="info">
            We highly recommend to send yourself a sample email for you to verify if everything is correct before
            inviting your list of contacts.
          </p>
          {this.props.user.email ? (
            <Button onClick={this.sendSample}>Send a sample email to {this.props.user.email}</Button>
          ) : (
            <Button disabled={true}>Your profile doesn't have an email</Button>
          )}
        </div>

        <div className="ui tiny header">Confirmation</div>

        <div className="field">
          <p className="info">Whenever you're ready, click the following button to send out these invites.</p>
          <Button onClick={this.sendInvites} color="green" disabled={this.state.numOfRecipients === 0}>
            Send {this.state.numOfRecipients} {this.state.numOfRecipients === 1 ? "invite" : "invites"}
          </Button>
        </div>
      </div>
    );
  }
}
