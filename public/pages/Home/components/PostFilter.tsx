import React, { useState } from "react"
import { PostStatus, Tag } from "@fider/models"
import { Checkbox, Dropdown, Icon } from "@fider/components"
import { HStack } from "@fider/components/layout"
import HeroIconFilter from "@fider/assets/images/heroicons-filter.svg"
import { useFider } from "@fider/hooks"
import { i18n } from "@lingui/core"
import { FilterState } from "./PostsContainer"

import "./PostFilter.scss"

type FilterType = "tag" | "status" | "myVotes"

interface OptionItem {
  value: string | boolean
  label: string
  count?: number
  type: FilterType
}

interface PostFilterProps {
  activeFilter: FilterState
  countPerStatus: { [key: string]: number }
  filtersChanged: (filter: FilterState) => void
  tags: Tag[]
}

export interface FilterItem {
  type: FilterType
  value: string | boolean
}

const FilterStateToFilterItems = (filterState: FilterState): FilterItem[] => {
  const filterItems: FilterItem[] = []
  filterState.statuses.forEach((s) => {
    filterItems.push({ type: "status", value: s })
  })
  filterState.tags.forEach((t) => {
    filterItems.push({ type: "tag", value: t })
  })
  if (filterState.myVotes) {
    filterItems.push({ type: "myVotes", value: true })
  }
  return filterItems
}

const FilterItemsToFilterState = (filterItems: FilterItem[]): FilterState => {
  const filterState: FilterState = { tags: [], statuses: [], myVotes: false }
  filterItems.forEach((i) => {
    if (i.type === "tag") {
      filterState.tags.push(i.value as string)
    } else if (i.type === "status") {
      filterState.statuses.push(i.value as string)
    } else if (i.type === "myVotes") {
      filterState.myVotes = true
    }
  })
  return filterState
}

export const PostFilter = (props: PostFilterProps) => {
  const fider = useFider()

  const filterItems: FilterItem[] = FilterStateToFilterItems(props.activeFilter)
  const [query, setQuery] = useState("")

  const handleChangeFilter = (item: OptionItem) => () => {
    const exists = filterItems.find((i) => i.type === item.type && i.value === item.value)
    const newFilter = exists
      ? filterItems.filter((i) => !(i.type === item.type && i.value === item.value))
      : [...filterItems, { type: item.type, value: item.value }]

    props.filtersChanged(FilterItemsToFilterState(newFilter))
    setQuery("")
  }
  const options: OptionItem[] = []

  if (fider.session.isAuthenticated) {
    options.push({ value: true, label: i18n._({ id: "home.postfilter.option.myvotes", message: "My Votes" }), type: "myVotes" })
  }

  PostStatus.All.filter((s) => s.filterable && props.countPerStatus[s.value]).forEach((s) => {
    const id = `enum.poststatus.${s.value.toString()}`
    options.push({
      label: i18n._(id, { message: s.title }),
      value: s.value,
      count: props.countPerStatus[s.value],
      type: "status",
    })
  })

  props.tags.forEach((t) => {
    options.push({
      label: t.name,
      value: t.slug,
      type: "tag",
    })
  })

  const filterCount = filterItems.length
  const filteredOptions = options.filter((option) => option.label.toLowerCase().includes(query.toLowerCase()))

  return (
    <HStack className="mr-4">
      <Dropdown
        onToggled={() => setQuery("")}
        renderHandle={
          <HStack className="h-10 text-medium text-xs rounded-md uppercase border border-gray-400 text-gray-800 p-2 px-3">
            <Icon sprite={HeroIconFilter} className="h-5 pr-1" />
            {i18n._({ id: "home.filter.label", message: "Filter" })}
            {filterCount > 0 && <div className="bg-gray-200 inline-block rounded-full px-2 py-1 w-min-4 text-2xs text-center">{filterCount}</div>}
          </HStack>
        }
      >
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="c-input filter-input"
          placeholder={i18n._({ id: "home.filter.search.label", message: "Search in filters..." })}
        />
        {filteredOptions.map((o) => {
          const isChecked = filterItems.find((f) => f.type === o.type && f.value === o.value) !== undefined

          return (
            <Dropdown.ListItem onClick={handleChangeFilter(o)} key={o.value.toString()}>
              <Checkbox field={o.value.toString()} checked={isChecked}>
                <HStack spacing={2}>
                  <span className={isChecked ? "text-semibold" : ""}>{o.label}</span>
                  <div className="">
                    {o.count && o.count > 0 && <span className="bg-gray-200 inline-block rounded-full px-1 w-min-4 text-2xs text-center">{o.count}</span>}
                  </div>
                </HStack>
              </Checkbox>
            </Dropdown.ListItem>
          )
        })}
      </Dropdown>
    </HStack>
  )
}
