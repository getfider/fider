import { Browser, BrowserTab, pageHasLoaded, ensure, elementIsNotVisible, delay } from "../lib";
import { HomePage } from "../pages";
import { setTenant } from "../context";

describe("E2E: Feedback Loop", () => {
  let browser1: Browser;
  let browser2: Browser;
  let tab1: BrowserTab;
  let tab2: BrowserTab;
  let tenantName: string;
  let tenantSubdomain: string;

  beforeAll(async () => {
    const now = new Date().getTime();
    tenantName = `Selenium ${now}`;
    tenantSubdomain = `selenium${now}`;
    setTenant(tenantSubdomain);

    browser1 = await Browser.launch();
    tab1 = await browser1.newTab();
    browser2 = await Browser.launch();
    tab2 = await browser2.newTab();
  });

  afterAll(async () => {
    await browser1.close();
    await browser2.close();
  });

  it("Tab1: User can sign up with facebook", async () => {
    await tab1.pages.signup.navigate();
    await tab1.pages.signup.signInWithFacebook();
    await tab1.pages.facebook.signInAsJonSnow();

    await tab1.pages.signup.signUpAs(tenantName, tenantSubdomain);

    await tab1.wait(pageHasLoaded(HomePage));
  });

  it("Tab1: User is authenticated after sign up", async () => {
    // Action
    await tab1.pages.home.navigate();
    await tab1.pages.home.UserMenu.click();

    // Assert
    await ensure(tab1.pages.home.UserName).textIs("JON SNOW");
  });

  it("Tab1: User doesn't lose what they typed", async () => {
    // Action
    await tab1.pages.home.navigate();
    await tab1.pages.home.IdeaTitle.type("My Great Idea");
    await tab1.pages.home.IdeaDescription.type("With an awesome description");
    await tab1.pages.home.IdeaTitle.clear();
    await tab1.pages.home.IdeaTitle.type("My Great Idea has a new title");

    // Assert
    await ensure(tab1.pages.home.IdeaTitle).textIs("My Great Idea has a new title");
    await ensure(tab1.pages.home.IdeaDescription).textIs("With an awesome description");

    // Action
    await tab1.pages.home.navigate();

    // Assert
    await ensure(tab1.pages.home.IdeaTitle).textIs("My Great Idea has a new title");
    await ensure(tab1.pages.home.IdeaDescription).textIs("With an awesome description");

    // Action
    await tab1.pages.home.IdeaDescription.clear();
    await tab1.pages.home.IdeaTitle.clear();
    await tab1.pages.home.navigate();

    // Assert
    await ensure(tab1.pages.home.IdeaTitle).textIs("");
    await tab1.wait(elementIsNotVisible(tab1.pages.home.IdeaDescription));
  });

  it("Tab1: Can submit ideas", async () => {
    // Action
    await tab1.pages.home.navigate();
    await tab1.pages.home.submitNewIdea(
      "Add support to TypeScript",
      "Because the language and community is awesome! :)"
    );

    // Assert
    await ensure(tab1.pages.showIdea.Title).textIs("Add support to TypeScript");
    await ensure(tab1.pages.showIdea.Description).textIs("Because the language and community is awesome! :)");
    await ensure(tab1.pages.showIdea.SupportCounter).textIs("1");
  });

  it("Tab1: Can edit title and description", async () => {
    // Action
    await tab1.pages.showIdea.edit("Support for TypeScript", "Because the language and community is awesome!");

    // Assert
    await ensure(tab1.pages.showIdea.Title).textIs("Support for TypeScript");
    await ensure(tab1.pages.showIdea.Description).textIs("Because the language and community is awesome!");
    await ensure(tab1.pages.showIdea.SupportCounter).textIs("1");
  });

  it("Tab2: Can login as another user", async () => {
    // Action
    await tab2.pages.home.navigate();
    await tab2.pages.home.signInWithFacebook();
    await tab2.pages.facebook.signInAsAryaStark();
    await tab2.wait(pageHasLoaded(HomePage));
    await tab2.pages.home.UserMenu.click();

    // Assert
    await ensure(tab2.pages.home.UserName).textIs("ARYA STARK");
  });

  it("Tab2: User can vote on idea", async () => {
    // Action
    await tab2.pages.home.navigate();
    const item = await tab2.pages.home.IdeaList.get("Support for TypeScript");
    await item.Vote.click();

    // Assert
    await ensure(item.Vote).textIs("2");
  });
});
