import dotenv from "dotenv"
dotenv.config()

import { parse as parseURL } from "url"
import * as http from "http"
import * as https from "https"
import { delay } from "./"

const httpGet = (endpoint: string): Promise<any> => {
  const url = parseURL(endpoint)
  return new Promise((resolve) => {
    const opts = {
      host: url.host,
      port: 443,
      method: "GET",
      auth: `api:${process.env.EMAIL_MAILGUN_API}`,
      path: url.path,
    }
    https.get(opts, (res: http.IncomingMessage) => {
      const content: any[] = []
      res.on("data", (chunk) => content.push(chunk))
      res.on("end", () => resolve(JSON.parse(content.join(""))))
    })
  })
}

export const mailgun = {
  getLinkFromLastEmailTo: async (tenant: string, subject: string, to: string): Promise<string> => {
    let messageURL = ""
    let count = 0

    do {
      count++

      let url = `https://api.mailgun.net/v3/${process.env.EMAIL_MAILGUN_DOMAIN}/events?to=${to}&subject=${subject}&event=accepted&limit=1&ascending=no`
      if (tenant !== "login") {
        url += `&tags=tenant:${tenant}`
      }

      const events = await httpGet(url)
      if (events.items.length > 0 && events.items[0].recipient === to) {
        messageURL = events.items[0].storage.url
      } else {
        await delay(500)
      }
    } while (!messageURL && count < 30)

    if (count === 30) {
      throw new Error(`Message not found for ${to}.`)
    }

    const messages = await httpGet(messageURL)

    const matches = /<a\s+(?:[^>]*?\s+)?href=(["'])(.*?)\1/.exec(messages["body-html"])
    if (matches) {
      return matches[2]
    }

    return ""
  },
}
