# go-search

基于 BleveSearch 的全文搜索引擎服务，支持通过 JSON API 创建索引、添加文档和执行搜索。

## 项目介绍

go-search 是一个轻量级搜索引擎服务，使用 Go 语言开发，基于 BleveSearch 实现高效的全文搜索功能。
系统提供 RESTful API 接口，支持通过 JSON 格式进行索引管理、文档操作和搜索查询。

## 功能特点

- 支持多索引管理
- JSON 格式的请求和响应
- 分页搜索结果
- 全文检索功能
- 简单易用的 API 接口

## 技术栈

- Go 1.16+
- BleveSearch v2
- Gin Web 框架

## 项目结构

```plainText
go-search/
├── main.go # 程序入口
├── config/ # 配置模块
├── handler/ # HTTP 处理器
├── service/ # 业务逻辑层
├── model/ # 数据模型
├── util/ # 工具函数
└── README.md # 项目文档
```

## 安装与启动

### 前提条件

- Go 1.16 或更高版本
- Git

### 安装步骤

1. 克隆代码库

```bash
git clone https://github.com/yourusername/go-search.git
cd go-search
```

2. 安装依赖

```bash
go mod download
```

3. 启动服务

```bash
go run main.go
```

服务将在 http://localhost:8080 启动

## API 接口文档

### 1. 创建索引

**请求**

- 方法: POST
- 路径: /api/index
- 内容类型: application/json

**请求体**

```json
{
  "index_name": "products"
}
```

**响应**

```json
{
  "message": "索引创建成功"
}
```

### 2. 添加文档

**请求**

- 方法: POST
- 路径: /api/document
- 内容类型: application/json

**请求体**

```json
{
  "index_name": "products",
  "id": "1",
  "fields": {
    "name": "iPhone 13",
    "description": "Apple iPhone 13 128GB 星光色",
    "price": 5999,
    "category": "智能手机"
  }
}
```

**响应**

```json
{
  "message": "文档添加成功"
}
```

### 3. 搜索文档

**请求**

- 方法: POST
- 路径: /api/search
- 内容类型: application/json

**请求体**

```json
{
  "index_name": "products",
  "query": "iPhone",
  "page": 1,
  "size": 10
}
```

**响应**

```json
{
  "total": 2,
  "page": 1,
  "size": 10,
  "hits": [
    {
      "id": "1",
      "score": 0.89,
      "fields": {
        "name": "iPhone 13",
        "description": "Apple iPhone 13 128GB 星光色",
        "price": 5999,
        "category": "智能手机"
      }
    },
    {
      "id": "2",
      "score": 0.75,
      "fields": {
        "name": "iPhone 14",
        "description": "Apple iPhone 14 128GB 午夜色",
        "price": 6999,
        "category": "智能手机"
      }
    }
  ]
}
```

## 错误码说明

- 400: 请求参数错误
- 500: 服务器内部错误
- 错误消息将在响应的 `error` 字段中返回

## 许可证

[MIT](LICENSE)