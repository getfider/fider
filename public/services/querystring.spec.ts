import * as qs from "./querystring"
import navigator from "./navigator"
;[
  { qs: "?name=John&age=30", name: "name", expected: "John" },
  { qs: "?name=John&age=30", name: "age", expected: "30" },
  { qs: "?name=", name: "name", expected: "" },
  { qs: "", name: "name", expected: "" },
].forEach((x) => {
  test(`get('${x.name}') should be ${x.expected}`, () => {
    navigator.url = () => `http://example.com${x.qs}`
    const result = qs.get(x.name)
    expect(result).toEqual(x.expected)
  })
})
;[
  { qs: "?name=John,Arya&age=30", name: "name", expected: ["John", "Arya"] },
  { qs: "?name=John&age=30", name: "age", expected: ["30"] },
  { qs: "?name=", name: "name", expected: [] },
  { qs: "", name: "name", expected: [] },
].forEach((x) => {
  test(`getArray('${x.name}') should be ${x.expected}`, () => {
    navigator.url = () => `http://example.com${x.qs}`
    const result = qs.getArray(x.name)
    expect(result).toEqual(x.expected)
  })
})
;[
  { qs: "?name=John&age=30", name: "age", expected: 30 },
  { qs: "?name=", name: "name", expected: undefined },
  { qs: "", name: "name", expected: undefined },
].forEach((x) => {
  test(`getNumber('${x.name}') should be ${x.expected}`, () => {
    navigator.url = () => `http://example.com${x.qs}`
    const result = qs.getNumber(x.name)
    expect(result).toEqual(x.expected)
  })
})
;[
  { qs: "?name=John&age=30", name: "age", value: 60, expected: "http://example.com?name=John&age=60" },
  { qs: "", name: "age", value: 50, expected: "http://example.com?age=50" },
  { qs: "?age=2", name: "age", value: 21, expected: "http://example.com?age=21" },
].forEach((x) => {
  test(`set(${x.name}, ${x.value}) should return ${x.expected}`, () => {
    navigator.url = () => `http://example.com${x.qs}`
    const newQueryString = qs.set(x.name, x.value)
    expect(newQueryString).toEqual(x.expected)
  })
})
;[
  { object: {}, expected: "" },
  { object: undefined, expected: "" },
  {
    object: {
      name: "John",
    },
    expected: "?name=John",
  },
  {
    object: {
      name: "John",
      age: 30,
    },
    expected: "?name=John&age=30",
  },
  {
    object: {
      tags: ["bug", "feature"],
      age: 30,
    },
    expected: "?tags=bug,feature&age=30",
  },
].forEach((x) => {
  test(`stringfy('${x.object}') should be ${x.expected}`, () => {
    const result = qs.stringify(x.object)
    expect(result).toEqual(x.expected)
  })
})
