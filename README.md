# Demtech SES API Mock

## Idea of the API Mocking

Instead of using the AWS SES API to send actual emails, we wanted to mock it to avoid billing costs and to make testing our dependent services easier and faster. The mocking can be implemented in several ways:

1. **Client-Controlled Success/Failure Ratio** – Allows the client to specify a success/failure ratio, introducing random failures within a defined percentage.
2. **Deterministic Error Distribution** – Predefines the percentage of each error, though this approach lacks realism.
3. **Validation-Based Failures (Chosen Approach)** – Implements various validations and randomly fails for most SES errors based on a configured percentage while mocking email sending (no emails are actually sent via SMTP).

## Special Rules

1. **Quota Validator** – The API enforces a quota limit. If a client tries to send more than X (configurable via environment variables) messages within the last N hours (configurable via environment variables), the API returns a `LimitExceededException` error until older messages move out of the time window. *(See: `QuotaValidator.go`)*
2. **SES Warming-Up Mechanism** – Amazon SES enforces gradual sending limits for new accounts to prevent spam and protect sender reputation:
   - **Initial Limits** – New accounts start with a low daily limit (e.g., 200 emails/day).
   - **Automatic Increase** – SES increases the limit based on good deliverability, low bounce, and complaint rates.
   - **Per-Second Throttling** – SES restricts the number of emails sent per second to prevent sudden spikes.
   
   *Note: This feature is not yet implemented due to time constraints.*

## Errors and Their Meanings

### 1. Request Errors (4xx)
Errors related to request issues, such as missing parameters or exceeding limits.

#### Authentication & Authorization Errors
- `AccessDeniedException` – Insufficient IAM permissions.
- `InvalidClientTokenId` – Invalid security token.
- `SignatureDoesNotMatch` – Incorrect AWS request signature.
- `RequestExpired` – Request timestamp is too old.

#### Validation Errors
- `ValidationError` – Request parameters failed validation.
- `InvalidParameterValue` – Invalid request parameters.
- `MissingParameter` – A required parameter is missing.
- `InvalidAction` – Requested action is not valid.

#### Quota & Limit Exceeded Errors
- `ThrottlingException` – Request throttled due to SES limits.
- `SendingQuotaExceededException` – Daily sending limit exceeded.
- `TooManyRequestsException` – Too many requests in a short time.
- `MessageRejected` – Message rejected due to spam filtering or size limits.

### 2. Email Sending Errors
Errors occurring while attempting to send an email.

#### Bounced Email Issues
- `MessageRejected` – High bounce rate or policy violation.
- `MailFromDomainNotVerifiedException` – MAIL FROM address not verified.
- `ConfigurationSetDoesNotExistException` – Configuration set does not exist.

#### Recipient Issues
- `RecipientBlacklisted` – Recipient is on AWS SES suppression list.
- `EmailAddressBlacklisted` – Email address is blacklisted.
- `TemplateDoesNotExistException` – Email template does not exist.

### 3. Internal & Service Errors (5xx)
Errors caused by AWS SES that require retrying.
- `InternalFailure` – Internal AWS SES error.
- `ServiceUnavailable` – AWS SES temporarily unavailable.
- `EndpointConnectionError` – AWS SES cannot connect to the endpoint.

## API Endpoints

### 1. Sending Email (Mock API Behavior)
The API validates requests similarly to AWS SES but lacks full error message fidelity due to time constraints. However, the design is **highly extensible**, allowing additional validations while maintaining SOLID principles.

#### Example Request
```sh
curl -X POST "http://localhost:8080/api/v1/send-email" \
  -H "Content-Type: application/json" \
  -d '{
    "Source": "sender@example.com",
    "Destination": {
      "ToAddresses": ["recipient@example.com"]
    },
    "Message": {
      "Subject": {
        "Data": "Test Email Subject"
      },
      "Body": {
        "Text": {
          "Data": "This is the email body."
        }
      }
    },
    "ReturnPath": "bounce@example.com"
  }'
```

#### Example Responses
- **Success**
  ```json
  {"message":"Email sent successfully","messageId":"a9cc1cc1-eb65-48ef-b092-ae0621c44498"}
  ```
- **Failure**
  ```json
  {"error":"LimitExceededException","message":"Sending quota exceeded"}
  ```

### 2. Reading Statistics
Retrieves current email sending statistics, including success and failure counts.

#### Example Request
```sh
curl localhost:8080/api/v1/email-stats
```

#### Example Response
```json
{
  "totalEmailsSent": 12,
  "successCount": 10,
  "totalErrCount": 2,
  "errors": {
    "LimitExceededException": 2
  }
}
```

## Running Tests

- **Unit Tests** *(Faster Execution)*
  ```sh
  make unit-test
  ```
- **All Tests (Unit + Integration)**
  ```sh
  make test
  ```
- **End-to-End Tests Only**
  ```sh
  make e2e
  ```

## Running the API Server in Docker
Runs the HTTP server and Redis storage in Docker.

```sh
make up
```

To stop and remove the Docker containers:
```sh
make down
```

*Note: Check `Makefile` for additional recipes and details.*

## Improvements

- More test cases can be added to cover further edge cases esp. for E2E tests.
- The AWS SES API is really extensive and the documentation is quite spread and quite time consuming. Remaining behaviour that are missing, I wish I could implement them.
- Instrumentation using OpenTelemetry.
- More statistic can be added as per demand.

