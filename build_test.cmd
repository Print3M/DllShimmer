@echo off
go build -o DllShimmer.exe --ldflags "-w -s" -trimpath main.go 
copy "%windir%\System32\version.dll" version.dll
mkdir shimload
DllShimmer -i version.dll -o shimload --original-path version.dll -m --static
choice /C Y /N /D Y /T 4 >nul
