# Fider

Fider is a feedback portal: people post ideas, vote on them, and discuss them.

## Language

### Discussion

**Comment**:
Anything someone posts in a post's discussion, whether at the top level or as a **Reply**.
_Avoid_: Message, note

**Reply**:
A **Comment** whose parent is a top-level **Comment**. It is a full Comment in every respect: reactions, editing, deleting, attachments, mentions, moderation and counts. A Reply's parent is never itself a Reply.
_Avoid_: Sub-comment, child comment, nested comment

**Thread**:
A top-level **Comment** plus all of its **Replies**. A Thread whose top-level Comment was deleted is **closed**: its Replies stay readable, but no new Replies can be added.
_Avoid_: Conversation, chain

## Example dialogue

> **Dev:** If Carol hits Reply on Bob's Reply, does her Comment nest under Bob's?
> **Domain expert:** No. Threads are one level deep, so Carol's Reply joins the end of the same Thread. We just pre-fill an @mention of Bob so it's clear who she's answering.
> **Dev:** And if Alice deletes the top-level Comment?
> **Domain expert:** The Thread is closed. Bob's and Carol's Replies stay visible under a "deleted" placeholder, but nobody can add to it.
