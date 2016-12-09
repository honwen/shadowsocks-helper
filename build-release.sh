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
LDFLAGS="-X main.version=$VERSION -s -w -linkmode external -extldflags -static"
GCFLAGS=""

OSES=(windows linux darwin freebsd)
ARCHS=(amd64 386)
rm -rf ./release
mkdir -p ./release
for os in ${OSES[@]}; do
	for arch in ${ARCHS[@]}; do
		suffix=""
		if [ "$os" == "windows" ]; then
			suffix=".exe"
			LDFLAGS="-X main.version=$VERSION -s -w"
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
LDFLAGS="-X main.version=$VERSION -s -w"
env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o ./release/ssr-helper_mipsle github.com/chenhw2/shadowsocks-helper
env CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o ./release/ssr-helper_mips github.com/chenhw2/shadowsocks-helper

if $UPX; then upx -9 ./release/ssr-helper_mips**;fi
tar -zcf ./release/ssr-helper_mipsle-$VERSION.tar.gz ./release/ssr-helper_mipsle
tar -zcf ./release/ssr-helper_mips-$VERSION.tar.gz ./release/ssr-helper_mips
$MD5 ./release/ssr-helper_mipsle-$VERSION.tar.gz
