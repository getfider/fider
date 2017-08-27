import * as React from 'react';
import { User, Comment, Idea } from '@fider/models';
import { Button, DisplayError } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { Session, IdeaService, Failure } from '@fider/services';

interface ResponseFormProps {
  idea: Idea;
}

interface ResponseFormState {
  error?: Failure;
}

export class ResponseForm extends React.Component<ResponseFormProps, ResponseFormState> {
  private user: User;
  private modal: HTMLDivElement;
  private text: HTMLTextAreaElement;
  private status: HTMLSelectElement;

  @inject(injectables.Session)
  public session: Session;

  @inject(injectables.IdeaService)
  public service: IdeaService;

  constructor(props: ResponseFormProps) {
    super(props);

    this.user = this.session.getCurrentUser();
    this.state = { };
  }

  private async submit() {
    const status = parseInt(this.status.value, 10);
    const result = await this.service.setResponse(this.props.idea.number, status, this.text.value);
    if (result.ok) {
      location.reload();
    } else {
      this.setState({
        error: result.error
      });
    }
  }

  private showModal() {
    $(this.modal).modal('show');
  }

  private closeModel() {
    $(this.modal).modal('hide');
  }

  public render() {
    const button = <Button className="primary" onClick={async () => this.showModal()}>Respond</Button>;

    const modal = <div className="ui form modal" ref={(e) => this.modal = e! }>

                  <div className="content">
                    <DisplayError fields={['status']} error={this.state.error} />
                    <div className="two fields">
                      <div className="field">
                        <label>Status</label>
                        <select className="ui dropdown"
                          ref={(input) => this.status = input!}>
                          <option selected={this.props.idea.status === 0} value="0">Open</option>
                          <option selected={this.props.idea.status === 1} value="1">Started</option>
                          <option selected={this.props.idea.status === 2} value="2">Completed</option>
                          <option selected={this.props.idea.status === 3} value="3">Declined</option>
                        </select>
                      </div>
                    </div>
                    <DisplayError fields={['text']} error={this.state.error} />
                    <div className="field">
                      <textarea
                        ref={(input) => this.text = input!}
                        defaultValue={this.props.idea.response && this.props.idea.response.text}
                        placeholder="What's going on with this idea? Let your users know what are your plans...">
                      </textarea>
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
                 </div>;

    return <div>
            { button }
            { modal }
          </div>;
  }
}
