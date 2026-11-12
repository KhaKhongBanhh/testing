# BÁO CÁO KIỂM THỬ HỆ THỐNG (SYSTEM TEST)

## Thông tin chung

| Mục | Nội dung |
|---|---|
| **Tên dự án** | Xây dựng website quản lý rạp chiếu phim (Cinema Management System – HQ Cinema) |
| **Nội dung kiểm thử** | System Test – Kiểm thử hệ thống cho chức năng Quản lý phim (Movie Management) và Quản lý ghế (Seat Management) |
| **Phiên bản** | 1.0 |
| **Người kiểm thử** | Phạm Minh Đức |
| **Ngày thực hiện** | 11/11/2026 |
| **Môi trường** | Docker containers (12 services), API Gateway port 8000, PostgreSQL, Redis |
| **Trình duyệt / Công cụ** | Trang quản trị (Admin) HQ Cinema trên trình duyệt; kết quả thực tế được thu thập qua API (PowerShell + Invoke-WebRequest) |

> **Lưu ý về cách mô tả:** Các cột *Dữ liệu đầu vào*, *Các bước thực hiện*, *Kết quả mong đợi*, *Kết quả thực tế* được mô tả theo thao tác và màn hình mà người dùng nhìn thấy. Kết quả thực tế do tester đo được ở tầng API (qua API Gateway), nên được diễn đạt lại thành hiện tượng tương ứng trên giao diện. Riêng các ca mà form giao diện có kiểm tra phía client (MOV-03, MOV-04, MOV-07, MOV-11, SEAT-07, SEAT-08, SEAT-17), giao diện có thể chặn dữ liệu trước khi gửi lên server; lỗi ở các ca này là lỗi của **backend** (xem cột Ghi chú) và cần chạy lại trực tiếp trên giao diện để xác nhận.

---

## Phần đánh giá

### Thông tin tổng hợp

| Trạng thái | Số lượng | Tỷ lệ |
|---|---|---|
| **TESTED** | 43 | 100% |
| **PASS** | 29 | 67,4% |
| **FAIL** | 14 | 32,6% |
| **N/A** | 0 | 0% |

> Hệ thống Cinema Management có tỷ lệ lỗi 32,6% (14/43), tập trung vào: thiếu kiểm tra dữ liệu đầu vào (validation), xử lý lỗi sai (trả HTTP 500 thay vì 400/404) và lỗ hổng nghiêm trọng trong phân quyền (RBAC).

---

## 1. QUẢN LÝ PHIM (Movie Management)

### 1.1. Thống kê

| Trạng thái | Số lượng | Định nghĩa |
|---|---|---|
| **TESTED** | 20 | Tổng số test case đã thực hiện |
| **PASS** | 13 | Test case cho kết quả đúng như mong đợi |
| **FAIL** | 7 | Test case cho kết quả sai so với mong đợi |
| **N/A** | 0 | Test case không áp dụng được |

### 1.2. Test case

| ID | Tên chức năng kiểm thử | Kỹ thuật áp dụng | Dữ liệu đầu vào | Các bước thực hiện | Kết quả mong đợi | Kết quả thực tế | Trạng thái | Ghi chú |
|---|---|---|---|---|---|---|---|---|
| MOV-01 | Xem danh sách phim | Equivalence Class Partitioning | Trang danh sách phim, hiển thị trang 1, mỗi trang 10 phim. Không cần đăng nhập. | 1. Mở trang danh sách phim của website.<br>2. Quan sát danh sách phim hiển thị ở trang đầu tiên. | Trang hiển thị danh sách phim, mỗi phim có đầy đủ thông tin: tên phim, đạo diễn, diễn viên, thể loại, thời lượng, trạng thái. | Trang hiển thị 7 phim với đầy đủ thông tin. Tuy nhiên dữ liệu trả về còn chứa trường debug thừa `"hello":"world_ok!"` (người dùng không nhìn thấy trên giao diện). | **PASS** | Chức năng hoạt động đúng. Trường debug thừa không ảnh hưởng chức năng nhưng cần loại bỏ trước khi deploy. |
| MOV-02 | Tạo phim với dữ liệu hợp lệ | Equivalence Class Partitioning | Form "Thêm phim mới":<br>- Tên phim: QA Test Movie 01<br>- Thời lượng: 120<br>- Đạo diễn: QA Director<br>- Diễn viên: Actor A, Actor B<br>- Thể loại: tích chọn Action và Drama<br>- Ngày phát hành: 20/12/2026<br>- Trạng thái: Sắp chiếu | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào menu "Quản lý phim", bấm nút "Thêm phim mới".<br>3. Nhập đầy đủ dữ liệu như cột Dữ liệu đầu vào.<br>4. Bấm nút "Tạo phim mới".<br>5. Mở lại phim vừa tạo để kiểm tra thông tin. | Hệ thống quay về trang "Quản lý phim", phim "QA Test Movie 01" xuất hiện trong danh sách với trạng thái "UPCOMING"; mở chi tiết phim thấy đúng các thông tin đã nhập. | Tạo phim thành công, phim mới có mã id (d2748cd9...) và trạng thái UPCOMING. Mở chi tiết phim thấy đúng dữ liệu đã nhập. | **PASS** | |
| MOV-03 | Tạo phim với tên phim để trống | Equivalence Class Partitioning | Form "Thêm phim mới":<br>- Tên phim: (để trống)<br>- Thời lượng: 120 | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý phim", bấm "Thêm phim mới".<br>3. Để trống ô "Tên phim", nhập Thời lượng = 120.<br>4. Bấm "Tạo phim mới". | Hệ thống không tạo phim, hiển thị thông báo lỗi (ví dụ "Tên phim và thời lượng là bắt buộc") và vẫn ở lại form. | Hệ thống từ chối tạo phim (server trả mã lỗi 400), phim không xuất hiện trong danh sách. | **PASS** | Form giao diện cũng có kiểm tra bắt buộc phía client. |
| MOV-04 | Tạo phim với tên phim chỉ chứa khoảng trắng | Equivalence Class Partitioning | Form "Thêm phim mới":<br>- Tên phim: "   " (3 dấu cách)<br>- Thời lượng: 120 | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý phim", bấm "Thêm phim mới".<br>3. Ở ô "Tên phim" chỉ nhập 3 dấu cách, Thời lượng = 120.<br>4. Bấm "Tạo phim mới". | Hệ thống không tạo phim và hiển thị thông báo lỗi vì tên phim chỉ gồm khoảng trắng là không hợp lệ. | Server chấp nhận và tạo phim thành công với tên phim chỉ gồm khoảng trắng (HTTP 201); trong danh sách xuất hiện một phim không có tên hiển thị. | **FAIL** | **BUG (Minor)**: Backend (hàm `IsValid()`) chỉ kiểm tra `Title == ""` mà không cắt khoảng trắng (trim). Form giao diện có trim phía client nên cần chạy lại trên giao diện; lỗi backend vẫn tồn tại khi gọi API trực tiếp. |
| MOV-05 | Tạo phim với tên phim dài 255 ký tự (biên trên hợp lệ) | Boundary Value Analysis | Form "Thêm phim mới":<br>- Tên phim: chuỗi 255 ký tự "Q"<br>- Thời lượng: 100 | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý phim", bấm "Thêm phim mới".<br>3. Nhập Tên phim gồm đúng 255 ký tự "Q", Thời lượng = 100.<br>4. Bấm "Tạo phim mới". | Tạo phim thành công vì 255 là độ dài tối đa cho phép; phim xuất hiện trong danh sách. | Tạo phim thành công với tên phim 255 ký tự. | **PASS** | Biên trên hợp lệ. |
| MOV-06 | Tạo phim với tên phim dài 256 ký tự (vượt biên trên) | Boundary Value Analysis | Form "Thêm phim mới":<br>- Tên phim: chuỗi 256 ký tự "Q"<br>- Thời lượng: 100 | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý phim", bấm "Thêm phim mới".<br>3. Nhập Tên phim gồm 256 ký tự "Q", Thời lượng = 100.<br>4. Bấm "Tạo phim mới". | Hệ thống từ chối, hiển thị thông báo lỗi vì tên phim vượt quá 255 ký tự; phim không được tạo. | Hệ thống từ chối tạo phim (server trả mã lỗi 400), phim không xuất hiện trong danh sách. | **PASS** | Biên trên không hợp lệ. |
| MOV-07 | Tạo phim với thời lượng bằng 0 | Boundary Value Analysis | Form "Thêm phim mới":<br>- Tên phim: QA Zero Duration<br>- Thời lượng: 0 | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý phim", bấm "Thêm phim mới".<br>3. Nhập Tên phim = "QA Zero Duration", Thời lượng = 0.<br>4. Bấm "Tạo phim mới". | Hệ thống từ chối, báo lỗi vì thời lượng phải từ 1 phút trở lên; phim không được tạo. | Hệ thống từ chối thời lượng = 0 (server trả mã lỗi 400), phim không được tạo. | **PASS** | Ô "Thời lượng (phút)" trên form có thuộc tính `min=1` nên giao diện cũng có thể chặn trước. |
| MOV-08 | Tạo phim với trạng thái không hợp lệ | Equivalence Class Partitioning | Gửi dữ liệu phim:<br>- Tên phim: QA Invalid Status<br>- Thời lượng: 100<br>- Trạng thái: "ABCXYZ" (không thuộc Sắp chiếu / Đang chiếu / Đã kết thúc) | 1. Đăng nhập bằng tài khoản Manager.<br>2. Vào "Thêm phim mới", nhập Tên phim = "QA Invalid Status", Thời lượng = 100.<br>3. Đặt Trạng thái = "ABCXYZ" (không có trong dropdown nên gửi trực tiếp request tạo phim).<br>4. Xác nhận tạo phim. | Hệ thống từ chối và báo lỗi vì trạng thái không thuộc 3 giá trị cho phép; phim không được tạo. | Hệ thống vẫn tạo phim thành công (HTTP 201) với trạng thái "ABCXYZ". | **FAIL** | **BUG (Major)**: `CreateMovieRequest.Status` có tag `omitempty` nhưng không có validation `oneof`; chỉ gán mặc định UPCOMING khi trống, không kiểm tra giá trị hợp lệ khi client gửi lên. |
| MOV-09 | Tạo phim với thể loại không tồn tại | Equivalence Class Partitioning | Gửi dữ liệu phim:<br>- Tên phim: QA Invalid Genre<br>- Thời lượng: 100<br>- Thể loại: "non-existent-uuid" (mã thể loại không có trong hệ thống) | 1. Đăng nhập bằng tài khoản Manager.<br>2. Vào "Thêm phim mới", nhập Tên phim = "QA Invalid Genre", Thời lượng = 100.<br>3. Chọn thể loại có mã không tồn tại (gửi trực tiếp request vì checkbox chỉ liệt kê thể loại hợp lệ).<br>4. Xác nhận tạo phim. | Hệ thống từ chối với thông báo thể loại không tồn tại (lỗi 400); phim không được tạo. | Hệ thống báo lỗi máy chủ (HTTP 500) nhưng phim "QA Invalid Genre" (id=741ae107...) vẫn được tạo và không có thể loại nào. | **FAIL** | **BUG (Major)**: 1) Lỗi khóa ngoại (foreign key) không được xử lý, trả 500 thay vì 400. 2) Phim vẫn được lưu dù lỗi – thiếu transaction rollback. |
| MOV-10 | Cập nhật phim với dữ liệu hợp lệ | Equivalence Class Partitioning | Form "Chỉnh sửa phim" của phim "QA Test Movie 01":<br>- Tên phim: QA Test Movie 01 Updated<br>- Thời lượng: 150<br>- Mô tả: Updated description<br>- Không thay đổi trạng thái | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý phim", bấm vào phim "QA Test Movie 01" rồi bấm "Edit Movie".<br>3. Sửa Tên phim, Thời lượng và Mô tả như cột Dữ liệu đầu vào (không đổi Trạng thái).<br>4. Bấm "Cập nhật phim". | Cập nhật thành công, hệ thống quay về danh sách phim và thông tin mới (tên, thời lượng, mô tả) được hiển thị. | Hệ thống hiển thị thông báo lỗi "Không thể cập nhật phim" (server trả lỗi 500), dữ liệu phim không thay đổi. | **FAIL** | **BUG (Major)**: Khi request không gửi trường `status`, giá trị rỗng "" bị so sánh với trạng thái hiện tại → gọi `CanTransitionTo("")` → trả `ErrInvalidStatusTransition`; handler không bắt lỗi này nên trả 500. Lỗi kép: logic update sai + handler thiếu ánh xạ lỗi. Form giao diện luôn gửi kèm trạng thái nên cần chạy lại trên giao diện để xác nhận. |
| MOV-11 | Cập nhật phim với thời lượng âm | Boundary Value Analysis | Form "Chỉnh sửa phim":<br>- Tên phim: QA Test<br>- Thời lượng: -5 | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Mở form chỉnh sửa của một phim bất kỳ.<br>3. Nhập Tên phim = "QA Test", Thời lượng = -5.<br>4. Bấm "Cập nhật phim". | Hệ thống từ chối, hiển thị thông báo lỗi vì thời lượng phải từ 1 phút trở lên; dữ liệu phim không đổi. | Hệ thống từ chối thời lượng âm (server trả mã lỗi 400), dữ liệu phim không đổi. | **PASS** | |
| MOV-12 | Chuyển trạng thái phim UPCOMING → SHOWING | State Transition | Phim đang ở trạng thái "UPCOMING", chọn đổi sang "SHOWING". | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý phim", bấm vào một phim đang có trạng thái UPCOMING.<br>3. Ở mục "Thay đổi trạng thái" trên trang chi tiết, bấm nút "showing". | Nhãn trạng thái của phim đổi từ UPCOMING sang SHOWING; nút "showing" chuyển sang trạng thái đang chọn (bị vô hiệu hóa). | Trạng thái phim chuyển thành công từ UPCOMING sang SHOWING. | **PASS** | Chuyển trạng thái hợp lệ theo quy tắc UPCOMING → SHOWING. |
| MOV-13 | Chuyển trạng thái phim SHOWING → UPCOMING (không hợp lệ) | State Transition | Phim đang ở trạng thái "SHOWING", chọn đổi ngược về "UPCOMING". | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý phim", bấm vào một phim đang có trạng thái SHOWING.<br>3. Ở mục "Thay đổi trạng thái", bấm nút "upcoming". | Hệ thống từ chối và hiện thông báo lỗi rõ ràng (chuyển ngược trạng thái không hợp lệ); trạng thái phim vẫn là SHOWING. | Hiện hộp thoại "Failed to update movie status" kèm lỗi máy chủ (HTTP 500) thay vì thông báo lỗi nghiệp vụ (400). Trạng thái phim vẫn là SHOWING. | **FAIL** | **BUG (Minor)**: Handler `UpdateMovieStatus` không bắt `ErrInvalidStatusTransition` để trả 400; lỗi rơi vào nhánh mặc định trả 500. |
| MOV-14 | Chuyển trạng thái phim SHOWING → ENDED | State Transition | Phim đang ở trạng thái "SHOWING", chọn đổi sang "ENDED". | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý phim", bấm vào một phim đang có trạng thái SHOWING.<br>3. Ở mục "Thay đổi trạng thái", bấm nút "ended". | Nhãn trạng thái của phim đổi từ SHOWING sang ENDED. | Trạng thái phim chuyển thành công từ SHOWING sang ENDED. | **PASS** | Chuyển trạng thái hợp lệ theo quy tắc SHOWING → ENDED. |
| MOV-15 | Chuyển trạng thái phim ENDED → SHOWING (không hợp lệ) | State Transition | Phim đang ở trạng thái "ENDED", chọn đổi lại thành "SHOWING". | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý phim", bấm vào một phim đang có trạng thái ENDED.<br>3. Ở mục "Thay đổi trạng thái", bấm nút "showing". | Hệ thống từ chối và hiện thông báo lỗi rõ ràng vì ENDED là trạng thái cuối, không thể chuyển ngược; trạng thái vẫn là ENDED. | Hiện hộp thoại "Failed to update movie status" kèm lỗi máy chủ (HTTP 500) thay vì thông báo lỗi nghiệp vụ (400). Trạng thái phim vẫn là ENDED. | **FAIL** | **BUG (Minor)**: Cùng nguyên nhân MOV-13 – handler không bắt `ErrInvalidStatusTransition`. |
| MOV-16 | Xóa phim | Equivalence Class Partitioning | Phim "QA Delete Test Movie" vừa tạo mới để xóa. | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Tạo phim mới tên "QA Delete Test Movie".<br>3. Mở chi tiết phim đó, bấm nút "Delete Movie".<br>4. Ở hộp thoại xác nhận "Are you sure you want to delete this movie?", bấm OK. | Phim bị xóa, hệ thống quay về trang "Quản lý phim" và phim "QA Delete Test Movie" không còn trong danh sách. | Xóa phim thành công (server trả HTTP 204, không có nội dung), phim biến mất khỏi danh sách. | **PASS** | |
| MOV-17 | Xóa phim đã bị xóa | Equivalence Class Partitioning | Mã phim (id) của phim đã xóa ở MOV-16. | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Mở lại trang chi tiết bằng đường dẫn của phim đã xóa (hoặc gửi lại yêu cầu xóa với cùng id).<br>3. Thực hiện xóa phim. | Hệ thống thông báo phim không tồn tại (không tìm thấy phim). | Hệ thống trả lỗi "không tìm thấy phim" (HTTP 404). | **PASS** | Chỉ kiểm tra được bằng đường dẫn/request trực tiếp vì phim đã biến mất khỏi danh sách. |
| MOV-18 | Tìm kiếm phim theo từ khóa | Equivalence Class Partitioning | Ô "Tìm kiếm phim theo tên...": nhập "Spider". | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý phim".<br>3. Nhập "Spider" vào ô "Tìm kiếm phim theo tên..." rồi bấm nút "Tìm kiếm". | Hệ thống hiển thị dòng "Kết quả tìm kiếm cho "Spider" (1 phim được tìm thấy)" và danh sách chỉ gồm các phim có chứa từ "Spider" trong tên. | Hiển thị 1 phim: "Spider-Man: No Way Home" (trạng thái SHOWING). Tìm kiếm chính xác. | **PASS** | |
| MOV-19 | Lọc phim theo trạng thái "Đang chiếu" | Equivalence Class Partitioning | Bộ lọc trạng thái: chọn "SHOWING" (Đang chiếu). | 1. Mở trang danh sách phim.<br>2. Chọn bộ lọc trạng thái "Đang chiếu" (SHOWING). | Danh sách chỉ hiển thị các phim có trạng thái SHOWING. | Hiển thị 5 phim: Avengers: Endgame, Spider-Man, Dune, The Batman, Inception – tất cả đều có trạng thái SHOWING. | **PASS** | Trang "Quản lý phim" của Admin hiện không có bộ lọc trạng thái; bộ lọc này dùng cho danh sách phim phía khách hàng / tham số `status` của API. |
| MOV-20 | Kiểm tra phân quyền: Customer tạo phim | Decision Table | Tài khoản Customer (alice@email.com, vai trò customer, không có quyền `movie_manage`); dữ liệu phim hợp lệ: Tên phim "Customer Should Not Create". | 1. Đăng nhập bằng tài khoản Customer (alice@email.com).<br>2. Gửi yêu cầu tạo phim "Customer Should Not Create" (giống thao tác "Thêm phim mới" của Manager) kèm token của Customer. | Hệ thống từ chối với thông báo không đủ quyền (403 Forbidden); phim không được tạo. | Hệ thống cho phép tạo phim thành công (HTTP 201), phim "Customer Should Not Create" (id=a417b8c7...) xuất hiện trong danh sách phim. | **FAIL** | **BUG NGHIÊM TRỌNG (Critical)**: API Gateway không kiểm tra quyền cho `POST /movies`. Customer không có quyền `movie_manage` vẫn tạo được phim – lỗ hổng RBAC. |

### 1.3. Kết luận chức năng Quản lý phim

**Tổng kết:**

| Chỉ số | Giá trị |
|---|---|
| Tổng test case | 20 |
| PASS | 13 (65%) |
| FAIL | 7 (35%) |
| N/A | 0 |

**Phân tích kết quả:** Chức năng Quản lý phim có **7 lỗi** trên 20 test case (tỷ lệ FAIL 35%), phân loại theo nguyên nhân chính (mỗi lỗi chỉ tính một nhóm):

1. **Thiếu validation dữ liệu đầu vào (2 lỗi):**
   - MOV-04: Chấp nhận tên phim chỉ chứa khoảng trắng
   - MOV-08: Chấp nhận giá trị trạng thái không hợp lệ (ví dụ "ABCXYZ")

2. **Xử lý lỗi sai – trả HTTP 500 thay vì 400 (3 lỗi):**
   - MOV-09: Thể loại không tồn tại trả 500 và vẫn lưu phim (thiếu rollback)
   - MOV-13, MOV-15: Chuyển trạng thái không hợp lệ trả 500 thay vì 400

3. **Lỗi logic nghiệp vụ (1 lỗi):**
   - MOV-10: Cập nhật phim hợp lệ bị lỗi do so sánh trạng thái rỗng trong business layer (đồng thời thiếu ánh xạ lỗi ở handler)

4. **Lỗi bảo mật RBAC (1 lỗi – Nghiêm trọng):**
   - MOV-20: Customer (không có quyền) vẫn tạo được phim

**Mức độ nghiêm trọng:**
- **Critical**: MOV-20
- **Major**: MOV-08, MOV-09, MOV-10
- **Minor**: MOV-04, MOV-13, MOV-15

---

## 2. QUẢN LÝ GHẾ (Seat Management)

### 2.1. Thống kê

> Thống kê tính cả 3 test case bổ sung ở mục 2.3 (SEAT-RBAC-01, SEAT-RBAC-02, SEAT-PW-01) vì đều là test case đã thực hiện.

| Trạng thái | Số lượng | Định nghĩa |
|---|---|---|
| **TESTED** | 23 | Tổng số test case đã thực hiện (20 chính + 3 bổ sung) |
| **PASS** | 16 | Test case cho kết quả đúng như mong đợi |
| **FAIL** | 7 | Test case cho kết quả sai so với mong đợi |
| **N/A** | 0 | Test case không áp dụng được |

### 2.2. Test case

| ID | Tên chức năng kiểm thử | Kỹ thuật áp dụng | Dữ liệu đầu vào | Các bước thực hiện | Kết quả mong đợi | Kết quả thực tế | Trạng thái | Ghi chú |
|---|---|---|---|---|---|---|---|---|
| SEAT-01 | Xem danh sách ghế theo phòng | Equivalence Class Partitioning | Màn hình "Quản lý Ghế ngồi", chế độ "Bảng":<br>- Phòng chiếu: Phòng 1<br>- Trang 1, 10 ghế/trang | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào menu "Quản lý Ghế ngồi", bấm nút "Bảng".<br>3. Nhập mã (Room ID) của Phòng 1 vào ô "Room ID". | Bảng hiển thị danh sách ghế của Phòng 1 với các cột: Phòng, Vị trí, Loại ghế, Trạng thái, Ngày tạo, Thao tác; có phân trang. | Bảng hiển thị danh sách ghế Phòng 1 đầy đủ thông tin (vị trí, loại Thường/VIP, trạng thái "Có sẵn"). | **PASS** | |
| SEAT-02 | Lọc ghế theo loại ghế VIP | Equivalence Class Partitioning | Chế độ "Bảng":<br>- Phòng chiếu: Phòng 1<br>- Dropdown "Tất cả loại ghế": chọn "VIP" | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý Ghế ngồi", bấm "Bảng", nhập Room ID của Phòng 1.<br>3. Ở dropdown "Tất cả loại ghế", chọn "VIP". | Bảng chỉ hiển thị các ghế có Loại ghế = "VIP". | Bảng hiển thị danh sách ghế VIP, tất cả các dòng đều có loại ghế "VIP". | **PASS** | |
| SEAT-03 | Lọc ghế theo trạng thái "Có sẵn" | Equivalence Class Partitioning | Chế độ "Bảng":<br>- Phòng chiếu: Phòng 1<br>- Dropdown "Tất cả trạng thái": chọn "Có sẵn" (AVAILABLE)<br>- 5 ghế/trang | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý Ghế ngồi", bấm "Bảng", nhập Room ID của Phòng 1.<br>3. Ở dropdown "Tất cả trạng thái", chọn "Có sẵn". | Bảng chỉ hiển thị các ghế có Trạng thái = "Có sẵn". | Bảng hiển thị danh sách ghế "Có sẵn", tất cả các dòng đều có trạng thái AVAILABLE. | **PASS** | |
| SEAT-04 | Xem chi tiết ghế theo ID | Equivalence Class Partitioning | Ghế A-01 của Phòng 1 (id = 136c7bb3...). | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý Ghế ngồi", bấm "Bảng", nhập Room ID của Phòng 1.<br>3. Ở dòng ghế A01, bấm biểu tượng "Chỉnh sửa" để mở form chi tiết ghế. | Form "Chỉnh sửa ghế ngồi" hiển thị đúng thông tin ghế: Phòng chiếu, Hàng ghế, Số ghế, Loại ghế, Trạng thái. | Form hiển thị đúng: Phòng 1, Hàng ghế A, Số ghế 01, Loại ghế Thường, Trạng thái "Có sẵn". | **PASS** | |
| SEAT-05 | Xem ghế với ID không tồn tại | Equivalence Class Partitioning | Mã ghế không tồn tại: 00000000-0000-0000-0000-000000000000. | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Truy cập đường dẫn chỉnh sửa ghế với id không tồn tại: `/admin/seats/00000000-0000-0000-0000-000000000000/edit`. | Hệ thống thông báo không tìm thấy ghế. | Hệ thống trả lỗi "ghế không tồn tại" (HTTP 404); form hiển thị thông báo lỗi tải dữ liệu. | **PASS** | |
| SEAT-06 | Tạo ghế với dữ liệu hợp lệ | Equivalence Class Partitioning | Form "Thêm ghế ngồi mới":<br>- Phòng chiếu: Phòng 1<br>- Hàng ghế: Z<br>- Số ghế: 99<br>- Loại ghế: Thường | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý Ghế ngồi", bấm "Bảng", bấm nút "Thêm ghế mới".<br>3. Chọn Phòng chiếu = "Phòng 1", nhập Hàng ghế = "Z", Số ghế = "99", chọn Loại ghế = "Thường".<br>4. Bấm "Tạo mới". | Hệ thống quay về danh sách ghế; ghế Z99 của Phòng 1 xuất hiện với loại "Thường" và trạng thái mặc định "Có sẵn". | Tạo ghế thành công (id=8f0ae512...): hàng Z, số 99, loại Thường, trạng thái "Có sẵn". | **PASS** | |
| SEAT-07 | Tạo ghế với số ghế để trống | Equivalence Class Partitioning | Form "Thêm ghế ngồi mới":<br>- Phòng chiếu: Phòng 1<br>- Hàng ghế: Z<br>- Số ghế: (để trống)<br>- Loại ghế: Thường | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý Ghế ngồi", bấm "Bảng", bấm "Thêm ghế mới".<br>3. Chọn Phòng 1, nhập Hàng ghế = "Z", để trống "Số ghế".<br>4. Bấm "Tạo mới". | Hệ thống không tạo ghế, hiển thị thông báo lỗi (ví dụ "Vui lòng điền đầy đủ thông tin bắt buộc") vì Số ghế là trường bắt buộc (*). | Hệ thống từ chối tạo ghế (server trả mã lỗi 400), ghế không được tạo. | **PASS** | Form giao diện cũng có kiểm tra bắt buộc phía client. |
| SEAT-08 | Tạo ghế khi chưa chọn phòng chiếu | Equivalence Class Partitioning | Form "Thêm ghế ngồi mới":<br>- Phòng chiếu: để mặc định "Chọn phòng"<br>- Hàng ghế: A<br>- Số ghế: 01<br>- Loại ghế: Thường | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý Ghế ngồi", bấm "Bảng", bấm "Thêm ghế mới".<br>3. Không chọn Phòng chiếu, nhập Hàng ghế = "A", Số ghế = "01".<br>4. Bấm "Tạo mới". | Hệ thống không tạo ghế và báo lỗi vì Phòng chiếu là trường bắt buộc (*). | Hệ thống từ chối tạo ghế (server trả mã lỗi 400), ghế không được tạo. | **PASS** | Form giao diện cũng có kiểm tra bắt buộc phía client. |
| SEAT-09 | Tạo ghế trùng vị trí | Decision Table | Form "Thêm ghế ngồi mới":<br>- Phòng chiếu: Phòng 1<br>- Hàng ghế: A<br>- Số ghế: 01 (vị trí A01 đã tồn tại)<br>- Loại ghế: Thường | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý Ghế ngồi", bấm "Bảng", bấm "Thêm ghế mới".<br>3. Chọn Phòng 1, nhập Hàng ghế = "A", Số ghế = "01".<br>4. Bấm "Tạo mới". | Hệ thống không tạo ghế và hiển thị thông báo "Vị trí ghế này đã tồn tại trong phòng". | Hệ thống từ chối tạo ghế trùng vị trí (server trả mã lỗi 400). | **PASS** | |
| SEAT-10 | Tạo ghế với phòng chiếu không tồn tại | Equivalence Class Partitioning | Gửi dữ liệu ghế:<br>- Mã phòng: 00000000-0000-0000-0000-000000000000 (không tồn tại)<br>- Hàng ghế: A<br>- Số ghế: 01<br>- Loại ghế: Thường | 1. Đăng nhập bằng tài khoản Manager.<br>2. Gửi yêu cầu tạo ghế với mã phòng không tồn tại (dropdown "Phòng chiếu" chỉ liệt kê phòng hợp lệ nên phải gửi trực tiếp request).<br>3. Xác nhận tạo ghế. | Hệ thống từ chối với thông báo phòng chiếu không tồn tại (lỗi 400). | Hệ thống báo lỗi máy chủ (HTTP 500); giao diện hiện thông báo chung "Có lỗi xảy ra khi tạo ghế". | **FAIL** | **BUG (Minor)**: Lỗi khóa ngoại (room_id không tồn tại) không được xử lý, trả 500 thay vì 400 với thông báo rõ ràng. |
| SEAT-11 | Tạo ghế với loại ghế không hợp lệ | Equivalence Class Partitioning | Gửi dữ liệu ghế:<br>- Phòng chiếu: Phòng 1<br>- Hàng ghế: Z<br>- Số ghế: 98<br>- Loại ghế: "INVALID_TYPE" (không thuộc Thường / VIP / Đôi) | 1. Đăng nhập bằng tài khoản Manager.<br>2. Gửi yêu cầu tạo ghế Z98 ở Phòng 1 với loại ghế "INVALID_TYPE" (dropdown "Loại ghế" chỉ có 3 giá trị hợp lệ nên phải gửi trực tiếp request).<br>3. Xác nhận tạo ghế. | Hệ thống từ chối và báo lỗi vì loại ghế phải là Thường, VIP hoặc Đôi; ghế không được tạo. | Hệ thống vẫn tạo ghế thành công (HTTP 201) với loại ghế "INVALID_TYPE". | **FAIL** | **BUG (Major)**: Không validate giá trị `seat_type`; cho phép lưu loại ghế tùy ý, gây dữ liệu không nhất quán và ảnh hưởng tính giá vé (hệ số giá theo loại ghế). |
| SEAT-12 | Cập nhật ghế với dữ liệu hợp lệ | Equivalence Class Partitioning | Form "Chỉnh sửa ghế ngồi" của ghế Z99, Phòng 1: đổi Loại ghế từ "Thường" sang "VIP". | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý Ghế ngồi", bấm "Bảng", nhập Room ID của Phòng 1, bấm biểu tượng "Chỉnh sửa" ở dòng ghế Z99.<br>3. Đổi Loại ghế = "VIP".<br>4. Bấm "Cập nhật". | Cập nhật thành công, quay về danh sách ghế; ghế Z99 hiển thị loại "VIP", trạng thái vẫn "Có sẵn". | Cập nhật thành công: loại ghế đổi từ Thường sang VIP, trạng thái vẫn là "Có sẵn". | **PASS** | |
| SEAT-13 | Cập nhật ghế sang vị trí đã có ghế khác | Decision Table | Form "Chỉnh sửa ghế ngồi" của ghế Z99, Phòng 1: đổi Hàng ghế = "A", Số ghế = "01" (A01 đã có ghế khác). | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý Ghế ngồi", bấm "Bảng", nhập Room ID của Phòng 1, bấm "Chỉnh sửa" ở dòng ghế Z99.<br>3. Đổi Hàng ghế = "A", Số ghế = "01".<br>4. Bấm "Cập nhật". | Hệ thống không cập nhật và hiển thị thông báo "Vị trí ghế này đã tồn tại trong phòng". | Hệ thống từ chối cập nhật do trùng vị trí (server trả mã lỗi 400). | **PASS** | |
| SEAT-14 | Chuyển trạng thái ghế "Có sẵn" → "Bảo trì" | State Transition | Ghế đang "Có sẵn", chọn đổi sang "Bảo trì" (MAINTENANCE). | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý Ghế ngồi", bấm "Bảng", nhập Room ID của Phòng 1.<br>3. Ở dòng ghế cần đổi, mở dropdown ở cột "Trạng thái" và chọn "Bảo trì". | Trạng thái của ghế đổi thành "Bảo trì" (nhãn màu vàng) và được giữ nguyên sau khi tải lại bảng. | Trạng thái ghế chuyển thành công từ "Có sẵn" sang "Bảo trì". | **PASS** | |
| SEAT-15 | Chuyển trạng thái ghế "Bảo trì" → "Có sẵn" | State Transition | Ghế đang "Bảo trì", chọn đổi lại "Có sẵn" (AVAILABLE). | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý Ghế ngồi", bấm "Bảng", nhập Room ID của Phòng 1.<br>3. Ở dòng ghế đang "Bảo trì", mở dropdown cột "Trạng thái" và chọn "Có sẵn". | Trạng thái ghế đổi lại thành "Có sẵn" (nhãn màu xanh). | Trạng thái ghế chuyển thành công từ "Bảo trì" sang "Có sẵn". | **PASS** | |
| SEAT-16 | Cập nhật trạng thái ghế với giá trị không hợp lệ | Equivalence Class Partitioning | Gửi trạng thái ghế: "INVALID_STATUS" (không thuộc Có sẵn / Đã đặt / Bảo trì / Bị khóa). | 1. Đăng nhập bằng tài khoản Manager.<br>2. Gửi yêu cầu đổi trạng thái của một ghế sang "INVALID_STATUS" (dropdown "Trạng thái" chỉ có 4 giá trị hợp lệ nên phải gửi trực tiếp request). | Hệ thống từ chối và báo lỗi vì trạng thái không hợp lệ; trạng thái ghế giữ nguyên. | Hệ thống chấp nhận và lưu trạng thái "INVALID_STATUS" (HTTP 200); ghế có trạng thái không thuộc danh sách chuẩn. | **FAIL** | **BUG (Major)**: API `UpdateSeatStatus` không validate trạng thái; chuỗi bất kỳ đều được lưu, làm hỏng dữ liệu và gây lỗi logic ở luồng đặt vé. |
| SEAT-17 | Tạo ghế với số ghế chỉ chứa khoảng trắng | Equivalence Class Partitioning | Form "Thêm ghế ngồi mới":<br>- Phòng chiếu: Phòng 1<br>- Hàng ghế: Y<br>- Số ghế: "   " (3 dấu cách)<br>- Loại ghế: Thường | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý Ghế ngồi", bấm "Bảng", bấm "Thêm ghế mới".<br>3. Chọn Phòng 1, nhập Hàng ghế = "Y", ở ô "Số ghế" chỉ nhập 3 dấu cách.<br>4. Bấm "Tạo mới". | Hệ thống không tạo ghế và báo lỗi vì số ghế chỉ gồm khoảng trắng là không hợp lệ. | Hệ thống chấp nhận và tạo ghế thành công (HTTP 201) với số ghế chỉ gồm khoảng trắng. | **FAIL** | **BUG (Minor)**: Tương tự MOV-04, backend không trim khoảng trắng khi validate. Form giao diện chỉ kiểm tra rỗng (không trim) nên cũng không chặn được. |
| SEAT-18 | Tạo ghế loại Đôi (COUPLE) | Equivalence Class Partitioning | Form "Thêm ghế ngồi mới":<br>- Phòng chiếu: Phòng 1<br>- Hàng ghế: Z<br>- Số ghế: 97<br>- Loại ghế: Đôi | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý Ghế ngồi", bấm "Bảng", bấm "Thêm ghế mới".<br>3. Chọn Phòng 1, nhập Hàng ghế = "Z", Số ghế = "97", chọn Loại ghế = "Đôi".<br>4. Bấm "Tạo mới". | Tạo ghế thành công; ghế Z97 xuất hiện trong danh sách với loại "Đôi". | Tạo ghế thành công (id=4e31de5c...): loại Đôi (COUPLE), trạng thái "Có sẵn". | **PASS** | Xác nhận hệ thống hỗ trợ ghế loại Đôi. Hệ số giá của loại ghế này cần đối chiếu riêng theo đặc tả. |
| SEAT-19 | Xóa ghế | Equivalence Class Partitioning | Ghế có số ghế chỉ gồm khoảng trắng (hàng Y) đã tạo ở SEAT-17. | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý Ghế ngồi", bấm "Bảng", nhập Room ID của Phòng 1, nhập hàng "Y" vào ô "Hàng ghế".<br>3. Bấm biểu tượng thùng rác ở cột "Thao tác" của ghế cần xóa.<br>4. Ở hộp thoại "Bạn có chắc chắn muốn xóa ghế này?", bấm OK. | Ghế bị xóa và biến mất khỏi bảng danh sách. | Xóa ghế thành công (server trả HTTP 204), ghế biến mất khỏi danh sách. | **PASS** | |
| SEAT-20 | Xóa ghế đã bị xóa | Equivalence Class Partitioning | Mã ghế (id) của ghế đã xóa ở SEAT-19. | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Gửi lại yêu cầu xóa ghế với id đã xóa ở SEAT-19 (ghế đã biến mất khỏi bảng nên phải gửi trực tiếp request). | Hệ thống thông báo ghế không tồn tại (lỗi 404). | Hệ thống báo lỗi máy chủ (HTTP 500); giao diện hiện hộp thoại chung "Có lỗi xảy ra khi xóa ghế". | **FAIL** | **BUG (Minor)**: Xóa ghế không tồn tại trả 500 thay vì 404; handler không kiểm tra tồn tại hoặc không ánh xạ lỗi Not Found. |

### 2.3. Test case bổ sung (Kiểm tra phân quyền RBAC và lọc đa điều kiện)

> Các test case dưới đây được thực hiện cùng đợt với 20 test case chính và **được tính vào thống kê** ở mục 2.1 và mục 3.

| ID | Tên chức năng kiểm thử | Kỹ thuật áp dụng | Dữ liệu đầu vào | Các bước thực hiện | Kết quả mong đợi | Kết quả thực tế | Trạng thái | Ghi chú |
|---|---|---|---|---|---|---|---|---|
| SEAT-RBAC-01 | Customer tạo ghế | Decision Table | Tài khoản Customer (không có quyền `seat_manage`); dữ liệu ghế hợp lệ (Phòng 1, hàng Z, số ghế bất kỳ chưa tồn tại, loại Thường). | 1. Đăng nhập bằng tài khoản Customer.<br>2. Gửi yêu cầu tạo ghế mới (giống thao tác "Thêm ghế mới" của Manager) kèm token của Customer. | Hệ thống từ chối với thông báo không đủ quyền (403 Forbidden); ghế không được tạo. | Hệ thống cho phép tạo ghế thành công (HTTP 201). | **FAIL** | **BUG NGHIÊM TRỌNG (Critical)**: Customer không có quyền `seat_manage` nhưng vẫn tạo được ghế. |
| SEAT-RBAC-02 | Ticket Staff tạo ghế | Decision Table | Tài khoản Ticket Staff (không có quyền `seat_manage`); dữ liệu ghế hợp lệ (Phòng 1, hàng Z, số ghế bất kỳ chưa tồn tại, loại Thường). | 1. Đăng nhập bằng tài khoản Ticket Staff.<br>2. Gửi yêu cầu tạo ghế mới kèm token của Ticket Staff. | Hệ thống từ chối với thông báo không đủ quyền (403 Forbidden); ghế không được tạo. | Hệ thống cho phép tạo ghế thành công (HTTP 201). | **FAIL** | **BUG NGHIÊM TRỌNG (Critical)**: Ticket Staff không có quyền `seat_manage` nhưng vẫn tạo được ghế. |
| SEAT-PW-01 | Lọc ghế đa điều kiện (phòng + loại ghế + hàng) | Pairwise Testing | Chế độ "Bảng":<br>- Room ID: Phòng 1<br>- Loại ghế: VIP<br>- Hàng ghế: F | 1. Đăng nhập trang quản trị bằng tài khoản Manager.<br>2. Vào "Quản lý Ghế ngồi", bấm "Bảng".<br>3. Nhập Room ID của Phòng 1, nhập "F" vào ô "Hàng ghế", chọn "VIP" ở dropdown "Tất cả loại ghế". | Bảng chỉ hiển thị các ghế VIP ở hàng F của Phòng 1. | Bảng hiển thị 8 ghế VIP ở hàng F (F01 đến F08). Lọc đa điều kiện chính xác. | **PASS** | |

### 2.4. Kết luận chức năng Quản lý ghế

**Tổng kết (23 test case, gồm 3 test case bổ sung):**

| Chỉ số | Giá trị |
|---|---|
| Tổng test case | 23 |
| PASS | 16 (69,6%) |
| FAIL | 7 (30,4%) |
| N/A | 0 |

**Phân tích kết quả:** Chức năng Quản lý ghế có **7 lỗi** trên 23 test case (tỷ lệ FAIL 30,4%), phân loại theo nguyên nhân chính (mỗi lỗi chỉ tính một nhóm):

1. **Thiếu validation dữ liệu đầu vào (3 lỗi):**
   - SEAT-11: Chấp nhận loại ghế không hợp lệ ("INVALID_TYPE")
   - SEAT-16: Chấp nhận trạng thái ghế không hợp lệ ("INVALID_STATUS")
   - SEAT-17: Chấp nhận số ghế chỉ chứa khoảng trắng

2. **Xử lý lỗi sai – trả HTTP 500 thay vì 400/404 (2 lỗi):**
   - SEAT-10: Phòng chiếu không tồn tại trả 500 thay vì 400
   - SEAT-20: Xóa ghế đã xóa trả 500 thay vì 404

3. **Lỗi bảo mật RBAC (2 lỗi – Nghiêm trọng):**
   - SEAT-RBAC-01: Customer tạo được ghế
   - SEAT-RBAC-02: Ticket Staff tạo được ghế

**Mức độ nghiêm trọng:**
- **Critical**: SEAT-RBAC-01, SEAT-RBAC-02
- **Major**: SEAT-11, SEAT-16 (dữ liệu không hợp lệ ảnh hưởng nghiệp vụ giá vé và đặt vé)
- **Minor**: SEAT-10, SEAT-17, SEAT-20

---

## 3. KẾT LUẬN TỔNG THỂ HỆ THỐNG

### 3.1. Bảng tổng hợp kết quả kiểm thử

| Chức năng | TESTED | PASS | FAIL | N/A | Tỷ lệ PASS |
|---|---|---|---|---|---|
| **Quản lý phim (Movie Management)** | 20 | 13 | 7 | 0 | 65% |
| **Quản lý ghế (Seat Management)** | 23 | 16 | 7 | 0 | 69,6% |
| **Tổng cộng** | **43** | **29** | **14** | **0** | **67,4%** |

### 3.2. Phân loại lỗi theo mức độ nghiêm trọng

| Mức độ | Số lượng | Mô tả |
|---|---|---|
| **Critical** | 3 | MOV-20, SEAT-RBAC-01, SEAT-RBAC-02 – Lỗ hổng phân quyền RBAC cho phép người dùng không có quyền thực hiện thao tác quản trị |
| **Major** | 5 | MOV-08, MOV-09, MOV-10, SEAT-11, SEAT-16 – Thiếu validation / xử lý sai dẫn đến dữ liệu không hợp lệ hoặc không thể cập nhật |
| **Minor** | 6 | MOV-04, MOV-13, MOV-15, SEAT-10, SEAT-17, SEAT-20 – Validation khoảng trắng và mã HTTP lỗi không chính xác |

### 3.3. Phân loại lỗi theo nguyên nhân gốc

| Nguyên nhân | Số lỗi | Lỗi liên quan | Ảnh hưởng |
|---|---|---|---|
| **Thiếu validation đầu vào** | 5 | MOV-04, MOV-08, SEAT-11, SEAT-16, SEAT-17 | Dữ liệu không hợp lệ được lưu vào DB (trạng thái, loại ghế, giá trị chỉ có khoảng trắng) |
| **Xử lý exception không đúng** | 5 | MOV-09, MOV-13, MOV-15, SEAT-10, SEAT-20 | Trả HTTP 500 thay vì 400/404, khó xử lý lỗi ở client; MOV-09 còn lưu phim dở dang |
| **Lỗ hổng RBAC / Phân quyền** | 3 | MOV-20, SEAT-RBAC-01, SEAT-RBAC-02 | Customer/Ticket Staff tạo được phim và ghế – vi phạm nguyên tắc least privilege |
| **Lỗi logic nghiệp vụ** | 1 | MOV-10 | API cập nhật phim không hoạt động khi không gửi trường status |
| **Tổng** | **14** | | |

> Chỉ thao tác **tạo** (POST) được kiểm thử phân quyền. Chưa có bằng chứng cho thao tác sửa/xóa; cần bổ sung test case RBAC cho PUT/PATCH/DELETE trước khi kết luận lỗ hổng bao trùm các thao tác này.

### 3.4. Đánh giá và phân tích tổng thể

**Hệ thống Cinema Management trong phiên bản kiểm thử hiện tại CHƯA ĐẠT yêu cầu chất lượng** với tỷ lệ PASS 67,4% (29/43 test case) và còn 3 lỗi Critical.

**Các vấn đề cần khắc phục theo thứ tự ưu tiên:**

1. **[Ưu tiên 1 – Critical] Sửa lỗi RBAC tại API Gateway**: Rà soát cấu hình `admin_paths` và `public_paths`; endpoint tạo phim (`POST /movies`) và tạo ghế (`POST /seats`) hiện không được bảo vệ, cho phép Customer và Ticket Staff thực hiện thao tác quản trị.

2. **[Ưu tiên 2 – Major] Bổ sung validation dữ liệu**: Thêm validation `oneof` cho các trường enum (trạng thái phim, loại ghế, trạng thái ghế); trim khoảng trắng trước khi kiểm tra tên phim và số ghế.

3. **[Ưu tiên 3 – Major] Sửa lỗi Update Movie**: Logic `UpdateMovie` cần giữ nguyên trạng thái hiện tại khi client không gửi trường status (thay vì so sánh với chuỗi rỗng); đồng thời thêm transaction rollback khi gắn thể loại lỗi (MOV-09).

4. **[Ưu tiên 4 – Minor] Cải thiện Error Handling**: Ánh xạ các lỗi nghiệp vụ (`ErrInvalidStatusTransition`, vi phạm khóa ngoại, Not Found) sang mã HTTP phù hợp (400, 404) thay vì mặc định 500.

**Khuyến nghị**: Không nên deploy phiên bản này lên môi trường production cho đến khi các lỗi Critical và Major được khắc phục. Cần kiểm thử hồi quy (regression test) sau khi sửa lỗi, và dọn các dữ liệu rác do kiểm thử tạo ra (phim "QA ...", phim do Customer tạo, ghế Z97–Z99 và ghế tạo trong SEAT-11/RBAC).
