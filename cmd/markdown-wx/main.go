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
    
    <!-- 引入 wenyan-3.0 的 JavaScript 库 -->
    <script src="static/marked/marked.min.js"></script>
    <script src="static/highlight/highlight.min.js"></script>
    <script src="static/marked/marked_hljs.umd.min.js"></script>
    <script src="static/csstree/csstree.js"></script>
    
    <script type="module">
        import frontMatter from './static/marked/front-matter+esm.js';
        window.frontMatter = frontMatter;
    </script>
    
    <script>
        function addContainer(math, doc) {
            const tag = math.display ? 'section' : 'span';
            const cls = math.display ? 'block-equation' : 'inline-equation';
            math.typesetRoot.setAttribute("math", math.math);
            math.typesetRoot = doc.adaptor.node(tag, {class: cls}, [math.typesetRoot]);
        }
        MathJax = {
            options: {
                renderActions: {
                    addContainer: [190, (doc) => {for (const math of doc.math) {addContainer(math, doc)}}, addContainer]
                }
            },
            svg: {
                fontCache: 'none'
            },
            tex: {
                inlineMath: [['$', '$'], ['\\(', '\\)']],
                displayMath: [['$$', '$$'], ['\\[', '\\]']],
                processEscapes: true
            }
        };
    </script>
    <script src="static/mathjax/tex-svg-full.min.js"></script>
    
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
        
        .controls {
            margin-bottom: 1rem;
        }
        
        .theme-selector {
            padding: 0.5rem;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 14px;
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
        
        /* 隐藏页面底部可能出现的wenyan元素 */
        #wenyan {
            display: none !important;
        }
        
        /* 确保预览区域的样式正确 */
        #preview-output > * {
            max-width: 100%;
        }
        
        /* 防止预览内容溢出 */
        .preview-container {
            position: relative;
            overflow: hidden;
        }
        
        #preview-output {
            overflow-y: auto;
            max-height: 468px; /* 500px - 2rem padding */
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>微信公众号 Markdown 编辑器</h1>
        <p>基于 WenYan 引擎的高级 Markdown 转微信公众号格式工具</p>
    </div>
    
    <div class="container">
        <div class="editor-container">
            <div class="input-section">
                <h2 class="section-title">📝 Markdown 编辑</h2>
                <div class="controls">
                    <label for="theme-select">主题样式：</label>
                    <select id="theme-select" class="theme-selector" onchange="changeTheme()">
                        <option value="gzh_default">公众号默认</option>
                        <option value="juejin_default">掘金</option>
                        <option value="zhihu_default">知乎</option>
                        <option value="lapis">青金石</option>
                        <option value="maize">玉米黄</option>
                        <option value="orangeheart">橙心</option>
                        <option value="phycat">物理猫</option>
                        <option value="pie">馅饼</option>
                        <option value="purple">紫色</option>
                        <option value="rainbow">彩虹</option>
                    </select>
                </div>
                <textarea id="markdown-input" placeholder="在此输入 Markdown 内容...

# 示例标题

这是一个段落示例，支持**粗体**和*斜体*文本。

## 二级标题

- 列表项 1
- 列表项 2
- 列表项 3

> 这是一个引用块，用于强调重要内容。

### 代码示例

行内代码：` + "`const hello = 'world'`" + `

代码块：
` + "```javascript" + `
function greet(name) {
    console.log('Hello, ' + name + '!');
}
greet('微信公众号');
` + "```" + `

### 数学公式

行内公式：$E = mc^2$

块级公式：
$$\\int_{-\\infty}^{\\infty} e^{-x^2} dx = \\sqrt{\\pi}$$

### 表格

| 特性 | 支持情况 | 说明 |
|------|---------|------|
| 标题 | ✅ | 支持 H1-H6 |
| 列表 | ✅ | 有序/无序列表 |
| 代码 | ✅ | 语法高亮 |
| 公式 | ✅ | LaTeX 数学公式 |

### 链接

这是一个[链接示例](https://example.com)，会自动转换为脚注格式。
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
                    <button class="btn" onclick="copyToClipboard()">📋 快速复制</button>
                    <button class="btn" onclick="previewWechat()">� 预览并复制</button>
                    <button class="btn btn-secondary" onclick="clearAll()">🗑️ 清空</button>
                </div>
            </div>
        </div>
        
        <div class="help-section">
            <h2>📖 功能特性</h2>
            <ul>
                <li><strong>完整 Markdown 支持</strong>：标题、段落、列表、引用、代码块等</li>
                <li><strong>数学公式渲染</strong>：支持 LaTeX 行内和块级数学公式</li>
                <li><strong>代码语法高亮</strong>：支持多种编程语言的语法高亮</li>
                <li><strong>多主题样式</strong>：内置多种精美主题，适配不同风格需求</li>
                <li><strong>脚注自动转换</strong>：链接自动转换为微信公众号规范的脚注格式</li>
                <li><strong>无背景色复制</strong>：生成的内容完全兼容微信公众号编辑器</li>
                <li><strong>实时预览</strong>：输入即时预览，所见即所得</li>
            </ul>
        </div>
    </div>
    
    <div id="toast" class="toast"></div>
    
    <!-- WenYan 转换引擎脚本 - 确保在所有依赖库之后加载 -->
    <script src="static/main.js"></script>
    
    <script>
        const markdownInput = document.getElementById('markdown-input');
        const previewOutput = document.getElementById('preview-output');
        let convertTimeout;
        let currentTheme = 'gzh_default';
        
        // 初始化 - 等待所有脚本加载完成
        function initializeApp() {
            console.log('Initializing app...');
            console.log('setContent function available:', typeof setContent);
            console.log('getContentForGzh function available:', typeof getContentForGzh);
            console.log('marked available:', typeof marked);
            console.log('hljs available:', typeof hljs);
            console.log('MathJax available:', typeof MathJax);
            
            loadTheme(currentTheme);
            convertMarkdown();
        }
        
        // 等待所有依赖加载完成
        document.addEventListener('DOMContentLoaded', function() {
            console.log('DOM loaded, checking dependencies...');
            
            // 检查所有必需的依赖
            const dependencies = {
                'marked': typeof marked,
                'hljs': typeof hljs,
                'MathJax': typeof MathJax,
                'frontMatter': typeof window.frontMatter,
                'csstree': typeof csstree
            };
            
            console.log('Dependencies status:', dependencies);
            
            // 检查关键函数是否加载
            if (typeof setContent === 'function' && typeof getContentForGzh === 'function') {
                console.log('转换引擎加载成功！');
                initializeApp();
            } else {
                console.log('转换引擎未加载，等待500ms后重试...');
                // 如果函数还没加载，等待一下再试
                setTimeout(() => {
                    if (typeof setContent === 'function' && typeof getContentForGzh === 'function') {
                        console.log('转换引擎延迟加载成功！');
                        initializeApp();
                    } else {
                        console.error('转换引擎加载失败！');
                        console.log('setContent:', typeof setContent);
                        console.log('getContentForGzh:', typeof getContentForGzh);
                        console.log('Available window functions:', Object.getOwnPropertyNames(window).filter(name => 
                            name.includes('Content') || name.includes('set') || name.includes('get')
                        ));
                        
                        // 显示错误信息给用户
                        previewOutput.innerHTML = '<p style="color: #dc3545; text-align: center; margin-top: 2rem;">⚠️ 转换引擎加载失败<br>请刷新页面重试</p>';
                    }
                }, 1000);
            }
        });
        
        // 实时转换
        markdownInput.addEventListener('input', function() {
            clearTimeout(convertTimeout);
            convertTimeout = setTimeout(convertMarkdown, 300);
        });
        
        // 加载主题
        async function loadTheme(themeName) {
            try {
                const response = await fetch('static/themes/' + themeName + '.css');
                const css = await response.text();
                setCustomTheme(css);
                convertMarkdown();
            } catch (error) {
                console.error('Failed to load theme:', error);
            }
        }
        
        // 切换主题
        function changeTheme() {
            const select = document.getElementById('theme-select');
            currentTheme = select.value;
            loadTheme(currentTheme);
        }
        
        // 转换 Markdown
        function convertMarkdown() {
            const markdown = markdownInput.value;
            
            if (!markdown.trim()) {
                previewOutput.innerHTML = '<p style="color: #999; text-align: center; margin-top: 2rem;">在左侧输入 Markdown，这里将显示转换后的微信公众号格式</p>';
                return;
            }
            
            try {
                // 使用 WenYan 引擎转换内容
                if (typeof setContent === 'function') {
                    setContent(markdown);
                    
                    // 等待 MathJax 渲染完成后显示内容
                    setTimeout(() => {
                        const wenyanElement = document.getElementById('wenyan');
                        if (wenyanElement) {
                            // 克隆内容到预览区域，而不是移动原始元素
                            previewOutput.innerHTML = wenyanElement.innerHTML;
                        } else {
                            // 如果 wenyan 元素不存在，直接获取转换后的HTML
                            const convertedHtml = getContentForGzh();
                            if (convertedHtml) {
                                previewOutput.innerHTML = convertedHtml;
                            } else {
                                previewOutput.innerHTML = '<p style="color: #dc3545;">转换失败: 无法获取转换后的内容</p>';
                            }
                        }
                    }, 200);
                } else {
                    previewOutput.innerHTML = '<p style="color: #dc3545;">转换引擎未加载</p>';
                }
            } catch (error) {
                console.error('Conversion error:', error);
                previewOutput.innerHTML = '<p style="color: #dc3545;">转换失败: ' + error.message + '</p>';
            }
        }
        
        // 复制到剪贴板
        function copyToClipboard() {
            try {
                // 检查是否有内容可以复制
                const markdown = markdownInput.value;
                if (!markdown.trim()) {
                    showToast('请先输入 Markdown 内容', 'error');
                    return;
                }
                
                // 使用专为微信公众号优化的 getContentForGzh 函数
                if (typeof getContentForGzh !== 'function') {
                    showToast('转换引擎未加载，请刷新页面重试', 'error');
                    return;
                }
                
                // 确保内容已经转换
                if (typeof setContent === 'function') {
                    setContent(markdown);
                }
                
                // 等待转换完成后获取微信公众号格式的HTML
                setTimeout(() => {
                    try {
                        const wechatHtml = getContentForGzh();
                        
                        if (!wechatHtml || wechatHtml.length < 50) {
                            showToast('获取转换内容失败，请重试', 'error');
                            return;
                        }
                        
                        // 创建临时元素用于复制
                        const tempDiv = document.createElement('div');
                        tempDiv.innerHTML = wechatHtml;
                        
                        // 添加到DOM中（但不显示）
                        tempDiv.style.position = 'absolute';
                        tempDiv.style.left = '-9999px';
                        document.body.appendChild(tempDiv);
                        
                        // 选择内容
                        const range = document.createRange();
                        range.selectNode(tempDiv);
                        const selection = window.getSelection();
                        selection.removeAllRanges();
                        selection.addRange(range);
                        
                        // 执行复制
                        const successful = document.execCommand('copy');
                        
                        // 清理DOM
                        selection.removeAllRanges();
                        document.body.removeChild(tempDiv);
                        
                        if (successful) {
                            showToast('✅ 已复制微信公众号格式！可直接粘贴到编辑器', 'success');
                        } else {
                            throw new Error('复制命令执行失败');
                        }
                        
                    } catch (innerErr) {
                        console.error('内部复制错误:', innerErr);
                        showToast('复制失败，请手动选择预览区域内容复制', 'error');
                    }
                }, 300);
                
            } catch (err) {
                console.error('Copy failed:', err);
                showToast('复制失败，请手动选择预览区域内容复制', 'error');
            }
        }
        
        // 预览微信格式（调试用）
        function previewWechat() {
            const markdown = markdownInput.value;
            if (!markdown.trim()) {
                showToast('请先输入 Markdown 内容', 'error');
                return;
            }
            
            if (typeof setContent === 'function' && typeof getContentForGzh === 'function') {
                setContent(markdown);
                setTimeout(() => {
                    const wechatHtml = getContentForGzh();
                    console.log('微信公众号格式HTML预览：');
                    console.log(wechatHtml);
                    
                    // 创建预览窗口
                    const newWindow = window.open('', '_blank', 'width=400,height=700,scrollbars=yes,resizable=yes');
                    
                    // 构建HTML内容
                    const doc = newWindow.document;
                    doc.write('<!DOCTYPE html>');
                    doc.write('<html lang="zh-CN">');
                    doc.write('<head>');
                    doc.write('<meta charset="UTF-8">');
                    doc.write('<title>微信公众号格式预览</title>');
                    doc.write('<style>');
                    doc.write('body{margin:0;padding:0;font-family:-apple-system,BlinkMacSystemFont,sans-serif;background:#f0f0f0}');
                    doc.write('.header{background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:white;padding:15px;text-align:center;box-shadow:0 2px 10px rgba(0,0,0,0.1)}');
                    doc.write('.header h2{margin:0;font-size:18px}');
                    doc.write('.copy-section{background:white;padding:15px;border-bottom:1px solid #e0e0e0;text-align:center}');
                    doc.write('.copy-btn{background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:white;border:none;padding:10px 20px;border-radius:6px;cursor:pointer;font-size:14px;font-weight:500;margin:0 5px;transition:all 0.3s}');
                    doc.write('.copy-btn:hover{transform:translateY(-1px);box-shadow:0 4px 12px rgba(102,126,234,0.3)}');
                    doc.write('.copy-btn.secondary{background:#6c757d}');
                    doc.write('.preview-container{background:white;margin:20px auto;padding:20px;border-radius:8px;box-shadow:0 2px 10px rgba(0,0,0,0.1);max-width:375px;border:1px solid #e0e0e0}');
                    doc.write('.status{padding:10px 15px;text-align:center;font-size:14px;background:#d4edda;color:#155724;border:1px solid #c3e6cb;margin:10px 20px;border-radius:4px;display:none}');
                    doc.write('.status.error{background:#f8d7da;color:#721c24;border-color:#f5c6cb}');
                    doc.write('.status.show{display:block}');
                    doc.write('.footer{text-align:center;padding:15px;color:#666;font-size:12px;background:white;border-top:1px solid #e0e0e0}');
                    doc.write('#preview-content{user-select:text;-webkit-user-select:text;-moz-user-select:text;-ms-user-select:text}');
                    doc.write('</style>');
                    doc.write('</head>');
                    doc.write('<body>');
                    doc.write('<div class="header"><h2>📱 微信公众号格式预览</h2></div>');
                    doc.write('<div class="copy-section">');
                    doc.write('<button class="copy-btn" onclick="copyContent()">📋 复制内容</button>');
                    doc.write('<button class="copy-btn secondary" onclick="selectAll()">🎯 全选</button>');
                    doc.write('<button class="copy-btn secondary" onclick="copyHTML()">💻 复制HTML</button>');
                    doc.write('</div>');
                    doc.write('<div id="status" class="status"></div>');
                    doc.write('<div class="preview-container">');
                    doc.write('<div id="preview-content">' + wechatHtml + '</div>');
                    doc.write('</div>');
                    doc.write('<div class="footer">💡 可以直接选择上方内容复制，或点击"复制内容"按钮<br>复制后可直接粘贴到微信公众号编辑器</div>');
                    
                    // 添加JavaScript
                    doc.write('<script>');
                    doc.write('function showStatus(message,isError){');
                    doc.write('const status=document.getElementById("status");');
                    doc.write('status.textContent=message;');
                    doc.write('status.className="status show"+(isError?" error":"");');
                    doc.write('setTimeout(()=>status.classList.remove("show"),3000);');
                    doc.write('}');
                    doc.write('function copyContent(){');
                    doc.write('try{');
                    doc.write('const content=document.getElementById("preview-content");');
                    doc.write('const range=document.createRange();');
                    doc.write('range.selectNode(content);');
                    doc.write('const selection=window.getSelection();');
                    doc.write('selection.removeAllRanges();');
                    doc.write('selection.addRange(range);');
                    doc.write('const successful=document.execCommand("copy");');
                    doc.write('if(successful){showStatus("✅ 内容已复制到剪贴板！可以粘贴到微信公众号编辑器了");}');
                    doc.write('else{throw new Error("复制命令失败");}');
                    doc.write('selection.removeAllRanges();');
                    doc.write('}catch(err){console.error("复制失败:",err);showStatus("❌ 复制失败，请手动选择内容复制",true);}');
                    doc.write('}');
                    doc.write('function selectAll(){');
                    doc.write('const content=document.getElementById("preview-content");');
                    doc.write('const range=document.createRange();');
                    doc.write('range.selectNode(content);');
                    doc.write('const selection=window.getSelection();');
                    doc.write('selection.removeAllRanges();');
                    doc.write('selection.addRange(range);');
                    doc.write('showStatus("📝 内容已全选，可以Ctrl+C复制");');
                    doc.write('}');
                    doc.write('function copyHTML(){');
                    doc.write('try{');
                    doc.write('const content=document.getElementById("preview-content").outerHTML;');
                    doc.write('const textArea=document.createElement("textarea");');
                    doc.write('textArea.value=content;');
                    doc.write('document.body.appendChild(textArea);');
                    doc.write('textArea.select();');
                    doc.write('const successful=document.execCommand("copy");');
                    doc.write('document.body.removeChild(textArea);');
                    doc.write('if(successful){showStatus("💻 HTML源码已复制到剪贴板");}');
                    doc.write('else{throw new Error("复制HTML失败");}');
                    doc.write('}catch(err){console.error("复制HTML失败:",err);showStatus("❌ 复制HTML失败",true);}');
                    doc.write('}');
                    doc.write('window.onload=function(){showStatus("🎉 预览页面加载完成！点击按钮或手动选择内容进行复制");};');
                    doc.write('</script>');
                    doc.write('</body></html>');
                    doc.close();
                    
                    newWindow.focus();
                }, 300);
            } else {
                showToast('转换引擎未加载', 'error');
            }
        }
        
        // 清空内容
        function clearAll() {
            markdownInput.value = '';
            previewOutput.innerHTML = '<p style="color: #999; text-align: center; margin-top: 2rem;">在左侧输入 Markdown，这里将显示转换后的微信公众号格式</p>';
        }
        
        // 显示提示
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
