import { http, Result } from "../http"

interface CheckoutPageLink {
  url: string
}

export const generateCheckoutLink = async (): Promise<Result<CheckoutPageLink>> => {
  return await http.post("/_api/billing/checkout-link")
}
