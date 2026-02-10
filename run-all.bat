@echo off
setlocal
set "SCRIPT_DIR=%~dp0"
set "SCRIPT_DIR=%SCRIPT_DIR:~0,-1%"

powershell -NoProfile -Command "& { param($d) [DateTime]::UtcNow.Ticks | Out-File -FilePath (Join-Path $d.TrimEnd('\') 'run-duration-start.txt') }" "%SCRIPT_DIR%"

start /min /wait cmd /c run-prometheus-collector.bat

powershell -NoProfile -Command "& { param($d); $p=$d.TrimEnd('\'); $s=[long](Get-Content (Join-Path $p 'run-duration-start.txt')); $start=[DateTime]::new($s); $end=[DateTime]::UtcNow; $dur=$end-$start; $line=\"$($end.ToString('yyyy-MM-dd HH:mm:ss')) - Duration: $dur\"; Add-Content -Path (Join-Path $p 'run-duration.log') -Value $line; Remove-Item (Join-Path $p 'run-duration-start.txt') -ErrorAction SilentlyContinue }" "%SCRIPT_DIR%"
endlocal
