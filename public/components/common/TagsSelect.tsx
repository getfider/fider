import React, { KeyboardEvent as ReactKeyboardEvent, useEffect, useRef, useState } from "react"
import { Tag } from "@fider/models"
import { sortTags } from "@fider/services"
import { Button, ShowTag } from "@fider/components"
import { useFider } from "@fider/hooks"

import { HStack, VStack } from "@fider/components/layout"
import { Trans } from "@lingui/react/macro"
import { i18n } from "@lingui/core"

import "./TagsSelect.scss"

export interface TagsSelectProps {
  tags: Tag[]
  selected: Tag[]
  selectionChanged: (selected: Tag[]) => void
  canEdit: boolean
  asLinks?: boolean
}

export const TagsSelect = (props: TagsSelectProps) => {
  const fider = useFider()
  const [isEditing, setIsEditing] = useState(false)
  const [query, setQuery] = useState("")

  const dropdownRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  const assignOrUnassignTag = async (tag: Tag) => {
    const idx = props.selected.indexOf(tag)
    const next = idx >= 0 ? props.selected.filter((x) => x !== tag) : props.selected.concat(tag)
    props.selectionChanged(next)
  }

  const onSubtitleClick = () => {
    if (props.canEdit) {
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

  const handleOptionKey = (event: ReactKeyboardEvent, tag: Tag) => {
    if (event.code !== "Enter" && event.code !== "Space") {
      return
    }

    event.preventDefault()
    assignOrUnassignTag(tag)
    if (inputRef.current) {
      inputRef.current.focus()
    }
  }

  const filteredOptions = props.tags.filter(
    (option) => option.name.toLowerCase().includes(query.toLowerCase()) && !props.selected.some((tag) => tag.slug === option.slug)
  )

  const handleEsc = (event: KeyboardEvent | ReactKeyboardEvent) => {
    if (event.code !== "Escape") {
      return
    }

    setIsEditing(false)
  }

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsEditing(false)
      }
    }

    document.addEventListener("mousedown", handleClickOutside)
    document.addEventListener("keydown", handleEsc)
    return () => {
      document.removeEventListener("mousedown", handleClickOutside)
      document.addEventListener("keydown", handleEsc)
    }
  }, [])

  useEffect(() => {
    if (isEditing && inputRef.current) {
      inputRef.current.focus()
    }
  }, [isEditing])

  if (!props.canEdit && props.selected.length === 0) {
    return null
  }

  const tagsList = (
    <div className="tags-list">
      {props.selected.length > 0 && sortTags(props.selected).map((tag) => <ShowTag key={tag.id} tag={tag} link={props.asLinks} />)}
      {props.canEdit && (
        <Button variant={"link"} size={"no-padding"} onClick={onSubtitleClick}>
          <Trans id="label.edittags">Edit tags</Trans>
        </Button>
      )}
    </div>
  )

  // Dynamic multiselect dropdown for tags selection
  const editTagsList = props.tags.length > 0 && (
    <div className="dropdown-wrapper" ref={dropdownRef}>
      {/* Selected options and search input */}
      <div className="selected-options-container">
        {props.selected.length === 0 && (
          <span className="text-muted">
            <Trans id="labels.notagsselected">No tags selected</Trans>
          </span>
        )}
        {sortTags(props.selected).map((tag) => (
          <div key={tag.id} className="selected-option">
            <ShowTag tag={tag} />
            <button onClick={() => handleOptionClick(tag)} className="remove-button">
              x
            </button>
          </div>
        ))}
      </div>

      {/* Dropdown options after items are filtered */}
      {isEditing && (
        <div className="options-container">
          {/* Search box to enter query string */}
          <input
            type="text"
            value={query}
            ref={inputRef}
            onChange={(e) => setQuery(e.target.value)}
            className="c-input search-input"
            placeholder={i18n._("label.selecttags", { message: "Select Tags..." })}
            onKeyDown={() => handleEsc}
          />
          {filteredOptions.length > 0 ? (
            sortTags(filteredOptions).map((tag) => (
              <div key={tag.id} className="option-item" onClick={() => handleOptionClick(tag)} onKeyDown={(event) => handleOptionKey(event, tag)}>
                <ShowTag tag={tag} />
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
