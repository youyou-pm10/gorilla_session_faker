package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/sessions"
)

// MockResponseWriter 是一个实现了 http.ResponseWriter 接口的结构体，
// 用于捕获 Set-Cookie 头。
type MockResponseWriter struct {
	headers http.Header
}

func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{headers: make(http.Header)}
}

func (mrw *MockResponseWriter) Header() http.Header {
	return mrw.headers
}

func (mrw *MockResponseWriter) Write([]byte) (int, error) {
	return 0, nil // 我们不需要写入任何内容到响应体中
}

func (mrw *MockResponseWriter) WriteHeader(statusCode int) {
	// 不需要实现，因为我们只关心响应头
}

// simulateSetSession 模拟设置 session 并返回 Set-Cookie 头的内容
func simulateSetSession(key string) string {
	// 创建新的 CookieStore 使用提供的密钥
	store := sessions.NewCookieStore([]byte(key))

	// 创建模拟的 HTTP 请求和响应写入器
	req, _ := http.NewRequest("GET", "/", nil)
	mrw := NewMockResponseWriter()

	// 获取session对象
	session, err := store.Get(req, "session-name")
	if err != nil {
		fmt.Println("Error getting session:", err)
		return ""
	}

	// 设置session值
	session.Values["name"] = "admin"

	// 保存session到响应中
	if err = session.Save(req, mrw); err != nil {
		fmt.Println("Error saving session:", err)
		return ""
	}

	// 返回第一个 Set-Cookie 头
	setCookies := mrw.headers["Set-Cookie"]
	if len(setCookies) > 0 {
		return setCookies[0]
	}
	return "" // 如果没有 Set-Cookie 头，则返回空字符串
}

// readKeysFromFile 从文件中读取密钥列表
func readKeysFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var keys []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		keys = append(keys, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

// writeResultsToFile 将模拟结果写入文件
func writeResultsToFile(filePath string, results []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, result := range results {
		_, err := fmt.Fprintln(writer, result)
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func main() {
	// 从文件读取密钥列表
	keys, err := readKeysFromFile("./keys.txt")
	if err != nil {
		log.Fatalf("Failed to read keys from file: %v", err)
	}

	// 准备存储模拟结果的切片
	var results []string

	for _, key := range keys {
		setCookie := simulateSetSession(key)
		results = append(results, setCookie)
	}

	// 将结果写入输出文件
	err = writeResultsToFile("./output.txt", results)
	if err != nil {
		log.Fatalf("Failed to write results to file: %v", err)
	}

	fmt.Println("Simulation completed and results have been written to output.txt.")
}