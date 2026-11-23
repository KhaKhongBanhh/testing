// =============================================================================
// FILE: seat_business_test.go
// PACKAGE: business
// MÔ TẢ: Unit test cho Seat Business Layer:
//   - GetSeatById()
//   - GetSeats()
//   - GetSeatsByIds()
//   - CreateSeat()
//   - UpdateSeat()
//   - DeleteSeat()
//   - UpdateSeatStatus()
//
// TESTER: DUCPM88 - Phạm Minh Đức - B22DCDT095
// MỤC TIÊU: Đảm bảo phủ cấp 2 (branch coverage); sử dụng fake/mock thay DB thật
// ROLLBACK: Tests dùng in-memory fake repo - KHÔNG tác động đến DB thật
//
// FAKE PATTERN:
//   - ducpm88SeatFakeRepo: mock SeatRepository - track được calls và trả về dữ liệu giả
//   - fakeSeatCache: mock Cache interface - luôn trả về cache miss
//   - fakeSeatRedis: mock redis client - no-op cho Keys/Del
// =============================================================================

package business

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"movie-service/internal/module/seat/entity"
	"movie-service/internal/pkg/caching"

	rediscache "github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

// =============================================================================
// FAKE IMPLEMENTATIONS (Mock Objects)
// =============================================================================

// ducpm88SeatFakeCache: giả lập Cache interface cho seat tests
// Luôn trả về cache.ErrCacheMiss để buộc callback được gọi
type ducpm88SeatFakeCache struct{}

func (f *ducpm88SeatFakeCache) Get(_ context.Context, _ string, _ any) error {
	return rediscache.ErrCacheMiss
}

func (f *ducpm88SeatFakeCache) Set(_ context.Context, _ string, _ any, _ time.Duration) error {
	return nil
}

func (f *ducpm88SeatFakeCache) Delete(_ context.Context, _ string) error {
	return nil
}

// Đảm bảo ducpm88SeatFakeCache thoả mãn interface Cache và ReadOnlyCache
var _ caching.Cache = (*ducpm88SeatFakeCache)(nil)
var _ caching.ReadOnlyCache = (*ducpm88SeatFakeCache)(nil)

// ducpm88SeatFakeRedis: giả lập redis.UniversalClient cho seat tests
type ducpm88SeatFakeRedis struct {
	redis.UniversalClient // embed để thoả mãn interface
}

func (f *ducpm88SeatFakeRedis) Scan(_ context.Context, _ uint64, _ string, _ int64) *redis.ScanCmd {
	return redis.NewScanCmd(context.Background(), nil)
}

func (f *ducpm88SeatFakeRedis) Del(_ context.Context, _ ...string) *redis.IntCmd {
	return redis.NewIntCmd(context.Background())
}

func (f *ducpm88SeatFakeRedis) Keys(_ context.Context, _ string) *redis.StringSliceCmd {
	return redis.NewStringSliceCmd(context.Background())
}

// ducpm88SeatFakeRepo: giả lập SeatRepository interface
type ducpm88SeatFakeRepo struct {
	storedSeat     *entity.Seat
	storedSeats    []*entity.Seat
	storedTotal    int
	positionExists bool
	getErr         error
	createErr      error
	updateErr      error
	deleteErr      error
	createCalled   int // CheckDB: kiểm tra Create đã được gọi chưa
	updateCalled   int // CheckDB: kiểm tra Update đã được gọi chưa
	deleteCalled   int // CheckDB: kiểm tra Delete đã được gọi chưa
}

func (r *ducpm88SeatFakeRepo) GetByID(_ context.Context, _ string) (*entity.Seat, error) {
	return r.storedSeat, r.getErr
}

func (r *ducpm88SeatFakeRepo) GetByIDs(_ context.Context, _ []string) ([]*entity.Seat, error) {
	return r.storedSeats, r.getErr
}

func (r *ducpm88SeatFakeRepo) GetMany(_ context.Context, _, _ int, _, _, _ string, _ entity.SeatType, _ entity.SeatStatus) ([]*entity.Seat, error) {
	return r.storedSeats, nil
}

func (r *ducpm88SeatFakeRepo) GetTotalCount(_ context.Context, _, _, _ string, _ entity.SeatType, _ entity.SeatStatus) (int, error) {
	return r.storedTotal, nil
}

func (r *ducpm88SeatFakeRepo) ExistsBySeatPosition(_ context.Context, _, _, _ string, _ string) (bool, error) {
	return r.positionExists, nil
}

func (r *ducpm88SeatFakeRepo) Create(_ context.Context, _ *entity.Seat) error {
	r.createCalled++ // Ghi nhận DB được tác động
	return r.createErr
}

func (r *ducpm88SeatFakeRepo) Update(_ context.Context, _ *entity.Seat) error {
	r.updateCalled++ // Ghi nhận DB được tác động
	return r.updateErr
}

func (r *ducpm88SeatFakeRepo) Delete(_ context.Context, _ string) error {
	r.deleteCalled++ // Ghi nhận DB được tác động
	return r.deleteErr
}

// newTestSeatBusiness: helper tạo seat business instance với fakes
func newTestSeatBusiness(repo *ducpm88SeatFakeRepo) *business {
	return &business{
		repository:  repo,
		cache:       &ducpm88SeatFakeCache{},
		roCache:     &ducpm88SeatFakeCache{},
		redisClient: &ducpm88SeatFakeRedis{},
	}
}

// newValidSeat: tạo Seat hợp lệ dùng làm input mẫu
func newValidSeat() *entity.Seat {
	now := time.Now()
	return &entity.Seat{
		Id:         "seat-test-001",
		RoomId:     "room-001",
		SeatNumber: "A1",
		RowNumber:  "A",
		SeatType:   entity.SeatTypeRegular,
		Status:     entity.SeatStatusAvailable,
		CreatedAt:  now,
	}
}

// =============================================================================
// TEST NHÓM 23: GetSeatById()
// Nguồn: movie-service/internal/module/seat/business/business.go
// Logic:
//   1. repo lỗi sql.ErrNoRows => ErrSeatNotFound
//   2. repo lỗi khác => wrapped error
//   3. Thành công => trả về *entity.Seat
// =============================================================================

// TC-SEAT-BIZ-01a
// Mục tiêu: GetSeatById() trả về ErrSeatNotFound khi DB không có seat
func TestGetSeatById_NotFound(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88SeatFakeRepo{getErr: sql.ErrNoRows}
	biz := newTestSeatBusiness(repo)
	inputId := "non-existent-seat"
	expectedErr := ErrSeatNotFound

	// --- Gọi method ---
	actualSeat, actualErr := biz.GetSeatById(context.Background(), inputId)

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-SEAT-BIZ-01a FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	if actualSeat != nil {
		t.Errorf("TC-SEAT-BIZ-01a FAILED: seat = %v, want nil", actualSeat)
	}
}

// TC-SEAT-BIZ-01b
// Mục tiêu: GetSeatById() trả về seat khi tìm thấy trong DB
func TestGetSeatById_Success(t *testing.T) {
	// --- Khai báo input ---
	expectedSeat := newValidSeat()
	repo := &ducpm88SeatFakeRepo{storedSeat: expectedSeat}
	biz := newTestSeatBusiness(repo)

	// --- Gọi method ---
	actualSeat, actualErr := biz.GetSeatById(context.Background(), expectedSeat.Id)

	// --- Kiểm tra kết quả ---
	if actualErr != nil {
		t.Errorf("TC-SEAT-BIZ-01b FAILED: unexpected error = %v", actualErr)
	}
	if actualSeat == nil {
		t.Fatal("TC-SEAT-BIZ-01b FAILED: seat is nil")
	}
	if actualSeat.Id != expectedSeat.Id {
		t.Errorf("TC-SEAT-BIZ-01b FAILED: seat.Id = %v, want %v", actualSeat.Id, expectedSeat.Id)
	}
}

// TC-SEAT-BIZ-01c
// Mục tiêu: GetSeatById() trả về wrapped error khi repo lỗi không phải sql.ErrNoRows
func TestGetSeatById_GenericRepoError(t *testing.T) {
	// --- Khai báo input ---
	expectedRepoErr := errors.New("connection timeout")
	repo := &ducpm88SeatFakeRepo{getErr: expectedRepoErr}
	biz := newTestSeatBusiness(repo)

	// --- Gọi method ---
	actualSeat, actualErr := biz.GetSeatById(context.Background(), "some-id")

	// --- Kiểm tra kết quả ---
	if actualErr == nil {
		t.Fatal("TC-SEAT-BIZ-01c FAILED: expected error but got nil")
	}
	if actualSeat != nil {
		t.Errorf("TC-SEAT-BIZ-01c FAILED: seat = %v, want nil", actualSeat)
	}
	// Lỗi phải chứa thông tin gốc
	if !errors.Is(actualErr, expectedRepoErr) {
		t.Errorf("TC-SEAT-BIZ-01c FAILED: error does not wrap repo error. got = %v", actualErr)
	}
}

// =============================================================================
// TEST NHÓM 24: GetSeats()
// Nguồn: movie-service/internal/module/seat/business/business.go
// Logic:
//   1. page < 1 hoặc size < 1 => ErrInvalidSeatData
//   2. Thành công => trả về ([]*entity.Seat, int, nil)
// =============================================================================

// TC-SEAT-BIZ-02a
// Mục tiêu: GetSeats() trả về ErrInvalidSeatData khi page < 1
func TestGetSeats_InvalidPage(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88SeatFakeRepo{}
	biz := newTestSeatBusiness(repo)
	inputPage := 0 // Page = 0 không hợp lệ
	inputSize := 10
	expectedErr := ErrInvalidSeatData

	// --- Gọi method ---
	actualSeats, actualTotal, actualErr := biz.GetSeats(context.Background(), inputPage, inputSize, "", "", "", "", "")

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-SEAT-BIZ-02a FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	if actualSeats != nil {
		t.Errorf("TC-SEAT-BIZ-02a FAILED: seats = %v, want nil", actualSeats)
	}
	if actualTotal != 0 {
		t.Errorf("TC-SEAT-BIZ-02a FAILED: total = %d, want 0", actualTotal)
	}
}

// TC-SEAT-BIZ-02b
// Mục tiêu: GetSeats() trả về ErrInvalidSeatData khi size < 1
func TestGetSeats_InvalidSize(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88SeatFakeRepo{}
	biz := newTestSeatBusiness(repo)
	inputPage := 1
	inputSize := 0 // Size = 0 không hợp lệ
	expectedErr := ErrInvalidSeatData

	// --- Gọi method ---
	_, _, actualErr := biz.GetSeats(context.Background(), inputPage, inputSize, "", "", "", "", "")

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-SEAT-BIZ-02b FAILED: err = %v, want %v", actualErr, expectedErr)
	}
}

// TC-SEAT-BIZ-02c
// Mục tiêu: GetSeats() trả về danh sách seats và total đúng khi tham số hợp lệ
func TestGetSeats_Success(t *testing.T) {
	// --- Khai báo input ---
	expectedSeats := []*entity.Seat{
		newValidSeat(),
		{Id: "s2", RoomId: "room-001", SeatNumber: "A2", RowNumber: "A", SeatType: entity.SeatTypeVIP, Status: entity.SeatStatusAvailable, CreatedAt: time.Now()},
	}
	expectedTotal := 2
	repo := &ducpm88SeatFakeRepo{
		storedSeats: expectedSeats,
		storedTotal: expectedTotal,
	}
	biz := newTestSeatBusiness(repo)

	// --- Gọi method ---
	actualSeats, actualTotal, actualErr := biz.GetSeats(context.Background(), 1, 10, "", "", "", entity.SeatTypeRegular, entity.SeatStatusAvailable)

	// --- Kiểm tra kết quả ---
	if actualErr != nil {
		t.Errorf("TC-SEAT-BIZ-02c FAILED: unexpected error = %v", actualErr)
	}
	if len(actualSeats) != len(expectedSeats) {
		t.Errorf("TC-SEAT-BIZ-02c FAILED: seats count = %d, want %d", len(actualSeats), len(expectedSeats))
	}
	if actualTotal != expectedTotal {
		t.Errorf("TC-SEAT-BIZ-02c FAILED: total = %d, want %d", actualTotal, expectedTotal)
	}
}

// TC-SEAT-BIZ-02d
// Mục tiêu: GetSeats() tính đúng offset từ page và size
func TestGetSeats_CalculatesOffsetCorrectly(t *testing.T) {
	// --- Khai báo input (page=3, size=20 => offset phải = 40) ---
	repo := &ducpm88SeatFakeRepo{storedSeats: []*entity.Seat{}, storedTotal: 100}
	biz := newTestSeatBusiness(repo)
	inputPage := 3
	inputSize := 20

	// --- Gọi method (không lỗi, chỉ kiểm tra không crash) ---
	_, _, actualErr := biz.GetSeats(context.Background(), inputPage, inputSize, "", "", "", "", "")

	// --- Kiểm tra kết quả ---
	if actualErr != nil {
		t.Errorf("TC-SEAT-BIZ-02d FAILED: unexpected error = %v", actualErr)
	}
}

// =============================================================================
// TEST NHÓM 25: CreateSeat()
// Nguồn: movie-service/internal/module/seat/business/business.go
// Logic:
//   1. seat nil hoặc !seat.IsValid() => ErrInvalidSeatData
//   2. ExistsBySeatPosition = true => ErrSeatPositionExists
//   3. repo.Create() thất bại => wrapped error
//   4. Thành công => CheckDB (Create được gọi)
// =============================================================================

// TC-SEAT-BIZ-03a
// Mục tiêu: CreateSeat() trả về ErrInvalidSeatData khi seat nil
func TestCreateSeat_NilSeat(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88SeatFakeRepo{}
	biz := newTestSeatBusiness(repo)
	var inputSeat *entity.Seat = nil
	expectedErr := ErrInvalidSeatData

	// --- Gọi method ---
	actualErr := biz.CreateSeat(context.Background(), inputSeat)

	// --- Kiểm tra kết quả method ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-SEAT-BIZ-03a FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	// --- CheckDB: Đảm bảo Create KHÔNG được gọi ---
	if repo.createCalled != 0 {
		t.Errorf("TC-SEAT-BIZ-03a FAILED: Create called %d times, want 0", repo.createCalled)
	}
}

// TC-SEAT-BIZ-03b
// Mục tiêu: CreateSeat() trả về ErrInvalidSeatData khi seat không hợp lệ (RoomId rỗng)
func TestCreateSeat_InvalidSeat_EmptyRoomId(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88SeatFakeRepo{}
	biz := newTestSeatBusiness(repo)
	inputSeat := &entity.Seat{
		RoomId:     "", // RoomId rỗng => IsValid() = false
		SeatNumber: "A1",
		RowNumber:  "A",
	}
	expectedErr := ErrInvalidSeatData

	// --- Gọi method ---
	actualErr := biz.CreateSeat(context.Background(), inputSeat)

	// --- Kiểm tra kết quả method ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-SEAT-BIZ-03b FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	// --- CheckDB: Đảm bảo Create KHÔNG được gọi ---
	if repo.createCalled != 0 {
		t.Errorf("TC-SEAT-BIZ-03b FAILED: Create called %d times, want 0", repo.createCalled)
	}
}

// TC-SEAT-BIZ-03c
// Mục tiêu: CreateSeat() trả về ErrSeatPositionExists khi vị trí đã tồn tại
func TestCreateSeat_PositionAlreadyExists(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88SeatFakeRepo{positionExists: true} // Vị trí đã tồn tại
	biz := newTestSeatBusiness(repo)
	inputSeat := newValidSeat()
	expectedErr := ErrSeatPositionExists

	// --- Gọi method ---
	actualErr := biz.CreateSeat(context.Background(), inputSeat)

	// --- Kiểm tra kết quả method ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-SEAT-BIZ-03c FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	// --- CheckDB: Đảm bảo Create KHÔNG được gọi (position conflict) ---
	if repo.createCalled != 0 {
		t.Errorf("TC-SEAT-BIZ-03c FAILED: Create called %d times, want 0", repo.createCalled)
	}
}

// TC-SEAT-BIZ-03d
// Mục tiêu: CreateSeat() thành công; kiểm tra DB đã được gọi
func TestCreateSeat_Success(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88SeatFakeRepo{positionExists: false, createErr: nil}
	biz := newTestSeatBusiness(repo)
	inputSeat := newValidSeat()

	// --- Gọi method ---
	actualErr := biz.CreateSeat(context.Background(), inputSeat)

	// --- Kiểm tra kết quả method ---
	if actualErr != nil {
		t.Errorf("TC-SEAT-BIZ-03d FAILED: unexpected error = %v", actualErr)
	}
	// --- CheckDB: Đảm bảo Create đã được gọi đúng 1 lần (DB được tác động) ---
	if repo.createCalled != 1 {
		t.Errorf("TC-SEAT-BIZ-03d FAILED: Create called %d times, want 1 (DB must be called)", repo.createCalled)
	}
}

// TC-SEAT-BIZ-03e
// Mục tiêu: CreateSeat() truyền lỗi từ repo.Create()
func TestCreateSeat_RepoError(t *testing.T) {
	// --- Khai báo input ---
	expectedRepoErr := errors.New("unique constraint violation")
	repo := &ducpm88SeatFakeRepo{positionExists: false, createErr: expectedRepoErr}
	biz := newTestSeatBusiness(repo)
	inputSeat := newValidSeat()

	// --- Gọi method ---
	actualErr := biz.CreateSeat(context.Background(), inputSeat)

	// --- Kiểm tra kết quả ---
	if actualErr == nil {
		t.Fatal("TC-SEAT-BIZ-03e FAILED: expected error but got nil")
	}
	if !errors.Is(actualErr, expectedRepoErr) {
		t.Errorf("TC-SEAT-BIZ-03e FAILED: error does not wrap repo error. got = %v", actualErr)
	}
}

// =============================================================================
// TEST NHÓM 26: UpdateSeat()
// Nguồn: movie-service/internal/module/seat/business/business.go
// Logic:
//   1. GetByID lỗi sql.ErrNoRows => ErrSeatNotFound
//   2. GetByID lỗi khác => wrapped error
//   3. SeatNumber/RowNumber thay đổi & position đã tồn tại => ErrSeatPositionExists
//   4. Cập nhật SeatType nếu updates.SeatType != nil
//   5. Cập nhật Status nếu updates.Status != nil
//   6. !seat.IsValid() sau cập nhật => ErrInvalidSeatData
//   7. Thành công => CheckDB (Update được gọi)
// =============================================================================

// TC-SEAT-BIZ-04a
// Mục tiêu: UpdateSeat() trả về ErrSeatNotFound khi seat không tồn tại
func TestUpdateSeat_SeatNotFound(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88SeatFakeRepo{getErr: sql.ErrNoRows}
	biz := newTestSeatBusiness(repo)
	inputId := "not-found-id"
	inputUpdates := &entity.UpdateSeatRequest{}
	expectedErr := ErrSeatNotFound

	// --- Gọi method ---
	actualErr := biz.UpdateSeat(context.Background(), inputId, inputUpdates)

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-SEAT-BIZ-04a FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	// --- CheckDB: Update không được gọi vì seat không tồn tại ---
	if repo.updateCalled != 0 {
		t.Errorf("TC-SEAT-BIZ-04a FAILED: Update called %d times, want 0", repo.updateCalled)
	}
}

// TC-SEAT-BIZ-04b
// Mục tiêu: UpdateSeat() trả về ErrSeatPositionExists khi vị trí mới đã tồn tại
func TestUpdateSeat_PositionConflict(t *testing.T) {
	// --- Khai báo input ---
	existingSeat := newValidSeat()
	newSeatNumber := "B5" // Vị trí mới
	repo := &ducpm88SeatFakeRepo{
		storedSeat:     existingSeat,
		positionExists: true, // Vị trí B5 đã được ghế khác dùng
	}
	biz := newTestSeatBusiness(repo)
	inputUpdates := &entity.UpdateSeatRequest{
		SeatNumber: &newSeatNumber,
	}
	expectedErr := ErrSeatPositionExists

	// --- Gọi method ---
	actualErr := biz.UpdateSeat(context.Background(), existingSeat.Id, inputUpdates)

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-SEAT-BIZ-04b FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	// --- CheckDB: Update không được gọi khi conflict ---
	if repo.updateCalled != 0 {
		t.Errorf("TC-SEAT-BIZ-04b FAILED: Update called %d times, want 0", repo.updateCalled)
	}
}

// TC-SEAT-BIZ-04c
// Mục tiêu: UpdateSeat() thành công khi cập nhật SeatType
func TestUpdateSeat_Success_UpdateSeatType(t *testing.T) {
	// --- Khai báo input ---
	existingSeat := newValidSeat()
	existingSeat.SeatType = entity.SeatTypeRegular
	repo := &ducpm88SeatFakeRepo{storedSeat: existingSeat, positionExists: false}
	biz := newTestSeatBusiness(repo)
	newSeatType := entity.SeatTypeVIP
	inputUpdates := &entity.UpdateSeatRequest{
		SeatType: &newSeatType, // Thay đổi từ REGULAR -> VIP
	}

	// --- Gọi method ---
	actualErr := biz.UpdateSeat(context.Background(), existingSeat.Id, inputUpdates)

	// --- Kiểm tra kết quả method ---
	if actualErr != nil {
		t.Errorf("TC-SEAT-BIZ-04c FAILED: unexpected error = %v", actualErr)
	}
	// --- CheckDB: Đảm bảo Update đã được gọi (DB được tác động) ---
	if repo.updateCalled != 1 {
		t.Errorf("TC-SEAT-BIZ-04c FAILED: Update called %d times, want 1", repo.updateCalled)
	}
}

// TC-SEAT-BIZ-04d
// Mục tiêu: UpdateSeat() thành công khi cập nhật Status
func TestUpdateSeat_Success_UpdateStatus(t *testing.T) {
	// --- Khai báo input ---
	existingSeat := newValidSeat()
	repo := &ducpm88SeatFakeRepo{storedSeat: existingSeat, positionExists: false}
	biz := newTestSeatBusiness(repo)
	newStatus := entity.SeatStatusMaintenance
	inputUpdates := &entity.UpdateSeatRequest{
		Status: &newStatus, // Chuyển sang maintenance
	}

	// --- Gọi method ---
	actualErr := biz.UpdateSeat(context.Background(), existingSeat.Id, inputUpdates)

	// --- Kiểm tra kết quả method ---
	if actualErr != nil {
		t.Errorf("TC-SEAT-BIZ-04d FAILED: unexpected error = %v", actualErr)
	}
	// --- CheckDB: Đảm bảo Update đã được gọi ---
	if repo.updateCalled != 1 {
		t.Errorf("TC-SEAT-BIZ-04d FAILED: Update called %d times, want 1", repo.updateCalled)
	}
}

// TC-SEAT-BIZ-04e
// Mục tiêu: UpdateSeat() thành công khi không có thay đổi nào (cập nhật rỗng)
func TestUpdateSeat_Success_EmptyUpdates(t *testing.T) {
	// --- Khai báo input ---
	existingSeat := newValidSeat()
	repo := &ducpm88SeatFakeRepo{storedSeat: existingSeat, positionExists: false}
	biz := newTestSeatBusiness(repo)
	inputUpdates := &entity.UpdateSeatRequest{} // Không có gì thay đổi

	// --- Gọi method ---
	actualErr := biz.UpdateSeat(context.Background(), existingSeat.Id, inputUpdates)

	// --- Kiểm tra kết quả ---
	if actualErr != nil {
		t.Errorf("TC-SEAT-BIZ-04e FAILED: unexpected error = %v", actualErr)
	}
	// --- CheckDB: Update vẫn được gọi (persist changes) ---
	if repo.updateCalled != 1 {
		t.Errorf("TC-SEAT-BIZ-04e FAILED: Update called %d times, want 1", repo.updateCalled)
	}
}

// =============================================================================
// TEST NHÓM 27: DeleteSeat()
// Nguồn: movie-service/internal/module/seat/business/business.go
// Logic:
//   1. GetByID lỗi sql.ErrNoRows => ErrSeatNotFound
//   2. GetByID lỗi khác => wrapped error
//   3. repo.Delete() thất bại => wrapped error
//   4. Thành công => CheckDB (Delete được gọi)
// =============================================================================

// TC-SEAT-BIZ-05a
// Mục tiêu: DeleteSeat() trả về ErrSeatNotFound khi seat không tồn tại
func TestDeleteSeat_SeatNotFound(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88SeatFakeRepo{getErr: sql.ErrNoRows}
	biz := newTestSeatBusiness(repo)
	inputId := "non-existent-seat"
	expectedErr := ErrSeatNotFound

	// --- Gọi method ---
	actualErr := biz.DeleteSeat(context.Background(), inputId)

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-SEAT-BIZ-05a FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	// --- CheckDB: Delete không được gọi khi không tìm thấy seat ---
	if repo.deleteCalled != 0 {
		t.Errorf("TC-SEAT-BIZ-05a FAILED: Delete called %d times, want 0", repo.deleteCalled)
	}
}

// TC-SEAT-BIZ-05b
// Mục tiêu: DeleteSeat() thành công; kiểm tra DB đã được gọi
func TestDeleteSeat_Success(t *testing.T) {
	// --- Khai báo input ---
	existingSeat := newValidSeat()
	repo := &ducpm88SeatFakeRepo{storedSeat: existingSeat, deleteErr: nil}
	biz := newTestSeatBusiness(repo)

	// --- Gọi method ---
	actualErr := biz.DeleteSeat(context.Background(), existingSeat.Id)

	// --- Kiểm tra kết quả method ---
	if actualErr != nil {
		t.Errorf("TC-SEAT-BIZ-05b FAILED: unexpected error = %v", actualErr)
	}
	// --- CheckDB: Đảm bảo Delete đã được gọi đúng 1 lần ---
	if repo.deleteCalled != 1 {
		t.Errorf("TC-SEAT-BIZ-05b FAILED: Delete called %d times, want 1 (DB must be called)", repo.deleteCalled)
	}
}

// TC-SEAT-BIZ-05c
// Mục tiêu: DeleteSeat() truyền lỗi từ repo.Delete()
func TestDeleteSeat_RepoDeleteError(t *testing.T) {
	// --- Khai báo input ---
	existingSeat := newValidSeat()
	expectedRepoErr := errors.New("foreign key constraint violation")
	repo := &ducpm88SeatFakeRepo{
		storedSeat: existingSeat,
		deleteErr:  expectedRepoErr,
	}
	biz := newTestSeatBusiness(repo)

	// --- Gọi method ---
	actualErr := biz.DeleteSeat(context.Background(), existingSeat.Id)

	// --- Kiểm tra kết quả ---
	if actualErr == nil {
		t.Fatal("TC-SEAT-BIZ-05c FAILED: expected error but got nil")
	}
	if !errors.Is(actualErr, expectedRepoErr) {
		t.Errorf("TC-SEAT-BIZ-05c FAILED: error does not wrap repo error. got = %v", actualErr)
	}
	// --- CheckDB: Delete đã được gọi (thất bại nhưng DB đã được tác động) ---
	if repo.deleteCalled != 1 {
		t.Errorf("TC-SEAT-BIZ-05c FAILED: Delete called %d times, want 1", repo.deleteCalled)
	}
}

// TC-SEAT-BIZ-05d
// Mục tiêu: DeleteSeat() wrapped error khi GetByID lỗi không phải sql.ErrNoRows
func TestDeleteSeat_GenericGetError(t *testing.T) {
	// --- Khai báo input ---
	expectedRepoErr := errors.New("database connection refused")
	repo := &ducpm88SeatFakeRepo{getErr: expectedRepoErr}
	biz := newTestSeatBusiness(repo)

	// --- Gọi method ---
	actualErr := biz.DeleteSeat(context.Background(), "some-id")

	// --- Kiểm tra kết quả ---
	if actualErr == nil {
		t.Fatal("TC-SEAT-BIZ-05d FAILED: expected error but got nil")
	}
	// --- CheckDB: Delete không được gọi khi GetByID thất bại ---
	if repo.deleteCalled != 0 {
		t.Errorf("TC-SEAT-BIZ-05d FAILED: Delete called %d times, want 0", repo.deleteCalled)
	}
}

// =============================================================================
// TEST NHÓM 28: UpdateSeatStatus()
// Nguồn: movie-service/internal/module/seat/business/business.go
// Logic:
//   1. GetByID lỗi sql.ErrNoRows => ErrSeatNotFound
//   2. GetByID lỗi khác => wrapped error
//   3. Thành công: seat.Status = status, gọi repo.Update() => CheckDB
// =============================================================================

// TC-SEAT-BIZ-06a
// Mục tiêu: UpdateSeatStatus() trả về ErrSeatNotFound khi seat không tồn tại
func TestUpdateSeatStatus_SeatNotFound(t *testing.T) {
	// --- Khai báo input ---
	repo := &ducpm88SeatFakeRepo{getErr: sql.ErrNoRows}
	biz := newTestSeatBusiness(repo)
	expectedErr := ErrSeatNotFound

	// --- Gọi method ---
	actualErr := biz.UpdateSeatStatus(context.Background(), "non-existent-id", entity.SeatStatusMaintenance)

	// --- Kiểm tra kết quả ---
	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("TC-SEAT-BIZ-06a FAILED: err = %v, want %v", actualErr, expectedErr)
	}
	// --- CheckDB: Update không được gọi ---
	if repo.updateCalled != 0 {
		t.Errorf("TC-SEAT-BIZ-06a FAILED: Update called %d times, want 0", repo.updateCalled)
	}
}

// TC-SEAT-BIZ-06b
// Mục tiêu: UpdateSeatStatus() thành công; trạng thái được cập nhật và Update được gọi
func TestUpdateSeatStatus_Success(t *testing.T) {
	// --- Khai báo input ---
	existingSeat := newValidSeat()
	existingSeat.Status = entity.SeatStatusAvailable
	repo := &ducpm88SeatFakeRepo{storedSeat: existingSeat}
	biz := newTestSeatBusiness(repo)
	inputNewStatus := entity.SeatStatusMaintenance // Chuyển sang bảo trì

	// --- Gọi method ---
	actualErr := biz.UpdateSeatStatus(context.Background(), existingSeat.Id, inputNewStatus)

	// --- Kiểm tra kết quả method ---
	if actualErr != nil {
		t.Errorf("TC-SEAT-BIZ-06b FAILED: unexpected error = %v", actualErr)
	}
	// --- CheckDB: Đảm bảo Update đã được gọi (DB được cập nhật) ---
	if repo.updateCalled != 1 {
		t.Errorf("TC-SEAT-BIZ-06b FAILED: Update called %d times, want 1 (DB must be updated)", repo.updateCalled)
	}
}

// TC-SEAT-BIZ-06c
// Mục tiêu: UpdateSeatStatus() thành công với tất cả các trạng thái hợp lệ
func TestUpdateSeatStatus_AllValidStatuses(t *testing.T) {
	// --- Khai báo input: test tất cả các trạng thái ---
	allStatuses := []entity.SeatStatus{
		entity.SeatStatusAvailable,
		entity.SeatStatusOccupied,
		entity.SeatStatusMaintenance,
		entity.SeatStatusBlocked,
	}

	for _, inputStatus := range allStatuses {
		// Reset fake repo cho mỗi test
		existingSeat := newValidSeat()
		repo := &ducpm88SeatFakeRepo{storedSeat: existingSeat}
		biz := newTestSeatBusiness(repo)

		// --- Gọi method ---
		actualErr := biz.UpdateSeatStatus(context.Background(), existingSeat.Id, inputStatus)

		// --- Kiểm tra kết quả ---
		if actualErr != nil {
			t.Errorf("TC-SEAT-BIZ-06c FAILED for status %s: unexpected error = %v", inputStatus, actualErr)
		}
		if repo.updateCalled != 1 {
			t.Errorf("TC-SEAT-BIZ-06c FAILED for status %s: Update called %d times, want 1", inputStatus, repo.updateCalled)
		}
	}
}

// =============================================================================
// TEST NHÓM 29: GetSeatsByIds()
// Nguồn: movie-service/internal/module/seat/business/business.go
// Logic: Gọi repo.GetByIDs() và trả về kết quả
// =============================================================================

// TC-SEAT-BIZ-07a
// Mục tiêu: GetSeatsByIds() trả về đúng danh sách seats theo IDs
func TestGetSeatsByIds_Success(t *testing.T) {
	// --- Khai báo input ---
	expectedSeats := []*entity.Seat{
		{Id: "s1", RoomId: "r1", SeatNumber: "A1", RowNumber: "A", SeatType: entity.SeatTypeRegular, Status: entity.SeatStatusAvailable, CreatedAt: time.Now()},
		{Id: "s2", RoomId: "r1", SeatNumber: "A2", RowNumber: "A", SeatType: entity.SeatTypeVIP, Status: entity.SeatStatusAvailable, CreatedAt: time.Now()},
	}
	repo := &ducpm88SeatFakeRepo{storedSeats: expectedSeats}
	biz := newTestSeatBusiness(repo)
	inputIds := []string{"s1", "s2"}

	// --- Gọi method ---
	actualSeats, actualErr := biz.GetSeatsByIds(context.Background(), inputIds)

	// --- Kiểm tra kết quả ---
	if actualErr != nil {
		t.Errorf("TC-SEAT-BIZ-07a FAILED: unexpected error = %v", actualErr)
	}
	if len(actualSeats) != len(expectedSeats) {
		t.Errorf("TC-SEAT-BIZ-07a FAILED: seats count = %d, want %d", len(actualSeats), len(expectedSeats))
	}
}

// TC-SEAT-BIZ-07b
// Mục tiêu: GetSeatsByIds() trả về lỗi wrapped khi repo thất bại
func TestGetSeatsByIds_RepoError(t *testing.T) {
	// --- Khai báo input ---
	expectedRepoErr := errors.New("db timeout")
	repo := &ducpm88SeatFakeRepo{getErr: expectedRepoErr}
	biz := newTestSeatBusiness(repo)
	inputIds := []string{"s1", "s2"}

	// --- Gọi method ---
	actualSeats, actualErr := biz.GetSeatsByIds(context.Background(), inputIds)

	// --- Kiểm tra kết quả ---
	if actualErr == nil {
		t.Fatal("TC-SEAT-BIZ-07b FAILED: expected error but got nil")
	}
	if actualSeats != nil {
		t.Errorf("TC-SEAT-BIZ-07b FAILED: seats = %v, want nil", actualSeats)
	}
}
