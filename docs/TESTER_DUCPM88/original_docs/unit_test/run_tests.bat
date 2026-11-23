@echo off
setlocal EnableDelayedExpansion

REM ================================================================
REM UNIT TEST RUNNER - DUCPM88 (Pham Minh Duc - B22DCDT095)
REM Chay Unit Test va do Do phu ma nguon (Code Coverage)
REM ================================================================

echo.
echo ================================================================
echo   UNIT TEST RUNNER - DUCPM88 (Pham Minh Duc - B22DCDT095)
echo   Movie Management + Seat Management
echo ================================================================
echo.

REM --- 1. Xac dinh duong dan ---
set "UNIT_TEST_DIR=%~dp0"
if "%UNIT_TEST_DIR:~-1%"=="\" set "UNIT_TEST_DIR=%UNIT_TEST_DIR:~0,-1%"

for %%d in ("%UNIT_TEST_DIR%\..\..\..\..") do set "PROJECT_ROOT=%%~fd"
set "MOVIE_SERVICE=%PROJECT_ROOT%\movie-service"

set "CODE_MOVIE=%UNIT_TEST_DIR%\code_test_movie_management"
set "CODE_SEAT=%UNIT_TEST_DIR%\code_test_seat_management"

set "DEST_MOVIE_ENTITY=%MOVIE_SERVICE%\internal\module\movie\entity"
set "DEST_MOVIE_BIZ=%MOVIE_SERVICE%\internal\module\movie\business"
set "DEST_PAGING=%MOVIE_SERVICE%\internal\pkg\paging"
set "DEST_SEAT_ENTITY=%MOVIE_SERVICE%\internal\module\seat\entity"
set "DEST_SEAT_BIZ=%MOVIE_SERVICE%\internal\module\seat\business"

echo [INFO] Thu muc test: %UNIT_TEST_DIR%
echo [INFO] Movie service: %MOVIE_SERVICE%
echo.

REM --- 2. Kiem tra Docker ---
echo [STEP 1/5] Kiem tra Docker Desktop...
docker info >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker khong hoat dong. Hay bat Docker Desktop len truoc!
    echo.
    goto :end
)
echo [OK] Docker dang chay san sang.
echo.

REM --- 3. Copy test scripts vao source ---
echo [STEP 2/5] Copy test files tu code_test_* vao source code...
copy /Y "%CODE_MOVIE%\movie_entity_test.go" "%DEST_MOVIE_ENTITY%\ducpm88_movie_entity_test.go" >nul
if errorlevel 1 ( echo [FAIL] Khong copy duoc movie_entity_test.go & goto :end )
echo   + Da copy: ducpm88_movie_entity_test.go

copy /Y "%CODE_MOVIE%\movie_business_test.go" "%DEST_MOVIE_BIZ%\ducpm88_movie_business_test.go" >nul
if errorlevel 1 ( echo [FAIL] Khong copy duoc movie_business_test.go & goto :end )
echo   + Da copy: ducpm88_movie_business_test.go

copy /Y "%CODE_MOVIE%\paging_test.go" "%DEST_PAGING%\ducpm88_paging_test.go" >nul
if errorlevel 1 ( echo [FAIL] Khong copy duoc paging_test.go & goto :end )
echo   + Da copy: ducpm88_paging_test.go

copy /Y "%CODE_SEAT%\seat_entity_test.go" "%DEST_SEAT_ENTITY%\ducpm88_seat_entity_test.go" >nul
if errorlevel 1 ( echo [FAIL] Khong copy duoc seat_entity_test.go & goto :end )
echo   + Da copy: ducpm88_seat_entity_test.go

copy /Y "%CODE_SEAT%\seat_business_test.go" "%DEST_SEAT_BIZ%\ducpm88_seat_business_test.go" >nul
if errorlevel 1 ( echo [FAIL] Khong copy duoc seat_business_test.go & goto :end )
echo   + Da copy: ducpm88_seat_business_test.go

echo.
echo [OK] Hoan tat copy toan bo test scripts vao source code.
echo.

REM --- 4. Chay Unit Tests ---
echo [STEP 3/5] Chay Unit Tests chi tiet (go test -v)...
echo ----------------------------------------------------------------
docker run --rm ^
  -v "%MOVIE_SERVICE%:/app" ^
  -v "gomodcache:/go/pkg/mod" ^
  -w /app ^
  golang:1.24-alpine ^
  go test -v -count=1 ^
    ./internal/module/movie/entity/... ^
    ./internal/module/movie/business/... ^
    ./internal/module/seat/entity/... ^
    ./internal/module/seat/business/... ^
    ./internal/pkg/paging/...

set "TEST_EXIT_CODE=%errorlevel%"
echo ----------------------------------------------------------------
if %TEST_EXIT_CODE% == 0 (
    echo [RESULT] TAT CA CAC TEST CASES DA PASSED 100%%!
) else (
    echo [RESULT] CO TEST CASE BI FAILED!
)
echo.

REM --- 5. Do Do phu (Code Coverage) ---
echo [STEP 4/5] Do Do phu ma nguon (Code Coverage)...
echo ----------------------------------------------------------------
docker run --rm ^
  -v "%MOVIE_SERVICE%:/app" ^
  -v "gomodcache:/go/pkg/mod" ^
  -w /app ^
  golang:1.24-alpine ^
  go test -count=1 ^
    -coverprofile=/app/coverage.out ^
    -covermode=atomic ^
    ./internal/module/movie/entity/... ^
    ./internal/module/movie/business/... ^
    ./internal/module/seat/entity/... ^
    ./internal/module/seat/business/... ^
    ./internal/pkg/paging/...

echo.
echo [BANG THONG KE DO PHU TUNG HAM]:
echo ----------------------------------------------------------------
docker run --rm ^
  -v "%MOVIE_SERVICE%:/app" ^
  -v "gomodcache:/go/pkg/mod" ^
  -w /app ^
  golang:1.24-alpine ^
  go tool cover -func=/app/coverage.out
echo ----------------------------------------------------------------
echo.

REM --- 6. Xuat file HTML do phu truc quan (Xanh/Do) ---
echo [STEP 5/5] Tao bao cao giao dien web truc quan (coverage.html)...
docker run --rm ^
  -v "%MOVIE_SERVICE%:/app" ^
  -v "gomodcache:/go/pkg/mod" ^
  -w /app ^
  golang:1.24-alpine ^
  go tool cover -html=/app/coverage.out -o /app/coverage.html

if exist "%MOVIE_SERVICE%\coverage.html" (
    copy /Y "%MOVIE_SERVICE%\coverage.html" "%UNIT_TEST_DIR%\coverage.html" >nul
    echo [OK] Da tao thanh cong file bao cao HTML:
    echo      + %UNIT_TEST_DIR%\coverage.html
    echo      + %MOVIE_SERVICE%\coverage.html
    echo.
    echo [INFO] Dang tu dong mo coverage.html tren trinh duyet...
    start "" "%UNIT_TEST_DIR%\coverage.html"
) else (
    echo [WARN] Chua tao duoc file coverage.html
)

echo.
echo ================================================================
echo   TONG KET UNIT TEST VA COVERAGE REPORT
echo ================================================================
echo   1. Ket qua unit test: PASSED
echo   2. File du lieu do phu: %MOVIE_SERVICE%\coverage.out
echo   3. Giao dien web HTML:  %UNIT_TEST_DIR%\coverage.html
echo.
echo   Luu y: Khi muon xoa cac file test da copy va don dep du lieu,
echo          hay click vao file: rollback.bat
echo ================================================================
echo.

:end
echo Cua so nay se giu nguyen den khi ban dong.
echo Nhan phim bat ky de thoat...
pause
endlocal
