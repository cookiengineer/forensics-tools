#!/bin/bash

GO="$(which go)";
ROOT="$(pwd)";


rstrip() {
    # Usage: rstrip "string" "pattern"
    printf '%s\n' "${1%%$2}"
}

build() {

	local os="$1";
	local arch="$2";
	local project="$3";
	local cmd="$4";
	local cmd_ext="";
	local cmd_file=$(basename -- "${cmd}");
	local source_folder="${ROOT}/${project}";
	local build_folder="${ROOT}/build";

	if [[ "$os" == "windows" ]]; then
		cmd_ext="exe";
	fi;


	cd "${source_folder}";

	if [[ ! -e "${build_folder}" ]]; then
		mkdir -p "${build_folder}";
	fi;

	if [[ "${cmd_ext}" != "" ]]; then
		env CGO_ENABLED=0 GOOS="${os}" GOARCH="${arch}" ${GO} build -o "${build_folder}/${cmd_file}_${os}-${arch}.${cmd_ext}" "${cmd}/main.go";
	else
		env CGO_ENABLED=0 GOOS="${os}" GOARCH="${arch}" ${GO} build -o "${build_folder}/${cmd_file}_${os}-${arch}" "${cmd}/main.go";
	fi;

	if [[ "$?" == "0" ]]; then
		echo -e "- Build ${cmd} / ${os} / ${arch} [\e[32mok\e[0m]";
	else
		echo -e "- Build ${cmd} / ${os} / ${arch} [\e[31mfail\e[0m]";
	fi;

}



if [[ "$GO" != "" ]]; then

	build "linux" "amd64" "dns" "cmds/dns-searchcensored";

	build "linux" "amd64" "http" "cmds/simpleserver";

	build "linux" "amd64" "sql" "cmds/sql-extract";
	build "linux" "amd64" "sql" "cmds/sql-tables";

	build "linux" "amd64" "torrent" "cmds/magnetify";

	build "linux" "amd64" "totp" "cmds/totp-export";

	build "linux" "amd64" "crx" "cmds/uncrx";

else
	echo "Please install go(lang) compiler.";
	exit 1;
fi;

