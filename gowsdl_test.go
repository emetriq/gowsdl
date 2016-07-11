// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

import (
	"bytes"
	"go/format"
	"path/filepath"
	"strings"
	"testing"
)

const (
	simpleHeader = `package myservice

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name
`
	simpleTypes = `
type GetInfo struct {
	XMLName xml.Name ` + "`xml:\"http://www.mnb.hu/webservices/ GetInfo\"`" + `
}

type GetInfoResponse struct {
	XMLName xml.Name ` + "`xml:\"http://www.mnb.hu/webservices/ GetInfoResponse\"`" + `

	// this is a comment

	GetInfoResult string ` + "`xml:\"GetInfoResult,omitempty\"`" + `
}
`
	simpleOps = `
type TestSoapPort struct {
	client *SOAPClient
}

func NewTestSoapPort(url string, tls bool, auth *BasicAuth, headers ...*HTTPHeader) *TestSoapPort {
	if url == "" {
		url = "http://www.mnb.hu/arfolyamok.asmx"
	}
	client := NewSOAPClient(url, tls, auth, headers)

	return &TestSoapPort{
		client: client,
	}
}
`
)

func TestElementGenerationDoesntCommentOutStructProperty(t *testing.T) {
	g := GoWSDL{
		file:         "testdata/test.wsdl",
		pkg:          "myservice",
		makePublicFn: makePublic,
	}

	resp, err := g.Start()
	if err != nil {
		t.Error(err)
	}

	if strings.Contains(string(resp["types"]), "// this is a comment  GetInfoResult string `xml:\"GetInfoResult,omitempty\"`") {
		t.Error("Type comment should not comment out struct type property")
		t.Error(string(resp["types"]))
	}
}

func TestVboxGeneratesWithoutSyntaxErrors(t *testing.T) {
	files, err := filepath.Glob("testdata/*.wsdl")
	if err != nil {
		t.Error(err)
	}

	for _, file := range files {
		g := GoWSDL{
			file:         file,
			pkg:          "myservice",
			makePublicFn: makePublic,
		}

		resp, err := g.Start()
		if err != nil {
			continue
			//t.Error(err)
		}

		data := new(bytes.Buffer)
		data.Write(resp["header"])
		data.Write(resp["types"])
		data.Write(resp["operations"])
		data.Write(resp["soap"])

		_, err = format.Source(data.Bytes())
		if err != nil {
			t.Error(err)
		}
	}
}

func TestSimpleHeader(t *testing.T) {
	file := "testdata/test.wsdl"

	g := GoWSDL{
		file:         file,
		pkg:          "myservice",
		makePublicFn: makePublic,
	}

	data, err := g.genHeader()
	header, err := format.Source(data)
	if err != nil {
		t.Error(err)
	}

	if strings.Compare(string(header), simpleHeader) != 0 {
		t.Errorf("Expexted\n---\n%s\n---\nbut got\n---\n%s\n---\n", simpleHeader, string(header))
	}
}

func TestSimpleTypes(t *testing.T) {
	file := "testdata/test.wsdl"

	g := GoWSDL{
		file:         file,
		pkg:          "myservice",
		makePublicFn: makePublic,
	}

	err := g.unmarshal()
	if err != nil {
		t.Error(err)
	}
	types, err := g.genTypes()
	if err != nil {
		t.Error(err)
	}
	data := new(bytes.Buffer)
	data.Write([]byte(simpleHeader))
	data.Write(types)

	code, err := format.Source(data.Bytes())
	if err != nil {
		t.Error(err)
	}
	if strings.Compare(string(code), simpleHeader+simpleTypes) != 0 {
		t.Errorf("Expected\n---\n%s\n---\nbut got\n---\n%s\n", simpleHeader+simpleTypes, code)
	}
}

func TestSimpleOperations(t *testing.T) {
	file := "testdata/test.wsdl"

	g := GoWSDL{
		file:         file,
		pkg:          "myservice",
		makePublicFn: makePublic,
	}

	err := g.unmarshal()
	if err != nil {
		t.Error(err)
	}
	types, err := g.genOperations()
	if err != nil {
		t.Error(err)
	}
	data := new(bytes.Buffer)
	data.Write([]byte(simpleHeader))
	data.Write([]byte(simpleTypes))
	data.Write(types)

	code, err := format.Source(data.Bytes())
	if err != nil {
		t.Error(err)
	}
	if strings.Compare(string(code), simpleHeader+simpleTypes+simpleOps) != 0 {
		t.Errorf("Expected\n---\n%s\n---\nbut got\n---\n%s\n", simpleHeader+simpleTypes+simpleOps, code)
	}
}
