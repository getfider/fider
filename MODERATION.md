# Moderation

We are going to implement the ability to moderate posts and comments added to fider. The idea is that when a new post is added, that post is then flagged as unmoderated. An admin will need to approve it.

This document explains everything that needs to change in Fider to facilitate this feature.

## Settings

This is a optional feature. Admins will be able to toggle this. This is done in the public/pages/Administration/pages/PrivacySettings.page.tsx page. Similar to how the other settings are controlled in this page. There needs to be a new column in the "tenants" database table called "is_moderation_enabled" to control this, so you're going to need a new migration file in migrations/

## New posts and comments

Posts and comments will need a new column "is_approved" to determine if the post or comment has been approved to be shown. Again, this will need adding to the migration.

When a new post or comment is added, if moderation is enabled then is_approved will be false, otherwise it will be true. New posts are added via public/pages/Home/components/ShareFeedback.tsx and comments via public/pages/ShowPost/components/CommentInput.tsx.

Once added, the post is only visible to the person who added it (and admins, see below). When you view the post (or the comment) via public/pages/ShowPost/ShowPost.page.tsx, there needs to be a message to tell you that it's awaiting moderation.

## Doing the moderation

The "admin" section of fider looks like this (public/pages/Administration/components/AdminBasePage.tsx):
![alt text](<CleanShot 2025-07-01 at 20.39.05@2x.png>)

There neeeds to be a new menu item on the left for "Moderation"

Clicking on that presents you with a tablular view of all non-moderated posts and comments.
For each row, display the following columns:
_ A checkbox to allow you to select multiple rows
_ User's name and date of post (e.g. "Matt, 10 minutes ago")
_ Wide column for the description.
_ If comment: "New comment: <comment>" (truncated to 200 chars)
_ If post: "New post: <post title>"
_ Thumbs up button to approve \* Thumbs down button to decline

If you click the description for a post or comment, it will take you to the post, and if you clicked a comment, will highlight the comment (this is already supporeted, see how the public/pages/ShowPost/ShowPost.page.tsx page highlights comments). When you are an admin, and it's a post that's awaiting moderation, the in place of the voting button, we need 2 buttons - one to approve, one to decline. The same is true for comments, there should be an approve / decline set of buttons udner the comment.

Declining a post or comment will delete it entirely. We should ask the user to confirm the action.

Approving should just set it as approved, and remove it from the list (might be easier to re-fetch the content for moderation)

Here is a screenshot showing the inspiration we found for the moderation page

![alt text](<CleanShot 2025-07-01 at 20.54.06@2x.png>)

You can see the bulk actions on this screenshot, we're interested in having bulk actions to approve or decline too to make it easier, plus a "select all" to highlight them all. The UI for the bulk actions should be like the "Sort by" options on the post listing in public/pages/Home/components/PostsSort.tsx
