package utils

import (
	"encoding/xml"
	"log"
)

// Log 结构体表示整个日志
type Log struct {
	XMLName    xml.Name   `xml:"log" json:"-"`
	LogEntries []LogEntry `xml:"logentry" json:"logEntries"`
}

// LogEntry 结构体表示每个日志条目
type LogEntry struct {
	Revision string `xml:"revision,attr" json:"revision,omitempty"`
	Author   string `xml:"author" json:"author,omitempty"`
	Date     string `xml:"date" json:"date,omitempty"`
	Paths    []Path `xml:"paths>path" json:"paths,omitempty"`
	Message  string `xml:"msg" json:"message,omitempty"`
}

// Path 结构体表示每个文件变更信息
type Path struct {
	TextMods bool   `xml:"text-mods,attr" json:"textMods,omitempty"`
	Kind     string `xml:"kind,attr" json:"kind,omitempty"`
	Action   string `xml:"action,attr" json:"action,omitempty"`
	PropMods bool   `xml:"prop-mods,attr" json:"propMods,omitempty"`
	FilePath string `xml:",chardata" json:"filePath,omitempty"`
}

func DecodeSvnXml(xmlData string) (result Log, err error) {
	// 解析 XML
	var logData Log
	err = xml.Unmarshal([]byte(xmlData), &logData)
	if err != nil {
		log.Fatalf("XML 解析失败: %v", err)
	}

	return logData, err
}
