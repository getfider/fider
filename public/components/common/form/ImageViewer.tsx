import React, { useState } from "react"
import { uploadedImageURL } from "@fider/services"
import { Modal, Button, Loader } from "@fider/components"

import "./ImageViewer.scss"
import { Trans } from "@lingui/macro"

interface ImageViewerModalProps {
  isOpen: boolean
  imgSrc: string | undefined
  loadedPreview: boolean
  onPreviewLoad: () => void
  onClose: () => void
}

const ImageViewerModal = (props: ImageViewerModalProps) => (
  <Modal.Window className="c-image-viewer-modal" isOpen={props.isOpen} onClose={props.onClose} center={false} size="fluid">
    <Modal.Content>
      {!props.loadedPreview && <Loader />}
      <img alt="" onLoad={props.onPreviewLoad} src={props.imgSrc} />
    </Modal.Content>

    <Modal.Footer>
      <Button variant="tertiary" onClick={props.onClose}>
        <Trans id="action.close">Close</Trans>
      </Button>
    </Modal.Footer>
  </Modal.Window>
)

interface ImageViewerProps {
  bkey: string
}

export const ImageViewer = (props: ImageViewerProps) => {
  const [showModal, setShowModal] = useState(false)
  const [loadedThumbnail, setLoadedThumbnail] = useState(false)
  const [loadedPreview, setLoadedPreview] = useState(false)

  const openModal = () => {
    if (loadedThumbnail) {
      setShowModal(true)
    }
  }

  const closeModal = () => {
    setShowModal(false)
  }

  const onThumbnailLoad = () => {
    setLoadedThumbnail(true)
  }

  const onPreviewLoad = () => {
    setLoadedPreview(true)
  }

  return (
    <div className="c-image-viewer">
      <ImageViewerModal
        onPreviewLoad={onPreviewLoad}
        isOpen={showModal}
        onClose={closeModal}
        imgSrc={uploadedImageURL(props.bkey, 1500)}
        loadedPreview={loadedPreview}
      />
      {!loadedThumbnail && <Loader />}
      <img alt="" onClick={openModal} onLoad={onThumbnailLoad} src={uploadedImageURL(props.bkey, 200)} />
    </div>
  )
}
