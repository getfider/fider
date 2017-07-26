import { Browser } from './lib';
import { AllPages } from './pages';
import { setPages, setBrowser, setTenant } from './context';

let browser: Browser;
let pages: AllPages;

before(async () => {
  const now = new Date().getTime();
  const tenantName = `Selenium ${now}`;
  const tenantSubdomain = `selenium${now}`;
  setTenant(tenantSubdomain);

  browser = new Browser('chrome');
  pages = new AllPages(browser);

  await pages.signup.navigate();
  await pages.signup.signInWithFacebook();
  await pages.facebook.signInAsJonSnow();

  await pages.signup.signUpAs(tenantName, tenantSubdomain);

  setBrowser(browser);
  setPages(pages);
});

after(async () => {
  await pages.dispose();
});
