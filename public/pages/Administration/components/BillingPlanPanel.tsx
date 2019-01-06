import "./BillingPlanPanel.scss";

import React from "react";
import { Segment, Button, Moment } from "@fider/components";
import { BillingPlan } from "@fider/models";
import { Fider, actions } from "@fider/services";

interface BillingPlanOptionProps {
  disabled: boolean;
  plan: BillingPlan;
  currentPlan?: BillingPlan;
  onSubscribe: (planID: string) => Promise<void>;
  onCancel: (planID: string) => Promise<void>;
}

const BillingPlanOption = (props: BillingPlanOptionProps) => {
  const billing = Fider.session.tenant.billing!;
  return (
    <div className="col-md-4">
      <Segment className="l-plan">
        <p className="l-title">{props.plan.name}</p>
        <p className="l-description">{props.plan.description}</p>
        <p>
          <span className="l-dollar">$</span>
          <span className="l-price">{props.plan.price / 100}</span>
          <span className="l-interval">/{props.plan.interval}</span>
        </p>
        {billing.stripePlanID === props.plan.id && !billing.subscriptionEndsAt && (
          <>
            <p>
              <Button disabled={props.disabled} color="danger" onClick={props.onCancel.bind(undefined, props.plan.id)}>
                Cancel
              </Button>
            </p>
            <p className="info">You are subscribed to this plan.</p>
          </>
        )}
        {(billing.stripePlanID !== props.plan.id || !!billing.subscriptionEndsAt) && (
          <>
            <p>
              <Button disabled={props.disabled} onClick={props.onSubscribe.bind(undefined, props.plan.id)}>
                Subscribe
              </Button>
            </p>
          </>
        )}
      </Segment>
    </div>
  );
};

interface BillingPlanPanelProps {
  disabled: boolean;
  plans: BillingPlan[];
}

interface BillingPlanPanelState {
  showConfirmation: boolean;
}

export class BillingPlanPanel extends React.Component<BillingPlanPanelProps, BillingPlanPanelState> {
  private subscribe = async (planID: string) => {
    const result = await actions.billingSubscribe(planID);
    if (result.ok) {
      location.reload();
    }
  };

  private cancel = async (planID: string) => {
    const result = await actions.cancelBillingSubscription(planID);
    if (result.ok) {
      location.reload();
    }
  };

  private getCurrentPlan(): BillingPlan | undefined {
    const filtered = this.props.plans.filter(x => x.id === Fider.session.tenant.billing!.stripePlanID);
    if (filtered.length > 0) {
      return filtered[0];
    }
  }

  public render() {
    const billing = Fider.session.tenant.billing!;
    const currentPlan = this.getCurrentPlan();
    const trialExpired = new Date(billing.trialEndsAt) <= new Date();

    return (
      <Segment className="l-billing-plans">
        <h4>Plans</h4>
        {!billing.stripePlanID && (
          <p className="info">
            You don't have any active subscription.
            {!trialExpired && (
              <>
                Subscribe to it before end of your trial:{" "}
                <strong>
                  <Moment date={billing.trialEndsAt} format="full" />
                </strong>
              </>
            )}
          </p>
        )}
        {currentPlan && !!billing.subscriptionEndsAt && (
          <p className="info">
            Your <strong>{currentPlan.name}</strong> ends on{" "}
            <strong>
              <Moment date={billing.subscriptionEndsAt} format="full" />
            </strong>
            . Subscribe to a new plan to avoid a service interruption.
          </p>
        )}
        <div className="row">
          {this.props.plans.map(x => (
            <BillingPlanOption
              key={x.id}
              plan={x}
              currentPlan={currentPlan}
              disabled={this.props.disabled}
              onSubscribe={this.subscribe}
              onCancel={this.cancel}
            />
          ))}
        </div>
      </Segment>
    );
  }
}
