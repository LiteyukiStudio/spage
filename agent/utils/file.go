package utils

import (
	"fmt"
	"github.com/LiteyukiStudio/spage/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var ReleaseName = "release.zip"
var HashPath = "release_hash.spage"

func UpdateSiteHash(sitePath string) (string, error) {
	// 构建 release.zip 完整路径
	siteZipPath := filepath.Join(sitePath, ReleaseName)

	// 检查文件是否存在
	if _, err := os.Stat(siteZipPath); os.IsNotExist(err) {
		return "", fmt.Errorf("release文件不存在: %s", siteZipPath)
	}

	// 计算文件哈希
	hash, err := utils.FileHash(siteZipPath)
	if err != nil {
		return "", fmt.Errorf("计算文件哈希失败: %v", err)
	}

	// 构建哈希文件路径
	hashFilePath := filepath.Join(sitePath, HashPath)

	// 将哈希值写入文件
	err = os.WriteFile(hashFilePath, []byte(hash), 0644)
	if err != nil {
		return "", fmt.Errorf("写入哈希文件失败: %v", err)
	}

	return hash, nil
}

func GetHash(sitePath string) (string, time.Time, error) {
	// 构建哈希文件完整路径
	hashFilePath := filepath.Join(sitePath, HashPath)

	// 获取文件信息以获取修改时间
	fileInfo, err := os.Stat(hashFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", time.Time{}, fmt.Errorf("哈希文件不存在: %s", hashFilePath)
		}
		return "", time.Time{}, fmt.Errorf("获取哈希文件信息失败: %v", err)
	}

	// 读取哈希文件内容
	hashBytes, err := ioutil.ReadFile(hashFilePath)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("读取哈希文件失败: %v", err)
	}

	// 返回哈希值和修改时间
	return string(hashBytes), fileInfo.ModTime(), nil
}
