package easyconf

import (
	"os"
	"path/filepath"
	"testing"
)

// TestJsonConfSaveAtomic 覆盖 Save 原子写改造的关键行为：
// 旧实现直开覆写且无 O_TRUNC——新内容比旧内容短时文件尾部残留旧字节，JSON 随即损坏。
func TestJsonConfSaveAtomic(t *testing.T) {
	dir := t.TempDir()
	jc := NewJsonConf(dir)

	// 先写长内容
	long := map[string]any{
		"name": "a-very-long-config",
		"list": []string{"a", "b", "c", "d", "e"},
		"desc": "padding-padding-padding-padding-padding",
	}
	if err := jc.Save(long, "conf.json"); err != nil {
		t.Fatalf("写入长内容失败: %v", err)
	}

	// 再写短内容：序列化后明显短于旧内容，旧实现必残留尾部旧字节
	short := map[string]any{"name": "x"}
	if err := jc.Save(short, "conf.json"); err != nil {
		t.Fatalf("写入短内容失败: %v", err)
	}
	var got map[string]any
	if err := jc.Read(&got, "conf.json"); err != nil {
		t.Fatalf("内容变短后读回失败（残留旧字节致 JSON 损坏）: %v", err)
	}
	if got["name"] != "x" {
		t.Fatalf("读回数据不一致: got %v", got)
	}

	// 反复覆盖后仍可正常读回
	again := map[string]any{"name": "y", "n": 2}
	if err := jc.Save(again, "conf.json"); err != nil {
		t.Fatalf("再次覆盖失败: %v", err)
	}
	var got2 map[string]any
	if err := jc.Read(&got2, "conf.json"); err != nil {
		t.Fatalf("反复覆盖后读回失败: %v", err)
	}
	if got2["name"] != "y" {
		t.Fatalf("反复覆盖后读回数据不一致: got %v", got2)
	}

	// 不残留临时文件
	matches, err := filepath.Glob(filepath.Join(dir, "*.tmp-*"))
	if err != nil {
		t.Fatalf("Glob 失败: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("Save 后残留临时文件: %v", matches)
	}
}

// TestJsonConfSaveEmptyDir 覆盖 dirPath 为空串的场景（Join 后为当前目录文件）。
func TestJsonConfSaveEmptyDir(t *testing.T) {
	jc := NewJsonConf("")
	name := "easyconf-jsonfile-test.json"
	defer func() { _ = os.Remove(name) }()
	if err := jc.Save(map[string]string{"k": "v"}, name); err != nil {
		t.Fatalf("空 dirPath 保存失败: %v", err)
	}
	var got map[string]string
	if err := jc.Read(&got, name); err != nil {
		t.Fatalf("空 dirPath 读回失败: %v", err)
	}
	if got["k"] != "v" {
		t.Fatalf("读回数据不一致: got %v", got)
	}
}
