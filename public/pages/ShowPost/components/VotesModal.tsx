import React, { useEffect, useState } from "react"
import { Post, Vote } from "@fider/models"
import { Modal, Button, Loader, Avatar, UserName, Moment, Input } from "@fider/components"
import { actions } from "@fider/services"
import { useFider } from "@fider/hooks"
import IconSearch from "@fider/assets/images/heroicons-search.svg"
import IconX from "@fider/assets/images/heroicons-x.svg"
import { HStack, VStack } from "@fider/components/layout"
import { t, Trans } from "@lingui/macro"

interface VotesModalProps {
  isOpen: boolean
  post: Post
  onClose?: () => void
}

export const VotesModal: React.FC<VotesModalProps> = (props) => {
  const [isLoading, setIsLoading] = useState(false)
  const [query, setQuery] = useState("")
  const [allVotes, setAllVotes] = useState<Vote[]>([])
  const [filteredVotes, setFilteredVotes] = useState<Vote[]>([])

  const fider = useFider()

  useEffect(() => {
    if (props.isOpen) {
      actions.listVotes(props.post.number).then((response) => {
        if (response.ok) {
          setAllVotes(response.data)
          setFilteredVotes(response.data)
          setIsLoading(false)
        }
      })
    }
  }, [props.isOpen])

  const closeModal = async () => {
    if (props.onClose) {
      props.onClose()
    }
  }

  const clearSearch = () => {
    handleSearchFilterChanged("")
  }

  const handleSearchFilterChanged = (query: string) => {
    const votes = allVotes.filter((x) => x.user.name.toLowerCase().indexOf(query.toLowerCase()) >= 0)
    setQuery(query)
    setFilteredVotes(votes)
  }

  return (
    <Modal.Window isOpen={props.isOpen} center={false} onClose={closeModal}>
      <Modal.Content>
        {isLoading && <Loader />}
        {!isLoading && (
          <>
            <Input
              field="query"
              icon={query ? IconX : IconSearch}
              onIconClick={query ? clearSearch : undefined}
              placeholder={t({ id: "modal.showvotes.query.placeholder", message: "Search for users by name..." })}
              value={query}
              onChange={handleSearchFilterChanged}
            />
            <VStack spacing={2} className="h-max-5xl overflow-scroll">
              {filteredVotes.map((x) => (
                <HStack key={x.user.id} justify="between">
                  <HStack>
                    <Avatar user={x.user} />
                    <VStack spacing={0}>
                      <UserName user={x.user} />
                      <span className="text-muted">{x.user.email}</span>
                    </VStack>
                  </HStack>
                  <span className="text-muted">
                    <Moment locale={fider.currentLocale} date={x.createdAt} />
                  </span>
                </HStack>
              ))}
              {filteredVotes.length === 0 && (
                <p className="text-muted">
                  <Trans id="modal.showvotes.message.zeromatches">
                    No users found matching <strong>{query}</strong>.
                  </Trans>
                </p>
              )}
            </VStack>
          </>
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
