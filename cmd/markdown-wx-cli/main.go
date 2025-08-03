package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"bilibili-uploader/internal/converter"
)

func main() {
	var (
		inputFile  = flag.String("input", "", "Input markdown file path")
		outputFile = flag.String("output", "", "Output HTML file path")
		interactive = flag.Bool("i", false, "Interactive mode")
	)
	flag.Parse()

	conv := converter.NewWechatConverterFixed()

	if *interactive {
		runInteractiveMode(conv)
		return
	}

	var markdown string
	var err error

	if *inputFile != "" {
		markdown, err = readFile(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
	} else {
		markdown, err = readStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
	}

	html := conv.ConvertMarkdownToWechat(markdown)

	if *outputFile != "" {
		err = writeFile(*outputFile, html)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Output written to: %s\n", *outputFile)
	} else {
		fmt.Println(html)
	}
}

func runInteractiveMode(conv *converter.WechatConverter) {
	fmt.Println("微信公众号 Markdown 转换器 - 交互模式")
	fmt.Println("输入 Markdown 文本，按 Ctrl+D (macOS/Linux) 或 Ctrl+Z (Windows) 结束输入")
	fmt.Println("或者输入 'quit' 退出")
	fmt.Println(strings.Repeat("-", 50))

	scanner := bufio.NewScanner(os.Stdin)
	var lines []string

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if line == "quit" {
			break
		}

		if line == "" && len(lines) > 0 {
			// 空行表示结束输入，开始转换
			markdown := strings.Join(lines, "\n")
			html := conv.ConvertMarkdownToWechat(markdown)
			
			fmt.Println("\n" + strings.Repeat("=", 50))
			fmt.Println("转换结果:")
			fmt.Println(strings.Repeat("=", 50))
			fmt.Println(html)
			fmt.Println(strings.Repeat("=", 50))
			
			lines = lines[:0] // 清空
			continue
		}

		lines = append(lines, line)
	}

	if len(lines) > 0 {
		markdown := strings.Join(lines, "\n")
		html := conv.ConvertMarkdownToWechat(markdown)
		
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("转换结果:")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println(html)
		fmt.Println(strings.Repeat("=", 50))
	}
}

func readFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func readStdin() (string, error) {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func writeFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}
