import React from "react"
import { render, screen, fireEvent } from "@testing-library/react"
import { FiderContext } from "@fider/services"
import { fiderMock } from "@fider/services/testing"
import { PortalDirectoryPage, PortalSummary } from "./PortalDirectory.page"

const avengers: PortalSummary = {
  name: "Avengers",
  url: "http://avengers.test.fider.io",
  host: "avengers.test.fider.io",
  logoURL: "http://avengers.test.fider.io/static/images/logos/avengers.png?size=200",
}

const demo: PortalSummary = {
  name: "Demonstration",
  url: "http://demo.test.fider.io",
  host: "demo.test.fider.io",
}

const renderPage = (portals: PortalSummary[]) =>
  render(
    <FiderContext.Provider value={fiderMock.notAuthenticated()}>
      <PortalDirectoryPage portals={portals} />
    </FiderContext.Provider>
  )

describe("<PortalDirectoryPage />", () => {
  test("renders a link to each portal", () => {
    renderPage([avengers, demo])

    const links = screen.getAllByRole("link")
    expect(links).toHaveLength(2)
    expect(links[0]).toHaveAttribute("href", "http://avengers.test.fider.io")
    expect(links[1]).toHaveAttribute("href", "http://demo.test.fider.io")
    expect(screen.getByText("Avengers")).toBeInTheDocument()
    expect(screen.getByText("avengers.test.fider.io")).toBeInTheDocument()
  })

  test("renders the logo when the portal has one", () => {
    const { container } = renderPage([avengers])

    const logo = container.querySelector("img")
    expect(logo).toHaveAttribute("src", avengers.logoURL)
  })

  test("renders an initial placeholder when the portal has no logo", () => {
    const { container } = renderPage([demo])

    expect(container.querySelector("img")).toBeNull()
    expect(container.querySelector(".c-portal-card__initial")).toHaveTextContent("D")
  })

  test("filters portals case-insensitively by name", () => {
    renderPage([avengers, demo])

    fireEvent.change(screen.getByRole("textbox"), { target: { value: "aVeN" } })

    const links = screen.getAllByRole("link")
    expect(links).toHaveLength(1)
    expect(links[0]).toHaveAttribute("href", "http://avengers.test.fider.io")
  })

  test("filters portals by host", () => {
    renderPage([avengers, demo])

    fireEvent.change(screen.getByRole("textbox"), { target: { value: "demo.test" } })

    const links = screen.getAllByRole("link")
    expect(links).toHaveLength(1)
    expect(links[0]).toHaveAttribute("href", "http://demo.test.fider.io")
  })

  // The two empty states are asserted by element rather than by copy: the lingui macro
  // hoists <Trans> children into a message prop, so translated text is not rendered under
  // the jest mock for @lingui/react.
  test("shows the no-match state when nothing matches the filter", () => {
    const { container } = renderPage([avengers, demo])

    fireEvent.change(screen.getByRole("textbox"), { target: { value: "nothing here" } })

    expect(screen.queryAllByRole("link")).toHaveLength(0)
    expect(container.querySelector(".c-portal-directory__nomatch")).not.toBeNull()
    expect(container.querySelector(".c-portal-directory__empty")).toBeNull()
  })

  test("shows the empty state, and no filter box, when there are no portals at all", () => {
    const { container } = renderPage([])

    expect(screen.queryAllByRole("link")).toHaveLength(0)
    expect(screen.queryByRole("textbox")).toBeNull()
    expect(container.querySelector(".c-portal-directory__empty")).not.toBeNull()
    expect(container.querySelector(".c-portal-directory__nomatch")).toBeNull()
  })
})
