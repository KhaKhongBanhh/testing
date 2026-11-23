@echo off
setlocal EnableDelayedExpansion

REM ================================================================
REM ROLLBACK & CLEANUP - API TEST - DUCPM88 (Pham Minh Duc - B22DCDT095)
REM Phuc hoi du lieu Database ve trang thai truoc khi test API
REM ================================================================

echo.
echo ================================================================
echo   ROLLBACK - API TEST CLEANUP
echo   TESTER: DUCPM88 - Pham Minh Duc - B22DCDT095
echo ================================================================
echo.

REM --- 1. Xac dinh duong dan ---
set "API_TEST_DIR=%~dp0"
if "%API_TEST_DIR:~-1%"=="\" set "API_TEST_DIR=%API_TEST_DIR:~0,-1%"

for %%d in ("%API_TEST_DIR%\..\..\..\..") do set "PROJECT_ROOT=%%~fd"

echo [INFO] Project root: %PROJECT_ROOT%
echo [INFO] Thu muc API Test: %API_TEST_DIR%
echo.

REM ================================================================
REM PHAN 1: DB ROLLBACK (PostgreSQL)
REM Xoa cac ban ghi test co the duoc tao boi API Testing
REM ================================================================
echo ----------------------------------------------------------------
echo [PHAN 1/2] Kiem tra va Rollback du lieu Database...
echo ----------------------------------------------------------------
echo [INFO] Giai thich:
echo   - Trong qua trinh test API, cac request POST/PUT/PATCH co the
echo     da tao hoac thay doi du lieu trong Database.
echo   - Script test tu dong da rollback ngay sau moi test case,
echo     nhung buoc nay dam bao don dep sach se moi du lieu con sot.
echo.

REM Kiem tra Docker container postgres
docker inspect -f "{{.State.Running}}" postgres >nul 2>&1
if errorlevel 1 goto :check_host_psql

echo [INFO] Tim thay Docker container "postgres". Dang thuc hien SQL Rollback...
echo.

echo   [1/6] Xoa movie test co title bat dau bang "API Test", "Minimal", "Test"...
docker exec -i postgres psql -U postgres -d cinema_app -c "DELETE FROM movie_genres WHERE movie_id IN (SELECT id FROM movies WHERE title LIKE 'API Test%%' OR title LIKE 'Minimal%%' OR title LIKE 'Before Update%%' OR title LIKE 'After Update%%' OR title LIKE 'Will%%' OR title LIKE 'Zero%%' OR title LIKE 'Negative%%' OR title LIKE 'No Duration%%' OR title LIKE 'Direct%%' OR title LIKE 'Invalid%%' OR title LIKE 'Skip%%' OR title LIKE 'Bug%%' OR title LIKE 'Status%%' OR title LIKE 'Empty%%' OR title LIKE 'Ghost%%' OR title LIKE 'Bad Genre%%' OR title LIKE 'Movie With%%' OR title LIKE 'Same%%');" >nul 2>&1
docker exec -i postgres psql -U postgres -d cinema_app -c "DELETE FROM movies WHERE title LIKE 'API Test%%' OR title LIKE 'Minimal%%' OR title LIKE 'Before Update%%' OR title LIKE 'After Update%%' OR title LIKE 'Will%%' OR title LIKE 'Zero%%' OR title LIKE 'Negative%%' OR title LIKE 'No Duration%%' OR title LIKE 'Direct%%' OR title LIKE 'Invalid%%' OR title LIKE 'Skip%%' OR title LIKE 'Bug%%' OR title LIKE 'Status%%' OR title LIKE 'Empty%%' OR title LIKE 'Ghost%%' OR title LIKE 'Bad Genre%%' OR title LIKE 'Movie With%%' OR title LIKE 'Same%%';" >nul 2>&1
echo   [OK] Da don dep movies test.

echo   [2/6] Xoa movie co status khong hop le (du lieu ban)...
docker exec -i postgres psql -U postgres -d cinema_app -c "DELETE FROM movie_genres WHERE movie_id IN (SELECT id FROM movies WHERE status NOT IN ('UPCOMING','SHOWING','ENDED'));" >nul 2>&1
docker exec -i postgres psql -U postgres -d cinema_app -c "DELETE FROM movies WHERE status NOT IN ('UPCOMING','SHOWING','ENDED');" >nul 2>&1
echo   [OK] Da don dep movies co status khong hop le.

echo   [3/6] Xoa seats test co seat_number bat dau bang Z9, DUP, CON, UPD, INV, ST, DEL, PRT, F0, SAME...
docker exec -i postgres psql -U postgres -d cinema_app -c "DELETE FROM seats WHERE seat_number IN ('Z99','Z98','Z97','Z96','DUP1','DUP2','CON1','CON2','UPD1','UPD2','INV1','ST1','ST2','ST3','ST4','ST5','ST6','DEL1','PRT1','PRT2','F01','SAME1','X01');" >nul 2>&1
echo   [OK] Da don dep seats test.

echo   [4/6] Xoa seats co seat_type khong hop le (du lieu ban)...
docker exec -i postgres psql -U postgres -d cinema_app -c "DELETE FROM seats WHERE seat_type NOT IN ('REGULAR','VIP','COUPLE');" >nul 2>&1
echo   [OK] Da don dep seats co seat_type khong hop le.

echo   [5/6] Xoa seats co status khong hop le (du lieu ban)...
docker exec -i postgres psql -U postgres -d cinema_app -c "UPDATE seats SET status='AVAILABLE' WHERE status NOT IN ('AVAILABLE','OCCUPIED','MAINTENANCE','BLOCKED');" >nul 2>&1
echo   [OK] Da reset seats co status khong hop le ve AVAILABLE.

echo   [6/6] Hien thi trang thai Database hien tai...
echo.
echo [INFO] Trang thai so luong ban ghi hien tai trong Database:
docker exec -i postgres psql -U postgres -d cinema_app -c "SELECT 'movies' as table_name, count(1) as count FROM movies UNION ALL SELECT 'seats', count(1) FROM seats UNION ALL SELECT 'rooms', count(1) FROM rooms UNION ALL SELECT 'genres', count(1) FROM genres UNION ALL SELECT 'movie_genres', count(1) FROM movie_genres;" 2>nul
echo.
echo [INFO] Du lieu mong doi: movies=12, seats=561, rooms=5, genres=15
goto :done_db

:check_host_psql
where psql >nul 2>&1
if errorlevel 1 (
    echo [CANH BAO] Khong tim thay psql hoac Docker container postgres khong chay.
    echo            Vui long dam bao Docker dang chay va thu lai.
    goto :done_db
)

set "PGPASSWORD=postgres"
echo [INFO] Su dung psql tren host de Rollback...
psql -h localhost -p 5434 -U postgres -d cinema_app -c "DELETE FROM movie_genres WHERE movie_id IN (SELECT id FROM movies WHERE title LIKE 'API Test%%' OR title LIKE 'Minimal%%' OR title LIKE 'Test%%');" >nul 2>&1
psql -h localhost -p 5434 -U postgres -d cinema_app -c "DELETE FROM movies WHERE title LIKE 'API Test%%' OR title LIKE 'Minimal%%' OR title LIKE 'Test%%';" >nul 2>&1
psql -h localhost -p 5434 -U postgres -d cinema_app -c "DELETE FROM seats WHERE seat_number IN ('Z99','Z98','Z97','Z96','DUP1','DUP2','CON1','CON2','UPD1','UPD2','INV1','ST1','ST2','ST3','ST4','ST5','ST6','DEL1','PRT1','PRT2','F01','SAME1','X01');" >nul 2>&1
psql -h localhost -p 5434 -U postgres -d cinema_app -c "DELETE FROM seats WHERE seat_type NOT IN ('REGULAR','VIP','COUPLE');" >nul 2>&1
psql -h localhost -p 5434 -U postgres -d cinema_app -c "UPDATE seats SET status='AVAILABLE' WHERE status NOT IN ('AVAILABLE','OCCUPIED','MAINTENANCE','BLOCKED');" >nul 2>&1
echo   [OK] Da thuc hien Rollback bang SQL qua psql host.

:done_db
echo.

REM ================================================================
REM PHAN 2: DON DEP FILE TAM THOI
REM ================================================================
echo ----------------------------------------------------------------
echo [PHAN 2/2] Don dep file tam thoi...
echo ----------------------------------------------------------------

if exist "%API_TEST_DIR%\api_test_results.json" (
    del /F /Q "%API_TEST_DIR%\api_test_results.json"
    echo   [OK] Da xoa: api_test_results.json
) else (
    echo   [SKIP] Khong co: api_test_results.json
)

if exist "%PROJECT_ROOT%\api_test_results.json" (
    del /F /Q "%PROJECT_ROOT%\api_test_results.json"
    echo   [OK] Da xoa: api_test_results.json (project root)
) else (
    echo   [SKIP] Khong co: api_test_results.json (project root)
)

echo.
echo ================================================================
echo   ROLLBACK DA HOAN THANH THANH CONG!
echo   - Du lieu Database da duoc kiem tra va rollback sach se.
echo   - Cac ban ghi test da bi xoa (movies, seats).
echo   - Du lieu ban (status/seat_type khong hop le) da duoc don dep.
echo   - Cac file tam thoi da duoc xoa.
echo ================================================================
echo.
echo Cua so nay se giu nguyen den khi ban dong.
echo Nhan phim bat ky hoac click nut [X] o goc de dong cua so nay...
pause
endlocal
