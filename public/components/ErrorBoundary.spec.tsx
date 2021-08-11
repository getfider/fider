import React from "react"
import { render } from "@testing-library/react"
import { ErrorBoundary } from "@fider/components"
import { FiderContext } from "@fider/services"
import { fiderMock } from "@fider/services/testing"

describe("<ErrorBoundary />", () => {
  let errorMethod: () => void

  // Stub out console.error to hide noisy Virtual DOM exceptions.
  beforeAll(() => {
    errorMethod = console.error
    console.error = () => null
  })

  afterAll(() => {
    console.error = errorMethod
  })

  test("when no error caught", () => {
    const errorSpy = jest.fn()
    render(
      <ErrorBoundary onError={errorSpy}>
        <div id="no-error">No Error!</div>
      </ErrorBoundary>
    )

    expect(errorSpy).not.toHaveBeenCalled()
  })

  describe("when error caught", () => {
    test("error should be passed to onError", () => {
      const errorSpy = jest.fn()
      const error = new Error("Whoops!")
      const ThrowError = () => {
        throw error
      }

      render(
        <FiderContext.Provider value={fiderMock.notAuthenticated()}>
          <ErrorBoundary onError={errorSpy}>
            <ThrowError />
          </ErrorBoundary>
        </FiderContext.Provider>
      )

      expect(errorSpy).toHaveBeenCalledWith(error)
    })
  })
})
