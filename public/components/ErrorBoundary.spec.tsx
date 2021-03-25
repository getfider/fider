import React from "react"
import { shallow } from "enzyme"
import { ErrorBoundary } from "@fider/components"

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
    shallow(
      <ErrorBoundary onError={errorSpy}>
        <div id="no-error">No Error!</div>
      </ErrorBoundary>
    )

    expect(errorSpy).not.toHaveBeenCalled()
  })

  describe("when error caught", () => {
    test("error should be passed to onError", () => {
      const error = new Error("Whoops!")
      const errorSpy = jest.fn()
      const wrapper = shallow(<ErrorBoundary onError={errorSpy} />)

      const componentDidCatch = wrapper.instance().componentDidCatch
      if (componentDidCatch) {
        componentDidCatch.bind(wrapper.instance())(error, {} as React.ErrorInfo)
      }

      expect(errorSpy).toHaveBeenCalledWith(error)
    })
  })
})
