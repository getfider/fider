import "./Error.page.scss";

import React from "react";
import { TenantLogo } from "@fider/components";
import { useFider } from "@fider/hooks";
import { useTranslation } from "react-i18next";

interface ErrorPageProps {
  error: Error;
  errorInfo: React.ErrorInfo;
  showDetails?: boolean;
}

export const ErrorPage = (props: ErrorPageProps) => {
  const fider = useFider();
  const { t } = useTranslation();
  return (
    <div id="p-error" className="container failure-page">
      <TenantLogo size={100} useFiderIfEmpty={true} />
      <h1> {t("error.title")} </h1>
      <p>{t("error.message")}</p>
      {fider.settings && (
        <span dangerouslySetInnerHTML={{ __html: t("error.back", { url: fider.settings.baseURL }) }} />
      )}
      {props.showDetails && (
        <pre className="error">
          {props.error.toString()}
          {props.errorInfo.componentStack}
        </pre>
      )}
    </div>
  );
};
