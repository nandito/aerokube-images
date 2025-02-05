package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

func main() {
	dir, err := downloadsDir()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":8090", mux(dir)))
}

type FFMpegParams struct {
	VideoSize    string `json:"video_size"`
	FrameRate    string `json:"frame_rate"`
	InputOptions string `json:"input_options"`
	Display      string `json:"display"`
	Codec        string `json:"codec"`
	Preset       string `json:"preset"`
	FileName     string `json:"file_name"`
}

var (
	ffmpegCmd   *exec.Cmd
	ffmpegMutex sync.Mutex
)

func ffmpegStartHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ffmpegStartHandler")
	ffmpegMutex.Lock()
	defer ffmpegMutex.Unlock()

	if ffmpegCmd != nil {
		http.Error(w, "FFmpeg is already running", http.StatusConflict)
		return
	}

	var params FFMpegParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	command := fmt.Sprintf(
		"ffmpeg -f pulse -thread_queue_size 2048 -i default -y -f x11grab -video_size %s -r %s %s -i %s -codec:v %s %s -filter:v \"pad=ceil(iw/2)*2:ceil(ih/2)*2\" /home/selenium/videooutput/%s",
		params.VideoSize, params.FrameRate, params.InputOptions, params.Display, params.Codec, params.Preset, params.FileName,
	)

	cmd := exec.Command("bash", "-c", command)
	if err := cmd.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to start FFmpeg: %v", err), http.StatusInternalServerError)
		return
	}

	ffmpegCmd = cmd
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("FFmpeg started successfully"))
}

func ffmpegStopHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ffmpegStopHandler")
	ffmpegMutex.Lock()
	defer ffmpegMutex.Unlock()

	if ffmpegCmd == nil {
		http.Error(w, "FFmpeg is not running", http.StatusNotFound)
		return
	}

	if err := ffmpegCmd.Process.Kill(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to stop FFmpeg: %v", err), http.StatusInternalServerError)
		return
	}

	ffmpegCmd = nil
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("FFmpeg stopped successfully"))
}

func downloadsDir() (string, error) {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = "/home/selenium"
	}
	dir := filepath.Join(homeDir, "Downloads")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create downloads dir: %v", err)
	}
	return dir, nil
}

const jsonParam = "json"

func mux(dir string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			deleteFileIfExists(w, r, dir)
			return
		}
		if _, ok := r.URL.Query()[jsonParam]; ok {
			listFilesAsJson(w, dir)
			return
		}
		http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
	})
	mux.HandleFunc("/ffmpeg-start", ffmpegStartHandler)
	mux.HandleFunc("/ffmpeg-stop", ffmpegStopHandler)
	return mux
}

func listFilesAsJson(w http.ResponseWriter, dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	files := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		files = append(files, info)
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ret := []string{}
	for _, f := range files {
		ret = append(ret, f.Name())
	}
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ret)
}

func deleteFileIfExists(w http.ResponseWriter, r *http.Request, dir string) {
	fileName := strings.TrimPrefix(r.URL.Path, "/")
	filePath := filepath.Join(dir, fileName)
	_, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unknown file %s", fileName), http.StatusNotFound)
		return
	}
	err = os.Remove(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete file %s: %v", fileName, err), http.StatusInternalServerError)
		return
	}
}
