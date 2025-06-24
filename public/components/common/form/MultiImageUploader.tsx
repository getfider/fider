import React from "react"
import { ImageUploader } from "./ImageUploader"
import { ImageUpload } from "@fider/models"
import { ValidationContext, hasError, DisplayError } from "@fider/components"
import { classSet } from "@fider/services"

import "./MultiImageUploader.scss"

interface MultiImageUploaderProps {
  field: string
  maxUploads: number
  bkeys?: string[]
  addImageButton?: React.ReactNode
  onChange?: (uploads: ImageUpload[]) => void
  noPadding?: boolean
}

interface MultiImageUploaderInstances {
  [key: string]: {
    element: React.JSX.Element
    upload?: ImageUpload
  }
}

interface MultiImageUploaderState {
  instances: MultiImageUploaderInstances
  removed: ImageUpload[]
}

export class MultiImageUploader extends React.Component<MultiImageUploaderProps, MultiImageUploaderState> {
  constructor(props: MultiImageUploaderProps) {
    super(props)

    const instances = {}
    if (props.bkeys) {
      for (const bkey of props.bkeys) {
        this.addNewElement(instances, bkey)
      }
    }

    if (Object.keys(instances).length <= this.props.maxUploads) {
      this.addNewElement(instances)
    }

    this.state = { instances, removed: [] }
  }

  private imageUploaded = (upload: ImageUpload, instanceID: string, ctx?: React.ContextType<typeof ValidationContext>) => {
    // Clear error for this field when user interacts with it
    if (ctx?.clearError && hasError(this.props.field, ctx.error)) {
      ctx.clearError(this.props.field)
    }

    const instances = { ...this.state.instances }
    const removed = [...this.state.removed]
    if (upload.remove) {
      if (upload.bkey) {
        removed.push(upload)
      }
      delete instances[instanceID]
    } else {
      instances[instanceID].upload = upload
    }
    if (Object.keys(instances).length < this.props.maxUploads && !this.hasUploadButton(instances)) {
      this.addNewElement(instances)
    }
    this.setState({ instances, removed }, this.triggerOnChange)
  }

  private triggerOnChange() {
    if (this.props.onChange) {
      const uploads = Object.keys(this.state.instances)
        .map((k) => this.state.instances[k].upload)
        .concat(this.state.removed)
        .filter((x) => !!x) as ImageUpload[]
      this.props.onChange(uploads)
    }
  }

  private addNewElement(instances: MultiImageUploaderInstances, bkey?: string) {
    const id = btoa(Math.random().toString())
    instances[id] = {
      element: (
        <ValidationContext.Consumer>
          {(ctx) => (
            <ImageUploader
              key={id}
              bkey={bkey}
              instanceID={id}
              field="attachment"
              onChange={(upload, instanceID) => this.imageUploaded(upload, instanceID || "", ctx)}
              addImageButton={this.props.addImageButton}
            />
          )}
        </ValidationContext.Consumer>
      ),
    }
  }

  private hasUploadButton(instances: MultiImageUploaderInstances) {
    for (const instance of Object.values(instances)) {
      if (instance.upload == null) {
        return true
      }
    }
    return false
  }

  public render() {
    const elements = Object.keys(this.state.instances).map((k) => this.state.instances[k].element)
    return (
      <ValidationContext.Consumer>
        {(ctx) => (
          <div
            className={classSet({
              "c-form-field": true,
              "c-multi-image-uploader": true,
              "m-error": hasError(this.props.field, ctx.error),
              "pt-0": this.props.noPadding,
            })}
          >
            <div className="c-multi-image-uploader-instances">{elements}</div>
            <DisplayError fields={[this.props.field]} error={ctx.error} />
          </div>
        )}
      </ValidationContext.Consumer>
    )
  }
}
