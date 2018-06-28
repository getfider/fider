import * as React from "react";
import { shallow } from "enzyme";
import { Idea, UserRole, UserStatus, IdeaStatus } from "@fider/models";
import { SupportCounter } from "@fider/components";
import { httpMock, fiderMock, rerender } from "@fider/services/testing";

let idea: Idea;

beforeEach(() => {
  idea = {
    id: 1,
    number: 10,
    slug: "add-typescript",
    title: "Add TypeScript",
    description: "",
    createdOn: "",
    status: IdeaStatus.Started.value,
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
    idea.viewerSupported = true;
    idea.totalSupporters = 9;
    const wrapper = shallow(<SupportCounter idea={idea} />);
    const button = wrapper.find("button");
    expect(button.text()).toBe("9");
    expect(button.hasClass("m-supported")).toBe(true);
    expect(button.hasClass("m-disabled")).toBe(false);
  });

  test("when viewerSupported === false", () => {
    idea.viewerSupported = false;
    idea.totalSupporters = 2;
    const wrapper = shallow(<SupportCounter idea={idea} />);
    const button = wrapper.find("button");
    expect(button.text()).toBe("2");
    expect(button.hasClass("m-supported")).toBe(false);
    expect(button.hasClass("m-disabled")).toBe(false);
  });

  test("when idea is closed", () => {
    idea.status = IdeaStatus.Completed.value;
    const wrapper = shallow(<SupportCounter idea={idea} />);
    const button = wrapper.find("button");
    expect(button.text()).toBe("5");
    expect(button.hasClass("m-supported")).toBe(false);
    expect(button.hasClass("m-disabled")).toBe(true);
  });

  test("click when unauthenticated", async () => {
    fiderMock.notAuthenticated();

    const mock = httpMock.alwaysOk();

    const wrapper = shallow(<SupportCounter idea={idea} />);
    wrapper.find("button").simulate("click");
    await rerender(wrapper);
    expect(wrapper.find("SignInModal").length).toBe(1);
    expect(mock).toHaveBeenCalledTimes(0);
  });

  test("click when authenticated and viewerSupported === false", async () => {
    fiderMock.authenticated();

    const mock = httpMock.alwaysOk();

    const wrapper = shallow(<SupportCounter idea={idea} />);
    wrapper.find("button").simulate("click");
    expect(mock).toHaveBeenCalledWith("/api/ideas/10/support");
    expect(mock).toHaveBeenCalledTimes(1);

    await rerender(wrapper);
    expect(wrapper.find("button").text()).toBe("6");
  });

  test("click when authenticated and viewerSupported === true", async () => {
    idea.viewerSupported = true;
    fiderMock.authenticated();

    const mock = httpMock.alwaysOk();

    const wrapper = shallow(<SupportCounter idea={idea} />);
    wrapper.find("button").simulate("click");
    expect(mock).toHaveBeenCalledWith("/api/ideas/10/unsupport");
    expect(mock).toHaveBeenCalledTimes(1);

    await rerender(wrapper);
    expect(wrapper.find("button").text()).toBe("4");
  });
});
