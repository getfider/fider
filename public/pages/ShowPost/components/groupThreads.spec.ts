import { Comment } from "@fider/models"
import { groupThreads } from "./groupThreads"

// Ids deliberately don't follow creation order (e.g. imported data)
const comment = (id: number, minute: number, parentId?: number) =>
  ({ id, parentId, createdAt: `2026-10-09T10:${String(minute).padStart(2, "0")}:00Z` }) as Comment

describe("groupThreads", () => {
  test("groups replies under their parent, newest thread first and replies oldest first", () => {
    const threads = groupThreads([comment(2, 50, 90), comment(80, 40), comment(5, 30, 90), comment(1, 20, 90), comment(90, 10)])

    expect(threads.map((t) => t.id)).toEqual([80, 90])
    expect(threads[0].root?.id).toEqual(80)
    expect(threads[0].replies).toEqual([])
    expect(threads[1].root?.id).toEqual(90)
    expect(threads[1].replies.map((c) => c.id)).toEqual([1, 5, 2])
  })

  test("places a thread whose top-level comment was deleted by its first reply", () => {
    const threads = groupThreads([comment(7, 50), comment(3, 40, 60), comment(4, 30, 60), comment(9, 10)])

    expect(threads.map((t) => t.id)).toEqual([7, 60, 9])
    expect(threads[1].root).toBeUndefined()
    expect(threads[1].replies.map((c) => c.id)).toEqual([4, 3])
  })
})
