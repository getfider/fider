import { useCache, useScript } from "@fider/hooks"
import { actions } from "@fider/services"
import { useEffect, useState } from "react"

interface UsePaddleParams {
  isSandbox: boolean
  vendorId: string
  planId: string
}

export function usePaddle(params: UsePaddleParams) {
  const status = useScript("https://cdn.paddle.com/paddle/paddle.js")
  const [price, setPrice] = useCache("price", "$30")
  const [isReady, setIsReady] = useState(false)

  useEffect(() => {
    if (status !== "ready" || !params) return

    if (params.isSandbox) {
      window.Paddle.Environment.set("sandbox")
    }

    const vendor = parseInt(params.vendorId, 10)
    window.Paddle.Setup({ vendor })
    setIsReady(true)

    const id = parseInt(params.planId, 10)
    window.Paddle.Product.Prices(id, (resp) => {
      setPrice(resp.price.net.replace(/\.00/g, ""))
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

  const openCheckoutUrl = async () => {
    const result = await actions.generateCheckoutLink()
    if (result.ok) {
      openUrl(result.data.url)
    }
  }

  return { isReady, price, openUrl, openCheckoutUrl }
}
