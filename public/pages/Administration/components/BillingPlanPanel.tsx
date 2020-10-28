import React from "react";
import { Segment, Button, Moment, Modal, ButtonClickEvent } from "@fider/components";
import { BillingPlan, InvoiceDue } from "@fider/models";
import { Fider, actions, classSet, currencySymbol } from "@fider/services";
import { useFider } from "@fider/hooks";

interface BillingPlanOptionProps {
  tenantUserCount: number;
  disabled: boolean;
  plan: BillingPlan;
  currentPlan?: BillingPlan;
  onSubscribe: (plan: BillingPlan) => Promise<void>;
  onCancel: (plan: BillingPlan) => Promise<void>;
}

const BillingPlanOption = (props: BillingPlanOptionProps) => {
  const fider = useFider();

  const billing = fider.session.tenant.billing!;
  const isSelected = billing.stripePlanID === props.plan.id && !billing.subscriptionEndsAt;
  const className = classSet({ "l-plan": true, selected: isSelected });

  return (
    <div className="col-md-4">
      <Segment className={className}>
        <p className="l-title">{props.plan.name}</p>
        <p className="l-description">{props.plan.description}</p>
        <p>
          <span className="l-currency">{currencySymbol(props.plan.currency)}</span>
          <span className="l-price">{props.plan.price / 100}</span>
          <span className="l-interval">/{props.plan.interval}</span>
        </p>
        {isSelected && (
          <>
            <p>
              <Button disabled={props.disabled} color="danger" onClick={props.onCancel.bind(undefined, props.plan)}>
                Cancel
              </Button>
            </p>
          </>
        )}
        {!isSelected && (
          <>
            <p>
              <Button
                disabled={props.disabled || (props.plan.maxUsers > 0 && props.tenantUserCount > props.plan.maxUsers)}
                onClick={props.onSubscribe.bind(undefined, props.plan)}
              >
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
  invoiceDue?: InvoiceDue;
  tenantUserCount: number;
  disabled: boolean;
  plans: BillingPlan[];
}

interface BillingPlanPanelState {
  confirmPlan?: BillingPlan;
  action?: "" | "subscribe" | "cancel";
}

export class BillingPlanPanel extends React.Component<BillingPlanPanelProps, BillingPlanPanelState> {
  constructor(props: BillingPlanPanelProps) {
    super(props);
    this.state = {};
  }

  private onSubscribe = async (plan: BillingPlan) => {
    this.setState({
      confirmPlan: plan,
      action: "subscribe"
    });
  };

  private onCancel = async (plan: BillingPlan) => {
    this.setState({
      confirmPlan: plan,
      action: "cancel"
    });
  };

  private confirm = async (e: ButtonClickEvent) => {
    e.preventEnable();

    if (this.state.action && this.state.confirmPlan) {
      const action = this.state.action === "subscribe" ? actions.billingSubscribe : actions.cancelBillingSubscription;
      const result = await action(this.state.confirmPlan.id);
      if (result.ok) {
        location.reload();
      }
    }
  };

  private closeModal = async () => {
    this.setState({
      action: "",
      confirmPlan: undefined
    });
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
      <>
        <Modal.Window isOpen={!!this.state.action} center={false} onClose={this.closeModal}>
          {this.state.action === "subscribe" && <Modal.Header>Subscribe</Modal.Header>}
          {this.state.action === "cancel" && <Modal.Header>Cancel Subscription</Modal.Header>}
          <Modal.Content>
            {this.state.action === "subscribe" && (
              <>
                <p>
                  You'll be billed a total of{" "}
                  <strong>
                    {currencySymbol(this.state.confirmPlan!.currency)}
                    {this.state.confirmPlan!.price / 100} per {this.state.confirmPlan!.interval}
                  </strong>{" "}
                  on your card.
                </p>
                <ul>
                  <li>You can cancel it at any time.</li>
                  <li>You can upgrade/downgrade it at any time.</li>
                </ul>
              </>
            )}
            {this.state.action === "cancel" && (
              <>
                <p>You're about to cancel your subscription. Please review the following before continuing.</p>
                <ul>
                  <li>Canceling the subscription will pause any further billing on your card.</li>
                  <li>You'll be able to use the service until the end of current period.</li>
                  <li>You can re-subscribe at any time.</li>
                  <li>No refunds will be given.</li>
                </ul>
                <strong>Are you sure?</strong>
              </>
            )}
          </Modal.Content>
          <Modal.Footer>
            <Button
              color={this.state.action === "subscribe" ? "positive" : "danger"}
              size="tiny"
              onClick={this.confirm}
            >
              {this.state.action === "subscribe" && "Confirm"}
              {this.state.action === "cancel" && "Yes"}
            </Button>
            <Button color="cancel" size="tiny" onClick={this.closeModal}>
              {this.state.action === "subscribe" && "Cancel"}
              {this.state.action === "cancel" && "No"}
            </Button>
          </Modal.Footer>
        </Modal.Window>

        <Segment className="l-billing-plans">
          <h4>Plans</h4>
          {currentPlan && (
            <p className="info">
              {billing.subscriptionEndsAt ? (
                <>
                  Your <strong>{currentPlan.name}</strong> subscription ends at{" "}
                  <strong>
                    <Moment date={billing.subscriptionEndsAt} useRelative={false} />
                  </strong>
                  . Subscribe to a new plan and avoid a service interruption.
                </>
              ) : this.props.invoiceDue ? (
                <>
                  Your upcoming invoice of{" "}
                  <strong>
                    {currencySymbol(this.props.invoiceDue.currency)}
                    {this.props.invoiceDue.amountDue / 100}
                  </strong>{" "}
                  is due on{" "}
                  <strong>
                    <Moment date={this.props.invoiceDue.dueDate} useRelative={false} />
                  </strong>
                  .
                </>
              ) : (
                <></>
              )}
            </p>
          )}
          {!billing.stripePlanID && (
            <p className="info">
              You don't have any active subscription.
              {!trialExpired && (
                <>
                  Your trial period ends at{" "}
                  <strong>
                    <Moment date={billing.trialEndsAt} useRelative={false} />
                  </strong>
                  . Subscribe to a plan and avoid a service interruption.
                </>
              )}
            </p>
          )}
          <div className="row">
            {this.props.plans.map(x => (
              <BillingPlanOption
                key={x.id}
                plan={x}
                tenantUserCount={this.props.tenantUserCount}
                currentPlan={currentPlan}
                disabled={this.props.disabled}
                onSubscribe={this.onSubscribe}
                onCancel={this.onCancel}
              />
            ))}
          </div>
          <div className="row">
            <div className="col-md-12">
              <p className="info">
                You have <strong>{this.props.tenantUserCount}</strong> tracked users.
              </p>
            </div>
          </div>
        </Segment>
      </>
    );
  }
}
