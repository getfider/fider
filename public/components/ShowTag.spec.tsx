import React from "react"
import { fireEvent, render, screen } from "@testing-library/react"
import { ShowTag } from "./ShowTag"

describe("<ShowTag />", () => {
  test("prevents navigation and preserves parent clicks without link", () => {
    const onClick = jest.fn()
    render(
      <div onClick={onClick}>
        <ShowTag tag={{ id: 1, name: "Bug", slug: "bug", color: "ff0000", isPublic: true }} />
      </div>
    )

    expect(fireEvent.click(screen.getByTitle("Bug"))).toBe(false)
    expect(onClick).toHaveBeenCalledTimes(1)
  })
})
