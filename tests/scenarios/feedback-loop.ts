import { pageHasLoaded, ensure, elementIsNotVisible, elementIsVisible, mailgun } from "../lib"
import { HomePage, ShowPostPage } from "../pages"
import { ctx } from "."

it("Tab1: User is authenticated after sign up", async () => {
  // Action
  await ctx.tab1.pages.home.navigate()

  // Assert
  await ensure(ctx.tab1.pages.home.UserName).textIs("Jon Snow")
})

it("Tab1: User doesn't lose what they typed", async () => {
  // Action
  await ctx.tab1.pages.home.navigate()
  await ctx.tab1.pages.home.PostTitle.type("My Great Post")
  await ctx.tab1.pages.home.PostDescription.type("With an awesome description")
  await ctx.tab1.pages.home.PostTitle.clear()
  await ctx.tab1.pages.home.PostTitle.type("My Great Post has a new title")

  // Assert
  await ensure(ctx.tab1.pages.home.PostTitle).textIs("My Great Post has a new title")
  await ensure(ctx.tab1.pages.home.PostDescription).textIs("With an awesome description")

  // Action
  await ctx.tab1.pages.home.navigate()

  // Assert
  await ensure(ctx.tab1.pages.home.PostTitle).textIs("My Great Post has a new title")
  await ensure(ctx.tab1.pages.home.PostDescription).textIs("With an awesome description")

  // Action
  await ctx.tab1.pages.home.PostDescription.clear()
  await ctx.tab1.pages.home.PostTitle.clear()
  await ctx.tab1.pages.home.navigate()

  // Assert
  await ensure(ctx.tab1.pages.home.PostTitle).textIs("")
  await ctx.tab1.wait(elementIsNotVisible(ctx.tab1.pages.home.PostDescription))
})

it("Tab1: Can submit Posts", async () => {
  // Action
  await ctx.tab1.pages.home.navigate()
  await ctx.tab1.pages.home.submitNewPost("Add support to TypeScript", "Because the language and community is awesome! :)")

  // Assert
  await ensure(ctx.tab1.pages.showPost.Title).textIs("Add support to TypeScript")
  await ensure(ctx.tab1.pages.showPost.Description).textIs("Because the language and community is awesome! :)")
  await ensure(ctx.tab1.pages.showPost.VoteCounter).textIs("1")
})

it("Tab1: Can edit title and description", async () => {
  // Action
  await ctx.tab1.pages.showPost.edit("Support for TypeScript", "Because the language and community is awesome!")

  // Assert
  await ensure(ctx.tab1.pages.showPost.Title).textIs("Support for TypeScript")
  await ensure(ctx.tab1.pages.showPost.Description).textIs("Because the language and community is awesome!")
  await ensure(ctx.tab1.pages.showPost.VoteCounter).textIs("1")
})

it("Tab2: Can login as another user", async () => {
  // Action
  await ctx.tab2.pages.home.navigate()
  await ctx.tab2.pages.home.signInWithFacebook()
  await ctx.tab2.pages.facebook.signInAsAryaStark()
  await ctx.tab2.wait(pageHasLoaded(HomePage))

  // Assert
  await ensure(ctx.tab2.pages.home.UserName).textIs("Arya Stark")
})

it("Tab2: User can vote on Post", async () => {
  // Action
  await ctx.tab2.pages.home.navigate()
  const item = await ctx.tab2.pages.home.PostList.get("Support for TypeScript")
  await item.Vote.click()

  // Assert
  await ensure(item.Vote).textIs("2")
})

it("Tab2: Open Post and vote on it", async () => {
  // Action
  const item = await ctx.tab2.pages.home.PostList.get("Support for TypeScript")
  await item.navigate()
  await ctx.tab2.pages.showPost.comment("I support this request!")

  // Assert
  await ensure(ctx.tab2.pages.showPost.Comments).countIs(1)
})

it("Tab1: Open Post and change status", async () => {
  // Action
  await ctx.tab1.pages.home.navigate()
  const item = await ctx.tab1.pages.home.PostList.get("Support for TypeScript")
  await item.navigate()
  await ctx.tab1.pages.showPost.changeStatus("Started", "This will be delivered on next release.")
  await ctx.tab1.wait(elementIsVisible(ctx.tab1.pages.showPost.Status))

  // Assert
  await ensure(ctx.tab1.pages.showPost.Status).textIs("Started")
  await ensure(ctx.tab1.pages.showPost.ResponseText).textIs("This will be delivered on next release.")
})

it("Tab2: Refresh and see status", async () => {
  // Action
  await ctx.tab2.reload(ShowPostPage)

  // Assert
  await ensure(ctx.tab1.pages.showPost.Status).textIs("Started")
  await ensure(ctx.tab1.pages.showPost.ResponseText).textIs("This will be delivered on next release.")
})

it("Tab2: Check notifications", async () => {
  // Action
  await ctx.tab2.pages.home.navigate()
  await ctx.tab2.pages.home.UserMenu.click()

  // Assert
  await ensure(ctx.tab2.pages.home.UnreadCounter).textIs("1")
})

it("Tab2: Logout and sign in with email", async () => {
  // Action
  await ctx.tab2.pages.home.navigate()
  await ctx.tab2.pages.home.signInWithEmail("darthvader.fider@gmail.com")
  const link = await mailgun.getLinkFromLastEmailTo(ctx.tenantSubdomain, `Sign in to ${ctx.tenantName}`, `darthvader.fider@gmail.com`)
  await ctx.tab2.navigate(link)
  await ctx.tab2.pages.home.completeSignIn("Darth Vader")

  // Assert
  await ensure(ctx.tab2.pages.home.UserName).textIs("Darth Vader")
})
