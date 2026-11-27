import React from "react"
import { fireEvent, render, screen } from "@testing-library/react"
import { Post, PostStatus, Tag, UserRole, UserStatus } from "@fider/models"
import { FiderContext } from "@fider/services"
import { fiderMock } from "@fider/services/testing"
import { ListPosts } from "./ListPosts"

const tag: Tag = { id: 1, name: "Bug", slug: "bug", color: "ff0000", isPublic: true }
const post: Post = {
  id: 1,
  number: 10,
  slug: "add-typescript",
  title: "Add TypeScript",
  description: "",
  createdAt: "",
  status: PostStatus.Open.value,
  user: {
    id: 5,
    name: "John",
    isTrusted: false,
    role: UserRole.Collaborator,
    status: UserStatus.Active,
    avatarURL: "/static/avatars/letter/5/John",
  },
  hasVoted: false,
  response: null,
  votesCount: 5,
  commentsCount: 0,
  tags: [tag.slug],
  isApproved: true,
}

describe("<ListPosts />", () => {
  const onPostClick = jest.fn()

  beforeEach(() => {
    onPostClick.mockClear()
    render(
      <FiderContext.Provider value={fiderMock.notAuthenticated()}>
        <ListPosts posts={[post]} tags={[tag]} emptyText="" onPostClick={onPostClick} />
      </FiderContext.Provider>
    )
  })

  test("links tags to the tag filter", () => {
    expect(screen.getByRole("link", { name: tag.name })).toHaveAttribute("href", "/?tags=bug")
  })

  test("allows tag navigation without opening the post", () => {
    const allowed = fireEvent.click(screen.getByRole("link", { name: tag.name }))

    expect(onPostClick).not.toHaveBeenCalled()
    expect(allowed).toBe(true)
  })

  test("opens the post when its title is clicked", () => {
    fireEvent.click(screen.getByText(post.title))

    expect(onPostClick).toHaveBeenCalledTimes(1)
    expect(onPostClick).toHaveBeenCalledWith(post.number, post.slug)
  })
})
