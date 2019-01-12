import { http, Result } from "@fider/services";
import { BillingPlan } from "@fider/models";

interface UpdatePaymentInfoRequest {
  name: string;
  card?: {
    type: string;
    token: string;
    country: string;
  };
  addressLine1: string;
  addressLine2: string;
  addressCity: string;
  addressState: string;
  addressPostalCode: string;
  addressCountry: string;
  vatNumber: string;
}

export const listBillingPlans = async (countryCode: string): Promise<Result<BillingPlan[]>> => {
  return http.get(`/_api/admin/billing/plans/${countryCode}`);
};

export const updatePaymentInfo = async (request: UpdatePaymentInfoRequest): Promise<Result> => {
  return http.post("/_api/admin/billing/paymentinfo", request).then(http.event("billing", "updatepaymentinfo"));
};

export const billingSubscribe = async (planID: string): Promise<Result> => {
  return http.post(`/_api/admin/billing/subscription/${planID}`).then(http.event("billing", "billingsubscribe"));
};

export const cancelBillingSubscription = async (planID: string): Promise<Result> => {
  return http
    .delete(`/_api/admin/billing/subscription/${planID}`)
    .then(http.event("billing", "cancelbillingsubscription"));
};
