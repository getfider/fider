import React from "react";
import { injectStripe, CardElement } from "react-stripe-elements";
import { Input, Field, Button, Form, Select, SelectOption, Modal, CardInfo } from "@fider/components";
import { Failure, actions } from "@fider/services";
import { PaymentInfo, Country } from "@fider/models";

type PatchedTokenResponse = stripe.TokenResponse & {
  error?: { decline_code?: string };
};

interface StripeProps {
  createToken(options?: stripe.TokenOptions): Promise<PatchedTokenResponse>;
  createSource(sourceData?: stripe.SourceOptions): Promise<stripe.SourceResponse>;
  paymentRequest: stripe.Stripe["paymentRequest"];
}

interface PaymentInfoModalProps {
  paymentInfo?: PaymentInfo;
  stripe?: StripeProps;
  countries: Country[];
  onClose: () => void;
}

interface PaymentInfoModalState {
  changingCard: boolean;
  name: string;
  email: string;
  addressLine1: string;
  addressLine2: string;
  addressCity: string;
  addressState: string;
  addressPostalCode: string;
  addressCountry: string;
  vatNumber: string;
  stripe: stripe.Stripe | null;
  error?: Failure;
}

class PaymentInfoModal extends React.Component<PaymentInfoModalProps, PaymentInfoModalState> {
  constructor(props: PaymentInfoModalProps) {
    super(props);
    this.state = {
      stripe: null,
      changingCard: false,
      name: this.props.paymentInfo ? this.props.paymentInfo.name : "",
      email: this.props.paymentInfo ? this.props.paymentInfo.email : "",
      addressLine1: this.props.paymentInfo ? this.props.paymentInfo.addressLine1 : "",
      addressLine2: this.props.paymentInfo ? this.props.paymentInfo.addressLine2 : "",
      addressCity: this.props.paymentInfo ? this.props.paymentInfo.addressCity : "",
      addressState: this.props.paymentInfo ? this.props.paymentInfo.addressState : "",
      addressPostalCode: this.props.paymentInfo ? this.props.paymentInfo.addressPostalCode : "",
      addressCountry: this.props.paymentInfo ? this.props.paymentInfo.addressCountry : "",
      vatNumber: this.props.paymentInfo ? this.props.paymentInfo.vatNumber : ""
    };
  }

  public handleSubmit = async () => {
    if (this.props.paymentInfo && !this.state.changingCard) {
      const response = await actions.updatePaymentInfo({
        ...this.state
      });

      if (response.ok) {
        location.reload();
      } else {
        this.setState({
          error: response.error
        });
      }

      return;
    }

    if (this.props.stripe) {
      const result = await this.props.stripe.createToken({
        name: this.state.name,
        address_line1: this.state.addressLine1,
        address_line2: this.state.addressLine2,
        address_city: this.state.addressCity,
        address_state: this.state.addressState,
        address_zip: this.state.addressPostalCode,
        address_country: this.state.addressCountry
      });

      if (result.token) {
        const response = await actions.updatePaymentInfo({
          ...this.state,
          card: {
            type: result.token.type,
            token: result.token.id,
            country: result.token.card ? result.token.card.country : ""
          }
        });

        if (response.ok) {
          location.reload();
        } else {
          this.setState({
            error: response.error
          });
        }
      } else if (result.error) {
        this.setState({
          error: {
            errors: [
              {
                field: "card",
                message: result.error.message!
              }
            ]
          }
        });
      }
    }
  };

  private setName = (name: string) => {
    this.setState({ name });
  };

  private setEmail = (email: string) => {
    this.setState({ email });
  };

  private setAddressLine1 = (addressLine1: string) => {
    this.setState({ addressLine1 });
  };

  private setAddressLine2 = (addressLine2: string) => {
    this.setState({ addressLine2 });
  };

  private setAddressCity = (addressCity: string) => {
    this.setState({ addressCity });
  };

  private setAddressState = (addressState: string) => {
    this.setState({ addressState });
  };

  private setAddressPostalCode = (addressPostalCode: string) => {
    this.setState({ addressPostalCode });
  };

  private setVATNumber = (vatNumber: string) => {
    this.setState({ vatNumber });
  };

  private setAddressCountry = (option: SelectOption | undefined) => {
    if (option) {
      this.setState({ addressCountry: option.value });
    }
  };

  private closeModal = async () => {
    this.props.onClose();
  };

  private changeCard = () => {
    this.setState({ changingCard: true });
  };

  private isEUCountry(): boolean {
    if (this.state.addressCountry) {
      const filtered = this.props.countries.filter(x => x.code === this.state.addressCountry);
      if (filtered && filtered.length > 0) {
        return filtered[0].isEU;
      }
    }
    return false;
  }

  public render() {
    return (
      <Modal.Window isOpen={true} center={false} size="large" onClose={this.closeModal}>
        <Modal.Content>
          <Form className="c-payment-info-modal" error={this.state.error}>
            <div className="row">
              {(!this.props.paymentInfo || this.state.changingCard) && (
                <div className="col-md-12">
                  <Field label="Credit/Debit Card" field="card">
                    <CardElement />
                    <p className="info">
                      We neither store nor see your card information. We integrate directly with Stripe.
                    </p>
                  </Field>
                </div>
              )}
              {this.props.paymentInfo && !this.state.changingCard && (
                <div className="col-md-12">
                  <Field
                    label="Credit/Debit Card"
                    field="card"
                    afterLabel={
                      <span className="ui info clickable" onClick={this.changeCard}>
                        change
                      </span>
                    }
                  >
                    <CardInfo
                      expMonth={this.props.paymentInfo.cardExpMonth}
                      expYear={this.props.paymentInfo.cardExpYear}
                      brand={this.props.paymentInfo.cardBrand}
                      last4={this.props.paymentInfo.cardLast4}
                    />
                  </Field>
                </div>
              )}
              <div className="col-md-12">
                <Input label="Name" field="name" value={this.state.name} onChange={this.setName} autoComplete="name" />
              </div>
              <div className="col-md-12">
                <Input
                  label="Email"
                  field="email"
                  value={this.state.email}
                  onChange={this.setEmail}
                  autoComplete="email"
                />
              </div>
              <div className="col-md-6">
                <Input
                  label="Address Line 1"
                  value={this.state.addressLine1}
                  field="addressLine1"
                  onChange={this.setAddressLine1}
                  autoComplete="address-line1"
                />
              </div>
              <div className="col-md-6">
                <Input
                  label="Address Line 2"
                  value={this.state.addressLine2}
                  field="addressLine2"
                  onChange={this.setAddressLine2}
                  autoComplete="address-line2"
                />
              </div>
              <div className="col-md-4">
                <Input
                  label="City"
                  field="addressCity"
                  value={this.state.addressCity}
                  onChange={this.setAddressCity}
                  autoComplete="address-level2"
                />
              </div>
              <div className="col-md-4">
                <Input
                  label="State / Region"
                  field="addressState"
                  value={this.state.addressState}
                  onChange={this.setAddressState}
                  autoComplete="address-level1"
                />
              </div>
              <div className="col-md-4">
                <Input
                  label="Postal Code"
                  field="addressPostalCode"
                  value={this.state.addressPostalCode}
                  onChange={this.setAddressPostalCode}
                  autoComplete="postal-code"
                />
              </div>
              <div className="col-md-6">
                <Select
                  label="Country"
                  field="addressCountry"
                  onChange={this.setAddressCountry}
                  defaultValue={this.state.addressCountry}
                  options={[
                    { value: "", label: "" },
                    ...this.props.countries.map(x => ({ value: x.code, label: x.name }))
                  ]}
                />
              </div>
              {this.isEUCountry() && (
                <div className="col-md-6">
                  <Input
                    label="VAT Number"
                    field="vatNumber"
                    value={this.state.vatNumber}
                    onChange={this.setVATNumber}
                  />
                </div>
              )}
              {!!this.state.addressCountry && (
                <div className="col-md-12">
                  <Field field="currency">
                    <p className="info">
                      Based on your Billing Address, your subscription will be charged in{" "}
                      {this.isEUCountry() ? "Euro (EUR)" : "US Dollar (USD)"}.
                    </p>
                  </Field>
                </div>
              )}
            </div>
          </Form>
        </Modal.Content>

        <Modal.Footer>
          <Button onClick={this.handleSubmit} color="positive">
            Save
          </Button>
          <Button color="cancel" onClick={this.closeModal}>
            Close
          </Button>
        </Modal.Footer>
      </Modal.Window>
    );
  }
}

export default injectStripe(PaymentInfoModal);
