import { jwt } from "./jwt"
;[
  {
    expected: {
      sub: "1234567890",
      name: "John Snow",
      age: "30",
      iat: 1516239022,
    },
    token:
      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gU25vdyIsImFnZSI6IjMwIiwiaWF0IjoxNTE2MjM5MDIyfQ.A4e171Ry70APUJ2a9uo9G9Aju9G08AJB_Cr9B9ivX-o",
  },
  {
    expected: undefined,
    token: "wrong",
  },
].forEach((x) => {
  test(`decode('${x.token}') should be ${x.expected}`, () => {
    const result = jwt.decode(x.token)
    expect(result).toEqual(x.expected)
  })
})
