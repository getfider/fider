import React from "react";
import { injectStripe, CardElement } from "react-stripe-elements";
import { Input, Field, Button, Form, Select, SelectOption } from "@fider/components";
import { Failure, actions } from "@fider/services";
import { PaymentInfo } from "@fider/models";

type PatchedTokenResponse = stripe.TokenResponse & {
  error?: { decline_code?: string };
};

interface StripeProps {
  createToken(options?: stripe.TokenOptions): Promise<PatchedTokenResponse>;
  createSource(sourceData?: stripe.SourceOptions): Promise<stripe.SourceResponse>;
  paymentRequest: stripe.Stripe["paymentRequest"];
}

interface BillingSourceFormProps {
  paymentInfo?: PaymentInfo;
  stripe?: StripeProps;
}

interface BillingSourceFormState {
  name: string;
  email: string;
  addressLine1: string;
  addressLine2: string;
  addressCity: string;
  addressState: string;
  addressPostalCode: string;
  addressCountry: string;
  error?: Failure;
}

class BillingSourceForm extends React.Component<BillingSourceFormProps, BillingSourceFormState> {
  constructor(props: BillingSourceFormProps) {
    super(props);
    this.state = {
      name: this.props.paymentInfo ? this.props.paymentInfo.name : "",
      email: this.props.paymentInfo ? this.props.paymentInfo.email : "",
      addressLine1: this.props.paymentInfo ? this.props.paymentInfo.addressLine1 : "",
      addressLine2: this.props.paymentInfo ? this.props.paymentInfo.addressLine2 : "",
      addressCity: this.props.paymentInfo ? this.props.paymentInfo.addressCity : "",
      addressState: this.props.paymentInfo ? this.props.paymentInfo.addressState : "",
      addressPostalCode: this.props.paymentInfo ? this.props.paymentInfo.addressPostalCode : "",
      addressCountry: this.props.paymentInfo ? this.props.paymentInfo.addressCountry : ""
    };
  }

  public handleSubmit = async (ev: React.FormEvent<HTMLFormElement>) => {
    ev.preventDefault();

    if (this.props.paymentInfo) {
      const response = await actions.updatePaymentInfo({
        ...this.state
      });

      this.setState({
        error: response.error
      });

      return;
    }

    if (this.props.stripe) {
      const result = await this.props.stripe.createToken({
        name: this.state.name,
        address_line1: this.state.addressLine1,
        address_line2: this.state.addressLine2,
        address_city: this.state.addressCity,
        address_state: this.state.addressState,
        address_zip: this.state.addressPostalCode,
        address_country: this.state.addressCountry
      });

      if (result.token) {
        const response = await actions.updatePaymentInfo({
          ...this.state,
          card: {
            type: result.token.type,
            token: result.token.id,
            country: result.token.card ? result.token.card.country : ""
          }
        });

        this.setState({
          error: response.error
        });
      } else if (result.error) {
        this.setState({
          error: {
            errors: [
              {
                field: "card",
                message: result.error.message!
              }
            ]
          }
        });
      }
    }
  };

  private setName = (name: string) => {
    this.setState({ name });
  };

  private setEmail = (email: string) => {
    this.setState({ email });
  };

  private setAddressLine1 = (addressLine1: string) => {
    this.setState({ addressLine1 });
  };

  private setAddressLine2 = (addressLine2: string) => {
    this.setState({ addressLine2 });
  };

  private setAddressCity = (addressCity: string) => {
    this.setState({ addressCity });
  };

  private setAddressState = (addressState: string) => {
    this.setState({ addressState });
  };

  private setAddressPostalCode = (addressPostalCode: string) => {
    this.setState({ addressPostalCode });
  };

  private setAddressCountry = (option: SelectOption | undefined) => {
    if (option) {
      this.setState({ addressCountry: option.value });
    }
  };

  public render() {
    return (
      <Form error={this.state.error} onSubmit={this.handleSubmit}>
        <div className="row">
          <div className="col-md-12">
            <Input label="Name" field="name" value={this.state.name} onChange={this.setName} />
          </div>
          <div className="col-md-12">
            <Input label="Email" field="email" value={this.state.email} onChange={this.setEmail} />
          </div>
          <div className="col-md-6">
            <Input
              label="Address Line 1"
              value={this.state.addressLine1}
              field="addressLine1"
              onChange={this.setAddressLine1}
            />
          </div>
          <div className="col-md-6">
            <Input
              label="Address Line 2"
              value={this.state.addressLine2}
              field="addressLine2"
              onChange={this.setAddressLine2}
            />
          </div>
          <div className="col-md-3">
            <Input label="City" field="addressCity" value={this.state.addressCity} onChange={this.setAddressCity} />
          </div>
          <div className="col-md-3">
            <Input
              label="State / Region"
              field="addressState"
              value={this.state.addressState}
              onChange={this.setAddressState}
            />
          </div>
          <div className="col-md-3">
            <Input
              label="Postal Code"
              field="addressPostalCode"
              value={this.state.addressPostalCode}
              onChange={this.setAddressPostalCode}
            />
          </div>
          <div className="col-md-3">
            <Select
              label="Country"
              field="addressCountry"
              onChange={this.setAddressCountry}
              defaultValue={this.state.addressCountry}
              options={[
                { value: "", label: "" },
                { value: "AF", label: "Afghanistan" },
                { value: "AX", label: "Åland Islands" },
                { value: "AL", label: "Albania" },
                { value: "DZ", label: "Algeria" },
                { value: "AS", label: "American Samoa" },
                { value: "AD", label: "Andorra" },
                { value: "AO", label: "Angola" },
                { value: "AI", label: "Anguilla" },
                { value: "AQ", label: "Antarctica" },
                { value: "AG", label: "Antigua and Barbuda" },
                { value: "AR", label: "Argentina" },
                { value: "AM", label: "Armenia" },
                { value: "AW", label: "Aruba" },
                { value: "AU", label: "Australia" },
                { value: "AT", label: "Austria" },
                { value: "AZ", label: "Azerbaijan" },
                { value: "BS", label: "Bahamas" },
                { value: "BH", label: "Bahrain" },
                { value: "BD", label: "Bangladesh" },
                { value: "BB", label: "Barbados" },
                { value: "BY", label: "Belarus" },
                { value: "BE", label: "Belgium" },
                { value: "BZ", label: "Belize" },
                { value: "BJ", label: "Benin" },
                { value: "BM", label: "Bermuda" },
                { value: "BT", label: "Bhutan" },
                { value: "BO", label: "Bolivia, Plurinational State of" },
                { value: "BQ", label: "Bonaire, Sint Eustatius and Saba" },
                { value: "BA", label: "Bosnia and Herzegovina" },
                { value: "BW", label: "Botswana" },
                { value: "BV", label: "Bouvet Island" },
                { value: "BR", label: "Brazil" },
                { value: "IO", label: "British Indian Ocean Territory" },
                { value: "BN", label: "Brunei Darussalam" },
                { value: "BG", label: "Bulgaria" },
                { value: "BF", label: "Burkina Faso" },
                { value: "BI", label: "Burundi" },
                { value: "KH", label: "Cambodia" },
                { value: "CM", label: "Cameroon" },
                { value: "CA", label: "Canada" },
                { value: "CV", label: "Cape Verde" },
                { value: "KY", label: "Cayman Islands" },
                { value: "CF", label: "Central African Republic" },
                { value: "TD", label: "Chad" },
                { value: "CL", label: "Chile" },
                { value: "CN", label: "China" },
                { value: "CX", label: "Christmas Island" },
                { value: "CC", label: "Cocos (Keeling) Islands" },
                { value: "CO", label: "Colombia" },
                { value: "KM", label: "Comoros" },
                { value: "CG", label: "Congo" },
                { value: "CD", label: "Congo, the Democratic Republic of the" },
                { value: "CK", label: "Cook Islands" },
                { value: "CR", label: "Costa Rica" },
                { value: "CI", label: "Côte d'Ivoire" },
                { value: "HR", label: "Croatia" },
                { value: "CU", label: "Cuba" },
                { value: "CW", label: "Curaçao" },
                { value: "CY", label: "Cyprus" },
                { value: "CZ", label: "Czech Republic" },
                { value: "DK", label: "Denmark" },
                { value: "DJ", label: "Djibouti" },
                { value: "DM", label: "Dominica" },
                { value: "DO", label: "Dominican Republic" },
                { value: "EC", label: "Ecuador" },
                { value: "EG", label: "Egypt" },
                { value: "SV", label: "El Salvador" },
                { value: "GQ", label: "Equatorial Guinea" },
                { value: "ER", label: "Eritrea" },
                { value: "EE", label: "Estonia" },
                { value: "ET", label: "Ethiopia" },
                { value: "FK", label: "Falkland Islands (Malvinas)" },
                { value: "FO", label: "Faroe Islands" },
                { value: "FJ", label: "Fiji" },
                { value: "FI", label: "Finland" },
                { value: "FR", label: "France" },
                { value: "GF", label: "French Guiana" },
                { value: "PF", label: "French Polynesia" },
                { value: "TF", label: "French Southern Territories" },
                { value: "GA", label: "Gabon" },
                { value: "GM", label: "Gambia" },
                { value: "GE", label: "Georgia" },
                { value: "DE", label: "Germany" },
                { value: "GH", label: "Ghana" },
                { value: "GI", label: "Gibraltar" },
                { value: "GR", label: "Greece" },
                { value: "GL", label: "Greenland" },
                { value: "GD", label: "Grenada" },
                { value: "GP", label: "Guadeloupe" },
                { value: "GU", label: "Guam" },
                { value: "GT", label: "Guatemala" },
                { value: "GG", label: "Guernsey" },
                { value: "GN", label: "Guinea" },
                { value: "GW", label: "Guinea-Bissau" },
                { value: "GY", label: "Guyana" },
                { value: "HT", label: "Haiti" },
                { value: "HM", label: "Heard Island and McDonald Islands" },
                { value: "VA", label: "Holy See (Vatican City State)" },
                { value: "HN", label: "Honduras" },
                { value: "HK", label: "Hong Kong" },
                { value: "HU", label: "Hungary" },
                { value: "IS", label: "Iceland" },
                { value: "IN", label: "India" },
                { value: "ID", label: "Indonesia" },
                { value: "IR", label: "Iran, Islamic Republic of" },
                { value: "IQ", label: "Iraq" },
                { value: "IE", label: "Ireland" },
                { value: "IM", label: "Isle of Man" },
                { value: "IL", label: "Israel" },
                { value: "IT", label: "Italy" },
                { value: "JM", label: "Jamaica" },
                { value: "JP", label: "Japan" },
                { value: "JE", label: "Jersey" },
                { value: "JO", label: "Jordan" },
                { value: "KZ", label: "Kazakhstan" },
                { value: "KE", label: "Kenya" },
                { value: "KI", label: "Kiribati" },
                { value: "KP", label: "Korea, Democratic People's Republic of" },
                { value: "KR", label: "Korea, Republic of" },
                { value: "KW", label: "Kuwait" },
                { value: "KG", label: "Kyrgyzstan" },
                { value: "LA", label: "Lao People's Democratic Republic" },
                { value: "LV", label: "Latvia" },
                { value: "LB", label: "Lebanon" },
                { value: "LS", label: "Lesotho" },
                { value: "LR", label: "Liberia" },
                { value: "LY", label: "Libya" },
                { value: "LI", label: "Liechtenstein" },
                { value: "LT", label: "Lithuania" },
                { value: "LU", label: "Luxembourg" },
                { value: "MO", label: "Macao" },
                { value: "MK", label: "Macedonia, the former Yugoslav Republic of" },
                { value: "MG", label: "Madagascar" },
                { value: "MW", label: "Malawi" },
                { value: "MY", label: "Malaysia" },
                { value: "MV", label: "Maldives" },
                { value: "ML", label: "Mali" },
                { value: "MT", label: "Malta" },
                { value: "MH", label: "Marshall Islands" },
                { value: "MQ", label: "Martinique" },
                { value: "MR", label: "Mauritania" },
                { value: "MU", label: "Mauritius" },
                { value: "YT", label: "Mayotte" },
                { value: "MX", label: "Mexico" },
                { value: "FM", label: "Micronesia, Federated States of" },
                { value: "MD", label: "Moldova, Republic of" },
                { value: "MC", label: "Monaco" },
                { value: "MN", label: "Mongolia" },
                { value: "ME", label: "Montenegro" },
                { value: "MS", label: "Montserrat" },
                { value: "MA", label: "Morocco" },
                { value: "MZ", label: "Mozambique" },
                { value: "MM", label: "Myanmar" },
                { value: "NA", label: "Namibia" },
                { value: "NR", label: "Nauru" },
                { value: "NP", label: "Nepal" },
                { value: "NL", label: "Netherlands" },
                { value: "NC", label: "New Caledonia" },
                { value: "NZ", label: "New Zealand" },
                { value: "NI", label: "Nicaragua" },
                { value: "NE", label: "Niger" },
                { value: "NG", label: "Nigeria" },
                { value: "NU", label: "Niue" },
                { value: "NF", label: "Norfolk Island" },
                { value: "MP", label: "Northern Mariana Islands" },
                { value: "NO", label: "Norway" },
                { value: "OM", label: "Oman" },
                { value: "PK", label: "Pakistan" },
                { value: "PW", label: "Palau" },
                { value: "PS", label: "Palestinian Territory, Occupied" },
                { value: "PA", label: "Panama" },
                { value: "PG", label: "Papua New Guinea" },
                { value: "PY", label: "Paraguay" },
                { value: "PE", label: "Peru" },
                { value: "PH", label: "Philippines" },
                { value: "PN", label: "Pitcairn" },
                { value: "PL", label: "Poland" },
                { value: "PT", label: "Portugal" },
                { value: "PR", label: "Puerto Rico" },
                { value: "QA", label: "Qatar" },
                { value: "RE", label: "Réunion" },
                { value: "RO", label: "Romania" },
                { value: "RU", label: "Russian Federation" },
                { value: "RW", label: "Rwanda" },
                { value: "BL", label: "Saint Barthélemy" },
                { value: "SH", label: "Saint Helena, Ascension and Tristan da Cunha" },
                { value: "KN", label: "Saint Kitts and Nevis" },
                { value: "LC", label: "Saint Lucia" },
                { value: "MF", label: "Saint Martin (French part)" },
                { value: "PM", label: "Saint Pierre and Miquelon" },
                { value: "VC", label: "Saint Vincent and the Grenadines" },
                { value: "WS", label: "Samoa" },
                { value: "SM", label: "San Marino" },
                { value: "ST", label: "Sao Tome and Principe" },
                { value: "SA", label: "Saudi Arabia" },
                { value: "SN", label: "Senegal" },
                { value: "RS", label: "Serbia" },
                { value: "SC", label: "Seychelles" },
                { value: "SL", label: "Sierra Leone" },
                { value: "SG", label: "Singapore" },
                { value: "SX", label: "Sint Maarten (Dutch part)" },
                { value: "SK", label: "Slovakia" },
                { value: "SI", label: "Slovenia" },
                { value: "SB", label: "Solomon Islands" },
                { value: "SO", label: "Somalia" },
                { value: "ZA", label: "South Africa" },
                { value: "GS", label: "South Georgia and the South Sandwich Islands" },
                { value: "SS", label: "South Sudan" },
                { value: "ES", label: "Spain" },
                { value: "LK", label: "Sri Lanka" },
                { value: "SD", label: "Sudan" },
                { value: "SR", label: "Suriname" },
                { value: "SJ", label: "Svalbard and Jan Mayen" },
                { value: "SZ", label: "Swaziland" },
                { value: "SE", label: "Sweden" },
                { value: "CH", label: "Switzerland" },
                { value: "SY", label: "Syrian Arab Republic" },
                { value: "TW", label: "Taiwan, Province of China" },
                { value: "TJ", label: "Tajikistan" },
                { value: "TZ", label: "Tanzania, United Republic of" },
                { value: "TH", label: "Thailand" },
                { value: "TL", label: "Timor-Leste" },
                { value: "TG", label: "Togo" },
                { value: "TK", label: "Tokelau" },
                { value: "TO", label: "Tonga" },
                { value: "TT", label: "Trinidad and Tobago" },
                { value: "TN", label: "Tunisia" },
                { value: "TR", label: "Turkey" },
                { value: "TM", label: "Turkmenistan" },
                { value: "TC", label: "Turks and Caicos Islands" },
                { value: "TV", label: "Tuvalu" },
                { value: "UG", label: "Uganda" },
                { value: "UA", label: "Ukraine" },
                { value: "AE", label: "United Arab Emirates" },
                { value: "GB", label: "United Kingdom" },
                { value: "US", label: "United States" },
                { value: "UM", label: "United States Minor Outlying Islands" },
                { value: "UY", label: "Uruguay" },
                { value: "UZ", label: "Uzbekistan" },
                { value: "VU", label: "Vanuatu" },
                { value: "VE", label: "Venezuela, Bolivarian Republic of" },
                { value: "VN", label: "Viet Nam" },
                { value: "VG", label: "Virgin Islands, British" },
                { value: "VI", label: "Virgin Islands, U.S." },
                { value: "WF", label: "Wallis and Futuna" },
                { value: "EH", label: "Western Sahara" },
                { value: "YE", label: "Yemen" },
                { value: "ZM", label: "Zambia" },
                { value: "ZW", label: "Zimbabwe" }
              ]}
            />
          </div>
          {!this.props.paymentInfo && (
            <div className="col-md-12">
              <Field label="Card" field="card">
                <CardElement />
              </Field>
            </div>
          )}
        </div>
        <Button type="submit" color="positive">
          Save
        </Button>
      </Form>
    );
  }
}

export default injectStripe(BillingSourceForm);
