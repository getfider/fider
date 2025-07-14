import React from "react"
import { Modal, Button, CloseIcon, Icon } from "@fider/components"
import { Trans } from "@lingui/react/macro"
import { HStack, VStack } from "./layout"
import IconRss from "@fider/assets/images/heroicons-rss.svg"
import IconClipboard from "@fider/assets/images/heroicons-clipboard.svg"

interface RSSModalProps {
  isOpen: boolean
  onClose: () => void
  url: string
}

export const RSSModal: React.FC<RSSModalProps> = (props) => {
  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(props.url)
    } catch (err) {
      // Fallback for older browsers
      const textArea = document.createElement("textarea")
      textArea.value = props.url
      document.body.appendChild(textArea)
      textArea.focus()
      textArea.select()
      document.execCommand("copy")
      document.body.removeChild(textArea)
    }
  }

  return (
    <Modal.Window isOpen={props.isOpen} onClose={props.onClose} size="small">
      <Modal.Header>
        <HStack justify="between">
          <HStack spacing={2}>
            <Icon sprite={IconRss} className="h-6 text-orange-500" />
            <h3 className="text-lg font-medium">
              <Trans id="modal.rss.title">Subscribe to ATOM feed</Trans>
            </h3>
          </HStack>
          <CloseIcon closeModal={props.onClose} />
        </HStack>
      </Modal.Header>
      <Modal.Content>
        <VStack spacing={2}>
          <p className="text-sm">
            <Trans id="modal.rss.description">To subscribe to this ATOM feed, copy and paste this URL into your RSS/ATOM reader.</Trans>
          </p>
          <div className="bg-gray-200 px-1 py-2 rounded border">
            <div className="flex items-center justify-between text-left">
              <p className="px-2 text-xs text-gray-800 mr-2 nowrap mb-0">{props.url}</p>
              <Button style={{ padding: "0 4px" }} className="text-center hover" variant="tertiary" onClick={copyToClipboard}>
                <Icon height="18" width="18" sprite={IconClipboard} />
              </Button>
            </div>
          </div>
        </VStack>
      </Modal.Content>
    </Modal.Window>
  )
}
