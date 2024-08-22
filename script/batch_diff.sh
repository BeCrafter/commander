#!/bin/bash

# set -x
# set -e

HOMEDIR=$(dirname $(dirname $(readlink -f $0)))
OUTDIR=${HOMEDIR}/output

cd ${HOMEDIR} && rm -rf ${OUTDIR}  2>/dev/null


index=-1
cmd_params=("$*")
if [[ $# -ge 0 ]]; then
    if [[ $1 =~ ^[0-9]+$ ]]; then
        index=$1
        shift 1
        cmd_params=("$*")
    fi
fi

# echo "Debug: ${cmd_params[@]} \t $index" && exit

# 编译命令
(go build -o "${OUTDIR}/commander" "./main.go") && (echo "Build OK") || (echo "Build Failed"; exit 1)


# 定义请求数据
first_host="http://127.0.0.1:8079"    # 第一个环境 
second_host="http://10.52.3.173:8000" # 第二个环境

url_list=(
    "/rest/info/get?product_id=11"
    "/rest/info/update"
)

method_list=(
    "GET"
    "POST"
)

data_list=(
    ""
    "{\"id\":2010,\"num\":1}"
)
header_list=(
    ""
    "Content-Type:application/json DEBUG_SWITCH:true"
)

run(){
    let i=$1
    params="-t ${first_host} -t ${second_host}"
    url_item=${url_list[$i]}
    if [ -n "${url_item}" ]; then
        params="${params} -u ${url_item}"
    fi
    method_item=${method_list[$i]}
    if [ -n "${method_item}" ]; then
        params="${params} -X ${method_item}"
    fi
    data_item=${data_list[$i]}
    if [ -n "${data_item}" ]; then
        params="${params} -d ${data_item}"
    fi

    header_item=${header_list[$i]}
    for header in $(echo "${header_item}"); do
        if [ -n "${header}" ]; then
            params="${params} -H ${header}"
        fi
    done
    echo "\n\n#======================# Path: ${url_item} #======================#"
    ${OUTDIR}/commander jsondiff ${params} ${cmd_params[@]}
}

if [[ $index -ge 0 ]]; then
    if [[ $index -lt ${#url_list[@]} ]]; then
        run $index
    else
        echo "\n\nIndex: $index is out of range.\n" && exit 1
    fi
else
    for i in $(seq 0 $((${#url_list[@]} -1))); do
        run $i
    done
fi