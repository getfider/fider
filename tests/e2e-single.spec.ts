process.env.HOST_MODE = "single";
jest.setTimeout(60000);
require("./testcases/entry");
