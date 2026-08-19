package service

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type LauncherService interface {
	OpenURL(rawURL string)
	OpenApplication(appName string, desktopEntry string)
	ExtractURLs(text string) []string
}

type DesktopLauncherService struct{}

func NewDesktopLauncherService() *DesktopLauncherService {
	return &DesktopLauncherService{}
}

func (s *DesktopLauncherService) OpenURL(rawURL string) {
	log.Printf("[LauncherService] 開始嘗試開啟網址: %s", rawURL)

	go func() {
		// 方法 1: xdg-open
		cmd1 := exec.Command("xdg-open", rawURL)
		cmd1.Env = os.Environ()
		out1, err1 := cmd1.CombinedOutput()
		if err1 == nil {
			log.Printf("[LauncherService] xdg-open 成功開啟: %s", rawURL)
			return
		}
		log.Printf("[LauncherService] xdg-open 失敗 (%v), 輸出: %s", err1, string(out1))

		// 方法 2: firefox
		cmd2 := exec.Command("firefox", rawURL)
		cmd2.Env = os.Environ()
		out2, err2 := cmd2.CombinedOutput()
		if err2 == nil {
			log.Printf("[LauncherService] firefox 成功開啟: %s", rawURL)
			return
		}
		log.Printf("[LauncherService] firefox 失敗 (%v), 輸出: %s", err2, string(out2))

		// 方法 3: gio open
		cmd3 := exec.Command("gio", "open", rawURL)
		cmd3.Env = os.Environ()
		out3, err3 := cmd3.CombinedOutput()
		if err3 == nil {
			log.Printf("[LauncherService] gio open 成功開啟: %s", rawURL)
			return
		}
		log.Printf("[LauncherService] gio open 失敗 (%v), 輸出: %s", err3, string(out3))
	}()
}

func (s *DesktopLauncherService) ExtractURLs(text string) []string {
	re := regexp.MustCompile(`https?://[^\s<>"']+|www\.[^\s<>"']+\.[^\s<>"']+`)
	matches := re.FindAllString(text, -1)
	var urls []string
	seen := make(map[string]bool)
	for _, m := range matches {
		url := m
		if strings.HasPrefix(m, "www.") {
			url = "https://" + m
		}
		if !seen[url] {
			seen[url] = true
			urls = append(urls, url)
		}
	}
	return urls
}

func (s *DesktopLauncherService) OpenApplication(appName string, desktopEntry string) {
	target := appName
	if desktopEntry != "" {
		target = desktopEntry
	}

	if isSystemApp(target) {
		return
	}

	go tryFocusOrLaunch(target)
}

func isSystemApp(appName string) bool {
	systemApps := map[string]bool{
		"system": true, "notify-send": true, "networkmanager": true,
		"power-profiles-daemon": true, "xdg-desktop-portal-hyprland": true,
	}
	return systemApps[appName]
}

func tryFocusOrLaunch(appName string) {
	target := strings.ToLower(appName)

	psCmd := exec.Command("sh", "-c", fmt.Sprintf("ps aux | grep -i %s | grep -v grep", target))
	err := psCmd.Run()

	if err == nil {
		log.Printf("[LauncherService] 發現正在執行的進程: %s，嘗試喚醒...", target)
	} else {
		log.Printf("[LauncherService] 未發現進程: %s，嘗試啟動...", target)
	}

	execCmd := findAndExecDesktop(target)
	if execCmd != "" {
		log.Printf("[LauncherService] 執行指令: %s", execCmd)
		go exec.Command("sh", "-c", execCmd+" &").Run()
	} else {
		go exec.Command("go-notify-manager").Run()
	}
}

func findAndExecDesktop(appName string) string {
	searchPaths := []string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		filepath.Join(os.Getenv("HOME"), ".local/share/applications"),
		"/var/lib/flatpak/exports/share/applications",
	}

	target := strings.ToLower(appName)
	exactName := target + ".desktop"

	for _, path := range searchPaths {
		fullPath := filepath.Join(path, exactName)
		if cmd := extractExecFromDesktopFile(fullPath); cmd != "" {
			log.Printf("[LauncherService] 從精確匹配 %s 找到執行指令: %s", exactName, cmd)
			return cmd
		}
	}

	for _, path := range searchPaths {
		files, err := os.ReadDir(path)
		if err != nil {
			continue
		}
		for _, f := range files {
			nameLower := strings.ToLower(f.Name())
			if !strings.HasSuffix(nameLower, ".desktop") {
				continue
			}
			base := strings.TrimSuffix(nameLower, ".desktop")
			if base == target || strings.HasSuffix(base, "."+target) {
				fullPath := filepath.Join(path, f.Name())
				if cmd := extractExecFromDesktopFile(fullPath); cmd != "" {
					log.Printf("[LauncherService] 從 ID 匹配 %s 找到執行指令: %s", f.Name(), cmd)
					return cmd
				}
			}
		}
	}

	for _, path := range searchPaths {
		files, err := os.ReadDir(path)
		if err != nil {
			continue
		}
		for _, f := range files {
			nameLower := strings.ToLower(f.Name())
			if !strings.HasSuffix(nameLower, ".desktop") {
				continue
			}
			if strings.Contains(nameLower, target) {
				if !strings.Contains(target, "developer") && strings.Contains(nameLower, "developer") {
					continue
				}
				if !strings.Contains(target, "nightly") && strings.Contains(nameLower, "nightly") {
					continue
				}
				fullPath := filepath.Join(path, f.Name())
				if cmd := extractExecFromDesktopFile(fullPath); cmd != "" {
					log.Printf("[LauncherService] 從模糊匹配 %s 找到執行指令: %s", f.Name(), cmd)
					return cmd
				}
			}
		}
	}

	return target
}

func extractExecFromDesktopFile(fullPath string) string {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Exec=") {
			execLine := strings.TrimPrefix(line, "Exec=")
			return cleanExecLine(execLine)
		}
	}
	return ""
}

func cleanExecLine(line string) string {
	re := regexp.MustCompile(`%[uUfFiIcKknNvV]`)
	cleaned := re.ReplaceAllString(line, "")
	return strings.TrimSpace(cleaned)
}
