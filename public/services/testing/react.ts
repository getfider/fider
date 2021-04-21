import { ShallowWrapper } from "enzyme"

const flushPromises = () => new Promise((resolve) => setImmediate(resolve))

export const rerender = async (component: ShallowWrapper<any, any, any>) => {
  await flushPromises()
  return component.update()
}
