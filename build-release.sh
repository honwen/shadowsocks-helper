#!/bin/bash
MD5='md5sum'
unamestr=`uname`
if [[ "$unamestr" == 'Darwin' ]]; then
	MD5='md5'
fi

UPX=false
if hash upx 2>/dev/null; then
	UPX=true
fi

VERSION=`date -u +%Y%m%d`
LDFLAGS="-X main.version=$VERSION -s -w"
GCFLAGS=""

OSES=(linux darwin windows freebsd)
ARCHS=(amd64 386)
rm -rf ./release
mkdir -p ./release
for os in ${OSES[@]}; do
	for arch in ${ARCHS[@]}; do
		suffix=""
		if [ "$os" == "windows" ]; then
			suffix=".exe"
		fi
		LDFLAGS="-X main.version=$VERSION -s -w"
		if [ "$os" == "linux" ]; then
			LDFLAGS="${LDFLAGS} -linkmode external -extldflags -static"
		fi
		env CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o ./release/ssr-helper_${os}_${arch}${suffix} github.com/chenhw2/shadowsocks-helper
		if $UPX; then upx -9 ./release/ssr-helper_${os}_${arch}${suffix};fi
		tar -zcf ./release/ssr-helper_${os}-${arch}-$VERSION.tar.gz ./release/ssr-helper_${os}_${arch}${suffix}
		$MD5 ./release/ssr-helper_${os}-${arch}-$VERSION.tar.gz
	done
done

# ARM
ARMS=(5 6 7)
for v in ${ARMS[@]}; do
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=$v go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o ./release/ssr-helper_arm$v  github.com/chenhw2/shadowsocks-helper
done
if $UPX; then upx -9 ./release/ssr-helper_arm*;fi
tar -zcf ./release/ssr-helper_arm-$VERSION.tar.gz ./release/ssr-helper_arm*
$MD5 ./release/ssr-helper_arm-$VERSION.tar.gz

#MIPS32LE
env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o ./release/ssr-helper_mipsle github.com/chenhw2/shadowsocks-helper
env CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o ./release/ssr-helper_mipsle github.com/chenhw2/shadowsocks-helper

if $UPX; then upx -9 client_linux_mips* server_linux_mips*;fi
tar -zcf ./release/ssr-helper_mipsle-$VERSION.tar.gz ./release/ssr-helper_mipsle
$MD5 ./release/ssr-helper_mipsle-$VERSION.tar.gz
