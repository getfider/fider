import { Comment } from "@fider/models"

export interface Thread {
  id: number
  // undefined when the top-level comment was deleted but its replies remain
  root?: Comment
  replies: Comment[]
}

const time = (c: Comment) => new Date(c.createdAt).getTime()

// Groups a flat comment list into one-level threads: newest thread first, replies oldest first.
// A thread whose top-level comment was deleted is placed by its first reply, the closest known time.
export const groupThreads = (comments: Comment[]): Thread[] => {
  const threads = new Map<number, Thread>()
  for (const comment of comments) {
    const id = comment.parentId ?? comment.id
    const thread = threads.get(id) ?? { id, replies: [] }
    if (comment.parentId) {
      thread.replies.push(comment)
    } else {
      thread.root = comment
    }
    threads.set(id, thread)
  }

  const result = Array.from(threads.values())
  result.forEach((thread) => thread.replies.sort((a, b) => time(a) - time(b)))
  const threadTime = (t: Thread) => time(t.root ?? t.replies[0])
  return result.sort((a, b) => threadTime(b) - threadTime(a))
}
