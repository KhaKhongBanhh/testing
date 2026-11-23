// =============================================================================
// FILE: paging_test.go
// PACKAGE: paging
// MÔ TẢ: Unit test cho package paging:
//   - NewPageInfo(): tạo PageInfo từ page, size, total
//   - stringToInt(): chuyển chuỗi sang số nguyên (unexported, test trong cùng package)
//
// TESTER: DUCPM88 - Phạm Minh Đức - B22DCDT095
// MỤC TIÊU: Phủ cấp 2 cho cả 2 nhánh của NewPageInfo (chia dư / không chia dư)
// LƯU Ý: Không có DB call, không có network call
// ĐƯỜNG DẪN KHI COPY: movie-service/internal/pkg/paging/
// =============================================================================

package paging

import "testing"

// =============================================================================
// TEST NHÓM 15: NewPageInfo()
// Nguồn: movie-service/internal/pkg/paging/paging.go
// Logic:
//   totalPages = total / size
//   if total%size > 0: totalPages++ (ceiling division)
// =============================================================================

// TC-PAGING-01a
// Mục tiêu: NewPageInfo() tính đúng totalPages khi tổng chia hết cho size
func TestNewPageInfo_ExactDivision(t *testing.T) {
	// --- Khai báo input ---
	inputPage := 1
	inputSize := 10
	inputTotal := 20        // 20 / 10 = 2 chính xác, không dư
	expectedTotalPages := 2 // Phải bằng 2

	// --- Gọi method ---
	actualPageInfo := NewPageInfo(inputPage, inputSize, inputTotal)

	// --- Kiểm tra kết quả ---
	if actualPageInfo.TotalPages != expectedTotalPages {
		t.Errorf("TC-PAGING-01a FAILED: TotalPages = %d, want %d", actualPageInfo.TotalPages, expectedTotalPages)
	}
	if actualPageInfo.Page != inputPage {
		t.Errorf("TC-PAGING-01a FAILED: Page = %d, want %d", actualPageInfo.Page, inputPage)
	}
	if actualPageInfo.Size != inputSize {
		t.Errorf("TC-PAGING-01a FAILED: Size = %d, want %d", actualPageInfo.Size, inputSize)
	}
	if actualPageInfo.Total != inputTotal {
		t.Errorf("TC-PAGING-01a FAILED: Total = %d, want %d", actualPageInfo.Total, inputTotal)
	}
}

// TC-PAGING-01b
// Mục tiêu: NewPageInfo() tính đúng totalPages khi tổng KHÔNG chia hết (ceiling division)
func TestNewPageInfo_CeilingDivision(t *testing.T) {
	// --- Khai báo input ---
	inputPage := 2
	inputSize := 10
	inputTotal := 25        // 25 / 10 = 2.5 => làm tròn lên = 3
	expectedTotalPages := 3 // Phải làm tròn lên

	// --- Gọi method ---
	actualPageInfo := NewPageInfo(inputPage, inputSize, inputTotal)

	// --- Kiểm tra kết quả ---
	if actualPageInfo.TotalPages != expectedTotalPages {
		t.Errorf("TC-PAGING-01b FAILED: TotalPages = %d, want %d", actualPageInfo.TotalPages, expectedTotalPages)
	}
}

// TC-PAGING-01c
// Mục tiêu: NewPageInfo() khi total = 0 => totalPages = 0 (không tăng lên vì 0%size = 0)
func TestNewPageInfo_ZeroTotal(t *testing.T) {
	// --- Khai báo input ---
	inputPage := 1
	inputSize := 10
	inputTotal := 0
	expectedTotalPages := 0 // 0/10=0, 0%10=0 => không increment

	// --- Gọi method ---
	actualPageInfo := NewPageInfo(inputPage, inputSize, inputTotal)

	// --- Kiểm tra kết quả ---
	if actualPageInfo.TotalPages != expectedTotalPages {
		t.Errorf("TC-PAGING-01c FAILED: TotalPages = %d, want %d", actualPageInfo.TotalPages, expectedTotalPages)
	}
	if actualPageInfo.Total != 0 {
		t.Errorf("TC-PAGING-01c FAILED: Total = %d, want 0", actualPageInfo.Total)
	}
}

// TC-PAGING-01d
// Mục tiêu: NewPageInfo() khi total < size => totalPages = 1 (ceiling(3/10) = 1)
func TestNewPageInfo_TotalLessThanSize(t *testing.T) {
	// --- Khai báo input ---
	inputPage := 1
	inputSize := 10
	inputTotal := 3         // Chỉ có 3 items nhưng page size là 10
	expectedTotalPages := 1 // 3/10=0, 3%10=3>0 => totalPages=0+1=1

	// --- Gọi method ---
	actualPageInfo := NewPageInfo(inputPage, inputSize, inputTotal)

	// --- Kiểm tra kết quả ---
	if actualPageInfo.TotalPages != expectedTotalPages {
		t.Errorf("TC-PAGING-01d FAILED: TotalPages = %d, want %d", actualPageInfo.TotalPages, expectedTotalPages)
	}
}

// TC-PAGING-01e
// Mục tiêu: NewPageInfo() khi size = 1 => totalPages bằng total
func TestNewPageInfo_SizeOne(t *testing.T) {
	// --- Khai báo input ---
	inputPage := 5
	inputSize := 1
	inputTotal := 7
	expectedTotalPages := 7 // 7/1=7, 7%1=0 => totalPages=7

	// --- Gọi method ---
	actualPageInfo := NewPageInfo(inputPage, inputSize, inputTotal)

	// --- Kiểm tra kết quả ---
	if actualPageInfo.TotalPages != expectedTotalPages {
		t.Errorf("TC-PAGING-01e FAILED: TotalPages = %d, want %d", actualPageInfo.TotalPages, expectedTotalPages)
	}
	if actualPageInfo.Page != inputPage {
		t.Errorf("TC-PAGING-01e FAILED: Page = %d, want %d", actualPageInfo.Page, inputPage)
	}
}

// =============================================================================
// TEST NHÓM 16: stringToInt() (unexported function)
// Nguồn: movie-service/internal/pkg/paging/paging.go
// Logic:
//   - Chuỗi hợp lệ => trả về số nguyên
//   - Chuỗi không hợp lệ => trả về 0
// Có thể test vì cùng package (package paging)
// =============================================================================

// TC-PAGING-02a
// Mục tiêu: stringToInt() chuyển đổi đúng chuỗi số hợp lệ
func TestStringToInt_ValidNumber(t *testing.T) {
	// --- Khai báo input ---
	inputStr := "42"
	expectedResult := 42

	// --- Gọi function ---
	actualResult := stringToInt(inputStr)

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-PAGING-02a FAILED: stringToInt(%q) = %d, want %d", inputStr, actualResult, expectedResult)
	}
}

// TC-PAGING-02b
// Mục tiêu: stringToInt() trả về 0 khi chuỗi không phải số
func TestStringToInt_InvalidString(t *testing.T) {
	// --- Khai báo input ---
	inputStr := "abc"
	expectedResult := 0 // Lỗi parse => trả về 0

	// --- Gọi function ---
	actualResult := stringToInt(inputStr)

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-PAGING-02b FAILED: stringToInt(%q) = %d, want %d", inputStr, actualResult, expectedResult)
	}
}

// TC-PAGING-02c
// Mục tiêu: stringToInt() trả về 0 khi chuỗi rỗng
func TestStringToInt_EmptyString(t *testing.T) {
	// --- Khai báo input ---
	inputStr := ""
	expectedResult := 0

	// --- Gọi function ---
	actualResult := stringToInt(inputStr)

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-PAGING-02c FAILED: stringToInt(%q) = %d, want %d", inputStr, actualResult, expectedResult)
	}
}

// TC-PAGING-02d
// Mục tiêu: stringToInt() chuyển đổi đúng chuỗi số âm
func TestStringToInt_NegativeNumber(t *testing.T) {
	// --- Khai báo input ---
	inputStr := "-5"
	expectedResult := -5

	// --- Gọi function ---
	actualResult := stringToInt(inputStr)

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-PAGING-02d FAILED: stringToInt(%q) = %d, want %d", inputStr, actualResult, expectedResult)
	}
}

// TC-PAGING-02e
// Mục tiêu: stringToInt() trả về 0 khi chuỗi là số thập phân (float)
func TestStringToInt_FloatString(t *testing.T) {
	// --- Khai báo input ---
	inputStr := "3.14"
	expectedResult := 0 // strconv.Atoi không parse được float => trả về 0

	// --- Gọi function ---
	actualResult := stringToInt(inputStr)

	// --- Kiểm tra kết quả ---
	if actualResult != expectedResult {
		t.Errorf("TC-PAGING-02e FAILED: stringToInt(%q) = %d, want %d", inputStr, actualResult, expectedResult)
	}
}
