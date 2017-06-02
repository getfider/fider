import axios from 'axios';
import * as React from 'react';
import * as storage from '../storage';
import { User, Comment, Idea } from '../models';
import { DisplayError } from './Common';

interface ResponseFormProps {
  idea: Idea;
}

interface ResponseFormState {
  active: boolean;
  error?: Error;
}

export class ResponseForm extends React.Component<ResponseFormProps, ResponseFormState> {
  private user: User;
  private text: HTMLTextAreaElement;
  private status: HTMLSelectElement;

  constructor(props: ResponseFormProps) {
    super(props);

    this.user = storage.getCurrentUser();
    this.state = {
      active: false
    };
  }

  private async submit() {
    try {
      await axios.post(`/api/ideas/${this.props.idea.number}/status`, {
          status: parseInt(this.status.value, 10),
          text: this.text.value
      });
      location.reload();
    } catch (ex) {
      this.setState({
        error: ex.response.data
      });
    }
  }

  private activate() {
    this.setState({ active: true });
  }

  private deactivate() {
    this.setState({ active: false });
  }

  public render() {
    const button = <div className="ui blue labeled submit icon button false"
                         onClick={() => this.activate()}>
                     <i className="icon announcement"></i>Respond
                   </div>;

    const form = <form className="ui reply form">
                  <DisplayError error={this.state.error} />
                  <div className="two fields">
                    <div className="field">
                      <label>Status</label>
                      <select className="ui dropdown"
                        ref={(input) => this.status = input}>
                        <option selected={this.props.idea.status === 0} value="0">New</option>
                        <option selected={this.props.idea.status === 1} value="1">Started</option>
                        <option selected={this.props.idea.status === 2} value="2">Completed</option>
                        <option selected={this.props.idea.status === 3} value="3">Declined</option>
                      </select>
                    </div>
                  </div>
                  <div className="field">
                    <textarea
                      ref={(input) => this.text = input}
                      defaultValue={this.props.idea.response && this.props.idea.response.text}
                      placeholder="What's going on with this idea? Let your users know what are your plans...">
                    </textarea>
                  </div>
                  <div className="ui blue labeled submit icon button false"
                        onClick={() => this.submit()}>
                    <i className="icon checkmark box"></i>Submit
                  </div>
                  <div className="ui submit button false"
                        onClick={() => this.deactivate()}>
                    Cancel
                  </div>
                 </form>;

    return this.state.active ? form : button;
  }
}
