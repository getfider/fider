import { ensure } from "../lib"
import { ctx } from "."

it("Tab1: User can change general settings", async () => {
  // Action
  await ctx.tab1.pages.generalSettings.navigate()
  await ctx.tab1.pages.generalSettings.changeSettings("Feedback for X", "Leave your feedback below", "Enter here...")

  // Assert
  await ensure(ctx.tab1.pages.home.MenuTitle).textIs("Feedback for X")
  await ensure(ctx.tab1.pages.home.PostTitle).attributeIs("placeholder", "Enter here...")
  await ensure(ctx.tab1.pages.home.WelcomeMessage).textIs("Leave your feedback below")
})
