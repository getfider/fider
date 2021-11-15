import { http, Result, navigator, analytics } from "@fider/services"

const ignoreErrors = [
  "http://gj.track.uc.cn/collect", // CSP error: UC Browser tries to use sendBeacon to this domain, which is blocked by our CSP rule
  "null is not an object (evaluating 'c.sheet.insertRule')", // CSP error: UC Browser throws this error even when page is loaded sucessfully
  "Refused to evaluate a string as JavaScript because 'unsafe-eval'", // CSP error: usually thrown because of bad Chrome Extensions
  "vid_mate_check is not defined", // CSP error: thrown by VidMate, an Android Browser
  "SecurityError: Failed to read the 'cssRules' property from 'CSSStyleSheet': Cannot access rules", // CSP error: usually thrown because of bad Chrome Extensions
]

export const logError = async (message: string, err?: Error): Promise<Result | undefined> => {
  for (const pattern of ignoreErrors) {
    if (message.indexOf(pattern) >= 0) {
      return
    }
  }

  const data = {
    url: navigator.url(),
    stack: err ? err.stack : "<not available>",
  }

  try {
    analytics.error(err)
    return await http.post("/_api/log-error", { message, data })
  } catch (err) {
    console.error(err)
  }
}
