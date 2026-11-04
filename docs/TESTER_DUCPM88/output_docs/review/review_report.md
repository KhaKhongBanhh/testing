# BÁO CÁO REVIEW — REQUIREMENTS SPECIFICATION & DESIGN

> **Dự án:** Hệ thống Quản lý Rạp chiếu phim (Cinema Management System) — Mã dự án: `PHI_2025`  
> **Phiên bản sản phẩm:** `v1.0`  
> **Phân hệ phụ trách:** Theater Manager (Movie Management, Room Management, Seat Management, Showtime Management)  
> **Người thực hiện (Reviewer):** Phạm Minh Đức — MSSV: `B22DCDT095` (Nhóm 05)  
> **Ngày thực hiện Review:** `04/11/2026`  
> **Tài liệu gốc đối chiếu:** `docs/graduation_thesis.pdf`, `docs/ARCHITECTURE.md`, `docs/API.md`

---

## MỤC LỤC

1. [Phần 1: Nâng cao SQA Plan cho hoạt động Review](#1-nâng-cao-sqa-plan-cho-hoạt-động-review)
2. [Phần 2: Review Checklist — Requirements Specification (J3)](#2-review-checklist--requirements-specification-j3)
3. [Phần 3: Review Checklist — Design (J5)](#3-review-checklist--design-j5)
4. [Phần 4: Tổng hợp kết quả Review](#4-tổng-hợp-kết-quả-review)

---

## 1. Nâng cao SQA Plan cho hoạt động Review

### 1.1. Xác định các tài liệu cần được Review

Tài liệu đặc tả yêu cầu và thiết kế của dự án Cinema Management System có quy mô lớn (tài liệu đồ án tốt nghiệp 85 trang, kèm các tài liệu bổ trợ `ARCHITECTURE.md`, `API.md`). Do đó, các tài liệu được tách thành các nhóm theo tính năng (feature-based groups) để review hiệu quả:

| STT | Tài liệu cần Review | Phạm vi / Mô tả | Mức ưu tiên |
|:---:|---|---|:---:|
| 1 | `graduation_thesis.pdf` — Chương 2, Mục 2.1 (Phân tích yêu cầu chức năng) | Biểu đồ Use Case tổng quan; Đặc tả 20 Use Cases (Đăng nhập, Đăng ký, Đặt vé, Quản lý phim, Quản lý phòng, Quản lý ghế, Quản lý suất chiếu, Quản lý nhân viên, Quản lý tin tức, Quản lý RAG chatbot, Chat, Xem thống kê, Xem tin tức, Xem thông báo, Bán vé tại quầy, Tra cứu vé, Tìm kiếm phim, Xem lịch sử đặt vé, Chỉnh sửa thông tin, Xuất vé) | High |
| 2 | `graduation_thesis.pdf` — Chương 2, Mục 2.2 (Phân tích yêu cầu phi chức năng) | Yêu cầu bảo mật (RBAC, JWT, CORS); Yêu cầu hiệu năng (Microservices, gRPC, Redis cache, Outbox Pattern) | High |
| 3 | `graduation_thesis.pdf` — Chương 2, Mục 2.3 (Thiết kế hệ thống) | Biểu đồ lớp (22 lớp); Thiết kế CSDL; Biểu đồ tuần tự (18 sơ đồ) | High |
| 4 | `docs/ARCHITECTURE.md` | Kiến trúc tổng thể Microservices (10 service); Mẫu giao tiếp (REST, gRPC, WebSocket, Outbox); Chiến lược bảo mật & hiệu năng | Medium |
| 5 | `docs/API.md` | Đặc tả REST API cho tất cả các endpoint (Auth, User, Movie, Booking, Payment, Notification, Analytics, Chatbot) | Medium |

### 1.2. Loại hình Review được thực hiện

| Tài liệu | Loại Review | Mô tả |
|---|---|---|
| Requirements Specification (Đặc tả yêu cầu) | **Guided Walkthrough** | Reviewer (Phạm Minh Đức) đọc kỹ từng Use Case trong phạm vi Theater Manager, đối chiếu với mã nguồn thực tế và checklist J3 để phát hiện các điểm thiếu sót, mâu thuẫn hoặc mơ hồ. |
| Design Documents (Tài liệu thiết kế) | **Meeting Review** | Reviewer kiểm tra biểu đồ lớp, sơ đồ CSDL, biểu đồ tuần tự và tài liệu kiến trúc theo checklist J5, đối chiếu với yêu cầu đặc tả và mã nguồn hiện hành. |

### 1.3. Thông tin bổ sung liên quan đến Review

| Thông tin | Chi tiết |
|---|---|
| **Phạm vi Review cụ thể** | 4 Use Cases thuộc Theater Manager: Quản lý phim (`Manage Movies`), Quản lý phòng chiếu (`Manage Rooms`), Quản lý ghế (`Manage Seats`), Quản lý suất chiếu (`Manage Showtimes`) |
| **Tài liệu tham chiếu chuẩn** | Sách "Mastering Software Quality Assurance" — Appendix J (Table J.3, J.4, J.5); IEEE 1028-2008 Standard for Software Reviews and Audits |
| **Công cụ hỗ trợ Review** | VS Code (đọc mã nguồn), PDF Reader (đọc thesis), Markdown editor (ghi chú) |
| **Effort dự kiến** | 4 MH (Man-Hours) |
| **Tiêu chí đánh giá** | Pass: Không tìm thấy lỗi; Fail: Tìm thấy ít nhất 1 lỗi; N/A: Không áp dụng được hoặc không thể trả lời |
| **Cơ chế báo cáo lỗi** | Tất cả các lỗi phát hiện qua Review được liệt kê chi tiết trong cột Notes của mỗi checklist, kèm giải thích và vị trí cụ thể trong tài liệu. |

---

## 2. Review Checklist — Requirements Specification (J3)

### 2.1. Thông tin Review

| Trường thông tin | Giá trị |
|---|---|
| **Project Name** | Cinema Management System (`PHI_2025`) |
| **Reviewed Document** | `graduation_thesis.pdf` — Chương 2: Phân tích yêu cầu và Thiết kế hệ thống (Mục 2.1, 2.2); `docs/API.md` |
| **Reviewed Feature / Scope** | Theater Manager — Quản lý phim, Quản lý phòng chiếu, Quản lý ghế ngồi, Quản lý suất chiếu |
| **Review Round** | 1 |
| **Reviewer** | Phạm Minh Đức (`B22DCDT095`) |
| **Review Date** | 04/11/2026 |

### 2.2. Tổng hợp kết quả

| Chỉ số | Giá trị |
|---|---|
| **Tổng số câu hỏi** | 34 |
| **Pass** | 18 |
| **Fail** | 13 |
| **N/A** | 3 |
| **Tổng số lỗi (Number of Errors)** | 21 |

### 2.3. Bảng câu hỏi Review — Requirements Specification

> **Nguồn:** Table J.3 & J.4 — "Mastering Software Quality Assurance" (J. Ross Publishing)  
> Các câu hỏi không phù hợp với dự án đã được loại bỏ hoặc đánh dấu N/A. Các câu hỏi đặc thù cho dự án Cinema đã được bổ sung ở cuối bảng.

| No. | Checklist Question | Review Result | Notes |
|:---:|---|:---:|---|
| 1 | Are the requirements in compliance with the contract? | **N/A** | Đây là đồ án tốt nghiệp, không có hợp đồng (contract) chính thức. Yêu cầu do nhóm tự đề xuất và được giảng viên hướng dẫn phê duyệt. |
| 2 | Have all the requirements been listed? | **Fail** | Tài liệu liệt kê 20 Use Cases nhưng thiếu một số yêu cầu ngầm định: (1) Không có Use Case cho quản lý thể loại phim (Genre CRUD) mặc dù mã nguồn có module genre riêng; (2) Không có Use Case cho cơ chế giữ ghế tạm thời (Seat Lock) — một tính năng cốt lõi của hệ thống đặt vé đồng thời; (3) Thiếu đặc tả Use Case cho API Health Check được cấu hình trong API Gateway. **Tổng: 3 lỗi.** |
| 3 | Are there any ambiguous requirements? | **Fail** | (1) UC "Quản lý ghế" ghi "hiển thị danh sách ghế theo chế độ bảng hoặc lưới" nhưng không rõ chế độ nào là mặc định, điều kiện chuyển đổi. (2) UC "Quản lý suất chiếu" ghi "thời lượng suất chiếu ngắn hơn thời lượng phim" trong Exception nhưng không giải thích cách tính thời lượng suất chiếu (`end_time - start_time`) có bao gồm thời gian buffer/dọn phòng hay không. **Tổng: 2 lỗi.** |
| 4 | Is each requirement described completely? | **Fail** | (1) UC "Quản lý phim" không mô tả các trường dữ liệu bắt buộc khi thêm phim (title, duration, release_date, genre_ids, poster_url...) — chỉ ghi "nhập đầy đủ thông tin phim". (2) UC "Quản lý suất chiếu" không đặc tả logic kiểm tra trùng lịch chiếu (overlap checking) — chỉ liệt kê trong Exception mà không mô tả thuật toán. (3) UC "Quản lý ghế" không nêu rõ các loại ghế hỗ trợ (REGULAR, VIP, COUPLE) và hệ số giá tương ứng (1.0x, 1.5x, 2.5x). **Tổng: 3 lỗi.** |
| 5 | Have the requirements been specified consistently throughout the document? | **Pass** | Các yêu cầu nhất quán giữa phần đặc tả Use Case (Chương 2.1) và phần thiết kế CSDL/biểu đồ lớp (Chương 2.3). Cấu trúc mô tả mỗi Use Case đồng nhất (Name, Description, Actor, Pre-condition, Post-condition, Flow, Exception). |
| 6 | Can the requirements be verified? | **Pass** | Tất cả 4 Use Cases thuộc Theater Manager đều có thể kiểm chứng bằng System Test (gọi API REST) và Unit Test (kiểm tra logic nghiệp vụ). |
| 7 | Has any additional functionality been included beyond the scope of the contract? | **N/A** | Không có hợp đồng. Phạm vi do nhóm tự xác định. |
| 8 | Are project management requirements included in the requirements? | **Pass** | Có mô tả về quy trình CI/CD (GitHub Actions), hạ tầng triển khai (Docker Compose), và chiến lược phân chia service. |
| 9 | Is the rationale for any derived requirements satisfactory? | **Pass** | Các yêu cầu phái sinh (như cơ chế cache Redis cho Movie, tách Read/Write DB) được giải thích hợp lý trong Chương 2.2 (Yêu cầu hiệu năng). |
| 10 | Are the specified external interfaces compatible? | **Pass** | Giao diện ngoài (REST API qua API Gateway, gRPC nội bộ) được mô tả chi tiết trong `API.md` và `ARCHITECTURE.md`, tương thích với kiến trúc Microservices. |
| 11 | Are the user interface requirements complete? | **Fail** | (1) Tài liệu có hình ảnh giao diện quản lý suất chiếu (Hình 48-49) và quản lý nhân viên (Hình 50-52) nhưng thiếu hình ảnh giao diện quản lý phim (Movie Form), quản lý phòng (Room Form) và quản lý ghế (Seat Form/Grid). (2) Không có wireframe hoặc mockup cho trạng thái lỗi/empty state của các form. **Tổng: 2 lỗi.** |
| 12 | Can the requirements be tested? Can the requirements be used directly for validation during acceptance testing? | **Fail** | Các Use Cases có thể test được, nhưng không thể dùng trực tiếp cho acceptance testing vì: (1) Thiếu Acceptance Criteria định lượng (thời gian phản hồi, số lượng concurrent users). (2) Exception chỉ mô tả thông báo lỗi chung ("Vui lòng điền đầy đủ thông tin") mà không gắn với mã lỗi HTTP cụ thể. **Tổng: 2 lỗi.** |
| 13 | Are the performance requirements adequate and feasible? | **Fail** | Chương 2.2 nêu "hệ thống đáp ứng yêu cầu hiệu năng" nhưng không định lượng: thiếu ngưỡng thời gian phản hồi API (ví dụ p95 < 200ms), thông lượng tối đa (requests/sec), số người dùng đồng thời tối đa. **Tổng: 1 lỗi.** |
| 14 | Have the security requirements been determined? | **Pass** | Chương 2.2.1 mô tả đầy đủ: RBAC với 4 vai trò, JWT (HMAC-SHA256), CORS, Redis session management. |
| 15 | Do any requirements conflict with or duplicate other requirements? | **Pass** | Không phát hiện mâu thuẫn hoặc trùng lặp giữa các Use Cases trong phạm vi Theater Manager. |
| 16 | Is each requirement written in clear, concise, unambiguous language? | **Pass** | Ngôn ngữ mô tả Use Case nhìn chung rõ ràng, dùng cấu trúc bước 1-2-3 cho Flow và Exception. |
| 17 | Is each requirement free of content and grammatical errors? | **Pass** | Không phát hiện lỗi nội dung hay ngữ pháp nghiêm trọng trong các Use Case thuộc phạm vi review. |
| 18 | Are the time-critical functions identified, and are the timing criteria for them specified? | **Fail** | Chức năng giữ ghế tạm thời (seat lock) được đề cập là time-critical (hết hạn sau 15 phút) nhưng không có Use Case riêng mô tả chi tiết luồng. Thời gian timeout cho booking PENDING (15 phút) chỉ được đề cập trong API.md mà không nằm trong SRS. **Tổng: 1 lỗi.** |
| 19 | Have internationalization issues been adequately addressed? | **N/A** | Hệ thống chỉ phục vụ thị trường Việt Nam, giao diện tiếng Việt. Không có yêu cầu đa ngôn ngữ. |
| 20 | Is the format in conformance with the format in the organizational process? | **Pass** | Tài liệu tuân theo định dạng chuẩn đồ án tốt nghiệp PTIT: Lời mở đầu, Chương 1-4, Tài liệu tham khảo. Mỗi Use Case theo chuẩn bảng (Name, Description, Actor, Pre/Post-condition, Flow, Exception). |
| 21 | Are all internal cross-references to other requirements correct? | **Pass** | Các tham chiếu chéo giữa Use Case quản lý suất chiếu (tham chiếu phim, phòng) và Use Case quản lý phòng (tham chiếu ghế) đều chính xác. |
| 22 | Do the requirements provide an adequate basis for software requirements specification? | **Pass** | Các Use Cases cung cấp đủ cơ sở cho thiết kế chi tiết (biểu đồ lớp, CSDL, API endpoint). |
| 23 | Have algorithms intrinsic to the functional requirements been defined? | **Fail** | (1) Thuật toán kiểm tra trùng lịch chiếu (time overlap checking) chưa được đặc tả — chỉ có exception "Lịch chiếu bị trùng với lịch chiếu khác trong phòng" mà không mô tả logic so sánh `start_time`/`end_time`. (2) Thuật toán làm tròn thời gian 30 phút (`TruncateToHalfHour`) không được đề cập trong SRS. (3) Công thức tính giá ghế (`seat_price = base_price * seat_type_multiplier`) không được đặc tả. **Tổng: 3 lỗi.** |
| 24 | Is each requirement in scope for the project? | **Pass** | Toàn bộ 4 Use Cases Theater Manager nằm trong phạm vi dự án. |
| 25 | Are all security and safety considerations properly specified? | **Pass** | Cơ chế bảo mật JWT, RBAC, CORS được đặc tả rõ trong Chương 2.2.1. |
| 26 | Have periodic processing procedures (for example, daily, monthly) been included? | **Fail** | Tài liệu không đề cập đến xử lý định kỳ như: (1) Tự động cập nhật trạng thái phim từ UPCOMING → SHOWING → ENDED theo thời gian; (2) Dọn dẹp booking PENDING hết hạn; (3) Tự động hủy suất chiếu quá hạn. Trong mã nguồn, worker-service xử lý các tác vụ này nhưng không có đặc tả trong SRS. **Tổng: 1 lỗi.** |
| 27 | Have all interfaces between components been identified? | **Pass** | `ARCHITECTURE.md` mô tả đầy đủ giao tiếp giữa 10 services qua REST, gRPC và Redis Pub/Sub. |
| 28 | Are the interfaces provided easy to use and consistent in format? | **Pass** | API tuân theo chuẩn RESTful, nhất quán về format JSON response, HTTP status codes, và URL pattern `/api/v1/{resource}`. |
| 29 | Have all external user interfaces been identified? | **Pass** | Giao diện: Web browser (React SPA), API Gateway (REST), WebSocket (Notifications). |
| 30 | Would this document be adequate to be able to proceed with software design description? | **Pass** | Tài liệu đủ cơ sở để thiết kế chi tiết — thực tế Chương 2.3 đã trình bày thiết kế CSDL và biểu đồ lớp dựa trên đặc tả này. |
| 31 | Has any additional functionality been included (exceeding the scope of the contract)? | **N/A** | Không áp dụng — không có hợp đồng. |
| 32 | Are all requirements in the user requirements specification included in this document? | **Pass** | Tất cả yêu cầu từ biểu đồ Use Case tổng quan (Hình 16) đều được đặc tả chi tiết. |
| 33 | Is any necessary information missing from a requirement? If so, is it identified as TBD? | **Fail** | Nhiều thông tin cần thiết bị thiếu nhưng không được đánh dấu TBD: (1) Định dạng phim hỗ trợ (2D, 3D, IMAX) — trong mã nguồn có trường `format` nhưng SRS không đề cập. (2) Trạng thái hợp lệ và quy tắc chuyển trạng thái (state machine) cho Movie, Room, Seat, Showtime. (3) Giới hạn số ghế tối đa mỗi phòng. **Tổng: 3 lỗi.** |
| 34 | Is the expected behavior documented for all anticipated error conditions? | **Fail** | Các Exception trong Use Cases chỉ mô tả thông báo lỗi chung. Thiếu: (1) Hành vi khi xóa phòng đang có suất chiếu active. (2) Hành vi khi cập nhật phim đang chiếu (SHOWING). (3) Hành vi khi tạo suất chiếu cho phim có status ENDED. (4) Hành vi khi xóa ghế đang được giữ (LOCKED/RESERVED). **Tổng: 4 lỗi (nhưng gộp thành 1 mục Fail).** |
| **P1** | **[Project-specific] Có đặc tả rõ các trạng thái hợp lệ (state machine) và quy tắc chuyển trạng thái cho các entity chính (Movie, Room, Seat, Showtime)?** | **Fail** | Tài liệu hoàn toàn không mô tả state machine. Trong mã nguồn, Movie có 3 trạng thái (UPCOMING → SHOWING → ENDED), Room có 2 (ACTIVE, INACTIVE), Seat có 3 (AVAILABLE, RESERVED, MAINTENANCE), Showtime có 3 (SCHEDULED, COMPLETED, CANCELLED) — nhưng không có đặc tả nào trong SRS. **Tổng: 1 lỗi.** |
| **P2** | **[Project-specific] Có đặc tả rõ cơ chế xử lý đồng thời (concurrency) khi nhiều người cùng đặt ghế?** | **Fail** | SRS chỉ nêu chung "ghế bạn chọn đã được đặt" nhưng không đặc tả: (1) Cơ chế lock ghế tạm thời (Redis-based distributed lock). (2) Thời gian timeout giữ ghế (15 phút — chỉ có trong API.md). (3) Hành vi khi 2 người cùng chọn 1 ghế. **Tổng: 1 lỗi.** |

### 2.4. Kết luận Review — Requirements Specification

**Kết luận: `Conditionally Approved`**

**Discussion:**  
Tài liệu đặc tả yêu cầu (SRS) phiên bản v1.0 của hệ thống Cinema Management đã đáp ứng các nội dung nền tảng quan trọng: mục đích hệ thống, phạm vi dự án, xác định đầy đủ 4 tác nhân, biểu đồ Use Case tổng quan, và đặc tả chi tiết 20 Use Cases theo cấu trúc chuẩn. Các yêu cầu phi chức năng về bảo mật (RBAC, JWT) và kiến trúc (Microservices) được trình bày rõ ràng. Tuy nhiên, review phát hiện **21 lỗi/thiếu sót** ở mức nghiêm trọng trung bình, tập trung vào 3 vấn đề chính: (1) **Thiếu đặc tả thuật toán nghiệp vụ** — logic kiểm tra trùng lịch chiếu, công thức tính giá ghế, cơ chế làm tròn thời gian 30 phút đều chỉ tồn tại trong mã nguồn mà không có trong SRS; (2) **Thiếu state machine** — quy tắc chuyển trạng thái cho Movie/Room/Seat/Showtime không được đặc tả dẫn đến nhiều lỗi nghiệp vụ trong mã nguồn (module Room, Seat, Showtime khai báo `ErrInvalidStatusTransition` nhưng không sử dụng); (3) **Thiếu yêu cầu phi chức năng định lượng** — không có ngưỡng hiệu năng, SLA, backup policy. Khuyến nghị bổ sung các nội dung trên trước khi chuyển sang phiên bản v2.0.

---

## 3. Review Checklist — Design (J5)

### 3.1. Thông tin Review

| Trường thông tin | Giá trị |
|---|---|
| **Project Name** | Cinema Management System (`PHI_2025`) |
| **Reviewed Document** | `graduation_thesis.pdf` — Chương 2, Mục 2.3 (Thiết kế hệ thống): Biểu đồ lớp, Thiết kế CSDL, Biểu đồ tuần tự; `docs/ARCHITECTURE.md` |
| **Reviewed Feature / Scope** | Theater Manager — Thiết kế module Movie, Room, Seat, Showtime (biểu đồ lớp, CSDL, API, kiến trúc phân tầng) |
| **Review Round** | 1 |
| **Reviewer** | Phạm Minh Đức (`B22DCDT095`) |
| **Review Date** | 04/11/2026 |

### 3.2. Tổng hợp kết quả

| Chỉ số | Giá trị |
|---|---|
| **Tổng số câu hỏi** | 35 |
| **Pass** | 18 |
| **Fail** | 13 |
| **N/A** | 4 |
| **Tổng số lỗi (Number of Errors)** | 22 |

### 3.3. Bảng câu hỏi Review — Design

> **Nguồn:** Table J.5 — "Mastering Software Quality Assurance" (J. Ross Publishing)  
> Các câu hỏi không phù hợp với dự án đã được loại bỏ hoặc đánh dấu N/A. Các câu hỏi đặc thù cho dự án Cinema đã được bổ sung ở cuối bảng.

| No. | Checklist Question | Review Result | Notes |
|:---:|---|:---:|---|
| 1 | Do the programming specifications fulfill the program objectives? | **Pass** | Thiết kế biểu đồ lớp (22 lớp) và kiến trúc Microservices (10 services) đáp ứng mục tiêu xây dựng hệ thống đặt vé phim có khả năng mở rộng và xử lý đồng thời cao. |
| 2 | Have the common and shared functions been defined? | **Pass** | Các hàm/module dùng chung được thiết kế rõ: `paging` package (phân trang), `caching` package (Redis cache utils), `middleware` (auth, CORS, rate limit), `Dependency Injection` qua `samber/do`. |
| 3 | Has each module of the software requirements specification been included? | **Fail** | (1) Thiết kế thiếu module quản lý thể loại phim (Genre Management) — trong mã nguồn có nhưng không xuất hiện trong biểu đồ lớp chi tiết. (2) Thiếu thiết kế chi tiết cho module Seat Lock (giữ ghế tạm thời) — một tính năng quan trọng của hệ thống đặt vé. **Tổng: 2 lỗi.** |
| 4 | Have all the exception conditions been handled? | **Fail** | (1) Thiết kế CSDL không có ràng buộc CHECK cho trường `status` trên bảng movies, rooms, seats, showtimes — cho phép lưu giá trị bất kỳ. (2) Thiếu thiết kế cho trường hợp xóa phòng đang có suất chiếu (ON DELETE restrict/cascade). (3) Thiếu thiết kế cho trường hợp cập nhật trạng thái phim khi đang có suất chiếu SCHEDULED. **Tổng: 3 lỗi.** |
| 5 | Does the design conform to structured methodologies? | **Pass** | Tuân theo Clean Architecture / Layered Architecture: `entity` → `repository` → `business` → `transport` (REST & gRPC). Mỗi module là một package Go độc lập, giao tiếp qua interface. |
| 6 | Are there safeguards against data overflow? | **Fail** | (1) Thiết kế CSDL không quy định giới hạn `capacity` tối đa cho bảng `rooms`. (2) Không có giới hạn số lượng ghế tối đa mỗi phòng. (3) Trường `base_price` trong `showtimes` không có ràng buộc `CHECK (base_price > 0)`. (4) Trường `duration` trong `movies` không có ràng buộc `CHECK (duration > 0)`. **Tổng: 4 lỗi.** |
| 7 | Can the program specification be easily coded? | **Pass** | Thiết kế sử dụng các pattern quen thuộc (Repository, DI, Middleware) với Go và PostgreSQL — dễ triển khai mã nguồn. Thực tế đã được cài đặt thành công. |
| 8 | Have the loop termination conditions been properly taken care of? | **Pass** | Thiết kế không có vòng lặp phức tạp. Các truy vấn danh sách đều có phân trang (`LIMIT/OFFSET`). Redis SCAN pattern có điều kiện dừng `cursor == 0`. |
| 9 | Do the nesting conditions conform to standards? | **Pass** | Thiết kế phân tầng rõ ràng, không có nesting quá sâu: Transport → Business → Repository → Database. |
| 10 | Are the modules independent and self-sufficient? | **Pass** | Mỗi module (movie, room, seat, showtime) có package riêng với entity, repository, business, transport. Giao tiếp qua interface `Biz` và `Repository`. |
| 11 | Are the test plans developed? | **Fail** | SQA Plan (Section 6.3) có bảng Proposed Tests nhưng thiếu: (1) Test plan chi tiết cho từng module Theater Manager. (2) Không có ma trận truy vết (Traceability Matrix) giữa Use Cases và Test Cases. **Tổng: 2 lỗi.** |
| 12 | Are the screen layouts, report layouts, and table and file usage in the program specifications identical to those specified in the software requirements specification? | **Fail** | (1) Schema CSDL thực tế trong mã nguồn có thêm các trường không có trong thiết kế: `format` (định dạng phim 2D/3D/IMAX) trong bảng `movies`, `status` trong bảng `rooms` và `seats`. (2) Tài liệu thiết kế bảng `seats` ghi `seat_type` nhưng thực tế mã nguồn có thêm `seat_number` là kết hợp `row_label + column_number`. **Tổng: 2 lỗi.** |
| 13 | Have the error messages been defined? | **Fail** | Thiết kế không có danh mục mã lỗi/thông báo lỗi chuẩn hóa. Các sentinel errors (`ErrMovieNotFound`, `ErrRoomNotFound`, `ErrInvalidStatusTransition`, `ErrMovieNotShowing`...) chỉ tồn tại trong mã nguồn mà không có trong tài liệu thiết kế. **Tổng: 1 lỗi.** |
| 14 | Are the program specifications consistent with the software requirements specification process logic? | **Fail** | (1) SRS ghi "Quản lý phim" cho Actor "Quản trị viên, nhân viên quản lý" nhưng thiết kế API (`API.md`) đặt endpoint tạo/sửa/xóa phim dưới `/api/v1/admin/movies` — chỉ Admin mới truy cập được, không phải Manager. (2) SRS ghi "nhân viên quản lý" có quyền quản lý nhân viên, nhưng `ARCHITECTURE.md` liệt kê roles chỉ có Admin, Staff, Customer — thiếu vai trò Manager Staff riêng biệt. **Tổng: 2 lỗi.** |
| 15 | Was each entity transformed into a table? | **Pass** | 22 lớp trong biểu đồ lớp đều được ánh xạ thành bảng CSDL tương ứng: `users`, `roles`, `permissions`, `movies`, `genres`, `movie_genres`, `rooms`, `seats`, `showtimes`, `bookings`, `tickets`, `payments`, `notifications`, `documents`, `document_chunks`, `chats`, `outbox_events`, v.v. |
| 16 | Has each attribute been mapped to a field in the database table of the entity? | **Fail** | (1) Lớp `Movie` trong biểu đồ lớp có thuộc tính `director`, `cast` nhưng schema CSDL trong `ARCHITECTURE.md` không có các trường này. (2) Lớp `Showtime` có thuộc tính `status` nhưng schema mẫu trong `ARCHITECTURE.md` (`showtimes` table) không liệt kê trường `status`. **Tổng: 2 lỗi.** |
| 17 | Has each key field been mapped to an index? | **Pass** | Primary keys (UUID) có B-tree index mặc định. Foreign keys (`movie_id`, `room_id`, `seat_id`, `showtime_id`) được index cho join performance. GIN index cho full-text search trên `movies.title`. |
| 18 | Have all foreign keys been specified NOT NULL (if applicable)? | **Pass** | Các foreign keys quan trọng (`showtimes.movie_id`, `showtimes.room_id`, `seats.room_id`, `tickets.seat_id`, `bookings.showtime_id`) đều là NOT NULL trong migration scripts. |
| 19 | Are there any data integrity rules that violate the referential integrity rules of the relational database management system? | **Pass** | Không phát hiện vi phạm ràng buộc toàn vẹn tham chiếu. Các quan hệ FK đều hợp lệ. |
| 20 | Have the triggers, events, objects such as stored procedures and functions, and actions been defined where the relational database management system supports them? | **N/A** | Hệ thống không sử dụng triggers/stored procedures. Logic nghiệp vụ được xử lý hoàn toàn trong application layer (Go business logic). |
| 21 | Do all the views have a valid purpose? | **N/A** | Không sử dụng database views. Analytics service truy vấn trực tiếp các bảng. |
| 22 | Have the views that can be updated been correctly designed? | **N/A** | Không áp dụng — không sử dụng views. |
| 23 | Have the free space parameters been set depending on table size and table use? Has enough free space been specified? | **N/A** | Thiết kế không đề cập free space parameters. PostgreSQL tự quản lý qua autovacuum. |
| 24 | Have stored procedures been designed for frequently executed transactions? | **Pass** | Không sử dụng stored procedures theo thiết kế có chủ đích — logic đặt trong application layer để dễ test và deploy theo kiến trúc Microservices. Tuy nhiên, các truy vấn phức tạp (kiểm tra trùng lịch chiếu) sử dụng raw SQL query tối ưu. |
| 25 | Have all the queries been analyzed to determine the most frequently used columns in the DISTINCT, GROUP BY, ORDER BY, and WHERE clauses? | **Fail** | Thiết kế không có phân tích query pattern. Thiếu index cho: (1) `showtimes.start_time` — truy vấn thường xuyên theo thời gian; (2) Composite index `(room_id, start_time, end_time)` trên `showtimes` cho overlap checking. **Tổng: 1 lỗi.** |
| 26 | Have all the indexes created been used? | **Pass** | Các index hiện có (PK, FK, GIN full-text) đều được sử dụng bởi các truy vấn tương ứng. |
| 27 | Have the bottlenecks in meeting performance requirements been identified? | **Fail** | (1) Thiết kế không nhận diện bottleneck tiềm ẩn khi sử dụng `redis.Keys()` (O(N) blocking) thay vì `SCAN` trong module Seat. (2) Không nhận diện vấn đề invalidate toàn bộ movie cache khi chỉ sửa 1 phim. (3) Không có thiết kế cho connection pool sizing. **Tổng: 3 lỗi.** |
| 28 | Does the database implicitly provide for locking? | **Pass** | PostgreSQL cung cấp row-level locking (FOR UPDATE) cho đặt vé đồng thời. Redis distributed lock cho giữ ghế tạm thời. Thiết kế tận dụng cả hai. |
| 29 | Have time-outs for locks and sessions been optimized? | **Fail** | (1) Thiết kế không quy định timeout cho database row lock khi đặt vé. (2) Timeout giữ ghế Redis (15 phút) chỉ xuất hiện trong mã nguồn, không có trong tài liệu thiết kế. (3) Cache TTL (5/15/30/60 phút) không được đặc tả trong thiết kế — chỉ tồn tại trong `const.go`. **Tổng: 1 lỗi.** |
| 30 | Have the roles and schemas been adequately defined to address the security issues? | **Pass** | Thiết kế RBAC rõ ràng: bảng `roles`, `permissions`, `role_permissions`. 4 vai trò (Admin, Manager Staff, Ticket Staff, Customer) với quyền hạn phân cấp. |
| 31 | Have privileges and authorities, data access control, etc. been implemented? | **Pass** | Thiết kế đã mô tả: JWT authentication, role-based authorization, API Gateway middleware, public/protected/admin routes. |
| **P1** | **[Project-specific] Thiết kế có mô tả đầy đủ kiến trúc phân tầng (Layered Architecture) cho từng module?** | **Pass** | Thiết kế nêu rõ kiến trúc phân tầng trong `ARCHITECTURE.md`: Transport (REST/gRPC) → Business Logic → Repository → Database. Mỗi module có cấu trúc package riêng biệt. |
| **P2** | **[Project-specific] Thiết kế biểu đồ tuần tự (Sequence Diagram) có bao phủ đầy đủ các luồng chính của Theater Manager?** | **Fail** | Chỉ có biểu đồ tuần tự cho luồng chung (Đăng nhập, Đặt vé, Thanh toán). Thiếu biểu đồ tuần tự riêng cho: (1) Tạo/Sửa/Xóa phim; (2) Tạo suất chiếu với kiểm tra trùng lịch; (3) Tạo ghế hàng loạt. **Tổng: 1 lỗi.** |
| **P3** | **[Project-specific] Thiết kế có mô tả cơ chế caching (Redis) cho từng module?** | **Fail** | `ARCHITECTURE.md` nêu chung "Redis caching for popular movies" nhưng không có thiết kế chi tiết: (1) Key pattern cho từng entity. (2) TTL cho từng loại cache. (3) Chiến lược invalidation (invalidate all vs. invalidate single). Toàn bộ thiết kế cache chỉ tồn tại trong mã nguồn (`const.go`). **Tổng: 1 lỗi.** |
| **P4** | **[Project-specific] Thiết kế CSDL có ràng buộc UNIQUE phù hợp để ngăn dữ liệu trùng lặp?** | **Pass** | Có UNIQUE constraint trên `users.email`, `genres.name`, composite unique `(room_id, row_label, column_number)` trên `seats` để ngăn trùng ghế trong phòng. |

### 3.4. Kết luận Review — Design

**Kết luận: `Conditionally Approved`**

**Discussion:**  
Thiết kế hệ thống Cinema Management thể hiện tư duy kiến trúc tốt: áp dụng Microservices với 10 services tách biệt, Clean Architecture cho từng module Go, CSDL PostgreSQL được chuẩn hóa với các quan hệ FK hợp lý, và tích hợp Redis cho caching/session/lock. Biểu đồ lớp 22 lớp phản ánh đầy đủ domain model. Tuy nhiên, review phát hiện **22 lỗi/thiếu sót**, tập trung vào 3 vấn đề chính: (1) **Thiếu ràng buộc CSDL** — không có CHECK constraint cho status, price, duration; không có ràng buộc cascade khi xóa entity cha có entity con phụ thuộc; (2) **Thiếu nhất quán giữa thiết kế và mã nguồn** — schema CSDL trong tài liệu thiếu nhiều trường thực tế có trong code (`format`, `status` trên rooms/seats), vai trò "Manager Staff" trong SRS không khớp với roles trong thiết kế API; (3) **Thiếu thiết kế chi tiết cho các thành phần vận hành** — caching strategy, error messages catalog, query performance analysis, timeout configuration đều chỉ tồn tại trong mã nguồn mà không có tài liệu thiết kế. Khuyến nghị bổ sung các thiết kế chi tiết trên và đồng bộ giữa tài liệu và mã nguồn.

---

## 4. Tổng hợp kết quả Review

### 4.1. Bảng tổng hợp

| Artifact | Pass | Fail | N/A | Tổng câu hỏi | Số lỗi phát hiện | Tỷ lệ đạt | Kết luận |
|---|:---:|:---:|:---:|:---:|:---:|:---:|---|
| Requirements Specification (J3) | 18 | 13 | 3 | 34 | 21 | 58.1% (18/31) | Conditionally Approved |
| Design (J5) | 18 | 13 | 4 | 35 | 22 | 58.1% (18/31) | Conditionally Approved |
| **Tổng cộng** | **36** | **26** | **7** | **69** | **43** | **58.1%** | **Conditionally Approved** |

### 4.2. Phân loại lỗi theo mức độ nghiêm trọng

| Mức độ | SRS | Design | Tổng | Mô tả tiêu biểu |
|---|:---:|:---:|:---:|---|
| **Critical** | 2 | 1 | 3 | Thiếu đặc tả state machine & concurrency handling; Thiếu ràng buộc CSDL cho status |
| **Major** | 8 | 10 | 18 | Thiếu thuật toán nghiệp vụ; Thiếu acceptance criteria; Không nhất quán giữa thiết kế và code |
| **Minor** | 11 | 11 | 22 | Thiếu hình ảnh UI; Thiếu mô tả xử lý định kỳ; Thiếu thiết kế cache chi tiết |

### 4.3. Danh sách các khuyến nghị ưu tiên

| # | Khuyến nghị | Artifact | Mức ưu tiên |
|:---:|---|---|:---:|
| 1 | Bổ sung đặc tả State Machine (biểu đồ chuyển trạng thái) cho Movie, Room, Seat, Showtime | SRS | **Critical** |
| 2 | Bổ sung CHECK constraints cho status, price, duration trong CSDL schema | Design | **Critical** |
| 3 | Bổ sung đặc tả thuật toán: kiểm tra trùng lịch chiếu, tính giá ghế, làm tròn 30 phút | SRS | **High** |
| 4 | Bổ sung Acceptance Criteria định lượng cho mỗi Use Case (thời gian phản hồi, concurrent users) | SRS | **High** |
| 5 | Đồng bộ schema CSDL giữa tài liệu thiết kế và mã nguồn (bổ sung `format`, `status`, `director`, `cast`) | Design | **High** |
| 6 | Bổ sung danh mục Error Messages / Error Codes chuẩn hóa | Design | **High** |
| 7 | Bổ sung thiết kế chi tiết Caching Strategy (key pattern, TTL, invalidation) | Design | **Medium** |
| 8 | Bổ sung biểu đồ tuần tự cho các luồng CRUD Theater Manager | Design | **Medium** |
| 9 | Bổ sung Use Case cho Genre Management và Seat Lock mechanism | SRS | **Medium** |
| 10 | Bổ sung Traceability Matrix giữa Use Cases ↔ Test Cases | Design | **Medium** |

---

*Báo cáo này được thực hiện bởi Phạm Minh Đức (B22DCDT095) theo yêu cầu môn học Đảm bảo Chất lượng Phần mềm, dựa trên checklist J3 và J5 từ sách "Mastering Software Quality Assurance" và đối chiếu với tài liệu dự án Cinema Management System.*

