#!/bin/bash

GO="$(which go)";
ROOT="$(pwd)";



build() {

	local os="$1";
	local arch="$2";
	local folder="${ROOT}/build";

	local ext="";

	if [[ "$os" == "windows" ]]; then
		ext="exe";
	fi;

	mkdir -p "$folder";

	cd "${ROOT}";

	if [[ "$ext" != "" ]]; then
		env CGO_ENABLED=0 GOOS="${os}" GOARCH="${arch}" ${GO} build -o "${folder}/uncrx_${os}-${arch}.${ext}" "${ROOT}/cmd/uncrx.go";
	else
		env CGO_ENABLED=0 GOOS="${os}" GOARCH="${arch}" ${GO} build -o "${folder}/uncrx_${os}-${arch}" "${ROOT}/cmd/uncrx.go";
	fi;

	if [[ "$?" == "0" ]]; then
		echo -e "- Build $os / $arch [\e[32mok\e[0m]";
	else
		echo -e "- Build $os / $arch [\e[31mfail\e[0m]";
	fi;

}



if [[ "$GO" != "" ]]; then

	# build "android" "arm";
	# build "android" "arm64";

	build "darwin" "amd64";
	build "darwin" "arm64";

	build "linux" "386";
	build "linux" "amd64";
	build "linux" "arm";
	build "linux" "arm64";

	build "windows" "386";
	build "windows" "amd64";
	build "windows" "arm";
	build "windows" "arm64";

else
	echo "Please install go(lang) compiler.";
	exit 1;
fi;

