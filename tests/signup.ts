import { Browser, ensure } from './lib';
import { Builder, ThenableWebDriver, WebElement, By, WebElementPromise, until } from 'selenium-webdriver';
import { AllPages } from './pages';

describe('Signup as a new tenant', () => {
  let pages: AllPages;

  before(async () => {
    pages = new AllPages(new Browser('chrome'));
  });

  it('Test Case #1: Can sign up using Google', async () => {
    // Action
    await pages.signup.navigate();
    await pages.signup.signInWithGoogle();
    await pages.google.signInAsDarthVader();
    const name = `Selenium${new Date().getTime()}`;
    await pages.signup.signUpAs(name, name.toLowerCase());
  });

  after(async () => {
    await pages.dispose();
  });
});
