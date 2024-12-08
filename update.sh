#!/usr/bin/env bash
# Written in [Amber](https://amber-lang.com/)
# version: 0.3.5-alpha
# date: 2024-12-08 23:41:27


is_command__76_v0() {
    local command=$1
    [ -x "$(command -v ${command})" ];
    __AS=$?;
if [ $__AS != 0 ]; then
        __AF_is_command76_v0=0;
        return 0
fi
    __AF_is_command76_v0=1;
    return 0
}
exit__80_v0() {
    local code=$1
    exit "${code}";
    __AS=$?
}
download__118_v0() {
    local url=$1
    local path=$2
    is_command__76_v0 "curl";
    __AF_is_command76_v0__10_9="$__AF_is_command76_v0";
    is_command__76_v0 "wget";
    __AF_is_command76_v0__13_9="$__AF_is_command76_v0";
    is_command__76_v0 "aria2c";
    __AF_is_command76_v0__16_9="$__AF_is_command76_v0";
    if [ "$__AF_is_command76_v0__10_9" != 0 ]; then
        curl -L -o "${path}" "${url}";
        __AS=$?
elif [ "$__AF_is_command76_v0__13_9" != 0 ]; then
        wget "${url}" -P "${path}";
        __AS=$?
elif [ "$__AF_is_command76_v0__16_9" != 0 ]; then
        aria2c "${url}" -d "${path}";
        __AS=$?
else
        __AF_download118_v0=0;
        return 0
fi
    __AF_download118_v0=1;
    return 0
}
__0_bool_true=1
__1_bool_false=0
download_latest_package__126_v0() {
    local target=$1
    echo "run download_latest_package"
    download__118_v0 "${target}" "./SgridCloudServer.tar.gz";
    __AF_download118_v0__11_5="$__AF_download118_v0";
    echo "$__AF_download118_v0__11_5" > /dev/null 2>&1
}
remove_old__127_v0() {
    echo "Removing old"
     rm -r ./app ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "app is not exist, ignore"
fi
     rm -rf ./dist ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "dist is not exist, ignore"
fi
}
tar__129_v0() {
     tar -xvf ./SgridCloudServer.tar.gz ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "dist is not exist, ignore"
fi
     chmod +x ./app ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "add permission failed"
fi
}
# http://124.220.19.199:17853/fm//sgrid/sgridcloud.tar.gz
args=("$0" "$@")
    download_path="${args[1]}"
    if [ $([ "_${download_path}" != "_" ]; echo $?) != 0 ]; then
        echo "error: miss download_path"
        exit__80_v0 1;
        __AF_exit80_v0__49_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__49_9" > /dev/null 2>&1
fi
    download_latest_package__126_v0 "${download_path}";
    __AF_download_latest_package126_v0__51_5="$__AF_download_latest_package126_v0";
    echo "$__AF_download_latest_package126_v0__51_5" > /dev/null 2>&1
    remove_old__127_v0 ;
    __AF_remove_old127_v0__52_5="$__AF_remove_old127_v0";
    echo "$__AF_remove_old127_v0__52_5" > /dev/null 2>&1
    tar__129_v0 ;
    __AF_tar129_v0__53_5="$__AF_tar129_v0";
    echo "$__AF_tar129_v0__53_5" > /dev/null 2>&1
