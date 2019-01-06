export interface PaymentInfo {
  cardBrand: string;
  cardLast4: string;
  cardExpMonth: number;
  cardExpYear: number;
  addressCity: string;
  addressCountry: string;
  name: string;
  email: string;
  addressLine1: string;
  addressLine2: string;
  addressState: string;
  addressPostalCode: string;
}

export interface BillingPlan {
  id: string;
  name: string;
  description: string;
  tier: number;
  price: number;
  interval: "month" | "year";
}
