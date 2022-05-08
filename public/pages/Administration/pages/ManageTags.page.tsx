import React from "react"
import { Button } from "@fider/components"

import { Tag } from "@fider/models"
import { actions, Failure, Fider } from "@fider/services"
import { AdminBasePage } from "../components/AdminBasePage"
import { TagFormState, TagForm } from "../components/TagForm"
import { TagListItem } from "../components/TagListItem"
import { VStack } from "@fider/components/layout"

interface ManageTagsPageProps {
  tags: Tag[]
}

interface ManageTagsPageState {
  isAdding: boolean
  allTags: Tag[]
  deleting?: number
  editing?: number
}

const tagSorter = (t1: Tag, t2: Tag) => {
  if (t1.name < t2.name) {
    return -1
  } else if (t1.name > t2.name) {
    return 1
  }
  return 0
}

export default class ManageTagsPage extends AdminBasePage<ManageTagsPageProps, ManageTagsPageState> {
  public id = "p-admin-tags"
  public name = "tags"
  public title = "Tags"
  public subtitle = "Manage your site tags"

  constructor(props: ManageTagsPageProps) {
    super(props)
    this.state = {
      isAdding: false,
      allTags: this.props.tags,
    }
  }

  private addNew = async () => {
    this.setState({
      isAdding: true,
      deleting: undefined,
      editing: undefined,
    })
  }

  private cancelAdd = () => {
    this.setState({ isAdding: false })
  }

  private saveNewTag = async (data: TagFormState): Promise<Failure | undefined> => {
    const result = await actions.createTag(data.name, data.color, data.isPublic)
    if (result.ok) {
      this.setState({
        isAdding: false,
        allTags: this.state.allTags.concat(result.data).sort(tagSorter),
      })
    } else {
      return result.error
    }
  }

  private handleTagDeleted = (tag: Tag) => {
    const idx = this.state.allTags.indexOf(tag)
    this.setState({
      allTags: this.state.allTags.splice(idx, 1) && this.state.allTags,
    })
  }

  private handleTagEdited = () => {
    this.setState({
      allTags: this.state.allTags.sort(tagSorter),
    })
  }

  private getTagList(filter: (tag: Tag) => boolean) {
    return this.state.allTags.filter(filter).map((t) => {
      return <TagListItem key={t.id} tag={t} onTagDeleted={this.handleTagDeleted} onTagEdited={this.handleTagEdited} />
    })
  }

  public content() {
    const publicTaglist = this.getTagList((t) => t.isPublic)
    const privateTagList = this.getTagList((t) => !t.isPublic)

    const form =
      Fider.session.user.isAdministrator &&
      (this.state.isAdding ? (
        <TagForm onSave={this.saveNewTag} onCancel={this.cancelAdd} />
      ) : (
        <Button variant="secondary" onClick={this.addNew}>
          Add new
        </Button>
      ))

    return (
      <VStack spacing={8}>
        <div>
          <h2 className="text-display">Public Tags</h2>
          <p className="text-muted">These tags are visible to all visitors.</p>
          <VStack spacing={4} divide={true}>
            {publicTaglist.length === 0 ? <p className="text-muted">There aren’t any public tags yet.</p> : publicTaglist}
          </VStack>
        </div>
        <div>
          <h2 className="text-display">Private Tags</h2>
          <p className="text-muted">These tags are only visible for members of this site.</p>
          <VStack spacing={4} divide={true}>
            {privateTagList.length === 0 ? <p className="text-muted">There aren’t any private tags yet.</p> : privateTagList}
          </VStack>
        </div>
        <div>{form}</div>
      </VStack>
    )
  }
}
