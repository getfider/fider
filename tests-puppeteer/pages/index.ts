import { SignUpPage } from "./SignUpPage";
import { Browser } from "../lib";

export { SignUpPage };

export class AllPages {
  public signup: SignUpPage;

  constructor(public browser: Browser) {
    this.signup = new SignUpPage(browser);
  }
}
