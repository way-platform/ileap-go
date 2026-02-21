---
name: clerk
description: Guide for authenticating a user via the Clerk Frontend API (FAPI) using username/password. Use when implementing the iLEAP Authentication Server Adapter or similar authentication backends that require programmatic username/password login via Clerk.
---

# Clerk Authentication Skill

This skill provides guidance on implementing an authentication backend using Clerk. Since iLEAP requires programmatic username/password authentication and doesn't support 3-legged OAuth flows, we must use the Clerk Frontend API (FAPI) to authenticate users directly.

## Using the Clerk Frontend API (FAPI) for Login

Clerk's design separates backend and frontend APIs. The Backend API (BAPI) is used for administrative tasks, while the Frontend API (FAPI) is used for authentication flows (logging in, signing up, etc.).

To programmatically authenticate a user with a username (or email) and password, you must interact with the FAPI `sign_ins` endpoint.

### 1. Initiate the Sign-In Flow

Send a `POST` request to the `/v1/client/sign_ins` endpoint on your Clerk Frontend API URL (e.g., `https://clerk.<your-domain>/v1/client/sign_ins`).

**Request Body (`application/x-www-form-urlencoded` or JSON):**
- `strategy`: Must be `"password"`
- `identifier`: The user's email address or username
- `password`: The user's password

**Example Request:**
```http
POST /v1/client/sign_ins
Content-Type: application/x-www-form-urlencoded

strategy=password&identifier=ileap-demo@way.cloud&password=HelloPrimaryData
```

### 2. Handle the Response

The API will return a `SignIn` object. You must check the `status` field to determine if the authentication was successful.

- If `status === "complete"`: The username and password are correct, and no further verification (like MFA) is required. The response will include a `created_session_id`.
- If `status === "needs_second_factor"`: The user has Multi-Factor Authentication enabled. (Note: For the iLEAP demo backend, MFA should ideally be disabled for the service account).
- If the credentials are invalid, the API will return a `4xx` error.

**Example Success Response:**
```json
{
  "status": "complete",
  "created_session_id": "sess_1234567890",
  "identifier": "ileap-demo@way.cloud",
  // ... other fields
}
```

### 3. Integrating with iLEAP

When adapting this for the iLEAP Authentication Server Adapter (e.g., the `POST /auth/token` route):
1. Extract the `username` and `password` from the incoming HTTP Basic Auth request.
2. Proxy these credentials to the Clerk FAPI `/v1/client/sign_ins` endpoint.
3. If Clerk returns `status: "complete"`, consider the user authenticated.
4. Generate the required iLEAP token (as seen in `demo/server.go`) and return it to the client.

## Reference Material

For further details on Clerk's APIs and custom flows, refer to the bundled reference documentation:

- **OpenAPI FAPI Spec:** `.agents/skills/clerk/references/openapi/frontend-api-2025-11-10.yml` (Look for `/v1/client/sign_ins`)
- **Custom Email/Password Flow Docs:** `.agents/skills/clerk/references/docs/guides/development/custom-flows/authentication/email-password.mdx`
- **Backend API Spec (for administrative tasks):** `.agents/skills/clerk/references/openapi/backend-api-2025-11-10.yml`

*Note: The Clerk Backend API also provides a `/users/{user_id}/verify_password` endpoint, but it requires knowing the `user_id` beforehand. The FAPI `/v1/client/sign_ins` endpoint allows direct authentication with just the identifier and password.*