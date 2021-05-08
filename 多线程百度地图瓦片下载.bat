@echo off
setlocal enabledelayedexpansion
cd "%~dp0"
title 百度地图瓦片爬虫
echo 地图风格
echo 常规 normal
echo 清新蓝 light
echo 黑夜 dark
echo 红色警戒 redalert
echo 精简 googlelite
echo 自然绿 grassgreen
echo 午夜蓝 midnight
echo 浪漫粉 pink
echo 青春绿 darkgreen
echo 清新蓝绿 bluish
echo 高端灰 grayscale
echo 强边界 hardedge
set style=midnight
set a=120.86441707191187
set b=31.87829317767047
set c=122.28637881960425
set d=30.68334074991772
set /a min=3
set /a max=19
set /p style=地图风格(默认 %style%):
set /p a=百度地图左上角经度(默认 %a%):
set /p b=百度地图左上角纬度(默认 %b%):
set /p c=百度地图右下角经度(默认 %c%):
set /p d=百度地图右下角纬度(默认 %d%):
set /p min=最小层级(默认 %min%):
set /p max=最大层级(默认 %max%):
set /a startH=%time:~0,2%
if %time:~3,1% EQU 0 (set /a startM=%time:~4,1%) else (set /a startM=%time:~3,2%)
if %time:~6,1% EQU 0 (set /a startS=%time:~7,1%) else (set /a startS=%time:~6,2%)
RD /S /Q map
mkdir map
cls
set /a n=0
set /a m=0
:start
if !min! LEQ !max! (set /a z=!min!) else (goto end)
mkdir map\!min!
set /a min+=1
:setx1
for /f "tokens=1 delims= " %%i in ('getid -x %a% -y %b% -z !z!') do (
set x1=%%i
mkdir map\!z!\!x1!
)
:sety1
for /f "tokens=2 delims= " %%i in ('getid -x %a% -y %b% -z !z!') do (
set y1=%%i
)
:setx2
for /f "tokens=1 delims= " %%i in ('getid -x %c% -y %d% -z !z!') do (
set x2=%%i
)
:sety2
for /f "tokens=2 delims= " %%i in ('getid -x %c% -y %d% -z !z!') do (
set y2=%%i
)
:loopy
if !y1! GEQ !y2! (start /b getpic !x1! !y2! !z! !style! && set /a n+=1 && set /a y2+=1) else (goto loopx)
set /a m+=1
if !m! EQU 100 (
set /a m=0
:loop
for /f "delims=" %%i in ('tasklist ^| find /c /i "curl.exe"') do (
set tasklist=%%i
if !tasklist! LEQ 30 (echo %time% 当前curl进程!tasklist!个) else (echo %time% 当前curl进程!tasklist!个 && goto loop)
)
) else (
title 第 !n! 个瓦片 层级 !z! 瓦片 X=!x1! 瓦片 Y=!y2! 地图风格 !style!
)
goto loopy
:loopx
if !x1! LSS !x2! (set /a x1+=1 && mkdir map\!z!\!x1! && goto sety1) else (goto start)
:end
for /f "delims=" %%i in ('tasklist ^| find /c /i "curl.exe"') do (
set tasklist=%%i
if !tasklist! EQU 0 (echo curl进程全部结束) else (echo %time% 等待curl进程全部结束,当前剩余进程!tasklist!个 && goto end)
)
set /a stopH=%time:~0,2%
if %time:~3,1% EQU 0 (set /a stopM=%time:~4,1%) else (set /a stopM=%time:~3,2%)
if %time:~6,1% EQU 0 (set /a stopS=%time:~7,1%) else (set /a stopS=%time:~6,2%)
set /a starttime=%startH%*3600+%startM%*60+%startS%
set /a stoptime=%stopH%*3600+%stopM%*60+%stopS%
if %starttime% GTR %stoptime% (set /a alltime=86400-%starttime%+%stoptime%) else (set /a alltime=%stoptime%-%starttime%)
set /a avg=n/alltime
echo 总计下载%n%个瓦片图,用时%alltime%秒,平均速度%avg%个/s
pause>nul