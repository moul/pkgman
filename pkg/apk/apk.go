package apk

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/avast/apkparser"
	xml2json "github.com/basgys/goxml2json"
)

type Package struct {
	r *zip.ReadCloser
}

func Open(path string) (*Package, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("open path: %w", err)
	}
	return &Package{r: r}, nil
}

func (p *Package) Close() error {
	return nil
}

func (p Package) Files() []*zip.File {
	return p.r.File
}

func (p Package) File(name string) *zip.File {
	for _, f := range p.r.File {
		if f.Name == name {
			return f
		}
	}
	return nil
}

func (p Package) FileBytes(name string) ([]byte, error) {
	f := p.File(name)
	if f == nil {
		return nil, fmt.Errorf("no such file: %s", name)
	}

	r, err := f.Open()
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(r)
}

func (p Package) ManifestXML() (string, error) {
	b, err := p.FileBytes("AndroidManifest.xml")
	if err != nil {
		return "", err
	}
	binbuf := bytes.NewBuffer(b)

	xmlbuf := bytes.Buffer{}
	xmlwriter := bufio.NewWriter(&xmlbuf)
	enc := xml.NewEncoder(xmlwriter)
	enc.Indent("  ", "    ")
	err = apkparser.ParseXml(binbuf, enc, nil)
	if err != nil {
		return "", err
	}

	return xmlbuf.String(), nil
}

func (p Package) ManifestJSON() (string, error) {
	xmlStr, err := p.ManifestXML()
	if err != nil {
		return "", err
	}

	xml := strings.NewReader(xmlStr)
	json, err := xml2json.Convert(xml)
	if err != nil {
		return "", err
	}

	return json.String(), nil
}

func (p Package) Manifest() (*Manifest, error) {
	xmlStr, err := p.ManifestXML()
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	err = xml.Unmarshal([]byte(xmlStr), &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}
