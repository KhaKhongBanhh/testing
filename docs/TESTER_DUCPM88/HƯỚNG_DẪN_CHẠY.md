🌐 Trạng thái hệ thống đang hoạt động
gõ 2 lệnh sau trong terminal:
cd FE
npm start
Frontend (Giao diện React): Đang chạy tại http://localhost:3000
Reverse Proxy (Nginx): Đang chạy tại http://localhost
API Gateway: Đang chạy tại http://localhost:8000
Các dịch vụ Microservices (Docker):
postgres (Port 5432 - Database cinema_app với pgvector)
redis (Port 6379)
api-gateway (Port 8000)
auth-service (Port 3001, gRPC 50052)
user-service (Port 8005, gRPC 50051)
movie-service (Port 8083, gRPC 50053)
booking-service (Port 8082, gRPC 50082)
payment-service (Port 8086)
notification-service (Port 8080)
analytics-service (Port 8087)
worker-service (Port 8088, gRPC 50083)
chatbot (Port 8089)
nginx (Port 80)
🔑 Tài khoản mẫu đã nạp sẵn vào Database (Demo Accounts)
Vai trò	Email đăng nhập	Mật khẩu	Ghi chú
Khách hàng (Customer)	alice@email.com	password	Đăng nhập trên giao diện đặt vé xem phim
Quản trị viên (Admin)	admin@cinema.com	password	Toàn quyền quản trị hệ thống
Quản lý rạp (Manager)	manager@cinema.com	password	Quản lý phim, phòng chiếu, lịch chiếu
Nhân viên vé (Ticket Staff)	ticket@cinema.com	password	Kiểm tra và xuất vé tại quầy
📝 Tóm tắt các thao tác đã hoàn thành
Tạo cấu hình .env: Thiết lập đồng bộ các biến môi trường cho 10 microservices, frontend và module migrate.
Tích hợp module Migrate vào Docker: Thêm Dockerfile và profile tools vào docker-compose.yml giúp chạy migration tự động qua Docker mà máy host không cần cài Go.
Migration & Seeding: Tạo toàn bộ schema và nạp 7 phim, 5 phòng chiếu, 555 ghế ngồi, 30 suất chiếu và các tài khoản mẫu.
Cập nhật tài liệu: Đã viết lại mục 🚀 Hướng Dẫn Cài Đặt và Khởi Chạy trong README.md đầy đủ 8 bước từ cài đặt, chạy backend, chạy migrate, chạy frontend đến bảng port và tài khoản demo.