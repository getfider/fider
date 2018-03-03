import './TagForm.scss';

import * as React from 'react';
import { Button, ButtonClickEvent, DisplayError } from '@fider/components/common';
import { ShowTag } from '@fider/components/ShowTag';
import { Tag } from '@fider/models';
import { Failure } from '@fider/services';

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
  constructor(props: TagFormProps) {
    super(props);
    this.state = {
      color: props.color || this.randomizeColor(),
      name: props.name || '',
      isPublic: props.isPublic || false
    };
  }

  private randomizeColor(): string {
    const letters = '0123456789ABCDEF';
    let color = '';
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
    return (
      <div id="tag-form" className="ui form">
        <div className="four fields">
          <div className="four wide field">
            <label>Name</label>
            <input
              className="name"
              onChange={(e) => this.setState({ name: e.currentTarget.value })}
              type="text"
              placeholder="New tag name"
              value={this.state.name}
            />
            <DisplayError fields={['name']} error={this.state.error} pointing="above" />
          </div>
          <div className="three wide field">
            <label>
              Color
              <span
                className="info clickable"
                onClick={() => this.setState({ color: this.randomizeColor() })}
              >
                randomize
              </span>
            </label>
            <input
              className="color"
              onChange={(e) => this.setState({ color: e.currentTarget.value })}
              type="text"
              value={this.state.color}
            />
            <DisplayError fields={['color']} error={this.state.error} pointing="above" />
          </div>
          <div className="two wide field">
            <div className="grouped fields">
              <label>Visibility</label>
              <div className="field">
                <div className="ui radio checkbox">
                  <input
                    id="visibility-public"
                    type="radio"
                    name="visibility"
                    checked={this.state.isPublic}
                    onChange={(e) => this.setState({ isPublic: true })}
                  />
                  <label htmlFor="visibility-public">Public</label>
                </div>
              </div>
              <div className="field">
                <div className="ui radio checkbox">
                  <input
                    id="visibility-private"
                    type="radio"
                    name="visibility"
                    checked={!this.state.isPublic}
                    onChange={(e) => this.setState({ isPublic: false })}
                  />
                  <label htmlFor="visibility-private">Private</label>
                </div>
              </div>
            </div>
          </div>
          <div className="field">
            <label>Preview</label>
            <ShowTag
              tag={{
                id: 0,
                slug: '',
                name: this.state.name,
                color: this.state.color,
                isPublic: this.state.isPublic
              }}
            />
          </div>
        </div>
        <Button
          onClick={async () => this.props.onCancel()}
        >
          Cancel
        </Button>
        <Button
          className="positive"
          onClick={(e) => this.onSave(e)}
        >
          Save
        </Button>
      </div>
    );
  }
}
