import React from "react"
import { PostStatus, Tag } from "@fider/models"
import { Checkbox, Dropdown, Icon } from "@fider/components"
import { HStack } from "@fider/components/layout"
import HeroIconFilter from "@fider/assets/images/heroicons-filter.svg"

import { useFider } from "@fider/hooks"
import { t } from "@lingui/macro"

interface OptionItem {
  value: string
  label: string
  count?: number
  tag?: Tag
}

interface PostFilterProps {
  activeFilter: FilterState
  countPerStatus: { [key: string]: number }
  filtersChanged: (filter: FilterState) => void
  tags: Tag[]
}
export interface FilterState {
  tags: string[]
  statuses: string[]
}

export const PostFilter = (props: PostFilterProps) => {
  const fider = useFider()

  const handleChangeFilter = (item: OptionItem) => () => {
    if (item.tag) {
      // Handle tag selection
      const newTags = props.activeFilter.tags.includes(item.value)
        ? props.activeFilter.tags.filter((t) => t !== item.value)
        : [...props.activeFilter.tags, item.value]
      const newFilters = { ...props.activeFilter, tags: newTags }
      props.filtersChanged(newFilters)
    } else {
      // Handle status selection
      const newStatuses = props.activeFilter.statuses.includes(item.value)
        ? props.activeFilter.statuses.filter((s) => s !== item.value)
        : [...props.activeFilter.statuses, item.value]
      const newFilters = { ...props.activeFilter, statuses: newStatuses }
      props.filtersChanged(newFilters)
    }
  }
  const options: OptionItem[] = []

  if (fider.session.isAuthenticated) {
    options.push({ value: "my-votes", label: t({ id: "home.postfilter.option.myvotes", message: "My Votes" }) })
  }

  PostStatus.All.filter((s) => s.filterable && props.countPerStatus[s.value]).forEach((s) => {
    const id = `enum.poststatus.${s.value.toString()}`
    options.push({
      label: t({ id, message: s.title }),
      value: s.value,
      count: props.countPerStatus[s.value],
    })
  })

  props.tags.forEach((t) => {
    options.push({
      label: t.name,
      value: t.slug,
      tag: t,
    })
  })

  const filterCount = props.activeFilter.tags.length + props.activeFilter.statuses.length

  return (
    <HStack className="mr-4">
      {/* <span className="text-category">
        <Trans id="home.postfilter.label.view">View</Trans>
      </span> */}
      <Dropdown
        renderHandle={
          <HStack className="h-10 text-medium text-xs rounded-md uppercase border border-gray-400 text-gray-800 p-2 px-3">
            <Icon sprite={HeroIconFilter} className="h-5 pr-1" />
            {t({ id: "home.filter.label", message: "Filter" })}
            {filterCount > 0 && <div className="bg-gray-200 inline-block rounded-full px-2 py-1 w-min-4 text-2xs text-center">{filterCount}</div>}
          </HStack>
        }
      >
        {options.map((o) => {
          const isChecked = o.tag ? props.activeFilter.tags.includes(o.value) : props.activeFilter.statuses.includes(o.value)
          return (
            <Dropdown.ListItem onClick={handleChangeFilter(o)} key={o.value}>
              <Checkbox field={o.value} checked={isChecked}>
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
