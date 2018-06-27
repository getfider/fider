import * as React from "react";
import { shallow } from "enzyme";
import { Idea, UserRole, UserStatus, IdeaStatus } from "@fider/models";
import { SupportCounter } from "@fider/components";
import { Fider as FiderImpl } from "../fider"; // TODO: remove this
import { actions, http, httpMock } from "@fider/services";

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

  test.only("click when authenticated and viewerSupported === false", async () => {
    // TODO: remove this hack
    (window as any).Fider = new FiderImpl();
    Object.defineProperty(Fider.session, "isAuthenticated", {
      get() {
        return true;
      }
    });

    const mock = httpMock.alwaysOk();
    http.post = mock;

    const wrapper = shallow(<SupportCounter idea={idea} />);
    wrapper.find("button").simulate("click");
    expect(mock).toHaveBeenCalledWith("/api/ideas/10/support");
    expect(mock).toHaveBeenCalledTimes(1);
  });

  test.only("click when authenticated and viewerSupported === true", async () => {
    // TODO: remove this hack
    idea.viewerSupported = true;
    (window as any).Fider = new FiderImpl();
    Object.defineProperty(Fider.session, "isAuthenticated", {
      get() {
        return true;
      }
    });

    const mock = httpMock.alwaysOk();
    http.post = mock;

    const wrapper = shallow(<SupportCounter idea={idea} />);
    wrapper.find("button").simulate("click");
    expect(mock).toHaveBeenCalledWith("/api/ideas/10/unsupport");
    expect(mock).toHaveBeenCalledTimes(1);
  });
});
