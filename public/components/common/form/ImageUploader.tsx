import "./ImageUploader.scss";

import React from "react";
import { ValidationContext } from "./Form";
import { DisplayError, hasError } from "./DisplayError";
import { classSet, fileToBase64 } from "@fider/services";
import { Button, ButtonClickEvent } from "@fider/components";

interface ImageUploaderProps {
  field: string;
  label?: string;
  defaultImageURL?: string;
  disabled?: boolean;
  previewMaxWidth: number;
  onChange(state: ImageUploaderState, previewURL?: string): void;
}

interface ImageUploaderState extends ImageUploadState {
  previewURL?: string;
}

export interface ImageUploadState {
  upload?: {
    content?: string;
    contentType?: string;
  };
  remove: boolean;
}

export class ImageUploader extends React.Component<ImageUploaderProps, ImageUploaderState> {
  private fileSelector?: HTMLInputElement | null;
  constructor(props: ImageUploaderProps) {
    super(props);
    this.state = {
      upload: undefined,
      remove: false,
      previewURL: this.props.defaultImageURL
    };
  }

  public fileChanged = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      const base64 = await fileToBase64(file);
      this.setState(
        {
          upload: {
            content: base64,
            contentType: file.type,
            action: "upload"
          },
          remove: false,
          previewURL: `data:${file.type};base64,${base64}`
        },
        () => {
          this.props.onChange(this.state, this.state.previewURL);
        }
      );
    }
  };

  public removeFile = async (e: ButtonClickEvent) => {
    this.setState(
      {
        remove: true,
        upload: undefined,
        previewURL: undefined
      },
      () => {
        this.props.onChange(this.state, this.state.previewURL);
      }
    );
  };

  public selectFile = async (e: ButtonClickEvent) => {
    if (this.fileSelector) {
      this.fileSelector.click();
    }
  };

  public render() {
    const isUploading = !!this.state.upload;
    const hasFile = (!this.state.remove && this.props.defaultImageURL) || isUploading;

    const imgStyles: React.CSSProperties = {
      maxWidth: `${this.props.previewMaxWidth}px`
    };

    return (
      <ValidationContext.Consumer>
        {ctx => (
          <div
            className={classSet({
              "c-form-field": true,
              "c-image-upload": true,
              "m-error": hasError(this.props.field, ctx.error)
            })}
          >
            <label htmlFor={`input-${this.props.field}`}>{this.props.label}</label>
            {hasFile && <img className="preview" src={this.state.previewURL} style={imgStyles} />}
            <input ref={e => (this.fileSelector = e)} type="file" onChange={this.fileChanged} />
            <DisplayError fields={[this.props.field]} error={ctx.error} />
            <div className="c-form-field-wrapper">
              <Button size="tiny" onClick={this.selectFile} disabled={this.props.disabled}>
                {hasFile ? "Change" : "Upload"}
              </Button>
              {hasFile && (
                <Button onClick={this.removeFile} size="tiny" disabled={this.props.disabled}>
                  Remove
                </Button>
              )}
            </div>
            {this.props.children}
          </div>
        )}
      </ValidationContext.Consumer>
    );
  }
}
