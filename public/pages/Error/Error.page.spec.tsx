import React from "react"
import { render } from "@testing-library/react"
import { ErrorPage } from "./Error.page"
import { FiderContext } from "@fider/services"
import { fiderMock } from "@fider/services/testing"

describe("<ErrorPage />", () => {
  const createFakeErrorInfo = () => ({ componentStack: "" } as React.ErrorInfo)

  test("it should show the error when showError returns true", () => {
    const error = new Error("Hello")
    const errorInfo = createFakeErrorInfo()
    const { container } = render(
      <FiderContext.Provider value={fiderMock.notAuthenticated()}>
        <ErrorPage error={error} errorInfo={errorInfo} showDetails={true} />
      </FiderContext.Provider>
    )

    const pre = container.querySelector("pre")
    expect(pre).not.toBeNull()
    expect(pre).toHaveTextContent("Error: Hello")
  })

  test("it should not show the error when showError returns false", () => {
    const error = new Error("Hello")
    const errorInfo = createFakeErrorInfo()

    const { container } = render(
      <FiderContext.Provider value={fiderMock.notAuthenticated()}>
        <ErrorPage error={error} errorInfo={errorInfo} showDetails={false} />
      </FiderContext.Provider>
    )

    expect(container.querySelector("pre")).toBeNull()
  })
})
