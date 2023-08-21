# Open Issues after _Handler Structs Refactor_

## VISA Specific
* We don't send in 3DS.
* Not sure but I think the fields F15, F49 and F126.9 are mandatory.
* Not sure which field to use for ecommerce indicator for Visa. Probably 60.8.
* We need to test refunds and we need to build and test reversals (not implemented yet).
* We need to implement MOTO.
* Response codes are copied over from MC but in the docs they are different.
* I don't trust the implementation of posEntryMode (and the rest) because I got an error on the length of this field.
* The card acceptor ID needs to be padded with spaces. THis is now done with fmt.Sprintf but should be done by the Visa lib.

## Generic missing things
* We need to reimplement cross field validations.
* We need to setup Insomnia in a proper way.
* The communication with the EAS and MIP should be moved to adapters (just like pubsub).
* Remove reference from `mastercard_authorizations` and move `authorization_id_response` to `authorizations` table. Both schemes use the latter field.
* `internal/authorization/ports/http.go:167` is resolved lean and mean. We need this as either an ecommerce indicator or a SLI with which we can use the string method.
* `pkg/visa/base1/f001.go:22` -> Should be part of the iso
* Low value <= €30,- must be implemented.
* TTC is set to "P" hardcoded in reversals but should be T. Better to fetch this from the original auth I guess? `internal/processing/scheme/mastercard/reversal.go:128`. This value should be stored in the DB and fetched upon reversal. MC will set this for us, we need to store it and send it in for reversal.
* AllowProductionCardNumbers has to be reimplemented. Maybe we can tackle this by limitted cardranges? Do these ranges overlap?
* Add authorizationType to refunds input. Should always be final for now.
* `Card.Info.IssuerCountryCode` should be of type `countryCode`.
* Countrycode on entity.PointOfServiceData should be of type countrycode Not string.
* Change amount and currency into money.Money.
* Function `internal/infrastructure/spanner/authorizationrepository.go:965` needs to be refactored. Maybe set the default value on the http handler default to time.Now?
* Circumvent timing attack on captures and reversals -> Right now I can send in an auth, partially capture it and reverse it. I don't know how Mastercard processes this, but I can imagine a decline on the presentment (because it is already reversed). We can circumvent this by implementing partial reversals.
* `LocalTransactionDateTime` needs to be changed. Should be moved to entities?
* Remove field `Entity.Recurring.Subsequent`. Figure out what to do with all cases where it's used.
* Ecommerce indicator implementation is a bit shitty because MC and Visa handle this differently.
* TraceID is with capital I and D. We should align this with the rest which is Id.
* createdAt is in the list of sortable columns but not in the result set of get authorizations. At least in the docs. Should be added.
* Move app `internal/infrastructure/spanner/**` code that is _vertical specific_ into the `adapters` package for that vertical.