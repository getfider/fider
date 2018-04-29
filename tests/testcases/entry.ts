import { Browser } from "../lib";
import { setPages, setBrowser, setTenant } from "../context";
import { AllPages } from "../pages";

describe("Using Fider", () => {
  let browser: Browser;
  let pages: AllPages;

  beforeAll(async () => {
    const now = new Date().getTime();
    const tenantName = `Selenium ${now}`;
    const tenantSubdomain = `selenium${now}`;
    setTenant(tenantSubdomain);

    browser = new Browser("chrome");
    pages = new AllPages(browser);

    setBrowser(browser);
    setPages(pages);

    await pages.signup.navigate();
    await pages.signup.signInWithFacebook();
    await pages.facebook.signInAsJonSnow();

    await pages.signup.signUpAs(tenantName, tenantSubdomain);
  });

  afterAll(async () => {
    await pages.dispose();
  });

  describe("As an Admin", () => {
    require("./admin");
  });

  describe("As an anonymous user", () => {
    beforeAll(async () => {
      await pages.home.navigate();
      await pages.home.signOut();
    });
    require("./anonymous");
  });

  describe("Alternative sign in", () => {
    require("./signin");
  });
});
