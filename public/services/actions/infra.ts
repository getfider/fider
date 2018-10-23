import { http, Result, navigator, analytics } from "@fider/services";

export const logError = async (message: string, err?: Error): Promise<Result> => {
  const data = {
    url: navigator.url(),
    stack: err ? err.stack : "<not available>"
  };

  return http.post("/_api/log-error", { message, data }).then(response => {
    analytics.error(err);
    return response;
  });
};
