import { ensure, Browser } from '../lib';
import { AllPages } from '../pages';
import { setPages, setBrowser, setTenant } from '../context';

export function initialize() {
  const browser = new Browser('chrome');
  const pages = new AllPages(browser);

  setBrowser(browser);
  setPages(pages);
}

export async function createTenant() {
  const now = new Date().getTime();
  const tenantName = `Selenium ${now}`;
  const tenantSubdomain = `selenium${now}`;
  setTenant(tenantSubdomain);

  const browser = new Browser('chrome');
  const pages = new AllPages(browser);

  setBrowser(browser);
  setPages(pages);

  await pages.signup.navigate();
  await pages.signup.signInWithFacebook();
  await pages.facebook.signInAsJonSnow();

  await pages.signup.signUpAs(tenantName, tenantSubdomain);
}
