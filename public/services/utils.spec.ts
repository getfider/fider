import { classSet, formatDate, timeSince, fileToBase64 } from "./utils";
import { readFileSync } from "fs";

[
  { input: null, expected: "" },
  { input: undefined, expected: "" },
  { input: {}, expected: "" },
  { input: { green: true }, expected: "green" },
  { input: { disabled: false, green: false }, expected: "" },
  { input: { disabled: true, green: true }, expected: "disabled green" }
].forEach(x => {
  test(`classSet of ${JSON.stringify(x.input)} should be ${x.expected}`, () => {
    const className = classSet(x.input);
    expect(className).toEqual(x.expected);
  });
});

[
  { input: new Date(2018, 4, 27, 10, 12, 59), expected: "May 27, 2018 · 10:12" },
  { input: new Date(2058, 12, 12, 23, 21, 53), expected: "January 12, 2059 · 23:21" },
  { input: "2018-04-11T18:13:33.128082Z", expected: "April 11, 2018 · 18:13" },
  { input: "2017-11-20T07:47:42.158142Z", expected: "November 20, 2017 · 07:47" }
].forEach(x => {
  test(`formatDate of ${x.input} should be ${x.expected}`, () => {
    const result = formatDate(x.input);
    expect(result).toEqual(x.expected);
  });
});

[
  { input: new Date(2018, 4, 27, 19, 51, 9), expected: "less than a minute ago" },
  { input: new Date(2018, 4, 27, 10, 12, 59), expected: "about 10 hours ago" },
  { input: new Date(2018, 4, 26, 10, 12, 59), expected: "a day ago" },
  { input: new Date(2018, 4, 22, 10, 12, 59), expected: "5 days ago" },
  { input: new Date(2018, 3, 22, 10, 12, 59), expected: "about a month ago" },
  { input: new Date(2018, 2, 22, 10, 12, 59), expected: "2 months ago" },
  { input: new Date(2017, 3, 22, 10, 12, 59), expected: "about a year ago" },
  { input: new Date(2013, 3, 22, 10, 12, 59), expected: "5 years ago" }
].forEach(x => {
  test(`timeSince ${x.input} should be ${x.expected}`, () => {
    const now = new Date(2018, 4, 27, 19, 51, 10);
    const result = timeSince(now, x.input);
    expect(result).toEqual(x.expected);
  });
});

[
  { input: new Date(2018, 4, 27, 19, 51, 9), expected: "less than a minute ago" },
  { input: new Date(2018, 4, 27, 10, 12, 59), expected: "about 10 hours ago" },
  { input: new Date(2018, 4, 26, 10, 12, 59), expected: "a day ago" },
  { input: new Date(2018, 4, 22, 10, 12, 59), expected: "5 days ago" },
  { input: new Date(2018, 3, 22, 10, 12, 59), expected: "about a month ago" },
  { input: new Date(2018, 2, 22, 10, 12, 59), expected: "2 months ago" },
  { input: new Date(2017, 3, 22, 10, 12, 59), expected: "about a year ago" },
  { input: new Date(2013, 3, 22, 10, 12, 59), expected: "5 years ago" }
].forEach(x => {
  test(`timeSince ${x.input} should be ${x.expected}`, () => {
    const now = new Date(2018, 4, 27, 19, 51, 10);
    const result = timeSince(now, x.input);
    expect(result).toEqual(x.expected);
  });
});

test("Can convert file to base64", async () => {
  const content = readFileSync("./favicon.ico");
  const favicon = new File([content], "favicon.ico");
  const base64 = await fileToBase64(favicon);
  expect(base64).toBe(
    "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//8BZ///MYT///z9//99fX12AAAAAAAAAAAAAAAAAAAAAHd3dwCZmZkAb29vZv7+//9Sj/v/Nkbp/zdF6f83Ren/N0Xp/zdF6f83Ren/oqv0//r9/v/6/f7/+v3+//r9/v/6/f7/+v3+//r9/v/6/P//+vz///r8///6/P//+vz///r8///6/P///v7//3x8fHYAAAAAAAAAAAAAAAAAAAAAd3d3ALGxsRbw8PDsstH//wFn//8eVPP/NkPq/zZD6v82Q+r/NkPq/1Rf7f/6+v7/xMTEupmZmYyZmZmMmZmZjJmZmYyZmZmMmZmZjImJiYCIiIiAiIiIgIiIiICIiIiAiIiIgImJiYBycnJsBgYGBgAAAAAAAAAAAAAAAAAAAAB3d3cA/f39mPj6//8mff//AWf//wRl/f8xR+z/NkPq/zZD6v82Q+r/xcn5/+zs7OgTExMQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAHd3dwD////8h7f//wFn//8BZ///AWf//xVa9/82Q+r/NkPq/3J78P/+/v7/cnJyagAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAd3d3AOTv//8Nbv//AWf//wFn//8BZ///AWb+/ypL7v87SOr/4uT8//T09PBqampiZmZmYGZmZmBmZmZgZmZmYGZmZmBmZmZgZmZmYGZmZmBmZmZgZmZmYGZmZmBmZmZgZmZmYGZmZmBmZmZgZmZmYGZmZmBaWlpUCwsLCgAAAAB3d3cAXZ7//wFn//8BZ///AWf//wFn//8BZ///DGD6/5Sc9P/////////////////////////////////////////////////////////////////////////////////////////////////////////////////AwMC4AQEBAnd3dwBCpf7/A4v+/wOL/v8Di/7/A4v+/wOL/v8EkP3/GeL9/xz9//8c/f//HP3//xz9//8c/f//K/rx/+uwS//1klX/9ZJV//WSVf/1klX/9ZJV/++SU/9opDf/ZbFJ/2WxSf9lsUn/ZbFJ/2WxSf9ksUr/dN61//////9aWlpSd3d3AM/o/v8Cj/7/AI///wCP//8Aj///AI7//wOt/P8Fwvv/Afn+/wH9//8B/f//Af3//wH9//+N33b//sIF//SKPP/0hUH/9IVB//SFQf/0hUH/l4oi/zaQAv9PpS7/U6g0/1OoNP9TqDT/U6g0/zG5X/8B0Zr/xPTn/+Tk5OCGhoYO/////2i7/v8Aj///AI///wCP//8Bmv3/Bbv7/wW8+/8E2fz/Af3//wH9//8B/f//KfTZ//nIB///xwH/+qwb//SFQf/0hUH/9IVB/9+GOv88jgL/NY8A/z+XEf9TqDT/U6g0/1OoNP9Nqzz/Bs6T/wDRmv8z2q7//P7+//Dw8Ib///+06/X+/xKW/v8Aj///AI/+/wS0+/8FvPv/Bbz7/wW9+/8C8/7/Af3//wH8/v+51kr//8cB///HAf/+xgH/9pI1//SFQf/0hUH/dosX/zWPAP81jwD/NY8A/0uhJ/9TqDT/U6g0/yLAcf8A0Zr/ANGa/wDRmv+Y7Nb/////9sbGxir6+vr4k87+/wCP//8Cov3/Bbz7/wW8+/8FvPv/Bbz7/wTO/P8B/f7/Tuy1///GAf//xwH//8cB///HAf/8txD/9IVA/8WHMf81jgD/NY8A/zWPAP81jwD/OpMJ/1KnM/9DsEj/AdCZ/wDRmv8A0Zr/ANGa/xPUof/s+/f/mZmZAI+Pj4j7/f7/LqX+/wS5+/8FvPv/Bbz7/wW8+/8FvPv/Bbz7/wvo9//czib//8cB///HAf//xwH//8cB///HAP/1nCr/WY0M/zWPAP81jwD/NY8A/zWPAP81jwD/Rp0f/xPHg/8A0Zr/ANGa/wDRmv8A0Zr/ANGa/2vkxP/j4+MAyMjIDv///+jk9v7/yPD+/8jw/v/I8P7/yPD+/8jw/v/I8P7/1vLx///zx///88f///PH///zx///88f///PH//bvyf/T5sf/0+bH/9Pmx//T5sf/0+bH/9Pmx//R6c3/x/Xp/8f16f/H9en/x/Xp/8f16f/H9en/y/Xq///4////8P////B////gP///wD///8Af//+AH///gB///wAf//8AP//+AH///AAD//wAAf/4AAD/+AAA//AAAH/gAAB/4AAAP8AAAD/AAAA/gAAAfwAf//8AP///AD///wAAAAMAAAADAAAAAQAAAAAAAAAAgAAAAIAAAADAAAAA"
  );
});
