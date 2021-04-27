import * as markdown from "./markdown"

const testCases = [
  {
    input: "Visit [GitHub](https://github.com) to learn more.",
    expectedFull: '<p>Visit <a target="_blank" rel="noopener" href="https://github.com" class="text-link">GitHub</a> to learn more.</p>',
    expectedPlainText: "Visit GitHub to learn more.",
  },
  {
    input: "My Picture ![](http://demo.dev.fider.io:3000/images/100/28)",
    expectedFull: "<p>My Picture </p>",
    expectedPlainText: "My Picture",
  },
  {
    input: "# Hello World",
    expectedFull: "<h1>Hello World</h1>",
    expectedSimple: "<p>Hello World</p>",
    expectedPlainText: "Hello World",
  },
  {
    input: "Hello <b>Beautiful</b> World",
    expectedFull: "<p>Hello &lt;b&gt;Beautiful&lt;/b&gt; World</p>",
    expectedPlainText: "Hello &lt;b&gt;Beautiful&lt;/b&gt; World",
  },
  {
    input: `[Uh oh...]("onerror="alert('XSS'))`,
    expectedFull: '<p><a target="_blank" rel="noopener" href="" class="text-link">Uh oh...</a></p>',
    expectedPlainText: "Uh oh...",
  },
  {
    input: "~~Option 3~~",
    expectedFull: "<p><del>Option 3</del></p>",
    expectedPlainText: "Option 3",
  },
  {
    input: "Check this out: `HEEEY`",
    expectedFull: "<p>Check this out: <code>HEEEY</code></p>",
    expectedPlainText: "Check this out: HEEEY",
  },
  {
    input: `# Hello World
How are you?`,
    expectedFull: `<h1>Hello World</h1>
<p>How are you?</p>`,
    expectedPlainText: "Hello World How are you?",
  },
  {
    input: `-123
-456
-789`,
    expectedFull: "<p>-123<br>-456<br>-789</p>",
    expectedPlainText: "-123 -456 -789",
  },
]

testCases.forEach((x) => {
  test(`Can parse markdown ${x.input} to ${x.expectedFull} (full mode)`, () => {
    const result = markdown.full(x.input)
    expect(result).toEqual(x.expectedFull)
  })
})

testCases.forEach((x) => {
  test(`Can parse markdown ${x.input} to ${x.expectedPlainText} (plain text)`, () => {
    const result = markdown.plainText(x.input)
    expect(result).toEqual(x.expectedPlainText)
  })
})
