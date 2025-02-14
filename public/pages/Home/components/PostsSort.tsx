import React from "react"
import { Dropdown } from "@fider/components"
import { HStack } from "@fider/components/layout"
import { i18n } from "@lingui/core"

interface PostsSortProps {
  value: string
  onChange: (value: string) => void
}

export const PostsSort: React.FC<PostsSortProps> = ({ value = "trending", onChange }) => {
  const options = [
    { value: "trending", label: i18n._("home.postfilter.option.trending", { message: "Trending" }) },
    { value: "most-wanted", label: i18n._("home.postfilter.option.mostwanted", { message: "Most Wanted" }) },
    { value: "most-discussed", label: i18n._("home.postfilter.option.mostdiscussed", { message: "Most Discussed" }) },
    { value: "recent", label: i18n._("home.postfilter.option.recent", { message: "Recent" }) },
  ]

  const selectedItem = options.find((x) => x.value === value) || options[0]

  return (
    <HStack>
      <Dropdown
        renderHandle={
          <div className="h-10 flex flex-items-center text-medium text-xs rounded-md uppercase border border-gray-400 text-gray-800 p-2 px-3">
            {i18n._("home.postsort.label", { message: "Sort by:" })} {selectedItem.label}
          </div>
        }
      >
        {options.map((o) => (
          <Dropdown.ListItem key={o.value} onClick={() => onChange(o.value)}>
            <span className={value === o.value ? "text-semibold" : ""}>{o.label}</span>
          </Dropdown.ListItem>
        ))}
      </Dropdown>
    </HStack>
  )
}
