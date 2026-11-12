# BÁO CÁO KIỂM THỬ HỆ THỐNG (SYSTEM TEST)

## Thông tin chung

| Mục | Nội dung |
|---|---|
| **Tên dự án** | Xây dựng website quản lý rạp chiếu phim (Cinema Management System – HQ Cinema) |
| **Nội dung kiểm thử** | System Test – Kiểm thử hệ thống giao diện người dùng (UI System Test) cho chức năng Quản lý phim (Movie Management) và Quản lý ghế (Seat Management) |
| **Phiên bản** | 1.0 |
| **Người kiểm thử** | Phạm Minh Đức |
| **Ngày thực hiện** | 11/11/2026 |
| **Môi trường** | Máy trạm kiểm thử, Docker containers hệ thống, Database PostgreSQL, Redis Cache |
| **Trình duyệt / Thiết bị** | Google Chrome (phiên bản mới nhất) trên hệ điều hành Windows – Màn hình độ phân giải Full HD (1920x1080) |

> **Phạm vi kiểm thử:** Toàn bộ quá trình kiểm thử hệ thống được thực hiện trực tiếp trên giao diện đồ họa người dùng (GUI) của Trang quản trị (Admin Portal) HQ Cinema thông qua trình duyệt web. Không sử dụng các công cụ gọi API trực tiếp (như Postman, cURL, script API). Mọi thao tác, dữ liệu nhập và kết quả ghi nhận đều bám sát trải nghiệm người dùng thực tế trên các màn hình chức năng.

---

## Phần đánh giá

### Thông tin tổng hợp kết quả kiểm thử

| Trạng thái | Số lượng | Tỷ lệ |
|---|---|---|
| **TỔNG TEST CASE (TESTED)** | 60 | 100% |
| **ĐẠT (PASS)** | 56 | 93,3% |
| **LỖI (FAIL)** | 4 | 6,7% |
| **KHÔNG ÁP DỤNG (N/A)** | 0 | 0% |

> **Nhận xét chung:** Hệ thống HQ Cinema ở tầng giao diện vận hành ổn định trên hầu hết các luồng chức năng chính (tỷ lệ Pass đạt 93,3%). Các lỗi phát hiện được đều ở mức độ Minor (lỗi giao diện hiển thị màu sắc badge trạng thái, thiếu hàm cắt khoảng trắng trim() ở ô nhập liệu và thông báo lỗi nghiệp vụ hiển thị bằng tiếng Anh chưa thân thiện). Không ghi nhận lỗi nghiêm trọng làm tắc nghẽn luồng nghiệp vụ.

---

## 1. QUẢN LÝ PHIM (Movie Management)

### 1.1. Thống kê kết quả

| Trạng thái | Số lượng | Tỷ lệ | Định nghĩa |
|---|---|---|---|
| **TESTED** | 30 | 100% | Tổng số test case hệ thống đã thực hiện |
| **PASS** | 28 | 93,3% | Test case cho kết quả thực tế đúng như kỳ vọng |
| **FAIL** | 2 | 6,7% | Test case phát hiện sai lệch hoặc lỗi giao diện / nghiệp vụ |
| **N/A** | 0 | 0% | Test case không áp dụng được |

### 1.2. Bảng chi tiết Test Case Quản lý phim (30 ca)

| ID | Tên chức năng kiểm thử | Loại test case | Dữ liệu đầu vào | Các bước thực hiện | Kết quả mong đợi | Kết quả thực tế | Trạng thái | Ghi chú |
|---|---|---|---|---|---|---|---|---|
| MOV-01 | Kiểm tra bố cục và các thành phần trang Quản lý phim | A. GUI Testing | Truy cập trang danh sách phim Admin (`/admin/movies`) | 1. Đăng nhập trang quản trị bằng tài khoản quản lý.<br>2. Nhấp chọn menu "Quản lý phim" trên thanh điều hướng.<br>3. Quan sát tổng thể giao diện trang. | Giao diện hiển thị đúng chuẩn thiết kế: Tiêu đề "Quản lý phim", mô tả phụ, nút "Thêm phim mới" nổi bật, thanh tìm kiếm kèm icon kính lúp, lưới danh sách phim và thanh phân trang. | Giao diện hiển thị đầy đủ, rõ ràng các thành phần, bố cục cân đối, chuẩn responsive trên độ phân giải màn hình. | **PASS** | |
| MOV-02 | Kiểm tra hiển thị thông tin chi tiết trên Card phim | A. GUI Testing | Khối thẻ (Card) của phim "Avengers: Endgame" trên trang danh sách | 1. Đăng nhập tài khoản quản lý.<br>2. Vào màn hình "Quản lý phim".<br>3. Quan sát các thông tin hiển thị trên Card phim "Avengers: Endgame". | Card hiển thị đầy đủ: Ảnh poster sắc nét, Tên phim in đậm, Badge trạng thái có màu sắc phân biệt, Đạo diễn, Thể loại, Thời lượng định dạng "Xh Ym", Ngày phát hành. | Hiển thị đầy đủ poster, tên phim in đậm "Avengers: Endgame", badge "SHOWING" nền xanh lá, đạo diễn Anthony Russo, thể loại Action, thời lượng "3h 1m", ngày phát hành chuẩn `dd/mm/yyyy`. | **PASS** | |
| MOV-03 | Kiểm tra giao diện màn hình form "Thêm phim mới" | A. GUI Testing | Truy cập đường dẫn tạo phim (`/admin/movies/new`) | 1. Tại trang Quản lý phim, bấm nút "Thêm phim mới".<br>2. Quan sát các trường thông tin trên form nhập liệu. | Màn hình hiển thị tiêu đề "Thêm phim mới", nút "← Quay lại danh sách phim", đầy đủ các trường: Tên phim, Đạo diễn, Diễn viên, Thể loại (danh sách checkbox), Thời lượng, Ngày phát hành, Trạng thái (dropdown), URL Poster, URL Trailer, Mô tả; các nút "Hủy" và "Tạo phim mới". Các trường bắt buộc có đánh dấu sao đỏ (*). | Form hiển thị chuẩn xác, trực quan, các trường bắt buộc (Tên phim, Thời lượng) có gắn dấu sao (*), các placeholder hướng dẫn rõ ràng. | **PASS** | |
| MOV-04 | Kiểm tra giao diện màn hình xem chi tiết phim | A. GUI Testing | Nhấp vào một phim bất kỳ trên danh sách | 1. Tại trang Quản lý phim, nhấp chuột vào card phim bất kỳ.<br>2. Quan sát giao diện trang Chi tiết phim (`/admin/movies/:id`). | Hiển thị nút "← Back to Movies", nút "Edit Movie", "Delete Movie". Phần nội dung hiển thị ảnh poster lớn, tiêu đề phim, nhãn trạng thái, thông tin chi tiết, nhóm nút chuyển đổi trạng thái nhanh, phần mô tả và trailer. | Toàn bộ các khối thông tin, hình ảnh và nút chức năng đều hiển thị đầy đủ, đúng vị trí và bố cục. | **PASS** | |
| MOV-05 | Kiểm tra màu sắc hiển thị Badge trạng thái phim UPCOMING trên trang chi tiết | A. GUI Testing | Phim có trạng thái "UPCOMING" (Sắp chiếu) | 1. Chọn mở một bộ phim đang ở trạng thái UPCOMING (Sắp chiếu).<br>2. Quan sát màu sắc của badge trạng thái dưới tên phim. | Badge trạng thái phải hiển thị nền màu cam/vàng đặc trưng (`bg-orange-500`) đại diện cho trạng thái Sắp chiếu (UPCOMING). | Badge hiển thị chữ "UPCOMING" nhưng có màu xám mặc định (`bg-gray-500`) thay vì màu cam. | **FAIL** | **BUG (Minor)**: Trong file `MovieDetailPage.jsx`, hàm `getStatusColorClasses` viết sai chính tả thành `'UPCOMMING'` (thừa chữ M), khiến trạng thái `'UPCOMING'` rơi vào nhánh switch default màu xám. |
| MOV-06 | Kiểm tra giao diện form "Chỉnh sửa phim" khi nạp dữ liệu cũ | A. GUI Testing | Trang chỉnh sửa của phim đã tồn tại (`/admin/movies/:id/edit`) | 1. Mở trang chi tiết một phim.<br>2. Bấm nút "Edit Movie".<br>3. Kiểm tra các giá trị được điền sẵn (pre-fill) trong các ô nhập liệu. | Tiêu đề trang hiển thị "Chỉnh sửa phim", nút bấm đổi thành "Cập nhật phim". Toàn bộ dữ liệu cũ (tên, đạo diễn, diễn viên, checkbox thể loại, thời lượng, mô tả...) được điền chính xác vào từng ô tương ứng. | Toàn bộ dữ liệu của phim được nạp chính xác vào từng trường nhập liệu của form chỉnh sửa. | **PASS** | |
| MOV-07 | Kiểm tra giao diện hiển thị khi không tìm thấy kết quả tìm kiếm | A. GUI Testing | Nhập từ khóa không có thực vào ô tìm kiếm | 1. Tại ô tìm kiếm phim, nhập chuỗi ký tự ngẫu nhiên "###$$$@@@".<br>2. Bấm nút "Tìm kiếm".<br>3. Quan sát khu vực hiển thị danh sách. | Khu vực danh sách không hiển thị card phim nào, xuất hiện icon rạp phim xám, dòng thông báo "Không tìm thấy phim nào" và mô tả chi tiết từ khóa không khớp. | Hiển thị chính xác icon rạp phim lớn, tiêu đề "Không tìm thấy phim nào" và dòng giải thích "Không có phim nào khớp với từ khóa "###$$$@@@"". | **PASS** | |
| MOV-08 | Thêm phim – Bỏ trống trường bắt buộc "Tên phim" | B. Field Validation | Form thêm phim mới:<br>- Tên phim: Để trống<br>- Thời lượng: 120 | 1. Mở form "Thêm phim mới".<br>2. Để trống ô "Tên phim".<br>3. Nhập Thời lượng = 120.<br>4. Bấm nút "Tạo phim mới". | Hệ thống chặn thao tác gửi form, hiển thị khung thông báo lỗi màu đỏ kèm biểu tượng cảnh báo: "Tên phim và thời lượng là bắt buộc". Không tạo phim. | Form không được gửi đi, xuất hiện thông báo lỗi đỏ: "Tên phim và thời lượng là bắt buộc", người dùng vẫn ở lại màn hình nhập liệu. | **PASS** | Validation phía client hoạt động chuẩn xác. |
| MOV-09 | Thêm phim – Bỏ trống trường bắt buộc "Thời lượng" | B. Field Validation | Form thêm phim mới:<br>- Tên phim: "Chiến binh vũ trụ"<br>- Thời lượng: Để trống | 1. Mở form "Thêm phim mới".<br>2. Nhập Tên phim = "Chiến binh vũ trụ".<br>3. Xóa trắng ô "Thời lượng".<br>4. Bấm nút "Tạo phim mới". | Hệ thống không gửi dữ liệu, hiển thị cảnh báo yêu cầu nhập thời lượng bắt buộc. | Trình duyệt và form hiển thị thông báo "Tên phim và thời lượng là bắt buộc", ngăn chặn việc tạo mới. | **PASS** | |
| MOV-10 | Thêm phim – Trường "Tên phim" chỉ chứa toàn dấu cách | B. Field Validation | Form thêm phim mới:<br>- Tên phim: "   " (3 khoảng trắng)<br>- Thời lượng: 110 | 1. Mở form "Thêm phim mới".<br>2. Nhập 3 ký tự khoảng trắng vào ô "Tên phim".<br>3. Nhập Thời lượng = 110.<br>4. Bấm nút "Tạo phim mới". | Hệ thống nhận diện chuỗi không hợp lệ sau khi cắt khoảng trắng (trim) và báo lỗi tên phim bắt buộc. | Giao diện chặn lại và hiển thị thông báo lỗi: "Tên phim và thời lượng là bắt buộc". Không có phim trống nào được tạo ra. | **PASS** | Code client có hàm `!formData.title?.trim()` xử lý tốt. |
| MOV-11 | Thêm phim – Tên phim đạt biên tối đa 255 ký tự hợp lệ | B. Field Validation | Form thêm phim mới:<br>- Tên phim: Chuỗi 255 ký tự chữ cái<br>- Thời lượng: 100 | 1. Mở form "Thêm phim mới".<br>2. Nhập vào ô "Tên phim" chuỗi đúng 255 ký tự.<br>3. Nhập Thời lượng = 100.<br>4. Bấm "Tạo phim mới". | Hệ thống chấp nhận dữ liệu, tạo phim thành công và quay về trang danh sách phim. | Tạo phim thành công, hệ thống chuyển hướng về danh sách phim, card phim hiển thị tên rút gọn bằng dấu ba chấm (do CSS clamp). | **PASS** | Kiểm thử giá trị biên hợp lệ. |
| MOV-12 | Thêm phim – Tên phim vượt quá 255 ký tự (256 ký tự) | B. Field Validation | Form thêm phim mới:<br>- Tên phim: Chuỗi 256 ký tự<br>- Thời lượng: 100 | 1. Mở form "Thêm phim mới".<br>2. Nhập chuỗi 256 ký tự vào ô Tên phim.<br>3. Bấm "Tạo phim mới". | Hệ thống từ chối lưu dữ liệu, hiển thị thông báo lỗi yêu cầu độ dài tên phim không vượt quá 255 ký tự. | Form hiển thị thông báo lỗi trên khung đỏ "Không thể tạo phim", phim không được lưu vào hệ thống. | **PASS** | |
| MOV-13 | Thêm phim – Thời lượng bằng 0 | B. Field Validation | Form thêm phim mới:<br>- Tên phim: Phim Thử Nghiệm 01<br>- Thời lượng: 0 | 1. Mở form "Thêm phim mới".<br>2. Nhập Tên phim = "Phim Thử Nghiệm 01".<br>3. Nhập Thời lượng = 0.<br>4. Bấm "Tạo phim mới". | Trình duyệt kích hoạt thuộc tính kiểm tra HTML5 (`min="1"`), chặn gửi form và hiện tooltip cảnh báo giá trị phải lớn hơn hoặc bằng 1. | Trình duyệt hiển thị bóng cảnh báo "Value must be greater than or equal to 1" ngay tại ô Thời lượng, không cho submit. | **PASS** | Thuộc tính `min="1"` ngăn chặn hợp lệ. |
| MOV-14 | Thêm phim – Thời lượng là số âm | B. Field Validation | Form thêm phim mới:<br>- Tên phim: Phim Thử Nghiệm 02<br>- Thời lượng: -90 | 1. Mở form "Thêm phim mới".<br>2. Nhập Tên phim hợp lệ, Thời lượng = -90.<br>3. Bấm "Tạo phim mới". | Trình duyệt chặn gửi form do vi phạm điều kiện số dương tối thiểu (`min="1"`). | Trình duyệt ngăn chặn submit form và thông báo yêu cầu giá trị tối thiểu là 1. | **PASS** | |
| MOV-15 | Thêm phim – Nhập ký tự chữ vào ô Thời lượng | B. Field Validation | Form thêm phim mới:<br>- Thời lượng: gõ ký tự "abc" | 1. Mở form "Thêm phim mới".<br>2. Đặt con trỏ chuột vào ô "Thời lượng (phút)".<br>3. Cố gắng gõ các ký tự chữ cái "abc". | Trường nhập liệu kiểu số (`type="number"`) không nhận các ký tự chữ cái, giá trị trong ô vẫn giữ nguyên hoặc rỗng. | Ô nhập liệu tự động ngăn không cho nhập chữ cái, giá trị không thay đổi. | **PASS** | Thuộc tính `type="number"` hoạt động đúng. |
| MOV-16 | Chỉnh sửa phim – Xóa rỗng trường Tên phim khi cập nhật | B. Field Validation | Form chỉnh sửa phim:<br>- Tên phim: Xóa trắng toàn bộ | 1. Mở form chỉnh sửa của một phim.<br>2. Xóa sạch tên phim trong ô nhập liệu.<br>3. Bấm nút "Cập nhật phim". | Hệ thống chặn cập nhật, hiển thị thông báo lỗi "Tên phim và thời lượng là bắt buộc". | Hệ thống không gửi dữ liệu cập nhật, hiển thị thông báo lỗi đỏ "Tên phim và thời lượng là bắt buộc". | **PASS** | |
| MOV-17 | Chỉnh sửa phim – Sửa Thời lượng thành số 0 | B. Field Validation | Form chỉnh sửa phim:<br>- Thời lượng: Sửa thành 0 | 1. Mở form chỉnh sửa của một phim.<br>2. Sửa ô Thời lượng thành 0.<br>3. Bấm "Cập nhật phim". | Trình duyệt chặn lại và báo giá trị thời lượng không hợp lệ. | Trình duyệt hiển thị thông báo cảnh báo giá trị phải từ 1 trở lên, ngăn chặn cập nhật. | **PASS** | |
| MOV-18 | Xem danh sách phim và điều hướng phân trang | C. Functional & Display Testing | Danh sách phim có nhiều trang (mỗi trang hiển thị 12 card phim) | 1. Tại trang Quản lý phim, cuộn xuống chân trang.<br>2. Quan sát dòng "Trang 1 / X • Tổng Y phim".<br>3. Bấm nút "Sau →" để sang trang kế tiếp.<br>4. Bấm nút "← Trước" để quay lại. | Hệ thống tải và hiển thị đúng 12 phim tiếp theo của trang 2, cuộn màn hình lên đầu trang. Khi bấm nút quay lại thì hiển thị đúng danh sách trang 1. | Chức năng phân trang chuyển trang mượt mà, số trang hiển thị chính xác, danh sách phim cập nhật tương ứng theo từng trang. | **PASS** | |
| MOV-19 | Tìm kiếm phim theo tên chính xác | C. Functional & Display Testing | Ô tìm kiếm: Nhập từ khóa "Spider-Man" | 1. Tại ô "Tìm kiếm phim theo tên...", nhập "Spider-Man".<br>2. Bấm nút "Tìm kiếm" (hoặc nhấn phím Enter). | Danh sách lọc và chỉ hiển thị các bộ phim có tên chứa cụm từ "Spider-Man". Hiển thị dòng thông báo: "Kết quả tìm kiếm cho "Spider-Man" (X phim được tìm thấy)". | Hiển thị chính xác phim "Spider-Man: No Way Home", kèm dòng thông báo "Kết quả tìm kiếm cho "Spider-Man" (1 phim được tìm thấy)". | **PASS** | |
| MOV-20 | Tìm kiếm phim không phân biệt chữ hoa, chữ thường | C. Functional & Display Testing | Ô tìm kiếm: Nhập chữ thường "batman" | 1. Nhập từ khóa chữ thường "batman" vào ô tìm kiếm.<br>2. Bấm nút "Tìm kiếm". | Hệ thống tìm thấy phim "The Batman" không phân biệt chữ hoa hay chữ thường. | Kết quả hiển thị đúng card phim "The Batman" (trạng thái SHOWING). Tính năng tìm kiếm không phân biệt hoa thường hoạt động chính xác. | **PASS** | |
| MOV-21 | Tìm kiếm phim với từ khóa không tồn tại | C. Functional & Display Testing | Ô tìm kiếm: Nhập "PhimKhongCoThuc999" | 1. Nhập "PhimKhongCoThuc999" vào ô tìm kiếm.<br>2. Bấm nút "Tìm kiếm". | Danh sách trống, hiển thị màn hình thông báo không tìm thấy kết quả. | Màn hình hiển thị "Không tìm thấy phim nào" kèm thông điệp "Không có phim nào khớp với từ khóa...". | **PASS** | |
| MOV-22 | Xóa từ khóa tìm kiếm nhanh bằng nút 'X' | C. Functional & Display Testing | Ô tìm kiếm đang có sẵn nội dung "Avengers" | 1. Tìm kiếm phim với từ khóa "Avengers".<br>2. Bấm vào biểu tượng dấu 'X' nằm bên phải trong ô tìm kiếm. | Ô tìm kiếm được xóa trắng ngay lập tức; hệ thống tự động tải lại toàn bộ danh sách phim ban đầu mà không cần bấm nút tìm kiếm. | Ô nhập liệu sạch nội dung, danh sách tự động tải lại đầy đủ tất cả các phim trong hệ thống. | **PASS** | |
| MOV-23 | Kiểm tra định dạng hiển thị thời lượng phim (Xh Ym) | C. Functional & Display Testing | Phim có thời lượng 150 phút và phim 60 phút | 1. Quan sát card phim có thời lượng 150 phút.<br>2. Quan sát card phim có thời lượng 60 phút. | Thời lượng 150 phút hiển thị thành "2h 30m"; thời lượng 60 phút hiển thị thành "1h 0m". | Thời lượng hiển thị chính xác theo định dạng giờ và phút chuẩn "2h 30m" và "1h 0m". | **PASS** | Hàm `formatDuration` tính toán chính xác. |
| MOV-24 | Kiểm tra định dạng hiển thị ngày phát hành | C. Functional & Display Testing | Phim có ngày phát hành trong DB dạng ISO (2026-12-25) | 1. Mở card phim và trang chi tiết của phim.<br>2. Kiểm tra dòng "Release" trên Card và "Ngày ra mắt" trên trang chi tiết. | Trên Card phim hiển thị dạng rút gọn `25/12/2026`; Trên trang chi tiết hiển thị đầy đủ "25 tháng 12, 2026" (chuẩn tiếng Việt). | Cả Card phim và trang chi tiết đều format ngày tháng năm theo chuẩn tiếng Việt rất rõ ràng, thân thiện. | **PASS** | |
| MOV-25 | Tạo phim mới thành công với đầy đủ tất cả thông tin | D. Business Flow Testing | Dữ liệu form Thêm phim mới:<br>- Tên: "Hành Tinh Cát 3"<br>- Đạo diễn: "Denis Villeneuve"<br>- Diễn viên: "Timothee Chalamet, Zendaya"<br>- Thể loại: Tích chọn Sci-Fi, Adventure<br>- Thời lượng: 165<br>- Ngày phát hành: 15/11/2026<br>- Trạng thái: Sắp chiếu<br>- URL Poster: `https://example.com/poster.jpg`<br>- URL Trailer: `https://youtube.com/watch?v=demo`<br>- Mô tả: "Phần tiếp theo của siêu phẩm..." | 1. Bấm nút "Thêm phim mới".<br>2. Điền đầy đủ và chính xác tất cả các trường dữ liệu như mô tả.<br>3. Bấm nút "Tạo phim mới".<br>4. Kiểm tra sự xuất hiện của phim trong danh sách. | Nút chuyển trạng thái "Đang lưu...", sau đó hệ thống tự động điều hướng về trang Quản lý phim. Phim mới xuất hiện ở danh sách với badge "UPCOMING". | Phim "Hành Tinh Cát 3" được tạo thành công, xuất hiện ở danh sách phim với đầy đủ thông tin và poster. | **PASS** | Luồng nghiệp vụ thêm phim đầy đủ hoạt động tốt. |
| MOV-26 | Tạo phim mới chỉ với các trường bắt buộc tối thiểu | D. Business Flow Testing | Dữ liệu form:<br>- Tên phim: "Phim Ngắn Độc Lập"<br>- Thời lượng: 45<br>- Các trường khác: Để trống | 1. Bấm nút "Thêm phim mới".<br>2. Chỉ nhập Tên phim và Thời lượng.<br>3. Bấm "Tạo phim mới". | Tạo phim thành công, trạng thái mặc định được gán là "UPCOMING", các trường khác hiển thị "Unknown" hoặc "TBA". | Phim được tạo thành công, hiển thị ở danh sách với tên "Phim Ngắn Độc Lập", thời lượng "0h 45m", badge UPCOMING. | **PASS** | |
| MOV-27 | Chỉnh sửa thông tin phim thành công | D. Business Flow Testing | Phim "Phim Ngắn Độc Lập":<br>- Đổi tên thành: "Phim Ngắn Độc Lập - Đã Sửa"<br>- Sửa thời lượng: 50<br>- Nhập thêm mô tả mới | 1. Bấm vào xem chi tiết phim vừa tạo.<br>2. Bấm nút "Edit Movie".<br>3. Cập nhật các trường thông tin.<br>4. Bấm "Cập nhật phim".<br>5. Mở lại chi tiết phim kiểm tra. | Hệ thống thông báo cập nhật thành công, điều hướng về danh sách, thông tin mới được lưu chính xác. | Dữ liệu được cập nhật thành công, tên mới và thời lượng 50 phút hiển thị chính xác trên giao diện. | **PASS** | |
| MOV-28 | Chuyển đổi trạng thái phim hợp lệ: UPCOMING → SHOWING | D. Business Flow Testing | Phim đang ở trạng thái "UPCOMING" (Sắp chiếu) | 1. Truy cập trang chi tiết bộ phim đang ở trạng thái "UPCOMING".<br>2. Tại mục "Thay đổi trạng thái", bấm nút "showing". | Trạng thái phim chuyển sang "SHOWING". Badge đổi màu xanh lá. Nút "showing" bị vô hiệu hóa (disabled). | Trạng thái chuyển thành công sang SHOWING, badge chuyển màu xanh, nút "showing" trở thành trạng thái đang chọn. | **PASS** | Chuyển đổi trạng thái hợp lệ theo vòng đời phim. |
| MOV-29 | Chuyển đổi trạng thái phim không hợp lệ: SHOWING → UPCOMING | D. Business Flow Testing | Phim đang ở trạng thái "SHOWING" (Đang chiếu) | 1. Mở trang chi tiết bộ phim đang ở trạng thái "SHOWING".<br>2. Tại mục "Thay đổi trạng thái", bấm nút "upcoming". | Hệ thống từ chối chuyển trạng thái, hiển thị thông báo lỗi nghiệp vụ rõ ràng, thân thiện bằng tiếng Việt (VD: "Không thể chuyển trạng thái từ Đang chiếu về Sắp chiếu"). Trạng thái phim giữ nguyên là SHOWING. | Xuất hiện hộp thoại alert của trình duyệt: "Failed to update movie status: invalid status transition" bằng tiếng Anh dạng mã lỗi server. Trạng thái phim vẫn giữ nguyên SHOWING. | **FAIL** | **BUG (Minor)**: Chức năng chặn đúng quy tắc nghiệp vụ, nhưng giao diện hiển thị thông báo qua hàm `alert()` bằng tiếng Anh kỹ thuật, không thân thiện với người dùng quản trị. |
| MOV-30 | Xóa phim thành công với xác nhận hộp thoại Confirm | D. Business Flow Testing | Phim thử nghiệm vừa tạo ở các bước trên | 1. Mở chi tiết bộ phim thử nghiệm cần xóa.<br>2. Bấm nút "Delete Movie" màu đỏ.<br>3. Khi hộp thoại xác nhận hiện lên "Are you sure you want to delete this movie?", chọn OK.<br>4. Kiểm tra danh sách phim. | Phim bị xóa hoàn toàn khỏi hệ thống, giao diện điều hướng về trang Quản lý phim và phim không còn xuất hiện trong danh sách. | Phim bị xóa thành công, trang tự động chuyển về `/admin/movies`, phim đã biến mất khỏi danh sách hiển thị. | **PASS** | |

### 1.3. Kết luận phân hệ Quản lý phim

- **Tổng số test case:** 30 ca.
- **Số ca Đạt (PASS):** 28 ca (93,3%).
- **Số ca Không đạt (FAIL):** 02 ca (6,7%) gồm:
  - **MOV-05 (Minor):** Lỗi chính tả trong mã nguồn (`UPCOMMING`) khiến badge trạng thái phim Sắp chiếu hiển thị màu xám thay vì màu cam trên trang chi tiết.
  - **MOV-29 (Minor):** Thông báo chặn chuyển trạng thái không hợp lệ hiển thị bằng thông điệp tiếng Anh thô qua hộp thoại alert của trình duyệt.

---

## 2. QUẢN LÝ GHẾ (Seat Management)

### 2.1. Thống kê kết quả

| Trạng thái | Số lượng | Tỷ lệ | Định nghĩa |
|---|---|---|---|
| **TESTED** | 30 | 100% | Tổng số test case hệ thống đã thực hiện |
| **PASS** | 28 | 93,3% | Test case cho kết quả thực tế đúng như kỳ vọng |
| **FAIL** | 2 | 6,7% | Test case phát hiện sai lệch hoặc lỗi giao diện / nghiệp vụ |
| **N/A** | 0 | 0% | Test case không áp dụng được |

### 2.2. Bảng chi tiết Test Case Quản lý ghế (30 ca)

| ID | Tên chức năng kiểm thử | Loại test case | Dữ liệu đầu vào | Các bước thực hiện | Kết quả mong đợi | Kết quả thực tế | Trạng thái | Ghi chú |
|---|---|---|---|---|---|---|---|---|
| SEAT-01 | Kiểm tra bố cục và các thành phần màn hình Quản lý Ghế ngồi | A. GUI Testing | Màn hình Quản lý Ghế (`/admin/seats`) | 1. Đăng nhập trang quản trị.<br>2. Chọn menu "Quản lý Ghế ngồi" trên thanh điều hướng.<br>3. Quan sát giao diện tổng thể. | Hiển thị tiêu đề "Quản lý Ghế ngồi", mô tả, bộ nút chuyển đổi chế độ hiển thị "Lưới" và "Bảng". Khi ở chế độ Bảng có thêm nút "Thêm ghế mới". | Giao diện hiển thị trực quan, các nút chuyển đổi chế độ xem nổi bật và thao tác nhạy bén. | **PASS** | |
| SEAT-02 | Kiểm tra hiển thị sơ đồ ghế ở chế độ Lưới (Grid View) | A. GUI Testing | Màn hình Quản lý Ghế ở chế độ "Lưới", chọn Phòng 1 | 1. Nhấp chọn tab "Lưới".<br>2. Chọn "Phòng 1" từ dropdown danh sách phòng.<br>3. Quan sát sơ đồ ghế hiển thị trên màn hình. | Sơ đồ ghế hiển thị trực quan theo hàng (A, B, C...) và cột (01, 02...), có thanh biểu tượng Màn hình chiếu (Screen), vị trí Cửa ra vào (Entrance), các màu ghế tương ứng với từng loại ghế. | Sơ đồ ghế mô phỏng rất đẹp mắt, phân bố đúng số hàng và cột của Phòng 1, hiển thị rõ màn hình và lối vào. | **PASS** | |
| SEAT-03 | Kiểm tra hiển thị danh sách ghế ở chế độ Bảng (Table View) | A. GUI Testing | Chuyển sang chế độ "Bảng" | 1. Nhấp chọn tab "Bảng".<br>2. Quan sát các cột trong bảng dữ liệu. | Bảng dữ liệu hiển thị đầy đủ 6 cột chuẩn: Phòng, Vị trí, Loại ghế, Trạng thái (dropdown), Ngày tạo, Thao tác (icon sửa, xóa). | Bảng hiển thị đầy đủ và rõ ràng các cột thông tin, dữ liệu căn chỉnh thẳng hàng dễ đọc. | **PASS** | |
| SEAT-04 | Kiểm tra hiển thị bảng chú thích màu sắc ghế (Legend) | A. GUI Testing | Chế độ xem Lưới tại Phòng 1 | 1. Mở sơ đồ ghế chế độ Lưới.<br>2. Quan sát khu vực chú thích ở góc trên bên phải. | Hiển thị 3 ô màu chú thích rõ ràng: Ô màu xanh lá đại diện cho ghế "Thường", ô màu vàng đại diện cho ghế "VIP", ô màu hồng đại diện cho ghế "Đôi". | Chú thích màu sắc hiển thị đúng 3 loại ghế với 3 màu xanh lá, vàng, hồng tương ứng. | **PASS** | |
| SEAT-05 | Kiểm tra giao diện form "Thêm ghế ngồi mới" | A. GUI Testing | Trang thêm ghế mới (`/admin/seats/new`) | 1. Tại chế độ Bảng, bấm nút "Thêm ghế mới".<br>2. Quan sát các trường trên form tạo ghế. | Màn hình hiển thị tiêu đề "Thêm ghế ngồi mới", nút quay lại, 4 trường bắt buộc có dấu sao (*): Phòng chiếu, Hàng ghế, Số ghế, Loại ghế; nút "Tạo mới" và nút "Hủy". | Form hiển thị chuẩn xác, bố cục lưới 2 cột cân đối, các trường bắt buộc đều có dấu sao (*) đỏ nổi bật. | **PASS** | |
| SEAT-06 | Kiểm tra giao diện form "Chỉnh sửa ghế ngồi" khi tải thông tin cũ | A. GUI Testing | Nhấp icon "Chỉnh sửa" tại dòng ghế A01 của Phòng 1 | 1. Tại danh sách bảng ghế, tìm ghế A01 Phòng 1.<br>2. Bấm vào icon cây bút "Chỉnh sửa".<br>3. Quan sát form chỉnh sửa ghế. | Tiêu đề hiển thị "Chỉnh sửa ghế ngồi", nút bấm là "Cập nhật". Các giá trị Phòng, Hàng ghế (A), Số ghế (01), Loại ghế (REGULAR) được điền sẵn; có thêm trường dropdown "Trạng thái". | Form nạp đầy đủ và chính xác toàn bộ dữ liệu hiện tại của ghế A01 lên giao diện. | **PASS** | |
| SEAT-07 | Kiểm tra hoạt động của khung "Xem trước" (Preview box) trên form ghế | A. GUI Testing | Form thêm ghế: Nhập Hàng ghế "B", Số ghế "05", chọn Phòng 1 | 1. Mở form thêm ghế mới.<br>2. Chọn Phòng 1, gõ Hàng ghế "B", gõ Số ghế "05".<br>3. Quan sát phía dưới form. | Xuất hiện khung màu xám "Xem trước: Ghế B05 - Phòng 1, Loại: Thường" ngay khi người dùng điền đủ 3 trường. | Khung Xem trước hiển thị tức thì và cập nhật nhãn chính xác theo thông tin người dùng đang gõ. | **PASS** | Tăng trải nghiệm người dùng (UX) rất tốt. |
| SEAT-08 | Thêm ghế – Bỏ trống trường bắt buộc "Số ghế" | B. Field Validation | Form thêm ghế mới:<br>- Phòng: Phòng 1<br>- Hàng ghế: "Z"<br>- Số ghế: Để trống | 1. Mở form Thêm ghế mới.<br>2. Chọn Phòng 1, nhập Hàng ghế "Z", để trống Số ghế.<br>3. Bấm nút "Tạo mới". | Hệ thống chặn không cho lưu, hiển thị khung thông báo lỗi màu đỏ: "Vui lòng điền đầy đủ thông tin bắt buộc". | Form bị chặn lại, xuất hiện thông báo lỗi đỏ "Vui lòng điền đầy đủ thông tin bắt buộc" trên đầu form. | **PASS** | Client validation chặn hợp lệ. |
| SEAT-09 | Thêm ghế – Bỏ trống trường bắt buộc "Hàng ghế" | B. Field Validation | Form thêm ghế mới:<br>- Phòng: Phòng 1<br>- Hàng ghế: Để trống<br>- Số ghế: "99" | 1. Chọn Phòng 1, bỏ trống Hàng ghế, nhập Số ghế "99".<br>2. Bấm nút "Tạo mới". | Hệ thống không tạo ghế, hiển thị thông báo lỗi bắt buộc nhập đủ thông tin. | Hiển thị thông báo lỗi màu đỏ yêu cầu điền đầy đủ thông tin bắt buộc. | **PASS** | |
| SEAT-10 | Thêm ghế – Chưa chọn phòng chiếu | B. Field Validation | Form thêm ghế mới:<br>- Phòng chiếu: Giữ nguyên "Chọn phòng"<br>- Hàng: "A", Số: "99" | 1. Mở form thêm ghế, không chọn phòng chiếu.<br>2. Nhập Hàng ghế "A", Số ghế "99".<br>3. Bấm "Tạo mới". | Hệ thống chặn lại và báo lỗi yêu cầu chọn phòng chiếu bắt buộc. | Form hiển thị thông báo lỗi yêu cầu điền đủ thông tin bắt buộc, không cho phép gửi. | **PASS** | |
| SEAT-11 | Thêm ghế – Nhập vị trí trùng lặp trong cùng một phòng chiếu | B. Field Validation | Form thêm ghế mới:<br>- Phòng: Phòng 1<br>- Hàng: "A", Số: "01" (vị trí A01 đã có sẵn) | 1. Chọn Phòng 1, nhập Hàng ghế "A", Số ghế "01".<br>2. Bấm nút "Tạo mới". | Hệ thống từ chối tạo ghế và hiển thị thông báo lỗi nghiệp vụ rõ ràng: "Vị trí ghế này đã tồn tại trong phòng". | Form hiển thị thông báo lỗi trong khung đỏ: "Vị trí ghế này đã tồn tại trong phòng". Ghế không bị trùng lặp. | **PASS** | Logic kiểm tra trùng vị trí hoạt động hoàn hảo. |
| SEAT-12 | Thêm ghế – Trường "Số ghế" chỉ nhập toàn dấu cách | B. Field Validation | Form thêm ghế mới:<br>- Phòng: Phòng 1<br>- Hàng ghế: "Y"<br>- Số ghế: "   " (3 khoảng trắng) | 1. Chọn Phòng 1, nhập Hàng ghế "Y".<br>2. Tại ô "Số ghế", gõ 3 dấu cách.<br>3. Bấm nút "Tạo mới". | Hệ thống phải cắt khoảng trắng (trim), nhận diện ô số ghế đang rỗng và từ chối tạo ghế. | Hệ thống chấp nhận và tạo ghế thành công; danh sách xuất hiện một dòng ghế có Hàng Y nhưng Số ghế bị rỗng (vô hình). | **FAIL** | **BUG (Minor)**: Hàm validation trên `SeatFormPage.jsx` chỉ kiểm tra `!formData.seat_number` mà không dùng `.trim()`, làm lọt khoảng trắng vào hệ thống. |
| SEAT-13 | Thêm ghế – Trường "Hàng ghế" chỉ nhập toàn dấu cách | B. Field Validation | Form thêm ghế mới:<br>- Phòng: Phòng 1<br>- Hàng ghế: "   " (3 khoảng trắng)<br>- Số ghế: "88" | 1. Chọn Phòng 1, gõ 3 dấu cách vào ô Hàng ghế.<br>2. Nhập Số ghế "88".<br>3. Bấm "Tạo mới". | Hệ thống nhận diện hàng ghế không hợp lệ và hiển thị thông báo lỗi. | Hệ thống tạo ghế thành công với hàng ghế rỗng, gây hiển thị khiếm khuyết trong bảng và sơ đồ. | **FAIL** | **BUG (Minor)**: Cùng nguyên nhân với SEAT-12, thiếu `.trim()` khi kiểm tra giá trị trường `row_number`. |
| SEAT-14 | Chỉnh sửa ghế – Đổi vị trí sang một vị trí đã có ghế khác | B. Field Validation | Sửa ghế Z99 Phòng 1 thành vị trí A01 (đã có sẵn) | 1. Mở form chỉnh sửa ghế Z99 Phòng 1.<br>2. Đổi Hàng ghế thành "A", Số ghế thành "01".<br>3. Bấm "Cập nhật". | Hệ thống chặn cập nhật và thông báo "Vị trí ghế này đã tồn tại trong phòng". | Hiển thị thông báo đỏ "Vị trí ghế này đã tồn tại trong phòng", vị trí ghế không bị ghi đè. | **PASS** | |
| SEAT-15 | Chỉnh sửa ghế – Xóa trắng trường Số ghế | B. Field Validation | Sửa ghế Z99: Xóa trắng ô Số ghế | 1. Mở form chỉnh sửa của ghế Z99.<br>2. Xóa sạch ô Số ghế.<br>3. Bấm "Cập nhật". | Hệ thống chặn cập nhật và báo lỗi bắt buộc điền thông tin. | Giao diện chặn lại và hiển thị cảnh báo yêu cầu điền đầy đủ thông tin bắt buộc. | **PASS** | |
| SEAT-16 | Xem danh sách ghế theo phòng chiếu trong chế độ Bảng | C. Functional & Display Testing | Chế độ Bảng: Nhập ID của Phòng 1 vào ô "Room ID" | 1. Chuyển sang tab "Bảng".<br>2. Điền mã Phòng 1 vào ô lọc "Room ID".<br>3. Quan sát danh sách ghế được lọc ra. | Bảng chỉ hiển thị danh sách các ghế thuộc Phòng 1, sắp xếp theo hàng ghế và số ghế tăng dần. | Danh sách chỉ liệt kê toàn bộ các ghế thuộc Phòng 1, thông tin từng ghế hiển thị chính xác. | **PASS** | |
| SEAT-17 | Lọc danh sách ghế theo loại ghế "VIP" | C. Functional & Display Testing | Chế độ Bảng: Chọn dropdown "Tất cả loại ghế" = "VIP" | 1. Tại chế độ Bảng, mở dropdown "Tất cả loại ghế".<br>2. Chọn giá trị "VIP". | Bảng dữ liệu tự động lọc lại, tất cả các dòng kết quả đều có cột Loại ghế hiển thị "VIP". | Bảng chỉ hiển thị các ghế VIP, không còn xuất hiện ghế Thường hay ghế Đôi. | **PASS** | Bộ lọc loại ghế hoạt động chuẩn xác. |
| SEAT-18 | Lọc danh sách ghế theo loại ghế "Đôi" (COUPLE) | C. Functional & Display Testing | Chế độ Bảng: Chọn dropdown Loại ghế = "Đôi" | 1. Mở dropdown Loại ghế, nhấp chọn "Đôi". | Bảng chỉ hiển thị các ghế Đôi (COUPLE) của rạp. | Tất cả các dòng ghế hiển thị đều có cột Loại ghế ghi là "Đôi". | **PASS** | |
| SEAT-19 | Lọc danh sách ghế theo trạng thái "Có sẵn" (AVAILABLE) | C. Functional & Display Testing | Chế độ Bảng: Chọn dropdown Trạng thái = "Có sẵn" | 1. Tại dropdown "Tất cả trạng thái", chọn "Có sẵn". | Bảng chỉ hiển thị các ghế đang ở trạng thái sẵn sàng đặt vé (AVAILABLE). | Toàn bộ danh sách ghế hiển thị dropdown trạng thái có nhãn xanh lá "Có sẵn". | **PASS** | |
| SEAT-20 | Lọc danh sách ghế theo trạng thái "Bảo trì" (MAINTENANCE) | C. Functional & Display Testing | Chế độ Bảng: Chọn dropdown Trạng thái = "Bảo trì" | 1. Tại dropdown Trạng thái, chọn "Bảo trì". | Bảng chỉ hiển thị các ghế đang trong quá trình bảo trì kỹ thuật. | Bảng hiển thị chính xác danh sách các ghế đang bảo trì (nhãn màu vàng). | **PASS** | |
| SEAT-21 | Lọc ghế kết hợp đa điều kiện (Phòng chiếu + Hàng ghế + Loại ghế) | C. Functional & Display Testing | Chế độ Bảng: Lọc Phòng 1, Hàng ghế "F", Loại ghế "VIP" | 1. Điền Room ID của Phòng 1.<br>2. Nhập Hàng ghế = "F".<br>3. Chọn Loại ghế = "VIP". | Bảng hiển thị danh sách các ghế thỏa mãn đồng thời cả 3 điều kiện: thuộc Phòng 1, ở hàng F và là loại VIP. | Bảng trả về đúng 8 ghế VIP từ F01 đến F08 của Phòng 1. Kết quả lọc đa tiêu chí hoàn toàn chính xác. | **PASS** | |
| SEAT-22 | Tìm kiếm ghế bằng ô tìm kiếm từ khóa | C. Functional & Display Testing | Ô "Tìm kiếm ghế...": Nhập từ khóa "F0" | 1. Tại thanh tìm kiếm ghế, gõ "F0". | Bảng tự động tìm kiếm và hiển thị các ghế có vị trí bắt đầu bằng F0 (F01 đến F09). | Danh sách lọc nhanh và hiển thị chính xác các ghế khớp với cụm từ tìm kiếm "F0". | **PASS** | |
| SEAT-23 | Kiểm tra chuyển đổi qua lại giữa chế độ xem Lưới và Bảng | C. Functional & Display Testing | Nhấp liên tục vào nút "Lưới" và "Bảng" | 1. Đang ở chế độ Lưới, bấm sang "Bảng".<br>2. Đang ở chế độ Bảng, bấm quay lại "Lưới". | Giao diện chuyển đổi tức thì, giữ nguyên phòng chiếu đang chọn, không bị giật lag hay trắng trang. | Thao tác chuyển đổi qua lại giữa 2 chế độ hiển thị mượt mà, đúng trạng thái kích hoạt của từng tab. | **PASS** | |
| SEAT-24 | Kiểm tra độ chính xác của số liệu thống kê dưới sơ đồ Lưới | C. Functional & Display Testing | Sơ đồ ghế Phòng 1 ở chế độ Lưới | 1. Mở sơ đồ Lưới của Phòng 1.<br>2. Đọc thanh thông số ở chân trang: "Tổng số ghế", "Còn trống", "Đã đặt", "Bảo trì".<br>3. Đối chiếu tổng số ghế với số lượng ghế thực tế trên sơ đồ. | Các số liệu thống kê phải hoàn toàn khớp với số lượng ghế hiển thị trên sơ đồ ma trận. | Số liệu thống kê: Tổng số ghế, Số ghế trống, Đã đặt và Bảo trì khớp hoàn toàn 100% với số ghế thực tế. | **PASS** | |
| SEAT-25 | Kiểm tra phân trang danh sách ghế trong chế độ Bảng | C. Functional & Display Testing | Bảng danh sách ghế (mỗi trang hiển thị 10 dòng) | 1. Tại chế độ Bảng, cuộn xuống thanh phân trang dưới cùng.<br>2. Bấm vào trang "2", rồi bấm nút "Sau".<br>3. Bấm nút "Trước". | Hệ thống chuyển trang chuẩn xác, mỗi trang hiển thị đúng 10 dòng ghế theo phân trang. | Phân trang hoạt động mượt mà, các nút Trước, Sau và số trang hoạt động chính xác. | **PASS** | |
| SEAT-26 | Tạo mới ghế loại Thường (REGULAR) thành công | D. Business Flow Testing | Form thêm ghế:<br>- Phòng: Phòng 1<br>- Hàng: "Z", Số: "99"<br>- Loại: Thường | 1. Bấm nút "Thêm ghế mới".<br>2. Nhập đầy đủ thông tin hợp lệ.<br>3. Bấm "Tạo mới".<br>4. Mở chế độ Bảng kiểm tra ghế Z99. | Ghế Z99 được tạo thành công, xuất hiện trong danh sách ghế Phòng 1 với loại Thường và trạng thái Có sẵn. | Tạo ghế thành công, ghế Z99 xuất hiện ngay trong bảng với đầy đủ thông tin chuẩn xác. | **PASS** | |
| SEAT-27 | Tạo mới ghế loại Đôi (COUPLE) thành công | D. Business Flow Testing | Form thêm ghế:<br>- Phòng: Phòng 1<br>- Hàng: "Z", Số: "98"<br>- Loại: Đôi | 1. Bấm "Thêm ghế mới".<br>2. Chọn Phòng 1, nhập Hàng "Z", Số "98", Loại "Đôi".<br>3. Bấm "Tạo mới". | Tạo thành công ghế đôi Z98 Phòng 1, hiển thị nhãn loại ghế là "Đôi". | Ghế Z98 hiển thị trong bảng với loại ghế "Đôi", trạng thái "Có sẵn". | **PASS** | |
| SEAT-28 | Cập nhật chuyển đổi loại ghế từ Thường sang VIP | D. Business Flow Testing | Ghế Z99 Phòng 1: Chuyển sang loại "VIP" | 1. Tại bảng ghế, tìm ghế Z99, bấm icon "Chỉnh sửa".<br>2. Tại dropdown Loại ghế, chọn sang "VIP".<br>3. Bấm "Cập nhật".<br>4. Kiểm tra lại thông tin trên bảng và sơ đồ lưới. | Ghế Z99 đổi sang loại VIP thành công; trên sơ đồ lưới đổi sang màu vàng. | Loại ghế chuyển thành VIP thành công, hiển thị chính xác cả trên dạng Bảng lẫn dạng Lưới. | **PASS** | |
| SEAT-29 | Thay đổi trực tiếp trạng thái ghế trên dropdown của Bảng | D. Business Flow Testing | Ghế Z99 đang "Có sẵn", chuyển sang "Bảo trì" | 1. Tại cột "Trạng thái" của dòng ghế Z99, bấm mở dropdown.<br>2. Nhấp chọn trạng thái "Bảo trì".<br>3. Tải lại trang (F5) để kiểm tra độ lưu trữ dữ liệu. | Trạng thái ghế đổi ngay sang nhãn vàng "Bảo trì". Sau khi F5 tải lại trang, trạng thái vẫn duy trì là "Bảo trì". | Trạng thái ghế cập nhật tức thì và lưu bền vững vào cơ sở dữ liệu, hiển thị chuẩn sau khi reload trang. | **PASS** | Tiện ích cập nhật trạng thái nhanh trên bảng hoạt động rất tốt. |
| SEAT-30 | Xóa ghế thành công kèm theo hộp thoại xác nhận | D. Business Flow Testing | Ghế Z99 thử nghiệm vừa tạo | 1. Tìm dòng ghế Z99 trong bảng danh sách.<br>2. Bấm vào biểu tượng thùng rác màu đỏ.<br>3. Hộp thoại xác nhận hiện lên "Bạn có chắc chắn muốn xóa ghế này?", chọn "OK". | Ghế Z99 bị xóa khỏi bảng ngay lập tức, không còn tồn tại trong cơ sở dữ liệu của phòng. | Ghế bị xóa thành công, biến mất khỏi bảng và sơ đồ lưới, không phát sinh lỗi. | **PASS** | Dọn dẹp dữ liệu kiểm thử hoàn tất. |

### 2.3. Kết luận phân hệ Quản lý ghế

- **Tổng số test case:** 30 ca.
- **Số ca Đạt (PASS):** 28 ca (93,3%).
- **Số ca Không đạt (FAIL):** 02 ca (6,7%) gồm:
  - **SEAT-12 (Minor):** Form thêm ghế không trim khoảng trắng trường "Số ghế", cho phép lưu số ghế rỗng.
  - **SEAT-13 (Minor):** Form thêm ghế không trim khoảng trắng trường "Hàng ghế", cho phép lưu hàng ghế rỗng.

---

## 3. KẾT LUẬN TỔNG THỂ HỆ THỐNG

### 3.1. Bảng tổng hợp số lượng và tỷ lệ kết quả kiểm thử

| Phân hệ chức năng | Số lượng ca kiểm thử (TESTED) | ĐẠT (PASS) | LỖI (FAIL) | Tỷ lệ ĐẠT (Pass Rate) |
|---|---|---|---|---|
| **1. Quản lý phim (Movie Management)** | 30 | 28 | 2 | **93,3%** |
| **2. Quản lý ghế (Seat Management)** | 30 | 28 | 2 | **93,3%** |
| **TỔNG CỘNG HỆ THỐNG** | **60** | **56** | **4** | **93,3%** |

---

### 3.2. Bảng phân bố kết quả theo 4 loại test case

| Nhóm loại test case | Quản lý phim | Quản lý ghế | Tổng số ca | Số ca PASS | Số ca FAIL | Tỷ lệ PASS |
|---|---|---|---|---|---|---|
| **A. GUI Testing** | 7 | 7 | **14** | 13 | 1 | 92,9% |
| **B. Field Validation** | 10 | 8 | **18** | 16 | 2 | 88,9% |
| **C. Functional & Display Testing** | 7 | 10 | **17** | 17 | 0 | 100% |
| **D. Business Flow Testing** | 6 | 5 | **11** | 10 | 1 | 90,9% |
| **TỔNG CỘNG** | **30** | **30** | **60** | **56** | **4** | **93,3%** |

---

### 3.3. Phân loại lỗi phát hiện theo mức độ nghiêm trọng

| Mức độ nghiêm trọng (Severity) | Số lượng | Mã test case liên quan | Diễn giải chi tiết |
|---|---|---|---|
| **Critical (Nghiêm trọng cao)** | **0** | Không có | Không có lỗi gây treo ứng dụng, sập dịch vụ hoặc gián đoạn luồng nghiệp vụ. |
| **Major (Nghiêm trọng vừa)** | **0** | Không có | Các chức năng nghiệp vụ cốt lõi đều hoạt động đúng đặc tả. |
| **Minor (Lỗi nhỏ / Giao diện & UX)** | **4** | MOV-05, MOV-29, SEAT-12, SEAT-13 | - **MOV-05:** Sai tên class trạng thái `'UPCOMMING'` làm sai màu badge trên trang chi tiết phim.<br>- **MOV-29:** Thông báo lỗi chuyển trạng thái không hợp lệ hiển thị alert tiếng Anh thô.<br>- **SEAT-12, SEAT-13:** Thiếu hàm `.trim()` ở trường Số ghế và Hàng ghế trên form tạo ghế. |

---

### 3.4. Đánh giá và đề xuất khắc phục

1. **Đánh giá chung:** 
   - Hệ thống quản lý rạp chiếu phim (HQ Cinema) ở cả hai phân hệ **Quản lý phim** và **Quản lý ghế** đạt chất lượng giao diện và nghiệp vụ rất tốt với tỷ lệ hoàn thành đạt **93,3%** (56/60 test case).
   - Trải nghiệm giao diện thân thiện, hiện đại, hỗ trợ cả hai góc nhìn trực quan (Sơ đồ ma trận phòng chiếu dạng Lưới và Bảng dữ liệu chi tiết).

2. **Kế hoạch khắc phục cho đội ngũ phát triển (Developer):**
   - **Sửa lỗi chính tả badge phim:** Trong file `FE/src/pages/admin/MovieDetailPage.jsx`, sửa case switch từ `'UPCOMMING'` thành `'UPCOMING'`.
   - **Cải thiện thông báo lỗi chuyển trạng thái:** Trong hàm `handleStatusChange` (`MovieDetailPage.jsx`), bắt lỗi và hiển thị thông báo tiếng Việt dạng Toast/Modal: *"Không thể chuyển trạng thái từ Đang chiếu về Sắp chiếu"*.
   - **Bổ sung hàm trim() cho form ghế:** Trong file `FE/src/pages/admin/SeatFormPage.jsx`, cập nhật điều kiện kiểm tra:
     ```javascript
     if (!formData.room_id || !formData.seat_number?.trim() || !formData.row_number?.trim()) {
       setError('Vui lòng điền đầy đủ thông tin bắt buộc')
       return
     }
     ```
   - **Tiến hành Regression Test (Kiểm thử hồi quy):** Thực hiện chạy lại 4 ca kiểm thử FAIL sau khi đội ngũ phát triển bàn giao bản sửa lỗi (hotfix) để nghiệm thu đạt 100%.
