import React, { useState } from "react"
import { Tag } from "@fider/models"
import { Dropdown, Icon, ShowTag } from "@fider/components"
import IconCheck from "@fider/assets/images/heroicons-check.svg"
import { HStack } from "@fider/components/layout"

import { plural } from "@lingui/core/macro"
import { Trans } from "@lingui/react/macro"
import { i18n } from "@lingui/core"
import ShieldCheck from "@fider/assets/images/heroicons-shieldcheck.svg"

interface TagsFilterProps {
  tags: Tag[]
  selected: string[]
  selectionChanged: (selected: string[]) => void
}

export const TagsFilter = (props: TagsFilterProps) => {
  const [selected, setSelected] = useState<string[]>(props.selected)

  const toggle = (item: Tag) => () => {
    const idx = selected.indexOf(item.slug)
    const next = idx >= 0 ? selected.filter((x) => x != item.slug) : selected.concat(item.slug)
    setSelected(next)
    props.selectionChanged(next)
  }

  const count = selected.length
  const label =
    count >= 1
      ? plural(count, {
          one: `# tag`,
          other: `# tags`,
        })
      : i18n._("home.tagsfilter.selected.none", { message: "Any tag" })

  return (
    <HStack>
      <span>
        <Trans id="home.tagsfilter.label.with">with</Trans>
      </span>
      <Dropdown renderHandle={<div className="text-medium text-xs uppercase rounded-md uppercase bg-gray-100 px-2 py-1">{label}</div>}>
        {props.tags.map((t) => (
          <Dropdown.ListItem onClick={toggle(t)} key={t.id}>
            <HStack spacing={1}>
              <ShowTag tag={t} circular noBackground />
              <span className={"flex flex-items-center"}>
                {!t.isPublic && <Icon height="14" width="14" sprite={ShieldCheck} className={"mr-1"} />}
                {t.name}
              </span>
              <Icon sprite={IconCheck} className={`h-4 text-green-600 ${!selected.includes(t.slug) && "invisible"}`} />
            </HStack>
          </Dropdown.ListItem>
        ))}
      </Dropdown>
    </HStack>
  )
}
