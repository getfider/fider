import { BrowserTab, WaitCondition } from "."

export type NewablePage<T extends Page> = new (tab: BrowserTab) => T

export abstract class Page {
  public async navigate(): Promise<void> {
    await this.tab.navigate(this.getURL())
    await this.tab.wait(this.loadCondition())
  }

  public abstract loadCondition(): WaitCondition

  protected getURL(): string {
    throw new Error("getURL not implemented")
  }

  public constructor(protected tab: BrowserTab) {}
}
