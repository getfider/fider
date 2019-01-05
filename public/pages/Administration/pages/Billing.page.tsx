import "./Billing.page.scss";

import React from "react";

import { FaFileInvoice } from "react-icons/fa";
import { AdminBasePage } from "../components/AdminBasePage";
import { Segment, Button, CardBrand } from "@fider/components";
import { PaymentInfo } from "@fider/models";
import { Fider } from "@fider/services";
import PaymentInfoModal from "../components/PaymentInfoModal";
import { StripeProvider, Elements } from "react-stripe-elements";

interface BillingPageProps {
  paymentInfo?: PaymentInfo;
  countries: Array<{ code: string; name: string }>;
}

interface BillingPageState {
  stripe: stripe.Stripe | null;
  showModal: boolean;
}

export default class BillingPage extends AdminBasePage<BillingPageProps, BillingPageState> {
  public id = "p-admin-billing";
  public name = "billing";
  public icon = FaFileInvoice;
  public title = "Billing";
  public subtitle = "Manage your subscription";

  constructor(props: BillingPageProps) {
    super(props);
    this.state = {
      stripe: null,
      showModal: false
    };
  }

  private openModal = async () => {
    if (!this.state.stripe) {
      const script = document.createElement("script");
      script.src = "https://js.stripe.com/v3/";
      script.onload = () => {
        this.setState({
          stripe: Stripe(Fider.settings.stripePublicKey!),
          showModal: true
        });
      };
      document.body.appendChild(script);
    } else {
      this.setState({
        showModal: true
      });
    }
  };

  private closeModal = async () => {
    this.setState({
      showModal: false
    });
  };

  public content() {
    return (
      <>
        {this.state.showModal && (
          <StripeProvider stripe={this.state.stripe}>
            <Elements>
              <PaymentInfoModal
                paymentInfo={this.props.paymentInfo}
                countries={this.props.countries}
                onClose={this.closeModal}
              />
            </Elements>
          </StripeProvider>
        )}
        <div className="row">
          <div className="col-md-3">
            <Segment className="l-payment-info">
              <h4>Payment Info</h4>
              {this.props.paymentInfo && (
                <>
                  <p>
                    <CardBrand brand={this.props.paymentInfo.cardBrand} />
                    <span>**** **** **** {this.props.paymentInfo.cardLast4}</span>
                  </p>
                  <Button onClick={this.openModal}>Edit</Button>
                </>
              )}
              {!this.props.paymentInfo && (
                <>
                  <p>You don't have any payment method set up. Start by adding one.</p>
                  <Button onClick={this.openModal}>Add</Button>
                </>
              )}
            </Segment>
          </div>
        </div>
      </>
    );
  }
}
