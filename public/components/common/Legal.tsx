import * as React from "react";
import { Modal, Field, Checkbox } from "@fider/components/common";
import { page } from "@fider/services";

const privacy = (
  <a href="/privacy" target="_blank">
    privacy policy
  </a>
);

const terms = (
  <a href="/terms" target="_blank">
    terms of service
  </a>
);

interface LegalAgreementProps {
  onChange: (agrred: boolean) => void;
}

export const LegalNotice: React.StatelessComponent<{}> = () => {
  if (page.systemSettings().hasLegal) {
    return (
      <p className="info">
        By signing in, you agree to the {privacy} and {terms}.
      </p>
    );
  }
  return null;
};

export const LegalFooter: React.StatelessComponent<{}> = () => {
  if (page.systemSettings().hasLegal) {
    return (
      <Modal.Footer align="center">
        <LegalNotice />
      </Modal.Footer>
    );
  }
  return null;
};

export const LegalAgreement: React.StatelessComponent<LegalAgreementProps> = props => {
  if (page.systemSettings().hasLegal) {
    return (
      <Checkbox field="legalAgreement" onChange={props.onChange}>
        I have read and agree to the {privacy} and {terms}.
      </Checkbox>
    );
  }
  return null;
};
