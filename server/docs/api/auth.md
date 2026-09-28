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
  * **`201 Created`**: Registration successful. Returns the generated user ID.
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
  * **`201 Created`**: Registration successful[cite: 12].
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
  * **`200 OK`**: Account verification successful.
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
  * **`200 OK`**: Login successful (returns access and refresh tokens).
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
  * **`200 OK`**: Password reset code sent to the email.
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
  * **`200 OK`**: Password reset successful[cite: 9].
  * **`400 Bad Request`**: Invalid or expired OTP / validation failure.