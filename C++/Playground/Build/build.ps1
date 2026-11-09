<#
.SYNOPSIS
    build.ps1 —— 简易 C++ 编译工具

.DESCRIPTION
    编译选项: -std=c++14 -O2 -Wl,--stack=1073741824
    输出:     与源文件同目录、同名的 .exe 文件

.PARAMETER Files
    要编译的 C++ 源文件；不传则编译当前目录下所有 .cpp

.EXAMPLE
    .\build.ps1 main.cpp

.EXAMPLE
    .\build.ps1

.NOTES
    如提示“禁止运行脚本”，可执行:
    powershell -ExecutionPolicy Bypass -File .\build.ps1 main.cpp
    也可通过环境变量指定编译器:  $env:CXX = "C:\mingw64\bin\g++.exe"
#>
[CmdletBinding()]
param(
    [Parameter(Position = 0, ValueFromRemainingArguments = $true)]
    [string[]]$Files
)

# 不要用 'Stop'：g++ 输出到 stderr 时可能被当成异常中断
$ErrorActionPreference = 'Continue'

$Cxx   = if ($env:CXX) { $env:CXX } else { 'g++' }
$Flags = @('-std=c++14', '-O2', '-Wl,--stack=1073741824')

# ---------- 检查编译器 ----------
if (-not (Get-Command $Cxx -ErrorAction SilentlyContinue)) {
    Write-Host "错误：找不到编译器 '$Cxx'。请安装 MinGW-w64 或设置 CXX 环境变量。" -ForegroundColor Red
    exit 1
}

# ---------- 收集源文件 ----------
if (-not $Files -or $Files.Count -eq 0) {
    $Files = @(Get-ChildItem -File -Filter *.cpp |
               Select-Object -ExpandProperty Name)
    if ($Files.Count -eq 0) {
        Write-Host "错误：当前目录下没有 .cpp 源文件。" -ForegroundColor Red
        exit 1
    }
}

# ---------- 逐个编译 ----------
$failed = 0
foreach ($src in $Files) {

    if (-not (Test-Path -LiteralPath $src -PathType Leaf)) {
        Write-Host "[跳过] 文件不存在: $src" -ForegroundColor Yellow
        $failed++
        continue
    }

    $ext = [System.IO.Path]::GetExtension($src).ToLowerInvariant()
    if ($ext -notin @('.cpp', '.cc', '.cxx', '.c++')) {
        Write-Host "[跳过] 不是 C++ 源文件: $src" -ForegroundColor Yellow
        $failed++
        continue
    }

    # 文件名不变，仅把扩展名换成 .exe
    $out = [System.IO.Path]::ChangeExtension($src, '.exe')

    Write-Host "==> 编译 $src  ->  $out"
    & $Cxx @Flags $src -o $out

    if ($LASTEXITCODE -ne 0) {
        Write-Host "    失败" -ForegroundColor Red
        $failed++
    }
    else {
        Write-Host "    成功" -ForegroundColor Green
    }
}

# ---------- 汇总 ----------
if ($failed -gt 0) {
    Write-Host "完成，但有 $failed 个文件编译失败。" -ForegroundColor Red
    exit 1
}
Write-Host "全部编译成功。" -ForegroundColor Green