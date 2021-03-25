import "./TagForm.scss"

import React from "react"
import { Button, Input, ShowTag, Form, RadioButton, Field, SelectOption } from "@fider/components"
import { Failure } from "@fider/services"

interface TagFormProps {
  name?: string
  color?: string
  isPublic?: boolean
  onSave: (data: TagFormState) => Promise<Failure | undefined>
  onCancel: () => void
}

export interface TagFormState {
  name: string
  color: string
  isPublic: boolean
  error?: Failure
}

export class TagForm extends React.Component<TagFormProps, TagFormState> {
  private visibilityPublic = { label: "Public", value: "public" }
  private visibilityPrivate = { label: "Private", value: "private" }

  constructor(props: TagFormProps) {
    super(props)
    this.state = {
      color: props.color || this.getRandomColor(),
      name: props.name || "",
      isPublic: props.isPublic || false,
    }
  }

  private getRandomColor(): string {
    const letters = "0123456789ABCDEF"
    let color = ""
    for (let i = 0; i < 6; i++) {
      color += letters[Math.floor(Math.random() * 16)]
    }
    return color
  }

  private handleSave = async () => {
    const error = await this.props.onSave(this.state)
    if (error) {
      this.setState({ error })
    }
  }

  private handleCancel = async () => {
    this.props.onCancel()
  }

  private setName = (name: string) => {
    this.setState({ name })
  }

  private setColor = (color: string) => {
    this.setState({ color })
  }

  private setVisibility = (option: SelectOption) => {
    this.setState({ isPublic: option === this.visibilityPublic })
  }

  private randomize = () => {
    this.setColor(this.getRandomColor())
  }

  public render() {
    const randomizer = (
      <span className="info clickable" onClick={this.randomize}>
        randomize
      </span>
    )

    return (
      <Form error={this.state.error}>
        <div className="row">
          <div className="col-lg-3">
            <Input field="name" label="Name" value={this.state.name} onChange={this.setName} />
          </div>
          <div className="col-lg-2">
            <Input field="color" label="Color" afterLabel={randomizer} value={this.state.color} onChange={this.setColor} />
          </div>
          <div className="col-lg-2">
            <RadioButton
              label="Visibility"
              field="visibility"
              defaultOption={this.state.isPublic ? this.visibilityPublic : this.visibilityPrivate}
              options={[this.visibilityPublic, this.visibilityPrivate]}
              onSelect={this.setVisibility}
            />
          </div>
          <div className="col-lg-5">
            <Field label="Preview">
              <ShowTag
                tag={{
                  id: 0,
                  slug: "",
                  name: this.state.name,
                  color: this.state.color,
                  isPublic: this.state.isPublic,
                }}
              />
            </Field>
          </div>
        </div>
        <Button color="positive" onClick={this.handleSave}>
          Save
        </Button>
        <Button onClick={this.handleCancel} color="cancel">
          Cancel
        </Button>
      </Form>
    )
  }
}
