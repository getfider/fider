import { ensure, elementIsVisible } from './lib';
import { pages, tenant, browser } from './context';

describe('Submit ideas', () => {

  it('Unauthenticated cannot submit ideas', async () => {
    // Action
    await pages.home.navigate();
    await pages.home.IdeaTitle.click();

    // Assert
    await browser.wait(elementIsVisible(() => pages.home.SignInModal));
  });

  it('Authenticated user can submit ideas', async () => {
    // Action
    await pages.home.navigate();
    await pages.home.signInWithFacebook();
    await pages.facebook.signInAsAryaStark();

    // Assert
    await ensure(pages.home.UserMenu).textIs('Arya Stark');

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

    // Assert
    await ensure(pages.home.UserMenu).textIs('Darth Vader');

    await pages.home.IdeaList.want(0);
    await ensure(await pages.home.IdeaList.at(0)).textIs('2');
  });
});
