import { ensure, elementIsVisible, mailgun } from '../lib';
import { pages, tenant, browser } from '../context';
import config from '../config';

describe('Sign up by e-mail', () => {
  it('User can sign up using e-mail', async () => {
    const now = new Date().getTime();

    // Action
    await pages.signup.navigate();
    await pages.signup.signInWithEmail('Darth Vader 1', `darthvader.fider+1@gmail.com`);
    await pages.signup.signUpAs(`Selenium ${now}`, `selenium${now}`);

    const link = await mailgun.getLinkFromLastEmailTo(`darthvader.fider+1@gmail.com`);

    await pages.goTo(link);
    browser.wait(pages.home.loadCondition());

    // Assert
    await pages.home.UserMenu.click();
    await ensure(pages.home.UserName).textIs('DARTH VADER 1');
  });
});
