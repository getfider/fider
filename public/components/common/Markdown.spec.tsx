import React from "react"
import { render } from "@testing-library/react"
import { Markdown } from "./Markdown"
import { fiderAllowedSchemes } from "@fider/hooks"
fiderAllowedSchemes.get = () => ""

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

describe("<Markdown style='full' /> with markdown-escaped values", () => {
  // Notification titles embed user names and post titles escaped by markdown.Escape (Go),
  // e.g. "**" + Escape(name) + "**". They must render as the literal text.
  const cases = [
    {
      escaped: "\\[Your account needs verifying\\]\\(https\\:\\/\\/evil\\.example\\)",
      text: "[Your account needs verifying](https://evil.example)",
    },
    { escaped: "https\\:\\/\\/evil\\.example", text: "https://evil.example" },
    { escaped: "www\\.evil\\.example", text: "www.evil.example" },
    { escaped: "jon\\@evil\\.example", text: "jon@evil.example" },
    { escaped: "\\@\\[Jon Snow\\]", text: "@[Jon Snow]" },
    { escaped: "\\!\\[i\\]\\(https\\:\\/\\/evil\\.example\\/x\\.png\\)", text: "![i](https://evil.example/x.png)" },
    { escaped: "\\_Urgent\\_ \\*\\*bold\\*\\* \\~\\~del\\~\\~ \\`code\\`", text: "_Urgent_ **bold** ~~del~~ `code`" },
    { escaped: "<a href\\=\\'https\\:\\/\\/evil\\.example\\'\\>x<\\/a\\>", text: "<a href='https://evil.example'>x</a>" },
    { escaped: "<https\\:\\/\\/evil\\.example\\>", text: "<https://evil.example>" },
    { escaped: "Tom \\& Jerry <3 \\&lt\\;b\\&gt\\;", text: "Tom & Jerry <3 &lt;b&gt;" },
  ]

  cases.forEach(({ escaped, text }) => {
    test(`renders literally: ${text}`, () => {
      const { container } = render(<Markdown style="full" text={`**${escaped}** left a comment`} />)
      expect(container.querySelector("a, img, em, del, code, span.mention")).toBeNull()
      expect(container.querySelector("strong")?.textContent).toEqual(text)
    })
  })
})
