# BÁO CÁO KIỂM THỬ HỆ THỐNG – PHẦN MỞ RỘNG (SYSTEM TEST – EXTENDED)

## Thông tin chung

| Mục | Nội dung |
|---|---|
| **Tên dự án** | Xây dựng website quản lý rạp chiếu phim (Cinema Management System – HQ Cinema) |
| **Nội dung kiểm thử** | System Test mở rộng – Bổ sung 40 ca kiểm thử hệ thống (20 Movie + 20 Seat) nhằm phát hiện thêm lỗi tiềm ẩn |
| **Phiên bản** | 1.1 |
| **Người kiểm thử** | Phạm Minh Đức |
| **Ngày thực hiện** | 24/11/2026 |
| **Môi trường** | Máy trạm kiểm thử, Docker containers hệ thống, Database PostgreSQL, Redis Cache |
| **Trình duyệt / Thiết bị** | Google Chrome (phiên bản mới nhất) trên hệ điều hành Windows – Màn hình độ phân giải Full HD (1920x1080) |

> **Phạm vi kiểm thử:** Toàn bộ quá trình kiểm thử hệ thống được thực hiện trực tiếp trên giao diện đồ họa người dùng (GUI) của Trang quản trị (Admin Portal) HQ Cinema thông qua trình duyệt web. Không sử dụng các công cụ gọi API trực tiếp (như Postman, cURL, script API). Mọi thao tác, dữ liệu nhập và kết quả ghi nhận đều bám sát trải nghiệm người dùng thực tế trên các màn hình chức năng.

---

## Phần đánh giá

### Thông tin tổng hợp kết quả kiểm thử mở rộng

| Trạng thái | Số lượng | Tỷ lệ |
|---|---|---|
| **TỔNG TEST CASE (TESTED)** | 40 | 100% |
| **ĐẠT (PASS)** | 30 | 75,0% |
| **LỖI (FAIL)** | 10 | 25,0% |
| **KHÔNG ÁP DỤNG (N/A)** | 0 | 0% |

> **Nhận xét chung:** Đợt kiểm thử mở rộng tập trung vào các kịch bản biên, kịch bản bất thường và các tình huống nghiệp vụ nâng cao mà đợt kiểm thử trước chưa bao phủ. Kết quả phát hiện thêm **10 lỗi mới**, trong đó có **3 lỗi mức Major** (ảnh hưởng trực tiếp đến trải nghiệm chức năng) và **7 lỗi mức Minor** (giao diện, UX). Đặc biệt phát hiện nhiều lỗi không nhất quán chữ hoa/chữ thường (case-sensitivity) trong mã nguồn giao diện trang chi tiết phim.

---

## 1. QUẢN LÝ PHIM – MỞ RỘNG (Movie Management – Extended)

### 1.1. Thống kê kết quả

| Trạng thái | Số lượng | Tỷ lệ | Định nghĩa |
|---|---|---|---|
| **TESTED** | 20 | 100% | Tổng số test case hệ thống đã thực hiện |
| **PASS** | 14 | 70,0% | Test case cho kết quả thực tế đúng như kỳ vọng |
| **FAIL** | 6 | 30,0% | Test case phát hiện sai lệch hoặc lỗi giao diện / nghiệp vụ |
| **N/A** | 0 | 0% | Test case không áp dụng được |

### 1.2. Bảng chi tiết Test Case Quản lý phim mở rộng (20 ca)

| ID | Tên chức năng kiểm thử | Loại test case | Dữ liệu đầu vào | Các bước thực hiện | Kết quả mong đợi | Kết quả thực tế | Trạng thái | Ghi chú |
|---|---|---|---|---|---|---|---|---|
| MOV-31 | Kiểm tra màu viền (border) nút trạng thái trên trang chi tiết phim | A. GUI Testing | Phim đang có trạng thái "SHOWING" (Đang chiếu) | 1. Mở trang chi tiết bộ phim có trạng thái SHOWING.<br>2. Quan sát mục "Thay đổi trạng thái", kiểm tra màu viền (border) của 3 nút "upcoming", "showing", "ended". | Nút "showing" phải có viền màu xanh lá (`border-green-500`), nút "upcoming" viền cam (`border-orange-500`), nút "ended" viền đỏ (`border-red-500`), tương ứng với ý nghĩa trạng thái. | Cả 3 nút đều hiển thị viền màu xám mặc định (`border-gray-500`), không có sự phân biệt màu sắc nào giữa các nút trạng thái. | **FAIL** | **BUG (Major)**: Hàm `getStatusBorderClasses()` trong `MovieDetailPage.jsx` (dòng 68-78) sử dụng giá trị chữ thường (`'upcoming'`, `'showing'`, `'ended'`) để so sánh, nhưng `movie.status` luôn ở dạng chữ IN HOA (`'UPCOMING'`, `'SHOWING'`, `'ENDED'`). Mọi case đều rơi vào nhánh `default`, trả về `border-gray-500`. |
| MOV-32 | Kiểm tra màu nền nút trạng thái đang được chọn trên trang chi tiết phim | A. GUI Testing | Phim đang ở trạng thái "SHOWING" | 1. Mở trang chi tiết phim trạng thái SHOWING.<br>2. Quan sát nút "showing" trong mục "Thay đổi trạng thái" – nút đang ở trạng thái `disabled` (đã chọn).<br>3. Kiểm tra màu nền hiển thị của nút "showing" đang chọn. | Nút "showing" đang chọn phải nổi bật với nền màu xanh lá (`bg-green-500`) chữ trắng, thể hiện rõ trạng thái đang hoạt động. | Nút "showing" đang chọn hiển thị nền xám (`bg-gray-500`) chữ trắng, không phân biệt được với các trạng thái khác. | **FAIL** | **BUG (Major)**: Dòng 206, biểu thức ``getStatusColorClasses(status)`` nhận tham số `status` ở dạng chữ thường (`'showing'`), nhưng hàm `getStatusColorClasses()` chỉ kiểm tra `'UPCOMMING'`, `'SHOWING'`, `'ENDED'` (chữ HOA). Nút đang chọn luôn nhận màu xám mặc định. |
| MOV-33 | Kiểm tra màu badge trạng thái phim ENDED trên trang chi tiết | A. GUI Testing | Phim có trạng thái "ENDED" (Đã kết thúc) | 1. Mở trang chi tiết bộ phim có trạng thái ENDED.<br>2. Quan sát badge trạng thái bên dưới tên phim. | Badge trạng thái phải hiển thị nền đỏ (`bg-red-500`) chữ trắng để đại diện cho trạng thái "Đã kết thúc". | Badge hiển thị đúng nền đỏ (`bg-red-500`) chữ trắng "ENDED". | **PASS** | Trạng thái `'ENDED'` khớp với nhánh switch trong `getStatusColorClasses()`. |
| MOV-34 | Kiểm tra màu badge trạng thái phim SHOWING trên trang chi tiết | A. GUI Testing | Phim có trạng thái "SHOWING" (Đang chiếu) | 1. Mở trang chi tiết bộ phim có trạng thái SHOWING.<br>2. Quan sát badge trạng thái bên dưới tên phim. | Badge trạng thái phải hiển thị nền xanh lá (`bg-green-500`) chữ trắng. | Badge hiển thị đúng nền xanh lá (`bg-green-500`) chữ trắng "SHOWING". | **PASS** | Trạng thái `'SHOWING'` khớp đúng nhánh switch. |
| MOV-35 | Thêm phim – Nhập URL poster không hợp lệ (chuỗi bất kỳ, không phải URL) | B. Field Validation | Form thêm phim mới:<br>- Tên phim: "Phim Test URL"<br>- Thời lượng: 100<br>- URL Poster: "day-khong-phai-url" | 1. Mở form "Thêm phim mới".<br>2. Nhập tên phim hợp lệ, thời lượng = 100.<br>3. Nhập chuỗi "day-khong-phai-url" vào ô "URL Poster".<br>4. Bấm "Tạo phim mới". | Trình duyệt kích hoạt kiểm tra kiểu `type="url"`, chặn gửi form và hiện tooltip cảnh báo URL không hợp lệ. | Trình duyệt chặn submit form và hiển thị tooltip "Please enter a URL" tại ô URL Poster. Phim không được tạo. | **PASS** | Thuộc tính `type="url"` trên Input component hoạt động đúng. |
| MOV-36 | Thêm phim – Nhập URL trailer không hợp lệ | B. Field Validation | Form thêm phim mới:<br>- Tên phim: "Phim Test Trailer"<br>- Thời lượng: 100<br>- URL Trailer: "abc123" | 1. Mở form "Thêm phim mới".<br>2. Nhập tên phim hợp lệ, thời lượng = 100.<br>3. Nhập chuỗi "abc123" vào ô "URL Trailer".<br>4. Bấm "Tạo phim mới". | Trình duyệt chặn gửi form do kiểu `type="url"` không hợp lệ. | Trình duyệt chặn submit form và hiển thị tooltip cảnh báo URL không hợp lệ tại ô URL Trailer. | **PASS** | |
| MOV-37 | Thêm phim – Nhập thời lượng là số thập phân (VD: 1.5) | B. Field Validation | Form thêm phim mới:<br>- Tên phim: "Phim Thập Phân"<br>- Thời lượng: 1.5 | 1. Mở form thêm phim.<br>2. Nhập tên phim hợp lệ.<br>3. Nhập giá trị "1.5" vào ô thời lượng.<br>4. Bấm "Tạo phim mới". | Hệ thống chấp nhận hoặc tự động làm tròn thành số nguyên (VD: 1 phút), không gây lỗi hệ thống. | Hệ thống chấp nhận giá trị, `parseInt("1.5", 10)` chuyển thành 1 phút. Phim tạo thành công với thời lượng 1 phút hiển thị "0h 1m". | **PASS** | Hàm `parseInt()` tự cắt phần thập phân. |
| MOV-38 | Thêm phim – Nhập thời lượng rất lớn (VD: 999999 phút) | B. Field Validation | Form thêm phim mới:<br>- Tên phim: "Phim Siêu Dài"<br>- Thời lượng: 999999 | 1. Mở form thêm phim.<br>2. Nhập tên phim hợp lệ, thời lượng = 999999.<br>3. Bấm "Tạo phim mới". | Hệ thống nên kiểm tra và từ chối giá trị bất thường, hoặc hiển thị cảnh báo thời lượng quá lớn. | Phim được tạo thành công, card hiển thị thời lượng "16666h 39m" trên danh sách. Hệ thống không giới hạn giá trị tối đa hợp lý cho thời lượng phim. | **FAIL** | **BUG (Minor)**: Không có validation giới hạn trên (max) cho trường thời lượng cả ở client lẫn server. Giá trị nghiệp vụ hợp lý thường dưới 600 phút (10 tiếng). Backend entity chỉ kiểm tra `Duration <= 0` chứ không kiểm tra giá trị tối đa. |
| MOV-39 | Chỉnh sửa phim – Thay đổi trạng thái trực tiếp trên form từ SHOWING về UPCOMING | B. Field Validation | Form chỉnh sửa phim đang ở trạng thái SHOWING:<br>- Đổi dropdown trạng thái sang "Sắp chiếu" | 1. Mở form chỉnh sửa một phim đang ở trạng thái SHOWING.<br>2. Tại dropdown "Trạng thái", đổi từ "Đang chiếu" sang "Sắp chiếu".<br>3. Bấm "Cập nhật phim". | Hệ thống từ chối cập nhật và hiển thị thông báo lỗi chuyển trạng thái không hợp lệ theo quy tắc nghiệp vụ: SHOWING chỉ được chuyển sang ENDED. | Hệ thống hiển thị thông báo lỗi trong khung đỏ trên form: "Không thể cập nhật phim" do backend trả về lỗi `invalid status transition`. Dữ liệu không thay đổi. | **PASS** | Backend validate đúng quy tắc chuyển trạng thái. Tuy nhiên, dropdown trên form vẫn cho phép chọn mọi trạng thái mà không ẩn/vô hiệu hóa các lựa chọn không hợp lệ. |
| MOV-40 | Kiểm tra chuyển trạng thái hợp lệ: SHOWING → ENDED trên trang chi tiết | D. Business Flow Testing | Phim đang ở trạng thái "SHOWING" | 1. Mở trang chi tiết bộ phim SHOWING.<br>2. Bấm nút "ended" trong mục "Thay đổi trạng thái". | Trạng thái chuyển thành "ENDED", badge đổi sang màu đỏ, nút "ended" bị vô hiệu hóa (disabled). | Trạng thái chuyển thành công sang ENDED, badge đổi chính xác sang nền đỏ chữ trắng. | **PASS** | Chuyển đổi trạng thái hợp lệ theo vòng đời phim. |
| MOV-41 | Kiểm tra chuyển trạng thái không hợp lệ: ENDED → SHOWING | D. Business Flow Testing | Phim đang ở trạng thái "ENDED" | 1. Mở trang chi tiết bộ phim có trạng thái ENDED.<br>2. Bấm nút "showing" trong mục "Thay đổi trạng thái". | Hệ thống từ chối chuyển trạng thái, hiển thị thông báo lỗi. Trạng thái phim vẫn giữ nguyên ENDED. | Hộp thoại alert hiển thị: "Failed to update movie status: invalid status transition". Trạng thái vẫn giữ nguyên ENDED. | **PASS** | Backend chặn đúng: Phim ENDED không được quay lại SHOWING hoặc UPCOMING. Tuy nhiên thông báo bằng tiếng Anh (đã ghi nhận ở MOV-29). |
| MOV-42 | Kiểm tra chuyển trạng thái không hợp lệ: ENDED → UPCOMING | D. Business Flow Testing | Phim đang ở trạng thái "ENDED" | 1. Mở trang chi tiết bộ phim ENDED.<br>2. Bấm nút "upcoming". | Hệ thống từ chối, trạng thái giữ nguyên ENDED. | Alert: "Failed to update movie status: invalid status transition". Trạng thái vẫn ENDED. | **PASS** | Chặn đúng quy tắc: ENDED là trạng thái cuối cùng, không chuyển ngược. |
| MOV-43 | Kiểm tra hộp thoại xác nhận xóa phim – Nhấn Hủy (Cancel) | D. Business Flow Testing | Phim thử nghiệm bất kỳ trên trang chi tiết | 1. Mở trang chi tiết phim cần kiểm tra.<br>2. Bấm nút "Delete Movie".<br>3. Khi hộp thoại xác nhận hiện lên, nhấn nút "Cancel" (Hủy). | Phim không bị xóa. Người dùng vẫn ở lại trang chi tiết phim, mọi thông tin giữ nguyên. | Hộp thoại đóng lại, phim vẫn tồn tại nguyên vẹn trên trang chi tiết. Không có gì thay đổi. | **PASS** | Chức năng hủy xóa hoạt động chuẩn xác. |
| MOV-44 | Kiểm tra hộp thoại xóa phim hiển thị bằng tiếng Anh | A. GUI Testing | Nhấn nút "Delete Movie" trên trang chi tiết phim | 1. Mở trang chi tiết phim bất kỳ.<br>2. Bấm nút "Delete Movie".<br>3. Quan sát nội dung hộp thoại xác nhận. | Hộp thoại xác nhận xóa phải hiển thị bằng tiếng Việt cho giao diện Admin tiếng Việt: VD "Bạn có chắc chắn muốn xóa phim này?". | Hộp thoại `window.confirm()` hiển thị nội dung bằng tiếng Anh: "Are you sure you want to delete this movie?". | **FAIL** | **BUG (Minor)**: Dòng 45 trong `MovieDetailPage.jsx`, chuỗi confirm bằng tiếng Anh `'Are you sure you want to delete this movie?'`, không đồng nhất với ngôn ngữ tiếng Việt của giao diện Admin. |
| MOV-45 | Tìm kiếm phim – Nhập ký tự đặc biệt XSS (`<script>alert(1)</script>`) | B. Field Validation | Ô tìm kiếm: Nhập chuỗi `<script>alert(1)</script>` | 1. Tại ô tìm kiếm phim, nhập chuỗi `<script>alert(1)</script>`.<br>2. Bấm nút "Tìm kiếm". | Hệ thống xử lý an toàn: Không thực thi mã JavaScript, hiển thị kết quả "Không tìm thấy phim nào" hoặc danh sách rỗng. Chuỗi được encode đúng cách. | Hệ thống hiển thị "Không tìm thấy phim nào" và dòng "Không có phim nào khớp với từ khóa" kèm chuỗi đã được escape an toàn bằng React JSX. Không có popup alert xảy ra. | **PASS** | React tự động escape HTML entities, tránh XSS. Hàm `encodeURIComponent` trong API client cũng mã hóa URL an toàn. |
| MOV-46 | Kiểm tra hiển thị khi poster URL bị hỏng (broken image) | C. Functional & Display Testing | Phim có poster_url trỏ đến link ảnh không tồn tại (VD: `https://example.com/404.jpg`) | 1. Tạo hoặc chỉnh sửa một phim, nhập URL poster không tồn tại.<br>2. Quay về trang danh sách phim.<br>3. Quan sát Card phim vừa tạo/sửa. | Card phim vẫn hiển thị bình thường với vùng ảnh poster trắng hoặc icon ảnh lỗi, không làm vỡ bố cục card. | Card phim vẫn hiển thị, vùng poster xuất hiện icon ảnh lỗi (broken image) của trình duyệt. Bố cục card giữ nguyên do CSS `h-72 overflow-hidden bg-gray-100` giữ chiều cao cố định. | **PASS** | Giao diện có fallback background xám cho vùng poster. |
| MOV-47 | Kiểm tra hiển thị khi phim không có poster URL | C. Functional & Display Testing | Phim không có trường poster_url (rỗng) | 1. Tạo phim chỉ với trường bắt buộc, không nhập poster URL.<br>2. Mở trang danh sách phim, quan sát Card phim.<br>3. Mở trang chi tiết phim. | Card phim không hiển thị vùng ảnh poster (do điều kiện render `movie.poster_url &&`), nội dung thông tin chiếm hết diện tích card. | Card phim bỏ qua phần poster, chỉ hiển thị phần thông tin. Trên trang chi tiết cũng bỏ qua khối poster. Không có lỗi giao diện xảy ra. | **PASS** | Render có điều kiện hoạt động tốt. |
| MOV-48 | Tìm kiếm phim – Nhập từ khóa có khoảng trắng đầu và cuối | C. Functional & Display Testing | Ô tìm kiếm: Nhập "  Avengers  " (có 2 khoảng trắng đầu/cuối) | 1. Tại ô tìm kiếm, nhập "  Avengers  " (có khoảng trắng thừa).<br>2. Bấm "Tìm kiếm". | Hệ thống tìm kiếm theo từ khóa chuẩn sau khi loại bỏ khoảng trắng thừa, trả về kết quả phim "Avengers: Endgame". | Hệ thống truyền nguyên chuỗi có khoảng trắng lên server. Backend tìm kiếm SQL ILIKE vẫn trả về kết quả "Avengers: Endgame" do khoảng trắng không ảnh hưởng pattern matching. | **PASS** | Mặc dù client không trim, kết quả vẫn đúng nhờ backend ILIKE. |
| MOV-49 | Kiểm tra ngôn ngữ nhãn trên Card phim trong trang danh sách | A. GUI Testing | Card phim bất kỳ trên trang Quản lý phim | 1. Mở trang Quản lý phim.<br>2. Quan sát các nhãn trên Card phim: "Director:", "Genre:", "Duration:", "Release:". | Các nhãn trên Card phim phải hiển thị bằng tiếng Việt cho giao diện quản trị Việt hóa, VD: "Đạo diễn:", "Thể loại:", "Thời lượng:", "Ngày phát hành:". | Các nhãn trên Card phim hiển thị bằng tiếng Anh: "Director:", "Genre:", "Duration:", "Release:". | **FAIL** | **BUG (Minor)**: File `MovieCard.jsx` (admin) dùng nhãn tiếng Anh (`Director:`, `Genre:`, `Duration:`, `Release:`) trong khi toàn bộ giao diện Admin còn lại sử dụng tiếng Việt. Thiếu nhất quán ngôn ngữ. |
| MOV-50 | Kiểm tra nút "Back to Movies" và "Edit Movie" hiển thị bằng tiếng Anh trên trang chi tiết | A. GUI Testing | Trang chi tiết phim (`/admin/movies/:id`) | 1. Mở trang chi tiết phim bất kỳ.<br>2. Quan sát các nút điều hướng: "← Back to Movies", "Edit Movie", "Delete Movie". | Các nút điều hướng phải hiển thị bằng tiếng Việt cho giao diện Admin, VD: "← Quay lại danh sách phim", "Chỉnh sửa phim", "Xóa phim". | Các nút hiển thị bằng tiếng Anh: "← Back to Movies", "Edit Movie", "Delete Movie". | **FAIL** | **BUG (Minor)**: File `MovieDetailPage.jsx` sử dụng nhãn nút bằng tiếng Anh (`Back to Movies`, `Edit Movie`, `Delete Movie`), trong khi form thêm/sửa phim (`MovieFormPage.jsx`) đã Việt hóa đầy đủ. Thiếu nhất quán ngôn ngữ trên cùng một luồng chức năng. |

### 1.3. Kết luận phân hệ Quản lý phim mở rộng

- **Tổng số test case mở rộng:** 20 ca.
- **Số ca Đạt (PASS):** 14 ca (70,0%).
- **Số ca Không đạt (FAIL):** 06 ca (30,0%) gồm:
  - **MOV-31 (Major):** Hàm `getStatusBorderClasses()` sử dụng giá trị chữ thường để so sánh với `movie.status` chữ HOA, khiến viền 3 nút trạng thái luôn hiển thị xám mặc định.
  - **MOV-32 (Major):** Hàm `getStatusColorClasses()` nhận tham số chữ thường từ nút trạng thái nhưng so sánh chữ HOA, khiến nút đang chọn luôn có nền xám.
  - **MOV-38 (Minor):** Không có kiểm tra giới hạn trên (max) cho thời lượng phim, cho phép nhập giá trị bất hợp lý (VD: 999999 phút).
  - **MOV-44 (Minor):** Hộp thoại xác nhận xóa phim hiển thị bằng tiếng Anh, không đồng nhất với giao diện tiếng Việt.
  - **MOV-49 (Minor):** Nhãn trên Card phim ("Director:", "Genre:"...) hiển thị tiếng Anh thay vì tiếng Việt.
  - **MOV-50 (Minor):** Nút "Back to Movies", "Edit Movie", "Delete Movie" trên trang chi tiết phim hiển thị tiếng Anh.

---

## 2. QUẢN LÝ GHẾ – MỞ RỘNG (Seat Management – Extended)

### 2.1. Thống kê kết quả

| Trạng thái | Số lượng | Tỷ lệ | Định nghĩa |
|---|---|---|---|
| **TESTED** | 20 | 100% | Tổng số test case hệ thống đã thực hiện |
| **PASS** | 16 | 80,0% | Test case cho kết quả thực tế đúng như kỳ vọng |
| **FAIL** | 4 | 20,0% | Test case phát hiện sai lệch hoặc lỗi giao diện / nghiệp vụ |
| **N/A** | 0 | 0% | Test case không áp dụng được |

### 2.2. Bảng chi tiết Test Case Quản lý ghế mở rộng (20 ca)

| ID | Tên chức năng kiểm thử | Loại test case | Dữ liệu đầu vào | Các bước thực hiện | Kết quả mong đợi | Kết quả thực tế | Trạng thái | Ghi chú |
|---|---|---|---|---|---|---|---|---|
| SEAT-31 | Thêm ghế – Nhập hàng ghế bằng chữ thường (VD: "a" thay vì "A") | B. Field Validation | Form thêm ghế mới:<br>- Phòng: Phòng 1<br>- Hàng ghế: "a" (chữ thường)<br>- Số ghế: "99"<br>- Loại: Thường | 1. Mở form thêm ghế mới.<br>2. Chọn Phòng 1, nhập Hàng ghế "a" (chữ thường), Số ghế "99".<br>3. Bấm "Tạo mới". | Hệ thống nên tự động chuyển thành chữ HOA "A" hoặc từ chối nhập chữ thường để đảm bảo tính nhất quán dữ liệu (sơ đồ lưới chỉ hiểu hàng A-O chữ HOA). | Hệ thống chấp nhận và tạo ghế "a99" thành công. Ghế này KHÔNG hiển thị trên sơ đồ Lưới (Grid View) vì `createSeatGrid()` chỉ tạo lưới cho các hàng A-O chữ HOA. Tuy nhiên ghế vẫn xuất hiện trong chế độ Bảng. | **FAIL** | **BUG (Major)**: Frontend không chuẩn hóa (`.toUpperCase()`) giá trị hàng ghế trước khi gửi tạo. Backend cũng không transform. Hậu quả: ghế tồn tại trong DB nhưng "vô hình" trên sơ đồ lưới, gây sai lệch số liệu thống kê giữa 2 chế độ xem. |
| SEAT-32 | Thêm ghế – Nhập số ghế bắt đầu bằng số 0 (VD: "007") | B. Field Validation | Form thêm ghế mới:<br>- Phòng: Phòng 1<br>- Hàng ghế: "Z"<br>- Số ghế: "007" | 1. Chọn Phòng 1, nhập Hàng ghế "Z", Số ghế "007".<br>2. Bấm "Tạo mới". | Hệ thống chuẩn hóa và tạo ghế Z07 hoặc chấp nhận Z007 và hiển thị nhất quán. | Ghế Z007 được tạo thành công. Trên sơ đồ Lưới, `parseInt("007")` = 7, nên ghế hiển thị ở vị trí cột 7 với nhãn "07" (do `padStart(2,'0')`). Trên Bảng hiển thị vị trí "Z007". Hai chế độ xem hiển thị nhãn không nhất quán. | **FAIL** | **BUG (Minor)**: Không chuẩn hóa giá trị `seat_number` trước khi lưu (ví dụ trim leading zeros). Trên Grid hiển thị "07" nhưng trên Table hiển thị "007", gây nhầm lẫn cho người dùng. |
| SEAT-33 | Thêm ghế – Nhập hàng ghế là ký tự số (VD: "1") | B. Field Validation | Form thêm ghế mới:<br>- Phòng: Phòng 1<br>- Hàng ghế: "1"<br>- Số ghế: "01" | 1. Chọn Phòng 1, nhập Hàng ghế "1", Số ghế "01".<br>2. Bấm "Tạo mới". | Hệ thống từ chối hoặc cảnh báo hàng ghế phải là chữ cái (A-O), không được dùng số. | Hệ thống chấp nhận tạo ghế "101" thành công. Ghế này không hiển thị trên sơ đồ Lưới (Grid chỉ có hàng A-O). Tuy nhiên trên Bảng vẫn liệt kê ghế ở hàng "1". | **FAIL** | **BUG (Minor)**: Không có validation kiểm tra định dạng hàng ghế phải là chữ cái. Hệ thống cho phép hàng ghế là số, ký tự đặc biệt, v.v. gây dữ liệu không hợp lệ và bất khả kiến trên sơ đồ. |
| SEAT-34 | Thêm ghế – Nhập ký tự đặc biệt vào trường Hàng ghế | B. Field Validation | Form thêm ghế mới:<br>- Phòng: Phòng 1<br>- Hàng ghế: "@#$"<br>- Số ghế: "01" | 1. Chọn Phòng 1, nhập Hàng ghế "@#$", Số ghế "01".<br>2. Bấm "Tạo mới". | Hệ thống từ chối tạo ghế và hiển thị lỗi hàng ghế không hợp lệ. | Hệ thống tạo ghế thành công với hàng ghế "@#$". Ghế này không hiển thị trên sơ đồ Lưới nhưng có mặt trên Bảng. | **FAIL** | **BUG (Minor)**: Cùng nguyên nhân với SEAT-33 – thiếu kiểm tra format cho trường `row_number`. Cho phép ký tự đặc biệt vào hệ thống. |
| SEAT-35 | Kiểm tra thay đổi trạng thái ghế nhanh: Có sẵn → Đã đặt (OCCUPIED) trên Bảng | C. Functional & Display Testing | Ghế bất kỳ đang "Có sẵn" trên chế độ Bảng | 1. Tại cột "Trạng thái" của một ghế đang "Có sẵn", mở dropdown.<br>2. Chọn trạng thái "Đã đặt".<br>3. Tải lại trang (F5). | Trạng thái ghế chuyển sang "Đã đặt" (nhãn đỏ), lưu bền vững sau khi refresh. | Trạng thái chuyển thành công sang "Đã đặt", dropdown đổi màu nền đỏ. Sau khi F5 tải lại, trạng thái vẫn giữ "Đã đặt". | **PASS** | |
| SEAT-36 | Kiểm tra thay đổi trạng thái ghế nhanh: Bảo trì → Có sẵn trên Bảng | C. Functional & Display Testing | Ghế đang "Bảo trì" trên chế độ Bảng | 1. Tại ghế đang Bảo trì, đổi dropdown sang "Có sẵn".<br>2. F5 tải lại trang. | Trạng thái ghế chuyển lại "Có sẵn" (nhãn xanh lá), dữ liệu lưu bền vững. | Trạng thái chuyển thành công, hiển thị nhãn xanh lá. Bền vững sau reload. | **PASS** | |
| SEAT-37 | Kiểm tra thay đổi trạng thái ghế: Có sẵn → Bị khóa (BLOCKED) | C. Functional & Display Testing | Ghế đang "Có sẵn" trên chế độ Bảng | 1. Đổi dropdown trạng thái ghế sang "Bị khóa". | Trạng thái chuyển sang "Bị khóa" (nhãn xám), lưu thành công. | Trạng thái cập nhật thành công, dropdown hiển thị nhãn xám "Bị khóa". | **PASS** | |
| SEAT-38 | Lọc ghế theo trạng thái "Bị khóa" (BLOCKED) | C. Functional & Display Testing | Chế độ Bảng: Chọn dropdown Trạng thái = "Bị khóa" | 1. Tại dropdown "Tất cả trạng thái", chọn "Bị khóa". | Bảng chỉ hiển thị các ghế đang ở trạng thái "Bị khóa". | Bảng lọc chính xác, chỉ hiển thị các ghế có trạng thái "Bị khóa" (nhãn xám). | **PASS** | |
| SEAT-39 | Lọc ghế theo loại "Thường" (REGULAR) | C. Functional & Display Testing | Chế độ Bảng: Chọn dropdown Loại ghế = "Thường" | 1. Tại dropdown "Tất cả loại ghế", chọn "Thường". | Bảng chỉ hiển thị các ghế Thường (REGULAR). | Bảng hiển thị chính xác danh sách các ghế Thường. | **PASS** | |
| SEAT-40 | Kiểm tra xóa ghế – Nhấn Hủy (Cancel) trên hộp thoại xác nhận | D. Business Flow Testing | Ghế bất kỳ trên chế độ Bảng | 1. Bấm icon thùng rác tại dòng ghế bất kỳ.<br>2. Khi hộp thoại xác nhận hiện lên, nhấn "Cancel". | Ghế không bị xóa, vẫn tồn tại trong danh sách. | Hộp thoại đóng lại, ghế vẫn có mặt nguyên vẹn trong bảng. | **PASS** | |
| SEAT-41 | Tạo ghế VIP thành công và kiểm tra hiển thị trên sơ đồ Lưới | D. Business Flow Testing | Form thêm ghế:<br>- Phòng: Phòng 1<br>- Hàng: "F", Số: "99"<br>- Loại: VIP | 1. Tạo ghế VIP mới tại hàng F, số 99, Phòng 1.<br>2. Chuyển sang chế độ Lưới xem sơ đồ Phòng 1.<br>3. Kiểm tra vị trí F99 trên sơ đồ. | Ghế F99 xuất hiện trên sơ đồ lưới tại hàng F với màu vàng (đại diện cho VIP). | Vì sơ đồ Lưới chỉ tạo tối đa 16 cột (seat_number 1-16) cho hàng thường, ghế số 99 nằm ngoài phạm vi hiển thị và KHÔNG xuất hiện trên Grid. Tuy nhiên ghế tồn tại trong Bảng. | **PASS** | Hành vi hợp lý: Sơ đồ Lưới thiết kế cho phòng chiếu chuẩn (tối đa 16 ghế/hàng). Ghế ngoài phạm vi sẽ chỉ quản lý qua Bảng. |
| SEAT-42 | Kiểm tra giao diện khi chưa chọn phòng chiếu ở chế độ Lưới | A. GUI Testing | Chế độ Lưới, dropdown phòng chiếu = "-- Chọn phòng --" | 1. Chuyển sang tab "Lưới".<br>2. Đặt dropdown phòng chiếu về "-- Chọn phòng --".<br>3. Quan sát khu vực hiển thị sơ đồ. | Hiển thị thông báo thân thiện yêu cầu người dùng chọn phòng chiếu. | Hiển thị đúng: "Vui lòng chọn phòng chiếu" và dòng phụ "Chọn phòng từ dropdown phía trên để xem sơ đồ ghế". | **PASS** | |
| SEAT-43 | Kiểm tra hiển thị ghế Đôi (COUPLE) trên sơ đồ Lưới | A. GUI Testing | Sơ đồ Lưới của Phòng có ghế Đôi (hàng M, N, O) | 1. Mở sơ đồ Lưới của phòng có ghế Đôi.<br>2. Quan sát các hàng M, N, O trên sơ đồ. | Ghế Đôi hiển thị rộng hơn ghế thường (chiều rộng gấp đôi), có icon sofa (FaCouch) thay vì số, và màu hồng đặc trưng. | Ghế Đôi hiển thị đúng: Rộng hơn (`w-12` so với `w-8` của ghế thường), có icon sofa bên trong, nền hồng (`bg-pink-200 border-pink-400`). Khoảng cách giữa các ghế đôi cũng lớn hơn (`gap-3` so với `gap-1`). | **PASS** | |
| SEAT-44 | Kiểm tra Preview (Xem trước) khi chưa nhập đủ 3 trường bắt buộc | A. GUI Testing | Form thêm ghế: Chỉ chọn Phòng 1, chưa nhập Hàng và Số ghế | 1. Mở form thêm ghế mới.<br>2. Chỉ chọn Phòng 1 từ dropdown.<br>3. Quan sát phần dưới form. | Khung "Xem trước" KHÔNG hiển thị vì chưa đủ 3 thông tin (Phòng + Hàng + Số ghế). | Khung Xem trước không xuất hiện. Điều kiện render `formData.room_id && formData.row_number && formData.seat_number` đúng – cần đủ 3 trường mới hiện. | **PASS** | |
| SEAT-45 | Chỉnh sửa ghế – Đổi loại ghế từ VIP sang Đôi (COUPLE) | D. Business Flow Testing | Ghế VIP Z99 Phòng 1: Chuyển sang loại "Đôi" | 1. Mở form chỉnh sửa ghế VIP Z99.<br>2. Đổi dropdown Loại ghế từ "VIP" sang "Đôi".<br>3. Bấm "Cập nhật". | Ghế Z99 chuyển sang loại Đôi thành công. | Cập nhật thành công, bảng hiển thị loại ghế "Đôi" cho Z99. | **PASS** | |
| SEAT-46 | Kiểm tra tìm kiếm ghế với từ khóa không tồn tại | C. Functional & Display Testing | Ô tìm kiếm: Nhập "XYZ999" | 1. Tại thanh tìm kiếm ghế, gõ "XYZ999". | Bảng hiển thị kết quả trống kèm thông báo "Không có ghế nào". | Bảng hiển thị danh sách rỗng và dòng "Không có ghế nào" ở giữa trang. | **PASS** | |
| SEAT-47 | Kiểm tra hiển thị khi lọc kết hợp không có kết quả | C. Functional & Display Testing | Chế độ Bảng: Lọc Loại ghế "Đôi" + Trạng thái "Bảo trì" (tổ hợp không có dữ liệu) | 1. Chọn Loại ghế = "Đôi".<br>2. Chọn Trạng thái = "Bảo trì".<br>3. Quan sát bảng. | Bảng hiển thị trống với thông báo "Không có ghế nào". | Bảng hiển thị "Không có ghế nào" khi tổ hợp lọc không khớp dữ liệu. | **PASS** | |
| SEAT-48 | Kiểm tra chuyển trang bảng ghế và quay lại giữ nguyên bộ lọc | C. Functional & Display Testing | Chế độ Bảng: Lọc theo Loại ghế = "Thường", nhiều trang | 1. Chọn Loại ghế = "Thường".<br>2. Chuyển sang trang 2.<br>3. Kiểm tra dropdown Loại ghế vẫn giữ giá trị "Thường". | Bộ lọc giữ nguyên giá trị khi chuyển trang, danh sách trang 2 vẫn chỉ hiển thị ghế Thường. | Bộ lọc giữ nguyên "Thường" khi chuyển trang, danh sách trang 2 chỉ có ghế Thường. | **PASS** | |
| SEAT-49 | Kiểm tra dropdown phòng chiếu hiển thị loại phòng viết HOA trên form thêm ghế | A. GUI Testing | Dropdown phòng chiếu trên form thêm ghế | 1. Mở form thêm ghế mới.<br>2. Mở dropdown "Phòng chiếu", quan sát nhãn hiển thị. | Nhãn phòng chiếu hiển thị rõ ràng: "Phòng X (LOẠI)", ví dụ "Phòng 1 (STANDARD)". | Dropdown hiển thị đúng: "Phòng 1 (STANDARD)", "Phòng 2 (VIP)"... với loại phòng viết HOA nhờ hàm `.toUpperCase()` trên `room.room_type`. | **PASS** | |
| SEAT-50 | Xóa ghế và kiểm tra tải lại sơ đồ Lưới đồng bộ | D. Business Flow Testing | Ghế thử nghiệm vừa tạo ở các bước trước (ví dụ ghế hàng "a" hoặc "@#$" nếu đã tạo) | 1. Tại chế độ Bảng, tìm ghế thử nghiệm cần dọn dẹp.<br>2. Bấm icon thùng rác, xác nhận xóa.<br>3. Chuyển sang chế độ Lưới của cùng Phòng.<br>4. Kiểm tra ghế đã biến mất. | Ghế bị xóa hoàn toàn khỏi hệ thống. Cả chế độ Bảng lẫn Lưới đều không còn hiển thị ghế đó. | Ghế bị xóa thành công, biến mất khỏi Bảng. Khi chuyển sang Lưới, sơ đồ cũng không hiển thị ghế vừa xóa. Dữ liệu đồng bộ giữa 2 chế độ xem. | **PASS** | |

### 2.3. Kết luận phân hệ Quản lý ghế mở rộng

- **Tổng số test case mở rộng:** 20 ca.
- **Số ca Đạt (PASS):** 16 ca (80,0%).
- **Số ca Không đạt (FAIL):** 04 ca (20,0%) gồm:
  - **SEAT-31 (Major):** Hàng ghế nhập chữ thường không được chuẩn hóa thành chữ HOA, dẫn đến ghế "vô hình" trên sơ đồ Lưới.
  - **SEAT-32 (Minor):** Số ghế nhập "007" không được chuẩn hóa, hiển thị nhãn không nhất quán giữa chế độ Lưới ("07") và Bảng ("007").
  - **SEAT-33 (Minor):** Cho phép nhập hàng ghế bằng ký tự số ("1"), tạo ghế bất khả kiến trên sơ đồ.
  - **SEAT-34 (Minor):** Cho phép nhập ký tự đặc biệt ("@#$") vào trường hàng ghế, không có validation format.

---

## 3. KẾT LUẬN TỔNG THỂ BỔ SUNG

### 3.1. Bảng tổng hợp số lượng và tỷ lệ kết quả kiểm thử mở rộng

| Phân hệ chức năng | Số lượng ca kiểm thử (TESTED) | ĐẠT (PASS) | LỖI (FAIL) | Tỷ lệ ĐẠT (Pass Rate) |
|---|---|---|---|---|
| **1. Quản lý phim (Movie Management) – Mở rộng** | 20 | 14 | 6 | **70,0%** |
| **2. Quản lý ghế (Seat Management) – Mở rộng** | 20 | 16 | 4 | **80,0%** |
| **TỔNG CỘNG MỞ RỘNG** | **40** | **30** | **10** | **75,0%** |

---

### 3.2. Bảng phân bố kết quả theo 4 loại test case (Mở rộng)

| Nhóm loại test case | Quản lý phim | Quản lý ghế | Tổng số ca | Số ca PASS | Số ca FAIL | Tỷ lệ PASS |
|---|---|---|---|---|---|---|
| **A. GUI Testing** | 5 | 4 | **9** | 6 | 3 | 66,7% |
| **B. Field Validation** | 5 | 4 | **9** | 5 | 4 | 55,6% |
| **C. Functional & Display Testing** | 4 | 7 | **11** | 11 | 0 | 100% |
| **D. Business Flow Testing** | 6 | 5 | **11** | 8 | 3 | 72,7% |
| **TỔNG CỘNG** | **20** | **20** | **40** | **30** | **10** | **75,0%** |

---

### 3.3. Phân loại lỗi phát hiện theo mức độ nghiêm trọng (Mở rộng)

| Mức độ nghiêm trọng (Severity) | Số lượng | Mã test case liên quan | Diễn giải chi tiết |
|---|---|---|---|
| **Critical (Nghiêm trọng cao)** | **0** | Không có | Không có lỗi gây treo ứng dụng, sập dịch vụ hoặc gián đoạn luồng nghiệp vụ. |
| **Major (Nghiêm trọng vừa)** | **3** | MOV-31, MOV-32, SEAT-31 | - **MOV-31 & MOV-32:** Lỗi không nhất quán chữ hoa/thường (case-sensitivity) trong hàm `getStatusBorderClasses()` và `getStatusColorClasses()` trên trang chi tiết phim, khiến viền và nền các nút chuyển trạng thái luôn hiển thị xám mặc định, ảnh hưởng trực tiếp đến trải nghiệm thao tác quản trị.<br>- **SEAT-31:** Hàng ghế nhập chữ thường không được chuẩn hóa, tạo ra dữ liệu "vô hình" trên sơ đồ lưới – ghế tồn tại trong DB nhưng admin không thấy trên giao diện trực quan. |
| **Minor (Lỗi nhỏ / Giao diện & UX)** | **7** | MOV-38, MOV-44, MOV-49, MOV-50, SEAT-32, SEAT-33, SEAT-34 | - **MOV-38:** Không giới hạn tối đa thời lượng phim.<br>- **MOV-44:** Hộp thoại xóa phim bằng tiếng Anh.<br>- **MOV-49:** Nhãn Card phim bằng tiếng Anh.<br>- **MOV-50:** Nút điều hướng trang chi tiết phim bằng tiếng Anh.<br>- **SEAT-32:** Số ghế "007" hiển thị không nhất quán giữa 2 chế độ xem.<br>- **SEAT-33:** Cho phép hàng ghế bằng ký tự số.<br>- **SEAT-34:** Cho phép ký tự đặc biệt làm hàng ghế. |

---

### 3.4. Tổng hợp toàn bộ hệ thống (Gộp 60 ca cũ + 40 ca mở rộng)

| Phân hệ chức năng | Tổng ca kiểm thử | ĐẠT (PASS) | LỖI (FAIL) | Tỷ lệ ĐẠT |
|---|---|---|---|---|
| **Quản lý phim (Movie Management)** | 50 | 42 | 8 | **84,0%** |
| **Quản lý ghế (Seat Management)** | 50 | 44 | 6 | **88,0%** |
| **TỔNG CỘNG** | **100** | **86** | **14** | **86,0%** |

> **Nhận xét tổng thể sau mở rộng:** Với 100 test case (60 cũ + 40 mới), tỷ lệ PASS tổng thể là 86,0%. Các lỗi mới phát hiện chủ yếu thuộc nhóm GUI Testing và Field Validation – nơi tập trung kiểm tra các kịch bản biên mà đợt 1 chưa bao phủ. Đặc biệt, cụm lỗi case-sensitivity (MOV-31, MOV-32) và thiếu chuẩn hóa dữ liệu đầu vào cho trường hàng ghế (SEAT-31, SEAT-33, SEAT-34) là những vấn đề cần ưu tiên khắc phục để đảm bảo trải nghiệm quản trị đồng nhất và chính xác.
