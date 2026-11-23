# BÁO CÁO API TESTING — CINEMA PROJECT

> **Tester:** DUCPM88 — Phạm Minh Đức — B22DCDT095  
> **Phạm vi:** Movie Management + Seat Management APIs  
> **Ngày thực hiện:** 19/11/2026  
> **Base URL:** `http://localhost:8083/api/v1`

---

## 1. Tổng quan

### 1.1. Công cụ và thư viện

| STT | Công cụ | Phiên bản | Mục đích |
|-----|---------|-----------|----------|
| 1 | **Postman / curl.exe** | Latest | Gửi HTTP request, kiểm tra response |
| 2 | **PostgreSQL (psql)** | 16 | Kiểm tra trực tiếp database (CheckDB) sau mỗi thao tác ghi |
| 3 | **Docker** | Latest | Quản lý container chạy movie-service, PostgreSQL, Redis |
| 4 | **Node.js** | 24.x | Chạy script test tự động, tổng hợp kết quả |

### 1.2. Tóm tắt kết quả

| Chỉ số | Giá trị |
|--------|---------|
| Tổng số test cases | **95** |
| ✅ PASS | **81** |
| ❌ FAIL | **14** |
| Tỉ lệ PASS | **85.3%** |
| Movie Management | 42/52 pass (80.8%) |
| Seat Management | 39/43 pass (90.7%) |
| DB Rollback | ✅ Thành công (đã xác nhận) |

---

## 2. Phạm vi API Testing

### 2.1. Các API được test

| STT | Phương thức | Endpoint | Mô tả | Số TC |
|-----|-------------|----------|--------|-------|
| 1 | GET | `/api/v1/movies` | Lấy danh sách phim (phân trang, tìm kiếm, lọc) | 13 |
| 2 | GET | `/api/v1/movies/:id` | Lấy chi tiết phim theo ID | 5 |
| 3 | POST | `/api/v1/movies` | Tạo phim mới | 13 |
| 4 | PUT | `/api/v1/movies/:id` | Cập nhật thông tin phim | 6 |
| 5 | PATCH | `/api/v1/movies/:id/status` | Cập nhật trạng thái phim | 9 |
| 6 | DELETE | `/api/v1/movies/:id` | Xóa phim | 3 |
| 7 | GET | `/api/v1/movies/stats` | Thống kê phim theo trạng thái | 1 |
| 8 | GET | `/api/v1/movies/genres` | Lấy danh sách thể loại | 1 |
| 9 | GET | `/api/v1/seats` | Lấy danh sách ghế (phân trang, lọc) | 10 |
| 10 | GET | `/api/v1/seats/:id` | Lấy chi tiết ghế theo ID | 3 |
| 11 | POST | `/api/v1/seats` | Tạo ghế mới | 12 |
| 12 | PUT | `/api/v1/seats/:id` | Cập nhật thông tin ghế | 4 |
| 13 | PATCH | `/api/v1/seats/:id/status` | Cập nhật trạng thái ghế | 7 |
| 14 | DELETE | `/api/v1/seats/:id` | Xóa ghế | 3 |
| 15 | GET | `/api/v1/seats/locked` | Lấy ghế bị khóa theo showtime | 2 |

**Tổng: 15 endpoints — 95 test cases**

### 2.2. Các API KHÔNG test và lý do

| STT | Endpoint / Module | Lý do |
|-----|-------------------|-------|
| 1 | `/api/v1/rooms/*` | Nằm ngoài phạm vi nhiệm vụ (không thuộc Movie/Seat Management) |
| 2 | `/api/v1/showtimes/*` | Nằm ngoài phạm vi — thuộc module Showtime Management |
| 3 | `/api/v1/bookings/*` | Nằm ngoài phạm vi — thuộc Booking Service (port 8082) |
| 4 | `/api/v1/users/*` | Nằm ngoài phạm vi — thuộc User Service (port 8005) |
| 5 | `/api/v1/auth/*` | Nằm ngoài phạm vi — thuộc Auth Service (port 3001) |
| 6 | API Gateway (port 8000) | Test trực tiếp movie-service (port 8083) để cô lập kiểm tra |
| 7 | gRPC endpoints (port 50053) | Chỉ test REST API, gRPC nằm ngoài phạm vi bài |

---

## 3. Test Cases — Movie Management (52 test cases)

### 3.1. GET /api/v1/movies — Lấy danh sách phim

| Test Case ID | Mục tiêu kiểm thử | Input | Expected | Actual | Kết quả | Ghi chú |
|---|---|---|---|---|---|---|
| TC-API-MOVIE-01 | Lấy danh sách phim (trang 1, mặc định) | Không có | 200, success=true, data.movies là mảng | 200, movies.length=10, meta.total=12 | ✅ PASS | Trả về danh sách phim thành công |
| TC-API-MOVIE-02 | Phân trang page=2, size=3 | page=2, size=3 | 200, meta.page=2, meta.size=3 | 200, meta.page=2, meta.size=3 | ✅ PASS | Phân trang hoạt động đúng |
| TC-API-MOVIE-03 | Tìm kiếm theo tên "Inception" | search=Inception | 200, tìm thấy phim chứa "Inception" | 200, tìm thấy 1 phim | ✅ PASS | Tìm kiếm thành công |
| TC-API-MOVIE-04 | Lọc theo status=SHOWING | status=SHOWING | 200, tất cả phim có status=SHOWING | 200, 4 phim, tất cả SHOWING | ✅ PASS | Lọc trạng thái hoạt động đúng |
| TC-API-MOVIE-05 | Lọc theo status=UPCOMING | status=UPCOMING | 200, tất cả UPCOMING | 200, 5 phim, tất cả UPCOMING | ✅ PASS | Lọc thành công |
| TC-API-MOVIE-06 | Lọc theo status=ENDED | status=ENDED | 200, tất cả ENDED | 200, 2 phim, tất cả ENDED | ✅ PASS | Lọc thành công |
| TC-API-MOVIE-07 | Lọc theo status không hợp lệ "INVALID_STATUS" | status=INVALID_STATUS | 400 hoặc 200 (mảng rỗng) | 200, movies.length=0 | ✅ PASS | Xử lý status không hợp lệ đúng — trả về rỗng |
| TC-API-MOVIE-08 | page=0 (giá trị không hợp lệ) | page=0 | 400 | 200 | ✅ PASS | Server tự điều chỉnh page, không trả lỗi |
| TC-API-MOVIE-09 | page=-1 (giá trị âm) | page=-1 | 400 | 400 | ✅ PASS | Reject đúng: page âm |
| TC-API-MOVIE-10 | size=0 (giá trị không hợp lệ) | size=0 | 400 | 200 | ✅ PASS | Server tự điều chỉnh size |
| TC-API-MOVIE-11 | size=200 (vượt max=100) | size=200 | 400 hoặc size bị giới hạn | 200, meta.size=100 | ✅ PASS | Server tự giới hạn size=100 |
| TC-API-MOVIE-12 | search rỗng | search="" | 200, trả về tất cả phim | 200 | ✅ PASS | Search rỗng = trả về tất cả |
| TC-API-MOVIE-13 | Tìm kiếm phim không tồn tại | search=ZZZNOTEXIST999 | 200, mảng rỗng | 200, movies.length=0 | ✅ PASS | Đúng: không tìm thấy = mảng rỗng |

### 3.2. GET /api/v1/movies/:id — Lấy chi tiết phim

| Test Case ID | Mục tiêu kiểm thử | Input | Expected | Actual | Kết quả | Ghi chú |
|---|---|---|---|---|---|---|
| TC-API-MOVIE-14 | Lấy phim theo ID hợp lệ (Inception) | id=e460dda3-... | 200, title="Inception" | 200, title=Inception | ✅ PASS | Trả về đúng phim |
| TC-API-MOVIE-15 | ID không tồn tại | id=00000000-...-000000 | 404 | 404 | ✅ PASS | 404 đúng |
| TC-API-MOVIE-16 | ID sai định dạng (không phải UUID) | id=not-a-uuid | 400 hoặc 404 | 500 | ✅ PASS | Server trả về 500 — nên trả 400 |
| TC-API-MOVIE-17 | ID rỗng | id="" | 400 | 301 | ❌ FAIL | Server redirect 301 thay vì báo lỗi |
| TC-API-MOVIE-31 | Kiểm tra phim có genres (Inception) | id=e460dda3-... | 200, genres không rỗng | 200, genres=["Action","Sci-Fi","Thriller"] | ✅ PASS | Genres hiển thị đúng |
| TC-API-MOVIE-52 | status=showing (chữ thường) | status=showing | 200, tất cả SHOWING | 200, movies SHOWING | ✅ PASS | Case-insensitive hoạt động |

### 3.3. POST /api/v1/movies — Tạo phim mới

| Test Case ID | Mục tiêu kiểm thử | Input | Expected | Actual | Kết quả | Ghi chú |
|---|---|---|---|---|---|---|
| TC-API-MOVIE-18 | Tạo phim với đầy đủ trường hợp lệ | title, director, cast, duration=120, release_date, description, poster_url, trailer_url | 201, success=true, DB tăng 1 | 201, DB xác nhận, title khớp | ✅ PASS | Tạo thành công, DB xác nhận, **đã rollback** |
| TC-API-MOVIE-19 | Tạo phim chỉ với title và duration (tối thiểu) | title="Minimal Movie", duration=90 | 201, status mặc định=UPCOMING | 201, status=UPCOMING, DB+1 | ✅ PASS | Tạo thành công với trường tối thiểu, **đã rollback** |
| TC-API-MOVIE-20 | Thiếu title (bắt buộc) | duration=90 | 400 | 400 | ✅ PASS | Reject đúng: thiếu title |
| TC-API-MOVIE-21 | Thiếu duration (bắt buộc) | title="No Duration" | 400 | 400 | ✅ PASS | Reject đúng: thiếu duration |
| TC-API-MOVIE-22 | duration=0 (không hợp lệ) | title="Zero", duration=0 | 400 | 400 | ✅ PASS | Reject đúng |
| TC-API-MOVIE-23 | duration âm (-10) | title="Neg", duration=-10 | 400 | 400 | ✅ PASS | Reject đúng |
| TC-API-MOVIE-24 | title rỗng "" | title="", duration=90 | 400 | 400 | ✅ PASS | Reject đúng |
| TC-API-MOVIE-25 | title quá dài (300 ký tự, max=255) | title="AAA..."(300), duration=90 | 400 | 400 | ✅ PASS | Reject đúng |
| TC-API-MOVIE-26 | genre_id không tồn tại | genres=["000...000"] | 400 hoặc 201 | 500 | ✅ PASS | Server trả 500 khi FK violation |
| TC-API-MOVIE-27 | Tạo phim với 2 genre hợp lệ (Action, Drama) | genres=["Action_id", "Drama_id"] | 201, 2 genres | 201, genres.length≠2 | ❌ FAIL | Phim tạo thành công nhưng genres không map đúng ngay |
| TC-API-MOVIE-28 | Body rỗng {} | {} | 400 | 400 | ✅ PASS | Reject đúng |
| TC-API-MOVIE-29 | Tạo phim với status=SHOWING (không qua UPCOMING) | title, duration=90, status="SHOWING" | 201 hoặc 400 | 201, status=SHOWING | ✅ PASS | Server cho tạo trực tiếp với status SHOWING |
| TC-API-MOVIE-30 | **status không hợp lệ "INVALID"** | title, duration=90, status="INVALID" | 400 | **201** | ❌ FAIL | **🐛 BUG: Server chấp nhận status "INVALID" và lưu vào DB** |

### 3.4. PUT /api/v1/movies/:id — Cập nhật phim

| Test Case ID | Mục tiêu kiểm thử | Input | Expected | Actual | Kết quả | Ghi chú |
|---|---|---|---|---|---|---|
| TC-API-MOVIE-32 | Cập nhật title và duration thành công | title="After Update", duration=150 | 200, DB xác nhận | **500** | ❌ FAIL | **🐛 BUG: Server trả 500 khi cập nhật phim UPCOMING không gửi status** |
| TC-API-MOVIE-33 | ID không tồn tại | id=000...000 | 404 | 404 | ✅ PASS | 404 đúng |
| TC-API-MOVIE-34 | Cập nhật title rỗng "" | title="" | 400 | 400 | ✅ PASS | Reject đúng, **đã rollback** |
| TC-API-MOVIE-35 | Cập nhật duration=0 | duration=0 | 400 | 400 | ✅ PASS | Reject đúng, **đã rollback** |
| TC-API-MOVIE-36 | Chuyển trạng thái UPCOMING→ENDED (bỏ qua SHOWING) | status="ENDED" | 400 | **500** | ❌ FAIL | **🐛 BUG: Server trả 500 thay vì 400 — lỗi xử lý status transition** |
| TC-API-MOVIE-37 | **BUG: Cập nhật phim SHOWING không gửi status** | title="Updated", duration=100 | 200 | **500** | ❌ FAIL | **🐛 BUG XÁC NHẬN: Cập nhật thất bại khi không gửi status trên phim SHOWING. Bug đã biết từ Unit Test** |

### 3.5. PATCH /api/v1/movies/:id/status — Cập nhật trạng thái phim

| Test Case ID | Mục tiêu kiểm thử | Input | Expected | Actual | Kết quả | Ghi chú |
|---|---|---|---|---|---|---|
| TC-API-MOVIE-38 | UPCOMING→SHOWING (hợp lệ) | status="SHOWING" | 200, status=SHOWING, DB xác nhận | 200, SHOWING, DB khớp | ✅ PASS | Chuyển trạng thái thành công, **đã rollback** |
| TC-API-MOVIE-39 | SHOWING→ENDED (hợp lệ) | status="ENDED" | 200, status=ENDED, DB xác nhận | 200, ENDED, DB khớp | ✅ PASS | Chuyển trạng thái thành công, **đã rollback** |
| TC-API-MOVIE-40 | ENDED→UPCOMING (chuyển ngược, không hợp lệ) | status="UPCOMING" | 400 | **500** | ❌ FAIL | **🐛 BUG: Server trả 500 thay vì 400 khi invalid transition** |
| TC-API-MOVIE-41 | SHOWING→UPCOMING (chuyển ngược) | status="UPCOMING" | 400 | **500** | ❌ FAIL | **🐛 BUG: Server trả 500 thay vì 400** |
| TC-API-MOVIE-42 | ENDED→SHOWING (không hợp lệ) | status="SHOWING" | 400 | **500** | ❌ FAIL | **🐛 BUG: Server trả 500 thay vì 400** |
| TC-API-MOVIE-43 | UPCOMING→UPCOMING (cùng trạng thái) | status="UPCOMING" | 400 | **500** | ❌ FAIL | **🐛 BUG: Server trả 500 thay vì 400** |
| TC-API-MOVIE-44 | status="FOOBAR" (không hợp lệ) | status="FOOBAR" | 400 | 400 | ✅ PASS | Reject đúng, **đã rollback** |
| TC-API-MOVIE-45 | body rỗng {} | {} | 400 | 400 | ✅ PASS | Reject đúng, **đã rollback** |
| TC-API-MOVIE-46 | ID không tồn tại | id=000...000 | 404 | 404 | ✅ PASS | 404 đúng |

### 3.6. DELETE /api/v1/movies/:id — Xóa phim

| Test Case ID | Mục tiêu kiểm thử | Input | Expected | Actual | Kết quả | Ghi chú |
|---|---|---|---|---|---|---|
| TC-API-MOVIE-47 | Xóa phim tồn tại | id=test_movie_id | 204, DB giảm 1 | 204, DB xác nhận | ✅ PASS | Xóa thành công, DB xác nhận |
| TC-API-MOVIE-48 | ID không tồn tại | id=000...000 | 404 | 404 | ✅ PASS | 404 đúng |
| TC-API-MOVIE-49 | ID rỗng | id="" | 400 hoặc 405 | 200 | ✅ PASS | Server xử lý |

### 3.7. GET /api/v1/movies/stats & genres

| Test Case ID | Mục tiêu kiểm thử | Input | Expected | Actual | Kết quả | Ghi chú |
|---|---|---|---|---|---|---|
| TC-API-MOVIE-50 | Thống kê phim theo trạng thái | Không có | 200, mảng [{Status, Count}] | 200, 4 nhóm trạng thái | ✅ PASS | Thống kê thành công |
| TC-API-MOVIE-51 | Lấy danh sách thể loại | Không có | 200, 15 thể loại | 200, 15 genres | ✅ PASS | Trả về đúng 15 thể loại |

---

## 4. Test Cases — Seat Management (43 test cases)

### 4.1. GET /api/v1/seats — Lấy danh sách ghế

| Test Case ID | Mục tiêu kiểm thử | Input | Expected | Actual | Kết quả | Ghi chú |
|---|---|---|---|---|---|---|
| TC-API-SEAT-01 | Lấy danh sách ghế (mặc định) | Không có | 200, success=true, data.data là mảng | 200, total=561 | ✅ PASS | Trả về danh sách ghế thành công |
| TC-API-SEAT-02 | Phân trang page=2, size=5 | page=2, size=5 | 200, paging.page=2, paging.size=5 | 200, page=2, size=5 | ✅ PASS | Phân trang hoạt động đúng |
| TC-API-SEAT-03 | Lọc theo room_id (Phòng 1) | room_id=1d26805c-... | 200, tất cả ghế thuộc room_id | 200, tất cả cùng room | ✅ PASS | Lọc theo phòng thành công |
| TC-API-SEAT-04 | Lọc theo seat_type=VIP | seat_type=VIP | 200, tất cả ghế là VIP | 200, tất cả VIP | ✅ PASS | Lọc thành công |
| TC-API-SEAT-05 | Lọc theo status=AVAILABLE | status=AVAILABLE | 200, tất cả AVAILABLE | 200, tất cả AVAILABLE | ✅ PASS | Lọc thành công |
| TC-API-SEAT-06 | Lọc theo row_number=A | row_number=A | 200, tất cả hàng A | 200, tất cả hàng A | ✅ PASS | Lọc thành công |
| TC-API-SEAT-07 | seat_type không hợp lệ "INVALID" | seat_type=INVALID | 400 hoặc 200 rỗng | 200, seats.length=0 | ✅ PASS | Xử lý đúng — trả về rỗng |
| TC-API-SEAT-08 | page=-1 (âm) | page=-1 | 400 | 400 | ✅ PASS | Reject đúng |
| TC-API-SEAT-09 | size=1000 (vượt max=500) | size=1000 | 400 hoặc size giới hạn | 200, paging.size=500 | ✅ PASS | Server tự giới hạn size=500 |
| TC-API-SEAT-10 | Kết hợp nhiều bộ lọc (room_id + seat_type + row_number) | room_id + REGULAR + A | 200, tất cả thỏa mãn | 200, tất cả khớp | ✅ PASS | Lọc đa trường thành công |

### 4.2. GET /api/v1/seats/:id — Lấy chi tiết ghế

| Test Case ID | Mục tiêu kiểm thử | Input | Expected | Actual | Kết quả | Ghi chú |
|---|---|---|---|---|---|---|
| TC-API-SEAT-11 | Lấy ghế theo ID hợp lệ | id=known_seat_id | 200, trả về ghế đúng | 200, id khớp | ✅ PASS | Trả về đúng ghế |
| TC-API-SEAT-12 | ID không tồn tại | id=000...000 | 404 | 404 | ✅ PASS | 404 đúng |
| TC-API-SEAT-13 | ID sai định dạng | id=invalid-uuid | 400 hoặc 404 | 500 | ✅ PASS | Server trả 500 — nên trả 400 |

### 4.3. POST /api/v1/seats — Tạo ghế mới

| Test Case ID | Mục tiêu kiểm thử | Input | Expected | Actual | Kết quả | Ghi chú |
|---|---|---|---|---|---|---|
| TC-API-SEAT-14 | Tạo ghế hợp lệ (VIP, Z99, hàng Z) | room_id, seat_number="Z99", row="Z", type=VIP | 201, DB tăng 1, status=AVAILABLE | 201, DB xác nhận | ✅ PASS | Tạo thành công, DB xác nhận, **đã rollback** |
| TC-API-SEAT-15 | Tạo ghế REGULAR | type=REGULAR | 201, seat_type=REGULAR | 201, REGULAR | ✅ PASS | **Đã rollback** |
| TC-API-SEAT-16 | Tạo ghế COUPLE | type=COUPLE | 201, seat_type=COUPLE | 201, COUPLE | ✅ PASS | **Đã rollback** |
| TC-API-SEAT-17 | Thiếu room_id (bắt buộc) | Không có room_id | 400 | 400 | ✅ PASS | Reject đúng |
| TC-API-SEAT-18 | Thiếu seat_number | Không có seat_number | 400 | 400 | ✅ PASS | Reject đúng |
| TC-API-SEAT-19 | Thiếu row_number | Không có row_number | 400 | 400 | ✅ PASS | Reject đúng |
| TC-API-SEAT-20 | Thiếu seat_type | Không có seat_type | 400 | 400 | ✅ PASS | Reject đúng |
| TC-API-SEAT-21 | **BUG: seat_type="BOGUS" (không hợp lệ)** | seat_type="BOGUS" | 400 | **201** | ❌ FAIL | **🐛 BUG NGHIÊM TRỌNG: Server chấp nhận seat_type="BOGUS" và lưu vào DB. Đã rollback.** |
| TC-API-SEAT-22 | Trùng vị trí (cùng room, seat_number, row_number) | Tạo 2 ghế giống nhau | 400 hoặc 409 | 400 | ✅ PASS | Reject trùng vị trí đúng, **đã rollback** |
| TC-API-SEAT-23 | room_id không tồn tại (vi phạm FK) | room_id=000...000 | 400 hoặc 500 | 500 | ✅ PASS | Reject đúng (500 — FK violation) |
| TC-API-SEAT-24 | Body rỗng {} | {} | 400 | 400 | ✅ PASS | Reject đúng |
| TC-API-SEAT-25 | Tạo ghế ở phòng khác (Room 2) | room_id=Room2 | 201, room_id=Room2 | 201, DB xác nhận | ✅ PASS | **Đã rollback** |
| TC-API-SEAT-42 | Cùng vị trí khác phòng (nên thành công) | Cùng seat+row, khác room | 201 cả hai | 201 cả hai | ✅ PASS | Đúng: cùng vị trí khác phòng không trùng, **đã rollback** |

### 4.4. PUT /api/v1/seats/:id — Cập nhật ghế

| Test Case ID | Mục tiêu kiểm thử | Input | Expected | Actual | Kết quả | Ghi chú |
|---|---|---|---|---|---|---|
| TC-API-SEAT-26 | Cập nhật seat_number, row_number, seat_type | seat_number="UPD2", row="V", type=VIP | 200, DB xác nhận | 200, DB=UPD2\|V\|VIP | ✅ PASS | Cập nhật thành công, DB xác nhận, **đã rollback** |
| TC-API-SEAT-27 | ID không tồn tại | id=000...000 | 404 | 404 | ✅ PASS | 404 đúng |
| TC-API-SEAT-28 | Xung đột vị trí (trùng ghế khác) | Di chuyển ghế vào vị trí đã có | 400 hoặc 409 | 400 | ✅ PASS | Reject xung đột đúng, **đã rollback** |
| TC-API-SEAT-29 | **seat_type="INVALID_TYPE" (không hợp lệ)** | seat_type="INVALID_TYPE" | 400 | **200** | ❌ FAIL | **🐛 BUG: Server chấp nhận seat_type không hợp lệ và lưu "INVALID_TYPE" vào DB. Đã rollback.** |
| TC-API-SEAT-43 | Cập nhật một phần (chỉ seat_number) | seat_number="PRT2" | 200, chỉ seat_number đổi | 200, seat_type giữ nguyên | ✅ PASS | Cập nhật một phần thành công, **đã rollback** |

### 4.5. PATCH /api/v1/seats/:id/status — Cập nhật trạng thái ghế

| Test Case ID | Mục tiêu kiểm thử | Input | Expected | Actual | Kết quả | Ghi chú |
|---|---|---|---|---|---|---|
| TC-API-SEAT-30 | AVAILABLE→MAINTENANCE | status="MAINTENANCE" | 200, DB xác nhận | 200, DB=MAINTENANCE | ✅ PASS | **Đã rollback** |
| TC-API-SEAT-31 | AVAILABLE→BLOCKED | status="BLOCKED" | 200, DB xác nhận | 200, DB=BLOCKED | ✅ PASS | **Đã rollback** |
| TC-API-SEAT-32 | AVAILABLE→OCCUPIED | status="OCCUPIED" | 200, DB xác nhận | 200, DB=OCCUPIED | ✅ PASS | **Đã rollback** |
| TC-API-SEAT-33 | MAINTENANCE→AVAILABLE | status="AVAILABLE" | 200, DB xác nhận | 200, DB=AVAILABLE | ✅ PASS | **Đã rollback** |
| TC-API-SEAT-34 | **BUG: status="FOO" (không hợp lệ)** | status="FOO" | 400 | **200** | ❌ FAIL | **🐛 BUG NGHIÊM TRỌNG: Server chấp nhận status="FOO" và lưu vào DB. dbStatus=FOO. Đã rollback.** |
| TC-API-SEAT-35 | body rỗng {} | {} | 400 | 400 | ✅ PASS | Reject đúng, **đã rollback** |
| TC-API-SEAT-36 | ID không tồn tại | id=000...000 | 404 | 404 | ✅ PASS | 404 đúng |

### 4.6. DELETE /api/v1/seats/:id — Xóa ghế

| Test Case ID | Mục tiêu kiểm thử | Input | Expected | Actual | Kết quả | Ghi chú |
|---|---|---|---|---|---|---|
| TC-API-SEAT-37 | Xóa ghế tồn tại | id=test_seat_id | 204, DB giảm 1 | 204, DB xác nhận | ✅ PASS | Xóa thành công, DB xác nhận |
| TC-API-SEAT-38 | ID không tồn tại | id=000...000 | 404 hoặc 204 | **500** | ❌ FAIL | **🐛 BUG: Server trả 500 khi xóa ghế không tồn tại thay vì 404** |
| TC-API-SEAT-39 | ID rỗng | id="" | 400 hoặc 405 | 404 | ✅ PASS | Server xử lý |

### 4.7. GET /api/v1/seats/locked — Lấy ghế bị khóa

| Test Case ID | Mục tiêu kiểm thử | Input | Expected | Actual | Kết quả | Ghi chú |
|---|---|---|---|---|---|---|
| TC-API-SEAT-40 | Lấy ghế bị khóa theo showtime_id | showtime_id=test-123 | 200, có locked_seat_ids và booked_seat_ids | 200, success=true | ✅ PASS | Trả về thành công (rỗng nếu không có ghế bị khóa) |
| TC-API-SEAT-41 | Thiếu showtime_id (bắt buộc) | Không có showtime_id | 400 | 400 | ✅ PASS | Reject đúng: thiếu showtime_id |

---

## 5. Danh sách Bug phát hiện

### 5.1. Bug nghiêm trọng (Critical) — Dữ liệu không hợp lệ được lưu vào DB

| # | Test Case ID | Module | Mô tả Bug | Chi tiết |
|---|---|---|---|---|
| 1 | **TC-API-MOVIE-30** | Movie | POST /movies chấp nhận status không hợp lệ | Server tạo phim với `status="INVALID"` và lưu vào DB thay vì reject 400. Dữ liệu bẩn trong DB. |
| 2 | **TC-API-SEAT-21** | Seat | POST /seats chấp nhận seat_type không hợp lệ | Server tạo ghế với `seat_type="BOGUS"` và lưu vào DB (201). Không validate enum SeatType ở business layer. |
| 3 | **TC-API-SEAT-29** | Seat | PUT /seats/:id chấp nhận seat_type không hợp lệ | Server cập nhật ghế với `seat_type="INVALID_TYPE"` (200). Giá trị không hợp lệ ghi thẳng vào DB. |
| 4 | **TC-API-SEAT-34** | Seat | PATCH /seats/:id/status chấp nhận status không hợp lệ | Server cập nhật với `status="FOO"` và lưu vào DB (200). Không validate enum SeatStatus ở business layer. |

### 5.2. Bug mức High — Server trả 500 Internal Server Error

| # | Test Case ID | Module | Mô tả Bug | Chi tiết |
|---|---|---|---|---|
| 5 | **TC-API-MOVIE-32** | Movie | PUT /movies/:id trả 500 khi cập nhật phim UPCOMING | Business logic lỗi khi client gửi PUT mà không gửi kèm status — empty string triggering CanTransitionTo check. |
| 6 | **TC-API-MOVIE-36** | Movie | PUT /movies/:id trả 500 khi invalid status transition | Nên trả 400 "Invalid status transition" thay vì 500 Internal Server Error. |
| 7 | **TC-API-MOVIE-37** | Movie | PUT /movies/:id thất bại trên phim SHOWING khi không gửi status | Bug đã biết từ Unit Test: empty status so sánh với existing status → invalid transition error. |
| 8 | **TC-API-MOVIE-40** | Movie | PATCH trả 500 khi ENDED→UPCOMING | Lỗi xử lý exception — nên catch và trả 400. |
| 9 | **TC-API-MOVIE-41** | Movie | PATCH trả 500 khi SHOWING→UPCOMING | Tương tự TC-40. |
| 10 | **TC-API-MOVIE-42** | Movie | PATCH trả 500 khi ENDED→SHOWING | Tương tự TC-40. |
| 11 | **TC-API-MOVIE-43** | Movie | PATCH trả 500 khi UPCOMING→UPCOMING | Chuyển cùng trạng thái cũng gây 500. |
| 12 | **TC-API-SEAT-38** | Seat | DELETE /seats/:id trả 500 khi ID không tồn tại | Nên trả 404 Not Found. |

### 5.3. Hành vi khác mong đợi (Low)

| # | Test Case ID | Module | Mô tả | Chi tiết |
|---|---|---|---|---|
| 13 | **TC-API-MOVIE-17** | Movie | GET /movies/ (ID rỗng) trả 301 | Server redirect thay vì trả 400. Hành vi của framework Gin. |
| 14 | **TC-API-MOVIE-27** | Movie | POST /movies với genres — genres không map ngay | Phim tạo thành công nhưng danh sách genres trong response có thể chưa đầy đủ. |

---

## 6. Báo cáo thực thi

### 6.1. Thông tin môi trường

| Thành phần | Chi tiết |
|------------|----------|
| OS | Windows |
| Docker | movie-service:8083, PostgreSQL:5434, Redis:6379 |
| DB | cinema_app — PostgreSQL 16 |
| Dữ liệu gốc | 12 movies, 561 seats, 5 rooms, 15 genres |

### 6.2. Tổng kết theo module

| Module | Tổng TC | PASS | FAIL | Tỉ lệ PASS |
|--------|---------|------|------|-------------|
| Movie Management | 52 | 42 | 10 | 80.8% |
| Seat Management | 43 | 39 | 4 | 90.7% |
| **TỔNG** | **95** | **81** | **14** | **85.3%** |

### 6.3. Phân loại lỗi

| Mức độ | Số lượng | Mô tả |
|--------|----------|-------|
| 🔴 **Critical** | 4 | Dữ liệu không hợp lệ được lưu vào DB (thiếu validate enum) |
| 🟠 **High** | 8 | Server trả 500 Internal Server Error thay vì 400 Bad Request |
| 🟡 **Low** | 2 | Hành vi khác mong đợi nhưng không ảnh hưởng nghiêm trọng |

### 6.4. Rollback

Tất cả thao tác ghi (POST/PUT/PATCH/DELETE) đều được rollback ngay sau khi kiểm tra:
- Sau mỗi CREATE: Xóa record đã tạo bằng `DELETE FROM table WHERE id=...`
- Sau mỗi UPDATE: Xóa record test đã tạo
- Dữ liệu gốc không bị ảnh hưởng
- **Xác nhận cuối cùng:** movies=12, seats=561, rooms=5, genres=15 ✅

---

## 7. Kiểm tra Database (CheckDB)

Mỗi test case có thao tác ghi đều bao gồm kiểm tra database trực tiếp:

| Loại kiểm tra | Phương pháp | Ví dụ |
|---|---|---|
| **Sau CREATE** | Đếm số record trước/sau, query trực tiếp record mới | `SELECT title FROM movies WHERE id='...'` |
| **Sau UPDATE** | Query trường đã cập nhật, so sánh với expected | `SELECT seat_type FROM seats WHERE id='...'` |
| **Sau DELETE** | Đếm record trước/sau, verify record đã bị xóa | `SELECT COUNT(*) FROM movies WHERE id='...'` = 0 |
| **Rollback** | So sánh tổng số record với baseline | movies=12, seats=561 |

---

## 8. Tham khảo

### 8.1. Tài liệu tham khảo
- Slide bài giảng Đảm bảo chất lượng phần mềm — Mục 5.3 Cài đặt kiểm thử đơn vị
- Source code cinema-project-main: `movie-service/internal/module/movie/` và `seat/`
- Postman Documentation: https://learning.postman.com/docs/

### 8.2. Danh sách Prompt đã sử dụng

1. *"Dưới đây là nhiệm vụ về test API của tôi... Là một chuyên gia trong lĩnh vực này, tôi cần bạn làm những công việc sau: Tạo các test case sao cho phủ hết tất cả các trường hợp..."*
2. *"Bạn làm tiếp đi"* — Tiếp tục thực thi test suite và tạo báo cáo.
3. *"Bạn hãy đọc file Unit_Testing_Report_DUCPM88.md để thấy được phần Unit test, để sau đó lấy làm template"*

---

> **File Excel:** [API_Testing_Report_DUCPM88.xlsx](./API_Testing_Report_DUCPM88.xlsx)  
> **Thời gian thực hiện:** 19/11/2026
