# 结构化（格式化/缩进）输出

`json.Marshal` 输出的是**紧凑单行** JSON。想要"结构化"（带缩进、可读）输出，Go 提供两种方式。

## 方式一：`json.MarshalIndent`

```go
func MarshalIndent(v any, prefix, indent string) ([]byte, error)
```

- `prefix`：每行前缀（常用 `""`）
- `indent`：每级缩进字符串（`"  "`、`"\t"`、`"    "`）

```go
b, _ := json.MarshalIndent(u, "", "  ")
fmt.Println(string(b))
```

输出：

```json
{
  "id": 1,
  "name": "Alice",
  "tags": [
    "dev",
    "go"
  ],
  "address": {
    "city": "Beijing"
  }
}
```

带前缀的玩法（每行加 `>>`）：

```go
json.MarshalIndent(u, ">>", "  ")
```

```
{
>>  "id": 1,
>>  "name": "Alice"
>>}
```

> ⚠️ 注意：`MarshalIndent` **不会**自动加尾部换行 `\n`，也**不会**改变 HTML 转义行为（默认仍转义 `<>&`）。

---

## 方式二：`json.Encoder`（流式 + 可配置）

适合直接写 HTTP 响应、文件、`stdout`，无需先拿到 `[]byte`：

```go
enc := json.NewEncoder(os.Stdout)
enc.SetIndent("", "  ")        // 开启缩进
enc.SetEscapeHTML(false)       // 关闭 HTML 转义（可选）
err := enc.Encode(u)           // 自动追加 '\n'
```

完整示例：

```go
func writeJSON(w http.ResponseWriter, v any) {
    w.Header().Set("Content-Type", "application/json")
    enc := json.NewEncoder(w)
    enc.SetIndent("", "  ")
    enc.SetEscapeHTML(false)
    if err := enc.Encode(v); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

### Encoder 的关键细节

| 特性 | `MarshalIndent` | `Encoder` |
|------|----------------|-----------|
| 输出 `[]byte` | ✅ | ❌ 直接写入 `io.Writer` |
| 末尾换行 | ❌ 需自己加 | ✅ 自动加 `\n` |
| 关闭 HTML 转义 | ❌ 不支持 | ✅ `SetEscapeHTML(false)` |
| 自定义缩进 | ✅ 参数传入 | ✅ `SetIndent` |
| 适合场景 | 拿到字符串、日志、返回给调用方 | HTTP 响应、文件、管道 |

---

## 两者输出是否一致？

**几乎一致，但有一个换行差异**：

```go
a, _ := json.MarshalIndent(u, "", "  ")
b := new(bytes.Buffer)
enc := json.NewEncoder(b)
enc.SetIndent("", "  ")
enc.Encode(u)

// a: {...}            （无尾换行）
// b: {...}\n          （有尾换行）
// 除尾换行外内容相同
```

若想 `MarshalIndent` 也带换行：

```go
b = append(b, '\n')
```

---

## 常见配合技巧

### 1. 输出到文件（缩进美化日志）

```go
f, _ := os.Create("out.json")
defer f.Close()
enc := json.NewEncoder(f)
enc.SetIndent("", "  ")
enc.SetEscapeHTML(false)
enc.Encode(v)
```

### 2. 调试时打印美化 JSON

```go
func pretty(v any) string {
    b, err := json.MarshalIndent(v, "", "  ")
    if err != nil {
        return fmt.Sprintf("json err: %v", err)
    }
    return string(b)
}
```

### 3. 生产环境别用缩进

缩进会显著增加输出体积（尤其数组），**HTTP API 默认应返回紧凑 JSON**，只有调试/配置/人类阅读场景才缩进。可加开关：

```go
if cfg.Pretty {
    enc.SetIndent("", "  ")
}
```

---

## 一句话总结

> `json.MarshalIndent(v, prefix, indent)` 拿到带缩进的 `[]byte`；`json.NewEncoder(w).SetIndent(prefix, indent).Encode(v)` 流式输出并自动加换行。需要关闭 HTML 转义就只能用 `Encoder`。

需要我对比一下 `MarshalIndent` 的**性能开销**（缩进对 CPU/内存的影响），或者 `SetEscapeHTML(false)` 的安全注意事项吗？