import React, { useState, useContext, memo } from "react"
import { Post, Tag } from "@fider/models"
import { actions } from "@fider/services"
import { ShowTag, Button, Icon } from "@fider/components"
import IconX from "@fider/assets/images/heroicons-x.svg"
import IconCheckCircle from "@fider/assets/images/heroicons-check-circle.svg"
import { TagListItem } from "./TagListItem"
import { Dropdown, DropdownContext } from "../../../components/common/Dropdown"
import { useFider } from "@fider/hooks"

import { HStack, VStack } from "@fider/components/layout"
import { Trans } from "@lingui/react/macro"

export interface TagsPanelProps {
  post: Post
  tags: Tag[]
}

const DropdownContent = memo(
  (props: { searchQuery: string; setSearchQuery: (q: string) => void; filteredTags: Tag[]; assignedTags: Tag[]; assignOrUnassignTag: (tag: Tag) => void }) => {
    const { searchQuery, setSearchQuery, filteredTags, assignedTags, assignOrUnassignTag } = props
    const dropdownContext = useContext(DropdownContext)

    return (
      <div
        style={{
          backgroundColor: "#fff",
          padding: "1rem",
          borderRadius: "8px",
          width: "20rem",
          boxSizing: "border-box",
        }}
      >
        <VStack>
          <div
            style={{
              display: "flex",
              alignItems: "center",
              borderBottom: "1px solid #eaeaea",
              paddingBottom: "0.5rem",
              marginBottom: "0.5rem",
              width: "12.5rem",
            }}
          >
            <h3 style={{ margin: 0, fontSize: "1rem", fontWeight: 500 }}>TAGS</h3>
            <Icon sprite={IconCheckCircle} className="h-4 text-yellow-500" />
          </div>
          <div style={{ position: "relative", width: "100%" }}>
            <input
              type="text"
              placeholder="Search tags..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              style={{
                width: "70%",
                padding: "0.75rem 2.5rem 0.75rem 0.75rem",
                border: "1px solid #ccc",
                borderRadius: "4px",
                boxSizing: "border-box",
              }}
            />
            {searchQuery && (
              <div
                style={{
                  position: "absolute",
                  right: "6rem",
                  top: "50%",
                  transform: "translateY(-50%)",
                  cursor: "pointer",
                }}
                onClick={() => setSearchQuery("")}
              >
                <Icon sprite={IconX} className="c-hint__close h-5" />
              </div>
            )}
          </div>
          <div style={{ maxHeight: "200px", overflowY: "auto" }}>
            {filteredTags.map((tag) => (
              <TagListItem key={tag.id} tag={tag} assigned={assignedTags.indexOf(tag) >= 0} onClick={assignOrUnassignTag} />
            ))}
          </div>
          <Button style={{ position: "relative", left: "70px", width: "50px" }} variant="danger" size="small" onClick={() => dropdownContext?.close()}>
            <Trans id="action.close">Close</Trans>
          </Button>
        </VStack>
      </div>
    )
  }
)
DropdownContent.displayName = "EditTags_DropdownContent"

export const TagsPanel = (props: TagsPanelProps) => {
  const fider = useFider()
  const canEdit = fider.session.isAuthenticated && fider.session.user.isCollaborator && props.tags.length > 0

  const [assignedTags, setAssignedTags] = useState(props.tags.filter((t) => props.post.tags.indexOf(t.slug) >= 0))
  const [searchQuery, setSearchQuery] = useState("")

  const assignOrUnassignTag = async (tag: Tag) => {
    const idx = assignedTags.indexOf(tag)
    let nextAssignedTags: Tag[] = []

    if (idx >= 0) {
      const response = await actions.unassignTag(tag.slug, props.post.number)
      if (response.ok) {
        nextAssignedTags = [...assignedTags]
        nextAssignedTags.splice(idx, 1)
      }
    } else {
      const response = await actions.assignTag(tag.slug, props.post.number)
      if (response.ok) {
        nextAssignedTags = [...assignedTags, tag]
      }
    }

    setAssignedTags(nextAssignedTags)
    setSearchQuery("")
  }

  if (!canEdit && assignedTags.length === 0) {
    return null
  }

  if (fider.isReadOnly) {
    return (
      <VStack>
        <HStack spacing={2} className="text-category">
          <Trans id="label.tags">Tags</Trans>
        </HStack>
        <HStack spacing={2} align="center">
          {assignedTags.length > 0 && assignedTags.map((tag) => <ShowTag key={tag.id} tag={tag} link />)}
        </HStack>
      </VStack>
    )
  }

  const appliedTags = (
    <HStack spacing={2} align="center">
      {assignedTags.length > 0 && assignedTags.map((tag) => <ShowTag key={tag.id} tag={tag} link />)}
    </HStack>
  )

  const filteredTags = props.tags.filter((tag) => tag.name.toLowerCase().includes(searchQuery.toLowerCase()))

  return (
    <VStack>
      {canEdit && (
        <Dropdown
          renderHandle={
            <div className="clickable" style={{ marginBottom: "0.5rem", fontWeight: 500 }}>
              <Trans id="label.edittags">Edit tags</Trans>
            </div>
          }
          position="right"
        >
          <DropdownContent
            searchQuery={searchQuery}
            setSearchQuery={setSearchQuery}
            filteredTags={filteredTags}
            assignedTags={assignedTags}
            assignOrUnassignTag={assignOrUnassignTag}
          />
        </Dropdown>
      )}
      {appliedTags}
    </VStack>
  )
}
