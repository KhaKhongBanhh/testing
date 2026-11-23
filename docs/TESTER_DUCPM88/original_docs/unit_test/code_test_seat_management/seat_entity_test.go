// =============================================================================
// FILE: seat_entity_test.go
// PACKAGE: entity
// MÔ TẢ: Unit test cho các phương thức của Seat entity:
//   - Seat.IsValid()
//   - GetSeatTypePriceMultiplier()
//   - Seat.CalculatePrice()
//   - CreateSeatRequest.ToSeat()
//   - ToSeatResponse()
//   - ToSeatsResponse()
//
// TESTER: DUCPM88 - Phạm Minh Đức - B22DCDT095
// MỤC TIÊU: Đảm bảo phủ cấp 2 (branch coverage) cho tất cả nhánh logic
// LƯU Ý: Không có DB call - thuần túy kiểm tra logic entity
// ĐƯỜNG DẪN KHI COPY: movie-service/internal/module/seat/entity/
// =============================================================================

package entity

import (
	"testing"
	"time"
)

// =============================================================================
// TEST NHÓM 17: Seat.IsValid()
// Nguồn: movie-service/internal/module/seat/entity/seat.go
// Logic: trả về false nếu RoomId, SeatNumber, hoặc RowNumber rỗng
// =============================================================================

// TC-SEAT-ENTITY-01a
// Mục tiêu: IsValid() trả về true khi tất cả các trường bắt buộc có giá trị
func TestSeatIsValid_AllFieldsProvided(t *testing.T) {
	// --- Khai báo input ---
	inputSeat := &Seat{
		RoomId:     "room-001",
		SeatNumber: "A1",
		RowNumber:  "A",
	}
	expectedResult := true

	// --- Gọi method ---
	actualResult := inputSeat.IsValid()

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-SEAT-ENTITY-01a FAILED: IsValid() = %v, want %v", actualResult, expectedResult)
	}
}

// TC-SEAT-ENTITY-01b
// Mục tiêu: IsValid() trả về false khi RoomId rỗng
func TestSeatIsValid_EmptyRoomId(t *testing.T) {
	// --- Khai báo input ---
	inputSeat := &Seat{
		RoomId:     "", // RoomId bị bỏ trống
		SeatNumber: "A1",
		RowNumber:  "A",
	}
	expectedResult := false

	// --- Gọi method ---
	actualResult := inputSeat.IsValid()

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-SEAT-ENTITY-01b FAILED: IsValid() = %v, want %v", actualResult, expectedResult)
	}
}

// TC-SEAT-ENTITY-01c
// Mục tiêu: IsValid() trả về false khi SeatNumber rỗng
func TestSeatIsValid_EmptySeatNumber(t *testing.T) {
	// --- Khai báo input ---
	inputSeat := &Seat{
		RoomId:     "room-001",
		SeatNumber: "", // SeatNumber bị bỏ trống
		RowNumber:  "A",
	}
	expectedResult := false

	// --- Gọi method ---
	actualResult := inputSeat.IsValid()

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-SEAT-ENTITY-01c FAILED: IsValid() = %v, want %v", actualResult, expectedResult)
	}
}

// TC-SEAT-ENTITY-01d
// Mục tiêu: IsValid() trả về false khi RowNumber rỗng
func TestSeatIsValid_EmptyRowNumber(t *testing.T) {
	// --- Khai báo input ---
	inputSeat := &Seat{
		RoomId:     "room-001",
		SeatNumber: "A1",
		RowNumber:  "", // RowNumber bị bỏ trống
	}
	expectedResult := false

	// --- Gọi method ---
	actualResult := inputSeat.IsValid()

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-SEAT-ENTITY-01d FAILED: IsValid() = %v, want %v", actualResult, expectedResult)
	}
}

// TC-SEAT-ENTITY-01e
// Mục tiêu: IsValid() trả về false khi tất cả các trường bắt buộc rỗng
func TestSeatIsValid_AllFieldsEmpty(t *testing.T) {
	// --- Khai báo input ---
	inputSeat := &Seat{
		RoomId:     "",
		SeatNumber: "",
		RowNumber:  "",
	}
	expectedResult := false

	// --- Gọi method ---
	actualResult := inputSeat.IsValid()

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-SEAT-ENTITY-01e FAILED: IsValid() = %v, want %v", actualResult, expectedResult)
	}
}

// =============================================================================
// TEST NHÓM 18: GetSeatTypePriceMultiplier()
// Nguồn: movie-service/internal/module/seat/entity/seat.go
// Logic:
//   REGULAR => 1.0
//   VIP     => 1.5
//   COUPLE  => 2.5
//   default => 1.0
// =============================================================================

// TC-SEAT-ENTITY-02a
// Mục tiêu: GetSeatTypePriceMultiplier() trả về 1.0 cho REGULAR
func TestGetSeatTypePriceMultiplier_Regular(t *testing.T) {
	// --- Khai báo input ---
	inputSeatType := SeatTypeRegular
	expectedMultiplier := 1.0

	// --- Gọi function ---
	actualMultiplier := GetSeatTypePriceMultiplier(inputSeatType)

	// --- Kiểm tra kết quả ---
	if actualMultiplier != expectedMultiplier {
		t.Errorf("TC-SEAT-ENTITY-02a FAILED: multiplier = %v, want %v", actualMultiplier, expectedMultiplier)
	}
}

// TC-SEAT-ENTITY-02b
// Mục tiêu: GetSeatTypePriceMultiplier() trả về 1.5 cho VIP
func TestGetSeatTypePriceMultiplier_VIP(t *testing.T) {
	// --- Khai báo input ---
	inputSeatType := SeatTypeVIP
	expectedMultiplier := 1.5

	// --- Gọi function ---
	actualMultiplier := GetSeatTypePriceMultiplier(inputSeatType)

	// --- Kiểm tra kết quả ---
	if actualMultiplier != expectedMultiplier {
		t.Errorf("TC-SEAT-ENTITY-02b FAILED: multiplier = %v, want %v", actualMultiplier, expectedMultiplier)
	}
}

// TC-SEAT-ENTITY-02c
// Mục tiêu: GetSeatTypePriceMultiplier() trả về 2.5 cho COUPLE
func TestGetSeatTypePriceMultiplier_Couple(t *testing.T) {
	// --- Khai báo input ---
	inputSeatType := SeatTypeCouple
	expectedMultiplier := 2.5

	// --- Gọi function ---
	actualMultiplier := GetSeatTypePriceMultiplier(inputSeatType)

	// --- Kiểm tra kết quả ---
	if actualMultiplier != expectedMultiplier {
		t.Errorf("TC-SEAT-ENTITY-02c FAILED: multiplier = %v, want %v", actualMultiplier, expectedMultiplier)
	}
}

// TC-SEAT-ENTITY-02d
// Mục tiêu: GetSeatTypePriceMultiplier() trả về 1.0 cho loại ghế không xác định (default)
func TestGetSeatTypePriceMultiplier_UnknownType(t *testing.T) {
	// --- Khai báo input ---
	inputSeatType := SeatType("PREMIUM") // Loại không có trong hệ thống
	expectedMultiplier := 1.0            // Default case

	// --- Gọi function ---
	actualMultiplier := GetSeatTypePriceMultiplier(inputSeatType)

	// --- Kiểm tra kết quả ---
	if actualMultiplier != expectedMultiplier {
		t.Errorf("TC-SEAT-ENTITY-02d FAILED: multiplier = %v, want %v", actualMultiplier, expectedMultiplier)
	}
}

// =============================================================================
// TEST NHÓM 19: Seat.CalculatePrice()
// Nguồn: movie-service/internal/module/seat/entity/seat.go
// Logic: basePrice * GetSeatTypePriceMultiplier(seat.SeatType)
// =============================================================================

// TC-SEAT-ENTITY-03a
// Mục tiêu: CalculatePrice() tính đúng giá cho ghế REGULAR
func TestCalculatePrice_RegularSeat(t *testing.T) {
	// --- Khai báo input ---
	inputSeat := &Seat{SeatType: SeatTypeRegular}
	inputBasePrice := 100000.0 // 100,000 VND
	expectedPrice := 100000.0  // 100,000 * 1.0 = 100,000

	// --- Gọi method ---
	actualPrice := inputSeat.CalculatePrice(inputBasePrice)

	// --- Kiểm tra kết quả ---
	if actualPrice != expectedPrice {
		t.Errorf("TC-SEAT-ENTITY-03a FAILED: price = %v, want %v", actualPrice, expectedPrice)
	}
}

// TC-SEAT-ENTITY-03b
// Mục tiêu: CalculatePrice() tính đúng giá cho ghế VIP (x1.5)
func TestCalculatePrice_VIPSeat(t *testing.T) {
	// --- Khai báo input ---
	inputSeat := &Seat{SeatType: SeatTypeVIP}
	inputBasePrice := 100000.0 // 100,000 VND
	expectedPrice := 150000.0  // 100,000 * 1.5 = 150,000

	// --- Gọi method ---
	actualPrice := inputSeat.CalculatePrice(inputBasePrice)

	// --- Kiểm tra kết quả ---
	if actualPrice != expectedPrice {
		t.Errorf("TC-SEAT-ENTITY-03b FAILED: price = %v, want %v", actualPrice, expectedPrice)
	}
}

// TC-SEAT-ENTITY-03c
// Mục tiêu: CalculatePrice() tính đúng giá cho ghế COUPLE (x2.5)
func TestCalculatePrice_CoupleSeat(t *testing.T) {
	// --- Khai báo input ---
	inputSeat := &Seat{SeatType: SeatTypeCouple}
	inputBasePrice := 100000.0 // 100,000 VND
	expectedPrice := 250000.0  // 100,000 * 2.5 = 250,000

	// --- Gọi method ---
	actualPrice := inputSeat.CalculatePrice(inputBasePrice)

	// --- Kiểm tra kết quả ---
	if actualPrice != expectedPrice {
		t.Errorf("TC-SEAT-ENTITY-03c FAILED: price = %v, want %v", actualPrice, expectedPrice)
	}
}

// TC-SEAT-ENTITY-03d
// Mục tiêu: CalculatePrice() trả về 0 khi basePrice = 0
func TestCalculatePrice_ZeroBasePrice(t *testing.T) {
	// --- Khai báo input ---
	inputSeat := &Seat{SeatType: SeatTypeVIP}
	inputBasePrice := 0.0
	expectedPrice := 0.0 // 0 * bất kỳ = 0

	// --- Gọi method ---
	actualPrice := inputSeat.CalculatePrice(inputBasePrice)

	// --- Kiểm tra kết quả ---
	if actualPrice != expectedPrice {
		t.Errorf("TC-SEAT-ENTITY-03d FAILED: price = %v, want %v", actualPrice, expectedPrice)
	}
}

// TC-SEAT-ENTITY-03e
// Mục tiêu: CalculatePrice() hoạt động đúng với số thập phân
func TestCalculatePrice_DecimalBasePrice(t *testing.T) {
	// --- Khai báo input ---
	inputSeat := &Seat{SeatType: SeatTypeVIP}
	inputBasePrice := 50000.5  // 50,000.5 VND
	expectedPrice := 75000.75  // 50,000.5 * 1.5 = 75,000.75

	// --- Gọi method ---
	actualPrice := inputSeat.CalculatePrice(inputBasePrice)

	// --- Kiểm tra kết quả ---
	if actualPrice != expectedPrice {
		t.Errorf("TC-SEAT-ENTITY-03e FAILED: price = %v, want %v", actualPrice, expectedPrice)
	}
}

// =============================================================================
// TEST NHÓM 20: CreateSeatRequest.ToSeat()
// Nguồn: movie-service/internal/module/seat/entity/types.go
// Logic: Chuyển đổi request sang Seat entity với status mặc định AVAILABLE
// =============================================================================

// TC-SEAT-ENTITY-04a
// Mục tiêu: ToSeat() chuyển đổi đúng request sang Seat entity
func TestCreateSeatRequestToSeat_AllFieldsMapped(t *testing.T) {
	// --- Khai báo input ---
	inputRequest := &CreateSeatRequest{
		RoomId:     "room-uuid-001",
		SeatNumber: "B5",
		RowNumber:  "B",
		SeatType:   SeatTypeVIP,
	}
	expectedStatus := SeatStatusAvailable // Status mặc định phải là AVAILABLE

	// --- Gọi method ---
	actualSeat := inputRequest.ToSeat()

	// --- Kiểm tra từng trường ---
	if actualSeat.RoomId != inputRequest.RoomId {
		t.Errorf("TC-SEAT-ENTITY-04a FAILED: RoomId = %v, want %v", actualSeat.RoomId, inputRequest.RoomId)
	}
	if actualSeat.SeatNumber != inputRequest.SeatNumber {
		t.Errorf("TC-SEAT-ENTITY-04a FAILED: SeatNumber = %v, want %v", actualSeat.SeatNumber, inputRequest.SeatNumber)
	}
	if actualSeat.RowNumber != inputRequest.RowNumber {
		t.Errorf("TC-SEAT-ENTITY-04a FAILED: RowNumber = %v, want %v", actualSeat.RowNumber, inputRequest.RowNumber)
	}
	if actualSeat.SeatType != inputRequest.SeatType {
		t.Errorf("TC-SEAT-ENTITY-04a FAILED: SeatType = %v, want %v", actualSeat.SeatType, inputRequest.SeatType)
	}
	// --- Kiểm tra Status mặc định ---
	if actualSeat.Status != expectedStatus {
		t.Errorf("TC-SEAT-ENTITY-04a FAILED: Status = %v, want %v", actualSeat.Status, expectedStatus)
	}
}

// TC-SEAT-ENTITY-04b
// Mục tiêu: ToSeat() luôn gán status AVAILABLE bất kể loại ghế
func TestCreateSeatRequestToSeat_DefaultStatusAlwaysAvailable(t *testing.T) {
	// --- Khai báo input ---
	requestRegular := &CreateSeatRequest{RoomId: "r1", SeatNumber: "A1", RowNumber: "A", SeatType: SeatTypeRegular}
	requestVIP := &CreateSeatRequest{RoomId: "r1", SeatNumber: "B1", RowNumber: "B", SeatType: SeatTypeVIP}
	requestCouple := &CreateSeatRequest{RoomId: "r1", SeatNumber: "C1", RowNumber: "C", SeatType: SeatTypeCouple}
	expectedStatus := SeatStatusAvailable

	// --- Gọi method cho từng loại ghế ---
	seatRegular := requestRegular.ToSeat()
	seatVIP := requestVIP.ToSeat()
	seatCouple := requestCouple.ToSeat()

	// --- Kiểm tra kết quả ---
	if seatRegular.Status != expectedStatus {
		t.Errorf("TC-SEAT-ENTITY-04b FAILED: REGULAR seat status = %v, want %v", seatRegular.Status, expectedStatus)
	}
	if seatVIP.Status != expectedStatus {
		t.Errorf("TC-SEAT-ENTITY-04b FAILED: VIP seat status = %v, want %v", seatVIP.Status, expectedStatus)
	}
	if seatCouple.Status != expectedStatus {
		t.Errorf("TC-SEAT-ENTITY-04b FAILED: COUPLE seat status = %v, want %v", seatCouple.Status, expectedStatus)
	}
}

// =============================================================================
// TEST NHÓM 21: ToSeatResponse()
// Nguồn: movie-service/internal/module/seat/entity/types.go
// Logic:
//   - Nếu seat nil => trả về nil
//   - Chuyển đổi CreatedAt sang chuỗi RFC3339
//   - Nếu UpdatedAt != nil => điền trường UpdatedAt trong response
// =============================================================================

// TC-SEAT-ENTITY-05a
// Mục tiêu: ToSeatResponse() trả về nil khi seat là nil
func TestToSeatResponse_NilSeat(t *testing.T) {
	// --- Khai báo input ---
	var inputSeat *Seat = nil

	// --- Gọi function ---
	actualResponse := ToSeatResponse(inputSeat)

	// --- Kiểm tra kết quả ---
	if actualResponse != nil {
		t.Errorf("TC-SEAT-ENTITY-05a FAILED: response = %v, want nil", actualResponse)
	}
}

// TC-SEAT-ENTITY-05b
// Mục tiêu: ToSeatResponse() chuyển đổi đúng tất cả các trường khi seat hợp lệ
func TestToSeatResponse_ValidSeat_NoUpdatedAt(t *testing.T) {
	// --- Khai báo input ---
	now := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	inputSeat := &Seat{
		Id:         "seat-001",
		RoomId:     "room-001",
		SeatNumber: "A1",
		RowNumber:  "A",
		SeatType:   SeatTypeRegular,
		Status:     SeatStatusAvailable,
		CreatedAt:  now,
		UpdatedAt:  nil, // Chưa có cập nhật
	}
	expectedCreatedAt := "2025-01-15T10:30:00Z" // RFC3339 format

	// --- Gọi function ---
	actualResponse := ToSeatResponse(inputSeat)

	// --- Kiểm tra kết quả ---
	if actualResponse == nil {
		t.Fatal("TC-SEAT-ENTITY-05b FAILED: response is nil")
	}
	if actualResponse.Id != inputSeat.Id {
		t.Errorf("TC-SEAT-ENTITY-05b FAILED: Id = %v, want %v", actualResponse.Id, inputSeat.Id)
	}
	if actualResponse.RoomId != inputSeat.RoomId {
		t.Errorf("TC-SEAT-ENTITY-05b FAILED: RoomId = %v, want %v", actualResponse.RoomId, inputSeat.RoomId)
	}
	if actualResponse.SeatNumber != inputSeat.SeatNumber {
		t.Errorf("TC-SEAT-ENTITY-05b FAILED: SeatNumber = %v, want %v", actualResponse.SeatNumber, inputSeat.SeatNumber)
	}
	if actualResponse.SeatType != inputSeat.SeatType {
		t.Errorf("TC-SEAT-ENTITY-05b FAILED: SeatType = %v, want %v", actualResponse.SeatType, inputSeat.SeatType)
	}
	if actualResponse.Status != inputSeat.Status {
		t.Errorf("TC-SEAT-ENTITY-05b FAILED: Status = %v, want %v", actualResponse.Status, inputSeat.Status)
	}
	if actualResponse.CreatedAt != expectedCreatedAt {
		t.Errorf("TC-SEAT-ENTITY-05b FAILED: CreatedAt = %v, want %v", actualResponse.CreatedAt, expectedCreatedAt)
	}
	if actualResponse.UpdatedAt != nil {
		t.Errorf("TC-SEAT-ENTITY-05b FAILED: UpdatedAt = %v, want nil", actualResponse.UpdatedAt)
	}
}

// TC-SEAT-ENTITY-05c
// Mục tiêu: ToSeatResponse() điền UpdatedAt khi seat có thời gian cập nhật
func TestToSeatResponse_WithUpdatedAt(t *testing.T) {
	// --- Khai báo input ---
	createdAt := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	updatedAt := time.Date(2025, 2, 1, 15, 0, 0, 0, time.UTC)
	inputSeat := &Seat{
		Id:         "seat-002",
		RoomId:     "room-001",
		SeatNumber: "B2",
		RowNumber:  "B",
		SeatType:   SeatTypeVIP,
		Status:     SeatStatusOccupied,
		CreatedAt:  createdAt,
		UpdatedAt:  &updatedAt, // Có thời gian cập nhật
	}
	expectedUpdatedAt := "2025-02-01T15:00:00Z"

	// --- Gọi function ---
	actualResponse := ToSeatResponse(inputSeat)

	// --- Kiểm tra kết quả ---
	if actualResponse.UpdatedAt == nil {
		t.Fatal("TC-SEAT-ENTITY-05c FAILED: UpdatedAt is nil, want non-nil")
	}
	if *actualResponse.UpdatedAt != expectedUpdatedAt {
		t.Errorf("TC-SEAT-ENTITY-05c FAILED: UpdatedAt = %v, want %v", *actualResponse.UpdatedAt, expectedUpdatedAt)
	}
}

// =============================================================================
// TEST NHÓM 22: ToSeatsResponse()
// Nguồn: movie-service/internal/module/seat/entity/types.go
// Logic: Chuyển đổi danh sách Seat sang SeatsResponse với paging info
// =============================================================================

// TC-SEAT-ENTITY-06a
// Mục tiêu: ToSeatsResponse() chuyển đổi đúng danh sách seats và tính paging
func TestToSeatsResponse_MultipleSeats(t *testing.T) {
	// --- Khai báo input ---
	now := time.Now()
	inputSeats := []*Seat{
		{Id: "s1", RoomId: "r1", SeatNumber: "A1", RowNumber: "A", SeatType: SeatTypeRegular, Status: SeatStatusAvailable, CreatedAt: now},
		{Id: "s2", RoomId: "r1", SeatNumber: "A2", RowNumber: "A", SeatType: SeatTypeVIP, Status: SeatStatusOccupied, CreatedAt: now},
		{Id: "s3", RoomId: "r1", SeatNumber: "B1", RowNumber: "B", SeatType: SeatTypeCouple, Status: SeatStatusAvailable, CreatedAt: now},
	}
	inputPage := 1
	inputSize := 10
	inputTotal := 3

	// --- Gọi function ---
	actualResponse := ToSeatsResponse(inputSeats, inputPage, inputSize, inputTotal)

	// --- Kiểm tra kết quả ---
	if len(actualResponse.Data) != len(inputSeats) {
		t.Errorf("TC-SEAT-ENTITY-06a FAILED: data count = %d, want %d", len(actualResponse.Data), len(inputSeats))
	}
	if actualResponse.Paging == nil {
		t.Fatal("TC-SEAT-ENTITY-06a FAILED: paging is nil")
	}
	if actualResponse.Paging.Page != inputPage {
		t.Errorf("TC-SEAT-ENTITY-06a FAILED: Paging.Page = %d, want %d", actualResponse.Paging.Page, inputPage)
	}
	if actualResponse.Paging.Total != inputTotal {
		t.Errorf("TC-SEAT-ENTITY-06a FAILED: Paging.Total = %d, want %d", actualResponse.Paging.Total, inputTotal)
	}
}

// TC-SEAT-ENTITY-06b
// Mục tiêu: ToSeatsResponse() trả về slice rỗng (không nil) khi không có seats
func TestToSeatsResponse_EmptySeats(t *testing.T) {
	// --- Khai báo input ---
	inputSeats := []*Seat{}
	inputPage := 1
	inputSize := 10
	inputTotal := 0

	// --- Gọi function ---
	actualResponse := ToSeatsResponse(inputSeats, inputPage, inputSize, inputTotal)

	// --- Kiểm tra kết quả ---
	if actualResponse.Data == nil {
		t.Error("TC-SEAT-ENTITY-06b FAILED: Data should not be nil (should be empty slice)")
	}
	if len(actualResponse.Data) != 0 {
		t.Errorf("TC-SEAT-ENTITY-06b FAILED: Data count = %d, want 0", len(actualResponse.Data))
	}
	if actualResponse.Paging.Total != 0 {
		t.Errorf("TC-SEAT-ENTITY-06b FAILED: Paging.Total = %d, want 0", actualResponse.Paging.Total)
	}
}
