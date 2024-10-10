export CGO_ENABLED=1
export WDIR=${PWD}

all: linux windows

linux:
	cp -r lib/Linux/ build/ 
	GOOS=linux  CGO_CFLAGS="-I${WDIR}/include"  CGO_LDFLAGS="-L${WDIR}/build -Wl,-rpath=${WDIR}/build -lhcnetsdk" go build -ldflags "-s -w" -o build/hik main.go
	build/hik

windows:
	cp -r lib/Windows/ build/
	CGO_LDFLAGS_ALLOW=".*" CGO_CFLAGS="-I$(CURDIR)/include" CGO_LDFLAGS="-L$(CURDIR)/build -Wl,--enable-stdcall-fixup,-rpath=$(CURDIR)/build -lHCNetSDK -lHCCore" GOOS=windows CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -ldflags "-s -w" -o build/hik.exe main.go
	build/hik.exe

clean:
	rm -r build

