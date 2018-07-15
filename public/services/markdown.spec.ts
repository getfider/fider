import * as markdown from "./markdown";

const testCases = [
  {
    input: "Visit [GitHub](https://github.com) to learn more.",
    expectedFull: '<p>Visit <a href="https://github.com" target="_blank">GitHub</a> to learn more.</p>',
    expectedSimple: '<p>Visit <a href="https://github.com" target="_blank">GitHub</a> to learn more.</p>'
  },
  {
    input: "My Picture ![](http://demo.dev.fider.io:3000/images/100/28)",
    expectedFull: '<p>My Picture <img src="http://demo.dev.fider.io:3000/images/100/28" alt="" /></p>',
    expectedSimple: '<p>My Picture !<a href="http://demo.dev.fider.io:3000/images/100/28" target="_blank"></a></p>'
  }
];

testCases.forEach(x => {
  test(`Can parse markdown ${x.input} to ${x.expectedFull} (full mode)`, () => {
    const result = markdown.full(x.input);
    expect(result).toEqual(x.expectedFull);
  });
});

testCases.forEach(x => {
  test(`Can parse markdown ${x.input} to ${x.expectedSimple} (simple mode)`, () => {
    const result = markdown.simple(x.input);
    expect(result).toEqual(x.expectedSimple);
  });
});
