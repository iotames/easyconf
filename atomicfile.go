package easyconf

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// writeFileAtomic 先写同目录临时文件、落盘同步、再改名替换目标文件。
//
// 为什么不用「O_TRUNC 直接覆写」：那是就地截断重写，写到一半掉电、磁盘满或进程被杀，
// 目标文件会留下半截内容。配置文件承载的是运行期状态的唯一落点，半截文件读回来是错的，
// 而错误往往到「重启后配置全丢」才暴露。临时文件 + rename 是原子的：
// 读者要么看到旧文件，要么看到新文件，不会看到半截。
//
// rename 后再对父目录 Sync 一次：POSIX 上 rename 的持久性依赖目录项落盘，
// 缺这一步掉电后可能回退到旧文件（本次更新丢失）。Windows 无目录句柄同步语义，跳过。
//
// 已存在的文件沿用原权限（配置文件可能含口令，不能因一次写入放宽）；新建文件使用 defaultMode。
// 临时文件建在目标同目录，要求目录可写——这与「写目标文件」的权限要求一致。
func writeFileAtomic(fpath string, data []byte, defaultMode os.FileMode) error {
	mode := defaultMode
	if st, err := os.Stat(fpath); err == nil {
		mode = st.Mode().Perm()
	}

	// 临时文件名带目标文件名，便于排障时看出归属；同目录保证 rename 不跨文件系统。
	tmp, err := os.CreateTemp(filepath.Dir(fpath), "."+filepath.Base(fpath)+".tmp-*")
	if err != nil {
		return fmt.Errorf("create temp file for(%s) err(%w)", fpath, err)
	}
	tmpName := tmp.Name()
	// 改名成功后该路径已不存在；这里是中途失败时的兜底清理，避免留下半截临时文件。
	defer os.Remove(tmpName)

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("write temp file(%s) err(%w)", tmpName, err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("sync temp file(%s) err(%w)", tmpName, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file(%s) err(%w)", tmpName, err)
	}
	if err := os.Chmod(tmpName, mode); err != nil {
		return fmt.Errorf("chmod temp file(%s) err(%w)", tmpName, err)
	}
	if err := os.Rename(tmpName, fpath); err != nil {
		return fmt.Errorf("replace file(%s) err(%w)", fpath, err)
	}
	syncParentDir(fpath)
	return nil
}

// syncParentDir 让父目录项落盘，保证 rename 持久。仅 POSIX 有意义；
// Windows 打开目录句柄即可写盘跟踪，无对应 API，直接跳过。失败只影响掉电后的持久性，
// 不影响本次写入的正确性，静默忽略。
func syncParentDir(fpath string) {
	if runtime.GOOS == "windows" {
		return
	}
	dir, err := os.Open(filepath.Dir(fpath))
	if err != nil {
		return
	}
	defer dir.Close()
	_ = dir.Sync()
}
