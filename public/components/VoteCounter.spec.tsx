import React from "react"
import { Post, UserRole, PostStatus, UserStatus } from "@fider/models"
import { VoteCounter } from "@fider/components"
import { screen, fireEvent, render } from "@testing-library/react"
import { fiderMock, httpMock, setupModalRoot } from "@fider/services/testing"
import { FiderContext } from "@fider/services"
import { act } from "react-dom/test-utils"

let post: Post

beforeEach(() => {
  setupModalRoot()

  post = {
    id: 1,
    number: 10,
    slug: "add-typescript",
    title: "Add TypeScript",
    description: "",
    createdAt: "",
    status: PostStatus.Started.value,
    user: {
      id: 5,
      name: "John",
      role: UserRole.Collaborator,
      status: UserStatus.Active,
      avatarURL: "/static/avatars/letter/5/John",
    },
    hasVoted: false,
    response: null,
    votesCount: 5,
    commentsCount: 2,
    tags: [],
  }
})

describe("<VoteCounter />", () => {
  test("when hasVoted === true", () => {
    post.hasVoted = true
    post.votesCount = 9

    const { container } = render(
      <FiderContext.Provider value={fiderMock.authenticated()}>
        <VoteCounter post={post} />
      </FiderContext.Provider>
    )
    const button = container.querySelector("button")

    expect(button).toHaveTextContent("9")
    expect(button).toHaveClass("c-vote-counter__button--voted")
    expect(button).not.toHaveClass("c-vote-counter__button--disabled")
  })

  test("when hasVoted === false", () => {
    post.hasVoted = false
    post.votesCount = 2
    const { container } = render(
      <FiderContext.Provider value={fiderMock.authenticated()}>
        <VoteCounter post={post} />
      </FiderContext.Provider>
    )
    const button = container.querySelector("button")
    expect(button).toHaveTextContent("2")
    expect(button).not.toHaveClass("c-vote-counter__button--voted")
    expect(button).not.toHaveClass("c-vote-counter__button--disabled")
  })

  test("when post is closed", () => {
    post.status = PostStatus.Completed.value
    const { container } = render(
      <FiderContext.Provider value={fiderMock.authenticated()}>
        <VoteCounter post={post} />
      </FiderContext.Provider>
    )
    const button = container.querySelector("button")
    expect(button).toHaveTextContent("5")
    expect(button).toHaveClass("c-vote-counter__button--disabled")
    expect(button).not.toHaveClass("c-vote-counter__button--voted")
  })

  test("click when unauthenticated", async () => {
    const mock = httpMock.alwaysOk()

    const { container } = render(
      <FiderContext.Provider value={fiderMock.notAuthenticated()}>
        <VoteCounter post={post} />
      </FiderContext.Provider>
    )
    const button = container.querySelector("button") || fail("button not found")
    fireEvent.click(button)

    expect(screen.queryByTestId("modal")).toBeInTheDocument()
    expect(mock.post).toHaveBeenCalledTimes(0)
    expect(mock.delete).toHaveBeenCalledTimes(0)
  })

  test("click when authenticated and hasVoted === false", async () => {
    const mock = httpMock.alwaysOk()

    const { container } = render(
      <FiderContext.Provider value={fiderMock.authenticated()}>
        <VoteCounter post={post} />
      </FiderContext.Provider>
    )

    const button = container.querySelector("button") || fail("button not found")
    expect(button).toHaveTextContent("5")
    await act(async () => {
      fireEvent.click(button)
    })

    expect(mock.post).toHaveBeenCalledWith("/api/v1/posts/10/votes")
    expect(mock.post).toHaveBeenCalledTimes(1)
    expect(button).toHaveTextContent("6")
  })

  test("click when authenticated and hasVoted === true", async () => {
    post.hasVoted = true

    const mock = httpMock.alwaysOk()

    const { container } = render(
      <FiderContext.Provider value={fiderMock.authenticated()}>
        <VoteCounter post={post} />
      </FiderContext.Provider>
    )

    const button = container.querySelector("button") || fail("button not found")
    expect(button).toHaveTextContent("5")
    await act(async () => {
      fireEvent.click(button)
    })

    expect(mock.delete).toHaveBeenCalledWith("/api/v1/posts/10/votes")
    expect(mock.delete).toHaveBeenCalledTimes(1)
    expect(button).toHaveTextContent("4")
  })
})
