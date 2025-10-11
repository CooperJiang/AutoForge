#!/bin/bash

# 检查embed文件系统中的文件
# 此脚本会构建一个临时程序来列出embed.FS中的所有文件

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

echo "=== 检查 embed.FS 中的文件 ==="
echo ""

# 创建临时Go程序
cat > /tmp/check_embed.go <<'EOF'
package main

import (
	"fmt"
	"io/fs"
	"template/internal/static"
)

func main() {
	fmt.Println("=== Web FS 文件列表 ===")
	webFS := static.GetWebDistFS()

	err := fs.WalkDir(webFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fmt.Printf("  %s\n", path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking web FS: %v\n", err)
	}

	fmt.Println("\n=== 检查特定文件 ===")
	testFiles := []string{
		"index.html",
		"assets/_plugin-vue_export-helper-xTsIm-oa.js",
	}

	for _, file := range testFiles {
		f, err := webFS.Open(file)
		if err != nil {
			fmt.Printf("  ✗ %s - NOT FOUND: %v\n", file, err)
		} else {
			f.Close()
			fmt.Printf("  ✓ %s - OK\n", file)
		}
	}
}
EOF

echo "构建检查程序..."
go run /tmp/check_embed.go

rm -f /tmp/check_embed.go

echo ""
echo "=== 物理文件检查 ==="
echo "internal/static/web/ 目录下的 assets 文件："
find internal/static/web/assets -name "*.js" | head -10
