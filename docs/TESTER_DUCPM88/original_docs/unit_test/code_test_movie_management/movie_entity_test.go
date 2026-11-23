// =============================================================================
// FILE: movie_entity_test.go
// PACKAGE: entity
// MÔ TẢ: Unit test cho các phương thức của Movie entity:
//   - Movie.IsValid()
//   - Movie.CanTransitionTo()
//   - CreateMovieRequest.ToEntity()
//   - UpdateMovieRequest.ToEntity()
//   - ToMovieResponse()
//   - ToMoviesResponse()
//
// TESTER: DUCPM88 - Phạm Minh Đức - B22DCDT095
// MỤC TIÊU: Đảm bảo phủ cấp 2 (branch coverage) cho tất cả nhánh logic
// LƯU Ý: Không có DB call trong test này - thuần túy kiểm tra logic entity
// =============================================================================

package entity

import (
	"testing"
	"time"
)

// =============================================================================
// TEST NHÓM 1: Movie.IsValid()
// Nguồn: movie-service/internal/module/movie/entity/movie.go
// Logic: trả về false nếu Title rỗng HOẶC Duration <= 0
// =============================================================================

// TC-MOVIE-ENTITY-01a
// Mục tiêu: IsValid() trả về true khi movie có Title và Duration hợp lệ
func TestMovieIsValid_ValidTitleAndDuration(t *testing.T) {
	// --- Khai báo input ---
	inputMovie := &Movie{
		Title:    "Avengers: Endgame",
		Duration: 182,
	}
	expectedResult := true

	// --- Gọi method cần test ---
	actualResult := inputMovie.IsValid()

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-MOVIE-ENTITY-01a FAILED: IsValid() = %v, want %v", actualResult, expectedResult)
	}
}

// TC-MOVIE-ENTITY-01b
// Mục tiêu: IsValid() trả về false khi Title rỗng (dù Duration hợp lệ)
func TestMovieIsValid_EmptyTitle(t *testing.T) {
	// --- Khai báo input ---
	inputMovie := &Movie{
		Title:    "", // Title bị bỏ trống
		Duration: 120,
	}
	expectedResult := false

	// --- Gọi method cần test ---
	actualResult := inputMovie.IsValid()

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-MOVIE-ENTITY-01b FAILED: IsValid() = %v, want %v", actualResult, expectedResult)
	}
}

// TC-MOVIE-ENTITY-01c
// Mục tiêu: IsValid() trả về false khi Duration = 0 (dù Title hợp lệ)
func TestMovieIsValid_ZeroDuration(t *testing.T) {
	// --- Khai báo input ---
	inputMovie := &Movie{
		Title:    "Test Movie",
		Duration: 0, // Duration = 0 không hợp lệ
	}
	expectedResult := false

	// --- Gọi method cần test ---
	actualResult := inputMovie.IsValid()

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-MOVIE-ENTITY-01c FAILED: IsValid() = %v, want %v", actualResult, expectedResult)
	}
}

// TC-MOVIE-ENTITY-01d
// Mục tiêu: IsValid() trả về false khi Duration âm
func TestMovieIsValid_NegativeDuration(t *testing.T) {
	// --- Khai báo input ---
	inputMovie := &Movie{
		Title:    "Test Movie",
		Duration: -10, // Duration âm không hợp lệ
	}
	expectedResult := false

	// --- Gọi method cần test ---
	actualResult := inputMovie.IsValid()

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-MOVIE-ENTITY-01d FAILED: IsValid() = %v, want %v", actualResult, expectedResult)
	}
}

// TC-MOVIE-ENTITY-01e
// Mục tiêu: IsValid() trả về false khi cả Title rỗng và Duration <= 0
func TestMovieIsValid_BothTitleEmptyAndZeroDuration(t *testing.T) {
	// --- Khai báo input ---
	inputMovie := &Movie{
		Title:    "",
		Duration: 0,
	}
	expectedResult := false

	// --- Gọi method cần test ---
	actualResult := inputMovie.IsValid()

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-MOVIE-ENTITY-01e FAILED: IsValid() = %v, want %v", actualResult, expectedResult)
	}
}

// =============================================================================
// TEST NHÓM 2: Movie.CanTransitionTo()
// Nguồn: movie-service/internal/module/movie/entity/movie.go
// Logic:
//   UPCOMING -> SHOWING: hợp lệ
//   UPCOMING -> ENDED:   không hợp lệ
//   SHOWING  -> ENDED:   hợp lệ
//   SHOWING  -> UPCOMING:không hợp lệ
//   ENDED    -> anything: không hợp lệ (default case)
// =============================================================================

// TC-MOVIE-ENTITY-02a
// Mục tiêu: UPCOMING -> SHOWING phải trả về true
func TestCanTransitionTo_UpcomingToShowing(t *testing.T) {
	// --- Khai báo input ---
	inputMovie := &Movie{Status: MovieStatusUpcoming}
	inputNewStatus := MovieStatusShowing
	expectedResult := true

	// --- Gọi method ---
	actualResult := inputMovie.CanTransitionTo(inputNewStatus)

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-MOVIE-ENTITY-02a FAILED: CanTransitionTo(SHOWING) from UPCOMING = %v, want %v", actualResult, expectedResult)
	}
}

// TC-MOVIE-ENTITY-02b
// Mục tiêu: UPCOMING -> ENDED phải trả về false (không được bỏ qua SHOWING)
func TestCanTransitionTo_UpcomingToEnded(t *testing.T) {
	// --- Khai báo input ---
	inputMovie := &Movie{Status: MovieStatusUpcoming}
	inputNewStatus := MovieStatusEnded
	expectedResult := false

	// --- Gọi method ---
	actualResult := inputMovie.CanTransitionTo(inputNewStatus)

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-MOVIE-ENTITY-02b FAILED: CanTransitionTo(ENDED) from UPCOMING = %v, want %v", actualResult, expectedResult)
	}
}

// TC-MOVIE-ENTITY-02c
// Mục tiêu: SHOWING -> ENDED phải trả về true
func TestCanTransitionTo_ShowingToEnded(t *testing.T) {
	// --- Khai báo input ---
	inputMovie := &Movie{Status: MovieStatusShowing}
	inputNewStatus := MovieStatusEnded
	expectedResult := true

	// --- Gọi method ---
	actualResult := inputMovie.CanTransitionTo(inputNewStatus)

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-MOVIE-ENTITY-02c FAILED: CanTransitionTo(ENDED) from SHOWING = %v, want %v", actualResult, expectedResult)
	}
}

// TC-MOVIE-ENTITY-02d
// Mục tiêu: SHOWING -> UPCOMING phải trả về false (không được quay về)
func TestCanTransitionTo_ShowingToUpcoming(t *testing.T) {
	// --- Khai báo input ---
	inputMovie := &Movie{Status: MovieStatusShowing}
	inputNewStatus := MovieStatusUpcoming
	expectedResult := false

	// --- Gọi method ---
	actualResult := inputMovie.CanTransitionTo(inputNewStatus)

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-MOVIE-ENTITY-02d FAILED: CanTransitionTo(UPCOMING) from SHOWING = %v, want %v", actualResult, expectedResult)
	}
}

// TC-MOVIE-ENTITY-02e
// Mục tiêu: ENDED -> bất kỳ status nào đều trả về false (trạng thái cuối)
func TestCanTransitionTo_EndedToShowing(t *testing.T) {
	// --- Khai báo input ---
	inputMovie := &Movie{Status: MovieStatusEnded}
	inputNewStatus := MovieStatusShowing
	expectedResult := false

	// --- Gọi method ---
	actualResult := inputMovie.CanTransitionTo(inputNewStatus)

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-MOVIE-ENTITY-02e FAILED: CanTransitionTo(SHOWING) from ENDED = %v, want %v", actualResult, expectedResult)
	}
}

// TC-MOVIE-ENTITY-02f
// Mục tiêu: ENDED -> UPCOMING phải trả về false
func TestCanTransitionTo_EndedToUpcoming(t *testing.T) {
	// --- Khai báo input ---
	inputMovie := &Movie{Status: MovieStatusEnded}
	inputNewStatus := MovieStatusUpcoming
	expectedResult := false

	// --- Gọi method ---
	actualResult := inputMovie.CanTransitionTo(inputNewStatus)

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-MOVIE-ENTITY-02f FAILED: CanTransitionTo(UPCOMING) from ENDED = %v, want %v", actualResult, expectedResult)
	}
}

// =============================================================================
// TEST NHÓM 3: CreateMovieRequest.ToEntity()
// Nguồn: movie-service/internal/module/movie/entity/types.go
// Logic: Chuyển đổi request DTO sang Movie entity
//   - Nếu Status rỗng => mặc định MovieStatusUpcoming
//   - Nếu Status có giá trị => sử dụng giá trị đó
// =============================================================================

// TC-MOVIE-ENTITY-03a
// Mục tiêu: ToEntity() tạo Movie với status mặc định UPCOMING khi Status rỗng
func TestCreateMovieRequestToEntity_DefaultStatusUpcoming(t *testing.T) {
	// --- Khai báo input ---
	releaseDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	inputRequest := &CreateMovieRequest{
		Title:       "Spider-Man: No Way Home",
		Director:    "Jon Watts",
		Cast:        "Tom Holland, Zendaya",
		Duration:    148,
		ReleaseDate: &releaseDate,
		Description: "Peter Parker seeks help from Doctor Strange",
		TrailerURL:  "https://youtube.com/trailer",
		PosterURL:   "https://cdn.example.com/poster.jpg",
		Status:      "", // Status rỗng => phải mặc định UPCOMING
	}
	expectedStatus := MovieStatusUpcoming

	// --- Gọi method ---
	actualMovie := inputRequest.ToEntity()

	// --- Kiểm tra kết quả ---
	if actualMovie.Status != expectedStatus {
		t.Errorf("TC-MOVIE-ENTITY-03a FAILED: Status = %v, want %v", actualMovie.Status, expectedStatus)
	}
	if actualMovie.Title != inputRequest.Title {
		t.Errorf("TC-MOVIE-ENTITY-03a FAILED: Title = %v, want %v", actualMovie.Title, inputRequest.Title)
	}
	if actualMovie.Duration != inputRequest.Duration {
		t.Errorf("TC-MOVIE-ENTITY-03a FAILED: Duration = %v, want %v", actualMovie.Duration, inputRequest.Duration)
	}
	if actualMovie.Director != inputRequest.Director {
		t.Errorf("TC-MOVIE-ENTITY-03a FAILED: Director = %v, want %v", actualMovie.Director, inputRequest.Director)
	}
}

// TC-MOVIE-ENTITY-03b
// Mục tiêu: ToEntity() sử dụng status được cung cấp khi Status != ""
func TestCreateMovieRequestToEntity_WithExplicitStatus(t *testing.T) {
	// --- Khai báo input ---
	inputRequest := &CreateMovieRequest{
		Title:    "Now Showing Movie",
		Duration: 100,
		Status:   "SHOWING", // Trạng thái được chỉ định rõ
	}
	expectedStatus := MovieStatus("SHOWING")

	// --- Gọi method ---
	actualMovie := inputRequest.ToEntity()

	// --- Kiểm tra kết quả ---
	if actualMovie.Status != expectedStatus {
		t.Errorf("TC-MOVIE-ENTITY-03b FAILED: Status = %v, want %v", actualMovie.Status, expectedStatus)
	}
}

// TC-MOVIE-ENTITY-03c
// Mục tiêu: ToEntity() khi ReleaseDate là nil => Movie.ReleaseDate cũng là nil
func TestCreateMovieRequestToEntity_NilReleaseDate(t *testing.T) {
	// --- Khai báo input ---
	inputRequest := &CreateMovieRequest{
		Title:       "No Release Date Movie",
		Duration:    90,
		ReleaseDate: nil, // Không có ngày phát hành
	}

	// --- Gọi method ---
	actualMovie := inputRequest.ToEntity()

	// --- Kiểm tra kết quả ---
	if actualMovie.ReleaseDate != nil {
		t.Errorf("TC-MOVIE-ENTITY-03c FAILED: ReleaseDate = %v, want nil", actualMovie.ReleaseDate)
	}
}

// =============================================================================
// TEST NHÓM 4: UpdateMovieRequest.ToEntity()
// Nguồn: movie-service/internal/module/movie/entity/types.go
// Logic: Chuyển đổi update request DTO sang Movie entity với id
//   - id phải được gán vào Movie.Id
//   - Nếu Status rỗng => Movie.Status = "" (không gán mặc định)
// =============================================================================

// TC-MOVIE-ENTITY-04a
// Mục tiêu: ToEntity(id) gán đúng Id và các trường khác
func TestUpdateMovieRequestToEntity_AssignsIdCorrectly(t *testing.T) {
	// --- Khai báo input ---
	inputMovieId := "movie-uuid-12345"
	inputRequest := &UpdateMovieRequest{
		Title:    "Updated Title",
		Duration: 150,
		Status:   "SHOWING",
	}
	expectedMovieId := inputMovieId
	expectedStatus := MovieStatus("SHOWING")

	// --- Gọi method ---
	actualMovie := inputRequest.ToEntity(inputMovieId)

	// --- Kiểm tra kết quả ---
	if actualMovie.Id != expectedMovieId {
		t.Errorf("TC-MOVIE-ENTITY-04a FAILED: Id = %v, want %v", actualMovie.Id, expectedMovieId)
	}
	if actualMovie.Status != expectedStatus {
		t.Errorf("TC-MOVIE-ENTITY-04a FAILED: Status = %v, want %v", actualMovie.Status, expectedStatus)
	}
	if actualMovie.Title != inputRequest.Title {
		t.Errorf("TC-MOVIE-ENTITY-04a FAILED: Title = %v, want %v", actualMovie.Title, inputRequest.Title)
	}
}

// TC-MOVIE-ENTITY-04b
// Mục tiêu: ToEntity(id) khi Status rỗng => Movie.Status phải là rỗng (không gán mặc định)
func TestUpdateMovieRequestToEntity_EmptyStatusNotAssigned(t *testing.T) {
	// --- Khai báo input ---
	inputRequest := &UpdateMovieRequest{
		Title:    "Updated Title",
		Duration: 100,
		Status:   "", // Status rỗng => không gán gì
	}

	// --- Gọi method ---
	actualMovie := inputRequest.ToEntity("some-id")

	// --- Kiểm tra kết quả: Status phải là rỗng (không phải UPCOMING như CreateRequest) ---
	if actualMovie.Status != "" {
		t.Errorf("TC-MOVIE-ENTITY-04b FAILED: Status = %v, want empty string", actualMovie.Status)
	}
}

// =============================================================================
// TEST NHÓM 5: ToMovieResponse()
// Nguồn: movie-service/internal/module/movie/entity/types.go
// Logic: Chuyển Movie entity sang MovieResponse DTO
//   - Lấy tên genres từ MovieGenres[i].Genre.Name
// =============================================================================

// TC-MOVIE-ENTITY-05a
// Mục tiêu: ToMovieResponse() chuyển đổi đúng tất cả các trường
func TestToMovieResponse_AllFieldsMappedCorrectly(t *testing.T) {
	// --- Khai báo input ---
	now := time.Now()
	genreName1 := "Action"
	genreName2 := "Adventure"
	inputMovie := &Movie{
		Id:          "movie-id-001",
		Title:       "Avengers: Endgame",
		Director:    "Russo Brothers",
		Cast:        "Robert Downey Jr.",
		Duration:    182,
		Description: "The Avengers assemble for the final battle",
		TrailerURL:  "https://youtube.com/avengers-trailer",
		PosterURL:   "https://cdn.example.com/avengers-poster.jpg",
		Status:      MovieStatusShowing,
		CreatedAt:   &now,
		MovieGenres: []*MovieGenre{
			{Genre: &Genre{Name: genreName1}},
			{Genre: &Genre{Name: genreName2}},
		},
	}
	// --- Gọi method ---
	actualResponse := ToMovieResponse(inputMovie)

	// --- Kiểm tra từng trường ---
	if actualResponse.Id != inputMovie.Id {
		t.Errorf("TC-MOVIE-ENTITY-05a FAILED: Id = %v, want %v", actualResponse.Id, inputMovie.Id)
	}
	if actualResponse.Title != inputMovie.Title {
		t.Errorf("TC-MOVIE-ENTITY-05a FAILED: Title = %v, want %v", actualResponse.Title, inputMovie.Title)
	}
	if actualResponse.Duration != inputMovie.Duration {
		t.Errorf("TC-MOVIE-ENTITY-05a FAILED: Duration = %v, want %v", actualResponse.Duration, inputMovie.Duration)
	}
	if actualResponse.Status != string(inputMovie.Status) {
		t.Errorf("TC-MOVIE-ENTITY-05a FAILED: Status = %v, want %v", actualResponse.Status, string(inputMovie.Status))
	}
	// Kiểm tra genres được lấy đúng tên
	if len(actualResponse.Genres) != 2 {
		t.Errorf("TC-MOVIE-ENTITY-05a FAILED: Genres count = %v, want 2", len(actualResponse.Genres))
	}
	if actualResponse.Genres[0] != genreName1 {
		t.Errorf("TC-MOVIE-ENTITY-05a FAILED: Genres[0] = %v, want %v", actualResponse.Genres[0], genreName1)
	}
	if actualResponse.Genres[1] != genreName2 {
		t.Errorf("TC-MOVIE-ENTITY-05a FAILED: Genres[1] = %v, want %v", actualResponse.Genres[1], genreName2)
	}
}

// TC-MOVIE-ENTITY-05b
// Mục tiêu: ToMovieResponse() khi Movie không có genres => trả về slice rỗng
func TestToMovieResponse_NoGenres(t *testing.T) {
	// --- Khai báo input ---
	inputMovie := &Movie{
		Id:          "movie-id-002",
		Title:       "Solo Movie",
		Duration:    90,
		Status:      MovieStatusUpcoming,
		MovieGenres: []*MovieGenre{}, // Không có genre
	}

	// --- Gọi method ---
	actualResponse := ToMovieResponse(inputMovie)

	// --- Kiểm tra kết quả ---
	if len(actualResponse.Genres) != 0 {
		t.Errorf("TC-MOVIE-ENTITY-05b FAILED: Genres count = %v, want 0", len(actualResponse.Genres))
	}
}

// =============================================================================
// TEST NHÓM 6: ToMoviesResponse()
// Nguồn: movie-service/internal/module/movie/entity/types.go
// Logic: Chuyển danh sách Movie sang GetMoviesResponse với phân trang
//   - TotalPages = ceil(total / size)
//   - Nếu totalPages = 0 => set thành 1
// =============================================================================

// TC-MOVIE-ENTITY-06a
// Mục tiêu: ToMoviesResponse() tính đúng TotalPages khi tổng chia hết cho size
func TestToMoviesResponse_TotalPagesDivisible(t *testing.T) {
	// --- Khai báo input ---
	inputMovies := []*Movie{
		{Id: "m1", Title: "Movie 1", Duration: 90, Status: MovieStatusUpcoming, MovieGenres: []*MovieGenre{}},
		{Id: "m2", Title: "Movie 2", Duration: 120, Status: MovieStatusShowing, MovieGenres: []*MovieGenre{}},
	}
	inputPage := 1
	inputSize := 10
	inputTotal := 20 // 20 / 10 = 2 trang đúng
	expectedTotalPages := 2

	// --- Gọi method ---
	actualResponse := ToMoviesResponse(inputMovies, inputPage, inputSize, inputTotal)

	// --- Kiểm tra kết quả ---
	if actualResponse.Meta.TotalPages != expectedTotalPages {
		t.Errorf("TC-MOVIE-ENTITY-06a FAILED: TotalPages = %v, want %v", actualResponse.Meta.TotalPages, expectedTotalPages)
	}
	if actualResponse.Meta.Page != inputPage {
		t.Errorf("TC-MOVIE-ENTITY-06a FAILED: Page = %v, want %v", actualResponse.Meta.Page, inputPage)
	}
	if actualResponse.Meta.Total != inputTotal {
		t.Errorf("TC-MOVIE-ENTITY-06a FAILED: Total = %v, want %v", actualResponse.Meta.Total, inputTotal)
	}
	if len(actualResponse.Movies) != len(inputMovies) {
		t.Errorf("TC-MOVIE-ENTITY-06a FAILED: Movies count = %v, want %v", len(actualResponse.Movies), len(inputMovies))
	}
}

// TC-MOVIE-ENTITY-06b
// Mục tiêu: ToMoviesResponse() tính đúng TotalPages khi tổng KHÔNG chia hết (ceiling division)
func TestToMoviesResponse_TotalPagesNotDivisible(t *testing.T) {
	// --- Khai báo input ---
	inputMovies := []*Movie{}
	inputPage := 2
	inputSize := 10
	inputTotal := 25        // 25 / 10 = 2.5 => ceil = 3 trang
	expectedTotalPages := 3 // Làm tròn lên

	// --- Gọi method ---
	actualResponse := ToMoviesResponse(inputMovies, inputPage, inputSize, inputTotal)

	// --- Kiểm tra kết quả ---
	if actualResponse.Meta.TotalPages != expectedTotalPages {
		t.Errorf("TC-MOVIE-ENTITY-06b FAILED: TotalPages = %v, want %v", actualResponse.Meta.TotalPages, expectedTotalPages)
	}
}

// TC-MOVIE-ENTITY-06c
// Mục tiêu: ToMoviesResponse() khi total = 0 => TotalPages = 1 (không được = 0)
func TestToMoviesResponse_ZeroTotalGivesOnePage(t *testing.T) {
	// --- Khai báo input ---
	inputMovies := []*Movie{}
	inputPage := 1
	inputSize := 10
	inputTotal := 0        // Không có phim nào
	expectedTotalPages := 1 // Phải trả về ít nhất 1 trang

	// --- Gọi method ---
	actualResponse := ToMoviesResponse(inputMovies, inputPage, inputSize, inputTotal)

	// --- Kiểm tra kết quả ---
	if actualResponse.Meta.TotalPages != expectedTotalPages {
		t.Errorf("TC-MOVIE-ENTITY-06c FAILED: TotalPages = %v, want %v", actualResponse.Meta.TotalPages, expectedTotalPages)
	}
}

// TC-MOVIE-ENTITY-06d
// Mục tiêu: ToMoviesResponse() field Hello luôn trả về "world_ok!"
func TestToMoviesResponse_HelloFieldIsWorldOk(t *testing.T) {
	// --- Khai báo input ---
	inputMovies := []*Movie{}
	expectedHello := "world_ok!"

	// --- Gọi method ---
	actualResponse := ToMoviesResponse(inputMovies, 1, 10, 0)

	// --- Kiểm tra kết quả ---
	if actualResponse.Hello != expectedHello {
		t.Errorf("TC-MOVIE-ENTITY-06d FAILED: Hello = %v, want %v", actualResponse.Hello, expectedHello)
	}
}
