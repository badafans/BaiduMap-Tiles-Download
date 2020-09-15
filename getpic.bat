@echo off
setlocal enabledelayedexpansion
cd "%~dp0"
:loop
curl "https://api.map.baidu.com/customimage/tile?&x=%1&y=%2&z=%3&customid=%4" -s -o map/%3/%1/%2.png
if not exist "map/%3/%1/%2.png" goto loop
exit