#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
微信公众号 Markdown 转换 API 测试脚本

使用方法:
1. 确保服务器正在运行: go run main.go
2. 运行此脚本: python test_api.py
3. 查看生成的 output.html 文件
"""

import requests
import json
import sys
from pathlib import Path

# API 配置
API_URL = "http://localhost:8080/api/convert"

# 测试用的 Markdown 内容
SAMPLE_MARKDOWN = """# 微信公众号文章标题

欢迎来到我的微信公众号！这篇文章将展示各种 Markdown 语法在微信公众号中的呈现效果。

## 📝 文本格式

在微信公众号中，我们可以使用各种文本格式：

- **粗体文本**: 用于强调重要内容
- *斜体文本*: 用于标注或引用
- `行内代码`: 用于显示代码片段或技术术语
- ~~删除线~~: 用于表示已删除或过时的内容

## 🔗 链接和引用

访问我的 [个人博客](https://example.com) 了解更多技术文章。

> 💡 **重要提示**: 这是一个引用块，通常用于强调重要信息或引用他人观点。在微信公众号中会有特别的样式处理。

## 📋 列表示例

### 技术栈
1. **前端**: HTML, CSS, JavaScript
2. **后端**: Go, Python, Node.js  
3. **数据库**: MySQL, Redis, MongoDB
4. **工具**: Git, Docker, VS Code

### 学习计划
- [ ] 完成 Go 基础教程
- [x] 学习微信公众号开发
- [ ] 掌握 Docker 容器技术
- [x] 搭建个人博客

## 💻 代码示例

### JavaScript 异步处理
```javascript
// 使用 async/await 处理异步操作
async function fetchUserData(userId) {
    try {
        const response = await fetch(`/api/users/${userId}`);
        const userData = await response.json();
        
        console.log('用户信息:', userData);
        return userData;
    } catch (error) {
        console.error('获取用户信息失败:', error);
        throw error;
    }
}

// 调用函数
fetchUserData(123).then(user => {
    console.log(`欢迎, ${user.name}!`);
});
```

### Python 数据处理
```python
import pandas as pd
import numpy as np

# 数据处理示例
def analyze_sales_data(file_path):
    # 读取数据
    df = pd.read_csv(file_path)
    
    # 数据清洗
    df = df.dropna()
    df['销售额'] = pd.to_numeric(df['销售额'], errors='coerce')
    
    # 统计分析
    summary = {
        '总销售额': df['销售额'].sum(),
        '平均销售额': df['销售额'].mean(),
        '最高销售额': df['销售额'].max(),
        '销售记录数': len(df)
    }
    
    return summary

# 使用示例
result = analyze_sales_data('sales_data.csv')
print(f"分析结果: {result}")
```

## 📊 表格数据

| 编程语言 | 难度等级 | 应用领域 | 推荐指数 |
|---------|---------|---------|---------|
| Python | ⭐⭐⭐ | 数据科学、AI | ⭐⭐⭐⭐⭐ |
| JavaScript | ⭐⭐ | Web开发 | ⭐⭐⭐⭐⭐ |
| Go | ⭐⭐⭐ | 后端服务 | ⭐⭐⭐⭐ |
| Rust | ⭐⭐⭐⭐⭐ | 系统编程 | ⭐⭐⭐ |

## 🧮 数学公式

### 行内公式
爱因斯坦质能方程: $E = mc^2$

### 块级公式
二次方程求根公式:
$$x = \\frac{-b \\pm \\sqrt{b^2 - 4ac}}{2a}$$

概率论中的贝叶斯定理:
$$P(A|B) = \\frac{P(B|A) \\cdot P(A)}{P(B)}$$

## 🎯 总结

通过这个 Markdown 编辑器，我们可以：

1. **轻松编写**: 使用熟悉的 Markdown 语法
2. **实时预览**: 即时查看微信公众号效果  
3. **一键复制**: 直接粘贴到公众号编辑器
4. **多种主题**: 选择合适的视觉风格

---

**关注我的公众号，获取更多技术干货！** 🚀

*本文使用微信公众号 Markdown 编辑器创建*"""

def convert_markdown(markdown_content):
    """
    调用 API 转换 Markdown 内容
    """
    headers = {
        "Content-Type": "application/json"
    }
    
    payload = {
        "markdown": markdown_content
    }
    
    try:
        print("🔄 正在调用转换 API...")
        response = requests.post(API_URL, 
                               data=json.dumps(payload), 
                               headers=headers,
                               timeout=30)
        response.raise_for_status()
        
        result = response.json()
        
        if result['success']:
            print("✅ 转换成功！")
            return result['html']
        else:
            print(f"❌ 转换失败: {result['error']}")
            return None
            
    except requests.exceptions.ConnectionError:
        print("🔌 连接失败: 请确保服务器正在运行 (http://localhost:8080)")
        print("   运行命令: go run main.go")
        return None
    except requests.exceptions.Timeout:
        print("⏰ 请求超时: 服务器响应时间过长")
        return None
    except requests.exceptions.RequestException as e:
        print(f"📡 请求失败: {e}")
        return None
    except json.JSONDecodeError:
        print("📄 响应格式错误: 服务器返回的不是有效的 JSON")
        return None

def save_html_file(html_content, filename="output.html"):
    """
    保存 HTML 内容到文件
    """
    try:
        # 创建完整的 HTML 文档
        full_html = f"""<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>微信公众号文章预览</title>
    <style>
        body {{
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background: #f5f5f5;
        }}
        .container {{
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }}
        .header {{
            text-align: center;
            margin-bottom: 30px;
            padding-bottom: 20px;
            border-bottom: 1px solid #eee;
        }}
        .header h1 {{
            color: #333;
            margin: 0;
        }}
        .header p {{
            color: #666;
            margin: 10px 0 0 0;
        }}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📱 微信公众号预览</h1>
            <p>以下是转换后的文章效果</p>
        </div>
        {html_content}
    </div>
</body>
</html>"""
        
        with open(filename, 'w', encoding='utf-8') as f:
            f.write(full_html)
        
        file_path = Path(filename).resolve()
        print(f"💾 HTML 文件已保存: {file_path}")
        print(f"🌐 在浏览器中打开: file://{file_path}")
        return True
        
    except Exception as e:
        print(f"💥 保存文件失败: {e}")
        return False

def main():
    """
    主函数
    """
    print("🚀 微信公众号 Markdown 转换 API 测试")
    print("=" * 50)
    
    # 检查是否有命令行参数指定 Markdown 文件
    if len(sys.argv) > 1:
        markdown_file = sys.argv[1]
        try:
            with open(markdown_file, 'r', encoding='utf-8') as f:
                markdown_content = f.read()
            print(f"📖 从文件读取内容: {markdown_file}")
        except FileNotFoundError:
            print(f"❌ 文件不存在: {markdown_file}")
            return
        except Exception as e:
            print(f"❌ 读取文件失败: {e}")
            return
    else:
        markdown_content = SAMPLE_MARKDOWN
        print("📝 使用内置示例内容")
    
    print(f"📏 Markdown 内容长度: {len(markdown_content)} 字符")
    print()
    
    # 调用转换 API
    html_content = convert_markdown(markdown_content)
    
    if html_content:
        print(f"📄 转换后 HTML 长度: {len(html_content)} 字符")
        
        # 保存到文件
        if save_html_file(html_content):
            print()
            print("🎉 转换完成！")
            print("📋 接下来你可以:")
            print("   1. 打开生成的 HTML 文件查看效果")
            print("   2. 复制 HTML 内容到微信公众号编辑器")
            print("   3. 或者使用浏览器测试页面: test_api.html")
        
        # 显示部分 HTML 内容预览
        print()
        print("📋 HTML 内容预览:")
        print("-" * 30)
        preview = html_content[:200] + "..." if len(html_content) > 200 else html_content
        print(preview)
        print("-" * 30)
    else:
        print()
        print("💡 故障排除:")
        print("   1. 检查服务器是否启动: go run main.go")
        print("   2. 确认端口 8080 未被占用")
        print("   3. 检查防火墙设置")

if __name__ == "__main__":
    main()
