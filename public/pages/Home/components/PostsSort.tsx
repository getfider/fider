import React from "react"
import { Dropdown } from "@fider/components"
import { i18n } from "@lingui/core"
import IconSparkles from "@fider/assets/images/heroicons-sparkles-outline.svg"
import IconThumbsUp from "@fider/assets/images/heroicons-thumbsup.svg"
import IconChat from "@fider/assets/images/heroicons-chat-alt-2.svg"
import IconClock from "@fider/assets/images/heroicons-clock.svg"
import { HStack } from "@fider/components/layout"

interface PostsSortProps {
  value: string
  onChange: (value: string) => void
}

export const PostsSort: React.FC<PostsSortProps> = ({ value = "trending", onChange }) => {
  const options = [
    { value: "trending", label: i18n._({ id: "home.postfilter.option.trending", message: "Trending" }), icon: IconSparkles },
    { value: "most-wanted", label: i18n._({ id: "home.postfilter.option.mostwanted", message: "Most Wanted" }), icon: IconThumbsUp },
    { value: "most-discussed", label: i18n._({ id: "home.postfilter.option.mostdiscussed", message: "Most Discussed" }), icon: IconChat },
    { value: "recent", label: i18n._({ id: "home.postfilter.option.recent", message: "Recent" }), icon: IconClock },
  ]

  const selectedItem = options.find((x) => x.value === value) || options[0]

  return (
    <div>
      <Dropdown
        renderHandle={
          <HStack className="c-post-sort-btn">
            {i18n._({ id: "home.postsort.label", message: "Sort by:" })} {selectedItem.label}
          </HStack>
        }
      >
        {options.map((o) => (
          <Dropdown.ListItem key={o.value} onClick={() => onChange(o.value)} icon={o.icon}>
            <span className={value === o.value ? "text-semibold" : ""}>{o.label}</span>
          </Dropdown.ListItem>
        ))}
      </Dropdown>
    </div>
  )
}
