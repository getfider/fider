import { Browser, pageHasLoaded, ensure, elementIsNotVisible, delay } from "../lib";
import { AllPages, HomePage } from "../pages";
import { setTenant } from "../context";

describe("E2E: Feedback Loop", () => {
  let browser: Browser;
  let pages: AllPages;
  let tenantName: string;
  let tenantSubdomain: string;

  beforeAll(async () => {
    const now = new Date().getTime();
    tenantName = `Selenium ${now}`;
    tenantSubdomain = `selenium${now}`;
    setTenant(tenantSubdomain);
    browser = await Browser.launch();
    pages = new AllPages(browser);
  });

  afterAll(async () => {
    await browser.close();
  });

  it("User can sign up with facebook", async () => {
    await pages.signup.navigate();
    await pages.signup.signInWithFacebook();
    await pages.facebook.signInAsJonSnow();

    await pages.signup.signUpAs(tenantName, tenantSubdomain);

    await browser.wait(pageHasLoaded(HomePage));
  });

  it("User is authenticated after sign up", async () => {
    // Action
    await pages.home.navigate();
    await pages.home.UserMenu.click();

    // Assert
    await ensure(pages.home.UserName).textIs("JON SNOW");
  });

  it("User doesn't lose what they typed", async () => {
    // Action
    await pages.home.navigate();
    await pages.home.IdeaTitle.type("My Great Idea");
    await pages.home.IdeaDescription.type("With an awesome description");
    await pages.home.IdeaTitle.clear();
    await pages.home.IdeaTitle.type("My Great Idea has a new title");

    // Assert
    await ensure(pages.home.IdeaTitle).textIs("My Great Idea has a new title");
    await ensure(pages.home.IdeaDescription).textIs("With an awesome description");

    // Action
    await pages.home.navigate();

    // Assert
    await ensure(pages.home.IdeaTitle).textIs("My Great Idea has a new title");
    await ensure(pages.home.IdeaDescription).textIs("With an awesome description");

    // Action
    await pages.home.IdeaDescription.clear();
    await pages.home.IdeaTitle.clear();
    await pages.home.navigate();

    // Assert
    await ensure(pages.home.IdeaTitle).textIs("");
    await browser.wait(elementIsNotVisible(pages.home.IdeaDescription));
  });

  it("Can submit ideas", async () => {
    // Action
    await pages.home.navigate();
    await pages.home.submitNewIdea("Add support to TypeScript", "Because the language and community is awesome! :)");

    // Assert
    await ensure(pages.showIdea.Title).textIs("Add support to TypeScript");
    await ensure(pages.showIdea.Description).textIs("Because the language and community is awesome! :)");
    await ensure(pages.showIdea.SupportCounter).textIs("1");
  });

  it.skip("Can edit title and description", async () => {
    return Promise.resolve(true);
  });
});
