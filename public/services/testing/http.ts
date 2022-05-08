import { http } from "@fider/services"

const createOkMock = () => {
  return jest.fn(() => {
    return Promise.resolve({
      ok: true,
      data: null as any,
    })
  })
}

export const httpMock = {
  alwaysOk: () => {
    http.get = createOkMock()
    http.post = createOkMock()
    http.put = createOkMock()
    http.delete = createOkMock()
    return http
  },
}
