import * as markdown from "./markdown"
import { fiderAllowedSchemes } from "@fider/hooks"
fiderAllowedSchemes.get = () => "^monero:[48]\n^bitcoin:(1|3|bc1)"

const testCases = [
  {
    input: "Visit [GitHub](https://github.com) to learn more.",
    expectedFull: '<p>Visit <a class="text-link" href="https://github.com" rel="noopener nofollow" target="_blank">GitHub</a> to learn more.</p>',
    expectedPlainText: "Visit GitHub to learn more.",
  },
  {
    input: "My Picture ![](http://demo.dev.fider.io:3000/images/100/28)",
    expectedFull: '<p>My Picture <img src="http://demo.dev.fider.io:3000/images/100/28" alt=""></p>',
    expectedPlainText: "My Picture",
  },
  {
    input: "My Fider Picture ![](fider-image:attachments/zy0hBtqrjQki7M56p26AuAXljRoaNUSwZO6MOky5gnYm2nW1rsMmrp3dwhjGk7ok-aden.jpeg)",
    expectedFull:
      '<p>My Fider Picture <img src="/static/images/attachments/zy0hBtqrjQki7M56p26AuAXljRoaNUSwZO6MOky5gnYm2nW1rsMmrp3dwhjGk7ok-aden.jpeg" alt="" class="fider-inline-image" data-bkey="attachments/zy0hBtqrjQki7M56p26AuAXljRoaNUSwZO6MOky5gnYm2nW1rsMmrp3dwhjGk7ok-aden.jpeg"></p>',
    expectedPlainText: "My Fider Picture",
  },
  {
    input: "# Hello World",
    expectedFull: "<h1>Hello World</h1>",
    expectedSimple: "<p>Hello World</p>",
    expectedPlainText: "Hello World",
  },
  // {
  //   input: "Hello <b>Beautiful</b> World",
  //   expectedFull: "<p>Hello &lt;b&gt;Beautiful&lt;/b&gt; World</p>",
  //   expectedPlainText: "Hello &lt;b&gt;Beautiful&lt;/b&gt; World",
  // },
  {
    input: `[Uh oh...]("onerror="alert('XSS'))`,
    expectedFull: '<p><a class="text-link" href="" rel="noopener nofollow" target="_blank">Uh oh...</a></p>',
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
  {
    input:
      "[monero](monero:83zJ2jMbBoxJkhtpaRLk6fQVrvfmGbd8gYUL7FSdLxU91JSpiWXoLUtAMGqmfvfq3qRS5gJUvMY7oLFSx71wxhKRGG6ypMt) [bitcoin](bitcoin:1CgLs6CxXMAY4Pj4edQq5vyaFoP9NdqVKH) [litecoin](litecoin:ltc1qg0elpp0hxguwlsapl68gvklt5ngemj8k8lu0f5)",
    expectedFull:
      '<p><a class="text-link" href="monero:83zJ2jMbBoxJkhtpaRLk6fQVrvfmGbd8gYUL7FSdLxU91JSpiWXoLUtAMGqmfvfq3qRS5gJUvMY7oLFSx71wxhKRGG6ypMt" rel="noopener nofollow" target="_blank">monero</a> <a class="text-link" href="bitcoin:1CgLs6CxXMAY4Pj4edQq5vyaFoP9NdqVKH" rel="noopener nofollow" target="_blank">bitcoin</a> <a class="text-link" rel="noopener nofollow" target="_blank">litecoin</a></p>',
    expectedPlainText: "monero bitcoin litecoin",
  },
  {
    input: "Jane's & Jim's > [Matt](https://example.com)",
    expectedFull: '<p>Jane\'s &amp; Jim\'s &gt; <a class="text-link" href="https://example.com" rel="noopener nofollow" target="_blank">Matt</a></p>',
    expectedPlainText: "Jane's & Jim's > Matt",
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
