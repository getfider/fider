import React from "react";
import { Modal, Checkbox } from "@fider/components/common";
import { Fider } from "@fider/services";

interface LegalAgreementProps {
  onChange: (agreed: boolean) => void;
}

export const TermsOfService: React.FunctionComponent<{}> = () => {
  if (Fider.settings.hasLegal) {
    return (
      <a href="/terms" target="_blank">
        Terms of Service
      </a>
    );
  }
  return null;
};

export const PrivacyPolicy: React.FunctionComponent<{}> = () => {
  if (Fider.settings.hasLegal) {
    return (
      <a href="/privacy" target="_blank">
        Privacy Policy
      </a>
    );
  }
  return null;
};

export const LegalNotice: React.FunctionComponent<{}> = () => {
  if (Fider.settings.hasLegal) {
    return (
      <p className="info">
        By signing in, you agree to the <PrivacyPolicy /> and <TermsOfService />.
      </p>
    );
  }
  return null;
};

export const LegalFooter: React.FunctionComponent<{}> = () => {
  if (Fider.settings.hasLegal) {
    return (
      <Modal.Footer align="center">
        <LegalNotice />
      </Modal.Footer>
    );
  }
  return null;
};

export const LegalAgreement: React.FunctionComponent<LegalAgreementProps> = props => {
  if (Fider.settings.hasLegal) {
    return (
      <Checkbox field="legalAgreement" onChange={props.onChange}>
        I have read and agree to the <PrivacyPolicy /> and <TermsOfService />.
      </Checkbox>
    );
  }
  return null;
};
