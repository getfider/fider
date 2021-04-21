import { classSet, formatDate, timeSince, fileToBase64 } from "./utils"
import { readFileSync } from "fs"
;[
  { input: null, expected: "" },
  { input: undefined, expected: "" },
  { input: {}, expected: "" },
  { input: { green: true }, expected: "green" },
  { input: { disabled: false, green: false }, expected: "" },
  { input: { disabled: true, green: true }, expected: "disabled green" },
  { input: { "ui button": true, green: true }, expected: "ui button green" },
  { input: { "ui button": true }, expected: "ui button" },
  { input: { "": true }, expected: "" },
  {
    input: {
      green: () => {
        throw Error()
      },
    },
    expected: "green",
  },
].forEach((x) => {
  test(`classSet of ${JSON.stringify(x.input)} should be ${x.expected}`, () => {
    const className = classSet(x.input)
    expect(className).toEqual(x.expected)
  })
})
;[
  { input: new Date(2018, 4, 27, 10, 12, 59), expected: "May 27, 2018 · 10:12" },
  { input: new Date(2058, 12, 12, 23, 21, 53), expected: "January 12, 2059 · 23:21" },
  { input: "2018-04-11T18:13:33.128082", expected: "April 11, 2018 · 18:13" },
  { input: "2017-11-20T07:47:42.158142", expected: "November 20, 2017 · 07:47" },
].forEach((x) => {
  test(`formatDate (full) of ${x.input} should be ${x.expected}`, () => {
    const result = formatDate(x.input, "full")
    expect(result).toEqual(x.expected)
  })
})
;[
  { input: new Date(2018, 4, 27, 10, 12, 59), expected: "May '18" },
  { input: new Date(2058, 12, 12, 23, 21, 53), expected: "Jan '59" },
  { input: "2018-04-11T18:13:33.128082", expected: "Apr '18" },
  { input: "2017-11-20T07:47:42.158142", expected: "Nov '17" },
].forEach((x) => {
  test(`formatDate (short) of ${x.input} should be ${x.expected}`, () => {
    const result = formatDate(x.input, "short")
    expect(result).toEqual(x.expected)
  })
})
;[
  { input: new Date(2018, 4, 27, 19, 51, 9), expected: "less than a minute ago" },
  { input: new Date(2018, 4, 27, 10, 12, 59), expected: "about 10 hours ago" },
  { input: new Date(2018, 4, 26, 10, 12, 59), expected: "a day ago" },
  { input: new Date(2018, 4, 22, 10, 12, 59), expected: "5 days ago" },
  { input: new Date(2018, 3, 22, 10, 12, 59), expected: "about a month ago" },
  { input: new Date(2018, 2, 22, 10, 12, 59), expected: "2 months ago" },
  { input: new Date(2017, 3, 22, 10, 12, 59), expected: "about a year ago" },
  { input: new Date(2013, 3, 22, 10, 12, 59), expected: "5 years ago" },
].forEach((x) => {
  test(`timeSince ${x.input} should be ${x.expected}`, () => {
    const now = new Date(2018, 4, 27, 19, 51, 10)
    const result = timeSince(now, x.input)
    expect(result).toEqual(x.expected)
  })
})
;[
  { input: new Date(2018, 4, 27, 19, 51, 9), expected: "less than a minute ago" },
  { input: new Date(2018, 4, 27, 10, 12, 59), expected: "about 10 hours ago" },
  { input: new Date(2018, 4, 26, 10, 12, 59), expected: "a day ago" },
  { input: new Date(2018, 4, 22, 10, 12, 59), expected: "5 days ago" },
  { input: new Date(2018, 3, 22, 10, 12, 59), expected: "about a month ago" },
  { input: new Date(2018, 2, 22, 10, 12, 59), expected: "2 months ago" },
  { input: new Date(2017, 3, 22, 10, 12, 59), expected: "about a year ago" },
  { input: new Date(2013, 3, 22, 10, 12, 59), expected: "5 years ago" },
].forEach((x) => {
  test(`timeSince ${x.input} should be ${x.expected}`, () => {
    const now = new Date(2018, 4, 27, 19, 51, 10)
    const result = timeSince(now, x.input)
    expect(result).toEqual(x.expected)
  })
})

test("Can convert file to base64", async () => {
  const content = readFileSync("./favicon.png")
  const favicon = new File([content], "favicon.png")
  const base64 = await fileToBase64(favicon)
  expect(base64).toBe(Buffer.from(content).toString("base64"))
})
