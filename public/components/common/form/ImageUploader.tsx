import "./ImageUploader.scss"

import React from "react"
import { ValidationContext } from "./Form"
import { DisplayError, hasError } from "./DisplayError"
import { classSet, fileToBase64, uploadedImageURL } from "@fider/services"
import { Button, Icon, Modal } from "@fider/components"
import { ImageUpload } from "@fider/models"
import IconPhotograph from "@fider/assets/images/heroicons-photograph.svg"

const hardFileSizeLimit = 5 * 1024 * 1024

interface ImageUploaderProps {
  instanceID?: string
  field: string
  label?: string
  bkey?: string
  disabled?: boolean
  onChange(state: ImageUpload, instanceID?: string, previewURL?: string): void
}

interface ImageUploaderState extends ImageUpload {
  previewURL?: string
  showModal: boolean
}

export class ImageUploader extends React.Component<ImageUploaderProps, ImageUploaderState> {
  private fileSelector?: HTMLInputElement | null

  constructor(props: ImageUploaderProps) {
    super(props)
    this.state = {
      upload: undefined,
      remove: false,
      showModal: false,
      previewURL: uploadedImageURL(this.props.bkey),
    }
  }

  public fileChanged = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0]
      if (file.size > hardFileSizeLimit) {
        alert("The image size must be smaller than 5MB.")
        return
      }

      const base64 = await fileToBase64(file)
      this.setState(
        {
          bkey: this.props.bkey,
          upload: {
            fileName: file.name,
            content: base64,
            contentType: file.type,
          },
          remove: false,
          previewURL: `data:${file.type};base64,${base64}`,
        },
        () => {
          this.props.onChange(this.state, this.props.instanceID, this.state.previewURL)
        }
      )
    }
  }

  public removeFile = async () => {
    if (this.fileSelector) {
      this.fileSelector.value = ""
    }

    this.setState(
      {
        bkey: this.props.bkey,
        remove: true,
        upload: undefined,
        previewURL: undefined,
      },
      () => {
        this.props.onChange(
          {
            bkey: this.state.bkey,
            remove: this.state.remove,
            upload: this.state.upload,
          },
          this.props.instanceID,
          this.state.previewURL
        )
      }
    )
  }

  public selectFile = async () => {
    if (this.fileSelector) {
      this.fileSelector.click()
    }
  }

  private openModal = () => {
    this.setState({ showModal: true })
  }

  private closeModal = async () => {
    this.setState({ showModal: false })
  }

  private modal() {
    return (
      <Modal.Window className="c-image-viewer-modal" isOpen={this.state.showModal} onClose={this.closeModal} center={false} size="fluid">
        <Modal.Content>{this.props.bkey ? <img alt="" src={uploadedImageURL(this.props.bkey)} /> : <img alt="" src={this.state.previewURL} />}</Modal.Content>

        <Modal.Footer>
          <Button variant="tertiary" onClick={this.closeModal}>
            Close
          </Button>
        </Modal.Footer>
      </Modal.Window>
    )
  }

  public render() {
    const isUploading = !!this.state.upload
    const hasFile = (!this.state.remove && this.props.bkey) || isUploading

    return (
      <ValidationContext.Consumer>
        {(ctx) => (
          <div
            className={classSet({
              "c-form-field": true,
              "c-image-upload": true,
              "m-error": hasError(this.props.field, ctx.error),
            })}
          >
            {this.modal()}
            {this.props.label && <label htmlFor={`input-${this.props.field}`}>{this.props.label}</label>}

            {hasFile && (
              <div className="preview h-20">
                <img alt="" onClick={this.openModal} src={this.state.previewURL} />
                {!this.props.disabled && (
                  <Button onClick={this.removeFile} variant="danger">
                    X
                  </Button>
                )}
              </div>
            )}

            <input ref={(e) => (this.fileSelector = e)} type="file" onChange={this.fileChanged} accept="image/*" />
            {!hasFile && (
              <Button onClick={this.selectFile} disabled={this.props.disabled}>
                <Icon sprite={IconPhotograph} />
              </Button>
            )}
            <DisplayError fields={[this.props.field]} error={ctx.error} />
            {this.props.children}
          </div>
        )}
      </ValidationContext.Consumer>
    )
  }
}
