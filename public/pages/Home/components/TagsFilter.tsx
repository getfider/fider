import React, { useEffect, useRef, useState } from "react"
import { Tag } from "@fider/models"
import { ShowTag } from "@fider/components"
import { HStack, VStack } from "@fider/components/layout"

import { Trans } from "@lingui/react/macro"
import { i18n } from "@lingui/core"
import { sortTags } from "@fider/services"
import { useFider } from "@fider/hooks"

import "./TagsFilter.scss"

interface TagsFilterProps {
  tags: Tag[]
  selected: Tag[]
  selectionChanged: (selected: Tag[]) => void
}

export const TagsFilter = (props: TagsFilterProps) => {
  const fider = useFider()
  const canEdit = fider.session.isAuthenticated && fider.settings.postWithTags && props.tags.length > 0

  const [isEditing, setIsEditing] = useState(false)
  const [query, setQuery] = useState("")
  const [selected, setSelected] = useState(props.selected)

  const dropdownRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  const assignOrUnassignTag = async (tag: Tag) => {
    const idx = selected.indexOf(tag)
    const next = idx >= 0 ? selected.filter((x) => x != tag) : selected.concat(tag)
    setSelected(next)
    props.selectionChanged(next)
  }

  const onSubtitleClick = () => {
    setIsEditing(!isEditing)
    // Immediately focus on the input element when editing starts
    if (inputRef.current) {
      inputRef.current.focus()
    }
    setQuery("")
  }

  const handleOptionClick = (tag: Tag) => {
    assignOrUnassignTag(tag)
    // Keep focus on the input element after selection
    if (inputRef.current) {
      inputRef.current.focus()
    }
  }

  const filteredOptions = props.tags.filter(
    (option) => option.name.toLowerCase().includes(query.toLowerCase()) && !selected.some((tag) => tag.slug === option.slug)
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

  if (!canEdit && selected.length === 0) {
    return null
  }

  const tagsList = (
    <div className="tags-list">
      {selected.length > 0 && sortTags(selected).map((tag) => <ShowTag key={tag.id} tag={tag} />)}
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
        {sortTags(selected).map((tag) => (
          <div key={tag.id} className="selected-option">
            <ShowTag key={tag.id} tag={tag} />
            <button onClick={() => handleOptionClick(tag)} className="remove-button">
              x
            </button>
          </div>
        ))}
        {/* Search box to enter query string */}
        <input
          type="text"
          value={query}
          ref={inputRef}
          onChange={(e) => setQuery(e.target.value)}
          className="search-input"
          placeholder={selected.length ? "" : i18n._("label.selecttags", { message: "Select Tags..." })}
        />
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
