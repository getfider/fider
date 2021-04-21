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
  previewMaxWidth: number
  onChange?: (uploads: ImageUpload[]) => void
}

interface MultiImageUploaderInstances {
  [key: string]: {
    element: JSX.Element
    upload?: ImageUpload
  }
}

interface MultiImageUploaderState {
  count: number
  instances: MultiImageUploaderInstances
  removed: ImageUpload[]
}

export class MultiImageUploader extends React.Component<MultiImageUploaderProps, MultiImageUploaderState> {
  constructor(props: MultiImageUploaderProps) {
    super(props)

    let count = 1
    const instances = {}
    if (props.bkeys) {
      for (const bkey of props.bkeys) {
        count++
        this.addNewElement(instances, bkey)
      }
    }

    if (count <= this.props.maxUploads) {
      count++
      this.addNewElement(instances)
    }

    this.state = { instances, count, removed: [] }
  }

  private imageUploaded = (upload: ImageUpload, instanceID: string) => {
    const instances = { ...this.state.instances }
    const removed = [...this.state.removed]
    let count = this.state.count
    if (upload.remove) {
      if (upload.bkey) {
        removed.push(upload)
      }
      delete instances[instanceID]
      if (--count === this.props.maxUploads) {
        this.addNewElement(instances)
      }
    } else {
      instances[instanceID].upload = upload
      if (count++ <= this.props.maxUploads) {
        this.addNewElement(instances)
      }
    }
    this.setState({ instances, count, removed }, this.triggerOnChange)
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
        <ImageUploader key={id} bkey={bkey} instanceID={id} field="attachment" previewMaxWidth={this.props.previewMaxWidth} onChange={this.imageUploaded} />
      ),
    }
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
