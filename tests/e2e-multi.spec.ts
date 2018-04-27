process.env.HOST_MODE = "multi";
jest.setTimeout(30000);
require("./testcases/email-signup");
require("./testcases/entry");
