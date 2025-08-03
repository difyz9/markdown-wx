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

### 使用示例

#### 1. 使用 curl 调用接口

```bash
curl -X POST http://localhost:8080/api/convert \
  -H "Content-Type: application/json" \
  -d '{
    "markdown": "# 微信公众号文章\n\n这是一篇**测试文章**，包含以下内容：\n\n## 主要特点\n\n- 支持*斜体*和**粗体**\n- 支持代码：`console.log(\"Hello World\")`\n- 支持列表和链接\n\n## 代码示例\n\n```javascript\nfunction hello(name) {\n  console.log(`Hello, ${name}!`);\n}\n```\n\n> 这是一段引用文字，用于强调重要内容。\n\n访问 [GitHub](https://github.com) 了解更多信息。"
  }'
```

#### 2. 使用 JavaScript fetch API

```javascript
async function convertMarkdown() {
  const markdownContent = `# 微信公众号文章

这是一篇**测试文章**，包含以下内容：

## 主要特点

- 支持*斜体*和**粗体**
- 支持代码：\`console.log("Hello World")\`
- 支持列表和链接

## 代码示例

\`\`\`javascript
function hello(name) {
  console.log(\`Hello, \${name}!\`);
}
\`\`\`

> 这是一段引用文字，用于强调重要内容。

访问 [GitHub](https://github.com) 了解更多信息。`;

  try {
    const response = await fetch('http://localhost:8080/api/convert', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        markdown: markdownContent
      })
    });

    const result = await response.json();
    
    if (result.success) {
      console.log('转换成功！');
      console.log('HTML 内容:', result.html);
      
      // 将结果显示在页面上
      document.getElementById('output').innerHTML = result.html;
    } else {
      console.error('转换失败:', result.error);
    }
  } catch (error) {
    console.error('请求失败:', error);
  }
}

// 调用函数
convertMarkdown();
```

#### 3. 使用 Python requests

```python
import requests
import json

def convert_markdown():
    url = "http://localhost:8080/api/convert"
    
    markdown_content = """# 微信公众号文章

这是一篇**测试文章**，包含以下内容：

## 主要特点

- 支持*斜体*和**粗体**
- 支持代码：`console.log("Hello World")`
- 支持列表和链接

## 代码示例

```python
def hello(name):
    print(f"Hello, {name}!")
```

> 这是一段引用文字，用于强调重要内容。

访问 [GitHub](https://github.com) 了解更多信息。"""

    payload = {
        "markdown": markdown_content
    }
    
    headers = {
        "Content-Type": "application/json"
    }
    
    try:
        response = requests.post(url, data=json.dumps(payload), headers=headers)
        response.raise_for_status()
        
        result = response.json()
        
        if result['success']:
            print("转换成功！")
            print("HTML 内容:")
            print(result['html'])
            
            # 可以将结果保存到文件
            with open('output.html', 'w', encoding='utf-8') as f:
                f.write(result['html'])
            print("结果已保存到 output.html")
        else:
            print(f"转换失败: {result['error']}")
            
    except requests.exceptions.RequestException as e:
        print(f"请求失败: {e}")

# 调用函数
if __name__ == "__main__":
    convert_markdown()
```

#### 4. 使用 Go 调用接口

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type ConvertRequest struct {
    Markdown string `json:"markdown"`
}

type ConvertResponse struct {
    HTML    string `json:"html"`
    Success bool   `json:"success"`
    Error   string `json:"error,omitempty"`
}

func main() {
    markdownContent := `# 微信公众号文章

这是一篇**测试文章**，包含以下内容：

## 主要特点

- 支持*斜体*和**粗体**
- 支持代码：` + "`console.log(\"Hello World\")`" + `
- 支持列表和链接

## 代码示例

` + "```go" + `
func hello(name string) {
    fmt.Printf("Hello, %s!\n", name)
}
` + "```" + `

> 这是一段引用文字，用于强调重要内容。

访问 [GitHub](https://github.com) 了解更多信息。`

    // 创建请求体
    reqBody := ConvertRequest{
        Markdown: markdownContent,
    }

    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        fmt.Printf("序列化请求失败: %v\n", err)
        return
    }

    // 发送 POST 请求
    resp, err := http.Post("http://localhost:8080/api/convert", 
        "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Printf("请求失败: %v\n", err)
        return
    }
    defer resp.Body.Close()

    // 读取响应
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("读取响应失败: %v\n", err)
        return
    }

    // 解析响应
    var result ConvertResponse
    if err := json.Unmarshal(body, &result); err != nil {
        fmt.Printf("解析响应失败: %v\n", err)
        return
    }

    if result.Success {
        fmt.Println("转换成功！")
        fmt.Println("HTML 内容:")
        fmt.Println(result.HTML)
    } else {
        fmt.Printf("转换失败: %s\n", result.Error)
    }
}
```

#### 5. 完整的 HTML 测试页面

创建一个 `test.html` 文件来测试接口：

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API 测试页面</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 1200px; margin: 0 auto; padding: 20px; }
        .container { display: flex; gap: 20px; }
        .input-section, .output-section { flex: 1; }
        textarea { width: 100%; height: 300px; font-family: monospace; }
        button { padding: 10px 20px; background: #007cba; color: white; border: none; cursor: pointer; }
        button:hover { background: #005a87; }
        .output { border: 1px solid #ddd; padding: 15px; min-height: 300px; background: #f9f9f9; }
    </style>
</head>
<body>
    <h1>微信公众号 Markdown 转换 API 测试</h1>
    
    <div class="container">
        <div class="input-section">
            <h3>输入 Markdown</h3>
            <textarea id="markdown-input" placeholder="在此输入 Markdown 内容..."># 微信公众号文章

这是一篇**测试文章**，包含以下内容：

## 主要特点

- 支持*斜体*和**粗体**
- 支持代码：`console.log("Hello World")`
- 支持列表和链接

## 代码示例

```javascript
function hello(name) {
  console.log(`Hello, ${name}!`);
}
```

> 这是一段引用文字，用于强调重要内容。

访问 [GitHub](https://github.com) 了解更多信息。</textarea>
            <br><br>
            <button onclick="convertMarkdown()">转换为微信公众号格式</button>
        </div>
        
        <div class="output-section">
            <h3>转换结果</h3>
            <div id="output" class="output">点击转换按钮查看结果...</div>
        </div>
    </div>

    <script>
        async function convertMarkdown() {
            const markdownContent = document.getElementById('markdown-input').value;
            const outputDiv = document.getElementById('output');
            
            if (!markdownContent.trim()) {
                outputDiv.innerHTML = '<p style="color: red;">请输入 Markdown 内容</p>';
                return;
            }
            
            outputDiv.innerHTML = '<p>转换中...</p>';
            
            try {
                const response = await fetch('http://localhost:8080/api/convert', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        markdown: markdownContent
                    })
                });

                const result = await response.json();
                
                if (result.success) {
                    outputDiv.innerHTML = result.html;
                } else {
                    outputDiv.innerHTML = `<p style="color: red;">转换失败: ${result.error}</p>`;
                }
            } catch (error) {
                outputDiv.innerHTML = `<p style="color: red;">请求失败: ${error.message}</p>`;
            }
        }
    </script>
</body>
</html>
```
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

