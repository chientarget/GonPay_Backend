
---

# GonPay API Documentation

## Mục lục
1. [Giới thiệu](#1-giới-thiệu)
2. [Thông tin chung](#2-thông-tin-chung)
3. [Xác thực (Authentication)](#3-xác-thực-authentication)
4. [Quản lý người dùng](#4-quản-lý-người-dùng)
5. [Quản lý ví điện tử](#5-quản-lý-ví-điện-tử)
6. [Quản lý giao dịch](#6-quản-lý-giao-dịch)
7. [Phương thức thanh toán](#7-phương-thức-thanh-toán)
8. [Người thụ hưởng](#8-người-thụ-hưởng)
9. [Hạn mức giao dịch](#9-hạn-mức-giao-dịch)
10. [Thông báo](#10-thông-báo)
11. [Nhật ký hệ thống](#11-nhật-ký-hệ-thống)
12. [Mã lỗi và xử lý](#12-mã-lỗi-và-xử-lý)


## 1. Giới thiệu

GonPay là hệ thống ví điện tử cung cấp các dịch vụ tài chính như chuyển tiền, nạp tiền, và quản lý tài khoản. API được thiết kế theo tiêu chuẩn RESTful.

## 2. Thông tin chung

### Base URL
```
https://gonpay-backend.onrender.com/v1
```

### Headers yêu cầu
```
Content-Type: application/json
Authorization: Bearer <jwt_token>
```

### Pagination Format
```json
{
    "data": [],
    "metadata": {
        "page": 1,
        "limit": 10,
        "total": 100,
        "total_pages": 10
    }
}
```

## 3. Xác thực (Authentication)

### 3.1. Đăng ký tài khoản
```http
POST /api/register
```

**Request Body:**
```json
{
    "username": "nguyenvana",
    "email": "nguyenvana@gmail.com",
    "phone_number": "+84912345678",
    "password": "Password123"
}
```

**Response (201 Created):**
```json
{
    "user": {
        "id": 1,
        "username": "nguyenvana",
        "email": "nguyenvana@gmail.com",
        "phone_number": "+84912345678",
        "status": "ACTIVE",
        "created_at": "2024-11-18T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### 3.2. Đăng nhập
```http
POST /api/login
```

**Request Body:**
```json
{
    "email": "nguyenvana@gmail.com",
    "password": "Password123"
}
```

**Response (200 OK):**
```json
{
    "user": {
        "id": 1,
        "username": "nguyenvana",
        "email": "nguyenvana@gmail.com",
        "role": "USER",
        "status": "ACTIVE"
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

## 4. Quản lý người dùng

### 4.1. Xem thông tin cá nhân
```http
GET /api/users/profile
```

**Response (200 OK):**
```json
{
    "id": 1,
    "username": "nguyenvana",
    "email": "nguyenvana@gmail.com",
    "phone_number": "+84912345678",
    "status": "ACTIVE",
    "preferences": {
        "language": "vi",
        "notification_enabled": true
    },
    "created_at": "2024-11-18T10:00:00Z",
    "updated_at": "2024-11-18T10:00:00Z"
}
```

### 4.2. Cập nhật thông tin cá nhân
```http
PUT /api/users/profile
```

**Request Body:**
```json
{
    "username": "nguyenvana_new",
    "email": "nguyenvana_new@gmail.com",
    "phone_number": "+84912345679"
}
```

**Response (200 OK):**
```json
{
    "id": 1,
    "username": "nguyenvana_new",
    "email": "nguyenvana_new@gmail.com",
    "phone_number": "+84912345679",
    "status": "ACTIVE",
    "updated_at": "2024-11-18T11:00:00Z"
}
```

### 4.3. Đổi mật khẩu
```http
PUT /api/users/password
```

**Request Body:**
```json
{
    "old_password": "Password123",
    "new_password": "NewPassword123"
}
```

**Response (200 OK):**
```json
{
    "message": "Mật khẩu đã được thay đổi thành công"
}
```

## 5. Quản lý ví điện tử

### 5.1. Tạo ví mới
```http
POST /api/wallets
```

**Response (201 Created):**
```json
{
    "id": 1,
    "user_id": 1,
    "wallet_number": "W123456789",
    "balance": 0,
    "status": "ACTIVE",
    "created_at": "2024-11-18T10:00:00Z"
}
```

### 5.2. Danh sách ví
```http
GET /api/wallets
```

**Response (200 OK):**
```json
{
    "data": [
        {
            "id": 1,
            "wallet_number": "W123456789",
            "balance": 1000000,
            "status": "ACTIVE",
            "created_at": "2024-11-18T10:00:00Z"
        }
    ]
}
```

### 5.3. Chi tiết ví
```http
GET /api/wallets/{id}
```

**Response (200 OK):**
```json
{
    "id": 1,
    "wallet_number": "W123456789",
    "balance": 1000000,
    "status": "ACTIVE",
    "created_at": "2024-11-18T10:00:00Z",
    "transactions_count": 10,
    "last_transaction_at": "2024-11-18T15:00:00Z"
}
```

### 5.4. Hủy kích hoạt ví
```http
POST /api/wallets/{id}/deactivate
```

**Response (200 OK):**
```json
{
    "message": "Ví đã được hủy kích hoạt thành công"
}
```

## 6. Quản lý giao dịch

### 6.1. Chuyển tiền
```http
POST /api/wallets/transfer
```

**Request Body:**
```json
{
    "source_wallet_id": 1,
    "destination_wallet_id": 2,
    "amount": 1000000,
    "description": "Chuyển tiền cho bạn"
}
```

**Response (200 OK):**
```json
{
    "transaction_id": 1,
    "reference_id": "TRX123456789",
    "source_wallet_id": 1,
    "destination_wallet_id": 2,
    "type": "TRANSFER",
    "amount": 1000000,
    "status": "COMPLETED",
    "description": "Chuyển tiền cho bạn",
    "created_at": "2024-11-18T10:00:00Z"
}
```

### 6.2. Nạp tiền
```http
POST /api/wallets/{id}/deposit
```

**Request Body:**
```json
{
    "amount": 1000000,
    "description": "Nạp tiền vào ví"
}
```

**Response (200 OK):**
```json
{
    "transaction_id": 2,
    "reference_id": "DEP123456789",
    "type": "DEPOSIT",
    "amount": 1000000,
    "status": "COMPLETED",
    "description": "Nạp tiền vào ví",
    "created_at": "2024-11-18T10:00:00Z"
}
```

### 6.3. Rút tiền
```http
POST /api/wallets/{id}/withdraw
```

**Request Body:**
```json
{
    "amount": 1000000,
    "description": "Rút tiền về tài khoản ngân hàng"
}
```

**Response (200 OK):**
```json
{
    "transaction_id": 3,
    "reference_id": "WIT123456789",
    "type": "WITHDRAW",
    "amount": 1000000,
    "status": "COM

PLETED",
    "description": "Rút tiền về tài khoản ngân hàng",
    "created_at": "2024-11-18T10:00:00Z"
}
```

## 7. Phương thức thanh toán

### 7.1. Danh sách phương thức thanh toán
```http
GET /api/payment-methods
```

**Response (200 OK):**
```json
{
    "data": [
        {
            "id": 1,
            "method": "Bank Transfer",
            "status": "ACTIVE",
            "created_at": "2024-11-18T10:00:00Z"
        },
        {
            "id": 2,
            "method": "Credit Card",
            "status": "ACTIVE",
            "created_at": "2024-11-18T10:00:00Z"
        }
    ]
}
```

## 8. Người thụ hưởng

### 8.1. Thêm người thụ hưởng
```http
POST /api/beneficiaries
```

**Request Body:**
```json
{
    "name": "Nguyễn Văn A",
    "account_number": "123456789",
    "bank": "Ngân hàng Vietcombank"
}
```

**Response (201 Created):**
```json
{
    "id": 1,
    "name": "Nguyễn Văn A",
    "account_number": "123456789",
    "bank": "Ngân hàng Vietcombank",
    "status": "ACTIVE",
    "created_at": "2024-11-18T10:00:00Z"
}
```

## 9. Hạn mức giao dịch

### 9.1. Xem hạn mức giao dịch
```http
GET /api/transaction-limits
```

**Response (200 OK):**
```json
{
    "daily_limit": 5000000,
    "monthly_limit": 10000000
}
```

## 10. Thông báo

### 10.1. Xem thông báo
```http
GET /api/notifications
```

**Response (200 OK):**
```json
{
    "data": [
        {
            "id": 1,
            "title": "Thông báo hệ thống",
            "message": "Hệ thống bảo trì từ 10:00 đến 12:00 ngày mai.",
            "created_at": "2024-11-18T10:00:00Z"
        }
    ]
}
```

## 11. Nhật ký hệ thống

### 11.1. Xem nhật ký
```http
GET /api/system-logs
```

**Response (200 OK):**
```json
{
    "data": [
        {
            "id": 1,
            "message": "Đăng nhập thành công",
            "timestamp": "2024-11-18T10:00:00Z"
        }
    ]
}
```

## 12. Mã lỗi và xử lý

| Mã lỗi | Thông báo                  | Miêu tả                       |
|--------|----------------------------|-------------------------------|
| 400    | Bad Request                | Dữ liệu không hợp lệ          |
| 401    | Unauthorized               | Chưa xác thực                 |
| 404    | Not Found                  | Tài nguyên không tìm thấy     |
| 500    | Internal Server Error      | Lỗi máy chủ                   |

Cám ơn !