@echo off
cd /d "c:\Users\harsh\Python\StockTrack"
(
    echo Timestamp: %date% %time%
    echo.
    echo === CHECKING GIT STATUS ===
    git status --short
    echo.
    echo === ADDING FILES ===
    git add -A
    echo Added files successfully
    echo.
    echo === COMMITTING ===
    git commit -m "Complete: Login screen, CORS middleware, comprehensive API docs, real-time debugging"
    echo.
    echo === PUSHING TO GITHUB ===
    git push origin main
    echo.
    echo === COMPLETE ===
    echo Push finished at: %date% %time%
) > push_output.txt 2>&1
type push_output.txt
