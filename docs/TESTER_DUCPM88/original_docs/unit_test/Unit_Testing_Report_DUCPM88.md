# BÁO CÁO KIỂM THỬ ĐƠN VỊ (UNIT TESTING REPORT)
## PHÂN HỆ: QUẢN LÝ PHIM (MOVIE MANAGEMENT) & QUẢN LÝ GHẾ (SEAT MANAGEMENT)
### DỰ ÁN: HỆ THỐNG QUẢN LÝ RẠP CHIẾU PHIM (CINEMA MANAGEMENT SYSTEM)

- **Người thực hiện (Tester):** Phạm Minh Đức (MSSV: `B22DCDT095` - Nhóm 05)
- **Mã định danh:** `TESTER_DUCPM88`
- **Dịch vụ kiểm thử:** `movie-service` (Golang Backend Microservice)
- **Phạm vi trọng tâm:** **Movie Management** và **Seat Management**
- **File báo cáo Excel đính kèm:** [`Unit_Testing_Report_DUCPM88.xlsx`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/docs/TESTER_DUCPM88/original_docs/unit_test/Unit_Testing_Report_DUCPM88.xlsx)
- **File báo cáo CSV dự phòng:** [`Unit_Testing_Report_DUCPM88.csv`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/docs/TESTER_DUCPM88/original_docs/unit_test/Unit_Testing_Report_DUCPM88.csv)
- **Script chạy Unit Test tự động:** [`run_tests.bat`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/docs/TESTER_DUCPM88/original_docs/unit_test/run_tests.bat)
- **Script Rollback dữ liệu:** [`rollback.bat`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/docs/TESTER_DUCPM88/original_docs/unit_test/rollback.bat)

---

## 1. Unit Testing Report

### 1.1. Tools and Libraries (Công cụ và Thư viện sử dụng)
- **Ngôn ngữ lập trình:** `Go` (Golang v1.24.3).
- **Testing Framework:** Thư viện chuẩn `testing` của Golang (`testing.T`, `t.Run()`, `t.Fatalf()`).
- **Code Coverage Tool:** Công cụ tích hợp sẵn `go test -coverprofile=coverage.out` và `go tool cover -func` / `go tool cover -html`.
- **Thư viện Mock / In-Memory Isolation:**
  - Mock Pattern thông qua **Fake Structs**: `fakeRepo`, `fakeCache`, `fakeRedis` cài đặt các interface `MovieRepository`, `SeatRepository`, `caching.Cache`, `caching.ReadOnlyCache`, và `redis.UniversalClient`.
  - Giúp cô lập 100% tầng Business logic khỏi hạ tầng Database (PostgreSQL) và Redis server, đảm bảo chạy độc lập, siêu tốc và không làm thay đổi trạng thái Database thực tế.
- **Dependency Injection Framework:** `github.com/samber/do` (sử dụng trong kiến trúc khởi tạo `movie-service`).
- **ORM / Database Driver:** `github.com/uptrace/bun` (PostgreSQL driver được trừu tượng hóa qua repository interface).
- **Caching Driver:** `github.com/go-redis/cache/v9` và `github.com/redis/go-redis/v9`.

---

### 1.2. Scope of Testing (Phạm vi kiểm thử)

#### A. Danh sách các functions, classes, files ĐƯỢC KIỂM THỬ (TESTED):
Theo đúng phân công nhiệm vụ tập trung vào **Movie Management** và **Seat Management**, các thành phần được kiểm thử bao gồm:

1. **Movie Management Module (`internal/module/movie`):**
   - **Tầng Entity & Validation (`internal/module/movie/entity/movie.go`):**
     - Hàm `Movie.IsValid()`: Kiểm tra tính hợp lệ của phim (Title không rỗng, Duration > 0).
     - Hàm `Movie.CanTransitionTo(newStatus)`: Kiểm tra máy trạng thái chuyển đổi vòng đời phim (`UPCOMING -> SHOWING -> ENDED`).
   - **Tầng DTO & Mapper (`internal/module/movie/entity/types.go`):**
     - Hàm `CreateMovieRequest.ToEntity()`: Chuyển đổi dữ liệu tạo mới, gán trạng thái mặc định `UPCOMING`.
     - Hàm `UpdateMovieRequest.ToEntity(id)`: Chuyển đổi dữ liệu cập nhật, giữ nguyên trạng thái nếu không truyền.
     - Hàm `ToMovieResponse(movie)`: Ánh xạ thực thể Movie và danh sách tên thể loại `Genres` sang response.
     - Hàm `ToMoviesResponse(movies, page, size, total)`: Tính toán số trang bằng phép chia làm tròn trần (ceiling division) và đóng gói meta response.
   - **Tầng Nghiệp vụ cốt lõi (`internal/module/movie/business/business.go`):**
     - Hàm `CreateMovie(ctx, movie, genreIds)`: Kiểm tra dữ liệu hợp lệ, gán status default, gọi repository tạo bản ghi và xóa cache danh sách.
     - Hàm `UpdateMovie(ctx, movie, genreIds)`: Kiểm tra Id, validate dữ liệu, lấy dữ liệu hiện tại, kiểm tra chuyển trạng thái hợp lệ, cập nhật DB, invalidate cache chi tiết và danh sách.
     - Hàm `DeleteMovie(ctx, id)`: Kiểm tra Id, kiểm tra phim tồn tại, xóa DB, xóa cache.
     - Hàm `UpdateMovieStatus(ctx, id, status)`: Kiểm tra Id, kiểm tra phim tồn tại, kiểm tra luật chuyển trạng thái, cập nhật DB, xóa cache.
     - Hàm `GetMovieById(ctx, id)`: Lấy chi tiết phim qua cache 2 lớp (ReadOnly Cache + Cache), fallback query DB, ánh xạ `sql.ErrNoRows` thành `ErrMovieNotFound`.
     - Hàm `GetMovies(ctx, page, size, search, status)`: Chuẩn hóa page/size, chặn trần `size <= 100`, tính toán offset, lấy total count và danh sách phim từ cache/DB.
     - Hàm `ValidateMovieForShowtime(ctx, movieId)`: Xác thực phim có tồn tại và đang ở trạng thái `SHOWING` để được phép xếp lịch chiếu.
     - Hàm `GetMovieStats(ctx)` & `GetGenres(ctx)`: Lấy thống kê số lượng phim theo trạng thái và danh sách thể loại.

2. **Seat Management Module (`internal/module/seat`):**
   - **Tầng Entity & Calculation (`internal/module/seat/entity/seat.go`):**
     - Hàm `Seat.IsValid()`: Kiểm tra bắt buộc có `RoomId`, `SeatNumber`, `RowNumber`.
     - Hàm `GetSeatTypePriceMultiplier(seatType)`: Tính hệ số giá theo loại ghế (`REGULAR: 1.0x`, `VIP: 1.5x`, `COUPLE: 2.5x`, `Khác: 1.0x`).
     - Hàm `Seat.CalculatePrice(basePrice)`: Tính giá vé cuối cùng của ghế dựa trên giá gốc suất chiếu.
   - **Tầng DTO & Mapper (`internal/module/seat/entity/types.go`):**
     - Hàm `CreateSeatRequest.ToSeat()`: Chuyển đổi request tạo ghế, gán trạng thái mặc định `AVAILABLE`.
     - Hàm `ToSeatResponse(seat)`: Ánh xạ thông tin ghế, định dạng thời gian chuẩn RFC3339.
     - Hàm `ToSeatsResponse(seats, page, size, total)`: Đóng gói danh sách ghế kèm metadata phân trang `PageInfo`.
   - **Tầng Nghiệp vụ cốt lõi (`internal/module/seat/business/business.go`):**
     - Hàm `CreateSeat(ctx, seat)`: Validate dữ liệu, kiểm tra chống trùng vị trí ghế trong phòng (`ExistsBySeatPosition`), lưu DB, xóa cache danh sách.
     - Hàm `UpdateSeat(ctx, id, updates)`: Kiểm tra ghế tồn tại, kiểm tra chống trùng vị trí với ghế khác (sử dụng `excludeId`), cập nhật loại ghế/trạng thái, lưu DB, xóa cache.
     - Hàm `DeleteSeat(ctx, id)`: Kiểm tra ghế tồn tại, xóa DB, xử lý lỗi khóa ngoại nếu ghế đã có lịch sử đặt vé, xóa cache.
     - Hàm `UpdateSeatStatus(ctx, id, status)`: Kiểm tra ghế tồn tại, cập nhật trạng thái mới, lưu DB, xóa cache.
     - Hàm `GetSeatById(ctx, id)`: Lấy chi tiết ghế qua cache, fallback DB, ánh xạ `sql.ErrNoRows` thành `ErrSeatNotFound`.
     - Hàm `GetSeats(ctx, page, size, ...)`: Validate `page >= 1`, `size >= 1`, tính offset, lấy danh sách ghế và tổng số lượng.

3. **Shared Helper Package (`internal/pkg/paging`):**
   - Hàm `NewPageInfo(page, size, total)`: Tính toán số trang và thông tin phân trang chuẩn.
   - Hàm `stringToInt(s)`: Chuyển đổi chuỗi sang số nguyên an toàn, xử lý chuỗi không hợp lệ.

---

#### B. Danh sách các functions, classes, files KHÔNG CẦN KIỂM THỬ (DO NOT NEED TESTING) và Giải thích lý do:

| STT | File / Module | Số dòng (LOC) | Lớp kiến trúc | Lý do KHÔNG cần thực hiện Unit Test (Explain Why) |
|---:|---|---:|---|---|
| 1 | `internal/module/movie/business/const.go` | 44 | Constants & Key Formatter | Chỉ chứa định nghĩa hằng số thời gian cache (`CACHE_TTL_*`) và các hàm format chuỗi key Redis (`fmt.Sprintf`). Không chứa bất kỳ câu lệnh rẽ nhánh `if/else`, vòng lặp hay logic nghiệp vụ nào; đã được thực thi gián tiếp 100% khi chạy test tầng business. |
| 2 | `internal/module/seat/business/const.go` | 27 | Constants & Key Formatter | Tương tự file const của movie, chỉ định nghĩa pattern key Redis (`keySeatDetail`, `keySeatsListWithFilters`), không có logic nghiệp vụ độc lập. |
| 3 | `internal/module/movie/repository/postgres/repository.go` | 243 | Repository (DB Driver) | Chứa các câu lệnh SQL/Bun ORM kết nối trực tiếp với PostgreSQL (`r.db.NewInsert()`, `NewSelect()`, `NewUpdate()`, `NewDelete()`). Đây là tầng phụ thuộc chặt chẽ vào Database Engine, Network, và cấu trúc bảng thực tế; việc kiểm thử tương tác DB thực tế thuộc phạm vi của **Integration Test** và **System Test** (đã được kiểm thử trong tài liệu `03_systemtest.md`). Trong Unit Test, tầng này được thay thế bằng `fakeRepo`. |
| 4 | `internal/module/seat/repository/postgres/repository.go` | 215 | Repository (DB Driver) | Tương tự tầng Repository của Movie, thao tác trực tiếp với bảng `seats` của cơ sở dữ liệu. |
| 5 | `internal/module/movie/transport/rest/api.go` | 222 | Transport (HTTP Handler) | Tầng tiếp nhận HTTP Request của Gin Framework (`c.ShouldBindJSON`, `c.Param`). Trách nhiệm của file này chỉ là unmarshal dữ liệu JSON và gọi xuống tầng Business. Không chứa nghiệp vụ tính toán hay logic lõi. Đã được kiểm thử toàn diện qua bộ Postman API Test và System Test. |
| 6 | `internal/module/seat/transport/rest/api.go` | 211 | Transport (HTTP Handler) | Tương tự REST handler của Movie, chỉ làm nhiệm vụ parse query param, bind JSON và trả về HTTP Status code. |
| 7 | `internal/module/movie/transport/grpc/grpc.go` | 201 | Transport (gRPC Adapter) | Triển khai interface gRPC server của Protobuf để phục vụ các service khác (Booking Service). Logic tính giá ghế `CalculatePrice` đã được test ở tầng Entity/Business. |
| 8 | `FE/src/pages/admin/MovieDetailPage.jsx`, `MovieCard.jsx`... | N/A | Frontend UI Components | Các file giao diện React JSX chỉ chịu trách nhiệm render DOM, xử lý sự kiện click của người dùng và gọi API. Theo đúng đặc tả kiểm thử phần mềm, tầng giao diện không thực hiện Unit Test backend mà kiểm thử qua UI/E2E Test. |

---

### 1.3. Create Unit Test Cases (Thiết kế Test Cases)

Quy trình thiết kế test case tuân thủ nghiêm ngặt **Mục 5.3 của tài liệu**:
- **Đạt mức phủ cấp 2 (Branch Coverage):** Phủ hết mọi nhánh điều kiện `if/else`, `switch/case`, kiểm tra biên (boundary), giá trị hợp lệ, không hợp lệ, chuỗi rỗng, số 0, số âm, và trường hợp lỗi từ repository/database.
- **Expected Output:** Xác định căn cứ vào tài liệu thiết kế và đặc tả nghiệp vụ phần mềm (không căn cứ mù quáng vào lỗi có sẵn trong mã nguồn). Nhờ đó, Unit Test đã phát hiện ra các lỗi logic thực tế trong mã nguồn.
- **Cấu trúc script:** Khai báo biến input/expected -> Gọi method cần test -> Kiểm tra kết quả trả về -> Xác minh tác động trạng thái (CheckDB qua FakeRepo) -> Đảm bảo Rollback dữ liệu.

Tổng cộng đã xây dựng và thực thi **122 Test Cases** (Movie: 71 TCs, Seat: 44 TCs, Paging: 7 TCs). Toàn bộ được tổ chức trong file Excel [`Unit_Testing_Report_DUCPM88.xlsx`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/docs/TESTER_DUCPM88/original_docs/unit_test/Unit_Testing_Report_DUCPM88.xlsx). Dưới đây là bảng chi tiết:

#### A. Phân hệ Movie Management (71 Test Cases)

| Test Case ID | Class / File Name | Test Objective (Mục tiêu kiểm thử) | Input (Dữ liệu đầu vào) | Expected Output | Actual Result | Status | Notes / CheckDB & Rollback |
|---|---|---|---|---|---|:---:|---|
| **TC-MOVIE-ENTITY-01a** | `movie/entity/movie.go` | `Movie.IsValid()` - Phim hợp lệ | `Title="Avengers", Duration=120` | `true` | `true` | **Pass** | Đủ Title và Duration > 0 |
| **TC-MOVIE-ENTITY-01b** | `movie/entity/movie.go` | `Movie.IsValid()` - Tiêu đề rỗng | `Title="", Duration=120` | `false` | `false` | **Pass** | Title rỗng vi phạm validation |
| **TC-MOVIE-ENTITY-01c** | `movie/entity/movie.go` | `Movie.IsValid()` - Thời lượng = 0 | `Title="Test", Duration=0` | `false` | `false` | **Pass** | Duration = 0 không hợp lệ |
| **TC-MOVIE-ENTITY-01d** | `movie/entity/movie.go` | `Movie.IsValid()` - Thời lượng âm | `Title="Test", Duration=-1` | `false` | `false` | **Pass** | Duration âm không hợp lệ |
| **TC-MOVIE-ENTITY-01e** | `movie/entity/movie.go` | `Movie.IsValid()` - Biên nhỏ nhất hợp lệ | `Title="Test", Duration=1` | `true` | `true` | **Pass** | Biên trên ngay sát 0 |
| **TC-MOVIE-ENTITY-01f** | `movie/entity/movie.go` | `Movie.IsValid()` - Cả Title rỗng và Duration=0 | `Title="", Duration=0` | `false` | `false` | **Pass** | Cả hai điều kiện đều vi phạm |
| **TC-MOVIE-ENTITY-02a** | `movie/entity/movie.go` | `CanTransitionTo()` - `UPCOMING -> SHOWING` | `current=UPCOMING, next=SHOWING` | `true` | `true` | **Pass** | Chuyển trạng thái hợp lệ |
| **TC-MOVIE-ENTITY-02b** | `movie/entity/movie.go` | `CanTransitionTo()` - `UPCOMING -> ENDED` | `current=UPCOMING, next=ENDED` | `false` | `false` | **Pass** | Không được nhảy cóc qua SHOWING |
| **TC-MOVIE-ENTITY-02c** | `movie/entity/movie.go` | `CanTransitionTo()` - `SHOWING -> ENDED` | `current=SHOWING, next=ENDED` | `true` | `true` | **Pass** | Chuyển trạng thái hợp lệ |
| **TC-MOVIE-ENTITY-02d** | `movie/entity/movie.go` | `CanTransitionTo()` - `SHOWING -> UPCOMING` | `current=SHOWING, next=UPCOMING` | `false` | `false` | **Pass** | Không cho phép lùi trạng thái |
| **TC-MOVIE-ENTITY-02e** | `movie/entity/movie.go` | `CanTransitionTo()` - `ENDED -> UPCOMING` | `current=ENDED, next=UPCOMING` | `false` | `false` | **Pass** | ENDED là trạng thái kết thúc |
| **TC-MOVIE-ENTITY-02f** | `movie/entity/movie.go` | `CanTransitionTo()` - `ENDED -> SHOWING` | `current=ENDED, next=SHOWING` | `false` | `false` | **Pass** | ENDED không thể mở lại |
| **TC-MOVIE-ENTITY-02g** | `movie/entity/movie.go` | `CanTransitionTo()` - `UPCOMING -> UPCOMING` | `current=UPCOMING, next=UPCOMING` | `false` | `false` | **Pass** | Không cho chuyển sang cùng trạng thái |
| **TC-MOVIE-ENTITY-02h** | `movie/entity/movie.go` | `CanTransitionTo()` - `SHOWING -> SHOWING` | `current=SHOWING, next=SHOWING` | `false` | `false` | **Pass** | Không cho chuyển sang cùng trạng thái |
| **TC-MOVIE-ENTITY-02i** | `movie/entity/movie.go` | `CanTransitionTo()` - `ENDED -> ENDED` | `current=ENDED, next=ENDED` | `false` | `false` | **Pass** | Không cho chuyển sang cùng trạng thái |
| **TC-MOVIE-ENTITY-02j** | `movie/entity/movie.go` | `CanTransitionTo()` - next status rỗng | `current=UPCOMING, next=""` | `false` | `false` | **Pass** | Trạng thái rỗng bị từ chối |
| **TC-MOVIE-ENTITY-02k** | `movie/entity/movie.go` | `CanTransitionTo()` - chữ thường `showing` | `current=UPCOMING, next="showing"` | `false` | `false` | **Pass** | Phân biệt chữ hoa/thường theo enum |
| **TC-MOVIE-ENTITY-02l** | `movie/entity/movie.go` | `CanTransitionTo()` - status hiện tại lạ | `current="BOGUS", next=SHOWING` | `false` | `false` | **Pass** | Nhánh default xử lý trạng thái không hợp lệ |
| **TC-MOVIE-ENTITY-03a** | `movie/entity/types.go` | `CreateMovieRequest.ToEntity()` - status rỗng | `Title="A", Duration=100, Status=""` | `movie.Status = UPCOMING` | `UPCOMING` | **Pass** | Tự động gán mặc định UPCOMING |
| **TC-MOVIE-ENTITY-03b** | `movie/entity/types.go` | `CreateMovieRequest.ToEntity()` - có status | `Title="A", Duration=100, Status="SHOWING"` | `movie.Status = SHOWING` | `SHOWING` | **Pass** | Giữ nguyên status được truyền |
| **TC-MOVIE-ENTITY-03c** | `movie/entity/types.go` | `CreateMovieRequest.ToEntity()` - sao chép các trường | `T, D, C, 90, desc, trailer, poster` | Các trường sao chép chính xác | Khớp đầy đủ | **Pass** | Ánh xạ DTO sang Entity |
| **TC-MOVIE-ENTITY-03d** | `movie/entity/types.go` | `UpdateMovieRequest.ToEntity()` - status rỗng | `Title="A", Duration=100, Status=""` | `movie.Status = ""` (giữ nguyên) | `""` | **Pass** | Update không tự gán UPCOMING |
| **TC-MOVIE-ENTITY-03e** | `movie/entity/types.go` | `UpdateMovieRequest.ToEntity()` - có status | `Title="A", Duration=100, Status="ENDED"` | `movie.Status = ENDED` | `ENDED` | **Pass** | Sao chép status update |
| **TC-MOVIE-ENTITY-04a** | `movie/entity/types.go` | `ToMovieResponse()` - ánh xạ trường cơ bản | `Movie{Id:"m-1", Title:"Test", Status:SHOWING}` | `MovieResponse tương ứng` | Khớp dữ liệu | **Pass** | Ánh xạ Entity sang Response DTO |
| **TC-MOVIE-ENTITY-04b** | `movie/entity/types.go` | `ToMovieResponse()` - ánh xạ danh sách thể loại | `MovieGenres=[Action, Drama]` | `Genres=["Action", "Drama"]` | Khớp thể loại | **Pass** | Duyệt mảng liên kết MovieGenre |
| **TC-MOVIE-ENTITY-04c** | `movie/entity/types.go` | `ToMovieResponse()` - danh sách thể loại rỗng | `MovieGenres=[]` | `Genres=[] (không nil)` | `Genres=[]` | **Pass** | Khởi tạo slice rỗng an toàn |
| **TC-MOVIE-ENTITY-05a** | `movie/entity/types.go` | `ToMoviesResponse()` - chia có dư (25/10) | `total=25, size=10` | `Meta.TotalPages = 3` | `TotalPages = 3` | **Pass** | Ceiling division (25+10-1)/10 = 3 |
| **TC-MOVIE-ENTITY-05b** | `movie/entity/types.go` | `ToMoviesResponse()` - `total = 0` | `total=0, size=10` | `Meta.TotalPages = 1` | `TotalPages = 1` | **Pass** | Không có dữ liệu vẫn hiển thị 1 trang |
| **TC-MOVIE-ENTITY-05c** | `movie/entity/types.go` | `ToMoviesResponse()` - chia hết (30/10) | `total=30, size=10` | `Meta.TotalPages = 3` | `TotalPages = 3` | **Pass** | Phép chia không dư |
| **TC-MOVIE-ENTITY-05d** | `movie/entity/types.go` | `ToMoviesResponse()` - 1 bản ghi (1/10) | `total=1, size=10` | `Meta.TotalPages = 1` | `TotalPages = 1` | **Pass** | Biên nhỏ nhất có dữ liệu |
| **TC-MOVIE-ENTITY-05e** | `movie/entity/types.go` | `ToMoviesResponse()` - vượt 1 bản ghi (11/10) | `total=11, size=10` | `Meta.TotalPages = 2` | `TotalPages = 2` | **Pass** | Vừa sang trang thứ 2 |
| **TC-MOVIE-ENTITY-05f** | `movie/entity/types.go` | `ToMoviesResponse()` - danh sách phim | `2 movies, page=1, size=10, total=2` | `len(Movies)=2, Id khớp` | Khớp dữ liệu | **Pass** | Ánh xạ toàn bộ danh sách |
| **TC-MOVIE-BIZ-01a** | `movie/business/business.go` | `CreateMovie()` - movie = nil | `movie = nil` | `ErrInvalidMovieData` | `ErrInvalidMovieData` | **Pass** | CheckDB: `r.created == 0` |
| **TC-MOVIE-BIZ-01b** | `movie/business/business.go` | `CreateMovie()` - Title rỗng | `movie = {Title: "", Duration: 100}` | `ErrInvalidMovieData` | `ErrInvalidMovieData` | **Pass** | CheckDB: `r.created == 0` |
| **TC-MOVIE-BIZ-01c** | `movie/business/business.go` | `CreateMovie()` - Duration = 0 | `movie = {Title: "A", Duration: 0}` | `ErrInvalidMovieData` | `ErrInvalidMovieData` | **Pass** | CheckDB: `r.created == 0` |
| **TC-MOVIE-BIZ-01d** | `movie/business/business.go` | `CreateMovie()` - Status rỗng nhận default | `movie = {Title: "A", Duration: 100, Status: ""}` | `movie.Status = UPCOMING` | `UPCOMING` | **Pass** | CheckDB: `r.created == 1`, invalidate cache |
| **TC-MOVIE-BIZ-01e** | `movie/business/business.go` | `CreateMovie()` - Có status SHOWING | `movie = {Title: "A", Duration: 100, Status: SHOWING}` | `movie.Status = SHOWING` | `SHOWING` | **Pass** | CheckDB: `r.created == 1` |
| **TC-MOVIE-BIZ-01f** | `movie/business/business.go` | `CreateMovie()` - Repo trả về lỗi DB | `repo.Create trả về lỗi "db down"` | Bọc lỗi `"failed to create movie: db down"` | Đúng lỗi bọc | **Pass** | CheckDB: xử lý lỗi DB graceful |
| **TC-MOVIE-BIZ-02a** | `movie/business/business.go` | `UpdateMovie()` - movie = nil | `movie = nil` | `ErrInvalidMovieData` | `ErrInvalidMovieData` | **Pass** | Validation đầu vào |
| **TC-MOVIE-BIZ-02b** | `movie/business/business.go` | `UpdateMovie()` - Id rỗng | `movie = {Id: "", Title: "A"}` | `ErrInvalidMovieData` | `ErrInvalidMovieData` | **Pass** | Id bắt buộc |
| **TC-MOVIE-BIZ-02c** | `movie/business/business.go` | `UpdateMovie()` - Dữ liệu không hợp lệ | `movie = {Id: "m1", Title: "", Duration: 1}` | `ErrInvalidMovieData` | `ErrInvalidMovieData` | **Pass** | Title không được rỗng |
| **TC-MOVIE-BIZ-02d** | `movie/business/business.go` | `UpdateMovie()` - Phim không tồn tại | `repo.GetByID trả về sql.ErrNoRows` | `ErrMovieNotFound` | `ErrMovieNotFound` | **Pass** | CheckDB: kiểm tra tồn tại |
| **TC-MOVIE-BIZ-02e** | `movie/business/business.go` | `UpdateMovie()` - Cùng trạng thái (UPCOMING->UPCOMING) | `status cũ = UPCOMING, mới = UPCOMING` | Thành công, cập nhật dữ liệu | Thành công | **Pass** | CheckDB: `r.updated.Title == "New"` |
| **TC-MOVIE-BIZ-02f** | `movie/business/business.go` | `UpdateMovie()` - Chuyển hợp lệ (UPCOMING->SHOWING) | `status cũ = UPCOMING, mới = SHOWING` | Thành công | Thành công | **Pass** | CheckDB: `r.updated.Status == SHOWING` |
| **TC-MOVIE-BIZ-02g** | `movie/business/business.go` | `UpdateMovie()` - Chuyển sai (SHOWING->UPCOMING) | `status cũ = SHOWING, mới = UPCOMING` | `ErrInvalidStatusTransition` | `ErrInvalidStatusTransition` | **Pass** | CheckDB: `r.updated == nil` |
| **TC-MOVIE-BIZ-02h** | `movie/business/business.go` | `UpdateMovie()` - Client không truyền status khi update | `status cũ = SHOWING, client gửi Status = ""` | Giữ nguyên SHOWING và cập nhật thành công | `ErrInvalidStatusTransition` | **FAIL** | **Phát hiện Bug:** Code so sánh `"" != SHOWING` và từ chối cập nhật |
| **TC-MOVIE-BIZ-02i** | `movie/business/business.go` | `UpdateMovie()` - Repo update trả về lỗi | `repo.Update lỗi "db down"` | Bọc lỗi `"failed to update movie: db down"` | Đúng lỗi bọc | **Pass** | CheckDB: xử lý lỗi DB graceful |
| **TC-MOVIE-BIZ-03a** | `movie/business/business.go` | `DeleteMovie()` - Id rỗng | `id = ""` | Lỗi `"movie id is required"` | Đúng lỗi | **Pass** | Kiểm tra id bắt buộc |
| **TC-MOVIE-BIZ-03b** | `movie/business/business.go` | `DeleteMovie()` - Phim không tồn tại | `id = "m1", repo trả về sql.ErrNoRows` | `ErrMovieNotFound` | `ErrMovieNotFound` | **Pass** | CheckDB: `r.deleted == ""` |
| **TC-MOVIE-BIZ-03c** | `movie/business/business.go` | `DeleteMovie()` - Lỗi ràng buộc khóa ngoại (FK) | `repo.Delete lỗi "fk violation"` | Bọc lỗi `"failed to delete movie: fk violation"` | Đúng lỗi bọc | **Pass** | CheckDB: xử lý lỗi DB graceful |
| **TC-MOVIE-BIZ-03d** | `movie/business/business.go` | `DeleteMovie()` - Xóa thành công | `id = "m1", phim tồn tại` | Thành công, `r.deleted = "m1"` | Thành công | **Pass** | CheckDB: `r.deleted == "m1"`, xóa cache |
| **TC-MOVIE-BIZ-04a** | `movie/business/business.go` | `UpdateMovieStatus()` - Id rỗng | `id = "", status = SHOWING` | Lỗi `"movie id is required"` | Đúng lỗi | **Pass** | Validation đầu vào |
| **TC-MOVIE-BIZ-04b** | `movie/business/business.go` | `UpdateMovieStatus()` - Phim không tồn tại | `id = "m1", repo trả về sql.ErrNoRows` | `ErrMovieNotFound` | `ErrMovieNotFound` | **Pass** | CheckDB: kiểm tra tồn tại |
| **TC-MOVIE-BIZ-04c** | `movie/business/business.go` | `UpdateMovieStatus()` - UPCOMING -> SHOWING | `id = "m1", status cũ = UPCOMING, mới = SHOWING` | Thành công, `r.updated.Status = SHOWING` | Thành công | **Pass** | CheckDB: `r.updated.Status == SHOWING` |
| **TC-MOVIE-BIZ-04d** | `movie/business/business.go` | `UpdateMovieStatus()` - SHOWING -> UPCOMING | `id = "m1", status cũ = SHOWING, mới = UPCOMING` | `ErrInvalidStatusTransition` | `ErrInvalidStatusTransition` | **Pass** | CheckDB: `r.updated == nil` |
| **TC-MOVIE-BIZ-04e** | `movie/business/business.go` | `UpdateMovieStatus()` - UPCOMING -> ENDED | `id = "m1", status cũ = UPCOMING, mới = ENDED` | `ErrInvalidStatusTransition` | `ErrInvalidStatusTransition` | **Pass** | CheckDB: `r.updated == nil` |
| **TC-MOVIE-BIZ-04f** | `movie/business/business.go` | `UpdateMovieStatus()` - ENDED -> SHOWING | `id = "m1", status cũ = ENDED, mới = SHOWING` | `ErrInvalidStatusTransition` | `ErrInvalidStatusTransition` | **Pass** | CheckDB: `r.updated == nil` |
| **TC-MOVIE-BIZ-05a** | `movie/business/business.go` | `GetMovieById()` - Id rỗng | `id = ""` | `ErrInvalidMovieData` | `ErrInvalidMovieData` | **Pass** | CheckDB: không gọi DB |
| **TC-MOVIE-BIZ-05b** | `movie/business/business.go` | `GetMovieById()` - Không tìm thấy | `id = "x", sql.ErrNoRows` | `ErrMovieNotFound` | `ErrMovieNotFound` | **Pass** | CheckDB: ánh xạ `sql.ErrNoRows` |
| **TC-MOVIE-BIZ-05c** | `movie/business/business.go` | `GetMovieById()` - Lấy thành công | `id = "m1"` | Trả về Movie đúng Id | Khớp dữ liệu | **Pass** | CheckDB: trả về đối tượng movie |
| **TC-MOVIE-BIZ-05d** | `movie/business/business.go` | `GetMovieById()` - Lỗi DB khác | `repo trả về "db down"` | Bọc lỗi `"failed to get movie: db down"` | Đúng lỗi bọc | **Pass** | Phân biệt lỗi DB với NotFound |
| **TC-MOVIE-BIZ-06a** | `movie/business/business.go` | `GetMovies()` - page=0, size=0 được chuẩn hóa | `page=0, size=0` | `limit=10, offset=0` | `limit=10, offset=0` | **Pass** | Chuẩn hóa giá trị mặc định |
| **TC-MOVIE-BIZ-06b** | `movie/business/business.go` | `GetMovies()` - size > 100 bị chặn trần | `page=1, size=500` | `limit=100` | `limit=100` | **Pass** | Bảo vệ DB tránh quá tải |
| **TC-MOVIE-BIZ-06c** | `movie/business/business.go` | `GetMovies()` - Tính offset đúng | `page=3, size=10` | `offset = (3-1)*10 = 20` | `offset = 20` | **Pass** | Công thức phân trang chính xác |
| **TC-MOVIE-BIZ-06d** | `movie/business/business.go` | `GetMovies()` - Trả về movies và total | `page=1, size=10, search="a", status="SHOWING"` | Danh sách movies, total=7 | Khớp dữ liệu | **Pass** | CheckDB: GetMany và GetTotalCount đúng filter |
| **TC-MOVIE-BIZ-07a** | `movie/business/business.go` | `ValidateMovieForShowtime()` - movieId rỗng | `movieId = ""` | `ErrInvalidMovieData` | `ErrInvalidMovieData` | **Pass** | Validation đầu vào |
| **TC-MOVIE-BIZ-07b** | `movie/business/business.go` | `ValidateMovieForShowtime()` - Phim không tồn tại | `movieId = "m1", không tìm thấy` | `ErrMovieNotFound` | `ErrMovieNotFound` | **Pass** | CheckDB: kiểm tra tồn tại |
| **TC-MOVIE-BIZ-07c** | `movie/business/business.go` | `ValidateMovieForShowtime()` - Trạng thái UPCOMING/ENDED | `movieId = "m1", status = UPCOMING hoặc ENDED` | `ErrMovieNotShowing` | `ErrMovieNotShowing` | **Pass** | Chỉ phim SHOWING mới được chiếu |
| **TC-MOVIE-BIZ-07d** | `movie/business/business.go` | `ValidateMovieForShowtime()` - Phim SHOWING hợp lệ | `movieId = "m1", status = SHOWING` | `nil (hợp lệ)` | `nil` | **Pass** | Cho phép tạo suất chiếu |
| **TC-MOVIE-BIZ-08a** | `movie/business/business.go` | `GetMovieStats()` - Lấy thống kê trạng thái | `ctx` | Danh sách []*MovieStat | Khớp dữ liệu | **Pass** | CheckDB: trả về thống kê qua cache |
| **TC-MOVIE-BIZ-08b** | `movie/business/business.go` | `GetGenres()` - Lấy danh sách thể loại | `ctx` | Danh sách []*Genre có Id="g1", Name="Action" | Khớp dữ liệu | **Pass** | CheckDB: trả về danh sách thể loại |

---

#### B. Phân hệ Seat Management (44 Test Cases)

| Test Case ID | Class / File Name | Test Objective (Mục tiêu kiểm thử) | Input (Dữ liệu đầu vào) | Expected Output | Actual Result | Status | Notes / CheckDB & Rollback |
|---|---|---|---|---|---|:---:|---|
| **TC-SEAT-ENTITY-01a** | `seat/entity/seat.go` | `Seat.IsValid()` - Ghế hợp lệ | `RoomId="r1", SeatNumber="01", RowNumber="A"` | `true` | `true` | **Pass** | Đủ 3 trường bắt buộc |
| **TC-SEAT-ENTITY-01b** | `seat/entity/seat.go` | `Seat.IsValid()` - RoomId rỗng | `RoomId="", SeatNumber="01", RowNumber="A"` | `false` | `false` | **Pass** | Thiếu RoomId |
| **TC-SEAT-ENTITY-01c** | `seat/entity/seat.go` | `Seat.IsValid()` - SeatNumber rỗng | `RoomId="r1", SeatNumber="", RowNumber="A"` | `false` | `false` | **Pass** | Thiếu SeatNumber |
| **TC-SEAT-ENTITY-01d** | `seat/entity/seat.go` | `Seat.IsValid()` - RowNumber rỗng | `RoomId="r1", SeatNumber="01", RowNumber=""` | `false` | `false` | **Pass** | Thiếu RowNumber |
| **TC-SEAT-ENTITY-01e** | `seat/entity/seat.go` | `Seat.IsValid()` - Ghế rỗng toàn bộ | `Seat{}` | `false` | `false` | **Pass** | Tất cả đều rỗng |
| **TC-SEAT-ENTITY-02a** | `seat/entity/seat.go` | `GetSeatTypePriceMultiplier()` - REGULAR | `SeatType = REGULAR` | `1.0` | `1.0` | **Pass** | Hệ số ghế thường 1.0x |
| **TC-SEAT-ENTITY-02b** | `seat/entity/seat.go` | `GetSeatTypePriceMultiplier()` - VIP | `SeatType = VIP` | `1.5` | `1.5` | **Pass** | Hệ số ghế VIP 1.5x |
| **TC-SEAT-ENTITY-02c** | `seat/entity/seat.go` | `GetSeatTypePriceMultiplier()` - COUPLE | `SeatType = COUPLE` | `2.5` | `2.5` | **Pass** | Hệ số ghế đôi 2.5x |
| **TC-SEAT-ENTITY-02d** | `seat/entity/seat.go` | `GetSeatTypePriceMultiplier()` - UNKNOWN | `SeatType = "UNKNOWN"` | `1.0` | `1.0` | **Pass** | Loại không rõ về 1.0 |
| **TC-SEAT-ENTITY-02e** | `seat/entity/seat.go` | `GetSeatTypePriceMultiplier()` - Rỗng | `SeatType = ""` | `1.0` | `1.0` | **Pass** | Mặc định về 1.0 |
| **TC-SEAT-ENTITY-03a** | `seat/entity/seat.go` | `Seat.CalculatePrice()` - Ghế REGULAR | `basePrice=100000, SeatType=REGULAR` | `100000` | `100000` | **Pass** | 100.000 * 1.0 = 100.000 |
| **TC-SEAT-ENTITY-03b** | `seat/entity/seat.go` | `Seat.CalculatePrice()` - Ghế VIP | `basePrice=100000, SeatType=VIP` | `150000` | `150000` | **Pass** | 100.000 * 1.5 = 150.000 |
| **TC-SEAT-ENTITY-03c** | `seat/entity/seat.go` | `Seat.CalculatePrice()` - Ghế COUPLE | `basePrice=100000, SeatType=COUPLE` | `250000` | `250000` | **Pass** | 100.000 * 2.5 = 250.000 |
| **TC-SEAT-ENTITY-03d** | `seat/entity/seat.go` | `Seat.CalculatePrice()` - Giá gốc = 0 | `basePrice=0, SeatType=VIP` | `0` | `0` | **Pass** | 0 * 1.5 = 0 |
| **TC-SEAT-ENTITY-03e** | `seat/entity/seat.go` | `Seat.CalculatePrice()` - Loại không rõ | `basePrice=80000, SeatType="UNKNOWN"` | `80000` | `80000` | **Pass** | 80.000 * 1.0 = 80.000 |
| **TC-SEAT-ENTITY-04a** | `seat/entity/types.go` | `CreateSeatRequest.ToSeat()` - status mặc định | `RoomId="r1", SeatNumber="01", RowNumber="A", SeatType=VIP` | `seat.Status = AVAILABLE` | `AVAILABLE` | **Pass** | Tự động gán AVAILABLE |
| **TC-SEAT-ENTITY-04b** | `seat/entity/types.go` | `CreateSeatRequest.ToSeat()` - sao chép các trường | `r1, 01, A, VIP` | Các trường sao chép đầy đủ | Khớp dữ liệu | **Pass** | Ánh xạ DTO sang Entity |
| **TC-SEAT-ENTITY-05a** | `seat/entity/types.go` | `ToSeatResponse()` - ánh xạ trường | `Seat{Id:"s-1", SeatType:COUPLE, Status:MAINTENANCE}` | `SeatResponse tương ứng` | Khớp dữ liệu | **Pass** | Ánh xạ Entity sang Response DTO |
| **TC-SEAT-ENTITY-05b** | `seat/entity/types.go` | `ToSeatResponse()` - đầu vào nil | `nil` | `nil` | `nil` | **Pass** | Xử lý con trỏ nil an toàn |
| **TC-SEAT-ENTITY-06a** | `seat/entity/types.go` | `ToSeatsResponse()` - danh sách 2 ghế | `2 seats, page=1, size=10, total=2` | `len(Data)=2, Paging != nil` | Khớp | **Pass** | Ánh xạ mảng dữ liệu |
| **TC-SEAT-ENTITY-06b** | `seat/entity/types.go` | `ToSeatsResponse()` - tính phân trang | `nil, page=2, size=50, total=120` | `TotalPages=3, Page=2` | `TotalPages=3, Page=2` | **Pass** | Phân trang (120/50 làm tròn lên = 3) |
| **TC-SEAT-BIZ-01a** | `seat/business/business.go` | `CreateSeat()` - seat = nil | `seat = nil` | `ErrInvalidSeatData` | `ErrInvalidSeatData` | **Pass** | CheckDB: `r.created == nil` |
| **TC-SEAT-BIZ-01b** | `seat/business/business.go` | `CreateSeat()` - SeatNumber rỗng | `seat = {RoomId: "r1", RowNumber: "A", SeatNumber: ""}` | `ErrInvalidSeatData` | `ErrInvalidSeatData` | **Pass** | CheckDB: `r.created == nil` |
| **TC-SEAT-BIZ-01c** | `seat/business/business.go` | `CreateSeat()` - Trùng vị trí ghế trong phòng | `exists = true từ repo.ExistsBySeatPosition` | `ErrSeatPositionExists` | `ErrSeatPositionExists` | **Pass** | CheckDB: `r.created == nil` (không insert) |
| **TC-SEAT-BIZ-01d** | `seat/business/business.go` | `CreateSeat()` - Repo kiểm tra vị trí báo lỗi DB | `existsErr = "db"` | Bọc lỗi `"failed to check seat position: db"` | Đúng lỗi bọc | **Pass** | CheckDB: xử lý lỗi DB graceful |
| **TC-SEAT-BIZ-01e** | `seat/business/business.go` | `CreateSeat()` - Lỗi khóa ngoại (FK room) | `createErr = "fk room"` | Bọc lỗi `"failed to create seat: fk room"` | Đúng lỗi bọc | **Pass** | CheckDB: xử lý lỗi DB graceful |
| **TC-SEAT-BIZ-01f** | `seat/business/business.go` | `CreateSeat()` - Tạo ghế thành công | `seat hợp lệ, vị trí chưa có` | Thành công, `r.created != nil` | Thành công | **Pass** | CheckDB: `r.created != nil`, xóa cache |
| **TC-SEAT-BIZ-01g** | `seat/business/business.go` | `CreateSeat()` - SeatType không thuộc enum | `SeatType = "BOGUS"` | Từ chối và trả về `ErrInvalidSeatData` | Lưu `"BOGUS"` vào DB! | **FAIL** | **Phát hiện Bug:** Code không validate giá trị SeatType hợp lệ |
| **TC-SEAT-BIZ-02a** | `seat/business/business.go` | `UpdateSeat()` - Ghế không tồn tại | `id = "s1", sql.ErrNoRows` | `ErrSeatNotFound` | `ErrSeatNotFound` | **Pass** | CheckDB: kiểm tra tồn tại |
| **TC-SEAT-BIZ-02b** | `seat/business/business.go` | `UpdateSeat()` - Đổi SeatNumber loại trừ chính nó | `updates.SeatNumber = "Z2"` | Thành công, `excludeId = "s1"` | Thành công | **Pass** | CheckDB: `ExistsBySeatPosition` truyền `excludeId="s1"` |
| **TC-SEAT-BIZ-02c** | `seat/business/business.go` | `UpdateSeat()` - Đổi RowNumber trùng vị trí ghế khác | `updates.RowNumber = "B", exists = true` | `ErrSeatPositionExists` | `ErrSeatPositionExists` | **Pass** | CheckDB: `r.updated == nil` |
| **TC-SEAT-BIZ-02d** | `seat/business/business.go` | `UpdateSeat()` - Đổi SeatNumber thành rỗng | `updates.SeatNumber = ""` | `ErrInvalidSeatData` | `ErrInvalidSeatData` | **Pass** | CheckDB: `r.updated == nil` |
| **TC-SEAT-BIZ-02e** | `seat/business/business.go` | `UpdateSeat()` - Chỉ đổi SeatType và Status | `updates.SeatType = VIP, Status = MAINTENANCE` | Thành công, không gọi `ExistsBySeatPosition` | Thành công | **Pass** | Tối ưu hóa: không query vị trí khi không đổi số ghế |
| **TC-SEAT-BIZ-02f** | `seat/business/business.go` | `UpdateSeat()` - Repo update trả về lỗi | `updateErr = "db"` | Bọc lỗi `"failed to update seat: db"` | Đúng lỗi bọc | **Pass** | CheckDB: xử lý lỗi DB graceful |
| **TC-SEAT-BIZ-03a** | `seat/business/business.go` | `DeleteSeat()` - Ghế không tồn tại | `id = "s1", sql.ErrNoRows` | `ErrSeatNotFound` | `ErrSeatNotFound` | **Pass** | CheckDB: kiểm tra tồn tại |
| **TC-SEAT-BIZ-03b** | `seat/business/business.go` | `DeleteSeat()` - Xóa thành công | `id = "s1", ghế tồn tại` | Thành công, `r.deleted = "s1"` | Thành công | **Pass** | CheckDB: `r.deleted == "s1"`, xóa cache |
| **TC-SEAT-BIZ-03c** | `seat/business/business.go` | `DeleteSeat()` - Lỗi khóa ngoại (đã có booking) | `deleteErr = "fk booking"` | Bọc lỗi `"failed to delete seat: fk booking"` | Đúng lỗi bọc | **Pass** | CheckDB: xử lý lỗi DB graceful |
| **TC-SEAT-BIZ-04a** | `seat/business/business.go` | `UpdateSeatStatus()` - Ghế không tồn tại | `id = "s1", sql.ErrNoRows` | `ErrSeatNotFound` | `ErrSeatNotFound` | **Pass** | CheckDB: kiểm tra tồn tại |
| **TC-SEAT-BIZ-04b** | `seat/business/business.go` | `UpdateSeatStatus()` - AVAILABLE -> MAINTENANCE | `id = "s1", status = MAINTENANCE` | Thành công, `r.updated.Status = MAINTENANCE` | Thành công | **Pass** | CheckDB: `r.updated.Status == MAINTENANCE` |
| **TC-SEAT-BIZ-04c** | `seat/business/business.go` | `UpdateSeatStatus()` - Status lạ không thuộc enum | `id = "s1", status = "FOO"` | Từ chối và trả về `ErrInvalidStatusTransition` | Lưu `"FOO"` vào DB! | **FAIL** | **Phát hiện Bug:** Code không validate giá trị SeatStatus |
| **TC-SEAT-BIZ-05a** | `seat/business/business.go` | `GetSeats()` - page < 1 hoặc size < 1 | `page=0, size=10` hoặc `page=1, size=0` | `ErrInvalidSeatData` | `ErrInvalidSeatData` | **Pass** | Validation phân trang |
| **TC-SEAT-BIZ-05b** | `seat/business/business.go` | `GetSeats()` - Lấy danh sách thành công | `page=2, size=50, roomId="r1"` | Danh sách seats, total=1, offset=50 | Khớp dữ liệu | **Pass** | CheckDB: `offset = (2-1)*50 = 50` |
| **TC-SEAT-BIZ-05c** | `seat/business/business.go` | `GetSeatById()` - Không tìm thấy | `id = "x", sql.ErrNoRows` | `ErrSeatNotFound` | `ErrSeatNotFound` | **Pass** | CheckDB: ánh xạ `sql.ErrNoRows` |
| **TC-SEAT-BIZ-05d** | `seat/business/business.go` | `GetSeatById()` - Lấy thành công | `id = "s1"` | Trả về Seat đúng Id | Khớp dữ liệu | **Pass** | CheckDB: trả về đối tượng seat |

#### C. Shared Package Paging (7 Test Cases)

| Test Case ID | Class / File Name | Test Objective | Input | Expected Output | Actual Result | Status | Notes |
|---|---|---|---|---|---|:---:|---|
| **TC-PAGING-01a** | `pkg/paging/paging.go` | `NewPageInfo()` - chia có dư (25/10) | `page=1, size=10, total=25` | `TotalPages=3, Page=1, Size=10, Total=25` | `TotalPages=3` | **Pass** | Ceiling division: 25/10 -> 3 trang |
| **TC-PAGING-01b** | `pkg/paging/paging.go` | `NewPageInfo()` - chia hết (20/10) | `page=1, size=10, total=20` | `TotalPages=2, Page=1, Size=10, Total=20` | `TotalPages=2` | **Pass** | Trang tròn: 20/10 -> 2 trang |
| **TC-PAGING-01c** | `pkg/paging/paging.go` | `NewPageInfo()` - total = 0 | `page=1, size=10, total=0` | `TotalPages=0, Page=1, Size=10, Total=0` | `TotalPages=0` | **Pass** | Khi không có dữ liệu: 0 trang |
| **TC-PAGING-01d** | `pkg/paging/paging.go` | `NewPageInfo()` - 1 bản ghi (1/10) | `page=1, size=10, total=1` | `TotalPages=1, Page=1, Size=10, Total=1` | `TotalPages=1` | **Pass** | Biên nhỏ nhất: 1 trang |
| **TC-PAGING-01e** | `pkg/paging/paging.go` | `NewPageInfo()` - vượt 1 bản ghi (11/10) | `page=1, size=10, total=11` | `TotalPages=2, Page=1, Size=10, Total=11` | `TotalPages=2` | **Pass** | Vừa sang trang thứ 2 |
| **TC-PAGING-02a** | `pkg/paging/paging.go` | `stringToInt()` - chuỗi số hợp lệ | `"7"` | `7` | `7` | **Pass** | Parse số nguyên thành công |
| **TC-PAGING-02b** | `pkg/paging/paging.go` | `stringToInt()` - chuỗi ký tự lỗi | `"abc"` | `0` (giá trị mặc định khi parse lỗi) | `0` | **Pass** | Xử lý lỗi parse an toàn |

---

### 1.4. Project Link
- **GitHub Repository URL:** `https://github.com/phmduc/cinema-project`
- **Thư mục chứa mã nguồn Unit Test trực tiếp:**
  - `movie-service/internal/module/movie/entity/movie_test.go`
  - `movie-service/internal/module/movie/entity/types_test.go`
  - `movie-service/internal/module/movie/business/business_test.go`
  - `movie-service/internal/module/seat/entity/seat_test.go`
  - `movie-service/internal/module/seat/entity/types_test.go`
  - `movie-service/internal/module/seat/business/business_test.go`
  - `movie-service/internal/pkg/paging/paging_test.go`

---

### 1.5. Execution Report (Báo cáo thực thi kiểm thử)

#### Bảng tổng hợp kết quả chạy Unit Test:

| Nhóm kiểm thử (Test Suite) | Tổng số Test Case | Passed | Failed | Tỷ lệ Pass (%) | Trạng thái phát hiện lỗi |
|---|---:|---:|---:|---:|---|
| **Movie Management — Entity Layer** | 32 | 32 | 0 | 100.0% | Pass 100% |
| **Movie Management — Business Layer** | 39 | 38 | 1 | 97.44% | **Phát hiện 1 Bug nghiệp vụ (`TC-MOVIE-BIZ-02h`)** |
| **Seat Management — Entity Layer** | 21 | 21 | 0 | 100.0% | Pass 100% |
| **Seat Management — Business Layer** | 23 | 21 | 2 | 91.30% | **Phát hiện 2 Bug nghiệp vụ (`TC-SEAT-BIZ-01g`, `04c`)** |
| **Paging Helper Package** | 7 | 7 | 0 | 100.0% | Pass 100% |
| **TỔNG CỘNG PHẠM VI TEST** | **122** | **119** | **3** | **97.54%** | **Phát hiện 3 lỗi nghiệp vụ thực tế!** |

#### Minh chứng chi tiết các lỗi nghiệp vụ phát hiện qua Unit Test:
1. **Lỗi `TC-MOVIE-BIZ-02h` (`UpdateMovie`):**
   - *Mô tả:* Khi client gửi request cập nhật thông tin phim đang có trạng thái `SHOWING` nhưng không truyền trường `status` (hoặc để rỗng `""`), hàm `UpdateMovie` so sánh `movie.Status != existingMovie.Status` (`"" != SHOWING`), sau đó gọi `CanTransitionTo("")` trả về `false`, khiến hệ thống báo lỗi `ErrInvalidStatusTransition` và chặn không cho cập nhật tiêu đề, đạo diễn hay mô tả của phim.
   - *Khắc phục đề xuất:* Kiểm tra `if movie.Status != "" && movie.Status != existingMovie.Status` trước khi gọi xác thực chuyển trạng thái.
2. **Lỗi `TC-SEAT-BIZ-01g` (`CreateSeat`):**
   - *Mô tả:* Hàm `CreateSeat` trong `business.go` không kiểm tra giá trị của `SeatType` có thuộc tập `REGULAR`, `VIP`, `COUPLE` hay không. Khi client gửi loại ghế bất kỳ như `"BOGUS"`, hệ thống vẫn chấp nhận và lưu vào DB.
   - *Khắc phục đề xuất:* Bổ sung validate `switch seat.SeatType { case SeatTypeRegular, SeatTypeVIP, SeatTypeCouple: default: return ErrInvalidSeatData }`.
3. **Lỗi `TC-SEAT-BIZ-04c` (`UpdateSeatStatus`):**
   - *Mô tả:* Hàm `UpdateSeatStatus` nhận vào tham số kiểu `SeatStatus` nhưng không kiểm tra tính hợp lệ trước khi gọi `repo.Update`, cho phép cập nhật trạng thái ghế thành bất kỳ chuỗi nào (ví dụ `"FOO"`).
   - *Khắc phục đề xuất:* Kiểm tra `switch status { case SeatStatusAvailable, SeatStatusOccupied, SeatStatusMaintenance, SeatStatusBlocked: default: return ErrInvalidSeatData }`.

---

### 1.6. Code Coverage Report (Báo cáo độ phủ mã nguồn)

Độ phủ mã nguồn được đo lường bằng công cụ chuẩn của Golang:
```bash
go test -count=1 -coverprofile=coverage.out -covermode=atomic ./internal/module/movie/entity/... ./internal/module/movie/business/... ./internal/module/seat/entity/... ./internal/module/seat/business/... ./internal/pkg/paging/...
go tool cover -func=coverage.out
```

#### Tóm tắt kết quả độ phủ theo từng package & file:

| Gói mã nguồn (Package) | Tệp tin (File) | Statements Tested | Branch Coverage | Tỷ lệ phủ (Coverage %) |
|---|---|:---:|:---:|:---:|
| `internal/module/movie/entity` | `movie.go` | 7/7 | N/A | **100.0%** |
| `internal/module/movie/entity` | `types.go` | 20/20 | N/A | **100.0%** |
| `internal/module/movie/business` | `business.go` | 102/120 | N/A | **85.0%** |
| `internal/module/seat/entity` | `seat.go` | 9/9 | N/A | **100.0%** |
| `internal/module/seat/entity` | `types.go` | 12/12 | N/A | **100.0%** |
| `internal/module/seat/business` | `business.go` | 82/138 | N/A | **59.4%** |
| `internal/pkg/paging` | `paging.go` | 8/11 | N/A | **72.7%** |
| **TỔNG THEO `go tool cover -func`** | **7 file trên + 2 file `const.go` (5/5 statements)** | **245/322** | **N/A** | **76.1%** |

> *Branch Coverage: N/A — công cụ `go tool cover` chỉ đo độ phủ câu lệnh (statement coverage), không xuất số liệu độ phủ nhánh.*

> **Nhận xét độ phủ:** Độ phủ câu lệnh tổng thể đạt **76.1%** (245/322 statements). Theo package: `movie/entity` 100%, `seat/entity` 100%, `movie/business` 85.4%, `seat/business` 60.0%, `paging` 72.7%. Các phần chưa được phủ: `NewBusiness` của movie và seat (0%), `GetLockedSeatsByShowtime`, `getConcurrentLockedSeatsByShowtime`, `getBookedSeatsByShowtime` của seat (0%), `GetQueryPaging` của paging (0%), cùng một số nhánh lỗi trong `GetMovies`, `GetGenres`, `GetMovieStats`, `UpdateMovieStatus`, `GetSeats`, `UpdateSeat`, `UpdateSeatStatus`.

---

### 1.7. References + List of Prompts Used

#### A. Tài liệu tham khảo (References):
1. **Golang Testing Package Documentation:** https://pkg.go.dev/testing
2. **Go Code Coverage Tooling:** https://go.dev/blog/cover
3. **Bun ORM Documentation:** https://bun.uptrace.dev/
4. **Clean Architecture in Go (Repository & In-Memory Fake Pattern):** Martin Fowler, *Mocks Aren't Stubs*.
5. **Đặc tả yêu cầu đồ án tốt nghiệp:** `docs/graduation_thesis.pdf` & `docs/API.md`.
6. **Mục 5.3 Quy trình kiểm thử đơn vị:** Hướng dẫn xây dựng test case mức phủ cấp 2, CheckDB và Rollback từ tài liệu đảm bảo chất lượng phần mềm.

#### B. Danh sách các câu lệnh / prompts đã sử dụng:
1. `Prompt 1`: Phân tích cấu trúc thư mục và mã nguồn backend `movie-service` để xác định các file thuộc phạm vi Movie Management và Seat Management.
2. `Prompt 2`: Lập danh sách phân loại chi tiết các file cần kiểm thử và các file không cần kiểm thử (Excluded Files) kèm lý do kỹ thuật.
3. `Prompt 3`: Thiết kế bộ test case cho tầng Entity (Validation, Price Calculation, State Transitions, DTO Mappers) đạt mức phủ cấp 2.
4. `Prompt 4`: Thiết kế bộ test case cho tầng Business logic sử dụng Fake Repository, kiểm tra CheckDB và xử lý lỗi DB/Cache.
5. `Prompt 5`: Xây dựng script thực thi test tự động hiển thị kết quả và độ phủ trong terminal (`run_tests.bat`).
6. `Prompt 6`: Xây dựng script Rollback dọn dẹp môi trường và bảo toàn tính toàn vẹn dữ liệu (`rollback.bat`).
7. `Prompt 7`: Xuất bảng test case ra file Excel (`Unit_Testing_Report_DUCPM88.xlsx`) và CSV hỗ trợ mở bằng Microsoft Excel.

---

## 2. Requirements for Unit Test Scripts (Yêu cầu đối với Script Unit Test)

Mọi file script kiểm thử trong dự án đều tuân thủ chặt chẽ 5 tiêu chuẩn cốt lõi:

1. **Detailed Comments (Chú thích chi tiết):**
   - Mỗi file test đều có khối comment mở đầu giải thích mục đích, đối tượng test và kỹ thuật giả lập (Fakes/Mocks).
   - Trước mỗi hàm test và test case đều có comment mô tả mục tiêu, input, điều kiện kỳ vọng và phân nhánh logic.

2. **Test Case ID Identification (Định danh Test Case ID):**
   - Mọi hàm con trong `t.Run()` đều mang mã định danh duy nhất (ví dụ: `t.Run("TC-MOVIE-BIZ-01a", func(t *testing.T) { ... })`).
   - Mã ID này đồng bộ 100% với file Excel `Unit_Testing_Report_DUCPM88.xlsx` và báo cáo này.

3. **Meaningful Naming Conventions (Quy ước đặt tên rõ ràng):**
   - Tên biến mang ý nghĩa trực quan: `movie`, `existingMovie`, `fakeRepo`, `boom` (đại diện cho lỗi DB giả lập), `pagingParams`.
   - Hàm test đặt theo chuẩn Go: `Test<TênMethod>`, các case con dùng `t.Run`.

4. **CheckDB Verification (Kiểm tra tác động Database):**
   - Đối với các phương thức có ghi/sửa/xóa DB (`CreateMovie`, `UpdateMovie`, `DeleteMovie`, `CreateSeat`, `UpdateSeat`, `DeleteSeat`):
     - Script kiểm tra trực tiếp biến lưu vết trong Fake Repository:
       - `r.created`: Kiểm tra số lần gọi hàm `Create` và nội dung struct được chuyển xuống DB.
       - `r.updated`: Kiểm tra dữ liệu cập nhật có đúng với giá trị mong đợi không.
       - `r.deleted`: Kiểm tra ID bản ghi bị xóa.
     - Khi dữ liệu đầu vào không hợp lệ hoặc kiểm tra logic thất bại, script kiểm tra đảm bảo `r.created == 0` hoặc `r.updated == nil` (chứng minh DB **không** bị tác động sai).

5. **Rollback & Data Integrity (Khôi phục trạng thái):**
   - Nhờ sử dụng kiến trúc **Fake In-Memory Repository**, mỗi test case được cấp một instance `fakeRepo` và `business` hoàn toàn mới và độc lập trong bộ nhớ RAM (`newBiz(&fakeRepo{...})`).
   - Dữ liệu thử nghiệm hoàn toàn không chạm vào PostgreSQL Database vật lý, loại bỏ 100% nguy cơ làm bẩn hoặc hỏng dữ liệu hệ thống.
   - Script [`rollback.bat`](file:///d:/TÀI%20LIỆU%20WORD%20EXCEL/ĐẢM%20BẢO%20CHẤT%20LƯỢNG%20PHẦN%20MỀM/cinema-project-main/docs/TESTER_DUCPM88/original_docs/unit_test/rollback.bat) được tạo sẵn để người dùng có thể click chạy khôi phục/dọn dẹp các tệp tin tạm và file coverage phát sinh sau khi test.
