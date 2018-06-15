import { Browser } from ".";
import { WaitCondition } from "./conditions";

export interface NewablePage<T extends Page> {
  new (browser: Browser): T;
}

export abstract class Page {
  public async navigate(): Promise<void> {
    await this.browser.navigate(this.getUrl());
    await this.browser.wait(this.loadCondition());
  }

  public abstract loadCondition(): WaitCondition;

  protected getUrl(): string {
    throw new Error("getUrl not implemented");
  }

  public constructor(protected browser: Browser) {}
}
