import React from "react"
import { shallow } from "enzyme"
import { ErrorPage } from "./Error.page"

describe("<ErrorPage />", () => {
  const createFakeErrorInfo = () => ({ componentStack: "" } as React.ErrorInfo)

  test("it should show the error when showError returns true", () => {
    const error = new Error("Hello")
    const errorInfo = createFakeErrorInfo()

    const wrapper = shallow(<ErrorPage error={error} errorInfo={errorInfo} showDetails={true} />)
    expect(wrapper.find("pre")).toHaveLength(1)
  })

  test("it should not show the error when showError returns false", () => {
    const error = new Error("Hello")
    const errorInfo = createFakeErrorInfo()

    const wrapper = shallow(<ErrorPage error={error} errorInfo={errorInfo} showDetails={false} />)

    expect(wrapper.find("pre")).toHaveLength(0)
  })
})
