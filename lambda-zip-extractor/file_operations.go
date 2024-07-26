package main

import (
	"bytes"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
)

func extractPDFContent(fileContent []byte) (string, error) {
	tmpFile, err := os.CreateTemp("", "temp-*.pdf")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write(fileContent); err != nil {
		return "", err
	}

	f, r, err := pdf.Open(tmpFile.Name())
	if err != nil {
		return "", err
	}
	defer f.Close()

	var content strings.Builder
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		text, _ := p.GetPlainText(nil)
		content.WriteString(text)
	}

	return content.String(), nil
}

func extractDOCXContent(fileContent []byte) (string, error) {
	reader := bytes.NewReader(fileContent)
	r, err := docx.ReadDocxFromMemory(reader, int64(len(fileContent)))
	if err != nil {
		return "", err
	}
	defer r.Close()

	return r.Editable().GetContent(), nil
}

func extractTextContent(fileContent []byte) (string, error) {
	return string(fileContent), nil
}

func sanitizeUTF8(s string) string {
	return strings.Map(func(r rune) rune {
		if r == 0x00 {
			return -1 
		}
		if r == utf8.RuneError {
			return -1
		}
		return r
	}, s)
}
