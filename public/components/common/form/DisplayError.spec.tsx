import React from "react"
import { DisplayError } from "./DisplayError"
import { render } from "@testing-library/react"
import { Failure } from "@fider/services"

describe("<DisplayError />", () => {
  test("when error is undefined", () => {
    const { container } = render(<DisplayError />)
    expect(container.querySelector("div")).toBeNull()
  })

  test("when error is empty", () => {
    const { container } = render(<DisplayError error={{}} />)
    expect(container.querySelector("div")).toBeNull()
  })

  test("when error has only top level messages and fields is empty", () => {
    const error: Failure = {
      errors: [{ message: "Something went wrong." }],
    }

    const { container } = render(<DisplayError error={error} />)
    const root = container.querySelector("div")
    expect(root).toHaveClass("c-form-error")
    const items = root?.querySelectorAll("ul li")
    expect(items).toHaveLength(1)
    expect(items?.item(0)).toHaveTextContent("Something went wrong.")
  })

  test("when error has only top level messages and fields is given", () => {
    const error: Failure = {
      errors: [{ message: "Something went wrong." }],
    }

    const { container } = render(<DisplayError error={error} fields={["name"]} />)
    expect(container.querySelector("div")).toBeNull()
  })

  test("when error has both field and top level messages and fields are given", () => {
    const error: Failure = {
      errors: [
        { message: "Something went wrong." },
        { field: "name", message: "Name is required" },
        { field: "name", message: "Name must have between 0 and 10 chars" },
        { field: "age", message: "Age must be >= 18" },
      ],
    }

    const { container: container1 } = render(<DisplayError error={error} fields={["name"]} />)
    const root1 = container1.querySelector("div")
    expect(root1).toHaveClass("c-form-error")
    const items1 = root1?.querySelectorAll("ul li")
    expect(items1).toHaveLength(2)
    expect(items1?.item(0)).toHaveTextContent("Name is required")
    expect(items1?.item(1)).toHaveTextContent("Name must have between 0 and 10 chars")

    const { container: container2 } = render(<DisplayError error={error} fields={["age"]} />)
    const root2 = container2.querySelector("div")
    expect(root2).toHaveClass("c-form-error")
    const items2 = root2?.querySelectorAll("ul li")
    expect(items2).toHaveLength(1)
    expect(items2?.item(0)).toHaveTextContent("Age must be >= 18")

    const { container: container3 } = render(<DisplayError error={error} fields={["name", "age"]} />)
    const root3 = container3.querySelector("div")
    expect(root3).toHaveClass("c-form-error")
    const items3 = root3?.querySelectorAll("ul li")
    expect(items3).toHaveLength(3)
    expect(items3?.item(0)).toHaveTextContent("Name is required")
    expect(items3?.item(1)).toHaveTextContent("Name must have between 0 and 10 chars")
    expect(items3?.item(2)).toHaveTextContent("Age must be >= 18")
  })
})
