import * as React from 'react';

import { User, Comment, Idea } from '@fider/models';
import { Button, DisplayError, Textarea } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { IdeaService, Failure } from '@fider/services';

interface ResponseFormProps {
  idea: Idea;
}

interface ResponseFormState {
  status: number;
  text: string;
  error?: Failure;
}

export class ResponseForm extends React.Component<ResponseFormProps, ResponseFormState> {
  private modal: HTMLDivElement;

  @inject(injectables.IdeaService)
  public service: IdeaService;

  constructor(props: ResponseFormProps) {
    super(props);

    this.state = {
      status: this.props.idea.status,
      text: this.props.idea.response && this.props.idea.response.text
    };
  }

  private async submit() {
    const result = await this.service.setResponse(this.props.idea.number, this.state.status, this.state.text);
    if (result.ok) {
      location.reload();
    } else {
      this.setState({
        error: result.error
      });
    }
  }

  private showModal() {
    $(this.modal).modal({
      blurring: true
    }).modal('show');
  }

  private closeModel() {
    $(this.modal).modal('hide');
  }

  public render() {
    const button = (
      <Button className="icon fluid text-left"  onClick={async () => this.showModal()}>
        <i className="announcement icon" /> Respond
      </Button>
    );

    const modal = (
      <div className="ui form modal" ref={(e) => this.modal = e!}>
        <div className="content">
          <DisplayError fields={['status']} error={this.state.error} />
          <div className="two fields">
            <div className="field">
              <label>Status</label>
              <select
                className="ui dropdown"
                defaultValue={this.props.idea.status.toString()}
                onChange={(e) => this.setState({ status: parseInt(e.currentTarget.value, 10) })}
              >
                <option value="0">Open</option>
                <option value="1">Started</option>
                <option value="2">Completed</option>
                <option value="3">Declined</option>
              </select>
            </div>
          </div>
          <DisplayError fields={['text']} error={this.state.error} />
          <div className="field">
            <Textarea
              onChange={(e) => this.setState({ text: e.currentTarget.value })}
              defaultValue={this.state.text}
              placeholder="What's going on with this idea? Let your users know what are your plans..."
            />
          </div>
        </div>

        <div className="actions">
          <Button className="primary" onClick={() => this.submit()}>
            Submit
          </Button>
          <Button className="basic" onClick={async () => this.closeModel()}>
            Cancel
          </Button>
        </div>
        </div>
      );

    return (
      <div>
        {button}
        {modal}
      </div>
    );
  }
}
