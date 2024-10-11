export CGO_ENABLED=1
export WDIR=${PWD}

all: linux windows

linux: clean
	cp -r lib/Linux/ build/ 
	GOOS=linux  CGO_CFLAGS="-I${WDIR}/include"  CGO_LDFLAGS="-L${WDIR}/build -Wl,-rpath=${WDIR}/build -lhcnetsdk" go build -ldflags "-s -w" -o build/hik main.go
	build/hik

windows: clean
	cp -r lib/Windows/ build/
	CGO_CFLAGS="-I$(CURDIR)/include" CGO_LDFLAGS="-L$(CURDIR)/build -lHCNetSDK -lHCCore"  go build -ldflags "-s -w" -o build/hik.exe main.go
	build/hik.exe

clean:
	rm -r build
	rm *.jpeg

