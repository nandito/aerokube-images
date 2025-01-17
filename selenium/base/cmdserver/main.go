package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sync"
)

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

func main() {
	http.HandleFunc("/ffmpeg-start", ffmpegStartHandler)
	http.HandleFunc("/ffmpeg-stop", ffmpegStopHandler)

	port := "9091"
	fmt.Printf("Starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
