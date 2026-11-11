# Bổ sung Mục 6 - Quality Assurance Activities (SQA Plan)

> **Ghi chú**: Nội dung dưới đây được bổ sung vào **Mục 6. Quality Assurance Activities**, cụ thể nằm giữa **6.2. Proposed testing strategy for the project** và **6.3. Proposed tests for the project**. Hai phần mới được thêm là **6.2.1** và **6.2.2**.

---

## 6.2.1. Xác định các tính năng cần kiểm thử và không cần kiểm thử

### A. Các tính năng được kiểm thử (System Test)

Nhóm lựa chọn thực hiện **System Test** cho các tính năng chính sau đây, bao gồm cả yêu cầu **chức năng (functional)** và **phi chức năng (non-functional)**:

| STT | Tính năng | Yêu cầu chức năng (Functional) | Yêu cầu phi chức năng (Non-functional) |
|-----|-----------|-------------------------------|----------------------------------------|
| 1 | **Login** | Đăng nhập bằng email & mật khẩu; xác thực JWT; xử lý sai mật khẩu, tài khoản không tồn tại; session timeout sau 30 phút | **Hiệu năng**: Thời gian phản hồi đăng nhập < 2 giây. **Bảo mật**: Mã hóa mật khẩu bcrypt, chống brute-force qua rate limiting (100 req/s), bảo vệ token JWT |
| 2 | **Register & Edit personal information** | Đăng ký tài khoản mới (email, password, tên, SĐT); validate email duy nhất, định dạng hợp lệ; chỉnh sửa thông tin cá nhân (tên, SĐT) | **Tính khả dụng (Usability)**: Form đăng ký rõ ràng, thông báo lỗi cụ thể cho từng trường. **Bảo mật**: Password tối thiểu 8 ký tự, mã hóa bcrypt |
| 3 | **Movie management** | Thêm/sửa/xóa phim (Admin); quản lý thể loại phim; tìm kiếm & lọc phim theo tên, thể loại; phân trang danh sách phim | **Hiệu năng**: Tìm kiếm phim phản hồi < 200ms nhờ Redis cache. **Tính khả dụng**: Hỗ trợ tìm kiếm full-text |
| 4 | **Seat management** | Quản lý ghế trong phòng chiếu (REGULAR, PREMIUM); thay đổi trạng thái ghế (AVAILABLE, RESERVED, OCCUPIED); quản lý phòng chiếu và sức chứa | **Tính toàn vẹn dữ liệu**: Đảm bảo không có xung đột trạng thái ghế khi nhiều người dùng truy cập đồng thời (row-level locking) |
| 5 | **Box Office Booking** | Nhân viên đặt vé tại quầy cho khách hàng; chọn suất chiếu, chọn ghế, xác nhận booking; xử lý thanh toán trực tiếp | **Hiệu năng**: Xử lý giao dịch đặt vé < 2 giây. **Độ tin cậy**: Cơ chế khóa ghế ngăn đặt trùng (concurrent booking prevention) |
| 6 | **Online Booking** | Khách hàng tự đặt vé trực tuyến; chọn suất chiếu → chọn ghế → thanh toán; booking hết hạn sau 15 phút nếu chưa thanh toán; tối đa 10 ghế/booking | **Hiệu năng**: Hỗ trợ 10.000 CCU với tỷ lệ lỗi giao dịch < 0.01%. **Tính khả dụng**: Hoàn thành quy trình đặt vé trong tối đa 4 bước (dưới 3 phút). **Bảo mật**: Mã hóa giao dịch, xác thực JWT bắt buộc |
| 7 | **Booking History** | Xem lịch sử đặt vé của người dùng; lọc theo trạng thái (PENDING, CONFIRMED, CANCELLED); phân trang; xem chi tiết từng booking | **Hiệu năng**: Truy vấn lịch sử phản hồi < 200ms. **Tính khả dụng**: Giao diện hiển thị rõ ràng trạng thái đơn hàng |
| 8 | **Ticket Export** | Xuất vé với mã QR code; hiển thị thông tin vé (mã vé, ghế, suất chiếu, giá); xác thực trạng thái vé (UNUSED, USED) | **Tính toàn vẹn**: Mã vé (ticket_code) là duy nhất, QR code chính xác. **Tính khả dụng**: Vé hiển thị đầy đủ thông tin cần thiết |
| 9 | **Staff & Permission Management** | Quản lý tài khoản nhân viên (Admin); gán/bỏ vai trò cho nhân viên; CRUD nhân viên | **Bảo mật**: Chỉ Admin mới có quyền quản lý nhân viên. **Tính khả dụng**: Giao diện quản lý nhân viên trực quan |
| 10 | **Revenue Analytics & Statistics** | Xem thống kê doanh thu theo ngày/tuần/tháng; thống kê theo phim, theo phòng chiếu; tỷ lệ chuyển đổi booking; phim phổ biến nhất | **Hiệu năng**: Truy vấn thống kê phức tạp được tối ưu bằng Redis cache và materialized views. **Bảo mật**: Chỉ Admin mới truy cập được |
| 11 | **Movie News Management** | Quản lý tài liệu chatbot (upload PDF); xử lý và lưu trữ knowledge base cho chatbot AI; quản lý nội dung hỗ trợ khách hàng | **Tính khả dụng**: Upload tài liệu PDF dễ dàng. **Hiệu năng**: Xử lý và tạo embedding vector cho document chunks |
| 12 | **Role-Based Access Control (RBAC)** | Phân quyền 3 cấp: Admin, Staff, Customer; kiểm soát truy cập tại API Gateway và service level; ánh xạ role–permission chi tiết (VD: `movies.create`, `users.delete`) | **Bảo mật**: 100% API được xác thực JWT và phân quyền RBAC. **Tính toàn vẹn**: Đảm bảo không có lỗ hổng leo thang quyền (privilege escalation). **Hiệu năng**: Xác thực/phân quyền không làm tăng đáng kể thời gian phản hồi API |

### B. Các tính năng / khía cạnh KHÔNG kiểm thử và lý do

| STT | Tính năng / Khía cạnh | Lý do không kiểm thử |
|-----|------------------------|----------------------|
| 1 | **Chatbot AI (RAG-based)** – Kiểm thử hiệu năng | Chatbot phụ thuộc hoàn toàn vào API bên ngoài (Google Gemini AI). Hiệu năng phản hồi của chatbot được quyết định bởi nhà cung cấp bên thứ ba, nhóm không có khả năng kiểm soát hoặc tối ưu. Nhóm chỉ kiểm thử chức năng tích hợp (integration test), không kiểm thử performance riêng cho chatbot. |
| 2 | **Payment Gateway bên ngoài (SePay, Ethereum)** – Kiểm thử hiệu năng | Đây là dịch vụ bên ngoài, hiệu năng xử lý thanh toán phụ thuộc vào hạ tầng của nhà cung cấp (blockchain Ethereum, SePay gateway). Nhóm chỉ kiểm thử khả năng tích hợp và xử lý kết quả trả về (timeout 5 giây, tỷ lệ kết nối > 99.5%), không kiểm thử tốc độ xử lý nội tại của các dịch vụ này. |
| 3 | **Hạ tầng phần cứng vật lý** | Nằm ngoài phạm vi dự án phần mềm. Hệ thống sử dụng Docker containers và triển khai trên môi trường ảo hóa, do đó kiểm thử phần cứng vật lý là không cần thiết. |
| 4 | **Movie News Management** – Kiểm thử hiệu năng (Performance test) | Tính năng quản lý tin tức/tài liệu chatbot có tần suất sử dụng thấp (chỉ Admin thỉnh thoảng upload tài liệu), không phải tính năng trọng yếu ảnh hưởng đến trải nghiệm đa người dùng. Do đó không cần kiểm thử hiệu năng dưới tải cao. |
| 5 | **Staff & Permission Management** – Kiểm thử hiệu năng (Performance test) | Quản lý nhân viên là tính năng quản trị nội bộ với số lượng người dùng rất ít (chỉ Admin), không có tình huống tải cao (high concurrency). Kiểm thử hiệu năng cho tính năng này là không cần thiết. |
| 6 | **Booking History** – Kiểm thử hiệu năng (Performance test) | Xem lịch sử đặt vé là thao tác chỉ đọc (read-only), truy vấn đơn giản với phân trang, không gây áp lực lớn lên hệ thống. Cơ chế caching Redis đã được thiết kế sẵn đảm bảo hiệu năng ổn định. |
| 7 | **Revenue Analytics & Statistics** – Kiểm thử hiệu năng (Performance test) | Tính năng thống kê chỉ được truy cập bởi Admin với tần suất thấp (không phải real-time). Hệ thống đã sử dụng materialized views và Redis cache để tối ưu. Không cần kiểm thử hiệu năng dưới tải cao. |
| 8 | **Hoạt động marketing bên ngoài** | Nằm ngoài phạm vi kỹ thuật của hệ thống phần mềm, không liên quan đến chức năng phần mềm cần kiểm thử. |

---

## 6.2.2. Kỹ thuật kiểm thử hộp đen (Black-box Testing Techniques) cho từng tính năng

Dựa trên 5 kỹ thuật kiểm thử hộp đen chính:
- **a. Equivalence Class Partitioning (ECP)** – Phân lớp tương đương
- **b. Boundary Value Analysis (BVA)** – Phân tích giá trị biên
- **c. Decision Tables** – Bảng quyết định
- **d. State Transition** – Chuyển đổi trạng thái
- **e. Pairwise Testing** – Kiểm thử cặp đôi

Dưới đây là bảng ánh xạ kỹ thuật kiểm thử cho từng tính năng và lý do lựa chọn:

### 1. Login

| Kỹ thuật | Áp dụng | Lý do |
|----------|---------|-------|
| **Equivalence Class Partitioning** | ✅ | Phân chia các lớp đầu vào: email hợp lệ/không hợp lệ/trống, password đúng/sai/trống, tài khoản tồn tại/không tồn tại. Mỗi lớp đại diện cho một nhóm dữ liệu có cùng hành vi kỳ vọng, giúp giảm số lượng test case mà vẫn đảm bảo coverage. |
| **Boundary Value Analysis** | ✅ | Kiểm tra giới hạn độ dài password (tối thiểu 8 ký tự): test các giá trị 7, 8, 9 ký tự. Kiểm tra giới hạn độ dài email. Kiểm tra session timeout đúng tại mốc 30 phút. |
| **Decision Tables** | ✅ | Tổ hợp điều kiện: email đúng/sai × password đúng/sai × tài khoản active/locked → xác định kết quả (login thành công, báo lỗi cụ thể). Bảng quyết định giúp liệt kê đầy đủ các tổ hợp điều kiện và hành động tương ứng. |
| **State Transition** | ✅ | Trạng thái phiên đăng nhập: Chưa đăng nhập → Đang đăng nhập → Đã đăng nhập → Session hết hạn → Chưa đăng nhập. Kiểm tra các chuyển đổi hợp lệ và không hợp lệ giữa các trạng thái. |

### 2. Register & Edit personal information

| Kỹ thuật | Áp dụng | Lý do |
|----------|---------|-------|
| **Equivalence Class Partitioning** | ✅ | Phân lớp các trường đầu vào: email (hợp lệ/không hợp lệ/đã tồn tại/trống), password (đủ mạnh/yếu/trống), tên (hợp lệ/trống/quá dài), SĐT (đúng định dạng/sai định dạng). Mỗi lớp tương đương cho kết quả xử lý giống nhau. |
| **Boundary Value Analysis** | ✅ | Kiểm tra biên của các trường: password đúng 8 ký tự (min), tên 100 ký tự (max), tên 101 ký tự (vượt max). Các giá trị biên thường là nơi dễ phát sinh lỗi nhất. |
| **Decision Tables** | ✅ | Kết hợp nhiều điều kiện validation đồng thời: email valid + password valid + tên valid → đăng ký thành công; email trùng + password valid → báo lỗi trùng email. Bảng quyết định giúp kiểm soát toàn bộ logic validation phức tạp. |

### 3. Movie management

| Kỹ thuật | Áp dụng | Lý do |
|----------|---------|-------|
| **Equivalence Class Partitioning** | ✅ | Phân lớp dữ liệu phim: title (hợp lệ/trống/quá dài), duration (số dương/âm/0/không phải số), release_date (ngày hợp lệ/quá khứ/tương lai/sai định dạng), genre_ids (tồn tại/không tồn tại/trống). |
| **Boundary Value Analysis** | ✅ | Kiểm tra biên phân trang: page = 0, 1, trang cuối, trang cuối + 1; per_page = 0, 1, 100 (max), 101. Kiểm tra biên duration phim (giá trị tối thiểu, tối đa). |
| **Decision Tables** | ✅ | Quyết định quyền thao tác: role (Admin/Staff/Customer) × action (Create/Read/Update/Delete) → cho phép/từ chối. Tổ hợp quyền và hành động tạo nên bảng quyết định rõ ràng cho kiểm thử phân quyền. |
| **Pairwise Testing** | ✅ | Tính năng tìm kiếm & lọc phim có nhiều tham số kết hợp: genre × search keyword × sort_by × order × page × per_page. Pairwise testing giúp giảm số lượng test case từ hàng trăm xuống còn vài chục mà vẫn phát hiện được lỗi do tương tác giữa các cặp tham số. |

### 4. Seat management

| Kỹ thuật | Áp dụng | Lý do |
|----------|---------|-------|
| **Equivalence Class Partitioning** | ✅ | Phân lớp loại ghế (REGULAR/PREMIUM), trạng thái ghế (AVAILABLE/RESERVED/OCCUPIED), loại phòng chiếu. Mỗi lớp có logic xử lý khác nhau cần được kiểm tra. |
| **Boundary Value Analysis** | ✅ | Kiểm tra biên số ghế trong phòng: phòng đầy (capacity = 0 ghế trống), phòng có đúng 1 ghế trống, phòng trống hoàn toàn. Kiểm tra giới hạn số hàng và số ghế mỗi hàng. |
| **State Transition** | ✅ | Ghế có vòng đời trạng thái rõ ràng: AVAILABLE → RESERVED → OCCUPIED (sau khi thanh toán) hoặc RESERVED → AVAILABLE (khi booking hết hạn/bị hủy). State Transition testing là kỹ thuật phù hợp nhất để kiểm tra tính đúng đắn của các chuyển đổi trạng thái ghế. |
| **Decision Tables** | ✅ | Quyết định cho phép/không cho phép đặt ghế dựa trên tổ hợp: trạng thái ghế hiện tại × loại ghế × suất chiếu còn hiệu lực → kết quả (đặt thành công/báo lỗi). |

### 5. Box Office Booking

| Kỹ thuật | Áp dụng | Lý do |
|----------|---------|-------|
| **Equivalence Class Partitioning** | ✅ | Phân lớp đầu vào: showtime_id (hợp lệ/không hợp lệ/đã chiếu), seat_ids (ghế trống/đã đặt/không tồn tại), số lượng ghế (1-10/0/>10). |
| **Boundary Value Analysis** | ✅ | Kiểm tra biên: đặt đúng 10 ghế (max), 11 ghế (vượt max), 1 ghế (min), 0 ghế. Kiểm tra biên thời gian: đặt vé cho suất chiếu sắp bắt đầu (sát giờ chiếu). |
| **Decision Tables** | ✅ | Tổ hợp điều kiện phức tạp: role nhân viên × trạng thái suất chiếu × trạng thái ghế × phương thức thanh toán → kết quả booking (thành công/từ chối/chờ thanh toán). Bảng quyết định giúp xử lý toàn diện các trường hợp nghiệp vụ. |
| **State Transition** | ✅ | Booking có các trạng thái: PENDING → CONFIRMED (sau thanh toán) hoặc PENDING → CANCELLED. Cần kiểm tra chuyển đổi trạng thái booking trong quy trình đặt vé tại quầy. |

### 6. Online Booking

| Kỹ thuật | Áp dụng | Lý do |
|----------|---------|-------|
| **Equivalence Class Partitioning** | ✅ | Phân lớp tương tự Box Office Booking, bổ sung thêm: user role (Customer/Guest không thể đặt), phương thức thanh toán online (CARD/ETHEREUM). |
| **Boundary Value Analysis** | ✅ | Kiểm tra biên tương tự Box Office, bổ sung: thời gian hết hạn booking 15 phút (test tại 14 phút 59 giây, 15 phút 00 giây, 15 phút 01 giây). Biên tổng giá trị đơn hàng. |
| **Decision Tables** | ✅ | Tổ hợp: trạng thái xác thực × trạng thái ghế × thời gian hết hạn × kết quả thanh toán → hành động hệ thống (xác nhận/hủy/trả ghế). Quy trình online booking phức tạp hơn với nhiều điều kiện kết hợp. |
| **State Transition** | ✅ | Quy trình online booking có chuỗi trạng thái phức tạp: Chọn ghế → PENDING → Thanh toán → CONFIRMED hoặc PENDING → Hết hạn (15 phút) → CANCELLED → Ghế trả về AVAILABLE. Kiểm tra kỹ tất cả đường chuyển trạng thái. |
| **Pairwise Testing** | ✅ | Nhiều tham số kết hợp: loại ghế (REGULAR/PREMIUM) × số lượng ghế (1-10) × phương thức thanh toán (CARD/ETHEREUM) × suất chiếu (sáng/chiều/tối). Pairwise giúp giảm tổ hợp test mà vẫn phát hiện lỗi tương tác cặp. |

### 7. Booking History

| Kỹ thuật | Áp dụng | Lý do |
|----------|---------|-------|
| **Equivalence Class Partitioning** | ✅ | Phân lớp bộ lọc: status (PENDING/CONFIRMED/CANCELLED/tất cả), phân trang (page hợp lệ/không hợp lệ), user có lịch sử booking/không có lịch sử. |
| **Boundary Value Analysis** | ✅ | Kiểm tra biên phân trang: page = 0, 1, trang cuối, trang cuối + 1. per_page = 0, 1, giá trị max. User có 0 booking, 1 booking, nhiều booking. |
| **Decision Tables** | ✅ | Tổ hợp bộ lọc: status filter × page × per_page → kết quả trả về (danh sách booking/rỗng/lỗi). Đảm bảo tất cả tổ hợp bộ lọc hoạt động chính xác. |

### 8. Ticket Export

| Kỹ thuật | Áp dụng | Lý do |
|----------|---------|-------|
| **Equivalence Class Partitioning** | ✅ | Phân lớp: booking có vé hợp lệ/không có vé, trạng thái vé (UNUSED/USED), quyền truy cập (chủ booking/người khác/Admin). |
| **State Transition** | ✅ | Vé có trạng thái: UNUSED → USED (khi quét QR tại rạp). Cần kiểm tra chuyển đổi trạng thái vé: không thể chuyển USED → UNUSED, không thể sử dụng vé USED lần thứ hai. |
| **Decision Tables** | ✅ | Tổ hợp: trạng thái booking (CONFIRMED/CANCELLED/PENDING) × trạng thái vé (UNUSED/USED) × quyền người dùng → cho phép xuất vé/từ chối. Bảng quyết định xác định rõ ràng khi nào vé được phép xuất. |

### 9. Staff & Permission Management

| Kỹ thuật | Áp dụng | Lý do |
|----------|---------|-------|
| **Equivalence Class Partitioning** | ✅ | Phân lớp: loại thao tác (thêm/sửa/xóa nhân viên), vai trò gán (Admin/Staff/Customer), dữ liệu nhân viên (hợp lệ/không hợp lệ). |
| **Decision Tables** | ✅ | Tổ hợp quyền phân quyền phức tạp: role hiện tại × role được gán × action → cho phép/từ chối. VD: Admin gán Staff → OK, Staff gán Admin → từ chối. Bảng quyết định là kỹ thuật tối ưu cho logic phân quyền. |
| **Boundary Value Analysis** | ✅ | Kiểm tra biên: thêm nhân viên khi đã đạt giới hạn hệ thống (nếu có), email nhân viên ở biên độ dài cho phép, tên ở độ dài max (100 ký tự). |

### 10. Revenue Analytics & Statistics

| Kỹ thuật | Áp dụng | Lý do |
|----------|---------|-------|
| **Equivalence Class Partitioning** | ✅ | Phân lớp: khoảng thời gian (ngày/tuần/tháng), loại thống kê (doanh thu/booking/phim phổ biến), quyền truy cập (Admin/không phải Admin). |
| **Boundary Value Analysis** | ✅ | Kiểm tra biên ngày: start_date = end_date (1 ngày), khoảng thời gian rất lớn (1 năm), start_date > end_date (không hợp lệ). Kiểm tra khi không có dữ liệu trong khoảng thời gian. |
| **Pairwise Testing** | ✅ | Nhiều tham số kết hợp: start_date × end_date × group_by (day/week/month/movie/room). Pairwise giảm số lượng test case cho các tổ hợp tham số thống kê mà vẫn đảm bảo phát hiện lỗi tương tác cặp. |

### 11. Movie News Management

| Kỹ thuật | Áp dụng | Lý do |
|----------|---------|-------|
| **Equivalence Class Partitioning** | ✅ | Phân lớp: loại file upload (PDF hợp lệ/file không phải PDF/file rỗng/file quá lớn), quyền upload (Admin/không phải Admin), nội dung tài liệu (có text/rỗng). |
| **Boundary Value Analysis** | ✅ | Kiểm tra biên kích thước file: file 0 byte, file ở giới hạn max cho phép, file vượt giới hạn. Số lượng document chunks tạo ra từ tài liệu ngắn nhất/dài nhất. |
| **Decision Tables** | ✅ | Tổ hợp: role × loại file × kích thước file → cho phép upload/từ chối. Xác định rõ ràng các trường hợp hệ thống chấp nhận hoặc từ chối upload. |

### 12. Role-Based Access Control (RBAC)

| Kỹ thuật | Áp dụng | Lý do |
|----------|---------|-------|
| **Decision Tables** | ✅ | Đây là kỹ thuật **quan trọng nhất** cho RBAC. Tổ hợp: role (Admin/Staff/Customer) × resource (movies/users/bookings/analytics/...) × action (create/read/update/delete) → allow/deny. Bảng quyết định giúp kiểm tra toàn diện ma trận phân quyền, đảm bảo không bỏ sót trường hợp nào. |
| **Equivalence Class Partitioning** | ✅ | Phân lớp role: Admin (toàn quyền), Staff (quyền hạn chế), Customer (quyền cơ bản), Chưa xác thực (chỉ public routes). Mỗi role đại diện cho một lớp tương đương với tập quyền khác nhau. |
| **State Transition** | ✅ | Chuyển đổi trạng thái quyền: User chưa đăng nhập → Đăng nhập (Customer) → Được nâng quyền (Staff/Admin) → Bị thu hồi quyền. Kiểm tra hệ thống phản ứng đúng khi vai trò người dùng thay đổi trong runtime. |

---

### Bảng tổng hợp kỹ thuật kiểm thử hộp đen cho từng tính năng

| Tính năng | ECP | BVA | Decision Tables | State Transition | Pairwise |
|-----------|:---:|:---:|:---------------:|:----------------:|:--------:|
| Login | ✅ | ✅ | ✅ | ✅ | |
| Register & Edit personal info | ✅ | ✅ | ✅ | | |
| Movie management | ✅ | ✅ | ✅ | | ✅ |
| Seat management | ✅ | ✅ | ✅ | ✅ | |
| Box Office Booking | ✅ | ✅ | ✅ | ✅ | |
| Online Booking | ✅ | ✅ | ✅ | ✅ | ✅ |
| Booking History | ✅ | ✅ | ✅ | | |
| Ticket Export | ✅ | | ✅ | ✅ | |
| Staff & Permission Management | ✅ | ✅ | ✅ | | |
| Revenue Analytics & Statistics | ✅ | ✅ | | | ✅ |
| Movie News Management | ✅ | ✅ | ✅ | | |
| Role-Based Access Control | ✅ | | ✅ | ✅ | |

**Giải thích nguyên tắc lựa chọn kỹ thuật:**

- **Equivalence Class Partitioning (ECP)**: Được áp dụng cho **tất cả** các tính năng vì mọi tính năng đều có các trường đầu vào cần phân lớp thành nhóm hợp lệ/không hợp lệ. Đây là kỹ thuật nền tảng giúp thiết kế test case hiệu quả.

- **Boundary Value Analysis (BVA)**: Áp dụng cho các tính năng có **dữ liệu đầu vào dạng số hoặc chuỗi có giới hạn** (password length, số ghế, phân trang, khoảng thời gian). Không áp dụng cho Ticket Export và RBAC vì các tính năng này chủ yếu xoay quanh trạng thái và quyền hạn, không có biên giá trị số rõ ràng cần kiểm tra.

- **Decision Tables**: Áp dụng cho các tính năng có **nhiều điều kiện kết hợp** ảnh hưởng đến kết quả. Đặc biệt quan trọng cho Login (tổ hợp email × password × trạng thái tài khoản), Booking (tổ hợp trạng thái ghế × suất chiếu × thanh toán), và RBAC (ma trận role × permission). Không áp dụng cho Revenue Analytics vì tính năng này chủ yếu là truy vấn đọc dữ liệu với các tham số lọc đơn giản.

- **State Transition**: Chỉ áp dụng cho các tính năng có **vòng đời trạng thái rõ ràng**: Login (trạng thái phiên), Seat management (trạng thái ghế), Booking (trạng thái đơn đặt), Ticket Export (trạng thái vé), RBAC (trạng thái quyền). Không áp dụng cho các tính năng CRUD đơn thuần (Register, Movie management, News management) vì không có chuỗi chuyển trạng thái phức tạp.

- **Pairwise Testing**: Chỉ áp dụng cho các tính năng có **nhiều tham số đầu vào kết hợp** (≥ 3 tham số): Movie management (tìm kiếm/lọc phim với nhiều filter), Online Booking (nhiều tùy chọn đặt vé), Revenue Analytics (nhiều tham số thống kê). Kỹ thuật này giúp giảm đáng kể số lượng test case so với kiểm thử tổ hợp đầy đủ mà vẫn đạt hiệu quả phát hiện lỗi cao.

