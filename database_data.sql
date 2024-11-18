-- Xóa dữ liệu cũ (nếu cần)
TRUNCATE TABLE audit_logs, notifications, transaction_limits, beneficiaries,
    payment_methods, transactions, wallets, users CASCADE;

-- Reset sequences
ALTER SEQUENCE users_user_id_seq RESTART WITH 1;
ALTER SEQUENCE wallets_wallet_id_seq RESTART WITH 1;
ALTER SEQUENCE transactions_transaction_id_seq RESTART WITH 1;
ALTER SEQUENCE payment_methods_payment_method_id_seq RESTART WITH 1;
ALTER SEQUENCE beneficiaries_beneficiary_id_seq RESTART WITH 1;
ALTER SEQUENCE transaction_limits_limit_id_seq RESTART WITH 1;
ALTER SEQUENCE notifications_notification_id_seq RESTART WITH 1;
ALTER SEQUENCE audit_logs_log_id_seq RESTART WITH 1;

-- Insert users (Mật khẩu mã hóa của "Password123")
INSERT INTO users (username, email, phone_number, password_hash, status, preferences)
VALUES ('nguyenvana', 'nguyenvana@gmail.com', '+84912345678', '$2a$10$xJ9Y8mZ5q5d5vP5jvz0K6OX5YZ5Z5YZ5Z5YZ5Z5YZ5Z5YZ5Z5Y', 'ACTIVE', '{"language": "vi", "notification_enabled": true}'),
       ('tranthib', 'tranthib@gmail.com', '+84923456789', '$2a$10$xJ9Y8mZ5q5d5vP5jvz0K6OX5YZ5Z5YZ5Z5YZ5Z5YZ5Z5YZ5Z5Y', 'ACTIVE', '{"language": "vi", "notification_enabled": true}'),
       ('levanc', 'levanc@gmail.com', '+84934567890', '$2a$10$xJ9Y8mZ5q5d5vP5jvz0K6OX5YZ5Z5YZ5Z5YZ5Z5YZ5Z5YZ5Z5Y', 'ACTIVE', '{"language": "vi", "notification_enabled": false}'),
       ('phamthid', 'phamthid@gmail.com', '+84945678901', '$2a$10$xJ9Y8mZ5q5d5vP5jvz0K6OX5YZ5Z5YZ5Z5YZ5Z5YZ5Z5YZ5Z5Y', 'ACTIVE', '{"language": "vi", "notification_enabled": true}'),
       ('hoangvane', 'hoangvane@gmail.com', '+84956789012', '$2a$10$xJ9Y8mZ5q5d5vP5jvz0K6OX5YZ5Z5YZ5Z5YZ5Z5YZ5Z5YZ5Z5Y', 'ACTIVE', '{"language": "vi", "notification_enabled": true}'),
       ('dangvanf', 'dangvanf@gmail.com', '+84967890123', '$2a$10$xJ9Y8mZ5q5d5vP5jvz0K6OX5YZ5Z5YZ5Z5YZ5Z5YZ5Z5YZ5Z5Y', 'ACTIVE', '{"language": "vi", "notification_enabled": true}'),
       ('vuthig', 'vuthig@gmail.com', '+84978901234', '$2a$10$xJ9Y8mZ5q5d5vP5jvz0K6OX5YZ5Z5YZ5Z5YZ5Z5YZ5Z5YZ5Z5Y', 'ACTIVE', '{"language": "vi", "notification_enabled": false}'),
       ('buivanh', 'buivanh@gmail.com', '+84989012345', '$2a$10$xJ9Y8mZ5q5d5vP5jvz0K6OX5YZ5Z5YZ5Z5YZ5Z5YZ5Z5YZ5Z5Y', 'ACTIVE', '{"language": "vi", "notification_enabled": true}'),
       ('ngothii', 'ngothii@gmail.com', '+84990123456', '$2a$10$xJ9Y8mZ5q5d5vP5jvz0K6OX5YZ5Z5YZ5Z5YZ5Z5YZ5Z5YZ5Z5Y', 'ACTIVE', '{"language": "vi", "notification_enabled": true}'),
       ('dothij', 'dothij@gmail.com', '+84901234567', '$2a$10$xJ9Y8mZ5q5d5vP5jvz0K6OX5YZ5Z5YZ5Z5YZ5Z5YZ5Z5YZ5Z5Y', 'INACTIVE', '{"language": "vi", "notification_enabled": false}');

-- Insert wallets
INSERT INTO wallets (user_id, balance, status)
VALUES (1, 5000000, 'ACTIVE'),
       (1, 2000000, 'ACTIVE'),
       (2, 10000000, 'ACTIVE'),
       (3, 15000000, 'ACTIVE'),
       (4, 8000000, 'ACTIVE'),
       (5, 12000000, 'ACTIVE'),
       (6, 6000000, 'ACTIVE'),
       (7, 9000000, 'ACTIVE'),
       (8, 7500000, 'ACTIVE'),
       (9, 3000000, 'ACTIVE');

-- Insert transactions
INSERT INTO transactions (source_wallet_id, destination_wallet_id, transaction_type, amount, status, description)
VALUES (1, 3, 'TRANSFER', 1000000, 'COMPLETED', 'Chuyển tiền mua hàng'),
       (3, NULL, 'WITHDRAW', 500000, 'COMPLETED', 'Rút tiền ATM'),
       (2, NULL, 'WITHDRAW', 1000000, 'COMPLETED', 'Rút tiền mặt'),
       (4, 5, 'TRANSFER', 2000000, 'COMPLETED', 'Thanh toán hóa đơn'),
       (5, NULL, 'DEPOSIT', 5000000, 'COMPLETED', 'Nạp tiền từ ngân hàng'),
       (6, 7, 'TRANSFER', 1500000, 'COMPLETED', 'Chuyển tiền cho bạn'),
       (8, NULL, 'DEPOSIT', 3000000, 'COMPLETED', 'Nạp tiền từ thẻ tín dụng'),
       (9, 1, 'TRANSFER', 500000, 'COMPLETED', 'Hoàn tiền'),
       (1, NULL, 'WITHDRAW', 200000, 'FAILED', 'Rút tiền thất bại'),
       (2, 4, 'TRANSFER', 1000000, 'COMPLETED', 'Thanh toán dịch vụ');

-- Insert payment_methods
INSERT INTO payment_methods (user_id, method_type, account_number, bank_name, is_default)
VALUES (1, 'BANK_ACCOUNT', '1234567890', 'Vietcombank', true),
       (2, 'BANK_ACCOUNT', '2345678901', 'Techcombank', true),
       (3, 'CREDIT_CARD', '4532456789012345', 'VPBank', false),
       (4, 'E_WALLET', '3456789012', 'MoMo', true),
       (5, 'DEBIT_CARD', '5432109876543210', 'BIDV', true),
       (6, 'BANK_ACCOUNT', '6789012345', 'ACB', true),
       (7, 'CREDIT_CARD', '4532456789012346', 'Sacombank', false),
       (8, 'E_WALLET', '7890123456', 'ZaloPay', true),
       (9, 'BANK_ACCOUNT', '8901234567', 'MB Bank', true),
       (10, 'DEBIT_CARD', '5432109876543211', 'Agribank', true);

-- Insert beneficiaries
INSERT INTO beneficiaries (user_id, beneficiary_name, account_identifier, account_type, bank_name)
VALUES (1, 'Nguyễn Văn Bình', '9876543210', 'BANK_ACCOUNT', 'Vietcombank'),
       (2, 'Trần Thị Cúc', '8765432109', 'BANK_ACCOUNT', 'Techcombank'),
       (3, 'Lê Văn Dũng', '7654321098', 'WALLET', NULL),
       (4, 'Phạm Thị Em', '6543210987', 'BANK_ACCOUNT', 'BIDV'),
       (5, 'Hoàng Văn Phú', '5432109876', 'WALLET', NULL),
       (6, 'Đặng Thị Giang', '4321098765', 'BANK_ACCOUNT', 'MB Bank'),
       (7, 'Vũ Văn Hùng', '3210987654', 'BANK_ACCOUNT', 'VPBank'),
       (8, 'Bùi Thị Lan', '2109876543', 'WALLET', NULL),
       (9, 'Ngô Văn Minh', '1098765432', 'BANK_ACCOUNT', 'Agribank'),
       (10, 'Đỗ Thị Nam', '0987654321', 'BANK_ACCOUNT', 'Sacombank');

-- Insert transaction_limits
INSERT INTO transaction_limits (user_id, transaction_type, daily_limit, monthly_limit)
VALUES (1, 'TRANSFER', 50000000, 1000000000),
       (1, 'WITHDRAW', 20000000, 500000000),
       (2, 'TRANSFER', 100000000, 2000000000),
       (3, 'WITHDRAW', 30000000, 600000000),
       (4, 'TRANSFER', 70000000, 1500000000),
       (5, 'WITHDRAW', 25000000, 550000000),
       (6, 'TRANSFER', 80000000, 1800000000),
       (7, 'WITHDRAW', 35000000, 700000000),
       (8, 'TRANSFER', 60000000, 1200000000),
       (9, 'WITHDRAW', 40000000, 800000000);

-- Insert notifications
INSERT INTO notifications (user_id, title, content, notification_type, is_read)
VALUES (1, 'Giao dịch thành công', 'Bạn đã chuyển 1.000.000đ cho số tài khoản 9876543210', 'TRANSACTION', false),
       (2, 'Nhận tiền thành công', 'Bạn đã nhận 2.000.000đ từ số tài khoản 1234567890', 'TRANSACTION', true),
       (3, 'Rút tiền thành công', 'Giao dịch rút 500.000đ đã hoàn tất', 'TRANSACTION', false),
       (4, 'Cảnh báo bảo mật', 'Phát hiện đăng nhập từ thiết bị mới', 'SECURITY', true),
       (5, 'Nạp tiền thành công', 'Tài khoản được nạp 5.000.000đ', 'TRANSACTION', false),
       (6, 'Thông báo giới hạn', 'Gần đạt giới hạn giao dịch ngày', 'LIMIT', true),
       (7, 'Khuyến mãi', 'Ưu đãi giảm phí chuyển khoản trong tháng', 'PROMOTION', false),
       (8, 'Giao dịch thất bại', 'Giao dịch chuyển tiền không thành công', 'TRANSACTION', true),
       (9, 'Xác thực tài khoản', 'Vui lòng xác thực để nâng cấp tài khoản', 'ACCOUNT', false),
       (10, 'Cập nhật hệ thống', 'Hệ thống sẽ bảo trì lúc 23:00 tối nay', 'SYSTEM', true);


-- Insert audit_logs
INSERT INTO audit_logs (user_id, action, entity_type, entity_id, old_value, new_value, ip_address, user_agent)
VALUES (1, 'LOGIN', 'USER', 1, NULL, NULL, '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/96.0.4664.110'),
       (2, 'TRANSFER', 'TRANSACTION', 1,
        '{"balance": 7000000}',
        '{"balance": 6000000}',
        '192.168.1.1',
        'Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15'),
       (3, 'UPDATE_PROFILE', 'USER', 3,
        '{"phone": "+84934567890"}',
        '{"phone": "+84934567891"}',
        '10.0.0.1',
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'),
       (4, 'ADD_BENEFICIARY', 'BENEFICIARY', 1,
        NULL,
        '{"account": "9876543210", "bank": "Vietcombank"}',
        '172.16.0.1',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Firefox/95.0'),
       (5, 'CHANGE_PASSWORD', 'USER', 5,
        NULL,
        '{"password_changed": true, "timestamp": "2024-11-17T10:30:00Z"}',
        '192.168.0.1',
        'Mozilla/5.0 (Linux; Android 11; SM-G991B)'),
       (6, 'ADD_PAYMENT_METHOD', 'PAYMENT_METHOD', 3,
        NULL,
        '{"type": "CREDIT_CARD", "bank": "VPBank"}',
        '10.10.0.1',
        'Mozilla/5.0 (iPad; CPU OS 15_1 like Mac OS X)'),
       (7, 'WITHDRAW', 'TRANSACTION', 3,
        '{"balance": 9000000}',
        '{"balance": 8000000, "fee": 11000}',
        '192.168.1.100',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edge/96.0.1054.43'),
       (8, 'UPDATE_LIMITS', 'TRANSACTION_LIMIT', 1,
        '{"daily": 40000000, "monthly": 800000000}',
        '{"daily": 50000000, "monthly": 1000000000}',
        '172.16.1.1',
        'Mozilla/5.0 (iPhone; CPU iPhone OS 15_1 like Mac OS X)'),
       (9, 'DEPOSIT', 'TRANSACTION', 7,
        '{"balance": 3000000}',
        '{"balance": 6000000, "method": "Bank Transfer"}',
        '10.0.0.100',
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Safari/605.1.15'),
       (10, 'LOGOUT', 'USER', 1,
        NULL,
        '{"session_duration": "01:30:15", "logout_type": "manual"}',
        '127.0.0.1',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/96.0.4664.110');

-- Thêm một số bản ghi audit cho các hoạt động bảo mật
INSERT INTO audit_logs (user_id, action, entity_type, entity_id, old_value, new_value, ip_address, user_agent)
VALUES (1, 'ENABLE_2FA', 'SECURITY', 1,
        '{"2fa_enabled": false}',
        '{"2fa_enabled": true, "method": "SMS"}',
        '127.0.0.1',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/96.0.4664.110'),
       (2, 'UPDATE_SECURITY', 'USER', 2,
        '{"login_notification": false}',
        '{"login_notification": true, "email_alert": true}',
        '192.168.1.1',
        'Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X)'),
       (3, 'FAILED_LOGIN', 'SECURITY', 3,
        NULL,
        '{"attempt": 1, "reason": "Invalid password"}',
        '10.0.0.1',
        'Mozilla/5.0 (Linux; Android 11; SM-G991B)'),
       (4, 'BLOCK_ACCOUNT', 'SECURITY', 4,
        '{"status": "ACTIVE"}',
        '{"status": "BLOCKED", "reason": "Multiple failed login attempts"}',
        '172.16.0.1',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64)'),
       (5, 'PASSWORD_RESET', 'USER', 5,
        NULL,
        '{"reset_method": "Email", "timestamp": "2024-11-17T15:45:00Z"}',
        '192.168.0.1',
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)');

-- Commit transaction nếu tất cả thành công
COMMIT;

-- Verify data
SELECT COUNT(*)
FROM users;
SELECT COUNT(*)
FROM wallets;
SELECT COUNT(*)
FROM transactions;
SELECT COUNT(*)
FROM payment_methods;
SELECT COUNT(*)
FROM beneficiaries;
SELECT COUNT(*)
FROM transaction_limits;
SELECT COUNT(*)
FROM notifications;
SELECT COUNT(*)
FROM audit_logs;