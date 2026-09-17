# Database Schema for Tickr Users
This document specifies the technical specification for tables and enums related to users of the multi-tenant authentication and Role-Based Access Control (RBAC) schema for **Tickr**.


### 1. `users` Table
Base authentication details for all users.


| Attribute | Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `id` | UUID | PRIMARY KEY, NOT NULL | User ID |
| `email` | Text | UNIQUE, NOT NULL | Login email |
| `phone_UUID` | Text | UNIQUE | Contact UUID |
| `password_hash` | Text | NOT NULL | Hashed password |
| `role` | Text | NOT NULL, DEFAULT `'user'` | Access role ("user" ,"venue_admin", "super_admin") |
| `is_active` | Boolean | DEFAULT `TRUE` | Active status |
| `is_email_verified` | Boolean | DEFAULT `FALSE` | Email verification status |
| `is_phone_verified` | Boolean | DEFAULT `FALSE` | Phone verification status |
| `last_login_at` | Time | NULL | Last login time |
| `created_at` | Time | DEFAULT `NOW()` | Creation time |
| `updated_at` | Time | DEFAULT `NOW()` | Update time |

---

### 2. `user_details` Table
Profiel details for end users with role `user`

| Attribute | Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `id` | UUID | PRIMARY KEY, NOT NULL | Profile ID |
| `user_id` | UUID | FK (`users`), UNIQUE, NOT NULL | Linked user |
| `full_name` | Text | NOT NULL | Customer name |
| `updated_at` | Time | DEFAULT `NOW()` | Update time |

---

### 3. `venue_details` Table
Profile details for venue or theaters with role `venue_admin`


| Attribute | Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `id` | UUID | PRIMARY KEY, NOT NULL | Venue ID |
| `user_id` | UUID | FK (`users`), UNIQUE, NOT NULL | Linked owner |
| `venue_name` | Text | NOT NULL | Venue name |
| `address` | Text | NOT NULL | Street address |
| `city` | Text | NOT NULL | City location |
| `total_screens` | NUMBER | NOT NULL, DEFAULT `1` | Screen count |
| `venue_status` | Text | NOT NULL, DEFAULT `'pending'` | Approval state ("pending", "approved", "rejected") |
| `approved_by` | UUID | FK (`users`), NULL | Approver admin |
| `approved_at` | Time | NULL | Approval time |
| `rejected_by` | UUID | FK (`users`), NULL | Rejecter admin |
| `rejected_at` | Time | NULL | Rejection time |
| `venue_created_at` | Time | DEFAULT `NOW()` | Creation time |
| `venue_updated_at` | Time | DEFAULT `NOW()` | Update time |

---

### 4. `admin_details` Table
Profile details for platform admin with role `super_admin`

| Attribute | Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `id` | UUID | PRIMARY KEY, NOT NULL | Admin ID |
| `user_id` | UUID | FK (`users`), UNIQUE, NOT NULL | Linked admin |
| `full_name` | Text | NOT NULL | Admin name |
| `is_root_user` | Boolean | DEFAULT `FALSE` | Root user flag |
| `admin_created_at` | Time | DEFAULT `NOW()` | Creation time |
| `admin_updated_at` | Time | DEFAULT `NOW()` | Update time |

---

### 5. `user_role` Enum

| Role Value | Description |
| :--- | :--- |
| `user` | Standard user account. Can browse movies, book tickets, and manage personal bookings. |
| `venue_admin` | Theater/venue owner account. Can manage theater screens, movie showtimes, and cinema details. |
| `super_admin` | Platform admin account. Has full administrative access, manages venue approvals/rejections. |

---

### 6. `venue_status` Enum

| Status Value | Description |
| :--- | :--- |
| `pending` | Awaiting review by a `super_admin`. |
| `approved` | Venue has passed platform compliance and can list movies and sell tickets. |
| `rejected` | Venue registration was denied by a `super_admin`. |

---