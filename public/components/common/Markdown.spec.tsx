import React from "react"
import { render } from "@testing-library/react"
import { Markdown } from "./Markdown"

describe("<Markdown style='plainText' />", () => {
  const payloads = ["&lt;img src=x onerror=alert(1)&gt;", "&lt;script&gt;alert(1)&lt;/script&gt;", "`&lt;img src=x onerror=alert(1)&gt;`"]

  payloads.forEach((text) => {
    test(`renders entity-encoded markup as text: ${text}`, () => {
      const { container } = render(<Markdown style="plainText" maxLength={300} text={text} />)
      expect(container.querySelector("img, script")).toBeNull()
    })
  })

  test("shows special characters as typed", () => {
    const { container } = render(<Markdown style="plainText" text={"Jane's & Jim's > [Matt](https://example.com) <b>hi</b>"} />)
    expect(container.textContent).toEqual("Jane's & Jim's > Matt <b>hi</b>")
  })
})
