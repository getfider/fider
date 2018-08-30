import * as React from "react";
import { shallow } from "enzyme";
import { Post, UserRole, PostStatus } from "@fider/models";
import { VoteCounter } from "@fider/components";
import { httpMock, fiderMock, rerender } from "@fider/services/testing";

let post: Post;

beforeEach(() => {
  post = {
    id: 1,
    number: 10,
    slug: "add-typescript",
    title: "Add TypeScript",
    description: "",
    createdOn: "",
    status: PostStatus.Started.value,
    user: {
      id: 5,
      name: "John",
      role: UserRole.Collaborator
    },
    hasVoted: false,
    response: null,
    totalVotes: 5,
    totalComments: 2,
    tags: []
  };
});

describe("<VoteCounter />", () => {
  test("when hasVoted === true", () => {
    post.hasVoted = true;
    post.totalVotes = 9;
    const wrapper = shallow(<VoteCounter post={post} />);
    const button = wrapper.find("button");
    expect(button.text()).toBe("9");
    expect(button.hasClass("m-voted")).toBe(true);
    expect(button.hasClass("m-disabled")).toBe(false);
  });

  test("when hasVoted === false", () => {
    post.hasVoted = false;
    post.totalVotes = 2;
    const wrapper = shallow(<VoteCounter post={post} />);
    const button = wrapper.find("button");
    expect(button.text()).toBe("2");
    expect(button.hasClass("m-voted")).toBe(false);
    expect(button.hasClass("m-disabled")).toBe(false);
  });

  test("when post is closed", () => {
    post.status = PostStatus.Completed.value;
    const wrapper = shallow(<VoteCounter post={post} />);
    const button = wrapper.find("button");
    expect(button.text()).toBe("5");
    expect(button.hasClass("m-voted")).toBe(false);
    expect(button.hasClass("m-disabled")).toBe(true);
  });

  test("click when unauthenticated", async () => {
    fiderMock.notAuthenticated();

    const mock = httpMock.alwaysOk();

    const wrapper = shallow(<VoteCounter post={post} />);
    wrapper.find("button").simulate("click");
    await rerender(wrapper);
    expect(wrapper.find("SignInModal").length).toBe(1);
    expect(mock.post).toHaveBeenCalledTimes(0);
    expect(mock.delete).toHaveBeenCalledTimes(0);
  });

  test("click when authenticated and hasVoted === false", async () => {
    fiderMock.authenticated();

    const mock = httpMock.alwaysOk();

    const wrapper = shallow(<VoteCounter post={post} />);
    wrapper.find("button").simulate("click");
    expect(mock.post).toHaveBeenCalledWith("/api/v1/posts/10/votes");
    expect(mock.post).toHaveBeenCalledTimes(1);

    await rerender(wrapper);
    expect(wrapper.find("button").text()).toBe("6");
  });

  test("click when authenticated and hasVoted === true", async () => {
    post.hasVoted = true;
    fiderMock.authenticated();

    const mock = httpMock.alwaysOk();

    const wrapper = shallow(<VoteCounter post={post} />);
    wrapper.find("button").simulate("click");
    expect(mock.delete).toHaveBeenCalledWith("/api/v1/posts/10/votes");
    expect(mock.delete).toHaveBeenCalledTimes(1);

    await rerender(wrapper);
    expect(wrapper.find("button").text()).toBe("4");
  });
});
