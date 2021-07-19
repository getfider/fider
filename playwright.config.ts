import { PlaywrightTestConfig } from "@playwright/test"
const config: PlaywrightTestConfig = {
  testDir: "e2e",
  use: {
    slowMo: 50,
    viewport: { width: 1280, height: 720 },
    ignoreHTTPSErrors: true,
  },
}
export default config
