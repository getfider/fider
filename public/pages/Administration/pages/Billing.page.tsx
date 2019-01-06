import "./Billing.page.scss";

import React from "react";

import { FaFileInvoice } from "react-icons/fa";
import { AdminBasePage } from "../components/AdminBasePage";
import { Segment, Button, CardInfo } from "@fider/components";
import { PaymentInfo, BillingPlan } from "@fider/models";
import { Fider } from "@fider/services";
import PaymentInfoModal from "../components/PaymentInfoModal";
import { StripeProvider, Elements } from "react-stripe-elements";

interface BillingPageProps {
  plans: BillingPlan[];
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
          <div className="col-md-12">
            <Segment className="l-payment-info">
              <h4>Payment Info</h4>
              {this.props.paymentInfo && (
                <>
                  <CardInfo
                    expMonth={this.props.paymentInfo.cardExpMonth}
                    expYear={this.props.paymentInfo.cardExpYear}
                    brand={this.props.paymentInfo.cardBrand}
                    last4={this.props.paymentInfo.cardLast4}
                  />
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
          <div className="col-md-12">
            <Segment className="l-billing-plans">
              <h4>Plans</h4>
              <p className="info">
                You don't have any active subscription. Subscribe to it before end of your trial:{" "}
                <strong>{Fider.session.tenant.billing!.trialEndsAt}</strong>
              </p>
              <div className="row">
                {this.props.plans.map(x => (
                  <div key={x.id} className="col-md-4">
                    <Segment className="l-plan">
                      <p className="l-title">{x.name}</p>
                      <p className="l-description">{x.description}</p>
                      <p>
                        <span className="l-dollar">$</span>
                        <span className="l-price">{x.price / 100}</span>/
                        <span className="l-interval">{x.interval}</span>
                      </p>
                      <Button>Subscribe</Button>
                    </Segment>
                  </div>
                ))}
              </div>
            </Segment>
          </div>
        </div>
      </>
    );
  }
}
