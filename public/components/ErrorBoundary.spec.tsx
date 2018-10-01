import * as React from "react";
import { mount } from "enzyme";
import { ErrorBoundary, ShowError } from "@fider/components";

describe("<ErrorBoundary />", () => {
  let errorMethod: () => void;

  // Stub out console.error to hide noisy Virtual DOM exceptions.
  beforeAll(() => {
    errorMethod = console.error; // tslint:disable-line
    console.error = () => null; // tslint:disable-line
  });

  afterAll(() => {
    console.error = errorMethod; // tslint:disable-line
  });

  test("when no error caught", () => {
    const wrapper = mount(
      <ErrorBoundary>
        <div id="no-error">No Error!</div>
      </ErrorBoundary>
    );

    expect(wrapper.find("#no-error").length).toEqual(1);
  });

  test("when error caught", () => {
    const BadComponent = () => {
      throw new Error("Whoops!");
    };

    const wrapper = mount(
      <ErrorBoundary>
        <BadComponent />
        <div id="no-error">No Error!</div>
      </ErrorBoundary>
    );

    expect(wrapper.find("#no-error").length).not.toEqual(1);
  });
});
