# UNIT TEST SUITE — TESTER DUCPM88 (PHẠM MINH ĐỨC - B22DCDT095)
## PHÂN HỆ: MOVIE MANAGEMENT & SEAT MANAGEMENT

Thư mục này chứa đầy đủ tài liệu, báo cáo Excel, và các script thực thi Unit Test cho dịch vụ `movie-service` của dự án Cinema Management System theo đúng yêu cầu đề bài.

---

### Danh mục các tệp tin trong thư mục:

1. **[`Unit_Testing_Report_DUCPM88.xlsx`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/docs/TESTER_DUCPM88/original_docs/unit_test/Unit_Testing_Report_DUCPM88.xlsx)**
   - Bảng tính Excel hoàn chỉnh gồm 5 sheets:
     - `1. Scope & Overview`: Tổng quan công cụ, phạm vi, kết quả và độ phủ.
     - `2. Excluded Files`: Danh sách các file không cần Unit test kèm giải thích kỹ thuật (1.2).
     - `3. Movie Management (71 TCs)`: Toàn bộ test case cho Movie (Entity + Business).
     - `4. Seat Management (44 TCs)`: Toàn bộ test case cho Seat (Entity + Business).
     - `5. Paging Helper (7 TCs)`: Test case cho package phân trang dùng chung.

2. **[`Unit_Testing_Report_DUCPM88.csv`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/docs/TESTER_DUCPM88/original_docs/unit_test/Unit_Testing_Report_DUCPM88.csv)**
   - File CSV định dạng chuẩn UTF-8 (có BOM), hỗ trợ mở trực tiếp bằng Microsoft Excel mà không bị lỗi font tiếng Việt.

3. **[`Unit_Testing_Report_DUCPM88.md`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/docs/TESTER_DUCPM88/original_docs/unit_test/Unit_Testing_Report_DUCPM88.md)**
   - Báo cáo chi tiết dạng Markdown trình bày đầy đủ từ mục 1.1 đến 1.7 và mục 2.

4. **[`run_tests.bat`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/docs/TESTER_DUCPM88/original_docs/unit_test/run_tests.bat)**
   - **Click đúp vào file này** để tự động chạy Unit Test, hiển thị kết quả Pass/Fail trong terminal và xuất báo cáo Code Coverage HTML.

5. **[`rollback.bat`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/docs/TESTER_DUCPM88/original_docs/unit_test/rollback.bat)**
   - **Click đúp vào file này** để thực hiện Rollback, dọn dẹp các tệp tin phát sinh sau khi chạy test và hiển thị giải thích về cơ chế toàn vẹn dữ liệu.

---

### Mã nguồn Unit Test trực tiếp trong dự án:
- [`movie-service/internal/module/movie/entity/movie_test.go`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/movie-service/internal/module/movie/entity/movie_test.go)
- [`movie-service/internal/module/movie/entity/types_test.go`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/movie-service/internal/module/movie/entity/types_test.go)
- [`movie-service/internal/module/movie/business/business_test.go`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/movie-service/internal/module/movie/business/business_test.go)
- [`movie-service/internal/module/seat/entity/seat_test.go`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/movie-service/internal/module/seat/entity/seat_test.go)
- [`movie-service/internal/module/seat/entity/types_test.go`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/movie-service/internal/module/seat/entity/types_test.go)
- [`movie-service/internal/module/seat/business/business_test.go`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/movie-service/internal/module/seat/business/business_test.go)
- [`movie-service/internal/pkg/paging/paging_test.go`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/movie-service/internal/pkg/paging/paging_test.go)
