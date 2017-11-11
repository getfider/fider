import { ensure, elementIsVisible } from '../lib';
import { pages, tenant, browser } from '../context';

describe('Sign up by e-mail', () => {
    it('User can sign up using e-mail', async () => {

      // Action
      await pages.inboxBear.navigate();
      await pages.inboxBear.clearInbox();

      const now = new Date().getTime();

      await pages.signup.navigate();
      await pages.signup.signInWithEmail('Selenium Tester', 'fiderselenium@inboxbear.com');
      await pages.signup.signUpAs(`Selenium ${now}`, `selenium${now}`);

      await pages.inboxBear.navigate();
      const link = await pages.inboxBear.getLinkFromEmail(0);

      await pages.goTo(link);
      browser.wait(pages.home.loadCondition());

      // Assert
      await pages.home.UserMenu.click();
      await ensure(pages.home.UserName).textIs('SELENIUM TESTER');
    });
});
