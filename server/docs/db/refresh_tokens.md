# Database Schema for Refresh Token
This document specifies the technical specification for tables and indexes related to the Refresh Token schema for **Tickr**.

---

## 1. `refresh_tokens` Table

| Attribute | Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `id` | UUID | PRIMARY KEY, NOT NULL, DEFAULT `gen_random_uuid()` | Token ID |
| `user_id` | UUID | NOT NULL, FK (`users`) | User of the session |
| `hashed_token` | Text | NOT NULL | Token Hash |
| `device_info` | Text | | Detail of the client device |
| `expires_at` | TIMESTAMPZ | NOT NULL | Expiration timestamp |
| `created_at` | TIMESTAMPZ | DEFAULT `NOW()` | Creation time |

---

## 2. `idx_refresh_tokens_user_id` Index 

| Index Name | Target Table | Indexed Columns |
| :--- | :--- | :--- |
| `idx_refresh_tokens_user_id` | `refresh_tokens` | `user_id` |

---

## 2. `idx_refresh_tokens_hashed_token` Index

| Index Name | Target Table | Indexed Columns |
| :--- | :--- | :--- |
| `idx_refresh_tokens_hashed_token` | `refresh_tokens` | `hashed_token` |

---