package main

import (
	"fmt"

	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/format"
)

func init() {
	format.RegisterAll()
}

func main() {
	file, _ := avutil.Open("test.mp4")
	streams, _ := file.Streams()
	for _, stream := range streams {
		if stream.Type().IsAudio() {
			astream := stream.(av.AudioCodecData)
			fmt.Println(astream.Type(), astream.SampleRate(), astream.SampleFormat(), astream.ChannelLayout())
		} else if stream.Type().IsVideo() {
			vstream := stream.(av.VideoCodecData)
			fmt.Println(vstream.Type(), vstream.Width(), vstream.Height())
		}
	}
	file.Close()
}
