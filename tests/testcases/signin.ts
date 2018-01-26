import { ensure, elementIsVisible, mailgun } from '../lib';
import { pages, browser } from '../context';

it('Can sign in with Google and support an idea', async () => {
  // Action & Assert
  await pages.home.navigate();
  await pages.home.signInWithGoogle();
  await pages.google.signInAsDarthVader();
  await pages.home.UserMenu.click();
  await ensure(pages.home.UserName).textIs('DARTH VADER');

  // Action & Assert
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
