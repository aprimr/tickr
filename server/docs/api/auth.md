# Auth API Reference

Base URL: `http://localhost:8000/v1`

---

### Register Regular User
Registers a new user account.

* **URL:** `/users`
* **Method:** `POST`
* **Rate Limit:** 10 requests per 5 minutes
* **Body:**
  ```json
  {
    "fullname": "Full name",
    "email": "example@gmail.com",
    "phone_number": "9811100000",
    "password": "Password@123"
  }
  ```

* **Responses:**
  * **`201 Created`**:
  ```json
  {
    "success": true,
    "status": 201,
    "message": "registration successful",
    "data": {
      "user_id": "16e63693-46e..",
    }
  }
  ```
  * **`400 Bad Request`**: Validation failed or invalid request body.
  * **`409 Conflict`**: An account with this email already exists.

---

### Register Venue Admin
Registers a new venue administrator account.

* **URL:** `/venues`
* **Method:** `POST`
* **Rate Limit:** 10 requests per 5 minutes
* **Body:**
  ```json
  {
    "email": "admin@gmail.com",
    "phone_number": "9911100099",
    "password": "Admin@123",
    "venue_name": "VFX CINEMAS",
    "address": "Ghorahi-13,Dang",
    "city": "ghorahi",
    "total_screens": 3
  }
  ```

* **Responses:**
  * **`201 Created`**: 
  ```json
  {
    "success": true,
    "status": 201,
    "message": "registration successful",
    "data": {
      "user_id": "16e63693-46e..",
    }
  }
  ```
  * **`400 Bad Request`**: Validation failed or invalid request body.
  * **`409 Conflict`**: Email already registered.

---

### Verify Account
Verifies a user's email account using the OTP code sent upon registration.

* **URL:** `/auth/verify-email`
* **Method:** `POST`
* **Rate Limit:** 10 requests per 5 minutes
* **Body:**
  ```json
  {
    "user_id": "6925aba9-2c...d8dcf5",
    "otp": "388236"
  }
  ```
* **Responses:**
  * **`200 OK`**: 
  ```json
  {
    "success": true,
    "status": 201,
    "message": "verification successful"
  }
  ```
  * **`400 Bad Request`**: Invalid or expired OTP / validation failure.
  * **`404 Not Found`**: User not found.

---

### User Login
Authenticates users and returns session tokens.

>Unverified users cannot login until email verification is completed.

>Inactive users cannot login until they are activated.

* **URL:** `/auth/login`
* **Method:** `POST`
* **Rate Limit:** 10 requests per 5 minutes
* **Body:**
  ```json
  {
    "email": "example@gmail.com",
    "password": "Password@123"
  }
  ```
* **Responses:**
  * **`200 OK`**: 
  ```json
  {
    "success": true,
    "status": 200,
    "message": "token rotation successful",
    "data": {
      "access_token": "eyJhbGciOiJIUzI1NiIsI...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsI..."
    }
  }
  ```
  * **`401 Unauthorized`**: Invalid credentials.
  * **`403 Forbidden`**: Email not verified (`"email is not verified, please verify it first"`) or account deactivated.

---

### Forgot Password
Triggers a password reset request, generating and emailing a recovery OTP to the user.

* **URL:** `/auth/forgot-password`
* **Method:** `POST`
* **Rate Limit:** 5 requests per 15 minutes
* **Body:**
  ```json
  {
    "email": "example@gmail.com"
  }
  ```
* **Responses:**
  * **`200 OK`**:
  ```json
  {
    "success": true,
    "status": 200,
    "message": "password reset code sent to the email"
  }
  ```
  * **`400 Bad Request`**: Validation error.
  * **`429 Too Many Requests`**: An active unexpired OTP already exists for this user.

---

### Reset Password
Verifies the forgot-password OTP, updates the user's password, marks the email verified, and invalidates the token.

* **URL:** `/auth/reset-password`
* **Method:** `POST`
* **Rate Limit:** 5 requests per 15 minutes
* **Body:**
  ```json
  {
    "email": "example@gmail.com",
    "otp": "499631",
    "new_password": "NewPassword@123"
  }
  ```
* **Responses:**
  * **`200 OK`**: 
  ```json
  {
    "success": true,
    "status": 200,
    "message": "password reset successful"
  }
  ```
  * **`400 Bad Request`**: Invalid or expired OTP / validation failure.

### Rotate Token
Rotate the access and refresh token of a user.

> Expired or invalid refresh tokens will be rejected and the user must re-authenticate.

* **URL:** `/auth/refresh-token`
* **Method:** `POST`
* **Rate Limit:** `20 requests per 5 minutes`
* **Body:**
  ```json
  {
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```
* **Responses**:
  * **`200 OK`:** Token rotation successful (returns new access and refresh tokens).
  ```json
  {
    "success": true,
    "status": 200,
    "message": "token rotation successful",
    "data": {
      "access_token": "eyJhbGciOiJIUzI1NiIsI...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsI..."
    }
  }
  ```
  * **`400 Bad Request`:** invalid request body
  * **`401 Unauthorized`:** unauthorized access
  * **`500 Internal Server Error`:** Something went wrong

### Logout
Logouts current users session

* **URL:** `/auth/logout`
* **Method**: `POST`
* **Rate Limit**: `10 requests per 5 minutes`
* **Headers:**
  * `Authorization`: Bearer `<access_token>`
* **Body**:
  ```json
  {
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
* **Responses**:
  * **`200 OK`:**
  ```json
  {
    "success": true,
    "status": 200,
    "message": "logout successful",
  }
  ```
  * **`400 Bad Request`:** Invalid request body or validation failure.
  * **`500 Internal Server Error`:** Unexpected server failure during logout.

### Logout All Devices
Revokes all active sessions for the current user.

* **URL:** `/auth/logout-all`
* **Method:** `POST`
* **Rate Limit:** 10 requests per 1 minute
* **Headers:**
  * `Authorization`: Bearer `<access_token>`
* **Responses:**
  * **`200 OK`**: 
  ```json
  {
    "success": true,
    "status": 200,
    "message": "logout successful"
  }
  ```
  * **`401 Unauthorized`**: Unauthorized access.
  * **`500 Internal Server Error`**: Something went wrong.

### Get Active Sessions
Returns all active sessions for the currently logged-in user.

> Requires the `x-refresh-token` header to validate the request context.

* **URL:** `/auth/sessions`
* **Method:** `GET`
* **Rate Limit:** 3 requests per 1 minute
* **Headers:**
  * `Authorization`: Bearer `<access_token>`
  * `x-refresh-token`: `<refresh_token>`
* **Responses:**
  * **`200 OK`**: 
  ```json
  {
    "success": true,
    "status": 200,
    "message": "sessions fetched successfully",
    "data": [
      {
      "session_id": "b7e28743-f5..",
      "device_info": "Chorme on Unknown OS",
      "is_current": true,
      "created_at": "2026-10-03T21:55:54.545251+05:45"
    },
    ]
  }
  ```
  * **`401 Unauthorized`**: Unauthorized access or missing refresh token.
  * **`500 Internal Server Error`**: Something went wrong.

### Delete Session
Deletes a specific session by its unique session ID.

* **URL:** `/auth/sessions/{sessionID}`
* **Method:** `DELETE`
* **Rate Limit:** 10 requests per 1 minute
* **URL Parameters:**
  * `sessionID`: UUID of the target session
* **Headers:**
  * `Authorization`: Bearer `<access_token>`
* **Responses:**
  * **`200 OK`**: 
  ```json
  {
    "success": true,
    "status": 200,
    "message": "session delete successful"
  }
  ```
  * **`400 Bad Request`**: Session ID is missing.
  * **`401 Unauthorized`**: Unauthorized access.
  * **`500 Internal Server Error`**: Something went wrong.