import * as React from "react";

import { Modal, Button, DisplayError, Textarea, Select, Form, TextArea, Field, SelectOption } from "@fider/components";
import { Comment, Post, IdeaStatus, User } from "@fider/models";
import { IdeaSearch } from "../";

import { actions, Failure } from "@fider/services";

interface ResponseFormProps {
  idea: Post;
}

interface ResponseFormState {
  showModal: boolean;
  status: number;
  text: string;
  originalNumber: number;
  error?: Failure;
}

export class ResponseForm extends React.Component<ResponseFormProps, ResponseFormState> {
  constructor(props: ResponseFormProps) {
    super(props);

    this.state = {
      showModal: false,
      status: this.props.idea.status,
      originalNumber: 0,
      text: this.props.idea.response ? this.props.idea.response.text : ""
    };
  }

  private submit = async () => {
    const result = await actions.respond(this.props.idea.number, this.state);
    if (result.ok) {
      location.reload();
    } else {
      this.setState({
        error: result.error
      });
    }
  };

  private showModal = async () => {
    this.setState({ showModal: true });
  };

  private closeModal = async () => {
    this.setState({ showModal: false });
  };

  private setStatus = (opt?: SelectOption) => {
    if (opt) {
      this.setState({ status: parseInt(opt.value, 10) });
    }
  };

  private setOriginalNumber = (originalNumber: number) => {
    this.setState({ originalNumber });
  };

  private setText = (text: string) => {
    this.setState({ text });
  };

  public render() {
    const button = (
      <Button className="respond" fluid={true} onClick={this.showModal}>
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
              onChange={this.setStatus}
            />
            {this.state.status === IdeaStatus.Duplicate.value ? (
              <>
                <Field>
                  <IdeaSearch exclude={[this.props.idea.number]} onChanged={this.setOriginalNumber} />
                </Field>
                <DisplayError fields={["originalNumber"]} error={this.state.error} />
                <span className="info">Votes from this idea will be merged into original idea.</span>
              </>
            ) : (
              <TextArea
                field="text"
                onChange={this.setText}
                value={this.state.text}
                minRows={5}
                placeholder="What's going on with this idea? Let your users know what are your plans..."
              />
            )}
          </Form>
        </Modal.Content>

        <Modal.Footer>
          <Button color="positive" onClick={this.submit}>
            Submit
          </Button>
          <Button color="cancel" onClick={this.closeModal}>
            Cancel
          </Button>
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
