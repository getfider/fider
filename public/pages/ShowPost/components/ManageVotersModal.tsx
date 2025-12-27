import React, { useEffect, useState } from "react"
import { Post, Vote } from "@fider/models"
import { Modal, Button, Loader, Avatar, UserName, Moment, Input, Form, Icon } from "@fider/components"
import { actions, notify, Failure } from "@fider/services"
import { useFider } from "@fider/hooks"
import IconSearch from "@fider/assets/images/heroicons-search.svg"
import IconX from "@fider/assets/images/heroicons-x.svg"
import { HStack, VStack } from "@fider/components/layout"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"

interface ManageVotersModalProps {
  isOpen: boolean
  post: Post
  onClose: () => void
  onVotesChanged: (delta: number) => void
}

export const ManageVotersModal: React.FC<ManageVotersModalProps> = (props) => {
  const fider = useFider()

  // Voter list state
  const [isLoading, setIsLoading] = useState(true)
  const [query, setQuery] = useState("")
  const [allVotes, setAllVotes] = useState<Vote[]>([])
  const [filteredVotes, setFilteredVotes] = useState<Vote[]>([])
  const [pendingRemoveUserId, setPendingRemoveUserId] = useState<number | null>(null)

  // Add voter form state
  const [email, setEmail] = useState("")
  const [name, setName] = useState("")
  const [error, setError] = useState<Failure>()

  useEffect(() => {
    if (props.isOpen) {
      setIsLoading(true)
      actions.listVotes(props.post.number).then((response) => {
        if (response.ok) {
          setAllVotes(response.data)
          setFilteredVotes(response.data)
        }
        setIsLoading(false)
      })
    }
  }, [props.isOpen])

  const closeModal = () => {
    setEmail("")
    setName("")
    setError(undefined)
    setQuery("")
    setPendingRemoveUserId(null)
    props.onClose()
  }

  const clearSearch = () => {
    handleSearchFilterChanged("")
  }

  const handleSearchFilterChanged = (query: string) => {
    const votes = allVotes.filter((x) => x.user.name.toLowerCase().indexOf(query.toLowerCase()) >= 0)
    setQuery(query)
    setFilteredVotes(votes)
  }

  const handleAddVoter = async () => {
    setError(undefined)
    const response = await actions.addVoteOnBehalf(props.post.number, email, name)
    if (response.ok) {
      notify.success(i18n._({ id: "managevoters.add.success", message: `Vote added for ${response.data.user.name}` }))
      // Add to local list
      const newVote: Vote = {
        user: {
          id: response.data.user.id,
          name: response.data.user.name,
          email: response.data.user.email,
          avatarURL: "",
        },
        createdAt: new Date(),
      }
      const newVotes = [newVote, ...allVotes]
      setAllVotes(newVotes)
      setFilteredVotes(newVotes.filter((x) => x.user.name.toLowerCase().indexOf(query.toLowerCase()) >= 0))
      setEmail("")
      setName("")
      props.onVotesChanged(1)
    } else if (response.error) {
      setError(response.error)
    }
  }

  const handleRemoveVoter = async (userId: number) => {
    const response = await actions.removeVoteOnBehalf(props.post.number, userId)
    if (response.ok) {
      notify.success(i18n._({ id: "managevoters.remove.success", message: "Vote removed" }))
      const newVotes = allVotes.filter((v) => v.user.id !== userId)
      setAllVotes(newVotes)
      setFilteredVotes(newVotes.filter((x) => x.user.name.toLowerCase().indexOf(query.toLowerCase()) >= 0))
      setPendingRemoveUserId(null)
      props.onVotesChanged(-1)
    }
  }

  return (
    <Modal.Window isOpen={props.isOpen} center={false} onClose={closeModal} size="large">
      <Modal.Header>
        <Trans id="managevoters.title">Manage Voters</Trans>
      </Modal.Header>
      <Modal.Content>
        {isLoading && <Loader />}
        {!isLoading && (
          <VStack spacing={4}>
            {/* Add Voter Form */}
            <div className="p-4 rounded bg-gray-50">
              <p className="text-muted mb-2">
                <Trans id="managevoters.add.description">Add a vote on behalf of someone</Trans>
              </p>
              <Form error={error}>
                <HStack spacing={2} align="end">
                  <div className="flex-1">
                    <Input
                      field="name"
                      label={i18n._({ id: "label.name", message: "Name" })}
                      value={name}
                      onChange={setName}
                      placeholder={i18n._({ id: "managevoters.name.placeholder", message: "John Doe" })}
                    />
                  </div>
                  <div className="flex-1">
                    <Input
                      field="email"
                      label={i18n._({ id: "label.email", message: "Email" })}
                      value={email}
                      onChange={setEmail}
                      placeholder={i18n._({ id: "managevoters.email.placeholder", message: "user@example.com" })}
                      inputMode="email"
                    />
                  </div>
                  <Button variant="primary" onClick={handleAddVoter} disabled={!email || !name}>
                    <Trans id="action.addvote">Add vote</Trans>
                  </Button>
                </HStack>
              </Form>
            </div>

            {/* Voter List */}
            <Input
              field="query"
              icon={query ? IconX : IconSearch}
              onIconClick={query ? clearSearch : undefined}
              placeholder={i18n._({ id: "modal.showvotes.query.placeholder", message: "Search for users by name..." })}
              value={query}
              onChange={handleSearchFilterChanged}
            />
            <VStack spacing={2} className="h-max-5xl overflow-auto">
              {filteredVotes.map((x) => (
                <HStack key={x.user.id} justify="between" className="p-2 rounded hover:bg-gray-50">
                  <HStack>
                    <Avatar user={x.user} />
                    <VStack spacing={0}>
                      <UserName user={x.user} />
                      <span className="text-muted">{x.user.email}</span>
                    </VStack>
                  </HStack>
                  <HStack spacing={2}>
                    <span className="text-muted">
                      <Moment locale={fider.currentLocale} date={x.createdAt} />
                    </span>
                    {pendingRemoveUserId === x.user.id ? (
                      <HStack spacing={1}>
                        <Button variant="danger" size="small" onClick={() => handleRemoveVoter(x.user.id)}>
                          <Trans id="action.confirm">Confirm</Trans>
                        </Button>
                        <Button variant="tertiary" size="small" onClick={() => setPendingRemoveUserId(null)}>
                          <Trans id="action.cancel">Cancel</Trans>
                        </Button>
                      </HStack>
                    ) : (
                      <Button variant="tertiary" size="small" onClick={() => setPendingRemoveUserId(x.user.id)}>
                        <Icon sprite={IconX} />
                      </Button>
                    )}
                  </HStack>
                </HStack>
              ))}
              {filteredVotes.length === 0 && (
                <p className="text-muted">
                  {query ? (
                    <Trans id="modal.showvotes.message.zeromatches">
                      No users found matching <strong>{query}</strong>.
                    </Trans>
                  ) : (
                    <Trans id="managevoters.novoters">No voters yet.</Trans>
                  )}
                </p>
              )}
            </VStack>
          </VStack>
        )}
      </Modal.Content>

      <Modal.Footer>
        <Button variant="tertiary" onClick={closeModal}>
          <Trans id="action.close">Close</Trans>
        </Button>
      </Modal.Footer>
    </Modal.Window>
  )
}
