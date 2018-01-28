import { ensure, elementIsVisible, elementIsNotVisible, } from '../lib';
import { pages, browser } from '../context';

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

it('Can submit ideas', async () => {
  // Action
  await pages.home.navigate();
  await pages.home.submitNewIdea('Add support to TypeScript', 'Because the language and community is awesome! :)');

  // Assert
  await ensure(pages.showIdea.Title).textIs('Add support to TypeScript');
  await ensure(pages.showIdea.Description).textIs('Because the language and community is awesome! :)');
  await ensure(pages.showIdea.SupportCounter).textIs('1');
});

it('Can change status', async () => {
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

it('Can change status to another one', async () => {
  // Action & Assert
  await pages.showIdea.RespondButton.click();
  await browser.wait(elementIsVisible(() => pages.showIdea.ResponseModal));
  await pages.showIdea.changeStatus('Planned', 'We will work on this soon...');
  await browser.wait(elementIsVisible(() => pages.showIdea.Status));
  await ensure(pages.showIdea.Status).textIs('PLANNED');
  await ensure(pages.showIdea.ResponseText).textIs('We will work on this soon...');
});

it('Can comment on an idea', async () => {
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

it('Can change Title, Invitation and Welcome Message', async () => {
  await pages.adminSettings.navigate();

  await pages.adminSettings.TitleInput.clear();
  await pages.adminSettings.TitleInput.type('Selenium Feedback');

  await pages.adminSettings.WelcomeMessageInput.clear();
  await pages.adminSettings.WelcomeMessageInput.type('Welcome to our Selenium Feedback Forum');

  await pages.adminSettings.InvitationInput.clear();
  await pages.adminSettings.InvitationInput.type('Say something...');

  await pages.adminSettings.ConfirmButton.click();

  await browser.wait(pages.home.loadCondition());

  await ensure(pages.home.MenuTitle).textIs('Selenium Feedback');
  await ensure(pages.home.WelcomeMessage).textIs('Welcome to our Selenium Feedback Forum');
  await ensure(pages.home.IdeaTitle).attributeIs('placeholder', 'Say something...');
});
