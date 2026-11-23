// =============================================================================
// FILE: movie_business_test.go
// PACKAGE: business
// MÔ TẢ: Unit test cho Movie Business Layer:
//   - CreateMovie()
//   - UpdateMovie()
//   - DeleteMovie()
//   - UpdateMovieStatus()
//   - GetMovieById()
//   - GetMovies()
//   - ValidateMovieForShowtime()
//   - GetGenres()
//   - GetMovieStats()
//
// TESTER: DUCPM88 - Phạm Minh Đức - B22DCDT095
// MỤC TIÊU: Đảm bảo phủ cấp 2 (branch coverage); sử dụng fake/mock thay DB thật
// ROLLBACK: Tests dùng in-memory fake repo - KHÔNG tác động đến DB thật
//
// FAKE PATTERN:
//   - fakeRepo: mock MovieRepository - track được calls và trả về dữ liệu giả
//   - fakeCache: mock Cache interface - luôn trả về cache miss (gọi callback)
//   - fakeRedis: mock redis client - no-op cho Scan/Del
// =============================================================================

package business

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"movie-service/internal/module/movie/entity"
	"movie-service/internal/pkg/caching"

	rediscache "github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

// =============================================================================
// FAKE IMPLEMENTATIONS (Mock Objects)
// =============================================================================

// ducpm88MovieFakeCache: giả lập Cache interface cho movie tests
type ducpm88MovieFakeCache struct{}

func (f *ducpm88MovieFakeCache) Get(_ context.Context, _ string, _ any) error {
	return rediscache.ErrCacheMiss
}

func (f *ducpm88MovieFakeCache) Set(_ context.Context, _ string, _ any, _ time.Duration) error {
	return nil
}

func (f *ducpm88MovieFakeCache) Delete(_ context.Context, _ string) error {
	return nil
}

// Đảm bảo ducpm88MovieFakeCache thoả mãn interface Cache và ReadOnlyCache
var _ caching.Cache = (*ducpm88MovieFakeCache)(nil)
var _ caching.ReadOnlyCache = (*ducpm88MovieFakeCache)(nil)

// ducpm88MovieFakeRedis: giả lập redis.UniversalClient
type ducpm88MovieFakeRedis struct {
	redis.UniversalClient // embed để thoả mãn interface
}

func (f *ducpm88MovieFakeRedis) Scan(_ context.Context, _ uint64, _ string, _ int64) *redis.ScanCmd {
	return redis.NewScanCmd(context.Background(), nil)
}

func (f *ducpm88MovieFakeRedis) Del(_ context.Context, _ ...string) *redis.IntCmd {
	return redis.NewIntCmd(context.Background())
}

func (f *ducpm88MovieFakeRedis) Keys(_ context.Context, _ string) *redis.StringSliceCmd {
	return redis.NewStringSliceCmd(context.Background())
}

// ducpm88MovieFakeRepo: giả lập MovieRepository interface
type ducpm88MovieFakeRepo struct {
	storedMovie  *entity.Movie
	storedMovies []*entity.Movie
	storedTotal  int
	storedStats  []*entity.MovieStat
	storedGenres []*entity.Genre
	getErr       error
	createErr    error
	updateErr    error
	deleteErr    error
	createCalled int // CheckDB: kiểm tra Create đã được gọi chưa
	updateCalled int // CheckDB: kiểm tra Update đã được gọi chưa
	deleteCalled int // CheckDB: kiểm tra Delete đã được gọi chưa
}

func (r *ducpm88MovieFakeRepo) GetByID(_ context.Context, _ string) (*entity.Movie, error) {
	return r.storedMovie, r.getErr
}

func (r *ducpm88MovieFakeRepo) GetMany(_ context.Context, _, _ int, _, _ string) ([]*entity.Movie, error) {
	return r.storedMovies, nil
}

func (r *ducpm88MovieFakeRepo) GetTotalCount(_ context.Context, _, _ string) (int, error) {
	return r.storedTotal, nil
}

func (r *ducpm88MovieFakeRepo) GetMovieStats(_ context.Context) ([]*entity.MovieStat, error) {
	return r.storedStats, nil
}

func (r *ducpm88MovieFakeRepo) GetGenres(_ context.Context) ([]*entity.Genre, error) {
	return r.storedGenres, nil
}

func (r *ducpm88MovieFakeRepo) Create(_ context.Context, _ *entity.Movie, _ []string) error {
	r.createCalled++ // Ghi nhận DB được gọi để tạo movie
	return r.createErr
}

func (r *ducpm88MovieFakeRepo) Update(_ context.Context, _ *entity.Movie, _ []string) error {
	r.updateCalled++ // Ghi nhận DB được gọi để cập nhật movie
	return r.updateErr
}

func (r *ducpm88MovieFakeRepo) Delete(_ context.Context, _ string) error {
	r.deleteCalled++ // Ghi nhận DB được gọi để xóa movie
	return r.deleteErr
}

// newTestBusiness: hàm helper tạo business instance với fakes thay thế dependencies thật
func newTestBusiness(repo *ducpm88MovieFakeRepo) *business {
	return &business{
		repository:  repo,
		cache:       &ducpm88MovieFakeCache{},
		roCache:     &ducpm88MovieFakeCache{},
		redisClient: &ducpm88MovieFakeRedis{},
	}
}

// newValidMovie: tạo Movie hợp lệ dùng làm input mẫu cho các test
func newValidMovie() *entity.Movie {
	now := time.Now()
	return &entity.Movie{
		Id:       "movie-test-001",
		Title:    "Test Movie Title",
		Duration: 120,
		Status:   entity.MovieStatusUpcoming,
		CreatedAt: &now,
	}
}

// =============================================================================
// TEST NHÓM 7: CreateMovie()
// Nguồn: movie-service/internal/module/movie/business/business.go
// Logic:
//   1. movie == nil => ErrInvalidMovieData
//   2. !movie.IsValid() => ErrInvalidMovieData
//   3. Status rỗng => mặc định UPCOMING
//   4. Gọi repo.Create() => CheckDB
//   5. Invalidate cache
// =============================================================================

// TC-MOVIE-BIZ-01a
// Mục tiêu: CreateMovie() trả về ErrInvalidMovieData khi truyền nil
func TestCreateMovie_NilMovie(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{}
	biz := newTestBusiness(repo)
	var inputMovie *entity.Movie = nil // Movie nil
	expectedErr := ErrInvalidMovieData

	// --- Gọi method ---
	actualErr := biz.CreateMovie(context.Background(), inputMovie, nil)

	// --- Kiểm tra kết quả method ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-MOVIE-BIZ-01a FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	// --- CheckDB: Đảm bảo repo.Create KHÔNG được gọi (không tác động DB) ---
	if repo.createCalled != 0 {
		t.Errorf("TC-MOVIE-BIZ-01a FAILED: Create called %d times, want 0 (should not hit DB)", repo.createCalled)
	}
}

// TC-MOVIE-BIZ-01b
// Mục tiêu: CreateMovie() trả về ErrInvalidMovieData khi movie không hợp lệ (Title rỗng)
func TestCreateMovie_InvalidMovie_EmptyTitle(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{}
	biz := newTestBusiness(repo)
	inputMovie := &entity.Movie{
		Title:    "", // Title rỗng => IsValid() = false
		Duration: 100,
	}
	expectedErr := ErrInvalidMovieData

	// --- Gọi method ---
	actualErr := biz.CreateMovie(context.Background(), inputMovie, nil)

	// --- Kiểm tra kết quả method ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-MOVIE-BIZ-01b FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	// --- CheckDB: Đảm bảo repo.Create KHÔNG được gọi ---
	if repo.createCalled != 0 {
		t.Errorf("TC-MOVIE-BIZ-01b FAILED: Create called %d times, want 0", repo.createCalled)
	}
}

// TC-MOVIE-BIZ-01c
// Mục tiêu: CreateMovie() thành công với movie hợp lệ; kiểm tra DB được gọi
func TestCreateMovie_Success(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{createErr: nil} // DB thành công
	biz := newTestBusiness(repo)
	inputMovie := newValidMovie()
	inputGenreIds := []string{"genre-id-01", "genre-id-02"}

	// --- Gọi method ---
	actualErr := biz.CreateMovie(context.Background(), inputMovie, inputGenreIds)

	// --- Kiểm tra kết quả method ---
	if actualErr != nil {
		t.Errorf("TC-MOVIE-BIZ-01c FAILED: unexpected error = %v", actualErr)
	}
	// --- CheckDB: Đảm bảo repo.Create đã được gọi đúng 1 lần (DB được tác động) ---
	if repo.createCalled != 1 {
		t.Errorf("TC-MOVIE-BIZ-01c FAILED: Create called %d times, want 1 (DB must be called)", repo.createCalled)
	}
}

// TC-MOVIE-BIZ-01d
// Mục tiêu: CreateMovie() gán status UPCOMING mặc định khi Status rỗng
func TestCreateMovie_DefaultStatusUpcoming(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{}
	biz := newTestBusiness(repo)
	inputMovie := &entity.Movie{
		Title:    "Movie Without Status",
		Duration: 100,
		Status:   "", // Status rỗng => phải được gán UPCOMING
	}
	expectedStatus := entity.MovieStatusUpcoming

	// --- Gọi method ---
	_ = biz.CreateMovie(context.Background(), inputMovie, nil)

	// --- Kiểm tra kết quả: Status phải được gán mặc định ---
	if inputMovie.Status != expectedStatus {
		t.Errorf("TC-MOVIE-BIZ-01d FAILED: Status = %v, want %v", inputMovie.Status, expectedStatus)
	}
}

// TC-MOVIE-BIZ-01e
// Mục tiêu: CreateMovie() truyền lỗi từ repo.Create()
func TestCreateMovie_RepoError(t *testing.T) {
	// --- Khai báo input ---
	expectedRepoErr := errors.New("db connection failed")
	repo := &ducpm88MovieFakeRepo{createErr: expectedRepoErr}
	biz := newTestBusiness(repo)
	inputMovie := newValidMovie()

	// --- Gọi method ---
	actualErr := biz.CreateMovie(context.Background(), inputMovie, nil)

	// --- Kiểm tra kết quả: lỗi phải chứa thông báo lỗi gốc ---
	if actualErr == nil {
		t.Fatal("TC-MOVIE-BIZ-01e FAILED: expected error but got nil")
	}
	if !errors.Is(actualErr, expectedRepoErr) {
		t.Errorf("TC-MOVIE-BIZ-01e FAILED: error does not wrap repo error. got = %v", actualErr)
	}
	// --- CheckDB: Đảm bảo Create đã được gọi (DB đã tác động, nhưng thất bại) ---
	if repo.createCalled != 1 {
		t.Errorf("TC-MOVIE-BIZ-01e FAILED: Create called %d times, want 1", repo.createCalled)
	}
}

// =============================================================================
// TEST NHÓM 8: UpdateMovie()
// Nguồn: movie-service/internal/module/movie/business/business.go
// Logic:
//   1. movie nil hoặc Id rỗng => ErrInvalidMovieData
//   2. !movie.IsValid() => ErrInvalidMovieData
//   3. GetByID thất bại => ErrMovieNotFound
//   4. Status transition không hợp lệ => ErrInvalidStatusTransition
//   5. Gọi repo.Update() => CheckDB
// =============================================================================

// TC-MOVIE-BIZ-02a
// Mục tiêu: UpdateMovie() trả về ErrInvalidMovieData khi movie nil
func TestUpdateMovie_NilMovie(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{}
	biz := newTestBusiness(repo)
	var inputMovie *entity.Movie = nil
	expectedErr := ErrInvalidMovieData

	// --- Gọi method ---
	actualErr := biz.UpdateMovie(context.Background(), inputMovie, nil)

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-MOVIE-BIZ-02a FAILED: err = %v, want %v", actualErr, expectedErr)
	}
}

// TC-MOVIE-BIZ-02b
// Mục tiêu: UpdateMovie() trả về ErrInvalidMovieData khi movie.Id rỗng
func TestUpdateMovie_EmptyId(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{}
	biz := newTestBusiness(repo)
	inputMovie := &entity.Movie{
		Id:       "", // Id rỗng
		Title:    "Some Title",
		Duration: 100,
	}
	expectedErr := ErrInvalidMovieData

	// --- Gọi method ---
	actualErr := biz.UpdateMovie(context.Background(), inputMovie, nil)

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-MOVIE-BIZ-02b FAILED: err = %v, want %v", actualErr, expectedErr)
	}
}

// TC-MOVIE-BIZ-02c
// Mục tiêu: UpdateMovie() trả về ErrMovieNotFound khi repo.GetByID() thất bại
func TestUpdateMovie_MovieNotFound(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{
		getErr: sql.ErrNoRows, // Giả lập không tìm thấy trong DB
	}
	biz := newTestBusiness(repo)
	inputMovie := &entity.Movie{
		Id:       "non-existent-id",
		Title:    "Some Title",
		Duration: 100,
	}
	expectedErr := ErrMovieNotFound

	// --- Gọi method ---
	actualErr := biz.UpdateMovie(context.Background(), inputMovie, nil)

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-MOVIE-BIZ-02c FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	// --- CheckDB: Update không được gọi vì movie không tồn tại ---
	if repo.updateCalled != 0 {
		t.Errorf("TC-MOVIE-BIZ-02c FAILED: Update called %d times, want 0", repo.updateCalled)
	}
}

// TC-MOVIE-BIZ-02d
// Mục tiêu: UpdateMovie() trả về ErrInvalidStatusTransition khi chuyển trạng thái không hợp lệ (ENDED->SHOWING)
func TestUpdateMovie_InvalidStatusTransition(t *testing.T) {
	// --- Khai báo input ---
	existingMovie := &entity.Movie{
		Id:       "movie-001",
		Title:    "Ended Movie",
		Duration: 90,
		Status:   entity.MovieStatusEnded, // Phim đã kết thúc
	}
	repo := &ducpm88MovieFakeRepo{storedMovie: existingMovie} // DB trả về movie đã ENDED
	biz := newTestBusiness(repo)
	inputMovie := &entity.Movie{
		Id:       "movie-001",
		Title:    "Ended Movie",
		Duration: 90,
		Status:   entity.MovieStatusShowing, // Cố chuyển về SHOWING => không hợp lệ
	}
	expectedErr := ErrInvalidStatusTransition

	// --- Gọi method ---
	actualErr := biz.UpdateMovie(context.Background(), inputMovie, nil)

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-MOVIE-BIZ-02d FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	// --- CheckDB: Update không được gọi vì transition không hợp lệ ---
	if repo.updateCalled != 0 {
		t.Errorf("TC-MOVIE-BIZ-02d FAILED: Update called %d times, want 0", repo.updateCalled)
	}
}

// TC-MOVIE-BIZ-02e
// Mục tiêu: UpdateMovie() thành công với chuyển trạng thái hợp lệ (UPCOMING->SHOWING)
func TestUpdateMovie_Success_ValidTransition(t *testing.T) {
	// --- Khai báo input ---
	existingMovie := &entity.Movie{
		Id:       "movie-001",
		Title:    "Upcoming Movie",
		Duration: 90,
		Status:   entity.MovieStatusUpcoming, // Phim đang upcoming
	}
	repo := &ducpm88MovieFakeRepo{storedMovie: existingMovie}
	biz := newTestBusiness(repo)
	inputMovie := &entity.Movie{
		Id:       "movie-001",
		Title:    "Updated Title",
		Duration: 100,
		Status:   entity.MovieStatusShowing, // Chuyển sang SHOWING => hợp lệ
	}

	// --- Gọi method ---
	actualErr := biz.UpdateMovie(context.Background(), inputMovie, nil)

	// --- Kiểm tra kết quả ---
	if actualErr != nil {
		t.Errorf("TC-MOVIE-BIZ-02e FAILED: unexpected error = %v", actualErr)
	}
	// --- CheckDB: Đảm bảo Update đã được gọi đúng 1 lần ---
	if repo.updateCalled != 1 {
		t.Errorf("TC-MOVIE-BIZ-02e FAILED: Update called %d times, want 1 (DB must be updated)", repo.updateCalled)
	}
}

// TC-MOVIE-BIZ-02f
// Mục tiêu: UpdateMovie() thành công khi status không thay đổi (giữ nguyên UPCOMING)
func TestUpdateMovie_Success_SameStatus(t *testing.T) {
	// --- Khai báo input ---
	existingMovie := &entity.Movie{
		Id:       "movie-001",
		Title:    "Movie",
		Duration: 90,
		Status:   entity.MovieStatusUpcoming,
	}
	repo := &ducpm88MovieFakeRepo{storedMovie: existingMovie}
	biz := newTestBusiness(repo)
	inputMovie := &entity.Movie{
		Id:       "movie-001",
		Title:    "Updated Title Only", // Chỉ update tiêu đề, giữ nguyên status
		Duration: 100,
		Status:   entity.MovieStatusUpcoming, // Status giống DB => không check transition
	}

	// --- Gọi method ---
	actualErr := biz.UpdateMovie(context.Background(), inputMovie, nil)

	// --- Kiểm tra kết quả ---
	if actualErr != nil {
		t.Errorf("TC-MOVIE-BIZ-02f FAILED: unexpected error = %v", actualErr)
	}
	// --- CheckDB: Update phải được gọi ---
	if repo.updateCalled != 1 {
		t.Errorf("TC-MOVIE-BIZ-02f FAILED: Update called %d times, want 1", repo.updateCalled)
	}
}

// =============================================================================
// TEST NHÓM 9: DeleteMovie()
// Nguồn: movie-service/internal/module/movie/business/business.go
// Logic:
//   1. id rỗng => error "movie id is required"
//   2. GetByID thất bại => ErrMovieNotFound
//   3. repo.Delete() thất bại => wrapped error
//   4. Thành công => CheckDB (Delete được gọi)
// =============================================================================

// TC-MOVIE-BIZ-03a
// Mục tiêu: DeleteMovie() trả về lỗi khi id rỗng
func TestDeleteMovie_EmptyId(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{}
	biz := newTestBusiness(repo)
	inputId := "" // Id rỗng

	// --- Gọi method ---
	actualErr := biz.DeleteMovie(context.Background(), inputId)

	// --- Kiểm tra kết quả ---
	if actualErr == nil {
		t.Fatal("TC-MOVIE-BIZ-03a FAILED: expected error but got nil")
	}
	// --- CheckDB: Đảm bảo không gọi DB khi id rỗng ---
	if repo.deleteCalled != 0 {
		t.Errorf("TC-MOVIE-BIZ-03a FAILED: Delete called %d times, want 0", repo.deleteCalled)
	}
}

// TC-MOVIE-BIZ-03b
// Mục tiêu: DeleteMovie() trả về ErrMovieNotFound khi movie không tồn tại
func TestDeleteMovie_MovieNotFound(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{getErr: sql.ErrNoRows}
	biz := newTestBusiness(repo)
	inputId := "non-existent-id"
	expectedErr := ErrMovieNotFound

	// --- Gọi method ---
	actualErr := biz.DeleteMovie(context.Background(), inputId)

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-MOVIE-BIZ-03b FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	// --- CheckDB: Delete không được gọi khi không tìm thấy ---
	if repo.deleteCalled != 0 {
		t.Errorf("TC-MOVIE-BIZ-03b FAILED: Delete called %d times, want 0", repo.deleteCalled)
	}
}

// TC-MOVIE-BIZ-03c
// Mục tiêu: DeleteMovie() thành công; kiểm tra DB đã được gọi
func TestDeleteMovie_Success(t *testing.T) {
	// --- Khai báo input ---
	existingMovie := newValidMovie()
	repo := &ducpm88MovieFakeRepo{storedMovie: existingMovie}
	biz := newTestBusiness(repo)
	inputId := existingMovie.Id

	// --- Gọi method ---
	actualErr := biz.DeleteMovie(context.Background(), inputId)

	// --- Kiểm tra kết quả method ---
	if actualErr != nil {
		t.Errorf("TC-MOVIE-BIZ-03c FAILED: unexpected error = %v", actualErr)
	}
	// --- CheckDB: Đảm bảo Delete đã được gọi đúng 1 lần ---
	if repo.deleteCalled != 1 {
		t.Errorf("TC-MOVIE-BIZ-03c FAILED: Delete called %d times, want 1 (DB must be called)", repo.deleteCalled)
	}
}

// TC-MOVIE-BIZ-03d
// Mục tiêu: DeleteMovie() truyền lỗi từ repo.Delete()
func TestDeleteMovie_RepoDeleteError(t *testing.T) {
	// --- Khai báo input ---
	existingMovie := newValidMovie()
	expectedRepoErr := errors.New("delete constraint violation")
	repo := &ducpm88MovieFakeRepo{
		storedMovie: existingMovie,
		deleteErr:   expectedRepoErr,
	}
	biz := newTestBusiness(repo)

	// --- Gọi method ---
	actualErr := biz.DeleteMovie(context.Background(), existingMovie.Id)

	// --- Kiểm tra kết quả ---
	if actualErr == nil {
		t.Fatal("TC-MOVIE-BIZ-03d FAILED: expected error but got nil")
	}
	if !errors.Is(actualErr, expectedRepoErr) {
		t.Errorf("TC-MOVIE-BIZ-03d FAILED: error does not wrap repo error. got = %v", actualErr)
	}
}

// =============================================================================
// TEST NHÓM 10: UpdateMovieStatus()
// Nguồn: movie-service/internal/module/movie/business/business.go
// Logic:
//   1. id rỗng => error "movie id is required"
//   2. GetByID thất bại => ErrMovieNotFound
//   3. !CanTransitionTo() => ErrInvalidStatusTransition
//   4. Gọi repo.Update() với status mới => CheckDB
// =============================================================================

// TC-MOVIE-BIZ-04a
// Mục tiêu: UpdateMovieStatus() trả về lỗi khi id rỗng
func TestUpdateMovieStatus_EmptyId(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{}
	biz := newTestBusiness(repo)
	inputId := ""
	inputStatus := entity.MovieStatusShowing

	// --- Gọi method ---
	actualErr := biz.UpdateMovieStatus(context.Background(), inputId, inputStatus)

	// --- Kiểm tra kết quả ---
	if actualErr == nil {
		t.Fatal("TC-MOVIE-BIZ-04a FAILED: expected error but got nil")
	}
}

// TC-MOVIE-BIZ-04b
// Mục tiêu: UpdateMovieStatus() trả về ErrMovieNotFound khi không tìm thấy
func TestUpdateMovieStatus_MovieNotFound(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{getErr: sql.ErrNoRows}
	biz := newTestBusiness(repo)
	expectedErr := ErrMovieNotFound

	// --- Gọi method ---
	actualErr := biz.UpdateMovieStatus(context.Background(), "some-id", entity.MovieStatusShowing)

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-MOVIE-BIZ-04b FAILED: err = %v, want %v", actualErr, expectedErr)
	}
}

// TC-MOVIE-BIZ-04c
// Mục tiêu: UpdateMovieStatus() trả về ErrInvalidStatusTransition khi chuyển không hợp lệ
func TestUpdateMovieStatus_InvalidTransition(t *testing.T) {
	// --- Khai báo input ---
	existingMovie := &entity.Movie{
		Id:       "m1",
		Title:    "Ended Movie",
		Duration: 90,
		Status:   entity.MovieStatusEnded, // ENDED không thể chuyển đi đâu
	}
	repo := &ducpm88MovieFakeRepo{storedMovie: existingMovie}
	biz := newTestBusiness(repo)
	expectedErr := ErrInvalidStatusTransition

	// --- Gọi method ---
	actualErr := biz.UpdateMovieStatus(context.Background(), "m1", entity.MovieStatusShowing)

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-MOVIE-BIZ-04c FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	// --- CheckDB: Update không được gọi khi transition sai ---
	if repo.updateCalled != 0 {
		t.Errorf("TC-MOVIE-BIZ-04c FAILED: Update called %d times, want 0", repo.updateCalled)
	}
}

// TC-MOVIE-BIZ-04d
// Mục tiêu: UpdateMovieStatus() thành công; movie.Status được cập nhật và Update được gọi
func TestUpdateMovieStatus_Success(t *testing.T) {
	// --- Khai báo input ---
	existingMovie := &entity.Movie{
		Id:       "m1",
		Title:    "Movie",
		Duration: 90,
		Status:   entity.MovieStatusUpcoming,
	}
	repo := &ducpm88MovieFakeRepo{storedMovie: existingMovie}
	biz := newTestBusiness(repo)
	inputStatus := entity.MovieStatusShowing // UPCOMING -> SHOWING: hợp lệ

	// --- Gọi method ---
	actualErr := biz.UpdateMovieStatus(context.Background(), "m1", inputStatus)

	// --- Kiểm tra kết quả method ---
	if actualErr != nil {
		t.Errorf("TC-MOVIE-BIZ-04d FAILED: unexpected error = %v", actualErr)
	}
	// --- CheckDB: Đảm bảo Update đã được gọi (DB đã được cập nhật) ---
	if repo.updateCalled != 1 {
		t.Errorf("TC-MOVIE-BIZ-04d FAILED: Update called %d times, want 1", repo.updateCalled)
	}
}

// =============================================================================
// TEST NHÓM 11: GetMovieById()
// Nguồn: movie-service/internal/module/movie/business/business.go
// Logic:
//   1. id rỗng => ErrInvalidMovieData
//   2. repo lỗi sql.ErrNoRows => ErrMovieNotFound
//   3. Thành công => trả về *entity.Movie
// =============================================================================

// TC-MOVIE-BIZ-05a
// Mục tiêu: GetMovieById() trả về ErrInvalidMovieData khi id rỗng
func TestGetMovieById_EmptyId(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{}
	biz := newTestBusiness(repo)
	inputId := ""
	expectedErr := ErrInvalidMovieData

	// --- Gọi method ---
	actualMovie, actualErr := biz.GetMovieById(context.Background(), inputId)

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-MOVIE-BIZ-05a FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	if actualMovie != nil {
		t.Errorf("TC-MOVIE-BIZ-05a FAILED: movie = %v, want nil", actualMovie)
	}
}

// TC-MOVIE-BIZ-05b
// Mục tiêu: GetMovieById() trả về ErrMovieNotFound khi DB không có movie
func TestGetMovieById_NotFound(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{getErr: sql.ErrNoRows}
	biz := newTestBusiness(repo)
	expectedErr := ErrMovieNotFound

	// --- Gọi method ---
	actualMovie, actualErr := biz.GetMovieById(context.Background(), "unknown-id")

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-MOVIE-BIZ-05b FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	if actualMovie != nil {
		t.Errorf("TC-MOVIE-BIZ-05b FAILED: movie = %v, want nil", actualMovie)
	}
}

// TC-MOVIE-BIZ-05c
// Mục tiêu: GetMovieById() trả về movie hợp lệ khi tìm thấy
func TestGetMovieById_Success(t *testing.T) {
	// --- Khai báo input ---
	expectedMovie := newValidMovie()
	repo := &ducpm88MovieFakeRepo{storedMovie: expectedMovie}
	biz := newTestBusiness(repo)

	// --- Gọi method ---
	actualMovie, actualErr := biz.GetMovieById(context.Background(), expectedMovie.Id)

	// --- Kiểm tra kết quả ---
	if actualErr != nil {
		t.Errorf("TC-MOVIE-BIZ-05c FAILED: unexpected error = %v", actualErr)
	}
	if actualMovie == nil {
		t.Fatal("TC-MOVIE-BIZ-05c FAILED: movie is nil")
	}
	if actualMovie.Id != expectedMovie.Id {
		t.Errorf("TC-MOVIE-BIZ-05c FAILED: Id = %v, want %v", actualMovie.Id, expectedMovie.Id)
	}
}

// =============================================================================
// TEST NHÓM 12: GetMovies()
// Nguồn: movie-service/internal/module/movie/business/business.go
// Logic:
//   - page<=0 => mặc định 1
//   - size<=0 => mặc định 10
//   - size>100 => giới hạn 100
// =============================================================================

// TC-MOVIE-BIZ-06a
// Mục tiêu: GetMovies() chuẩn hóa page<=0 thành 1
func TestGetMovies_NegativePageDefaultsTo1(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{
		storedMovies: []*entity.Movie{},
		storedTotal:  0,
	}
	biz := newTestBusiness(repo)
	inputPage := -1  // Page không hợp lệ => phải chuẩn hóa thành 1
	inputSize := 10

	// --- Gọi method ---
	actualMovies, actualTotal, actualErr := biz.GetMovies(context.Background(), inputPage, inputSize, "", "")

	// --- Kiểm tra kết quả ---
	if actualErr != nil {
		t.Errorf("TC-MOVIE-BIZ-06a FAILED: unexpected error = %v", actualErr)
	}
	if actualMovies == nil {
		t.Error("TC-MOVIE-BIZ-06a FAILED: movies should not be nil")
	}
	_ = actualTotal // total không cần check ở test này
}

// TC-MOVIE-BIZ-06b
// Mục tiêu: GetMovies() giới hạn size > 100 thành 100
func TestGetMovies_SizeOver100CappedAt100(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{
		storedMovies: []*entity.Movie{},
		storedTotal:  5,
	}
	biz := newTestBusiness(repo)
	inputPage := 1
	inputSize := 999 // Vượt quá 100 => phải giới hạn lại

	// --- Gọi method (không lỗi, chỉ cần chạy được) ---
	_, _, actualErr := biz.GetMovies(context.Background(), inputPage, inputSize, "", "")

	// --- Kiểm tra kết quả ---
	if actualErr != nil {
		t.Errorf("TC-MOVIE-BIZ-06b FAILED: unexpected error = %v", actualErr)
	}
}

// TC-MOVIE-BIZ-06c
// Mục tiêu: GetMovies() trả về đúng danh sách movies và total từ repo
func TestGetMovies_ReturnsMoviesAndTotal(t *testing.T) {
	// --- Khai báo input ---
	expectedMovies := []*entity.Movie{
		{Id: "m1", Title: "Movie 1", Duration: 90, Status: entity.MovieStatusUpcoming},
		{Id: "m2", Title: "Movie 2", Duration: 120, Status: entity.MovieStatusShowing},
	}
	expectedTotal := 15
	repo := &ducpm88MovieFakeRepo{
		storedMovies: expectedMovies,
		storedTotal:  expectedTotal,
	}
	biz := newTestBusiness(repo)

	// --- Gọi method ---
	actualMovies, actualTotal, actualErr := biz.GetMovies(context.Background(), 1, 10, "", "")

	// --- Kiểm tra kết quả ---
	if actualErr != nil {
		t.Errorf("TC-MOVIE-BIZ-06c FAILED: unexpected error = %v", actualErr)
	}
	if len(actualMovies) != len(expectedMovies) {
		t.Errorf("TC-MOVIE-BIZ-06c FAILED: movies count = %d, want %d", len(actualMovies), len(expectedMovies))
	}
	if actualTotal != expectedTotal {
		t.Errorf("TC-MOVIE-BIZ-06c FAILED: total = %d, want %d", actualTotal, expectedTotal)
	}
}

// =============================================================================
// TEST NHÓM 13: ValidateMovieForShowtime()
// Nguồn: movie-service/internal/module/movie/business/business.go
// Logic:
//   1. movieId rỗng => ErrInvalidMovieData
//   2. movie không tìm thấy => ErrMovieNotFound
//   3. movie.Status != SHOWING => ErrMovieNotShowing
//   4. movie.Status == SHOWING => nil (hợp lệ)
// =============================================================================

// TC-MOVIE-BIZ-07a
// Mục tiêu: ValidateMovieForShowtime() trả về ErrInvalidMovieData khi id rỗng
func TestValidateMovieForShowtime_EmptyId(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88MovieFakeRepo{}
	biz := newTestBusiness(repo)
	expectedErr := ErrInvalidMovieData

	// --- Gọi method ---
	actualErr := biz.ValidateMovieForShowtime(context.Background(), "")

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-MOVIE-BIZ-07a FAILED: err = %v, want %v", actualErr, expectedErr)
	}
}

// TC-MOVIE-BIZ-07b
// Mục tiêu: ValidateMovieForShowtime() trả về ErrMovieNotShowing khi movie là UPCOMING
func TestValidateMovieForShowtime_UpcomingMovie(t *testing.T) {
	// --- Khai báo input ---
	upcomingMovie := &entity.Movie{
		Id:       "m1",
		Title:    "Upcoming Movie",
		Duration: 90,
		Status:   entity.MovieStatusUpcoming, // Chưa chiếu
	}
	repo := &ducpm88MovieFakeRepo{storedMovie: upcomingMovie}
	biz := newTestBusiness(repo)
	expectedErr := ErrMovieNotShowing

	// --- Gọi method ---
	actualErr := biz.ValidateMovieForShowtime(context.Background(), "m1")

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-MOVIE-BIZ-07b FAILED: err = %v, want %v", actualErr, expectedErr)
	}
}

// TC-MOVIE-BIZ-07c
// Mục tiêu: ValidateMovieForShowtime() trả về nil khi movie đang SHOWING
func TestValidateMovieForShowtime_ShowingMovie_Success(t *testing.T) {
	// --- Khai báo input ---
	showingMovie := &entity.Movie{
		Id:       "m1",
		Title:    "Now Showing",
		Duration: 90,
		Status:   entity.MovieStatusShowing, // Đang chiếu => hợp lệ cho showtime
	}
	repo := &ducpm88MovieFakeRepo{storedMovie: showingMovie}
	biz := newTestBusiness(repo)

	// --- Gọi method ---
	actualErr := biz.ValidateMovieForShowtime(context.Background(), "m1")

	// --- Kiểm tra kết quả ---
	if actualErr != nil {
		t.Errorf("TC-MOVIE-BIZ-07c FAILED: unexpected error = %v", actualErr)
	}
}

// TC-MOVIE-BIZ-07d
// Mục tiêu: ValidateMovieForShowtime() trả về ErrMovieNotShowing khi movie là ENDED
func TestValidateMovieForShowtime_EndedMovie(t *testing.T) {
	// --- Khai báo input ---
	endedMovie := &entity.Movie{
		Id:       "m1",
		Title:    "Ended Movie",
		Duration: 90,
		Status:   entity.MovieStatusEnded, // Đã kết thúc
	}
	repo := &ducpm88MovieFakeRepo{storedMovie: endedMovie}
	biz := newTestBusiness(repo)
	expectedErr := ErrMovieNotShowing

	// --- Gọi method ---
	actualErr := biz.ValidateMovieForShowtime(context.Background(), "m1")

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-MOVIE-BIZ-07d FAILED: err = %v, want %v", actualErr, expectedErr)
	}
}

// =============================================================================
// TEST NHÓM 14: GetGenres() và GetMovieStats()
// =============================================================================

// TC-MOVIE-BIZ-08a
// Mục tiêu: GetGenres() trả về danh sách genres từ repo
func TestGetGenres_ReturnsGenres(t *testing.T) {
	// --- Khai báo input ---
	expectedGenres := []*entity.Genre{
		{Id: "g1", Name: "Action", Slug: "action"},
		{Id: "g2", Name: "Comedy", Slug: "comedy"},
	}
	repo := &ducpm88MovieFakeRepo{storedGenres: expectedGenres}
	biz := newTestBusiness(repo)

	// --- Gọi method ---
	actualGenres, actualErr := biz.GetGenres(context.Background())

	// --- Kiểm tra kết quả ---
	if actualErr != nil {
		t.Errorf("TC-MOVIE-BIZ-08a FAILED: unexpected error = %v", actualErr)
	}
	if len(actualGenres) != len(expectedGenres) {
		t.Errorf("TC-MOVIE-BIZ-08a FAILED: genres count = %d, want %d", len(actualGenres), len(expectedGenres))
	}
}

// TC-MOVIE-BIZ-08b
// Mục tiêu: GetMovieStats() trả về danh sách thống kê từ repo
func TestGetMovieStats_ReturnsStats(t *testing.T) {
	// --- Khai báo input ---
	expectedStats := []*entity.MovieStat{
		{Status: "UPCOMING", Count: 5},
		{Status: "SHOWING", Count: 3},
		{Status: "ENDED", Count: 10},
	}
	repo := &ducpm88MovieFakeRepo{storedStats: expectedStats}
	biz := newTestBusiness(repo)

	// --- Gọi method ---
	actualStats, actualErr := biz.GetMovieStats(context.Background())

	// --- Kiểm tra kết quả ---
	if actualErr != nil {
		t.Errorf("TC-MOVIE-BIZ-08b FAILED: unexpected error = %v", actualErr)
	}
	if len(actualStats) != len(expectedStats) {
		t.Errorf("TC-MOVIE-BIZ-08b FAILED: stats count = %d, want %d", len(actualStats), len(expectedStats))
	}
	if actualStats[0].Count != expectedStats[0].Count {
		t.Errorf("TC-MOVIE-BIZ-08b FAILED: stats[0].Count = %d, want %d", actualStats[0].Count, expectedStats[0].Count)
	}
}
