import { Browser, pageHasLoaded, ensure, elementIsNotVisible, delay } from "../lib";
import { AllPages, HomePage } from "../pages";
import { setTenant } from "../context";

describe("E2E: Feedback Loop", () => {
  let browser1: Browser;
  let page1: AllPages;
  let browser2: Browser;
  let page2: AllPages;
  let tenantName: string;
  let tenantSubdomain: string;

  beforeAll(async () => {
    const now = new Date().getTime();
    tenantName = `Selenium ${now}`;
    tenantSubdomain = `selenium${now}`;
    setTenant(tenantSubdomain);

    browser1 = await Browser.launch();
    page1 = new AllPages(browser1);

    browser2 = await Browser.launch();
    page2 = new AllPages(browser2);
  });

  afterAll(async () => {
    await browser1.close();
    await browser2.close();
  });

  it("Browser1: User can sign up with facebook", async () => {
    await page1.signup.navigate();
    await page1.signup.signInWithFacebook();
    await page1.facebook.signInAsJonSnow();

    await page1.signup.signUpAs(tenantName, tenantSubdomain);

    await browser1.wait(pageHasLoaded(HomePage));
  });

  it("Browser1: User is authenticated after sign up", async () => {
    // Action
    await page1.home.navigate();
    await page1.home.UserMenu.click();

    // Assert
    await ensure(page1.home.UserName).textIs("JON SNOW");
  });

  it("Browser1: User doesn't lose what they typed", async () => {
    // Action
    await page1.home.navigate();
    await page1.home.IdeaTitle.type("My Great Idea");
    await page1.home.IdeaDescription.type("With an awesome description");
    await page1.home.IdeaTitle.clear();
    await page1.home.IdeaTitle.type("My Great Idea has a new title");

    // Assert
    await ensure(page1.home.IdeaTitle).textIs("My Great Idea has a new title");
    await ensure(page1.home.IdeaDescription).textIs("With an awesome description");

    // Action
    await page1.home.navigate();

    // Assert
    await ensure(page1.home.IdeaTitle).textIs("My Great Idea has a new title");
    await ensure(page1.home.IdeaDescription).textIs("With an awesome description");

    // Action
    await page1.home.IdeaDescription.clear();
    await page1.home.IdeaTitle.clear();
    await page1.home.navigate();

    // Assert
    await ensure(page1.home.IdeaTitle).textIs("");
    await browser1.wait(elementIsNotVisible(page1.home.IdeaDescription));
  });

  it("Browser1: Can submit ideas", async () => {
    // Action
    await page1.home.navigate();
    await page1.home.submitNewIdea("Add support to TypeScript", "Because the language and community is awesome! :)");

    // Assert
    await ensure(page1.showIdea.Title).textIs("Add support to TypeScript");
    await ensure(page1.showIdea.Description).textIs("Because the language and community is awesome! :)");
    await ensure(page1.showIdea.SupportCounter).textIs("1");
  });

  it("Browser1: Can edit title and description", async () => {
    // Action
    await page1.showIdea.edit("Support for TypeScript", "Because the language and community is awesome!");

    // Assert
    await ensure(page1.showIdea.Title).textIs("Support for TypeScript");
    await ensure(page1.showIdea.Description).textIs("Because the language and community is awesome!");
    await ensure(page1.showIdea.SupportCounter).textIs("1");
  });

  it("Browser2: Can login as another user", async () => {
    // Action
    await page2.home.navigate();
    await page2.home.signInWithFacebook();
    await page2.facebook.signInAsAryaStark();
    await browser2.wait(pageHasLoaded(HomePage));
    await page2.home.UserMenu.click();

    // Assert
    await ensure(page2.home.UserName).textIs("ARYA STARK");
  });

  it("Browser2: User can vote on idea", async () => {
    // Action
    await page2.home.navigate();
    const item = await page2.home.IdeaList.get("Support for TypeScript");
    await item.Vote.click();

    // Assert
    await ensure(item.Vote).textIs("2");
  });
});
