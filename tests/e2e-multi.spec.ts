process.env.HOST_MODE = "multi";
jest.setTimeout(60000);
require("./testcases/email-signup");
require("./testcases/entry");
