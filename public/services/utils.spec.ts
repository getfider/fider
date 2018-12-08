import { classSet, formatDate, timeSince, fileToBase64 } from "./utils";
import { readFileSync } from "fs";

[
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
        throw Error();
      }
    },
    expected: "green"
  }
].forEach(x => {
  test(`classSet of ${JSON.stringify(x.input)} should be ${x.expected}`, () => {
    const className = classSet(x.input);
    expect(className).toEqual(x.expected);
  });
});

[
  { input: new Date(2018, 4, 27, 10, 12, 59), expected: "May 27, 2018 · 10:12" },
  { input: new Date(2058, 12, 12, 23, 21, 53), expected: "January 12, 2059 · 23:21" },
  { input: "2018-04-11T18:13:33.128082", expected: "April 11, 2018 · 18:13" },
  { input: "2017-11-20T07:47:42.158142", expected: "November 20, 2017 · 07:47" }
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
    "AAABAAEAEBAAAAEAIABoBAAAFgAAACgAAAAQAAAAIAAAAAEAIAAAAAAAAAQAABILAAASCwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB0YiwAfIqMAHBeHeR8bk6MgHJkQHxyWAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACoAggA2H8QAOB3FKi0dsNsiHZ38Ix6mYCIdnwAmILEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABBXPAAR4T/BTw32Jg3Ic3/KR/G/yYgtZAnIMYAJiCxAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYzPYASEjvADhz8kpCW+7wOj/l/ygg4/4nINVsKCDpACYgsQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAG8v2ABvS9xMcwfa9N4v0/05r8f8+UenPJBjaHCce2wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYMT/AADg7QA8y/xxGNL5/BjK+f9Kn/f/YaP4xVWb9zVDkfYCTpn2AAAAAAAAAAAAAAAAAAAAAAAAAAAAj///AEDO4wBDzuMrKdHyxxLU+vckyvv/QLz9/lOo+f5QlfbnVZv3Y4zH/AFhpfgAAAAAAAAAAAAAAAAAAAAAAJaZUQC/gSIFN9C9mynWy88W0fSaH8n47D/D/f9Pw/3/UKX3/2Cj+ONgpvguYaX4AAAAAAAAAAAAAAAAAMuKGwDAgh4AxoAYQY+dVu85z7f1KNrMZxvJ+zNBwv+EQsP40y217PpIrPT/Ra36Y0is+QAAAAAAAAAAAAAAAACOTxoAxIIcAKtqG3+5eBT/h5hR/z7KpOg9yI9xUrtwGCS27hcbtu1ZJrTorzCz8FYss+0AAAAAAAAAAAAAAAAAhEUZAIM9KQCCPSRollMc/rt9FP6ZoUj/R8OI+0jBgtJ6mi96kHoAJDfBwwgRzv8HE8z/AAAAAAAAAAAAAAAAAAAAAAB9NisAezQsG4M5Kb2oYSP/xIUa/4qXQv5Hv3v/V6tH/nGBGOBYWwySYmcUNuztXgXQ0E4AAAAAAAAAAAAAAAAAfjcqAJA8JgCIOSgekUIleaxpFMuwcw/4f4gv/0yySv9OmTT/T14R/2JrFOufpTR60c9CBra3NgAAAAAAAAAAAAAAAAAAAAAAlDsRAEkAHwCnXhQUpWoTUqFnFaZ/jifnTqs6/k6MIP9VZw3/c30Y8YySHkWKkR4AaXcXAAAAAAAAAAAAAAAAAAAAAAAAAAAAn2YVALdyCgC3agoGsXkQLmCYJHxMmRzLXJIc9W19Fv9zexSAbXgVAHd+EwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE+WGgA5dAoAUp4eGGKzKV9wkRvAfYASbXp/EwB/ghIA/j8AAPw/AAD4PwAA+D8AAPA/AADwHwAA4A8AAMAPAADADwAAwA8AAMAPAADABwAA4AMAAPADAAD+AwAA/4MAAA=="
  );
});
