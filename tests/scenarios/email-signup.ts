import { Browser, mailgun, pageHasLoaded, ensure } from "../lib";
import { AllPages, HomePage } from "../pages";

describe("E2E: Sign up with e-mail", () => {
  let browser: Browser;
  let pages: AllPages;

  beforeAll(async () => {
    browser = await Browser.launch();
    pages = new AllPages(browser);
  });

  afterAll(async () => {
    await browser.close();
  });

  it("User can sign up using email", async () => {
    const now = new Date().getTime();

    // Action
    await pages.signup.navigate();
    await pages.signup.signInWithEmail(`Darth Vader ${now}`, `darthvader.fider@gmail.com`);
    await pages.signup.signUpAs(`Selenium ${now}`, `selenium${now}`);

    const link = await mailgun.getLinkFromLastEmailTo(`selenium${now}`, `darthvader.fider@gmail.com`);

    await pages.goTo(link);
    await browser.wait(pageHasLoaded(HomePage));

    // Assert
    await pages.home.UserMenu.click();
    await ensure(pages.home.UserName).textIs(`DARTH VADER ${now}`);
  });
});
