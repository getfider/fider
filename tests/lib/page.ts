import { BrowserTab, WaitCondition } from ".";

export interface NewablePage<T extends Page> {
  new (tab: BrowserTab): T;
}

export abstract class Page {
  public async navigate(): Promise<void> {
    await this.tab.navigate(this.getUrl());
    await this.tab.wait(this.loadCondition());
  }

  public abstract loadCondition(): WaitCondition;

  protected getUrl(): string {
    throw new Error("getUrl not implemented");
  }

  public constructor(protected tab: BrowserTab) {}
}
