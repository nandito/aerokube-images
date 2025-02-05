package main

import (
	"bufio"
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

    // save something to /Users/nandito/Workspace/oss/aerokube/images/selenium/base/cmdserver
    command := "ffmpeg -f lavfi -i anoisesrc=d=0:c=pink -c:a pcm_s16le /Users/nandito/Workspace/oss/aerokube/images/selenium/base/cmdserver/output.wav"
	// command := fmt.Sprintf(
	// 	"ffmpeg -f pulse -thread_queue_size 2048 -i default -y -f x11grab -video_size %s -r %s %s -i %s -codec:v %s %s -filter:v \"pad=ceil(iw/2)*2:ceil(ih/2)*2\" /home/selenium/videooutput/%s",
	// 	params.VideoSize, params.FrameRate, params.InputOptions, params.Display, params.Codec, params.Preset, params.FileName,
	// )

	cmd := exec.Command("bash", "-c", command)
	// if err := cmd.Start(); err != nil {
	// 	http.Error(w, fmt.Sprintf("Failed to start FFmpeg: %v", err), http.StatusInternalServerError)
	// 	return
	// }

    // Create pipes to capture stdout and stderr
    stdoutPipe, err := cmd.StdoutPipe()
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to create stdout pipe: %v", err), http.StatusInternalServerError)
        return
    }

    stderrPipe, err := cmd.StderrPipe()
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to create stderr pipe: %v", err), http.StatusInternalServerError)
        return
    }

    // Start the command
    if err := cmd.Start(); err != nil {
        http.Error(w, fmt.Sprintf("Failed to start FFmpeg: %v", err), http.StatusInternalServerError)
        return
    }

    // Log output in goroutines
    go func() {
        scanner := bufio.NewScanner(stdoutPipe)
        for scanner.Scan() {
            fmt.Printf("[FFmpeg stdout] %s\n", scanner.Text())
        }
    }()

    go func() {
        scanner := bufio.NewScanner(stderrPipe)
        for scanner.Scan() {
            fmt.Printf("[FFmpeg stderr] %s\n", scanner.Text())
        }
    }()

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
