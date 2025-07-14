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
  // If true, you always see the tags edit box, rather than having to put the tags list into "edit mode"
  alwaysEditing?: boolean
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

  const viewModeTagsList = (
    <div className="c-tags-select__container">
      <div className="c-tags-select__list">
        {props.selected.length > 0 && sortTags(props.selected).map((tag) => <ShowTag key={tag.id} tag={tag} link={props.asLinks} />)}
        {props.canEdit && (
          <div>
            <Button variant={"link"} size={"no-padding"} onClick={onSubtitleClick}>
              {props.selected.length ? <Trans id="label.edittags">Edit tags</Trans> : <Trans id="label.addtags">Add tags...</Trans>}
            </Button>
          </div>
        )}
      </div>
    </div>
  )

  // Dynamic multiselect dropdown for tags selection
  const editTagsList = props.tags.length > 0 && (
    <div className="c-tags-select__container" ref={dropdownRef} onClick={props.canEdit && !isEditing ? onSubtitleClick : undefined}>
      <div className="c-tags-select__selected-container">
        {props.selected.length === 0 && props.canEdit && (
          <Button className="text-gray-600" variant={"link"} size={"no-padding"} onClick={onSubtitleClick}>
            <Trans id="label.addtags">Add tags...</Trans>
          </Button>
        )}
        {sortTags(props.selected).map((tag) => (
          <div key={tag.id} className="c-tags-select__selected-item">
            <ShowTag tag={tag} />
            <button onClick={() => handleOptionClick(tag)} className="c-tags-select__remove-button">
              x
            </button>
          </div>
        ))}
      </div>

      {/* Dropdown options after items are filtered */}
      {isEditing && (
        <div className="c-tags-select__options">
          {/* Search box to enter query string */}
          <input
            type="text"
            value={query}
            ref={inputRef}
            onChange={(e) => setQuery(e.target.value)}
            className="c-input c-tags-select__search-input"
            placeholder={i18n._({ id: "label.searchtags", message: "Search tags..." })}
            onKeyDown={() => handleEsc}
          />
          {filteredOptions.length > 0 ? (
            sortTags(filteredOptions).map((tag) => (
              <div key={tag.id} className="c-tags-select__option" onClick={() => handleOptionClick(tag)} onKeyDown={(event) => handleOptionKey(event, tag)}>
                <ShowTag tag={tag} />
              </div>
            ))
          ) : (
            <div className="c-tags-select__no-options">
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
        {viewModeTagsList}
      </VStack>
    )
  }

  return (
    <VStack className="c-tags-select">
      <HStack spacing={2} align="center" className="text-primary-base text-xs">
        {isEditing || props.alwaysEditing ? editTagsList : viewModeTagsList}
      </HStack>
    </VStack>
  )
}
