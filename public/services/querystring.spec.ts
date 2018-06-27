import * as qs from "./querystring";

jest
  .spyOn(window.location, "href")
  .mockImplementation(() => "abc")

[
  { qs: "?name=John&age=30", name: "name", expected: "John" }, 
  { qs: "", name: "name", expected: "" }
].forEach(x => {
  test.skip(`get('${x.name}') should be ${x.expected}`, () => {
    const result = qs.get(x.name);
    expect(result).toEqual(x.expected);
  });
});
