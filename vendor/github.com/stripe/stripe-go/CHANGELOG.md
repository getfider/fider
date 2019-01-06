# Changelog

## 55.3.0 - 2019-01-02
* [#742](https://github.com/stripe/stripe-go/pull/742) Add field `FraudType` to `IssuerFraudRecord`

## 55.2.0 - 2018-12-31
* [#741](https://github.com/stripe/stripe-go/pull/741) Add missing parameters `InvoiceNow` and `Prorate` for subscription cancellation

## 55.1.0 - 2018-12-27
* [#743](https://github.com/stripe/stripe-go/pull/743) Add support for `clear_usage` on `SubscriptionItem` deletion

## 55.0.0 - 2018-12-13
* [#739](https://github.com/stripe/stripe-go/pull/739) Use `ApplicationFee` struct for `FeeRefund.Fee` (minor breaking change)

## 54.2.0 - 2018-11-30
* [#734](https://github.com/stripe/stripe-go/pull/734) Put `/v1/` prefix as part of all paths instead of URL

## 54.1.1 - 2018-11-30
* [#733](https://github.com/stripe/stripe-go/pull/733) Fix malformed URL generated for the uploads API when using `NewBackends`

## 54.1.0 - 2018-11-28
* [#730](https://github.com/stripe/stripe-go/pull/730) Add support for the Review resource
* [#731](https://github.com/stripe/stripe-go/pull/731) Add missing properties on the Refund resource

## 54.0.0 - 2018-11-27
* [#721](https://github.com/stripe/stripe-go/pull/721) Add support for `RadarValueList` and `RadarValueListItem`
* [#721](https://github.com/stripe/stripe-go/pull/721) Remove `Closed` and `Forgiven` from `InvoiceParams`
* [#721](https://github.com/stripe/stripe-go/pull/721) Add `PaidOutOfBand` to `InvoicePayParams`

## 53.4.0 - 2018-11-26
* [#728](https://github.com/stripe/stripe-go/pull/728) Add `IssuingCard` to `EphemeralKeyParams`

## 53.3.0 - 2018-11-26
* [#727](https://github.com/stripe/stripe-go/pull/727) Add support for `TransferData` on payment intent create and update

## 53.2.0 - 2018-11-21
* [#725](https://github.com/stripe/stripe-go/pull/725) Improved error deserialization

## 53.1.0 - 2018-11-15
* [#723](https://github.com/stripe/stripe-go/pull/723) Add support for `last_payment_error` on `PaymentIntent`.
* [#724](https://github.com/stripe/stripe-go/pull/724) Add support for `transfer_data[destination]` on `PaymentIntent`.

## 53.0.1 - 2018-11-12
* [#714](https://github.com/stripe/stripe-go/pull/714) Fix bug in retry logic that would cause the client to panic

## 53.0.0 - 2018-11-08
* [#716](https://github.com/stripe/stripe-go/pull/716) Drop support for Go 1.8.
* [#715](https://github.com/stripe/stripe-go/pull/715) Ship changes to the `PaymentIntent` resource to match the final layout.
* [#717](https://github.com/stripe/stripe-go/pull/717) Add support for `flat_amount` on `Plan` tiers.
* [#718](https://github.com/stripe/stripe-go/pull/718) Add support for `supported_transfer_countries` on `CountrySpec`.
* [#720](https://github.com/stripe/stripe-go/pull/720) Add support for `review` on `PaymentIntent`.
* [#707](https://github.com/stripe/stripe-go/pull/707) Add new invoice methods and fixes to the Issuing Cardholder resource (multiple breaking changes)
    * Move to API version 2018-11-08.
    * Add support for new API methods, properties and parameters for `Invoice`.
    * Add support for `default_source` on `Subscription` and `Invoice`.

## 52.1.0 - 2018-10-31
* [#705](https://github.com/stripe/stripe-go/pull/705) Add support for the `Person` resource
* [#706](https://github.com/stripe/stripe-go/pull/706) Add support for the `WebhookEndpoint` resource

## 52.0.0 - 2018-10-29
* [#711](https://github.com/stripe/stripe-go/pull/711) Set `Request.GetBody` when making requests
* [#711](https://github.com/stripe/stripe-go/pull/711) Drop support for Go 1.7 (hasn't been supported by Go core since the release of Go 1.9 in August 2017)

## 51.4.0 - 2018-10-19
* [#708](https://github.com/stripe/stripe-go/pull/708) Add Stripe Terminal endpoints to master to `client.API` 

## 51.3.0 - 2018-10-09
* [#704](https://github.com/stripe/stripe-go/pull/704) Add support for `subscription_cancel_at_period_end` on the Upcoming Invoice API.

## 51.2.0 - 2018-10-09
* [#702](https://github.com/stripe/stripe-go/pull/702) Add support for `delivery_success` filter when listing Events.

## 51.1.0 - 2018-10-03
* [#700](https://github.com/stripe/stripe-go/pull/700) Add support for `on_behalf_of` on Subscription and Charge resources.

## 51.0.0 - 2018-09-27
* [#698](https://github.com/stripe/stripe-go/pull/698) Move to API version 2018-09-24
    * Rename `FileUpload` to `File` (and all `FileUpload*` structs to `File*`)
	* Fix file links client

## 50.0.0 - 2018-09-24
* [#695](https://github.com/stripe/stripe-go/pull/695) Rename `Transaction` to `DisputedTransaction` in `IssuingDisputeParams` (minor breaking change)
* [#695](https://github.com/stripe/stripe-go/pull/695) Add support for Stripe Terminal

## 49.2.0 - 2018-09-24
* [#697](https://github.com/stripe/stripe-go/pull/697) Fix `number` JSON tag on the `IssuingCardDetails` resource.

## 49.1.0 - 2018-09-11
* [#694](https://github.com/stripe/stripe-go/pull/694) Add `ErrorCodeResourceMissing` error code constant

## 49.0.0 - 2018-09-11
* [#693](https://github.com/stripe/stripe-go/pull/693) Change `Product` under `Plan` from a string to a full `Product` struct pointer (this is a minor breaking change -- upgrade by changing to `plan.Product.ID`)

## 48.3.0 - 2018-09-06
* [#691](https://github.com/stripe/stripe-go/pull/691) Add `InvoicePrefix` to `Customer` and `CustomerParams`

## 48.2.0 - 2018-09-05
* [#690](https://github.com/stripe/stripe-go/pull/690) Add support for reporting resources

## 48.1.0 - 2018-09-05
* [#683](https://github.com/stripe/stripe-go/pull/683) Add `StatusTransitions` filter parameters to `OrderListParams`

## 48.0.0 - 2018-09-05
* [#681](https://github.com/stripe/stripe-go/pull/681) Handle deserialization of `OrderItem` parent into an object if expanded (minor breaking change)

## 47.0.0 - 2018-09-04
* New major version for better compatibility with Go's new module system (no breaking changes)

## 46.1.0 - 2018-09-04
* [#688](https://github.com/stripe/stripe-go/pull/688) Encode `Params` in `AppendToAsSourceOrExternalAccount` (bug fix)
* [#689](https://github.com/stripe/stripe-go/pull/689) Add `go.mod` for the new module system

## 46.0.0 - 2018-09-04
* [#686](https://github.com/stripe/stripe-go/pull/686) Add `Mandate` and `Receiver` to `SourceObjectParams` and change `Date` on `SourceMandateAcceptance` to `int64` (minor breaking change)

## 45.0.0 - 2018-08-30
* [#680](https://github.com/stripe/stripe-go/pull/680) Change `SubscriptionTaxPercent` on `Invoice` from `int64` to `float64` (minor breaking change)

## 44.0.0 - 2018-08-28
* [#678](https://github.com/stripe/stripe-go/pull/678) Allow payment intent capture to take its own parameters

## 43.1.1 - 2018-08-28
* [#675](https://github.com/stripe/stripe-go/pull/675) Fix incorrectly encoded parameter in `UsageRecordSummaryListParams`

## 43.1.0 - 2018-08-28
* [#669](https://github.com/stripe/stripe-go/pull/669) Add `AuthorizationCode` to `Charge`
* [#671](https://github.com/stripe/stripe-go/pull/671) Fix deserialization of `TaxID` on `CustomerTaxInfo`

## 43.0.0 - 2018-08-23
* [#668](https://github.com/stripe/stripe-go/pull/668) Move to API version 2018-08-23
    * Add `TaxInfo` and `TaxInfoVerification` to `Customer`
	* Rename `Amount` to `UnitAmount` on `PlanTierParams`
	* Remove `BusinessVATID` from `Customer`
	* Remove `AtPeriodEnd` from `SubscriptionCancelParams`

## 42.3.0 - 2018-08-23
* [#667](https://github.com/stripe/stripe-go/pull/667) Add `Forgive` to `InvoicePayParams`

## 42.2.0 - 2018-08-22
* [#666](https://github.com/stripe/stripe-go/pull/666) Add `Subscription` to `SubscriptionItem`

## 42.1.0 - 2018-08-22
* [#664](https://github.com/stripe/stripe-go/pull/664) Add `AvailablePayoutMethods` to `Card`

## 42.0.0 - 2018-08-20
* [#663](https://github.com/stripe/stripe-go/pull/663) Add support for usage record summaries and rename `Live` on `IssuerFraudRecord, `SourceTransaction`, and `UsageRecord` to `Livemode` (a minor breaking change)

## 41.0.0 - 2018-08-17
* [#659](https://github.com/stripe/stripe-go/pull/659) Remove mutating Bitcoin receiver API calls (these were no longer functional anyway)
* [#661](https://github.com/stripe/stripe-go/pull/661) Correct `IssuingCardShipping`'s type to `int64`
* [#662](https://github.com/stripe/stripe-go/pull/662) Rename `IssuingCardShipping`'s `Eta` to `ETA`

## 40.2.0 - 2018-08-15
* [#657](https://github.com/stripe/stripe-go/pull/657) Use integer-indexed encoding for all arrays

## 40.1.0 - 2018-08-10
* [#656](https://github.com/stripe/stripe-go/pull/656) Expose new `ValidatePayload` functions for validating incoming payloads without constructing an event

## 40.0.2 - 2018-08-07
* [#652](https://github.com/stripe/stripe-go/pull/652) Change the type of `FileUpload.Links` to `FileLinkList` (this is a bug fix given that the previous type would never have worked)

## 40.0.1 - 2018-08-07
* [#653](https://github.com/stripe/stripe-go/pull/653) All `BackendImplementation`s should sleep by default on retries

## 40.0.0 - 2018-08-06
* [#648](https://github.com/stripe/stripe-go/pull/648) Introduce buffers so a request's body can be read multiple times (this modifies the interface of a few exported internal functions so it's technically breaking, but it will probably not be breaking for most users)
* [#649](https://github.com/stripe/stripe-go/pull/649) Rename `BackendConfiguration` to `BackendImplementation` (likewise, technically breaking, but minor)
* [#650](https://github.com/stripe/stripe-go/pull/650) Export `webhook.ComputeSignature`

## 39.0.0 - 2018-08-04
* [#646](https://github.com/stripe/stripe-go/pull/646) Set request body before every retry (this modifies the interface of a few exported internal functions so it's technically breaking, but it will probably not be breaking for most users)

## 38.2.0 - 2018-08-03
* [#644](https://github.com/stripe/stripe-go/pull/644) Add support for file links
* [#645](https://github.com/stripe/stripe-go/pull/645) Add support for `Cancel` to topups

## 38.1.0 - 2018-08-01
* [#643](https://github.com/stripe/stripe-go/pull/643) Bug fix and various code/logging improvements to retry code

## 38.0.0 - 2018-07-30
* [#641](https://github.com/stripe/stripe-go/pull/641) Minor breaking changes to correct a few naming inconsistencies:
    * `IdentityVerificationDetailsCodeScanIdCountryNotSupported` becomes `IdentityVerificationDetailsCodeScanIDCountryNotSupported`
    * `IdentityVerificationDetailsCodeScanIdTypeNotSupported` becomes `IdentityVerificationDetailsCodeScanIDTypeNotSupported`
    * `BitcoinUri` on `BitcoinReceiver` becomes `BitcoinURI`
    * `NetworkId` on `IssuingAuthorization` becomes `NetworkID`

## 37.0.0 - 2018-07-30
* [#637](https://github.com/stripe/stripe-go/pull/637) Add support for Sigma scheduled query runs
* [#639](https://github.com/stripe/stripe-go/pull/639) Move to API version `2018-07-27` (breaking)
    * Remove `SKUs` from `Product`
    * Subscription creation and update can no longer take a source
    * Change `PercentOff` on coupon struct and params from integer to float
* [#640](https://github.com/stripe/stripe-go/pull/640) Add missing field `Created` to `Account`

## 36.3.0 - 2018-07-27
* [#636](https://github.com/stripe/stripe-go/pull/636) Add `RiskScore` to `ChargeOutcome`

## 36.2.0 - 2018-07-26
* [#635](https://github.com/stripe/stripe-go/pull/635) Add support for Stripe Issuing

## 36.1.2 - 2018-07-24
* [#633](https://github.com/stripe/stripe-go/pull/633) Fix encoding of list params for bank accounts and cards

## 36.1.1 - 2018-07-17
* [#627](https://github.com/stripe/stripe-go/pull/627) Wire an `http.Client` from `NewBackends` through to backends

## 36.1.0 - 2018-07-11
* [#624](https://github.com/stripe/stripe-go/pull/624) Add `AutoAdvance` for `Invoice`

## 36.0.0 - 2018-07-09
* [#606](https://github.com/stripe/stripe-go/pull/606) Add support for payment intents
* [#623](https://github.com/stripe/stripe-go/pull/623) Changed `Payout.Destination` from `string` to `*PayoutDestination` to support expanding (minor breaking change)

## 35.13.0 - 2018-07-06
* [#622](https://github.com/stripe/stripe-go/pull/622) Correct position of `DeclineChargeOn` (it was added accidentally on `LegalEntityParams` when it should have been on `AccountParams`)

## 35.12.0 - 2018-07-05
* [#620](https://github.com/stripe/stripe-go/pull/620) Add support for `Quantity` and `UnitAmount` to `InvoiceItemParams` and `Quantity` to `InvoiceItem`

## 35.11.0 - 2018-07-05
* [#618](https://github.com/stripe/stripe-go/pull/618) Add support for `DeclineChargeOn` to `Account` and `AccountParams`

## 35.10.0 - 2018-07-04
* [#616](https://github.com/stripe/stripe-go/pull/616) Adding missing clients to the `API` struct including a `UsageRecords` entry

## 35.9.0 - 2018-07-03
* [#611](https://github.com/stripe/stripe-go/pull/611) Introduce `GetBackendWithConfig` and make logging configurable per backend

## 35.8.0 - 2018-06-28
* [#607](https://github.com/stripe/stripe-go/pull/607) Add support for `PartnerID` from `stripe.SetAppInfo`

## 35.7.0 - 2018-06-26
* [#604](https://github.com/stripe/stripe-go/pull/604) Add extra parameters `CustomerReference` and `ShippingFromZip` to `ChargeLevel3Params` and `ChargeLevel3`

## 35.6.0 - 2018-06-25
* [#603](https://github.com/stripe/stripe-go/pull/603) Add support for Level III data on charge creation

## 35.5.0 - 2018-06-22
* [#601](https://github.com/stripe/stripe-go/pull/601) Add missing parameters for retrieving an upcoming invoice

## 35.4.0 - 2018-06-21
* [#599](https://github.com/stripe/stripe-go/pull/599) Add `ExchangeRate` to `BalanceTransaction`

## 35.3.0 - 2018-06-20
* [#596](https://github.com/stripe/stripe-go/pull/596) Add `Type` to `ProductListParams` so that products can be listed by type

## 35.2.0 - 2018-06-19
* [#595](https://github.com/stripe/stripe-go/pull/595) Add `Product` to `PlanListParams` so that plans can be listed by product

## 35.1.0 - 2018-06-17
* [#592](https://github.com/stripe/stripe-go/pull/592) Add `Name` field to `Coupon` and `CouponParams`

## 35.0.0 - 2018-06-15
* [#557](https://github.com/stripe/stripe-go/pull/557) Add automatic retries for intermittent errors (enabling using `BackendConfiguration.SetMaxNetworkRetries`)
* [#589](https://github.com/stripe/stripe-go/pull/589) Fix all `Get` methods to support standardized parameter structs + remove some deprecated functions
	* `IssuerFraudRecordListParams` now uses `*string` for `Charge` (set it using `stripe.String` like elsewhere)
	* `event.Get` now takes `stripe.EventParams` instead of `Params` for consistency
	* The `Get` method for `countryspec`, `exchangerate`, `issuerfraudrecord` now take an extra params struct parameter to be consistent and allow setting a connected account (use `stripe.CountrySpecParams`, `stripe.ExchangeRateParams`, and `IssuerFraudRecordParams`)
	* `charge.MarkFraudulent` and `charge.MarkSafe` have been removed; use `charge.Update` instead
	* `charge.CloseDispute` and `charge.UpdateDispute` have been removed; use `dispute.Update` or `dispute.Close` instead
	* `loginlink.New` now properly passes its params struct into its API call

## 34.3.0 - 2018-06-14
* [#587](https://github.com/stripe/stripe-go/pull/587) Use `net/http` constants instead of string literals for HTTP verbs (this is an internal cleanup and should not affect library behavior)

## 34.2.0 - 2018-06-14
* [#581](https://github.com/stripe/stripe-go/pull/581) Push parameter encoding into `BackendConfiguration.Call` (this is an internal cleanup and should not affect library behavior)

## 34.1.0 - 2018-06-13
* [#586](https://github.com/stripe/stripe-go/pull/586) Add `AmountPaid`, `AmountRemaining`, `BillingReason` (including new `InvoiceBillingReason` and constants), and `SubscriptionProrationDate` to `Invoice`

## 34.0.0 - 2018-06-12
* [#585](https://github.com/stripe/stripe-go/pull/585) Remove `File` in favor of `FileUpload`, and consolidating both classes which were already nearly identical except `MIMEType` has been replaced by `Type` (this is technically a breaking change, but quite a small one)

## 33.1.0 - 2018-06-12
* [#578](https://github.com/stripe/stripe-go/pull/578) Improve expansion parsing by not discarding unmarshal errors

## 33.0.0 - 2018-06-11
* [#583](https://github.com/stripe/stripe-go/pull/583) Add new account constants, rename one, and fix `DueBy` (this is technically a breaking change, but quite a small one)

## 32.4.1 - 2018-06-11
* [#582](https://github.com/stripe/stripe-go/pull/582) Fix unmarshaling of `LegalEntity` (specifically when we have `legal_entity[additional_owners][][verification]`) so that it comes out as a struct

## 32.4.0 - 2018-06-07
* [#577](https://github.com/stripe/stripe-go/pull/577) Add `DocumentBack` to account legal entity identity verification parameters and response

## 32.3.0 - 2018-06-07
* [#576](https://github.com/stripe/stripe-go/pull/576) Fix plan transform usage to use `BucketSize` instead of `DivideBy`; note this is technically a breaking API change, but we've released it as a minor because the previous manifestation didn't work

## 32.2.0 - 2018-06-06
* [#571](https://github.com/stripe/stripe-go/pull/571) Add `HostedInvoiceURL` and `InvoicePDF` to `Invoice`
* [#573](https://github.com/stripe/stripe-go/pull/573) Add `FormatURLPath` helper to allow safer URL path building

## 32.1.0 - 2018-06-06
* [#572](https://github.com/stripe/stripe-go/pull/572) Add `Active` to plan parameters and response

## 32.0.1 - 2018-06-06
* [#569](https://github.com/stripe/stripe-go/pull/569) Fix unmarshaling of expanded transaction sources in balance transactions

## 32.0.0 - 2018-06-06
* [#544](https://github.com/stripe/stripe-go/pull/544) **MAJOR** changes that make all fields on parameter structs pointers, and rename many fields on parameter and response structs to be consistent with naming in the REST API; we've written [a migration guide with complete details](https://github.com/stripe/stripe-go/blob/master/v32_migration_guide.md) to help with the upgrade

## 31.0.0 - 2018-06-06
* [#566](https://github.com/stripe/stripe-go/pull/566) Support `DisputeParams` in `dispute.Close`

## 30.8.1 - 2018-05-24
* [#562](https://github.com/stripe/stripe-go/pull/562) Add `go.mod` for vgo support

## 30.8.0 - 2018-05-22
* [#558](https://github.com/stripe/stripe-go/pull/558) Add `SubscriptionItem` to `InvoiceLine`

## 30.7.0 - 2018-05-09
* [#552](https://github.com/stripe/stripe-go/pull/552) Add support for issuer fraud records

## 30.6.1 - 2018-05-04
* [#550](https://github.com/stripe/stripe-go/pull/550) Append standard `Params` as well as card options when encoding `CardParams`

## 30.6.0 - 2018-04-17
* [#546](https://github.com/stripe/stripe-go/pull/546) Add `SubParams.TrialFromPlan` and `SubItemsParams.ClearUsage`

## 30.5.0 - 2018-04-09
* [#543](https://github.com/stripe/stripe-go/pull/543) Support listing orders by customer (add `Customer` to `OrderListParams`)

## 30.4.0 - 2018-04-06
* [#541](https://github.com/stripe/stripe-go/pull/541) Add `Mandate` on `Source` (and associated mandate structs)

## 30.3.0 - 2018-04-02
* [#538](https://github.com/stripe/stripe-go/pull/538) Introduce flexible billing primitives for subscriptions

## 30.2.0 - 2018-03-23
* [#535](https://github.com/stripe/stripe-go/pull/535) Add constant for redirect status `not_required` (`RedirectFlowStatusNotRequired`)

## 30.1.0 - 2018-03-17
* [#534](https://github.com/stripe/stripe-go/pull/534) Add `AmountZero` to `InvoiceItemParams`

## 30.0.0 - 2018-03-14
* [#533](https://github.com/stripe/stripe-go/pull/533) Make `DestPayment` under `Transfer` expandable by changing it from a string to a `Charge`

## 29.3.1 - 2018-03-08
* [#530](https://github.com/stripe/stripe-go/pull/530) Fix mixed up types in `CountrySpec.SupportedBankAccountCurrencies`

## 29.3.0 - 2018-03-01
* [#527](https://github.com/stripe/stripe-go/pull/527) Add `MaidenName`, `PersonalIDNumber`, `PersonalIDNumberProvided` fields to `Owner` struct

## 29.2.0 - 2018-02-26
* [#525](https://github.com/stripe/stripe-go/pull/525) Support shipping carrier and tracking number in orders
* [#526](https://github.com/stripe/stripe-go/pull/526) Fix ignored `commonParams` when returning an order

## 29.1.1 - 2018-02-21
* [#522](https://github.com/stripe/stripe-go/pull/522) Bump API version and fix creating plans with a product

## 29.1.0 - 2018-02-21
* [#520](https://github.com/stripe/stripe-go/pull/520) Add support for topups

## 29.0.1 - 2018-02-16
**WARNING:** Please use 29.1.1 instead.
* [#519](https://github.com/stripe/stripe-go/pull/519) Correct the implementation of `PaymentSource.MarshalJSON` to also handle bank account sources

## 29.0.0 - 2018-02-14
**WARNING:** Please use 29.1.1 instead.
* [#518](https://github.com/stripe/stripe-go/pull/518) Bump API version to 2018-02-06 and add support for Product & Plan API

## 28.12.0 - 2018-02-09
* [#517](https://github.com/stripe/stripe-go/pull/517) Add `BillingCycleAnchor` to `Sub` and `BillingCycleAnchorUnchanged` to `SubParams`

## 28.11.0 - 2018-01-29
* [#516](https://github.com/stripe/stripe-go/pull/516) Add `AmountZero` to `PlanParams` to it's possible to send zero values when creating or updating a plan

## 28.10.1 - 2018-01-18
* [#512](https://github.com/stripe/stripe-go/pull/512) Encode empty values found in maps (like `Meta`)

## 28.10.0 - 2018-01-09
* [#509](https://github.com/stripe/stripe-go/pull/509) Plumb through additional possible errors when unmarshaling polymorphic types (please test your integrations while upgrading)

## 28.9.0 - 2018-01-08
* [#506](https://github.com/stripe/stripe-go/pull/506) Add support for recursing into slices in `event.GetObjValue`

## 28.8.0 - 2017-12-12
* [#500](https://github.com/stripe/stripe-go/pull/500) Support sharing for bank accounts and cards (adds `ID` field to bank account and charge parameters)

## 28.7.0 - 2017-12-05
* [#494](https://github.com/stripe/stripe-go/pull/494) Add `Automatic` to `Payout` struct

## 28.6.1 - 2017-11-02
* [#492](https://github.com/stripe/stripe-go/pull/492) Correct name of user agent header used to send Go version to Stripe's API

## 28.6.0 - 2017-10-31
* [#491](https://github.com/stripe/stripe-go/pull/491) Support for exchange rates APIs

## 28.5.0 - 2017-10-27
* [#488](https://github.com/stripe/stripe-go/pull/488) Support for listing source transactions

## 28.4.2 - 2017-10-25
* [#486](https://github.com/stripe/stripe-go/pull/486) Send the required `object=bank_account` parameter when adding a bank account through an account
* [#487](https://github.com/stripe/stripe-go/pull/487) Make bank account's `account_holder_name` and `account_holder_type` parameters truly optional

## 28.4.1 - 2017-10-24
* [#484](https://github.com/stripe/stripe-go/pull/484) Error early when params not specified for card-related API calls

## 28.4.0 - 2017-10-19
* [#477](https://github.com/stripe/stripe-go/pull/477) Support context on API requests with `Params.Context` and `ListParams.Context`

## 28.3.2 - 2017-10-19
* [#479](https://github.com/stripe/stripe-go/pull/479) Pass token in only one of `external_account` *or* source when appending card

## 28.3.1 - 2017-10-17
* [#476](https://github.com/stripe/stripe-go/pull/476) Make initializing new backends concurrency-safe

## 28.3.0 - 2017-10-10
* [#359](https://github.com/stripe/stripe-go/pull/359) Add support for verify sources (added `Values` on `SourceVerifyParams`)

## 28.2.0 - 2017-10-09
* [#472](https://github.com/stripe/stripe-go/pull/472) Add support for `statement_descriptor` in source objects
* [#473](https://github.com/stripe/stripe-go/pull/473) Add support for detaching sources from customers

## 28.1.0 - 2017-10-05
* [#471](https://github.com/stripe/stripe-go/pull/471) Add support for `RedirectFlow.FailureReason` for sources

## 28.0.1 - 2017-10-03
* [#468](https://github.com/stripe/stripe-go/pull/468) Fix encoding of pointer-based scalars (e.g. `Active *bool` in `Product`)
* [#470](https://github.com/stripe/stripe-go/pull/470) Fix concurrent race in `form` package's encoding caches

## 28.0.0 - 2017-09-27
* [#467](https://github.com/stripe/stripe-go/pull/467) Change `Product.Get` to include `ProductParams` for request metadata
* [#467](https://github.com/stripe/stripe-go/pull/467) Fix sending extra parameters on product and SKU requests

## 27.0.2 - 2017-09-26
* [#465](https://github.com/stripe/stripe-go/pull/465) Fix encoding of `CVC` parameter in `CardParams`

## 27.0.1 - 2017-09-20
* [#461](https://github.com/stripe/stripe-go/pull/461) Fix encoding of `TypeData` under sources

## 27.0.0 - 2017-09-19
* [#458](https://github.com/stripe/stripe-go/pull/458) Remove `ChargeParams.Token` (this seems like it was added accidentally)

## 26.0.0 - 2017-09-17
* Introduce `form` package so it's no longer necessary to build conditional structures to encode parameters -- this may result in parameters that were set but previously not encoded to now be encoded so **PLEASE TEST CAREFULLY WHEN UPGRADING**!
* Alphabetize all struct fields -- this may result in position-based struct initialization to fail if it was being used
* Switch to stripe-mock for testing (test suite now runs completely!)
* Remote Displayer interface and Display implementations
* Add `FraudDetails` to `ChargeParams`
* Remove `FraudReport` from `ChargeParams` (use `FraudDetails` instead)

## 25.2.0 - 2017-09-13
* Add `OnBehalfOf` to charge parameters.
* Add `OnBehalfOf` to subscription parameters.

## 25.1.0 - 2017-09-06
* Use bearer token authentication for API requests

## 25.0.0 - 2017-08-21
* All `Del` methods now take params as second argument (which may be `nil`)
* Product `Delete` has been renamed to `Del` for consistency
* Product `Delete` now returns `(*Product, error)` for consistency
* SKU `Delete` has been renamed to `Del` for consistency
* SKU `Delete` now returns `(*SKU, error)` for consistency

## 24.3.0 - 2017-08-08
* Add `FeeZero` to invoice and `TaxPercentZero` to subscription for zeroing values

## 24.2.0 - 2017-07-25
* Add "range queries" for supported parameters (e.g. `created[gte]=123`)

## 24.1.0 - 2017-07-17
* Add metadata to subscription items

## 24.0.0 - 2017-06-27
	`Pay` on invoice now takes specific pay parameters

## 23.2.1 - 2017-06-26
* Fix bank account retrieval when using a customer ID

## 23.2.0 - 2017-06-26
* Support sharing path while creating a source

## 23.1.0 - 2017-06-26
* Add LoginLinks to client list

## 23.0.0 - 2017-06-23
	plan.Del now takes `stripe.PlanParams` as a second argument

## 22.6.0 - 2017-06-19
* Support for ephemeral keys

## 22.5.0 - 2017-06-15
* Support for checking webhook signatures

## 22.4.1 - 2017-06-15
* Fix returned type of subscription items list
* Note: I meant to release this as 22.3.1, but I'm leaving it as it was released

## 22.3.0 - 2017-06-14
* Fix parameters for subscription items list

## 22.2.0 - 2017-06-13
* Support subscription items when getting upcoming invoice
* Support setting subscription's quantity to zero when getting upcoming invoice

## 22.1.1 - 2017-06-12
* Handle `deleted` parameter when updating subscription items in a subscription

## 22.1.0 - 2017-05-25
* Change `Logger` to a `log.Logger`-like interface so other loggers are usable

## 22.0.0 - 2017-05-25
* Add support for login links
* Add support for new `Type` for accounts
* Make `Event` `Request` (renamed from `Req`) a struct with a new idempotency key
* Rename `Event` `UserID` to `Account`

## 21.5.1 - 2017-05-23
* Fix plan update so `TrialPeriod` parameter is sent

## 21.5.0 - 2017-05-15
* Implement `Get` for `RequestValues`

## 21.4.1 - 2017-05-11
* Pass extra parameters to API calls on bank account deletion

## 21.4.0 - 2017-05-04
* Add `Billing` and `DueDate` filters to invoice listing
* Add `Billing` filter to subscription listing

## 21.3.0 - 2017-05-02
* Add `DetailsCode` to `IdentityVerification`

## 21.2.0 - 2017-04-19
* Send user agent information with `X-Stripe-Client-User-Agent`
* Add `stripe.SetAppInfo` for plugin authors to register app information

## 21.1.0 - 2017-04-12
* Allow coupon to be specified when creating orders
* No longer require that items have descriptions when creating orders

## 21.0.0 - 2017-04-07
* Balances are now retrieved by payout instead of by transfer

## 20.0.0 - 2017-04-06
* Bump API version to 2017-04-06: https://stripe.com/docs/upgrades#2017-04-06
* Add support for payouts and recipient transfers
* Change the transfer resource to support its new format
* Deprecate recipient creation
* Disputes under charges are now expandable and collapsed by default
* Rules under charge outcomes are now expandable and collapsed by default

## 19.17.0 - 2017-04-06
* Please see 20.0.0 (bad release)	

## 19.16.0 - 2017-03-23
* Allow the ID of an identity document to be passed into an account owner update

## 19.15.0 - 2017-03-22
* Add `ShippingCarrier` to dispute evidence

## 19.14.0 - 2017-03-20
* Add `Period`, `Plan`, and `Quantity` to `InvoiceItem`

## 19.13.0 - 2017-03-20
* Add `AdditionalOwnersEmpty` to allow additional owners to be unset

## 19.12.0 - 2017-03-17
* Add new form of file upload using `io.FileReader` and filename

## 19.11.0 - 2017-03-13
* Add `Token` to `SourceObjectParams`

## 19.10.0 - 2017-03-13
* Add `CouponEmpty` (allowing a coupon to be cleared) to customer parameters
* Add `CouponEmpty` (allowing a coupon to be cleared) to subscription parameters

## 19.9.0 - 2017-03-08
* Add missing value "all" to subscription statuses

## 19.8.0 - 2017-03-02
* Add subscription items client to main `client.API` struct

## 19.7.0 - 2017-03-01
* Add `Statement` (statement descriptor) to `CaptureParams`

## 19.6.0 - 2017-02-22
* Add new parameters for invoices and subscriptions

## 19.5.0 - 2017-02-13
* Add new rich `Destination` type to `ChargeParams`

## 19.4.0 - 2017-02-03
* Support Connect account as payment source

## 19.3.0 - 2017-02-02
* Add transfer group to charges and transfers

## 19.2.0 - 2017-01-23
* Add `Rule` to `ChargeOutcome`

## 19.1.0 - 2017-01-18
* Add support for updating sources

## 19.0.2 - 2017-01-04
* Fix subscription `trial_period_days` to be populated by the right value

## 19.0.1 - 2016-12-08
* Include verification document details when persisting `LegalEntity`

## 19.0.0 - 2016-12-07
* Remote `SubProrationDateNow` field from `InvoiceParams`

## 18.14.1 - 2016-12-05
* Truncate `tax_percent` at four decimals (e.g. 3.9750%) instead of two

## 18.14.0 - 2016-11-23
* Add retrieve method for 3-D Secure resources

## 18.13.0 - 2016-11-15
* Add `PaymentSource` to `API`

## 18.12.0 - 2016-11-14
* Allow bank accounts to be created as a customer source

## 18.11.0 - 2016-11-14
* Add `TrialPeriodEnd` to `SubParams`

## 18.10.0 - 2016-11-09
* Add `StatusTransitions` to `Order`

## 18.9.0 - 2016-11-04
* Add `Application` to `Charge`

## 18.8.0 - 2016-10-24
* Add `Review` to `Charge` for the charge reviews

## 18.7.0 - 2016-10-18
* Add `RiskLevel` to `ChargeOutcome`

## 18.6.0 - 2016-10-18
* Support for 403 status codes (permission denied)

## 18.5.0 - 2016-10-18
* Add `Status` to `SubListParams` to allow filtering subscriptions by status

## 18.4.0 - 2016-10-14
* Add `HasEvidence` and `PastDue` to `EvidenceDetails`

## 18.3.0 - 2016-10-10
* Add `NoDiscountable` to `InvoiceItemParams`

## 18.2.0 - 2016-10-10
* Add `BusinessLogo` to `Account`
* Add `ReceiptNumber` to `Charge`
* Add `DestPayment` to `Transfer`

## 18.1.0 - 2016-10-04
* Support for Apple Pay domains

## 18.0.0 - 2016-10-03
* Support for subscription items
* Correct `SourceTx` on `Transfer` to be a `SourceTransaction`
* Change `Charge` on `Resource` to be expandable (now a struct instead of string)

## 17.5.0 - 2016-09-22
* Support customer-related operations for bank accounts

## 17.4.2 - 2016-09-19
* Fix but where some parameters were not being included on order update

## 17.4.1 - 2016-09-15
* Fix bug that required a date of birth to be included on account update

## 17.4.0 - 2016-09-13
* Add missing Kana and Kanji address and name fields to account's legal entity
* Add `ReceiptNumber` and `Status` to `Refund`

## 17.3.0 - 2016-09-07
* Add support for sources endpoint

## 17.2.0 - 2016-08-29
* Add order returns to `API`

## 17.1.0 - 2016-08-22
* Add `DeactiveOn` to `Product`

## 17.0.0 - 2016-08-18
* Allow expansion of destination on transfers
* Allow expansion of sources on balance transactions

## 16.8.0 - 2016-08-17
* Add `OriginatingTransaction` to `Fee`

## 16.7.1 - 2016-08-17
* Allow params to be nil when retrieving a refund

## 16.7.0 - 2016-08-11
* Add support for 3-D Secure

## 16.6.0 - 2016-08-09
* Add `ReceiptNumber` to `Invoice`

## 16.5.0 - 2016-08-08
* Add `Meta` to `Account`

## 16.4.0 - 2016-08-05
* Allow the migration of recipients to accounts
* Add `MigratedTo` to `Recipient`

## 16.3.1 - 2016-07-25
* URL-escape the IDs of coupons and plans when making API requests

## 16.3.0 - 2016-07-19
* Add `NoClosed` to `InvoiceParams` to allow an invoice to be reopened

## 16.2.1 - 2016-07-11
* Consider `SubParams.QuantityZero` when updating a subscription

## 16.2.0 - 2016-07-07
* Upgrade API version to 2016-07-06

## 16.1.0 - 2016-07-07
* Add `Returns` field to `Order`

## 16.0.0 - 2016-06-30
* Remove `Name` field on `SKU`; it's not actually supported
* Support updating `Product` on `SKU`

## 15.6.0 - 2016-06-24
* Allow product and SKU attributes to be updated

## 15.5.0 - 2016-06-24
* Add `TaxPercent` and `TaxPercentZero` to `CustomerParams`

## 15.4.0 - 2016-06-20
* Add `TokenizationMethod` to `Card` struct

## 15.3.0 - 2016-06-15
* Add `BalanceZero` to `CustomerParams` so that balance can be zeroed out

## 15.2.0 - 2016-06-03
* Add `ToValues` to `RequestValues` struct

## 15.1.0 - 2016-05-26
* Add `BusinessVatID` to customer creation parameters

## 15.0.0 - 2016-05-24
* Fix handling of nested objects in arrays in request parameters

## 14.4.0 - 2016-05-24
* Add granular error types in new `Err` field on `stripe.Error`

## 14.3.0 - 2016-05-20
* Allow Relay orders to be returned and add associated types

## 14.2.3 - 2016-05-20
* When creating a bank account token, only send routing number if it's been set

## 14.2.2 - 2016-05-17
* When creating a bank account, only send routing number if it's been set

## 14.2.1 - 2016-05-17
* Add missing SKU clinet to client API type

## 14.2.0 - 2016-05-11
* Add `Reversed` and `AmountReversed` fields to `Transfer`

## 14.1.0 - 2016-05-05
* Allow `default_for_currency` to be set when creating a card

## 14.0.0 - 2016-05-04
* Change the signature for `sub.Delete`. The customer ID is no longer required.

## 13.12.0 - 2016-04-28
* Add `Currency` to `Card`

## 13.11.1 - 2016-04-22
* Fix bug where new external accounts could not be marked default from token

## 13.11.0 - 2016-04-21
* Expose a number of list types that were previously internal (full list below)
* Expose `stripe.AccountList`
* Expose `stripe.TransactionList`
* Expose `stripe.BitcoinReceiverList`
* Expose `stripe.ChargeList`
* Expose `stripe.CountrySpecList`
* Expose `stripe.CouponList`
* Expose `stripe.CustomerList`
* Expose `stripe.DisputeList`
* Expose `stripe.EventList`
* Expose `stripe.FeeList`
* Expose `stripe.FileUploadList`
* Expose `stripe.InvoiceList`
* Expose `stripe.OrderList`
* Expose `stripe.ProductList`
* Expose `stripe.RecipientList`
* Expose `stripe.TransferList`
* Switch to use of `stripe.BitcoinTransactionList`
* Switch to use of `stripe.SKUList`

## 13.10.1 - 2016-04-20
* Add support for `TaxPercentZero` to invoice and subscription updates

## 13.10.0 - 2016-04-19
* Expose `stripe.PlanList` (previously an internal type)

## 13.9.0 - 2016-04-18
* Add `TaxPercentZero` struct to `InvoiceParams`
* Add `TaxPercentZero` to `SubParams`

## 13.8.0 - 2016-04-12
* Add `Outcome` struct to `Charge`

## 13.7.0 - 2016-04-06
* Add `Description`, `IIN`, and `Issuer` to `Card`

## 13.6.0 - 2016-04-05
* Add `SourceType` (and associated constants) to `Transfer`

## 13.5.0 - 2016-03-29
* Add `Meta` (metadata) to `BankAccount`

## 13.4.0 - 2016-03-29
* Add `Meta` (metadata) to `Card`

## 13.3.0 - 2016-03-29
* Add `DefaultCurrency` to `CountrySpec`

## 13.2.0 - 2016-03-18
* Add `SourceTransfer` to `Charge`
* Add `SourceTx` to `Transfer`

## 13.1.0 - 2016-03-15
* Add `Reject` on `Account` to support the new API feature

## 13.0.0 - 2016-03-15
* Upgrade API version to 2016-03-07
* Remove `Account.BankAccounts` in favor of `ExternalAccounts`
* Remove `Account.Currencies` in favor of `CountrySpec`

## 12.1.0 - 2016-02-04
* Add `ListParams.StripeAccount` for making list calls on behalf of connected accounts
* Add `Params.StripeAccount` for symmetry with `ListParams.StripeAccount`
* Deprecate `Params.Account` in favor of `Params.StripeAccount`

## 12.0.0 - 2016-02-02
* Add support for fetching events for managed accounts (`event.Get` now takes `Params`)

## 11.5.0 - 2016-02-26
* Allow a `PII.PersonalIDNumber` number to be used to create a token

## 11.4.0 - 2016-02-24
* Add missing subscription fields to `InvoiceParams` for use with `invoice.GetNext`

## 11.3.0 - 2016-02-19
* Add `AccountHolderName` and `AccountHolderType` to bank accounts

## 11.2.0 - 2016-02-11
* Add support for `CountrySpec`
* Add `SSNProvided`, `PersonalIDProvided` and `BusinessTaxIDProvided` to `LegalEntity`

## 11.1.2 - 2016-02-02
* Fix card update method to correctly take expiration date

## 11.1.1 - 2016-02-01
* Fix recipient update so that it can take a bank token (like create)

## 11.0.1 - 2016-01-11
* Add missing field `country` to shipping details of `Charge` and `Customer`

## 11.0.0 - 2016-01-07
* Add missing field `Default` to `BankAccount`
* Add `OrderParams` parameter to `Order` retrieval
* Fix parameter bug when creating a new `Order`
* Support special value of 'now' for trial end when updating subscriptions

## 10.3.0 - 2015-12-10
* Allow an account to be referenced when creating a card

## 10.2.0 - 2015-12-04
* Add `Update` function on `Coupon` client so that metadata can be set

## 10.1.0 - 2015-12-01
* Add a verification routine for external accounts

## 10.0.0 - 2015-11-30
* Return models along with `error` when deleting resources with `Del`
* Fix bug where country parameter wasn't included for some account creation

## 9.0.0 - 2015-11-13
* Return model (`Sub`) when cancelling a subscription (`sub.Cancel`)

## 8.0.0 - 2015-08-17
* Add ability to list and retrieve refunds without a Charge

## 7.0.0 - 2015-08-03
* Add ability to list and retrieve disputes

## 6.8.0 - 2015-07-29
* Add ability to delete an account

## 6.7.1 - 2015-07-17
* Bug fixes

## 6.7.0 - 2015-07-16
* Expand logging object
* Move proration date to subscription update
* Send country when creating/updating account

## 6.6.0 - 2015-07-06
* Add request ID to errors

## 6.5.0 - 2015-07-06
* Update bank account creation API
* Add destination, application fee, transfer to Charge struct
* Add missing fields to invoice line item
* Rename deprecated customer param value

## 6.4.2 - 2015-06-23
* Add BusinessUrl, BusinessUrl, BusinessPrimaryColor, SupportEmail, and
* SupportUrl to Account.

## 6.4.1 - 2015-06-16
* Change card.dynamic_last_four to card.dynamic_last4

## 6.4.0 - 2015-05-28
* Rename customer.default_card -> default_source

## 6.3.0 - 2015-05-19
* Add shipping address to charges
* Expose card.dynamic_last_four
* Expose account.tos_acceptance
* Bug fixes
* Bump API version to most recent one

## 6.2.0 - 2015-04-09
* Bug fixes
* Add Extra to parameters

## 6.1.0 - 2015-03-17
* Add TaxPercent for subscriptions
* Event bug fixes

## 6.0.0 - 2015-03-15
* Add more operations for /accounts endpoint
* Add /transfers/reversals endpoint
* Add /accounts/bank_accounts endpoint
* Add support for Stripe-Account header

## 5.1.0 - 2015-02-25
* Add new dispute status `warning_closed`
* Add SubParams.TrialEndNow to support `trial_end = "now"`

## 5.0.1 - 2015-02-25
* Fix URL for upcoming invoices

## 5.0.0 - 2015-02-19
* Bump to API version 2014-02-18
* Change Card, DefaultCard, Cards to Source, DefaultSource, Sources in Stripe response objects
* Add paymentsource package for manipulating Customer's sources
* Support Update action for Bitcoin Receivers

## 4.4.3 - 2015-02-08
* Modify NewIdempotencyKey() algorithm to increase likelihood of randomness

## 4.4.2 - 2015-01-24
* Add BankAccountParams.Token
* Add Token.ClientIP
* Add LogLevel

## 4.4.0 - 2015-01-20
* Add Bitcoin support

## 4.3.0 - 2015-01-13
* Added support for listing FileUploads
* Mime parameter on FileUpload has been changed to Type

## 4.2.1 - 2014-12-28
* Handle charges with customer card tokens

## 4.2.0 - 2014-12-18
* Add idempotency support

## 4.1.0 - 2014-12-17
* Bump to API version 2014-12-17.

## 4.0.0 - 2014-12-16
* Add FileUpload resource. This brings in a new endpoint (uploads.stripe.com) and thus makes changes to some of the existing interfaces.
* This also adds support for multipart content.

## 3.1.0 - 2014-12-16
* Add Charge.FraudDetails

## 3.0.1 - 2014-12-15
* Add timeout value to HTTP requests

## 3.0.0 - 2014-12-05
* Add Dispute.EvidenceDetails
* Remove Dispute.DueDate
* Change Dispute.Evidence from string to struct

## 2.0.0 - 2014-11-26
* Change List interface to .Next() and .Resource()
* Better error messages for Get() methods
* EventData.Raw contains the raw event message
* SubParams.QuantityZero can be used for free subscriptions

## 1.0.3 - 2014-10-22
* Add AddMeta method

## 1.0.2 - 2014-09-23
* Minor fixes

## 1.0.1 - 2014-09-23
* Linter-based updates

## 1.0.0 - 2014-09-22
* Initial version
