# 微信公众号 Markdown 编辑器

一个专为微信公众号设计的 Markdown 编辑器，支持实时预览和一键复制到微信公众号编辑器。

## 特性

- 🚀 **实时预览**: 支持 Markdown 实时渲染预览
- 📱 **微信公众号样式**: 专门优化的微信公众号排版样式
- 🎨 **多种主题**: 内置多种精美主题，包括默认、掘金、知乎等风格
- 💡 **语法高亮**: 支持代码块语法高亮
- 🧮 **数学公式**: 支持 LaTeX 数学公式渲染
- 📋 **一键复制**: 转换后可直接复制到微信公众号编辑器
- 🔧 **自定义样式**: 支持 CSS 样式自定义编辑
- 📖 **Markdown 扩展**: 支持表格、脚注、任务列表等扩展语法

## 快速开始

### 环境要求

- Go 1.23.4 或更高版本
- 现代浏览器（Chrome、Firefox、Safari、Edge）

### 安装和运行

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd markdown-wx
   ```

2. **安装依赖**
   ```bash
   go mod tidy
   ```

3. **启动服务器**
   ```bash
   go run main.go
   ```

4. **打开浏览器**
   
   访问 `http://localhost:8080` 开始使用

### 自定义端口

可以通过环境变量设置自定义端口：

```bash
PORT=3000 go run main.go
```

## 使用方法

1. 在左侧编辑器中输入 Markdown 内容
2. 右侧会实时显示微信公众号样式预览
3. 选择合适的主题样式
4. 点击"复制"按钮，然后粘贴到微信公众号编辑器中

## 支持的 Markdown 语法

### 基础语法
- **标题**: `# ## ###`
- **强调**: `**粗体**` `*斜体*`
- **列表**: 有序列表和无序列表
- **链接**: `[文本](URL)`
- **图片**: `![alt](URL)`
- **代码**: `inline code` 和 ```代码块```

### 扩展语法
- **表格**: 支持表格渲染
- **引用**: `> 引用内容`
- **任务列表**: `- [x] 已完成` `- [ ] 待完成`
- **脚注**: `[^1]` 语法
- **数学公式**: `$inline math$` 和 `$$block math$$`

## 主题样式

内置多种主题样式：

- **默认主题** (gzh_default.css): 经典微信公众号风格
- **掘金主题** (juejin_default.css): 掘金社区风格
- **知乎主题** (zhihu_default.css): 知乎专栏风格
- **Medium主题** (medium_default.css): Medium 平台风格
- **头条主题** (toutiao_default.css): 今日头条风格
- **其他精美主题**: lapis、maize、orangeheart、phycat、pie、purple、rainbow

## 项目结构

```
markdown-wx/
├── main.go                 # 主程序入口
├── go.mod                  # Go 模块定义
├── internal/
│   └── converter/
│       └── markdown_wx.go  # Markdown 转换核心逻辑
└── web/
    └── static/             # 静态资源
        ├── themes/         # 主题样式文件
        ├── codemirror/     # 代码编辑器
        ├── highlight/      # 语法高亮
        ├── marked/         # Markdown 解析器
        ├── mathjax/        # 数学公式渲染
        └── prettier/       # 代码格式化
```

## API 接口

### POST /api/convert

将 Markdown 内容转换为微信公众号 HTML 格式。

**请求体:**
```json
{
  "markdown": "# 标题\n\n内容..."
}
```

**响应:**
```json
{
  "html": "<div>转换后的HTML</div>",
  "success": true,
  "error": ""
}
```

## 技术栈

### 后端
- **Go**: 主要编程语言
- **net/http**: HTTP 服务器
- **html/template**: HTML 模板引擎

### 前端
- **Vanilla JavaScript**: 原生 JavaScript
- **CodeMirror**: 代码编辑器
- **Marked.js**: Markdown 解析器
- **Highlight.js**: 语法高亮
- **MathJax**: 数学公式渲染
- **Prettier**: 代码格式化

## 开发

### 本地开发

```bash
# 启动开发服务器
go run main.go

# 构建项目
go build -o markdown-wx main.go
```

### 添加新主题

1. 在 `web/static/themes/` 目录下创建新的 CSS 文件
2. 按照现有主题的样式结构编写样式
3. 在前端页面中添加主题选择选项

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 更新日志

### v1.0.2
- 优化微信公众号样式适配
- 增加多种主题支持
- 改进代码高亮显示
- 优化数学公式渲染

---

如果这个项目对你有帮助，请给个 ⭐️ Star 支持一下！

