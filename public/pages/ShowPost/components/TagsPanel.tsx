import React, { useEffect, useRef, useState } from "react"
import { Post, Tag } from "@fider/models"
import { actions, sortTags } from "@fider/services"
import { ShowTag } from "@fider/components"
import { useFider } from "@fider/hooks"

import { HStack, VStack } from "@fider/components/layout"
import { Trans } from "@lingui/react/macro"

import "./TagsPanel.scss"

export interface TagsPanelProps {
  post: Post
  tags: Tag[]
}

export const TagsPanel = (props: TagsPanelProps) => {
  const fider = useFider()
  const canEdit = fider.session.isAuthenticated && fider.session.user.isCollaborator && props.tags.length > 0

  const [isEditing, setIsEditing] = useState(false)
  const [assignedTags, setAssignedTags] = useState(props.tags.filter((t) => props.post.tags.indexOf(t.slug) >= 0))
  const [query, setQuery] = useState("")

  const dropdownRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)

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
  }

  const onSubtitleClick = () => {
    if (canEdit) {
      setIsEditing(!isEditing)
      // Immediately focus on the input element when editing starts
      if (inputRef.current) {
        inputRef.current.focus()
      }
      setQuery("")
    }
  }

  const handleOptionClick = (tag: Tag) => {
    assignOrUnassignTag(tag)
    // Keep focus on the input element after selection
    if (inputRef.current) {
      inputRef.current.focus()
    }
  }

  const filteredOptions = props.tags.filter(
    (option) => option.name.toLowerCase().includes(query.toLowerCase()) && !assignedTags.some((tag) => tag.slug === option.slug)
  )

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: { target: any }) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsEditing(false)
      }
    }

    document.addEventListener("mousedown", handleClickOutside)
    return () => {
      document.removeEventListener("mousedown", handleClickOutside)
    }
  }, [])

  useEffect(() => {
    if (isEditing && inputRef.current) {
      inputRef.current.focus()
    }
  }, [isEditing])

  if (!canEdit && assignedTags.length === 0) {
    return null
  }

  const tagsList = (
    <div className="tags-list">
      {assignedTags.length > 0 && sortTags(assignedTags).map((tag) => <ShowTag key={tag.id} tag={tag} link />)}
      <HStack spacing={1} align="center" className="clickable" onClick={onSubtitleClick}>
        <span>
          <Trans id="label.edittags">Edit tags</Trans>
        </span>
      </HStack>
    </div>
  )

  // Dynamic multiselect dropdown for tags selection
  const editTagsList = props.tags.length > 0 && (
    <div className="dropdown-wrapper" ref={dropdownRef}>
      {/* Selected options and search input */}
      <div className="selected-options-container">
        {sortTags(assignedTags).map((tag) => (
          <div key={tag.id} className="selected-option">
            <ShowTag key={tag.id} tag={tag} link />
            <button onClick={() => handleOptionClick(tag)} className="remove-button">
              x
            </button>
          </div>
        ))}
        {/* Search box to enter query string */}
        <Trans id="label.selecttags">
          <input
            type="text"
            value={query}
            ref={inputRef}
            onChange={(e) => setQuery(e.target.value)}
            className="search-input"
            placeholder={assignedTags.length ? "" : "Select tags..."}
          />
        </Trans>
      </div>

      {/* Dropdown options after items are filtered */}
      {isEditing && (
        <div className="options-container">
          {filteredOptions.length > 0 ? (
            sortTags(filteredOptions).map((tag) => (
              <div key={tag.id} className="option-item" onClick={() => handleOptionClick(tag)}>
                <ShowTag key={tag.id} tag={tag} />
              </div>
            ))
          ) : (
            <div className="no-options">
              <Trans id="labels.notagsavailable">No tags available</Trans>
            </div>
          )}
        </div>
      )}
    </div>
  )

  if (fider.isReadOnly) {
    return (
      <VStack>
        <HStack spacing={2} className="text-category">
          <Trans id="label.tags">Tags</Trans>
        </HStack>
        {tagsList}
      </VStack>
    )
  }

  return (
    <VStack>
      <HStack spacing={2} align="center" className="text-primary-base text-xs">
        {!isEditing && tagsList}
        {isEditing && editTagsList}
      </HStack>
    </VStack>
  )
}
