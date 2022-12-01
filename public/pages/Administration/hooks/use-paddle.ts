import { useCache, useScript } from "@fider/hooks"
import { actions } from "@fider/services"
import { useEffect, useState } from "react"

interface UsePaddleParams {
  isSandbox: boolean
  vendorId: string
  monthlyPlanId: string
  yearlyPlanId: string
}

export function usePaddle(params: UsePaddleParams) {
  const status = useScript("https://cdn.paddle.com/paddle/paddle.js")
  const [monthlyPrice, setMonthlyPrice] = useCache("monthlyPrice", "$30")
  const [yearlyPrice, setYearlyPrice] = useCache("yearlyPrice", "$300")
  const [isReady, setIsReady] = useState(false)

  useEffect(() => {
    if (status !== "ready" || !params) return

    if (params.isSandbox) {
      window.Paddle.Environment.set("sandbox")
    }

    const vendor = parseInt(params.vendorId, 10)
    window.Paddle.Setup({ vendor })
    setIsReady(true)

    window.Paddle.Product.Prices(parseInt(params.monthlyPlanId, 10), (resp) => {
      setMonthlyPrice(resp.price.net.replace(/\.00/g, ""))
    })
    window.Paddle.Product.Prices(parseInt(params.yearlyPlanId, 10), (resp) => {
      setYearlyPrice(resp.price.net.replace(/\.00/g, ""))
    })
  }, [status])

  const openUrl = (url: string) => {
    window.Paddle.Checkout.open({
      override: url,
      closeCallback: () => {
        location.reload()
      },
    })
  }

  const subscribeMonthly = async () => {
    const result = await actions.generateCheckoutLink(params.monthlyPlanId)
    if (result.ok) {
      openUrl(result.data.url)
    }
  }

  const subscribeYearly = async () => {
    const result = await actions.generateCheckoutLink(params.yearlyPlanId)
    if (result.ok) {
      openUrl(result.data.url)
    }
  }

  return { isReady, monthlyPrice, yearlyPrice, openUrl, subscribeMonthly, subscribeYearly }
}
