package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"bilibili-uploader/internal/converter"
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
	// 创建转换器
	conv := converter.NewWechatConverterFixed()

	// 静态文件服务
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	// 主页
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		serveHomePage(w, r)
	})

	// API接口
	http.HandleFunc("/api/convert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ConvertRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// 转换Markdown
		html := conv.ConvertMarkdownToWechat(req.Markdown)

		response := ConvertResponse{
			HTML:    html,
			Success: true,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on http://localhost:%s\n", port)
	fmt.Printf("Open your browser and visit: http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>微信公众号 Markdown 编辑器</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background-color: #f5f5f5;
            line-height: 1.6;
        }
        
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 2rem 0;
            text-align: center;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        
        .header h1 {
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
        }
        
        .header p {
            font-size: 1.1rem;
            opacity: 0.9;
        }
        
        .container {
            max-width: 1400px;
            margin: 2rem auto;
            padding: 0 1rem;
        }
        
        .editor-container {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 2rem;
            background: white;
            border-radius: 12px;
            box-shadow: 0 4px 20px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        
        .input-section, .output-section {
            padding: 2rem;
        }
        
        .section-title {
            font-size: 1.3rem;
            font-weight: 600;
            margin-bottom: 1rem;
            color: #333;
            border-bottom: 2px solid #667eea;
            padding-bottom: 0.5rem;
        }
        
        #markdown-input {
            width: 100%;
            height: 500px;
            border: 2px solid #e1e5e9;
            border-radius: 8px;
            padding: 1rem;
            font-family: 'SF Mono', Monaco, monospace;
            font-size: 14px;
            line-height: 1.5;
            resize: vertical;
            transition: border-color 0.3s;
        }
        
        #markdown-input:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }
        
        .preview-container {
            border: 2px solid #e1e5e9;
            border-radius: 8px;
            height: 500px;
            overflow-y: auto;
            background: white;
        }
        
        #preview-output {
            padding: 1rem;
        }
        
        .button-group {
            margin-top: 1rem;
            text-align: center;
        }
        
        .btn {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            padding: 0.8rem 2rem;
            border-radius: 6px;
            cursor: pointer;
            font-size: 16px;
            font-weight: 500;
            transition: all 0.3s;
            margin: 0 0.5rem;
        }
        
        .btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
        }
        
        .btn:active {
            transform: translateY(0);
        }
        
        .btn-secondary {
            background: #6c757d;
        }
        
        .loading {
            opacity: 0.6;
            cursor: not-allowed;
        }
        
        .help-section {
            margin-top: 2rem;
            background: white;
            border-radius: 12px;
            padding: 2rem;
            box-shadow: 0 4px 20px rgba(0,0,0,0.1);
        }
        
        .help-section h2 {
            color: #333;
            margin-bottom: 1rem;
        }
        
        .help-section ul {
            list-style-position: inside;
            color: #666;
        }
        
        .help-section li {
            margin-bottom: 0.5rem;
        }
        
        @media (max-width: 768px) {
            .editor-container {
                grid-template-columns: 1fr;
                gap: 1rem;
            }
            
            .header h1 {
                font-size: 2rem;
            }
            
            #markdown-input, .preview-container {
                height: 300px;
            }
        }
        
        .toast {
            position: fixed;
            top: 20px;
            right: 20px;
            background: #28a745;
            color: white;
            padding: 1rem 1.5rem;
            border-radius: 6px;
            box-shadow: 0 4px 12px rgba(0,0,0,0.2);
            z-index: 1000;
            transform: translateX(100%);
            transition: transform 0.3s;
        }
        
        .toast.show {
            transform: translateX(0);
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>微信公众号 Markdown 编辑器</h1>
        <p>简单好用的在线、免费、实时的 Markdown 转微信公众号文章格式工具</p>
    </div>
    
    <div class="container">
        <div class="editor-container">
            <div class="input-section">
                <h2 class="section-title">📝 Markdown 编辑</h2>
                <textarea id="markdown-input" placeholder="在此输入 Markdown 内容...

# 示例标题

这是一个段落示例。

## 二级标题

- 列表项 1
- 列表项 2

> 这是一个引用

**粗体文本** 和 *斜体文本*

` + "`行内代码`" + `

[链接文本](https://example.com)

## 命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| -login | 登录bilibili账号 | false |
| -upload | 指定要上传的视频文件路径 | &quot;&quot; |
| -title | 视频标题 | &quot;&quot; |
| -desc | 视频描述 | &quot;&quot; |
| -tags | 视频标签(逗号分隔) | &quot;&quot; |
| -tid | 视频分区ID | 122 |
| -config | 配置文件路径 | config.json |
| -cookie | Cookie文件路径 | cookies.json |

` + "```javascript" + `
console.log('Hello World');
` + "```" + `
"></textarea>
            </div>
            
            <div class="output-section">
                <h2 class="section-title">👁️ 微信预览</h2>
                <div class="preview-container">
                    <div id="preview-output">
                        <p style="color: #999; text-align: center; margin-top: 2rem;">
                            在左侧输入 Markdown，这里将显示转换后的微信公众号格式
                        </p>
                    </div>
                </div>
                <div class="button-group">
                    <button class="btn" onclick="copyToClipboard()">📋 复制到剪贴板</button>
                    <button class="btn btn-secondary" onclick="clearAll()">🗑️ 清空</button>
                </div>
            </div>
        </div>
        
        <div class="help-section">
            <h2>📖 使用说明</h2>
            <ul>
                <li>在左侧输入框中输入或粘贴 Markdown 文本</li>
                <li>右侧会实时显示转换后的微信公众号格式预览</li>
                <li>点击"复制到剪贴板"按钮复制格式化内容</li>
                <li>在微信公众号编辑器中粘贴即可使用</li>
                <li>链接会自动转换为脚注形式，符合微信公众号规范</li>
                <li>支持标题、列表、引用、代码块、粗体、斜体等格式</li>
            </ul>
        </div>
    </div>
    
    <div id="toast" class="toast"></div>
    
    <script>
        const markdownInput = document.getElementById('markdown-input');
        const previewOutput = document.getElementById('preview-output');
        let convertTimeout;
        
        // 实时转换
        markdownInput.addEventListener('input', function() {
            clearTimeout(convertTimeout);
            convertTimeout = setTimeout(convertMarkdown, 300);
        });
        
        // 初始转换
        convertMarkdown();
        
        async function convertMarkdown() {
            const markdown = markdownInput.value;
            
            if (!markdown.trim()) {
                previewOutput.innerHTML = '<p style="color: #999; text-align: center; margin-top: 2rem;">在左侧输入 Markdown，这里将显示转换后的微信公众号格式</p>';
                return;
            }
            
            try {
                const response = await fetch('/api/convert', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ markdown: markdown })
                });
                
                const result = await response.json();
                
                if (result.success) {
                    previewOutput.innerHTML = result.html;
                } else {
                    previewOutput.innerHTML = '<p style="color: #dc3545;">转换失败: ' + result.error + '</p>';
                }
            } catch (error) {
                previewOutput.innerHTML = '<p style="color: #dc3545;">网络错误: ' + error.message + '</p>';
            }
        }
        
        function copyToClipboard() {
            const html = previewOutput.innerHTML;
            
            if (!html || html.includes('在左侧输入 Markdown')) {
                showToast('请先输入 Markdown 内容', 'error');
                return;
            }
            
            // 创建临时元素
            const tempDiv = document.createElement('div');
            tempDiv.innerHTML = html;
            document.body.appendChild(tempDiv);
            
            // 选择内容
            const range = document.createRange();
            range.selectNode(tempDiv);
            const selection = window.getSelection();
            selection.removeAllRanges();
            selection.addRange(range);
            
            try {
                document.execCommand('copy');
                showToast('已复制到剪贴板！现在可以粘贴到微信公众号编辑器中', 'success');
            } catch (err) {
                showToast('复制失败，请手动选择并复制', 'error');
            }
            
            // 清理
            selection.removeAllRanges();
            document.body.removeChild(tempDiv);
        }
        
        function clearAll() {
            markdownInput.value = '';
            previewOutput.innerHTML = '<p style="color: #999; text-align: center; margin-top: 2rem;">在左侧输入 Markdown，这里将显示转换后的微信公众号格式</p>';
        }
        
        function showToast(message, type = 'success') {
            const toast = document.getElementById('toast');
            toast.textContent = message;
            toast.style.background = type === 'success' ? '#28a745' : '#dc3545';
            toast.classList.add('show');
            
            setTimeout(() => {
                toast.classList.remove('show');
            }, 3000);
        }
    </script>
</body>
</html>`

	t, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ConvertResponse{
		Success: false,
		Error:   message,
	})
}
