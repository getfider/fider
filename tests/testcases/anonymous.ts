import { ensure, elementIsVisible } from "../lib";
import { pages, browser } from "../context";

it("Cannot submit ideas", async () => {
  // Action
  await pages.home.navigate();
  await pages.home.IdeaTitle.click();

  // Assert
  await browser.wait(elementIsVisible(() => pages.home.SignInModal));
});

it("Cannot see actions or comment on an idea", async () => {
  // Action & Assert
  await pages.home.navigate();
  await pages.home.IdeaList.navigateAndWait(0, pages.showIdea.loadCondition());
  await ensure(pages.showIdea.RespondButton).isNotVisible();

  // Action & Assert
  await pages.showIdea.CommentInput.click();
  await browser.wait(elementIsVisible(() => pages.home.SignInModal));
});
