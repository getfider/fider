import React from "react";
import { ImageUploader } from "./ImageUploader";
import "./MultiImageUploader.scss";
import { ImageUpload } from "@fider/models";

interface MultiImageUploaderProps {
  maxUploads: number;
  previewMaxWidth: number;
  onChange?: (uploads: ImageUpload[]) => void;
}

interface MultiImageUploaderInstances {
  [key: string]: {
    element: JSX.Element;
    upload?: ImageUpload;
  };
}

interface MultiImageUploaderState {
  count: number;
  instances: MultiImageUploaderInstances;
}

export class MultiImageUploader extends React.Component<MultiImageUploaderProps, MultiImageUploaderState> {
  constructor(props: MultiImageUploaderProps) {
    super(props);

    const instances = {};
    this.addNewElement(instances);
    this.state = {
      instances,
      count: 1
    };
  }

  private imageUploaded = (upload: ImageUpload, instanceId: string) => {
    const instances = { ...this.state.instances };
    let count = this.state.count;
    if (upload.remove) {
      delete instances[instanceId];
      if (--count === this.props.maxUploads) {
        this.addNewElement(instances);
      }
    } else {
      instances[instanceId].upload = upload;
      count++;
      if (count <= this.props.maxUploads) {
        this.addNewElement(instances);
      }
    }
    this.setState({ instances, count }, this.triggerOnChange);
  };

  private triggerOnChange() {
    if (this.props.onChange) {
      const uploads = Object.keys(this.state.instances)
        .map(k => this.state.instances[k].upload)
        .filter(x => !!x) as ImageUpload[];
      this.props.onChange(uploads);
    }
  }

  private addNewElement(instances: MultiImageUploaderInstances) {
    const id = btoa(Math.random().toString());
    instances[id] = {
      element: (
        <ImageUploader
          key={id}
          instanceId={id}
          field="attachment"
          previewMaxWidth={this.props.previewMaxWidth}
          onChange={this.imageUploaded}
        />
      )
    };
  }

  public render() {
    const elements = Object.keys(this.state.instances).map(k => this.state.instances[k].element);
    return <div className="c-multi-image-uploader">{elements}</div>;
  }
}
