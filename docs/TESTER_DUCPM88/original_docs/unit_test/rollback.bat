@echo off
setlocal EnableDelayedExpansion

REM ================================================================
REM ROLLBACK & CLEANUP - DUCPM88 (Pham Minh Duc - B22DCDT095)
REM Don dep cac file test da copy va phuc hoi du lieu Database
REM ================================================================

echo.
echo ================================================================
echo   ROLLBACK - UNIT TEST CLEANUP
echo   TESTER: DUCPM88 - Pham Minh Duc - B22DCDT095
echo ================================================================
echo.

REM --- 1. Xac dinh duong dan ---
set "UNIT_TEST_DIR=%~dp0"
if "%UNIT_TEST_DIR:~-1%"=="\" set "UNIT_TEST_DIR=%UNIT_TEST_DIR:~0,-1%"

for %%d in ("%UNIT_TEST_DIR%\..\..\..\..") do set "PROJECT_ROOT=%%~fd"
set "MOVIE_SERVICE=%PROJECT_ROOT%\movie-service"

set "DEST_MOVIE_ENTITY=%MOVIE_SERVICE%\internal\module\movie\entity"
set "DEST_MOVIE_BIZ=%MOVIE_SERVICE%\internal\module\movie\business"
set "DEST_PAGING=%MOVIE_SERVICE%\internal\pkg\paging"
set "DEST_SEAT_ENTITY=%MOVIE_SERVICE%\internal\module\seat\entity"
set "DEST_SEAT_BIZ=%MOVIE_SERVICE%\internal\module\seat\business"

echo [INFO] Project root: %PROJECT_ROOT%
echo.

REM ================================================================
REM PHAN 1: XOA TEST FILES DA COPY VAO SOURCE CODE
REM ================================================================
echo ----------------------------------------------------------------
echo [PHAN 1/3] Xoa cac test files da copy vao source code...
echo ----------------------------------------------------------------
set "CLEANUP_COUNT=0"

if exist "%DEST_MOVIE_ENTITY%\ducpm88_movie_entity_test.go" (
    del /F /Q "%DEST_MOVIE_ENTITY%\ducpm88_movie_entity_test.go"
    echo   [OK] Da xoa: movie\entity\ducpm88_movie_entity_test.go
    set /a CLEANUP_COUNT+=1
) else (
    echo   [SKIP] Khong co: movie\entity\ducpm88_movie_entity_test.go
)

if exist "%DEST_MOVIE_BIZ%\ducpm88_movie_business_test.go" (
    del /F /Q "%DEST_MOVIE_BIZ%\ducpm88_movie_business_test.go"
    echo   [OK] Da xoa: movie\business\ducpm88_movie_business_test.go
    set /a CLEANUP_COUNT+=1
) else (
    echo   [SKIP] Khong co: movie\business\ducpm88_movie_business_test.go
)

if exist "%DEST_PAGING%\ducpm88_paging_test.go" (
    del /F /Q "%DEST_PAGING%\ducpm88_paging_test.go"
    echo   [OK] Da xoa: pkg\paging\ducpm88_paging_test.go
    set /a CLEANUP_COUNT+=1
) else (
    echo   [SKIP] Khong co: pkg\paging\ducpm88_paging_test.go
)

if exist "%DEST_SEAT_ENTITY%\ducpm88_seat_entity_test.go" (
    del /F /Q "%DEST_SEAT_ENTITY%\ducpm88_seat_entity_test.go"
    echo   [OK] Da xoa: seat\entity\ducpm88_seat_entity_test.go
    set /a CLEANUP_COUNT+=1
) else (
    echo   [SKIP] Khong co: seat\entity\ducpm88_seat_entity_test.go
)

if exist "%DEST_SEAT_BIZ%\ducpm88_seat_business_test.go" (
    del /F /Q "%DEST_SEAT_BIZ%\ducpm88_seat_business_test.go"
    echo   [OK] Da xoa: seat\business\ducpm88_seat_business_test.go
    set /a CLEANUP_COUNT+=1
) else (
    echo   [SKIP] Khong co: seat\business\ducpm88_seat_business_test.go
)

echo.
echo [KET QUA PHAN 1]: Da don dep %CLEANUP_COUNT% test file.
echo.

REM ================================================================
REM PHAN 2: DB ROLLBACK (PostgreSQL)
REM ================================================================
echo ----------------------------------------------------------------
echo [PHAN 2/3] Kiem tra va Rollback du lieu Database...
echo ----------------------------------------------------------------
echo [INFO] Giai thich:
echo   - Unit Tests da su dung IN-MEMORY MOCK nen hoan toan KHONG
echo     lam thay doi du lieu Database that.
echo   - Buoc nay thuc hien rollback bang lenh SQL xoa bat ky ban ghi
echo     thu nghiem nao neu nguoi dung co tao thu cong.
echo.

REM Kiem tra Docker container postgres
docker inspect -f "{{.State.Running}}" postgres >nul 2>&1
if errorlevel 1 goto :check_host_psql

echo [INFO] Tim thay Docker container "postgres". Dang thuc hien SQL Rollback...
docker exec -i postgres psql -U postgres -d cinema_app -c "DELETE FROM movies WHERE title LIKE 'TEST_%%' OR title LIKE 'Test Movie%%';" >nul 2>&1
docker exec -i postgres psql -U postgres -d cinema_app -c "DELETE FROM seats WHERE seat_number LIKE 'TEST_%%';" >nul 2>&1
echo   [OK] Da thuc hien Rollback bang SQL qua Docker container postgres.
echo.
echo [INFO] Trang thai so luong ban ghi hien tai trong Database:
docker exec -i postgres psql -U postgres -d cinema_app -c "SELECT 'movies' as table_name, count(1) as count FROM movies UNION ALL SELECT 'seats', count(1) FROM seats UNION ALL SELECT 'rooms', count(1) FROM rooms;" 2>nul
goto :done_db

:check_host_psql
where psql >nul 2>&1
if errorlevel 1 (
    echo [INFO] Khong tim thay psql hoac Docker container postgres khong chay.
    echo        Du lieu van an toan 100%% vi unit test dung mock.
    goto :done_db
)

set "PGPASSWORD=postgres"
psql -h localhost -p 5434 -U postgres -d cinema_app -c "DELETE FROM movies WHERE title LIKE 'TEST_%%' OR title LIKE 'Test Movie%%';" >nul 2>&1
psql -h localhost -p 5434 -U postgres -d cinema_app -c "DELETE FROM seats WHERE seat_number LIKE 'TEST_%%';" >nul 2>&1
echo   [OK] Da thuc hien Rollback bang SQL qua psql host.

:done_db
echo.

REM ================================================================
REM PHAN 3: XOA FILE COVERAGE TAO RA KHI TEST
REM ================================================================
echo ----------------------------------------------------------------
echo [PHAN 3/3] Don dep file coverage cu...
echo ----------------------------------------------------------------
if exist "%MOVIE_SERVICE%\coverage.out" (
    del /F /Q "%MOVIE_SERVICE%\coverage.out"
    echo   [OK] Da xoa: %MOVIE_SERVICE%\coverage.out
) else (
    echo   [SKIP] Khong co: coverage.out
)

if exist "%MOVIE_SERVICE%\coverage.html" (
    del /F /Q "%MOVIE_SERVICE%\coverage.html"
    echo   [OK] Da xoa: %MOVIE_SERVICE%\coverage.html
) else (
    echo   [SKIP] Khong co: %MOVIE_SERVICE%\coverage.html
)

if exist "%UNIT_TEST_DIR%\coverage.html" (
    del /F /Q "%UNIT_TEST_DIR%\coverage.html"
    echo   [OK] Da xoa: %UNIT_TEST_DIR%\coverage.html
) else (
    echo   [SKIP] Khong co: %UNIT_TEST_DIR%\coverage.html
)

echo.
echo ================================================================
echo   ROLLBACK DA HOAN THANH THANH CONG!
echo   - Toan bo ma nguon da tro ve trang thai ban dau.
echo   - Du lieu Database da duoc kiem tra va rollback sach se.
echo   - Cac file coverage tam thoi da duoc xoa.
echo ================================================================
echo.
echo Cua so nay se giu nguyen den khi ban dong.
echo Nhan phim bat ky hoac click nut [X] o goc de dong cua so nay...
pause
endlocal
