import { ensure, elementIsVisible, elementIsNotVisible, mailgun, delay } from '../lib';
import { pages, tenant, browser } from '../context';
import { createTenant } from './setup';
import {
  equal as assertEqual
} from 'assert';

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

  it('User doesn\'t lose what they typed', async () => {
    // Action
    await pages.home.navigate();
    await pages.home.IdeaTitle.type('My Great Idea');
    await pages.home.IdeaDescription.type('With an awesome description');

    await pages.home.IdeaTitle.clear();
    await pages.home.IdeaTitle.type('My Great Idea has a new title');

    // Assert
    await ensure(pages.home.IdeaTitle).textIs('My Great Idea has a new title');
    await ensure(pages.home.IdeaDescription).textIs('With an awesome description');

    // Action
    await pages.home.navigate();

    // Assert
    await ensure(pages.home.IdeaTitle).textIs('My Great Idea has a new title');
    await ensure(pages.home.IdeaDescription).textIs('With an awesome description');

    // Action
    await pages.home.IdeaDescription.clear();
    await pages.home.IdeaTitle.clear();
    await pages.home.navigate();

    // Assert
    await ensure(pages.home.IdeaTitle).textIs('');
    await browser.wait(elementIsNotVisible(() => pages.home.IdeaDescription));
  });

  it('Unauthenticated users cannot submit ideas', async () => {
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
    await pages.facebook.signInAsJonSnow();
    await pages.home.UserMenu.click();

    // Assert
    await ensure(pages.home.UserName).textIs('JON SNOW');

    // Action
    await pages.home.submitNewIdea('Add support to TypeScript', 'Because the language and community is awesome! :)');

    // Assert
    await ensure(pages.showIdea.Title).textIs('Add support to TypeScript');
    await ensure(pages.showIdea.Description).textIs('Because the language and community is awesome! :)');
    await ensure(pages.showIdea.SupportCounter).textIs('1');
  });

  it('Admin can change status', async () => {
    // Action & Assert
    await pages.home.navigate();
    await pages.home.IdeaList.navigateAndWait(0, pages.showIdea.loadCondition());
    await pages.showIdea.RespondButton.click();
    await browser.wait(elementIsVisible(() => pages.showIdea.ResponseModal));

    // Action & Assert
    await pages.showIdea.changeStatus('Planned', 'We will work on this soon...');
    await browser.wait(elementIsVisible(() => pages.showIdea.Status));
    await ensure(pages.showIdea.Status).textIs('PLANNED');
    await ensure(pages.showIdea.ResponseText).textIs('We will work on this soon...');
  });

  it('Admin can change status to another one', async () => {
    // Action & Assert
    await pages.showIdea.RespondButton.click();
    await browser.wait(elementIsVisible(() => pages.showIdea.ResponseModal));
    await pages.showIdea.changeStatus('Planned', 'We will work on this soon...');
    await browser.wait(elementIsVisible(() => pages.showIdea.Status));
    await ensure(pages.showIdea.Status).textIs('PLANNED');
    await ensure(pages.showIdea.ResponseText).textIs('We will work on this soon...');
  });

  it('Authenticated user can comment on an idea', async () => {
    // Action & Assert
    await pages.home.navigate();
    await pages.home.IdeaList.navigateAndWait(0, pages.showIdea.loadCondition());
    await ensure(pages.showIdea.SubmitCommentButton).isNotVisible();

    // Action & Assert
    await pages.showIdea.CommentInput.type('This is my first comment');
    await ensure(pages.showIdea.SubmitCommentButton).isVisible();

    // Action & Assert
    await pages.showIdea.SubmitCommentButton.click();
    await browser.wait(async () => await pages.showIdea.CommentList.count() === 1);
  });

  it('Unauthenticated user cannot comment on an idea', async () => {
    // Action & Assert
    await pages.home.navigate();
    await pages.home.signOut();
    await pages.home.IdeaList.navigateAndWait(0, pages.showIdea.loadCondition());
    assertEqual(await pages.showIdea.CommentList.count(), 1);

    // Action & Assert
    await pages.showIdea.CommentInput.click();
    await browser.wait(elementIsVisible(() => pages.home.SignInModal));
  });

  it('Unauthenticated user cannot see actions', async () => {
    // Action & Assert
    await pages.home.navigate();
    await pages.home.signOut();
    await pages.home.IdeaList.navigateAndWait(0, pages.showIdea.loadCondition());
    await ensure(pages.showIdea.RespondButton).isNotVisible();
  });

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

  it.skip('Admin can change Title, Invitation and Welcom Message', async () => {
    return true;
  });
});
