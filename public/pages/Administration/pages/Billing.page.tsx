import "./Billing.page.scss";

import React from "react";

import { FaFileInvoice } from "react-icons/fa";
import { AdminBasePage } from "../components/AdminBasePage";
import { Loader } from "@fider/components";
import { PaymentInfo } from "@fider/models";
import { Fider } from "@fider/services";
import { StripeProvider, Elements } from "react-stripe-elements";
import BillingSourceForm from "../components/BillingSourceForm";

interface BillingPageProps {
  paymentInfo?: PaymentInfo;
}

interface BillingPageState {
  stripe: stripe.Stripe | null;
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
      stripe: null
    };
  }

  public componentDidMount() {
    const script = document.createElement("script");
    script.src = "https://js.stripe.com/v3/";
    script.onload = () => {
      this.setState({
        stripe: Stripe(Fider.settings.stripePublicKey!)
      });
    };
    document.body.appendChild(script);
  }

  public content() {
    if (!this.state.stripe) {
      return <Loader />;
    }

    return (
      <>
        <StripeProvider stripe={this.state.stripe}>
          <Elements>
            <BillingSourceForm paymentInfo={this.props.paymentInfo} />
          </Elements>
        </StripeProvider>
      </>
    );
  }
}
