# Business Rules

The `risk` package takes care of all business validation rules.

For the Authorization process we have the following rules.

* `insecureAuthorization` Validates that if an authorization is not 3DS not exemption and not recurring should be rejected.
* `validateInitialRecurring` Validates that an initial authorization recurring must be merchant or full authenticated.
* `lowValueExemption` Validates that when the low value exemption is used the amount cannot be higher than EUR30 (if higher reject).