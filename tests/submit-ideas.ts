import { ensure, elementIsVisible } from './lib';
import { pages, tenant, browser } from './context';

describe('Submit ideas', () => {

  it('User is authenticated after sign up', async () => {
    // Action
    await pages.home.navigate();
    await pages.home.UserMenu.click();

    // Assert
    await ensure(pages.home.UserName).textIs('JON SNOW');
  });

  it('Unauthenticated cannot submit ideas', async () => {
    // Action
    await pages.home.navigate();
    await pages.home.signOut();
    await pages.home.IdeaTitle.click();

    // Assert
    await browser.wait(elementIsVisible(() => pages.home.SignInModal));
  });

  it('Authenticated user can submit ideas', async () => {
    // Action
    await pages.home.navigate();
    await pages.home.signInWithFacebook();
    await pages.facebook.signInAsAryaStark();
    await pages.home.UserMenu.click();

    // Assert
    await ensure(pages.home.UserName).textIs('ARYA STARK');

    // Action
    await pages.home.submitNewIdea('Add support to TypeScript', 'Because the language and community is awesome! :)');

    // Assert
    await Promise.all([
      ensure(pages.showIdea.Title).textIs('Add support to TypeScript'),
      ensure(pages.showIdea.Description).textIs('Because the language and community is awesome! :)'),
      ensure(pages.showIdea.SupportCounter).textIs('1'),
    ]);
  });

  it('Can sign in with Google and support an idea', async () => {
    // Action
    await pages.home.navigate();
    await pages.home.signInWithGoogle();
    await pages.google.signInAsDarthVader();
    await pages.home.UserMenu.click();

    // Assert
    await ensure(pages.home.UserName).textIs('DARTH VADER');

    await pages.home.IdeaList.want(0);
    await ensure(await pages.home.IdeaList.at(0)).textIs('2');
  });
});
