# Database Schema for OTP
This document specifies the technical specification for tables and enums related to the OTP schema for **Tickr**.

---

## 1. `otps` Table

| Attribute | Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `id` | UUID | PRIMARY KEY, NOT NULL, DEFAULT `gen_random_uuid()` | OTP ID |
| `user_id` | UUID | NOT NULL, FK (`users`) | User ID |
| `hashed_otp` | Text | NOT NULL | OTP HASH |
| `type` | otp_type | NOT NULL | OTP type |
| `is_used` | Boolean | DEFAULT `FALSE` | Flacg for tracking used otps |
| `expires_at` | TIMESTAMPZ | NOT NULL | Expiration timestamp |
| `created_at` | TIMESTAMPZ | DEFAULT `NOW()` | Creation time |

---

## 2. `otp_type` Enum

| Enum Value | Description |
| :--- | :--- |
| `account_verification` | Verification code sent to the user when user sign-up to verify their email address. |
| `forgot_password` | Verification code sent when a user requests to reset their password. |

---

## 3. `idx_otps_user_type` Index

| Index Name | Target Table | Indexed Columns |
| :--- | :--- | :--- |
| `idx_otps_user_type` | `otps` | `user_id` `type` |

---