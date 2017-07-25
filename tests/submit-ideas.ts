import { Browser, ensure } from './lib';
import { Builder, ThenableWebDriver, WebElement, By, WebElementPromise, until } from 'selenium-webdriver';
import { AllPages } from './pages';

describe('Submit ideas', () => {
  let pages: AllPages;

  before(async () => {
    pages = new AllPages(new Browser('chrome'));
  });

  it('Test Case #1: Unauthenticated cannot submit ideas', async () => {
    // Action
    await pages.home.navigate();
    await pages.home.IdeaTitle.type('Add support to TypeScript');

    // Assert
    await Promise.all([
      ensure(pages.home.UserMenu).textIs('Sign in'),
      ensure(pages.home.SubmitIdea).isNotVisible(),
      ensure(pages.home.IdeaDescription).isNotVisible(),
    ]);
  });

  it('Test Case #2: Authenticated can submit ideas', async () => {
    // Action
    await pages.home.navigate();
    await pages.home.signInWithGoogle();
    await pages.google.signInAsDarthVader();

    // Assert
    await ensure(pages.home.UserMenu).textIs('Darth Vader');

    // Action
    await pages.home.submitNewIdea('Add support to TypeScript', 'Because the language and community is awesome! :)');

    // Assert
    await Promise.all([
      ensure(pages.showIdea.Title).textIs('Add support to TypeScript'),
      ensure(pages.showIdea.Description).textIs('Because the language and community is awesome! :)'),
      ensure(pages.showIdea.SupportCounter).textIs('1'),
    ]);
  });

  after(async () => {
    await pages.dispose();
  });
});
