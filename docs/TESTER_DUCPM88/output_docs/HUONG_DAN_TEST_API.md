Tôi đã tạo một file Markdown hướng dẫn chi tiết toàn bộ cách test API cho dự án tại: 📂 HUONG_DAN_TEST_API.md

Dưới đây là tóm tắt nhanh những điểm quan trọng nhất giúp bạn thực hiện test ngay trên Postman / Thunder Client / cURL:

1. Quy ước chung
Base URL qua API Gateway: http://localhost:8000/api/v1
Header chung khi gửi body: Content-Type: application/json
Header xác thực khi đăng nhập: Authorization: Bearer <chuỗi_token>
2. Bước 1: Đăng nhập để lấy Token (Bắt buộc)

Hệ thống phân quyền (RBAC) chia làm 2 cổng đăng nhập:

A. Đăng nhập với quyền Khách hàng (Customer)
Method: POST
URL: http://localhost:8000/api/v1/auth/login
Body (JSON):
json
{
  "email": "alice@email.com",
  "password": "password"
}
Kết quả: Nhận chuỗi token trong trường "token".
B. Đăng nhập với quyền Quản lý / Admin (Manager / Admin)
Method: POST
URL: http://localhost:8000/api/v1/auth/admin/login (Lưu ý: Bắt buộc dùng URL /admin/login)
Body (JSON):
json
{
  "email": "manager@cinema.com",
  "password": "password"
}
Kết quả: Nhận chuỗi token quản lý để thao tác thêm/sửa/xóa phim, suất chiếu.

💡 Cách dùng Token ở các request tiếp theo:
Vào tab Headers của Postman/Thunder Client, thêm header:
Authorization: Bearer <dán_chuỗi_token_vào_đây>

3. Bước 2: Test Xem Danh sách Phim & Lịch chiếu (Trọng tâm)
🎬 A. Xem danh sách phim
Method: GET
URL: http://localhost:8000/api/v1/movies
Header: Authorization: Bearer <token> (Tùy chọn: có hay không đều xem được)
Body: Không có (để trống).
Các tham số lọc (Query Params điền trên URL):
page: Trang (vd: 1)
size: Số phim mỗi trang (vd: 10)
search: Tìm theo tên phim (vd: Spider, Batman, Endgame)
status: Lọc theo trạng thái (SHOWING - đang chiếu, UPCOMING - sắp chiếu, ENDED - ngừng chiếu)

Ví dụ thực tế:

Xem tất cả phim trang 1:
GET http://localhost:8000/api/v1/movies?page=1&size=10
Lọc các phim đang chiếu:
GET http://localhost:8000/api/v1/movies?page=1&size=10&status=SHOWING
Tìm kiếm phim "Spider":
GET http://localhost:8000/api/v1/movies?search=Spider
🕒 B. Xem lịch chiếu (Suất chiếu)
Method: GET
URL: http://localhost:8000/api/v1/showtimes
Header: Authorization: Bearer <token> (Tùy chọn)
Body: Không có (để trống).
Các tham số lọc (Query Params điền trên URL):
movie_id: UUID của bộ phim cần xem lịch chiếu
room_id: UUID của phòng chiếu
format: Định dạng chiếu (2D, 3D, IMAX)
date_from và date_to: Lọc theo khoảng ngày (YYYY-MM-DD, vd: 2026-11-01 đến 2026-11-15)
exclude_ended: Đặt true để ẩn các suất chiếu đã kết thúc trong quá khứ

Ví dụ thực tế:

Xem tất cả các suất chiếu còn hiệu lực:
GET http://localhost:8000/api/v1/showtimes?exclude_ended=true
Xem lịch chiếu của 1 phim cụ thể (Ví dụ Avengers: 2b6c929a-fb0f-4fa4-a746-1fa5a56d9532):
GET http://localhost:8000/api/v1/showtimes?movie_id=2b6c929a-fb0f-4fa4-a746-1fa5a56d9532&exclude_ended=true
Xem các suất chiếu sắp tới (Widget trang chủ):
GET http://localhost:8000/api/v1/showtimes/upcoming?limit=5
4. Bước 3: Xem sơ đồ ghế & Tình trạng ghế của suất chiếu

Khi bấm vào một suất chiếu để đặt vé:

Lấy toàn bộ sơ đồ ghế của phòng:
Method: GET
URL: http://localhost:8000/api/v1/seats?room_id=1d26805c-8d37-4370-ad0b-1afb972b2671
Kiểm tra ghế đang bị khóa (Locked Seats) theo Suất chiếu:
Method: GET
URL: http://localhost:8000/api/v1/seats/locked?showtime_id=e301901a-da5f-4091-bf96-b08fcf572a11
5. Chi tiết các module khác trong file tài liệu

Trong file HUONG_DAN_TEST_API.md đã bổ sung đầy đủ:

Tạo & Cập nhật phim: Cú pháp POST /movies và lưu ý không được bỏ trống trường status khi PUT /movies/:id.
Tạo & Hủy suất chiếu: Cú pháp POST /showtimes, các điều kiện về start_time, end_time và ràng buộc không trùng phòng.
Quy trình Đặt vé & Thanh toán: POST /bookings ➔ POST /payments.
Soát vé tại quầy (Box Office): GET /tickets/search và PATCH /tickets/:id/mark-used.
Chatbot AI (Google Gemini RAG): POST /chatbot/chat.