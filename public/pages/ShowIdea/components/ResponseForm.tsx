import * as React from "react";

import { Modal, Button, DisplayError, Textarea, Select, Form, TextArea } from "@fider/components";
import { Comment, Idea, IdeaStatus, User } from "@fider/models";
import { IdeaSearch } from "../";

import { actions, Failure } from "@fider/services";

interface ResponseFormProps {
  idea: Idea;
}

interface ResponseFormState {
  showModal: boolean;
  status: number;
  text: string;
  originalNumber: number;
  error?: Failure;
}

export class ResponseForm extends React.Component<ResponseFormProps, ResponseFormState> {
  private modal!: HTMLDivElement;

  constructor(props: ResponseFormProps) {
    super(props);

    this.state = {
      showModal: false,
      status: this.props.idea.status,
      originalNumber: 0,
      text: this.props.idea.response && this.props.idea.response.text
    };
  }

  private async submit() {
    const result = await actions.respond(this.props.idea.number, this.state);
    if (result.ok) {
      location.reload();
    } else {
      this.setState({
        error: result.error
      });
    }
  }

  public render() {
    const button = (
      <Button className="respond" fluid={true} onClick={async () => this.setState({ showModal: true })}>
        <i className="announcement icon" /> Respond
      </Button>
    );

    const options = IdeaStatus.All.map(s => ({
      value: s.value.toString(),
      label: s.title
    }));

    const modal = (
      <Modal.Window isOpen={this.state.showModal} center={false} size="large">
        <Modal.Content>
          <Form error={this.state.error} className="c-response-form">
            <Select
              field="status"
              label="Status"
              defaultValue={this.props.idea.status.toString()}
              options={options}
              onChange={opt => opt && this.setState({ status: parseInt(opt.value, 10) })}
            />
            {this.state.status === IdeaStatus.Duplicate.value ? (
              <>
                <DisplayError fields={["originalNumber"]} error={this.state.error} />
                <IdeaSearch
                  exclude={[this.props.idea.number]}
                  onChanged={originalNumber => this.setState({ originalNumber })}
                />
                <span className="info">Votes from this idea will be merged into original idea.</span>
              </>
            ) : (
              <TextArea
                field="text"
                onChange={text => this.setState({ text })}
                value={this.state.text}
                minRows={5}
                placeholder="What's going on with this idea? Let your users know what are your plans..."
              />
            )}
          </Form>
        </Modal.Content>

        <Modal.Footer>
          <Button color="positive" onClick={() => this.submit()}>
            Submit
          </Button>
          <Button onClick={async () => this.setState({ showModal: false })}>Cancel</Button>
        </Modal.Footer>
      </Modal.Window>
    );

    return (
      <>
        {button}
        {modal}
      </>
    );
  }
}
