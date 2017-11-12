import { ensure, elementIsVisible, mailgun, delay } from '../lib';
import { pages, tenant, browser } from '../context';
import { initialize } from './setup';

describe('Sign up by e-mail', () => {
  before(async () => {
    initialize();
  });

  after(async () => {
    await pages.dispose();
  });

  it('User can sign up using e-mail', async () => {
    const now = new Date().getTime();

    // Action
    await pages.signup.navigate();
    await pages.signup.signInWithEmail(`Darth Vader ${now}`, `darthvader.fider+${now}@gmail.com`);
    await pages.signup.signUpAs(`Selenium ${now}`, `selenium${now}`);

    const link = await mailgun.getLinkFromLastEmailTo(`darthvader.fider+${now}@gmail.com`);

    await pages.goTo(link);
    await browser.wait(pages.home.loadCondition());

    // Assert
    await pages.home.UserMenu.click();
    await ensure(pages.home.UserName).textIs(`DARTH VADER ${now}`);
  });
});
