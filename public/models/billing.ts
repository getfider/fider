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
  vatNumber: string;
}

export interface BillingPlan {
  id: string;
  name: string;
  description: string;
  maxUsers: number;
  price: number;
  interval: "month" | "year";
}
