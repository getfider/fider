process.env.HOST_MODE = "multi"
jest.setTimeout(30000)
require("./scenarios/email-signup")
require("./scenarios/index")
