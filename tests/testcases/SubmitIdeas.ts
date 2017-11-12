import { ensure, elementIsVisible, mailgun } from '../lib';
import { pages, tenant, browser } from '../context';
import { createTenant } from './setup';

describe('Submit ideas', () => {

  before(async () => {
    await createTenant();
  });

  after(async () => {
    await pages.dispose();
  });

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

  it('User can sign in using e-mail on existing tenant', async () => {
    const now = new Date().getTime();

    // Action
    await pages.home.navigate();
    await pages.home.signInWithEmail(`darthvader.fider+${now}@gmail.com`);

    const link = await mailgun.getLinkFromLastEmailTo(`darthvader.fider+${now}@gmail.com`);

    await pages.goTo(link);
    await browser.wait(pages.home.loadCondition());

    await pages.home.completeSignIn(`Darth Vader ${now}`);

    // Assert
    await pages.home.UserMenu.click();
    await ensure(pages.home.UserName).textIs(`DARTH VADER ${now}`);
  });
});
