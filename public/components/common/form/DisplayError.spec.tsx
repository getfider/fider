import React from "react"
import { shallow } from "enzyme"
import { DisplayError } from "./DisplayError"
import { Failure } from "@fider/services"

describe("<DisplayError />", () => {
  test("when error is undefined", () => {
    const wrapper = shallow(<DisplayError />)
    expect(wrapper.getElement()).toBeNull()

    const wrapper1 = shallow(<DisplayError error={{}} />)
    expect(wrapper1.getElement()).toBeNull()
  })

  test("when error has only top level messages and fields is empty", () => {
    const error: Failure = {
      errors: [{ message: "Something went wrong." }],
    }

    const wrapper = shallow(<DisplayError error={error} />)
    const root = wrapper.find("div")
    expect(root.hasClass("c-form-error")).toBe(true)
    const items = root.find("ul li")
    expect(items).toHaveLength(1)
    expect(items.at(0).key()).toBe("Something went wrong.")
    expect(items.at(0).text()).toBe("Something went wrong.")
  })

  test("when error has only top level messages and fields is given", () => {
    const error: Failure = {
      errors: [{ message: "Something went wrong." }],
    }

    const wrapper = shallow(<DisplayError error={error} fields={["name"]} />)
    expect(wrapper.getElement()).toBeNull()
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

    const wrapper1 = shallow(<DisplayError error={error} fields={["name"]} />)
    const root1 = wrapper1.find("div")
    expect(root1.hasClass("c-form-error")).toBe(true)
    const items1 = root1.find("ul li")
    expect(items1).toHaveLength(2)
    expect(items1.at(0).key()).toBe("Name is required")
    expect(items1.at(0).text()).toBe("Name is required")
    expect(items1.at(1).key()).toBe("Name must have between 0 and 10 chars")
    expect(items1.at(1).text()).toBe("Name must have between 0 and 10 chars")

    const wrapper2 = shallow(<DisplayError error={error} fields={["age"]} />)
    const root2 = wrapper2.find("div")
    expect(root2.hasClass("c-form-error")).toBe(true)
    const items2 = root2.find("ul li")
    expect(items2).toHaveLength(1)
    expect(items2.at(0).key()).toBe("Age must be >= 18")
    expect(items2.at(0).text()).toBe("Age must be >= 18")

    const wrapper3 = shallow(<DisplayError error={error} fields={["name", "age"]} />)
    const root3 = wrapper3.find("div")
    expect(root3.hasClass("c-form-error")).toBe(true)
    const items3 = root3.find("ul li")
    expect(items3).toHaveLength(3)
    expect(items3.at(0).key()).toBe("Name is required")
    expect(items3.at(0).text()).toBe("Name is required")
    expect(items3.at(1).key()).toBe("Name must have between 0 and 10 chars")
    expect(items3.at(1).text()).toBe("Name must have between 0 and 10 chars")
    expect(items3.at(2).key()).toBe("Age must be >= 18")
    expect(items3.at(2).text()).toBe("Age must be >= 18")
  })
})
