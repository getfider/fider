import {
  Browser,
  BrowserTab,
  pageHasLoaded,
  ensure,
  elementIsNotVisible,
  delay,
  elementIsVisible,
  mailgun
} from "../lib";
import { HomePage, ShowIdeaPage } from "../pages";
import { ctx } from ".";

it("Tab1: User is authenticated after sign up", async () => {
  // Action
  await ctx.tab1.pages.home.navigate();

  // Assert
  await ensure(ctx.tab1.pages.home.UserName).textIs("Jon Snow");
});

it("Tab1: User doesn't lose what they typed", async () => {
  // Action
  await ctx.tab1.pages.home.navigate();
  await ctx.tab1.pages.home.IdeaTitle.type("My Great Idea");
  await ctx.tab1.pages.home.IdeaDescription.type("With an awesome description");
  await ctx.tab1.pages.home.IdeaTitle.clear();
  await ctx.tab1.pages.home.IdeaTitle.type("My Great Idea has a new title");

  // Assert
  await ensure(ctx.tab1.pages.home.IdeaTitle).textIs("My Great Idea has a new title");
  await ensure(ctx.tab1.pages.home.IdeaDescription).textIs("With an awesome description");

  // Action
  await ctx.tab1.pages.home.navigate();

  // Assert
  await ensure(ctx.tab1.pages.home.IdeaTitle).textIs("My Great Idea has a new title");
  await ensure(ctx.tab1.pages.home.IdeaDescription).textIs("With an awesome description");

  // Action
  await ctx.tab1.pages.home.IdeaDescription.clear();
  await ctx.tab1.pages.home.IdeaTitle.clear();
  await ctx.tab1.pages.home.navigate();

  // Assert
  await ensure(ctx.tab1.pages.home.IdeaTitle).textIs("");
  await ctx.tab1.wait(elementIsNotVisible(ctx.tab1.pages.home.IdeaDescription));
});

it("Tab1: Can submit ideas", async () => {
  // Action
  await ctx.tab1.pages.home.navigate();
  await ctx.tab1.pages.home.submitNewIdea(
    "Add support to TypeScript",
    "Because the language and community is awesome! :)"
  );

  // Assert
  await ensure(ctx.tab1.pages.showIdea.Title).textIs("Add support to TypeScript");
  await ensure(ctx.tab1.pages.showIdea.Description).textIs("Because the language and community is awesome! :)");
  await ensure(ctx.tab1.pages.showIdea.SupportCounter).textIs("1");
});

it("Tab1: Can edit title and description", async () => {
  // Action
  await ctx.tab1.pages.showIdea.edit("Support for TypeScript", "Because the language and community is awesome!");

  // Assert
  await ensure(ctx.tab1.pages.showIdea.Title).textIs("Support for TypeScript");
  await ensure(ctx.tab1.pages.showIdea.Description).textIs("Because the language and community is awesome!");
  await ensure(ctx.tab1.pages.showIdea.SupportCounter).textIs("1");
});

it("Tab2: Can login as another user", async () => {
  // Action
  await ctx.tab2.pages.home.navigate();
  await ctx.tab2.pages.home.signInWithFacebook();
  await ctx.tab2.pages.facebook.signInAsAryaStark();
  await ctx.tab2.wait(pageHasLoaded(HomePage));

  // Assert
  await ensure(ctx.tab2.pages.home.UserName).textIs("Arya Stark");
});

it("Tab2: User can vote on idea", async () => {
  // Action
  await ctx.tab2.pages.home.navigate();
  const item = await ctx.tab2.pages.home.IdeaList.get("Support for TypeScript");
  await item.Vote.click();

  // Assert
  await ensure(item.Vote).textIs("2");
});

it("Tab2: Open idea and vote on it", async () => {
  // Action
  const item = await ctx.tab2.pages.home.IdeaList.get("Support for TypeScript");
  await item.navigate();
  await ctx.tab2.pages.showIdea.comment("I support this request!");

  // Assert
  await ensure(ctx.tab2.pages.showIdea.Comments).countIs(1);
});

it("Tab1: Open idea and change status", async () => {
  // Action
  await ctx.tab1.pages.home.navigate();
  const item = await ctx.tab1.pages.home.IdeaList.get("Support for TypeScript");
  await item.navigate();
  await ctx.tab1.pages.showIdea.changeStatus("Started", "This will be delivered on next release.");
  await ctx.tab1.wait(elementIsVisible(ctx.tab1.pages.showIdea.Status));

  // Assert
  await ensure(ctx.tab1.pages.showIdea.Status).textIs("Started");
  await ensure(ctx.tab1.pages.showIdea.ResponseText).textIs("This will be delivered on next release.");
});

it("Tab2: Refresh and see status", async () => {
  // Action
  await ctx.tab2.reload(ShowIdeaPage);

  // Assert
  await ensure(ctx.tab1.pages.showIdea.Status).textIs("Started");
  await ensure(ctx.tab1.pages.showIdea.ResponseText).textIs("This will be delivered on next release.");
});

it("Tab2: Check notifications", async () => {
  // Action
  await ctx.tab2.pages.home.navigate();
  await ctx.tab2.pages.home.UserMenu.click();

  // Assert
  await ensure(ctx.tab2.pages.home.UnreadCounter).textIs("1");
});

it("Tab2: Logout and sign in with email", async () => {
  // Action
  await ctx.tab2.pages.home.navigate();
  await ctx.tab2.pages.home.signInWithEmail("darthvader.fider@gmail.com");
  const link = await mailgun.getLinkFromLastEmailTo(ctx.tenantSubdomain, `darthvader.fider@gmail.com`);
  await ctx.tab2.navigate(link);
  await ctx.tab2.pages.home.completeSignIn("Darth Vader");

  // Assert
  await ensure(ctx.tab2.pages.home.UserName).textIs("Darth Vader");
});
