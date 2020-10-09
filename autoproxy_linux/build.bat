
call :build_all darwin 386
call :build_all darwin amd64

call :build_all linux 386
call :build_all linux amd64

call :build_all linux arm
call :build_all linux arm64

exit /b 0

:build_all
    set GOOS=%1
    set GOARCH=%2

    echo build %GOOS% %GOARCH%

    mkdir output
	
	copy client.yaml output\
	copy server.yaml output\
	
	copy start.sh output\
	
	go build -ldflags="-w -s" -o output\autoproxy .

	cd output
    tar -zcf ../autoproxy_%GOOS%_%GOARCH%.tar.gz *
	cd ..
	
	rmdir /q/s output
	
goto :eof

