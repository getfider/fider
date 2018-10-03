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

  describe("when error caught", () => {
    test("error should be passed to onError", () => {
      let error;
      const BadComponent = () => {
        error = new Error("Whoops!");
        throw error;
      };

      const errorSpy = jest.fn();
      const wrapper = mount(
        <ErrorBoundary onError={errorSpy}>
          <BadComponent />
          <div id="no-error">No Error!</div>
        </ErrorBoundary>
      );

      expect(errorSpy.mock.calls[0][0]).toEqual(error);
    });

    test("error should not be shown when showError is false", () => {
      const message = "My Error";
      const BadComponent = () => {
        throw new Error(message);
      };

      const wrapper = mount(
        <ErrorBoundary showError={false}>
          <BadComponent />
          <div id="no-error">No Error!</div>
        </ErrorBoundary>
      );

      expect(
        wrapper
          .render()
          .text()
          .indexOf(message) === -1
      ).toBeTruthy();
    });

    test("error should be hidden when showError is false", () => {
      const message = "My Error 12345";
      const BadComponent = () => {
        throw new Error(message);
      };

      const wrapper = mount(
        <ErrorBoundary showError={true}>
          <BadComponent />
          <div id="no-error">No Error!</div>
        </ErrorBoundary>
      );

      expect(
        wrapper
          .render()
          .text()
          .indexOf(message) > -1
      ).toBeTruthy();
    });
  });
});
