import * as React from "react";
import { shallow } from "enzyme";
import { Post, UserRole, UserStatus, PostStatus } from "@fider/models";
import { SupportCounter } from "@fider/components";
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
      role: UserRole.Collaborator,
      status: UserStatus.Active
    },
    viewerSupported: false,
    response: null,
    totalSupporters: 5,
    totalComments: 2,
    tags: []
  };
});

describe("<SupportCounter />", () => {
  test("when viewerSupported === true", () => {
    post.viewerSupported = true;
    post.totalSupporters = 9;
    const wrapper = shallow(<SupportCounter post={post} />);
    const button = wrapper.find("button");
    expect(button.text()).toBe("9");
    expect(button.hasClass("m-supported")).toBe(true);
    expect(button.hasClass("m-disabled")).toBe(false);
  });

  test("when viewerSupported === false", () => {
    post.viewerSupported = false;
    post.totalSupporters = 2;
    const wrapper = shallow(<SupportCounter post={post} />);
    const button = wrapper.find("button");
    expect(button.text()).toBe("2");
    expect(button.hasClass("m-supported")).toBe(false);
    expect(button.hasClass("m-disabled")).toBe(false);
  });

  test("when post is closed", () => {
    post.status = PostStatus.Completed.value;
    const wrapper = shallow(<SupportCounter post={post} />);
    const button = wrapper.find("button");
    expect(button.text()).toBe("5");
    expect(button.hasClass("m-supported")).toBe(false);
    expect(button.hasClass("m-disabled")).toBe(true);
  });

  test("click when unauthenticated", async () => {
    fiderMock.notAuthenticated();

    const mock = httpMock.alwaysOk();

    const wrapper = shallow(<SupportCounter post={post} />);
    wrapper.find("button").simulate("click");
    await rerender(wrapper);
    expect(wrapper.find("SignInModal").length).toBe(1);
    expect(mock).toHaveBeenCalledTimes(0);
  });

  test("click when authenticated and viewerSupported === false", async () => {
    fiderMock.authenticated();

    const mock = httpMock.alwaysOk();

    const wrapper = shallow(<SupportCounter post={post} />);
    wrapper.find("button").simulate("click");
    expect(mock).toHaveBeenCalledWith("/api/posts/10/support");
    expect(mock).toHaveBeenCalledTimes(1);

    await rerender(wrapper);
    expect(wrapper.find("button").text()).toBe("6");
  });

  test("click when authenticated and viewerSupported === true", async () => {
    post.viewerSupported = true;
    fiderMock.authenticated();

    const mock = httpMock.alwaysOk();

    const wrapper = shallow(<SupportCounter post={post} />);
    wrapper.find("button").simulate("click");
    expect(mock).toHaveBeenCalledWith("/api/posts/10/unsupport");
    expect(mock).toHaveBeenCalledTimes(1);

    await rerender(wrapper);
    expect(wrapper.find("button").text()).toBe("4");
  });
});
