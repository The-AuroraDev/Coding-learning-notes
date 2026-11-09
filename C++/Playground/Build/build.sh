#!/usr/bin/env bash
#
# build.sh —— 简易 C++ 编译工具
#
# 用法:
#   ./build.sh a.cpp [b.cpp ...]   编译指定源文件（可多个）
#   ./build.sh                     编译当前目录下所有 .cpp
#
# 编译选项: -std=c++14 -O2 -Wl,--stack=1073741824
#           （--stack 仅 MinGW/MSYS/Cygwin 下的链接器支持，其他平台自动省略）
# 输出:     与源文件同目录、同名的可执行文件（Windows 平台追加 .exe）
#
# 可通过环境变量指定编译器:  CXX=clang++ ./build.sh main.cpp

set -uo pipefail

CXX="${CXX:-g++}"
STD_FLAG="-std=c++14"
OPT_FLAG="-O2"
STACK_FLAG="-Wl,--stack=1073741824"

# ---------- 平台判断 ----------
EXT=""
LINK_OPTS=""
case "$(uname -s 2>/dev/null || echo unknown)" in
    MINGW*|MSYS*|CYGWIN*)
        LINK_OPTS="$STACK_FLAG"   # Windows 链接器才认识 --stack
        EXT=".exe"
        ;;
    *)
        ;;                        # Linux/macOS: 不加 --stack，否则 ld 报错
esac

# ---------- 检查编译器 ----------
if ! command -v "$CXX" >/dev/null 2>&1; then
    echo "错误：找不到编译器 '$CXX'。请安装 g++ 或设置 CXX 环境变量。" >&2
    exit 1
fi

# ---------- 单个文件编译 ----------
build_one() {
    local src="$1"

    if [[ ! -f "$src" ]]; then
        echo "[跳过] 文件不存在: $src" >&2
        return 1
    fi

    case "$src" in
        *.cpp|*.cc|*.cxx|*.C|*.c++) ;;
        *) echo "[跳过] 不是 C++ 源文件: $src" >&2; return 1 ;;
    esac

    local base="${src%.*}"          # 去掉扩展名 -> 文件名不变
    local out="${base}${EXT}"

    echo "==> 编译 $src  ->  $out"
    if "$CXX" "$STD_FLAG" "$OPT_FLAG" "$src" -o "$out" $LINK_OPTS; then
        echo "    成功"
        return 0
    else
        echo "    失败" >&2
        return 1
    fi
}

# ---------- 主流程 ----------
targets=()
if [[ $# -gt 0 ]]; then
    targets=("$@")
else
    shopt -s nullglob
    targets=(*.cpp *.cc *.cxx)
    if [[ ${#targets[@]} -eq 0 ]]; then
        echo "错误：当前目录下没有 .cpp 源文件。" >&2
        exit 1
    fi
fi

failed=0
for f in "${targets[@]}"; do
    build_one "$f" || failed=$((failed + 1))
done

if [[ $failed -gt 0 ]]; then
    echo "完成，但有 $failed 个文件编译失败。" >&2
    exit 1
fi
echo "全部编译成功。"
```