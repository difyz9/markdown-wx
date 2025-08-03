package converter

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

// WechatConverter 微信公众号Markdown转换器
type WechatConverter struct {
	footnotes  []string
	styles     WechatStyles
	codeBlocks map[string]string
}

// WechatStyles 微信公众号样式定义
type WechatStyles struct {
	H1Style         string
	H2Style         string
	H3Style         string
	ParagraphStyle  string
	QuoteStyle      string
	CodeBlockStyle  string
	InlineCodeStyle string
	ListStyle       string
	LinkStyle       string
	ImageStyle      string
	TableStyle      string
	TableHeaderStyle string
	TableCellStyle   string
}

// NewWechatConverterFixed 创建新的转换器
func NewWechatConverterFixed() *WechatConverter {
	return &WechatConverter{
		footnotes:  make([]string, 0),
		styles:     getDefaultStyles(),
		codeBlocks: make(map[string]string),
	}
}

// getDefaultStyles 获取默认微信公众号样式
func getDefaultStyles() WechatStyles {
	return WechatStyles{
		H1Style: `style="display: table; text-align: center; color: #3f3f3f; line-height: 1.75; font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif; font-size: 18px; font-weight: bold; margin: 2em auto 1em; padding: 0 1em; border-bottom: 3px solid #009874; margin-top: 0;"`,
		
		H2Style: `style="display: table; text-align: center; color: #fff; line-height: 1.75; font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif; font-size: 16px; font-weight: bold; margin: 4em auto 2em; padding: 0 0.3em; background: #009874;"`,
		
		H3Style: `style="text-align: left; color: #3f3f3f; line-height: 1.2; font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif; font-size: 14px; font-weight: bold; margin: 2em 8px 0.75em 0; padding-left: 8px; border-left: 5px solid #009874;"`,
		
		ParagraphStyle: `style="font-size: 16px; line-height: 1.5em; padding: 0.5em 0; margin: 0; color: initial;"`,
		
		QuoteStyle: `style="text-align: left; font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif; font-size: 14px; font-style: normal; border-left: none; padding: 0.5em 1em; background: rgba(27, 31, 35, 0.05); margin: 1em 0;"`,
		
		CodeBlockStyle: `style="display: block; padding: 1em; color: rgb(51, 51, 51); background: rgb(248, 248, 248); font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; font-weight: 400; letter-spacing: normal; orphans: 2; text-indent: 0px; text-transform: none; widows: 2; word-spacing: 0px; text-decoration-style: initial; text-decoration-color: initial; text-align: left; line-height: 1.5; font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif; margin: 0.9rem 0; white-space: pre;"`,
		
		InlineCodeStyle: `style="text-align: left; line-height: 1; white-space: initial; color: #333; background: rgba(27, 31, 35, 0.05); padding: 0.1em 0.3em; font-weight: bold; font-size: 1em; top: -0.1em; position: relative;"`,
		
		ListStyle: `style="padding-left: 1.2em;"`,
		
		LinkStyle: `style="color: #009874; text-decoration: none; font-size: 14px;"`,
		
		ImageStyle: `style="display: initial; max-width: 100%;"`,
		
		TableStyle: `style="width: 100%; border-collapse: collapse; line-height: 1.35; font-size: 14px;"`,
		
		TableHeaderStyle: `style="background: rgb(0 0 0 / 5%); border: 1px solid #ddd; padding: 0.25em 0.5em;"`,
		
		TableCellStyle: `style="border: 1px solid #ddd; padding: 0.25em 0.5em;"`,
	}
}

// ConvertMarkdownToWechat 将Markdown转换为微信公众号格式
func (c *WechatConverter) ConvertMarkdownToWechat(markdown string) string {
	// 重置脚注
	c.footnotes = make([]string, 0)
	
	// 预处理：处理特殊格式
	html := c.preprocessText(markdown)
	
	// 预处理：处理代码块，避免其他规则干扰
	html = c.extractCodeBlocks(html)
	
	// 转换各种元素
	html = c.processHeaders(html)
	html = c.processQuotes(html)
	html = c.processTables(html)
	html = c.processLists(html)
	html = c.processLinks(html)
	html = c.processImages(html)
	html = c.processInlineCode(html)
	html = c.processBoldItalic(html)
	html = c.processParagraphs(html)
	
	// 恢复代码块
	html = c.restoreCodeBlocks(html)
	
	// 添加脚注
	if len(c.footnotes) > 0 {
		html += c.generateFootnotes()
	}
	
	return html
}

// extractCodeBlocks 提取代码块并用占位符替换
func (c *WechatConverter) extractCodeBlocks(text string) string {
	codeBlockRegex := regexp.MustCompile("(?s)```(.*?)\n(.*?)```")
	counter := 0
	
	result := codeBlockRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholder := fmt.Sprintf("__CODE_BLOCK_%d__", counter)
		
		// 提取语言和代码
		matches := codeBlockRegex.FindStringSubmatch(match)
		code := matches[2]
		
		// 生成微信公众号代码块HTML（不显示语言标识符）
		escapedCode := html.EscapeString(code)
		codeHTML := fmt.Sprintf(`<section %s>%s</section>`,
			c.styles.CodeBlockStyle, escapedCode)
		
		c.codeBlocks[placeholder] = codeHTML
		counter++
		return placeholder
	})
	
	return result
}

// restoreCodeBlocks 恢复代码块
func (c *WechatConverter) restoreCodeBlocks(text string) string {
	for placeholder, codeHTML := range c.codeBlocks {
		text = strings.ReplaceAll(text, placeholder, codeHTML)
	}
	return text
}

// processHeaders 处理标题
func (c *WechatConverter) processHeaders(text string) string {
	// H1
	h1Regex := regexp.MustCompile(`(?m)^# (.+)$`)
	text = h1Regex.ReplaceAllStringFunc(text, func(match string) string {
		title := strings.TrimSpace(strings.TrimPrefix(match, "#"))
		return fmt.Sprintf(`<h1 %s>%s</h1>`, c.styles.H1Style, html.EscapeString(title))
	})
	
	// H2
	h2Regex := regexp.MustCompile(`(?m)^## (.+)$`)
	text = h2Regex.ReplaceAllStringFunc(text, func(match string) string {
		title := strings.TrimSpace(strings.TrimPrefix(match, "##"))
		return fmt.Sprintf(`<h2 %s>%s</h2>`, c.styles.H2Style, html.EscapeString(title))
	})
	
	// H3
	h3Regex := regexp.MustCompile(`(?m)^### (.+)$`)
	text = h3Regex.ReplaceAllStringFunc(text, func(match string) string {
		title := strings.TrimSpace(strings.TrimPrefix(match, "###"))
		return fmt.Sprintf(`<h3 %s>%s</h3>`, c.styles.H3Style, html.EscapeString(title))
	})
	
	return text
}

// processTables 处理表格
func (c *WechatConverter) processTables(text string) string {
	lines := strings.Split(text, "\n")
	var result []string
	var tableLines []string
	inTable := false
	
	for i, line := range lines {
		line = strings.TrimSpace(line)
		
		// 检测表格开始（包含 | 分隔符的行）
		if strings.Contains(line, "|") && !inTable {
			// 检查下一行是否是分隔符行
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				if isTableSeparator(nextLine) {
					inTable = true
					tableLines = []string{line} // 只保存表头，不保存分隔符行
					i++ // 跳过分隔符行
					continue
				}
			}
		}
		
		// 在表格中
		if inTable {
			if strings.Contains(line, "|") && !isTableSeparator(line) {
				tableLines = append(tableLines, line)
			} else {
				// 表格结束
				result = append(result, c.convertTable(tableLines))
				tableLines = []string{}
				inTable = false
				
				// 添加当前行
				if line != "" {
					result = append(result, line)
				}
			}
		} else {
			result = append(result, line)
		}
	}
	
	// 处理末尾的表格
	if inTable && len(tableLines) > 0 {
		result = append(result, c.convertTable(tableLines))
	}
	
	return strings.Join(result, "\n")
}

// isTableSeparator 检查是否是表格分隔符行
func isTableSeparator(line string) bool {
	// 移除空格
	line = strings.ReplaceAll(line, " ", "")
	// 检查是否只包含 |, -, : 字符
	for _, char := range line {
		if char != '|' && char != '-' && char != ':' {
			return false
		}
	}
	return strings.Contains(line, "-") && strings.Contains(line, "|")
}

// convertTable 转换表格
func (c *WechatConverter) convertTable(tableLines []string) string {
	if len(tableLines) < 1 {
		return strings.Join(tableLines, "\n")
	}
	
	var tableHTML strings.Builder
	tableHTML.WriteString(fmt.Sprintf(`<table %s>`, c.styles.TableStyle))
	
	// 处理表头
	headerLine := strings.Trim(tableLines[0], "|")
	headerCells := strings.Split(headerLine, "|")
	
	tableHTML.WriteString("<tr>")
	for _, cell := range headerCells {
		cell = strings.TrimSpace(cell)
		if cell != "" {
			tableHTML.WriteString(fmt.Sprintf(`<th %s>%s</th>`, c.styles.TableHeaderStyle, html.EscapeString(cell)))
		}
	}
	tableHTML.WriteString("</tr>")
	
	// 处理数据行（从第1行开始，因为第0行是表头）
	for i := 1; i < len(tableLines); i++ {
		line := strings.Trim(tableLines[i], "|")
		cells := strings.Split(line, "|")
		
		tableHTML.WriteString("<tr>")
		for j, cell := range cells {
			cell = strings.TrimSpace(cell)
			
			// 处理特殊样式的单元格
			cellStyle := c.styles.TableCellStyle
			
			// 如果是代码类型的单元格，使用特殊样式
			if strings.HasPrefix(cell, "-") || strings.HasPrefix(cell, "`") {
				cellStyle = strings.Replace(cellStyle, "color: #333", "color: #d63384; font-family: 'SFMono-Regular', Consolas, monospace", 1)
			}
			
			// 如果是第一列且包含特殊标记，添加特殊样式
			if j == 0 && (strings.HasPrefix(cell, "-") || strings.Contains(cell, "login") || strings.Contains(cell, "upload")) {
				cellStyle = strings.Replace(cellStyle, "color: #333", "color: #d63384; font-weight: bold; font-family: 'SFMono-Regular', Consolas, monospace", 1)
			}
			
			tableHTML.WriteString(fmt.Sprintf(`<td %s>%s</td>`, cellStyle, html.EscapeString(cell)))
		}
		tableHTML.WriteString("</tr>")
	}
	
	tableHTML.WriteString("</table>")
	return tableHTML.String()
}

// processQuotes 处理引用
func (c *WechatConverter) processQuotes(text string) string {
	quoteRegex := regexp.MustCompile(`(?m)^> (.+)$`)
	text = quoteRegex.ReplaceAllStringFunc(text, func(match string) string {
		content := strings.TrimSpace(strings.TrimPrefix(match, ">"))
		return fmt.Sprintf(`<blockquote %s>%s</blockquote>`, c.styles.QuoteStyle, html.EscapeString(content))
	})
	return text
}

// processLists 处理列表
func (c *WechatConverter) processLists(text string) string {
	// 处理 • 符号开头的列表项
	bulletRegex := regexp.MustCompile(`(?m)^•(.+)$`)
	text = bulletRegex.ReplaceAllStringFunc(text, func(match string) string {
		content := strings.TrimSpace(strings.TrimPrefix(match, "•"))
		return fmt.Sprintf(`<ul %s><li style="margin: 0; line-height: 1.5em; font-size: 14px;">%s</li></ul>`, c.styles.ListStyle, html.EscapeString(content))
	})
	
	// 无序列表（- * + 开头）
	unorderedRegex := regexp.MustCompile(`(?m)^[-*+] (.+)$`)
	text = unorderedRegex.ReplaceAllStringFunc(text, func(match string) string {
		content := regexp.MustCompile(`^[-*+] `).ReplaceAllString(match, "")
		return fmt.Sprintf(`<ul %s><li style="margin: 0; line-height: 1.5em; font-size: 14px;">%s</li></ul>`, c.styles.ListStyle, html.EscapeString(content))
	})
	
	// 有序列表（数字. 开头）
	orderedRegex := regexp.MustCompile(`(?m)^\d+\.(.+)$`)
	text = orderedRegex.ReplaceAllStringFunc(text, func(match string) string {
		content := regexp.MustCompile(`^\d+\.`).ReplaceAllString(match, "")
		content = strings.TrimSpace(content)
		return fmt.Sprintf(`<ol %s><li style="margin: 0; line-height: 1.5em; font-size: 14px;">%s</li></ol>`, c.styles.ListStyle, html.EscapeString(content))
	})
	
	return text
}

// processLinks 处理链接，转换为脚注
func (c *WechatConverter) processLinks(text string) string {
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	
	return linkRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := linkRegex.FindStringSubmatch(match)
		linkText := matches[1]
		linkURL := matches[2]
		
		// 添加到脚注
		footnoteIndex := len(c.footnotes) + 1
		c.footnotes = append(c.footnotes, linkURL)
		
		// 返回带脚注标记的文本
		return fmt.Sprintf(`<span %s>%s</span><sup>[%d]</sup>`, 
			c.styles.LinkStyle, html.EscapeString(linkText), footnoteIndex)
	})
}

// processImages 处理图片
func (c *WechatConverter) processImages(text string) string {
	imageRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	
	return imageRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := imageRegex.FindStringSubmatch(match)
		altText := matches[1]
		imageURL := matches[2]
		
		return fmt.Sprintf(`<img %s src="%s" alt="%s" />`, 
			c.styles.ImageStyle, imageURL, html.EscapeString(altText))
	})
}

// processInlineCode 处理行内代码
func (c *WechatConverter) processInlineCode(text string) string {
	inlineCodeRegex := regexp.MustCompile("`([^`]+)`")
	
	return inlineCodeRegex.ReplaceAllStringFunc(text, func(match string) string {
		code := strings.Trim(match, "`")
		return fmt.Sprintf(`<code %s>%s</code>`, c.styles.InlineCodeStyle, html.EscapeString(code))
	})
}

// processBoldItalic 处理粗体和斜体
func (c *WechatConverter) processBoldItalic(text string) string {
	// 粗体
	boldRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	text = boldRegex.ReplaceAllString(text, `<strong>$1</strong>`)
	
	// 斜体
	italicRegex := regexp.MustCompile(`\*([^*]+)\*`)
	text = italicRegex.ReplaceAllString(text, `<em>$1</em>`)
	
	return text
}

// processParagraphs 处理段落
func (c *WechatConverter) processParagraphs(text string) string {
	lines := strings.Split(text, "\n")
	var result []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// 跳过已经是HTML标签的行
		if strings.HasPrefix(line, "<") {
			result = append(result, line)
		} else {
			// 包装为段落
			result = append(result, fmt.Sprintf(`<p %s>%s</p>`, c.styles.ParagraphStyle, line))
		}
	}
	
	return strings.Join(result, "\n")
}

// generateFootnotes 生成脚注
func (c *WechatConverter) generateFootnotes() string {
	if len(c.footnotes) == 0 {
		return ""
	}
	
	var footnoteHTML strings.Builder
	footnoteHTML.WriteString(`<hr style="margin: 30px 0; border: none; border-top: 1px solid #eee;" />`)
	
	// 脚注标题样式，参考 wxmp 项目
	footnoteHTML.WriteString(`<h2 style="display: table; font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif; font-size: 14px; font-weight: bold; margin: 3em 0 0.6em 0; padding-left: 0.2em;">参考</h2>`)
	
	for i, footnote := range c.footnotes {
		footnoteHTML.WriteString(fmt.Sprintf(`<p style="font-size: 10px; font-style: italic; line-height: 1.2; margin: 0.4rem 0;">[%d] %s</p>`, i+1, html.EscapeString(footnote)))
	}
	
	return footnoteHTML.String()
}

// preprocessText 预处理文本，处理特殊格式
func (c *WechatConverter) preprocessText(text string) string {
	lines := strings.Split(text, "\n")
	var result []string
	
	for i, line := range lines {
		// 处理带空行的编号列表格式，如 "1.\n\n内容"
		if regexp.MustCompile(`^\d+\.$`).MatchString(strings.TrimSpace(line)) {
			// 查找下一个非空行
			for j := i + 1; j < len(lines); j++ {
				nextLine := strings.TrimSpace(lines[j])
				if nextLine != "" {
					// 合并编号和内容
					result = append(result, strings.TrimSpace(line)+" "+nextLine)
					// 跳过已处理的行
					for k := i + 1; k <= j; k++ {
						if k < len(lines) {
							lines[k] = "" // 标记为已处理
						}
					}
					break
				}
			}
		} else if line != "" { // 只添加非空行或未被标记的行
			result = append(result, line)
		}
	}
	
	return strings.Join(result, "\n")
}
