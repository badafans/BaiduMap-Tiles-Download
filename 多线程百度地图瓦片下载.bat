@echo off
setlocal enabledelayedexpansion
cd "%~dp0"
title �ٶȵ�ͼ��Ƭ����
echo ��ͼ���
echo ���� normal
echo ������ light
echo ��ҹ dark
echo ��ɫ���� redalert
echo ���� googlelite
echo ��Ȼ�� grassgreen
echo ��ҹ�� midnight
echo ������ pink
echo �ഺ�� darkgreen
echo �������� bluish
echo �߶˻� grayscale
echo ǿ�߽� hardedge
set style=midnight
set a=120.86441707191187
set b=31.87829317767047
set c=122.28637881960425
set d=30.68334074991772
set /a min=3
set /a max=19
set /p style=��ͼ���(Ĭ�� %style%):
set /p a=�ٶȵ�ͼ���ϽǾ���(Ĭ�� %a%):
set /p b=�ٶȵ�ͼ���Ͻ�γ��(Ĭ�� %b%):
set /p c=�ٶȵ�ͼ���½Ǿ���(Ĭ�� %c%):
set /p d=�ٶȵ�ͼ���½�γ��(Ĭ�� %d%):
set /p min=��С�㼶(Ĭ�� %min%):
set /p max=���㼶(Ĭ�� %max%):
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
if !tasklist! LEQ 30 (echo %time% ��ǰcurl����!tasklist!��) else (echo %time% ��ǰcurl����!tasklist!�� && goto loop)
)
) else (
title �� !n! ����Ƭ �㼶 !z! ��Ƭ X=!x1! ��Ƭ Y=!y2! ��ͼ��� !style!
)
goto loopy
:loopx
if !x1! LSS !x2! (set /a x1+=1 && mkdir map\!z!\!x1! && goto sety1) else (goto start)
:end
for /f "delims=" %%i in ('tasklist ^| find /c /i "curl.exe"') do (
set tasklist=%%i
if !tasklist! EQU 0 (echo curl����ȫ������) else (echo %time% �ȴ�curl����ȫ������,��ǰʣ�����!tasklist!�� && goto end)
)
set /a stopH=%time:~0,2%
if %time:~3,1% EQU 0 (set /a stopM=%time:~4,1%) else (set /a stopM=%time:~3,2%)
if %time:~6,1% EQU 0 (set /a stopS=%time:~7,1%) else (set /a stopS=%time:~6,2%)
set /a starttime=%startH%*3600+%startM%*60+%startS%
set /a stoptime=%stopH%*3600+%stopM%*60+%stopS%
if %starttime% GTR %stoptime% (set /a alltime=86400-%starttime%+%stoptime%) else (set /a alltime=%stoptime%-%starttime%)
set /a avg=n/alltime
echo �ܼ�����%n%����Ƭͼ,��ʱ%alltime%��,ƽ���ٶ�%avg%��/s
pause>nul