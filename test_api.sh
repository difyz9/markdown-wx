#!/bin/bash

# 微信公众号 Markdown 转换 API 测试脚本
# 使用方法: ./test_api.sh

API_URL="http://localhost:8080/api/convert"
OUTPUT_FILE="curl_output.html"

echo "🚀 微信公众号 Markdown 转换 API 测试"
echo "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "="

# 检查服务器是否运行
echo "🔍 检查服务器状态..."
if curl -s --connect-timeout 3 "$API_URL" > /dev/null 2>&1; then
    echo "✅ 服务器运行正常"
else
    echo "❌ 无法连接到服务器 ($API_URL)"
    echo "💡 请先启动服务器: go run main.go"
    exit 1
fi

echo ""
echo "📝 准备测试内容..."

# Markdown 测试内容
read -r -d '' MARKDOWN_CONTENT << 'EOF'
# 🎯 微信公众号技术分享

大家好！今天和大家分享一个实用的技术工具。

## 🚀 项目介绍

这是一个**微信公众号 Markdown 编辑器**，主要特点包括：

- ✨ **实时预览**: 所见即所得的编辑体验
- 🎨 **多种主题**: 内置多款精美主题样式  
- 💻 **语法高亮**: 完美支持代码块展示
- 📱 **一键复制**: 直接粘贴到公众号编辑器

## 💡 使用场景

### 技术博主
对于经常写技术文章的朋友，这个工具可以：
1. 提高写作效率
2. 保证排版美观
3. 支持代码高亮

### 产品经理  
撰写产品文档时，可以：
- 快速制作**需求文档**
- 生成*项目报告*  
- 创建`流程说明`

## 🔧 技术实现

### 后端架构
```go
// Go 语言实现的转换服务
func ConvertMarkdownToWechat(markdown string) string {
    // 解析 Markdown 内容
    parsed := parseMarkdown(markdown)
    
    // 应用微信公众号样式
    styled := applyWechatStyles(parsed)
    
    return styled
}
```

### 前端技术栈
```javascript
// 实时预览功能
const editor = new CodeMirror(container, {
    mode: 'markdown',
    theme: 'default',
    lineNumbers: true
});

editor.on('change', (cm) => {
    const content = cm.getValue();
    updatePreview(content);
});
```

## 📊 性能数据

| 指标 | 数值 | 说明 |
|------|------|------|
| 转换速度 | < 100ms | 毫秒级响应 |
| 支持主题 | 10+ | 多样化选择 |
| 浏览器兼容 | 95%+ | 主流浏览器 |
| 文件大小 | < 2MB | 轻量级工具 |

## 🧮 数学公式支持

### 基础公式
质能方程: $E = mc^2$

### 复杂公式  
$$
\int_{-\infty}^{\infty} e^{-x^2} dx = \sqrt{\pi}
$$

## 📝 最佳实践

> **重要提示**: 在使用过程中，建议遵循以下规范：
> 
> 1. 保持文章结构清晰
> 2. 合理使用标题层级
> 3. 适当添加表情符号增加趣味性

## 🎉 总结

通过这个工具，我们可以：

- [x] 提升内容创作效率
- [x] 获得专业的排版效果
- [x] 节省大量格式调整时间
- [ ] 持续优化用户体验

---

**如果觉得有用，欢迎分享给更多朋友！** 🌟

*— 技术分享，让创作更简单*
EOF

echo "📏 内容长度: $(echo "$MARKDOWN_CONTENT" | wc -c) 字符"
echo ""

# 创建 JSON 请求体
JSON_PAYLOAD=$(cat << EOF
{
  "markdown": $(echo "$MARKDOWN_CONTENT" | jq -Rs .)
}
EOF
)

echo "🔄 正在调用转换 API..."

# 发送请求并获取响应
RESPONSE=$(curl -s -w "\n%{http_code}" \
  -X POST \
  -H "Content-Type: application/json" \
  -d "$JSON_PAYLOAD" \
  "$API_URL")

# 分离响应体和状态码
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$RESPONSE" | head -n -1)

echo "📡 HTTP 状态码: $HTTP_CODE"

if [ "$HTTP_CODE" -eq 200 ]; then
    echo "✅ 请求成功！"
    
    # 解析 JSON 响应
    SUCCESS=$(echo "$RESPONSE_BODY" | jq -r '.success')
    
    if [ "$SUCCESS" = "true" ]; then
        echo "🎉 转换成功！"
        
        # 提取 HTML 内容
        HTML_CONTENT=$(echo "$RESPONSE_BODY" | jq -r '.html')
        HTML_LENGTH=$(echo "$HTML_CONTENT" | wc -c)
        
        echo "📄 生成的 HTML 长度: $HTML_LENGTH 字符"
        
        # 创建完整的 HTML 文件
        cat > "$OUTPUT_FILE" << EOF
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>微信公众号文章预览 - curl 测试</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
            padding-bottom: 20px;
            border-bottom: 1px solid #eee;
        }
        .header h1 {
            color: #333;
            margin: 0;
        }
        .header p {
            color: #666;
            margin: 10px 0 0 0;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #999;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📱 微信公众号预览</h1>
            <p>通过 curl 命令调用 API 生成</p>
        </div>
        
        $HTML_CONTENT
        
        <div class="footer">
            <p>生成时间: $(date '+%Y-%m-%d %H:%M:%S')</p>
            <p>测试工具: curl + bash</p>
        </div>
    </div>
</body>
</html>
EOF
        
        ABS_PATH=$(realpath "$OUTPUT_FILE")
        echo "💾 HTML 文件已保存: $ABS_PATH"
        echo "🌐 在浏览器中打开: file://$ABS_PATH"
        
        # 显示 HTML 内容预览
        echo ""
        echo "📋 HTML 内容预览:"
        echo "----------------------------------------"
        echo "$HTML_CONTENT" | head -c 300
        echo "..."
        echo "----------------------------------------"
        
        echo ""
        echo "🎯 测试完成！可以进行以下操作:"
        echo "   1. 打开 $OUTPUT_FILE 查看效果"
        echo "   2. 复制 HTML 内容到微信公众号编辑器"
        echo "   3. 使用不同的 Markdown 内容继续测试"
        
    else
        ERROR_MSG=$(echo "$RESPONSE_BODY" | jq -r '.error')
        echo "❌ 转换失败: $ERROR_MSG"
    fi
    
else
    echo "❌ 请求失败 (HTTP $HTTP_CODE)"
    echo "响应内容: $RESPONSE_BODY"
    
    echo ""
    echo "💡 故障排除建议:"
    echo "   1. 检查服务器是否正常启动"
    echo "   2. 确认 API 地址是否正确: $API_URL"
    echo "   3. 检查网络连接"
fi

echo ""
echo "📋 其他测试方式:"
echo "   • 浏览器测试: 打开 test_api.html"
echo "   • Python 测试: python test_api.py"
echo "   • 手动测试: 访问 http://localhost:8080"
