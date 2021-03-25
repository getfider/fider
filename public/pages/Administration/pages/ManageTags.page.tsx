import "./ManageTags.page.scss"

import React from "react"
import { Button, Segment, List, ListItem, Heading } from "@fider/components"

import { Tag } from "@fider/models"
import { actions, Failure, Fider } from "@fider/services"
import { FaTags } from "react-icons/fa"
import { AdminBasePage } from "../components/AdminBasePage"
import { TagFormState, TagForm } from "../components/TagForm"
import { TagListItem } from "../components/TagListItem"

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
  public icon = FaTags
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
        <Segment>
          <TagForm onSave={this.saveNewTag} onCancel={this.cancelAdd} />
        </Segment>
      ) : (
        <Button color="positive" onClick={this.addNew}>
          Add new
        </Button>
      ))

    return (
      <>
        <Segment>
          <List divided={true}>
            <ListItem>
              <Heading size="small" title="Public Tags" subtitle="These tags are visible to all visitors." />
            </ListItem>
            {publicTaglist.length === 0 ? <ListItem>There aren’t any public tags yet.</ListItem> : publicTaglist}
          </List>
        </Segment>

        <Segment>
          <List divided={true}>
            <ListItem>
              <Heading size="small" title="Private Tags" subtitle="These tags are only visible for members of this site." />
            </ListItem>
            {privateTagList.length === 0 ? <ListItem>There aren’t any private tags yet.</ListItem> : privateTagList}
          </List>
        </Segment>
        {form}
      </>
    )
  }
}
