import React, { useEffect, useState } from "react"
import { Trans } from "@lingui/react/macro"
import { Modal, Button, Loader } from "@fider/components"
import { notify } from "@fider/services"
import { listDisposableUsers, bulkDeleteDisposableUsers, DisposableUserRow } from "@fider/services/actions/disposable"

import "./DisposableUsersModal.scss"

interface Props {
  onClose: () => void
}

export const DisposableUsersModal: React.FC<Props> = ({ onClose }) => {
  const [loading, setLoading] = useState(true)
  const [users, setUsers] = useState<DisposableUserRow[]>([])
  const [total, setTotal] = useState(0)
  const [selected, setSelected] = useState<Record<number, boolean>>({})
  const [confirming, setConfirming] = useState(false)
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      const r = await listDisposableUsers()
      if (cancelled) return
      if (r.ok) {
        setUsers(r.data.users)
        setTotal(r.data.total)
        const sel: Record<number, boolean> = {}
        r.data.users.forEach((u) => (sel[u.userID] = true))
        setSelected(sel)
      } else {
        notify.error("Failed to load disposable users.")
      }
      setLoading(false)
    })()
    return () => {
      cancelled = true
    }
  }, [])

  const toggle = (id: number) => setSelected({ ...selected, [id]: !selected[id] })

  const toggleAll = () => {
    const allChecked = users.every((u) => selected[u.userID])
    const sel: Record<number, boolean> = {}
    users.forEach((u) => (sel[u.userID] = !allChecked))
    setSelected(sel)
  }

  const selectedIds = Object.entries(selected)
    .filter(([, v]) => v)
    .map(([k]) => Number(k))

  const summary = users.reduce(
    (acc, u) => ({
      votes: acc.votes + u.voteCount,
      posts: acc.posts + u.postCount,
      comments: acc.comments + u.commentCount,
    }),
    { votes: 0, posts: 0, comments: 0 }
  )

  const allChecked = users.length > 0 && users.every((u) => selected[u.userID])

  const performDelete = async () => {
    setBusy(true)
    const r = await bulkDeleteDisposableUsers(selectedIds)
    setBusy(false)
    if (r.ok) {
      notify.success(`Deleted ${r.data.deleted} users.`)
      onClose()
    } else {
      notify.error("Bulk delete failed.")
    }
  }

  return (
    <Modal.Window isOpen={true} size="large" canClose={!busy} onClose={onClose}>
      <Modal.Header>
        <Trans id="admin.members.disposableModal.title">Delete disposable accounts</Trans>
      </Modal.Header>
      <Modal.Content>
        {loading && <Loader />}
        {!loading && (
          <div className="c-disposable-modal">
            <p>
              {total} users matched ({summary.votes} votes, {summary.posts} posts, {summary.comments} comments).
            </p>
            {total > users.length && (
              <p className="text-muted text-sm">
                Showing first {users.length} of {total}. Run again after deleting these.
              </p>
            )}
            {users.length === 0 ? (
              <p className="text-muted">No disposable-email accounts found.</p>
            ) : (
              <table className="c-disposable-modal__table">
                <thead>
                  <tr>
                    <th>
                      <input type="checkbox" checked={allChecked} onChange={toggleAll} />
                    </th>
                    <th>Email</th>
                    <th>Name</th>
                    <th>Votes</th>
                    <th>Posts</th>
                    <th>Comments</th>
                  </tr>
                </thead>
                <tbody>
                  {users.map((u) => (
                    <tr key={u.userID}>
                      <td>
                        <input type="checkbox" checked={!!selected[u.userID]} onChange={() => toggle(u.userID)} />
                      </td>
                      <td>{u.email}</td>
                      <td>{u.name}</td>
                      <td>{u.voteCount}</td>
                      <td>{u.postCount}</td>
                      <td>{u.commentCount}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
            {confirming && (
              <div className="c-disposable-modal__confirm mt-4">
                <p>
                  This will permanently delete {selectedIds.length} users and their votes. Posts and comments will remain but be anonymized. This cannot be undone.
                </p>
              </div>
            )}
          </div>
        )}
      </Modal.Content>
      <Modal.Footer>
        {!confirming && (
          <>
            <Button variant="tertiary" disabled={busy} onClick={onClose}>
              <Trans id="action.cancel">Cancel</Trans>
            </Button>
            <Button variant="danger" disabled={loading || selectedIds.length === 0} onClick={() => setConfirming(true)}>
              <Trans id="admin.members.disposableModal.deleteSelected">Delete selected</Trans>
            </Button>
          </>
        )}
        {confirming && (
          <>
            <Button variant="tertiary" disabled={busy} onClick={() => setConfirming(false)}>
              <Trans id="action.cancel">Cancel</Trans>
            </Button>
            <Button variant="danger" disabled={busy} onClick={performDelete}>
              <Trans id="action.confirm">Confirm</Trans>
            </Button>
          </>
        )}
      </Modal.Footer>
    </Modal.Window>
  )
}
