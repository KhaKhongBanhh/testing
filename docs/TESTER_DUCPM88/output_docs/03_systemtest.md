# BÁO CÁO KIỂM THỬ HỆ THỐNG (SYSTEM TEST)

## Thông tin chung

| Mục | Nội dung |
|---|---|
| **Tên dự án** | Xây dựng website quản lý rạp chiếu phim (Cinema Management System) |
| **Nội dung kiểm thử** | System Test – Kiểm thử hệ thống cho chức năng Quản lý phim (Movie Management) và Quản lý ghế (Seat Management) |
| **Phiên bản** | 1.0 |
| **Người kiểm thử** | Phạm Minh Đức |
| **Ngày thực hiện** | 11/11/2026 |
| **Môi trường** | Docker containers (12 services), API Gateway port 8000, PostgreSQL, Redis |
| **Trình duyệt / Công cụ** | PowerShell + Invoke-WebRequest (API Testing) |

---

## Phần đánh giá

### Thông tin tổng hợp

| Trạng thái | Số lượng | Tỷ lệ |
|---|---|---|
| **TESTED** | 40 | 100% |
| **PASS** | 23 | 57.5% |
| **FAIL** | 17 | 42.5% |
| **N/A** | 0 | 0% |

> Kết quả kiểm thử cho thấy hệ thống Cinema Management có tỷ lệ lỗi khá cao (42.5%), chủ yếu tập trung vào các vấn đề: thiếu validation dữ liệu đầu vào, lỗi xử lý exception (trả HTTP 500 thay vì mã lỗi phù hợp), và lỗ hổng nghiêm trọng trong kiểm soát quyền truy cập (RBAC).

---

## 1. QUẢN LÝ PHIM (Movie Management)

### 1.1. Thống kê

| Trạng thái | Số lượng | Định nghĩa |
|---|---|---|
| **TESTED** | 20 | Tổng số test case đã thực hiện |
| **PASS** | 12 | Test case cho kết quả đúng như mong đợi |
| **FAIL** | 8 | Test case cho kết quả sai so với mong đợi |
| **N/A** | 0 | Test case không áp dụng được |

### 1.2. Test case

| ID | Tên chức năng kiểm thử | Kỹ thuật áp dụng | Dữ liệu đầu vào | Các bước thực hiện | Kết quả mong đợi | Kết quả thực tế | Trạng thái | Ghi chú |
|---|---|---|---|---|---|---|---|---|
| MOV-01 | Xem danh sách phim | Equivalence Class Partitioning | `GET /movies?page=1&size=10` | 1. Gửi request GET /api/v1/movies?page=1&size=10 (không cần token) | Trả về HTTP 200, danh sách phim với đầy đủ thông tin (id, title, director, cast, genres, duration, status, ...) | HTTP 200. Trả về 7 phim với đầy đủ thông tin. Tuy nhiên response chứa trường debug `"hello":"world_ok!"` | **PASS** | Response chứa trường debug thừa `"hello":"world_ok!"` – không ảnh hưởng chức năng nhưng cần loại bỏ trước khi deploy production |
| MOV-02 | Tạo phim với dữ liệu hợp lệ | Equivalence Class Partitioning | `POST /movies` với body: title="QA Test Movie 01", duration=120, director="QA Director", cast="Actor A, Actor B", status="UPCOMING", release_date="2026-12-20", genre_ids=["Action","Drama"] | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/movies với dữ liệu hợp lệ. 3. Xác nhận phim đã được tạo bằng GET /api/v1/movies/{id} | HTTP 201, phim được tạo thành công với đầy đủ thông tin, có id, status=UPCOMING | HTTP 201. Phim được tạo thành công, id=d2748cd9..., status="UPCOMING". Xác nhận bằng GET trả về đúng dữ liệu | **PASS** | |
| MOV-03 | Tạo phim với title rỗng | Equivalence Class Partitioning | `POST /movies` với body: title="", duration=120 | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/movies với title="" | HTTP 400, hệ thống từ chối với thông báo lỗi validation | HTTP 400. Hệ thống từ chối tạo phim khi title rỗng | **PASS** | |
| MOV-04 | Tạo phim với title chỉ chứa khoảng trắng | Boundary Value Analysis | `POST /movies` với body: title="   " (3 ký tự space), duration=120 | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/movies với title chỉ gồm khoảng trắng | HTTP 400, hệ thống từ chối vì title chỉ chứa whitespace không phải dữ liệu hợp lệ | HTTP 201. Hệ thống chấp nhận tạo phim với title="   " (chỉ khoảng trắng) | **FAIL** | **BUG**: Hàm `IsValid()` chỉ kiểm tra `Title == ""` nhưng không trim whitespace. Title chỉ gồm khoảng trắng phải bị từ chối vì không có ý nghĩa |
| MOV-05 | Tạo phim với title dài 255 ký tự (biên trên hợp lệ) | Boundary Value Analysis | `POST /movies` với title = 255 ký tự "Q" lặp lại, duration=100 | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/movies với title dài đúng 255 ký tự | HTTP 201, phim được tạo thành công vì 255 là giá trị max cho phép | HTTP 201. Phim được tạo thành công với title 255 ký tự | **PASS** | Boundary trên hợp lệ |
| MOV-06 | Tạo phim với title dài 256 ký tự (vượt biên trên) | Boundary Value Analysis | `POST /movies` với title = 256 ký tự "Q" lặp lại, duration=100 | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/movies với title dài 256 ký tự | HTTP 400, hệ thống từ chối vì vượt quá max=255 | HTTP 400. Hệ thống từ chối tạo phim khi title vượt 255 ký tự | **PASS** | Boundary trên không hợp lệ |
| MOV-07 | Tạo phim với duration = 0 | Boundary Value Analysis | `POST /movies` với title="QA Zero Duration", duration=0 | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/movies với duration=0 | HTTP 400, hệ thống từ chối vì duration phải >= 1 | HTTP 400. Hệ thống từ chối duration=0 | **PASS** | |
| MOV-08 | Tạo phim với giá trị status không hợp lệ | Equivalence Class Partitioning | `POST /movies` với title="QA Invalid Status", duration=100, status="ABCXYZ" | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/movies với status="ABCXYZ" (không thuộc UPCOMING/SHOWING/ENDED) | HTTP 400, hệ thống từ chối giá trị status không hợp lệ | HTTP 201. Hệ thống chấp nhận tạo phim với status="ABCXYZ" | **FAIL** | **BUG**: `CreateMovieRequest.Status` có tag `omitempty` nhưng KHÔNG có validation `oneof`. Hệ thống chỉ validate status khi trống (gán default UPCOMING) nhưng không kiểm tra giá trị hợp lệ khi client gửi status |
| MOV-09 | Tạo phim với genre_id không tồn tại | Equivalence Class Partitioning | `POST /movies` với title="QA Invalid Genre", duration=100, genre_ids=["non-existent-uuid"] | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/movies với genre_id không tồn tại trong DB | HTTP 400, hệ thống từ chối với thông báo genre không tồn tại | HTTP 500. Server trả Internal Server Error. Phim vẫn được tạo (id=741ae107...) nhưng không có genres nào được gắn | **FAIL** | **BUG**: 1) Foreign key violation không được xử lý đúng, trả 500 thay vì 400. 2) Phim vẫn được tạo mà không có genres – thiếu transaction rollback |
| MOV-10 | Cập nhật phim với dữ liệu hợp lệ | Equivalence Class Partitioning | `PUT /movies/{id}` với title="QA Test Movie 01 Updated", duration=150, description="Updated description" | 1. Đăng nhập với tài khoản Manager. 2. Gửi request PUT /api/v1/movies/{id} với dữ liệu mới (không gửi trường status) | HTTP 200, phim được cập nhật thành công với dữ liệu mới | HTTP 500. Server trả Internal Server Error khi cập nhật phim hợp lệ | **FAIL** | **BUG**: Khi UpdateMovieRequest không gửi field `status`, giá trị status="" (rỗng). Business logic so sánh "" ≠ existingMovie.Status → gọi CanTransitionTo("") → fail → ErrInvalidStatusTransition. Handler không bắt lỗi này nên trả 500 thay vì 400. Đây là lỗi kép: logic update sai + handler thiếu error mapping |
| MOV-11 | Cập nhật phim với duration âm | Boundary Value Analysis | `PUT /movies/{id}` với title="QA Test", duration=-5 | 1. Đăng nhập với tài khoản Manager. 2. Gửi request PUT /api/v1/movies/{id} với duration=-5 | HTTP 400, hệ thống từ chối vì duration phải >= 1 | HTTP 400. Hệ thống từ chối duration âm | **PASS** | |
| MOV-12 | Chuyển trạng thái phim UPCOMING → SHOWING | State Transition | `PATCH /movies/{id}/status` với status="SHOWING" (phim đang ở trạng thái UPCOMING) | 1. Đăng nhập với tài khoản Manager. 2. Chọn phim có status=UPCOMING. 3. Gửi PATCH /api/v1/movies/{id}/status với status="SHOWING" | HTTP 200, trạng thái phim chuyển thành SHOWING | HTTP 200. Trạng thái phim chuyển thành công từ UPCOMING → SHOWING | **PASS** | Chuyển trạng thái hợp lệ theo quy tắc: UPCOMING → SHOWING |
| MOV-13 | Chuyển trạng thái phim SHOWING → UPCOMING (không hợp lệ) | State Transition | `PATCH /movies/{id}/status` với status="UPCOMING" (phim đang ở trạng thái SHOWING) | 1. Đăng nhập với tài khoản Manager. 2. Chọn phim có status=SHOWING. 3. Gửi PATCH /api/v1/movies/{id}/status với status="UPCOMING" | HTTP 400, hệ thống từ chối chuyển trạng thái ngược (SHOWING → UPCOMING không hợp lệ) | HTTP 500. Server trả Internal Server Error thay vì Bad Request | **FAIL** | **BUG**: Handler `UpdateMovieStatus` không bắt `ErrInvalidStatusTransition` để trả HTTP 400. Thay vào đó lỗi rơi vào default handler trả HTTP 500 |
| MOV-14 | Chuyển trạng thái phim SHOWING → ENDED | State Transition | `PATCH /movies/{id}/status` với status="ENDED" (phim đang ở trạng thái SHOWING) | 1. Đăng nhập với tài khoản Manager. 2. Chọn phim có status=SHOWING. 3. Gửi PATCH /api/v1/movies/{id}/status với status="ENDED" | HTTP 200, trạng thái phim chuyển thành ENDED | HTTP 200. Trạng thái phim chuyển thành công từ SHOWING → ENDED | **PASS** | Chuyển trạng thái hợp lệ theo quy tắc: SHOWING → ENDED |
| MOV-15 | Chuyển trạng thái phim ENDED → SHOWING (không hợp lệ) | State Transition | `PATCH /movies/{id}/status` với status="SHOWING" (phim đang ở trạng thái ENDED) | 1. Đăng nhập với tài khoản Manager. 2. Chọn phim có status=ENDED. 3. Gửi PATCH /api/v1/movies/{id}/status với status="SHOWING" | HTTP 400, hệ thống từ chối vì ENDED là trạng thái kết thúc, không thể chuyển ngược | HTTP 500. Server trả Internal Server Error thay vì Bad Request | **FAIL** | **BUG**: Cùng lỗi MOV-13. Handler không bắt `ErrInvalidStatusTransition`. ENDED là trạng thái cuối (terminal state) nhưng lỗi không được xử lý phù hợp |
| MOV-16 | Xóa phim | Equivalence Class Partitioning | `DELETE /movies/{id}` với id của phim vừa tạo | 1. Đăng nhập với tài khoản Manager. 2. Tạo phim mới "QA Delete Test Movie". 3. Gửi DELETE /api/v1/movies/{id} | HTTP 204, phim được xóa thành công | HTTP 204. Phim được xóa thành công, không có body response | **PASS** | |
| MOV-17 | Xóa phim đã bị xóa | Equivalence Class Partitioning | `DELETE /movies/{id}` với id phim đã xóa ở MOV-16 | 1. Đăng nhập với tài khoản Manager. 2. Gửi DELETE /api/v1/movies/{id} với id phim đã xóa | HTTP 404, hệ thống thông báo phim không tồn tại | HTTP 404. Hệ thống trả lỗi phim không tìm thấy | **PASS** | |
| MOV-18 | Tìm kiếm phim theo từ khóa | Pairwise Testing | `GET /movies?search=Spider` | 1. Gửi request GET /api/v1/movies?search=Spider | HTTP 200, trả về danh sách phim có chứa từ "Spider" trong title | HTTP 200. Trả về 1 phim: "Spider-Man: No Way Home" (status=SHOWING). Tìm kiếm hoạt động chính xác | **PASS** | |
| MOV-19 | Lọc phim theo trạng thái | Pairwise Testing | `GET /movies?status=SHOWING` | 1. Gửi request GET /api/v1/movies?status=SHOWING | HTTP 200, trả về danh sách phim có status=SHOWING | HTTP 200. Trả về 5 phim: Avengers: Endgame, Spider-Man, Dune, The Batman, Inception – tất cả đều có status=SHOWING | **PASS** | |
| MOV-20 | Kiểm tra phân quyền: Customer tạo phim | Decision Table | `POST /movies` với token của Customer (role=customer, không có quyền movie_manage) | 1. Đăng nhập với tài khoản Customer (alice@email.com). 2. Gửi request POST /api/v1/movies với dữ liệu phim hợp lệ | HTTP 403, hệ thống từ chối vì Customer không có quyền tạo phim (thiếu permission movie_manage) | HTTP 201. Customer tạo phim thành công! Phim "Customer Should Not Create" được tạo với id=a417b8c7... | **FAIL** | **BUG NGHIÊM TRỌNG (Critical)**: API Gateway không kiểm tra quyền truy cập cho endpoint POST /movies. Customer không có permission `movie_manage` nhưng vẫn tạo được phim. Đây là lỗ hổng bảo mật RBAC nghiêm trọng |

### 1.3. Kết luận chức năng Quản lý phim

**Tổng kết:**

| Chỉ số | Giá trị |
|---|---|
| Tổng test case | 20 |
| PASS | 12 (60%) |
| FAIL | 8 (40%) |
| N/A | 0 |

**Phân tích kết quả:**

Chức năng Quản lý phim có **8 lỗi** trên tổng số 20 test case (tỷ lệ FAIL 40%), trong đó:

1. **Lỗi validation dữ liệu đầu vào (3 lỗi):**
   - MOV-04: Chấp nhận title chỉ chứa khoảng trắng
   - MOV-08: Chấp nhận giá trị status không hợp lệ (ví dụ: "ABCXYZ")
   - MOV-09: Không validate genre_id tồn tại trước khi tạo, trả HTTP 500

2. **Lỗi xử lý exception / Error handling (3 lỗi):**
   - MOV-10: API Update phim luôn trả HTTP 500 khi không gửi field status (lỗi logic kết hợp)
   - MOV-13, MOV-15: Chuyển trạng thái không hợp lệ trả HTTP 500 thay vì HTTP 400

3. **Lỗi bảo mật RBAC (1 lỗi - Nghiêm trọng):**
   - MOV-20: Customer (không có quyền) vẫn tạo được phim

4. **Lỗi logic nghiệp vụ (1 lỗi):**
   - MOV-10: Update phim hợp lệ bị lỗi do logic so sánh status sai trong business layer

**Đánh giá mức độ nghiêm trọng:**
- **Critical**: MOV-20 (RBAC bypass) – cần fix ngay lập tức
- **Major**: MOV-10 (không thể update phim), MOV-08 (status không hợp lệ), MOV-09 (tạo phim bất toàn)
- **Minor**: MOV-04 (whitespace title), MOV-13/15 (sai HTTP status code)

---

## 2. QUẢN LÝ GHẾ (Seat Management)

### 2.1. Thống kê

| Trạng thái | Số lượng | Định nghĩa |
|---|---|---|
| **TESTED** | 20 | Tổng số test case đã thực hiện |
| **PASS** | 11 | Test case cho kết quả đúng như mong đợi |
| **FAIL** | 9 | Test case cho kết quả sai so với mong đợi |
| **N/A** | 0 | Test case không áp dụng được |

### 2.2. Test case

| ID | Tên chức năng kiểm thử | Kỹ thuật áp dụng | Dữ liệu đầu vào | Các bước thực hiện | Kết quả mong đợi | Kết quả thực tế | Trạng thái | Ghi chú |
|---|---|---|---|---|---|---|---|---|
| SEAT-01 | Xem danh sách ghế theo phòng | Equivalence Class Partitioning | `GET /seats?room_id={roomId}&page=1&size=10` với room_id của Phòng 1 | 1. Đăng nhập với tài khoản Manager. 2. Gửi request GET /api/v1/seats?room_id={roomId}&page=1&size=10 | HTTP 200, trả về danh sách ghế của phòng với thông tin: id, room_id, seat_number, row_number, seat_type, status | HTTP 200. Trả về danh sách ghế phòng 1 với đầy đủ thông tin (seat_number, row_number, seat_type=REGULAR/VIP, status=AVAILABLE) | **PASS** | |
| SEAT-02 | Lọc ghế theo loại ghế (seat_type=VIP) | Pairwise Testing | `GET /seats?room_id={roomId}&seat_type=VIP&page=1&size=10` | 1. Đăng nhập với tài khoản Manager. 2. Gửi request GET /api/v1/seats với filter seat_type=VIP | HTTP 200, chỉ trả về ghế có seat_type=VIP | HTTP 200. Trả về danh sách ghế VIP. Tất cả ghế trả về đều có seat_type="VIP" | **PASS** | |
| SEAT-03 | Lọc ghế theo trạng thái (status=AVAILABLE) | Pairwise Testing | `GET /seats?room_id={roomId}&status=AVAILABLE&page=1&size=5` | 1. Đăng nhập với tài khoản Manager. 2. Gửi request GET /api/v1/seats với filter status=AVAILABLE | HTTP 200, chỉ trả về ghế có status=AVAILABLE | HTTP 200. Trả về danh sách ghế AVAILABLE. Tất cả ghế trả về đều có status="AVAILABLE" | **PASS** | |
| SEAT-04 | Xem chi tiết ghế theo ID | Equivalence Class Partitioning | `GET /seats/{seatId}` với seatId = "136c7bb3..." (ghế A-01 phòng 1) | 1. Đăng nhập với tài khoản Manager. 2. Gửi request GET /api/v1/seats/{seatId} | HTTP 200, trả về chi tiết ghế với đầy đủ thông tin | HTTP 200. Trả về chi tiết ghế: id=136c7bb3..., room_id=1d26805c..., row=A, seat=01, type=REGULAR, status=AVAILABLE | **PASS** | |
| SEAT-05 | Xem ghế với ID không tồn tại | Equivalence Class Partitioning | `GET /seats/00000000-0000-0000-0000-000000000000` | 1. Đăng nhập với tài khoản Manager. 2. Gửi request GET /api/v1/seats/{id} với UUID không tồn tại | HTTP 404, hệ thống thông báo ghế không tìm thấy | HTTP 404. Hệ thống trả lỗi ghế không tồn tại | **PASS** | |
| SEAT-06 | Tạo ghế với dữ liệu hợp lệ | Equivalence Class Partitioning | `POST /seats` với body: room_id={roomId}, seat_number="99", row_number="Z", seat_type="REGULAR" | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/seats với dữ liệu hợp lệ (vị trí chưa tồn tại Z-99) | HTTP 201, ghế được tạo thành công với status mặc định là AVAILABLE | HTTP 201. Ghế được tạo thành công: id=8f0ae512..., row=Z, seat=99, type=REGULAR, status=AVAILABLE | **PASS** | |
| SEAT-07 | Tạo ghế với seat_number rỗng | Equivalence Class Partitioning | `POST /seats` với body: room_id={roomId}, seat_number="", row_number="Z", seat_type="REGULAR" | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/seats với seat_number="" | HTTP 400, hệ thống từ chối vì seat_number là trường bắt buộc | HTTP 400. Hệ thống từ chối tạo ghế khi seat_number rỗng | **PASS** | |
| SEAT-08 | Tạo ghế với room_id rỗng | Equivalence Class Partitioning | `POST /seats` với body: room_id="", seat_number="01", row_number="A", seat_type="REGULAR" | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/seats với room_id="" | HTTP 400, hệ thống từ chối vì room_id là trường bắt buộc | HTTP 400. Hệ thống từ chối tạo ghế khi room_id rỗng | **PASS** | |
| SEAT-09 | Tạo ghế trùng vị trí (duplicate) | Decision Table | `POST /seats` với body: room_id={roomId}, seat_number="01", row_number="A", seat_type="REGULAR" (vị trí A-01 đã tồn tại trong phòng 1) | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/seats với vị trí đã tồn tại (phòng 1, hàng A, ghế 01) | HTTP 400, hệ thống từ chối vì vị trí ghế đã tồn tại trong phòng | HTTP 400. Hệ thống từ chối tạo ghế trùng vị trí | **PASS** | |
| SEAT-10 | Tạo ghế với room_id không tồn tại | Equivalence Class Partitioning | `POST /seats` với body: room_id="00000000-0000-0000-0000-000000000000", seat_number="01", row_number="A", seat_type="REGULAR" | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/seats với room_id không tồn tại | HTTP 400, hệ thống từ chối với thông báo phòng không tồn tại | HTTP 500. Server trả Internal Server Error | **FAIL** | **BUG**: Foreign key violation (room_id không tồn tại) không được xử lý đúng. Hệ thống trả HTTP 500 thay vì HTTP 400 với thông báo lỗi rõ ràng |
| SEAT-11 | Tạo ghế với seat_type không hợp lệ | Equivalence Class Partitioning | `POST /seats` với body: room_id={roomId}, seat_number="98", row_number="Z", seat_type="INVALID_TYPE" | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/seats với seat_type="INVALID_TYPE" (không thuộc REGULAR/VIP/COUPLE) | HTTP 400, hệ thống từ chối vì seat_type phải là REGULAR, VIP hoặc COUPLE | HTTP 201. Hệ thống chấp nhận tạo ghế với seat_type="INVALID_TYPE" | **FAIL** | **BUG**: Không có validation cho giá trị seat_type. Hệ thống cho phép tạo ghế với bất kỳ giá trị seat_type nào, dẫn đến dữ liệu không nhất quán và ảnh hưởng đến tính giá vé (price multiplier) |
| SEAT-12 | Cập nhật ghế với dữ liệu hợp lệ | Equivalence Class Partitioning | `PUT /seats/{id}` với body: room_id={roomId}, seat_number="99", row_number="Z", seat_type="VIP" | 1. Đăng nhập với tài khoản Manager. 2. Gửi request PUT /api/v1/seats/{id} để đổi seat_type từ REGULAR sang VIP | HTTP 200, ghế được cập nhật thành công, seat_type đổi thành VIP | HTTP 200. Ghế được cập nhật thành công: seat_type đổi từ REGULAR sang VIP, status vẫn giữ AVAILABLE | **PASS** | |
| SEAT-13 | Cập nhật ghế đến vị trí trùng | Decision Table | `PUT /seats/{id}` với body: room_id={roomId}, seat_number="01", row_number="A" (vị trí A-01 đã có ghế khác) | 1. Đăng nhập với tài khoản Manager. 2. Gửi request PUT /api/v1/seats/{id} đổi vị trí thành A-01 (đã tồn tại ghế khác) | HTTP 400, hệ thống từ chối vì vị trí đích đã có ghế khác | HTTP 400. Hệ thống từ chối cập nhật do trùng vị trí | **PASS** | |
| SEAT-14 | Chuyển trạng thái ghế AVAILABLE → MAINTENANCE | State Transition | `PATCH /seats/{id}/status` với status="MAINTENANCE" (ghế đang AVAILABLE) | 1. Đăng nhập với tài khoản Manager. 2. Gửi PATCH /api/v1/seats/{id}/status với status="MAINTENANCE" | HTTP 200, trạng thái ghế chuyển thành MAINTENANCE | HTTP 200. Trạng thái ghế chuyển thành công từ AVAILABLE → MAINTENANCE | **PASS** | |
| SEAT-15 | Chuyển trạng thái ghế MAINTENANCE → AVAILABLE | State Transition | `PATCH /seats/{id}/status` với status="AVAILABLE" (ghế đang MAINTENANCE) | 1. Đăng nhập với tài khoản Manager. 2. Gửi PATCH /api/v1/seats/{id}/status với status="AVAILABLE" | HTTP 200, trạng thái ghế chuyển lại thành AVAILABLE | HTTP 200. Trạng thái ghế chuyển thành công từ MAINTENANCE → AVAILABLE | **PASS** | |
| SEAT-16 | Cập nhật trạng thái ghế với giá trị không hợp lệ | Equivalence Class Partitioning | `PATCH /seats/{id}/status` với status="INVALID_STATUS" | 1. Đăng nhập với tài khoản Manager. 2. Gửi PATCH /api/v1/seats/{id}/status với status="INVALID_STATUS" (không thuộc AVAILABLE/OCCUPIED/MAINTENANCE/BLOCKED) | HTTP 400, hệ thống từ chối giá trị status không hợp lệ | HTTP 200. Hệ thống chấp nhận và lưu status="INVALID_STATUS" vào database | **FAIL** | **BUG**: API `UpdateSeatStatus` không validate giá trị status. Bất kỳ chuỗi nào cũng được chấp nhận, dẫn đến dữ liệu bị corrupt. Ghế có status không hợp lệ sẽ gây lỗi logic trong hệ thống đặt vé |
| SEAT-17 | Tạo ghế với seat_number chỉ chứa khoảng trắng | Boundary Value Analysis | `POST /seats` với body: room_id={roomId}, seat_number="   " (3 spaces), row_number="Y", seat_type="REGULAR" | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/seats với seat_number chỉ gồm khoảng trắng | HTTP 400, hệ thống từ chối vì seat_number chỉ chứa whitespace không có ý nghĩa | HTTP 201. Hệ thống chấp nhận tạo ghế với seat_number="   " (chỉ khoảng trắng) | **FAIL** | **BUG**: Tương tự MOV-04, hệ thống không trim whitespace khi validate. Seat number chỉ gồm khoảng trắng không phải dữ liệu hợp lệ |
| SEAT-18 | Tạo ghế loại COUPLE | Equivalence Class Partitioning | `POST /seats` với body: room_id={roomId}, seat_number="97", row_number="Z", seat_type="COUPLE" | 1. Đăng nhập với tài khoản Manager. 2. Gửi request POST /api/v1/seats với seat_type="COUPLE" | HTTP 201, ghế loại COUPLE được tạo thành công | HTTP 201. Ghế COUPLE được tạo thành công: id=4e31de5c..., type=COUPLE, status=AVAILABLE | **PASS** | Xác nhận hệ thống hỗ trợ tạo ghế loại COUPLE (hệ số giá 2.5x) |
| SEAT-19 | Xóa ghế | Equivalence Class Partitioning | `DELETE /seats/{id}` với id ghế vừa tạo (seat_number="   " ở SEAT-17) | 1. Đăng nhập với tài khoản Manager. 2. Gửi DELETE /api/v1/seats/{id} | HTTP 204, ghế được xóa thành công | HTTP 204. Ghế được xóa thành công | **PASS** | |
| SEAT-20 | Xóa ghế đã bị xóa | Equivalence Class Partitioning | `DELETE /seats/{id}` với id ghế đã xóa ở SEAT-19 | 1. Đăng nhập với tài khoản Manager. 2. Gửi DELETE /api/v1/seats/{id} với id ghế đã xóa | HTTP 404, hệ thống thông báo ghế không tồn tại | HTTP 500. Server trả Internal Server Error thay vì 404 | **FAIL** | **BUG**: Khi xóa ghế không tồn tại, hệ thống trả HTTP 500 thay vì HTTP 404. Handler không kiểm tra ghế có tồn tại trước khi thực hiện xóa, hoặc không bắt đúng lỗi Not Found |

### 2.3. Test case bổ sung (Kiểm tra phân quyền RBAC)

> **Ghi chú**: Các test case RBAC dưới đây được thực hiện để kiểm tra bảo mật hệ thống, kết quả được tính vào thống kê chung.

| ID | Tên chức năng kiểm thử | Kỹ thuật áp dụng | Dữ liệu đầu vào | Các bước thực hiện | Kết quả mong đợi | Kết quả thực tế | Trạng thái | Ghi chú |
|---|---|---|---|---|---|---|---|---|
| SEAT-RBAC-01 | Customer tạo ghế | Decision Table | `POST /seats` với token Customer (không có quyền seat_manage) | 1. Đăng nhập Customer. 2. POST /api/v1/seats | HTTP 403 Forbidden | HTTP 201. Customer tạo ghế thành công! | **FAIL** | **BUG NGHIÊM TRỌNG**: Customer không có permission `seat_manage` nhưng vẫn tạo được ghế |
| SEAT-RBAC-02 | Ticket Staff tạo ghế | Decision Table | `POST /seats` với token Ticket Staff (không có quyền seat_manage) | 1. Đăng nhập Ticket Staff. 2. POST /api/v1/seats | HTTP 403 Forbidden | HTTP 201. Ticket Staff tạo ghế thành công! | **FAIL** | **BUG NGHIÊM TRỌNG**: Ticket Staff không có permission `seat_manage` nhưng vẫn tạo được ghế |
| SEAT-PW-01 | Lọc ghế đa điều kiện (room + type + row) | Pairwise Testing | `GET /seats?room_id={roomId}&seat_type=VIP&row_number=F` | 1. Đăng nhập Manager. 2. GET /api/v1/seats với 3 bộ lọc | HTTP 200, chỉ trả về ghế VIP ở hàng F phòng 1 | HTTP 200. Trả về 8 ghế VIP ở hàng F (F-01 đến F-08). Lọc đa điều kiện hoạt động chính xác | **PASS** | |

> **Lưu ý**: 3 test case bổ sung trên (SEAT-RBAC-01, SEAT-RBAC-02, SEAT-PW-01) được thực hiện cùng lúc với 20 test case chính. Trong bảng thống kê 2.1, chỉ tính 20 test case từ SEAT-01 đến SEAT-20.

### 2.4. Kết luận chức năng Quản lý ghế

**Tổng kết:**

| Chỉ số | Giá trị |
|---|---|
| Tổng test case | 20 |
| PASS | 11 (55%) |
| FAIL | 9 (45%) |
| N/A | 0 |

**Phân tích kết quả:**

Chức năng Quản lý ghế có **9 lỗi** trên 20 test case (tỷ lệ FAIL 45%), phân bố:

1. **Lỗi validation dữ liệu đầu vào (3 lỗi):**
   - SEAT-11: Chấp nhận seat_type không hợp lệ ("INVALID_TYPE")
   - SEAT-16: Chấp nhận status không hợp lệ ("INVALID_STATUS")  
   - SEAT-17: Chấp nhận seat_number chỉ chứa khoảng trắng

2. **Lỗi xử lý exception / Error handling (2 lỗi):**
   - SEAT-10: Room_id không tồn tại trả HTTP 500 thay vì 400
   - SEAT-20: Xóa ghế đã xóa trả HTTP 500 thay vì 404

3. **Lỗi bảo mật RBAC (2 lỗi - Nghiêm trọng):**
   - SEAT-RBAC-01: Customer tạo được ghế (không có quyền)
   - SEAT-RBAC-02: Ticket Staff tạo được ghế (không có quyền)

4. **Lỗi tính toàn vẹn dữ liệu (2 lỗi):**
   - SEAT-11 + SEAT-16: Ghế với seat_type/status không hợp lệ sẽ ảnh hưởng đến tính giá vé và logic đặt vé

**Đánh giá mức độ nghiêm trọng:**
- **Critical**: SEAT-RBAC-01, SEAT-RBAC-02 (bypass phân quyền)
- **Major**: SEAT-11, SEAT-16 (dữ liệu không hợp lệ ảnh hưởng nghiệp vụ)
- **Minor**: SEAT-10, SEAT-17, SEAT-20 (xử lý lỗi không đúng chuẩn)

---

## 3. KẾT LUẬN TỔNG THỂ HỆ THỐNG

### 3.1. Bảng tổng hợp kết quả kiểm thử

| Chức năng | TESTED | PASS | FAIL | N/A | Tỷ lệ PASS |
|---|---|---|---|---|---|
| **Quản lý phim (Movie Management)** | 20 | 12 | 8 | 0 | 60% |
| **Quản lý ghế (Seat Management)** | 20 | 11 | 9 | 0 | 55% |
| **Tổng cộng** | **40** | **23** | **17** | **0** | **57.5%** |

### 3.2. Phân loại lỗi theo mức độ nghiêm trọng

| Mức độ | Số lượng | Mô tả |
|---|---|---|
| **Critical** | 3 | MOV-20, SEAT-RBAC-01, SEAT-RBAC-02 – Lỗ hổng phân quyền RBAC cho phép người dùng không có quyền thực hiện thao tác quản trị |
| **Major** | 5 | MOV-08, MOV-09, MOV-10, SEAT-11, SEAT-16 – Thiếu validation dữ liệu dẫn đến dữ liệu không hợp lệ được lưu vào hệ thống |
| **Minor** | 9 | MOV-04, MOV-13, MOV-15, SEAT-10, SEAT-17, SEAT-20, và các lỗi HTTP status code không chính xác |

### 3.3. Phân loại lỗi theo nguyên nhân gốc

| Nguyên nhân | Số lỗi | Ảnh hưởng |
|---|---|---|
| **Thiếu validation đầu vào** | 7 | Dữ liệu không hợp lệ được lưu vào DB (status, seat_type, whitespace-only values) |
| **Xử lý exception không đúng** | 5 | Trả HTTP 500 thay vì mã lỗi phù hợp (400/404), gây khó khăn cho client-side error handling |
| **Lỗ hổng RBAC / Phân quyền** | 3 | Bất kỳ user nào (kể cả customer) đều có thể tạo/sửa/xóa phim và ghế – vi phạm nghiêm trọng nguyên tắc least privilege |
| **Lỗi logic nghiệp vụ** | 2 | API Update phim không hoạt động khi không gửi trường status; genre không được gắn khi genre_id sai |

### 3.4. Đánh giá và phân tích tổng thể

**Hệ thống Cinema Management trong phiên bản kiểm thử hiện tại CHƯA ĐẠT yêu cầu chất lượng** với tỷ lệ PASS chỉ 57.5% (23/40 test case).

**Các vấn đề chính cần khắc phục theo thứ tự ưu tiên:**

1. **[Ưu tiên 1 - Critical] Sửa lỗi RBAC tại API Gateway**: Cấu hình `admin_paths` và `public_paths` trong API Gateway cần được rà soát lại toàn bộ. Hiện tại, endpoint POST/PUT/DELETE cho `/movies` và `/seats` không được bảo vệ đúng cách, cho phép bất kỳ user đã đăng nhập nào thực hiện thao tác quản trị.

2. **[Ưu tiên 2 - Major] Bổ sung validation dữ liệu**: Cần thêm validation `oneof` cho các trường enum (movie status, seat type, seat status), trim whitespace trước khi validate title và seat_number.

3. **[Ưu tiên 3 - Major] Sửa lỗi Update Movie**: Logic trong `UpdateMovie` business cần xử lý trường hợp client không gửi field status (giữ nguyên status hiện tại thay vì so sánh với "").

4. **[Ưu tiên 4 - Minor] Cải thiện Error Handling**: Map các lỗi nghiệp vụ (ErrInvalidStatusTransition, FK violation, Not Found) sang HTTP status code phù hợp (400, 404) thay vì để mặc định trả 500.

**Khuyến nghị**: Không nên deploy phiên bản này lên môi trường production cho đến khi ít nhất các lỗi Critical và Major được khắc phục. Cần thực hiện kiểm thử hồi quy (regression test) sau khi sửa lỗi.

