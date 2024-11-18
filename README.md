# GonPay - Hệ Thống API Ví Điện Tử

<div align="center">
    <img src="/api/placeholder/200/200" alt="GonPay Logo" width="200">
    <p><em>Giải pháp thanh toán điện tử toàn diện</em></p>
</div>

## 📑 Mục lục
- [Tổng quan](#-tổng-quan)
- [Kiến trúc hệ thống](#-kiến-trúc-hệ-thống)
- [Yêu cầu hệ thống](#-yêu-cầu-hệ-thống)
- [Cài đặt và Phát triển](#-cài-đặt-và-phát-triển)
- [Chi tiết API](#-chi-tiết-api)
- [Bảo mật](#-bảo-mật)
- [Đóng góp](#-đóng-góp)
- [Hỗ trợ](#-hỗ-trợ)

## 🌟 Tổng quan

GonPay là hệ thống API ví điện tử được xây dựng với Go, cung cấp các dịch vụ:
- Quản lý tài khoản và xác thực
- Quản lý ví điện tử
- Chuyển tiền và thanh toán
- Quản lý người thụ hưởng
- Theo dõi giao dịch
- Thông báo và báo cáo

### Công nghệ sử dụng
- Go 1.21+
- PostgreSQL 15+
- JWT Authentication
- Clean Architecture
- RESTful API

## 📚 Chi tiết API

### Thông tin chung

**Base URL:**
```
https://api.gonpay.com/v1
```

**Headers chung:**
```http
Content-Type: application/json
Accept: application/json
```

**Headers cho API cần xác thực:**
```http
Authorization: Bearer <jwt_token>
```

### 1. Xác thực (Authentication)

#### 1.1. Đăng ký tài khoản [`POST /api/register`]

**Request Body:**
```json
{
  "username": "nguyenvana",
  "email": "nguyenvana@gmail.com",
  "phone_number": "+84912345678",
  "password": "Password123"
}
```

**Validation:**
- `username`: 3-50 ký tự, chỉ chữ và số
- `email`: định dạng email hợp lệ
- `phone_number`: định dạng E.164
- `password`: ít nhất 6 ký tự, chứa chữ hoa, thường và số

**Success Response (201 Created):**
```json
{
  "data": {
    "user": {
      "id": 1,
      "username": "nguyenvana",
      "email": "nguyenvana@gmail.com",
      "phone_number": "+84912345678",
      "status": "ACTIVE",
      "created_at": "2024-11-18T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
  },
  "message": "Đăng ký tài khoản thành công"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Dữ liệu không hợp lệ",
    "details": {
      "email": "Email không đúng định dạng",
      "password": "Mật khẩu phải có ít nhất 6 ký tự"
    }
  }
}
```

#### 1.2. Đăng nhập [`POST /api/login`]

**Request Body:**
```json
{
  "email": "nguyenvana@gmail.com",
  "password": "Password123"
}
```

**Success Response (200 OK):**
```json
{
  "data": {
    "user": {
      "id": 1,
      "username": "nguyenvana",
      "email": "nguyenvana@gmail.com",
      "role": "USER",
      "last_login": "2024-11-18T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
  },
  "message": "Đăng nhập thành công"
}
```

### 2. Quản lý người dùng

#### 2.1. Xem thông tin cá nhân [`GET /api/users/profile`]

**Headers required:**
```http
Authorization: Bearer <jwt_token>
```

**Success Response (200 OK):**
```json
{
  "data": {
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
}
```

#### 2.2. Cập nhật thông tin cá nhân [`PUT /api/users/profile`]

**Headers required:**
```http
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "username": "nguyenvana_new",
  "email": "new_email@gmail.com",
  "phone_number": "+84987654321",
  "preferences": {
    "language": "en",
    "notification_enabled": false
  }
}
```

### 3. Quản lý ví điện tử

#### 3.1. Tạo ví mới [`POST /api/wallets`]

**Headers required:**
```http
Authorization: Bearer <jwt_token>
```

**Success Response (201 Created):**
```json
{
  "data": {
    "id": 1,
    "wallet_number": "W123456789",
    "balance": 0,
    "status": "ACTIVE",
    "created_at": "2024-11-18T10:00:00Z"
  },
  "message": "Tạo ví thành công"
}
```

#### 3.2. Danh sách ví [`GET /api/wallets`]

**Headers required:**
```http
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `status` (optional): ACTIVE, INACTIVE
- `sort` (optional): created_at, balance
- `order` (optional): asc, desc
- `page` (optional): Default 1
- `limit` (optional): Default 10

**Success Response (200 OK):**
```json
{
  "data": {
    "wallets": [
      {
        "id": 1,
        "wallet_number": "W123456789",
        "balance": 1000000,
        "status": "ACTIVE",
        "created_at": "2024-11-18T10:00:00Z",
        "transaction_count": 15,
        "last_transaction": "2024-11-18T15:00:00Z"
      }
    ],
    "metadata": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

### 4. Giao dịch

#### 4.1. Chuyển tiền [`POST /api/wallets/transfer`]

**Headers required:**
```http
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "source_wallet_id": 1,
  "destination_wallet_id": 2,
  "amount": 1000000,
  "description": "Chuyển tiền cho bạn",
  "reference_number": "TRX123456"
}
```

**Success Response (200 OK):**
```json
{
  "data": {
    "transaction_id": 1,
    "reference_id": "TRX123456",
    "source_wallet": {
      "id": 1,
      "number": "W123456789",
      "balance_after": 9000000
    },
    "destination_wallet": {
      "id": 2,
      "number": "W987654321",
      "balance_after": 11000000
    },
    "amount": 1000000,
    "type": "TRANSFER",
    "status": "COMPLETED",
    "description": "Chuyển tiền cho bạn",
    "created_at": "2024-11-18T10:00:00Z"
  },
  "message": "Giao dịch thành công"
}
```

#### 4.2. Nạp tiền [`POST /api/wallets/{id}/deposit`]

**Headers required:**
```http
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "amount": 1000000,
  "payment_method_id": 1,
  "description": "Nạp tiền vào ví"
}
```

**Success Response (200 OK):**
```json
{
  "data": {
    "transaction_id": 2,
    "reference_id": "DEP123456",
    "wallet": {
      "id": 1,
      "balance_after": 2000000
    },
    "amount": 1000000,
    "type": "DEPOSIT",
    "status": "COMPLETED",
    "created_at": "2024-11-18T10:00:00Z"
  }
}
```

### 5. Phương thức thanh toán

#### 5.1. Thêm phương thức thanh toán [`POST /api/payment-methods`]

**Headers required:**
```http
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "method_type": "BANK_ACCOUNT",
  "account_number": "123456789",
  "bank_name": "Vietcombank",
  "branch": "Hà Nội",
  "account_holder": "NGUYEN VAN A",
  "is_default": true
}
```

### 6. Người thụ hưởng

#### 6.1. Thêm người thụ hưởng [`POST /api/beneficiaries`]

**Headers required:**
```http
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "name": "Nguyễn Văn B",
  "account_identifier": "987654321",
  "account_type": "BANK_ACCOUNT",
  "bank_name": "Vietcombank",
  "bank_branch": "Ho Chi Minh",
  "relationship": "FRIEND"
}
```

### 7. Hạn mức và bảo mật

#### 7.1. Thiết lập hạn mức [`POST /api/limits`]

**Headers required:**
```http
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "transaction_type": "TRANSFER",
  "daily_limit": 50000000,
  "monthly_limit": 1000000000,
  "enabled": true
}
```

### 8. Thông báo

#### 8.1. Danh sách thông báo [`GET /api/notifications`]

**Headers required:**
```http
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `type`: TRANSACTION, SECURITY, SYSTEM
- `is_read`: true, false
- `page`: Default 1
- `limit`: Default 20

**Success Response (200 OK):**
```json
{
  "data": {
    "notifications": [
      {
        "id": 1,
        "type": "TRANSACTION",
        "title": "Giao dịch thành công",
        "content": "Bạn đã chuyển 1.000.000đ cho số tài khoản 9876543210",
        "is_read": false,
        "created_at": "2024-11-18T10:00:00Z"
      }
    ],
    "metadata": {
      "page": 1,
      "limit": 20,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

## 🔒 Bảo mật

### Xác thực và Phân quyền
- Sử dụng JWT (JSON Web Token)
- Token hết hạn sau 24 giờ
- Refresh token có thời hạn 30 ngày
- Role-based access control (RBAC):
    - USER: Người dùng thông thường
    - ADMIN: Quản trị viên
    - SYSTEM: Hệ thống

### Rate Limiting
```
PUBLIC APIs:
- 60 requests/phút/IP
- 1000 requests/ngày/IP

AUTHENTICATED APIs:
- 1000 requests/phút/user
- 10000 requests/ngày/user

ADMIN APIs:
- 2000 requests/phút/admin
- 50000 requests/ngày/admin
```

### Mã hóa và Bảo mật
- Tất cả kết nối phải sử dụng HTTPS
- Mật khẩu được hash bằng bcrypt
- Sensitive data được mã hóa trong database
- Access logs được lưu trữ 90 ngày

## 📊 Giám sát và Logging

### Monitoring
- Uptime monitoring
- Performance metrics
- Error tracking
- Resource usage

### Logging
```json
{
    "timestamp": "2024-11-18T10:00:00Z",
    "level": "INFO",
    "method": "POST",
    "path": "/api/wallets/transfer",
    "user_id": 1,
    "ip": "127.0.0.1",
    "duration": 235,
    "status": 200
}
```

## 🤝 Đóng góp và Phát triển

### Quy trình phát triển
1. Fork repository
2. Tạo feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

### Coding Standards
- Tuân thủ Go standards
- Documentation bắt buộc
- Unit tests coverage > 80%
- Integration tests cho APIs

## 📞 Hỗ trợ

### Kênh hỗ trợ
- Email: chientarget@gmail.com 

<div align="center">
    <p>Copyright © 2024 GonPay. All rights reserved.</p>
</div>