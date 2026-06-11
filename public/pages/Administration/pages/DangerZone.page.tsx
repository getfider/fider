import React from "react"

import { Button, Modal, Input, Message } from "@fider/components"
import { AdminBasePage } from "@fider/pages/Administration/components/AdminBasePage"
import { actions, notify, Fider } from "@fider/services"
import { Icon } from "@fider/components"
import IconExclamation from "@fider/assets/images/heroicons-exclamation.svg"

interface DangerZonePageProps {
  isOwner: boolean
  scheduledDeletionAt?: string | null
}

interface DangerZonePageState {
  showModal: boolean
  confirmation: string
  scheduledDeletionAt?: string | null
}

export default class DangerZonePage extends AdminBasePage<DangerZonePageProps, DangerZonePageState> {
  public id = "p-admin-danger"
  public name = "danger-zone"
  public title = "Danger Zone"
  public subtitle = "Permanently delete this site"

  constructor(props: DangerZonePageProps) {
    super(props)
    this.state = {
      showModal: false,
      confirmation: "",
      scheduledDeletionAt: props.scheduledDeletionAt,
    }
  }

  private openModal = () => this.setState({ showModal: true })
  private closeModal = () => this.setState({ showModal: false, confirmation: "" })

  private confirmDelete = async () => {
    const response = await actions.requestTenantDeletion(this.state.confirmation)
    if (response.ok) {
      this.setState({ showModal: false, confirmation: "", scheduledDeletionAt: response.data.scheduledDeletionAt })
      notify.success("Your site is scheduled for deletion. We've emailed you a link to cancel — you have 1 hour.")
    } else {
      notify.error("Failed to schedule deletion. Please try again later.")
    }
  }

  private cancelDeletion = async () => {
    const response = await actions.cancelTenantDeletion()
    if (response.ok) {
      this.setState({ scheduledDeletionAt: null })
      notify.success("The scheduled deletion has been cancelled.")
    } else {
      notify.error("Failed to cancel the scheduled deletion. Please try again later.")
    }
  }

  private renderScheduled() {
    const when = this.state.scheduledDeletionAt ? new Date(this.state.scheduledDeletionAt).toLocaleString() : ""
    return (
      <Message type="error">
        <h4 className="text-title mb-1">This site is scheduled for deletion</h4>
        <p className="mb-2">
          Everything for this site will be permanently deleted at <strong>{when}</strong>. This cannot be undone once it runs.
        </p>
        {this.props.isOwner ? (
          <Button variant="danger" size="small" onClick={this.cancelDeletion}>
            Cancel deletion
          </Button>
        ) : (
          <p className="text-muted">Only the account owner can cancel the deletion.</p>
        )}
      </Message>
    )
  }

  private renderDelete() {
    const subdomain = Fider.session.tenant.subdomain

    if (!this.props.isOwner) {
      return (
        <Message type="warning">
          <h4 className="text-title mb-1">Delete this fider board</h4>
          <p className="text-muted">
            Only the account owner (the person who created this fider board) can delete it. Please contact them if this site needs to be deleted.
          </p>
        </Message>
      )
    }

    return (
      <div>
        <Modal.Window isOpen={this.state.showModal} center={false} onClose={this.closeModal}>
          <Modal.Header>Are you sure you want to delete everything?</Modal.Header>
          <Modal.Content>
            <p>
              To confirm, This will <strong>permanently delete</strong> this fider board and everything in it — all users, posts, comments and votes — and
              cancel any active subscription.
            </p>
            <h4 className="text-title mb-2">How it works</h4>
            <p>
              We'll put the delete request in a queue and send an email to confirm what's happening. After an hour, the delete will be processed and everything
              will be deleted. No going back. No restore from backup. It's all gone.
            </p>
            <p className="mt-2">
              To confirm, type the subdomain of this site (<strong>{subdomain}</strong>) below:
            </p>
            <Input field="confirmation" value={this.state.confirmation} placeholder={subdomain} onChange={(confirmation) => this.setState({ confirmation })} />
          </Modal.Content>
          <Modal.Footer>
            <Button variant="danger" disabled={this.state.confirmation !== subdomain} onClick={this.confirmDelete}>
              Yes, delete this fider board.
            </Button>
            <Button variant="tertiary" onClick={this.closeModal}>
              Cancel
            </Button>
          </Modal.Footer>
        </Modal.Window>

        <h4 className="text-title mb-1 flex items-center gap-1 text-red-700">
          <Icon sprite={IconExclamation} height="24" />
          Delete this fider board
        </h4>
        <p className="text-muted text-red">
          Careful! This will permanently delete this fider board and <strong>everything</strong> in it — all users, posts, comments and votes — and cancels any
          active subscription to fider Pro.
        </p>
        <Button variant="danger" size="small" onClick={this.openModal}>
          I understand, delete this fider board
        </Button>
      </div>
    )
  }

  public content() {
    return this.state.scheduledDeletionAt ? this.renderScheduled() : this.renderDelete()
  }
}
