# GoUtils

[![Go Reference](https://pkg.go.dev/badge/github.com/jelech/goutils.svg)](https://pkg.go.dev/github.com/jelech/goutils)
[![Go Report Card](https://goreportcard.com/badge/github.com/jelech/goutils)](https://goreportcard.com/report/github.com/jelech/goutils)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

GoUtils 是一个轻量级的 Go 工具库，提供了常用的工具函数和组件，帮助开发者提高开发效率。

## Features

### Core Packages

#### Retry
- Exponential backoff retry mechanism
- Context-aware retries with cancellation support
- Customizable retry attempts and delays
- Support for retry conditions

#### HTTP Client
- Built-in retry capabilities for HTTP requests
- Configurable timeouts and headers
- Support for context cancellation
- Automatic retry on transient failures

#### Cache
- Memory cache with TTL support
- LRU (Least Recently Used) cache implementation
- Thread-safe operations
- Customizable capacity and expiration

#### String Utilities
- Case conversions (camelCase, PascalCase, kebab-case)
- String reversal and manipulation
- Random string generation
- Email validation
- String padding utilities

#### Convert
- Type conversions with error handling
- Struct to map conversion
- JSON marshaling/unmarshaling helpers
- Time parsing utilities
- Slice type conversions

#### Parquet
- **Go 1.17 Compatible**: Placeholder implementation for future development
- **Current Status**: API structure defined, actual I/O functionality requires Go 1.18+ for generics
- Writer with buffering capability (placeholder)
- Error messages indicating not-yet-implemented features
- Future: High-performance Parquet file reading and writing with type-safe operations

## 安装

```bash
go get github.com/jelech/goutils
```

## 快速开始

### Retry 重试机制

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/jelech/goutils/retry"
)

func main() {
    // 基本重试
    err := retry.Do(func() error {
        // 你的业务逻辑
        return doSomething()
    })
    
    // 自定义重试策略
    err = retry.Do(
        func() error {
            return doSomething()
        },
        retry.WithMaxAttempts(5),
        retry.WithBackoff(retry.ExponentialBackoff),
        retry.WithDelay(time.Second),
    )
}
```

### HTTP 客户端

```go
package main

import (
    "github.com/jelech/goutils/http"
)

func main() {
    client := http.NewClient()
    
    resp, err := client.Get("https://api.example.com/data")
    if err != nil {
        // 处理错误
    }
    
    // 带重试的请求
    resp, err = client.GetWithRetry("https://api.example.com/data", 3)
}
```

### 缓存

```go
package main

import (
    "time"
    
    "github.com/jelech/goutils/cache"
)

func main() {
    // 内存缓存
    c := cache.NewMemoryCache()
    c.Set("key", "value", time.Minute*10)
    
    value, ok := c.Get("key")
    if ok {
        fmt.Println(value)
    }
}
```

### Parquet 文件操作

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/jelech/goutils/parquet"
)

type Person struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`
    Age      int     `json:"age"`
    Salary   float64 `json:"salary"`
    IsActive bool    `json:"is_active"`
}

func main() {
    // 注意：当前实现是 Go 1.17 兼容的占位符版本
    // 完整的 Parquet 功能需要 Go 1.18+ 以支持泛型
    
    // Writer API（缓冲功能正常，实际文件写入为占位符）
    writer, err := parquet.NewWriter("data.parquet")
    if err != nil {
        log.Fatal(err)
    }
    
    people := []Person{
        {ID: 1, Name: "Alice", Age: 30, Salary: 50000.0, IsActive: true},
        {ID: 2, Name: "Bob", Age: 25, Salary: 45000.0, IsActive: true},
    }
    
    // 缓冲数据（正常工作）
    err = writer.Write(people)
    if err != nil {
        log.Fatal(err)
    }
    
    // 关闭 writer（占位符）
    err = writer.Close()
    if err != nil {
        log.Printf("占位符实现: %v", err)
    }
    
    // 直接文件操作（Go 1.17 版本为占位符）
    err = parquet.WriteFile("data.parquet", people)
    if err != nil {
        log.Printf("占位符实现: %v", err) // 当前版本中的预期行为
    }
    
    var readData []Person
    err = parquet.ReadFile("data.parquet", &readData)
    if err != nil {
        log.Printf("占位符实现: %v", err) // 当前版本中的预期行为
    }
    
    fmt.Println("Parquet 包已为未来实现做好准备")
}
```

## API 文档

详细的 API 文档请查看 [pkg.go.dev](https://pkg.go.dev/github.com/jelech/goutils)

## 贡献

欢迎提交 Pull Request 和 Issue！

## 许可证

本项目采用 MIT 许可证。详情请查看 [LICENSE](LICENSE) 文件。
