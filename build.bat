
call :build_all linux 386
call :build_all linux amd64

call :build_all linux arm
call :build_all linux arm64

call :build_all windows 386 .exe
call :build_all windows amd64 .exe


exit /b 0

:build_all
    set GOOS=%1
    set GOARCH=%2
    set TAG=%3

    echo build %GOOS% %GOARCH%

    mkdir output
	
	copy client.yaml output\
	copy server.yaml output\
	
	go build -ldflags="-w -s" -o output\autoproxy%TAG% .

	cd output
    tar -zcf ../%GOOS%_%GOARCH%.tar.gz *
	cd ..
	
	rmdir /q/s output
	
goto :eof

