import "./VotesModal.scss"

import React from "react"
import { Post, Vote } from "@fider/models"
import { Modal, Button, Loader, List, ListItem, Avatar, UserName, Moment, Input } from "@fider/components"
import { actions } from "@fider/services"
import { FaTimes, FaSearch } from "react-icons/fa"

interface VotesModalProps {
  isOpen: boolean
  post: Post
  onClose?: () => void
}

interface VotesModalState {
  searchText: string
  allVotes: Vote[]
  filteredVotes: Vote[]
  isLoading: boolean
  query: string
}

export class VotesModal extends React.Component<VotesModalProps, VotesModalState> {
  constructor(props: VotesModalProps) {
    super(props)
    this.state = {
      searchText: "",
      query: "",
      allVotes: [],
      filteredVotes: [],
      isLoading: true,
    }
  }

  public componentDidUpdate(prevProps: VotesModalProps) {
    if (this.props.isOpen && !prevProps.isOpen) {
      actions.listVotes(this.props.post.number).then((response) => {
        if (response.ok) {
          this.setState({
            allVotes: response.data,
            filteredVotes: response.data,
            isLoading: false,
          })
        }
      })
    }
  }

  private closeModal = async () => {
    if (this.props.onClose) {
      this.props.onClose()
    }
  }

  private clearSearch = () => {
    this.handleSearchFilterChanged("")
  }

  private handleSearchFilterChanged = (query: string) => {
    const votes = this.state.allVotes.filter((x) => x.user.name.toLowerCase().indexOf(query.toLowerCase()) >= 0)
    this.setState({ query, filteredVotes: votes })
  }

  public render() {
    return (
      <Modal.Window className="c-votes-modal" isOpen={this.props.isOpen} center={false} onClose={this.closeModal}>
        <Modal.Content>
          {this.state.isLoading && <Loader />}
          {!this.state.isLoading && (
            <>
              <Input
                field="query"
                icon={this.state.query ? FaTimes : FaSearch}
                onIconClick={this.state.query ? this.clearSearch : undefined}
                placeholder="Search for users by name..."
                value={this.state.query}
                onChange={this.handleSearchFilterChanged}
              />
              <List hover={true}>
                {this.state.filteredVotes.map((x) => (
                  <ListItem key={x.user.id}>
                    <Avatar user={x.user} />
                    <span className="l-user">
                      <UserName user={x.user} />
                      <span className="info">{x.user.email}</span>
                    </span>
                    <span className="l-date info">
                      <Moment date={x.createdAt} />
                    </span>
                  </ListItem>
                ))}
              </List>
            </>
          )}
        </Modal.Content>

        <Modal.Footer>
          <Button color="cancel" onClick={this.closeModal}>
            Close
          </Button>
        </Modal.Footer>
      </Modal.Window>
    )
  }
}
