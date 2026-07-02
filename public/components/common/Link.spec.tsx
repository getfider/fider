import React from "react"
import { render } from "@testing-library/react"
import { Link } from "./Link"
import { Fider } from "@fider/services"

beforeAll(() => {
  Fider.initialize({ settings: { baseURL: "" }, props: {} })
})

describe("<Link />", () => {
  test("renders bare href when hosted at the domain root", () => {
    Fider.settings.baseURL = "https://example.com"
    const { container } = render(<Link href="/posts/1">One</Link>)
    expect(container.querySelector("a")).toHaveAttribute("href", "/posts/1")
  })

  test("prepends basePath when hosted under a sub-path", () => {
    Fider.settings.baseURL = "https://example.com/feedback"
    const { container } = render(<Link href="/posts/1">One</Link>)
    expect(container.querySelector("a")).toHaveAttribute("href", "/feedback/posts/1")
  })

  test("passes additional props through to the underlying anchor", () => {
    Fider.settings.baseURL = "https://example.com"
    const { container } = render(
      <Link href="/x" className="my-class" target="_blank" rel="noopener">
        link
      </Link>
    )
    const a = container.querySelector("a")
    expect(a).toHaveClass("my-class")
    expect(a).toHaveAttribute("target", "_blank")
    expect(a).toHaveAttribute("rel", "noopener")
  })
})
