@echo off
cd /d "c:\Users\harsh\Python\StockTrack"
echo Checking git status...
git status --short
echo.
echo Adding all files...
git add -A
echo All files added
echo.
echo Committing changes...
git commit -m "Complete: Login screen, CORS middleware, comprehensive API docs, real-time debug logging"
echo.
echo Pushing to GitHub...
git push origin main
echo.
echo Push complete!
pause
