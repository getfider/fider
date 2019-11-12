import React from "react";
import "./CardInfo.scss";
import { useFider } from "@fider/hooks";

const visa = require("@fider/assets/images/card-visa.svg");
const diners = require("@fider/assets/images/card-diners.svg");
const americanExpress = require("@fider/assets/images/card-americanexpress.svg");
const discover = require("@fider/assets/images/card-discover.svg");
const jcb = require("@fider/assets/images/card-jcb.svg");
const unknown = require("@fider/assets/images/card-unknown.svg");
const masterCard = require("@fider/assets/images/card-mastercard.svg");

interface CardBrandProps {
  brand: string;
  last4: string;
  expMonth: number;
  expYear: number;
}

export const CardInfo = (props: CardBrandProps) => {
  const fider = useFider();

  return (
    <p className="c-card-info">
      <img src={`${fider.settings.globalAssetsURL}${brandImage(props.brand)}`} alt={props.brand} />
      <span>
        **** **** **** {props.last4}{" "}
        <span className="c-card-info-exp">
          Exp. {props.expMonth}/{props.expYear}
        </span>
      </span>
    </p>
  );
};

const brandImage = (brand: string) => {
  switch (brand) {
    case "Visa":
      return visa;
    case "American Express":
      return americanExpress;
    case "MasterCard":
      return masterCard;
    case "Discover":
      return discover;
    case "JCB":
      return jcb;
    case "Diners Club":
      return diners;
  }
  return unknown;
};
