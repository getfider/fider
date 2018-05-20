import "./TagForm.scss";

import * as React from "react";
import { Button, Input, ButtonClickEvent, ShowTag, DisplayError, Form2, RadioButton, Label } from "@fider/components";
import { Tag } from "@fider/models";
import { Failure } from "@fider/services";

interface TagFormProps {
  name?: string;
  color?: string;
  isPublic?: boolean;
  onSave: (data: TagFormState) => Promise<Failure | undefined>;
  onCancel: () => void;
}

export interface TagFormState {
  name: string;
  color: string;
  isPublic: boolean;
  error?: Failure;
}

export class TagForm extends React.Component<TagFormProps, TagFormState> {
  private visibilityPublic = { label: "Public", value: "public" };
  private visibilityPrivate = { label: "Private", value: "private" };

  constructor(props: TagFormProps) {
    super(props);
    this.state = {
      color: props.color || this.randomizeColor(),
      name: props.name || "",
      isPublic: props.isPublic || false
    };
  }

  private randomizeColor(): string {
    const letters = "0123456789ABCDEF";
    let color = "";
    for (let i = 0; i < 6; i++) {
      color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
  }

  private async onSave(e: ButtonClickEvent) {
    const error = await this.props.onSave(this.state);
    if (error) {
      this.setState({ error });
    }
  }

  public render() {
    const randomizer = (
      <span className="info clickable" onClick={() => this.setState({ color: this.randomizeColor() })}>
        randomize
      </span>
    );

    return (
      <>
        <Form2 error={this.state.error}>
          <div className="row">
            <div className="col-lg-3">
              <Input field="name" label="Name" value={this.state.name} onChange={name => this.setState({ name })} />
            </div>
            <div className="col-lg-2">
              <Input
                field="color"
                label="Color"
                afterLabel={randomizer}
                value={this.state.color}
                onChange={color => this.setState({ color })}
              />
            </div>
            <div className="col-lg-2">
              <RadioButton
                label="Visibility"
                field="visibility"
                defaultOption={this.state.isPublic ? this.visibilityPublic : this.visibilityPrivate}
                options={[this.visibilityPublic, this.visibilityPrivate]}
                onSelect={opt => this.setState({ isPublic: opt === this.visibilityPublic })}
              />
            </div>
            <div className="col-lg-5">
              <Label>Preview</Label>
              <ShowTag
                tag={{
                  id: 0,
                  slug: "",
                  name: this.state.name,
                  color: this.state.color,
                  isPublic: this.state.isPublic
                }}
              />
            </div>
          </div>
          <Button color="positive" onClick={e => this.onSave(e)}>
            Save
          </Button>
          <Button onClick={async () => this.props.onCancel()}>Cancel</Button>
        </Form2>
      </>
    );
  }
}
