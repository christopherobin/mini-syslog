package main

import (
	"log"
	"os"
	"text/template"

	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/fatih/color"
	"gopkg.in/mcuadros/go-syslog.v2"
)

// colors we use in the app
var (
	yellow = color.New(color.FgHiYellow).SprintFunc()
	blue   = color.New(color.FgHiBlue).SprintFunc()
	red    = color.New(color.FgHiRed).SprintFunc()
	gray   = color.New(color.FgHiBlack).SprintFunc()
	green  = color.New(color.FgHiGreen).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
)

// app flags
var (
	app  = kingpin.New("mini-syslog", "Mini-Syslog with papertrail-like output")
	bind = app.Flag("bind", "What IP:Port to bind to").
		Short('b').
		Default(":1514").
		String()
	protocol = app.Flag("protocol", "Which protocol to use").
			Short('p').
			Default("udp").
			Enum("udp", "tcp", "dgram")
	useColors = app.Flag("color", "Enable colors").
			Default("true").
			Bool()
	customTemplate = app.Flag("template", "Override output template").
			Short('t').
			String()
)

var defaultTemplate = "{{ .timestamp }} [{{ severity .severity }}] {{ yellow .hostname }} {{ or .app_name .tag | blue }}: {{ or .message .content }}"

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.Automatic)
	server.SetHandler(handler)

	var err error
	switch *protocol {
	case "tcp":
		err = server.ListenTCP(*bind)
	case "udp":
		err = server.ListenUDP(*bind)
	case "dgram":
		err = server.ListenUnixgram(*bind)
	default:
		err = server.ListenUDP(*bind)
	}

	if err != nil {
		panic(err)
	}

	server.Boot()

	if !*useColors {
		color.NoColor = true
	}

	t := template.New("output")

	// maps syslog severity to colored strings
	var severityMap = map[int]string{
		7: gray("debug"),
		6: blue("info"),
		5: blue("notice"),
		4: yellow("warning"),
		3: red("error"),
		2: bold(red("critical")),
		1: bold(red("alert")),
		0: bold(red("emergency")),
	}

	t.Funcs(template.FuncMap{
		"severity": func(severity int) string {
			if val, ok := severityMap[severity]; ok {
				return val
			}
			return fmt.Sprintf("unknown (%d)", severity)
		},
		// color functions
		"bold":   bold,
		"blue":   blue,
		"gray":   gray,
		"green":  green,
		"red":    red,
		"yellow": yellow,
	})

	if *customTemplate != "" {
		t = template.Must(t.Parse(*customTemplate))
	} else {
		t = template.Must(t.Parse(defaultTemplate))
	}

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			t.Execute(os.Stdout, logParts)
			os.Stdout.Write([]byte("\n"))
		}
	}(channel)

	log.Printf("server listening on %s (%s) \n", *bind, *protocol)
	server.Wait()
}
