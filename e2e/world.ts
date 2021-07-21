import { World as CucumberWorld } from "@cucumber/cucumber"
import { BrowserContext, Page } from "playwright"

export interface FiderWorld extends CucumberWorld {
  tenantName: string
  context: BrowserContext
  page: Page
}
