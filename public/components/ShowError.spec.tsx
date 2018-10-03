import * as React from "react";
import { shallow } from "enzyme";
import { ShowError } from "@fider/components";

describe("<ShowError />", () => {
  const createFakeErrorInfo = () => ({ componentStack: "" } as React.ErrorInfo);

  test("it should show the error when showError returns true", () => {
    const error = new Error("Hello");
    const errorInfo = createFakeErrorInfo();

    const showError = () => true;
    const wrapper = shallow(<ShowError error={error} errorInfo={errorInfo} showError={showError} />);
    expect(wrapper.find("pre")).toHaveLength(1);
  });

  test("it should not show the error when showError returns false", () => {
    const error = new Error("Hello");
    const errorInfo = createFakeErrorInfo();

    const showError = () => false;
    const wrapper = shallow(<ShowError error={error} errorInfo={errorInfo} showError={showError} />);

    expect(wrapper.find("pre")).toHaveLength(0);
  });
});
