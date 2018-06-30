package core

import (
	"bytes"
	"encoding/json"
	"os"
)

//FileType EnumType for filetype
type FileType int

// Enum for FileType
const (
	FILE = 1
	DIR  = 2
)

// FTypeName Map of fileType to Its String name
var FTypeName = map[FileType]string{
	FILE: "file",
	DIR:  "dir",
}

// FTypeIds Map of String Name to Its FileType
var FTypeIds = map[string]FileType{
	"file": FILE,
	"dir":  DIR,
}

// String Implementation of Stringer interface
func (ftype FileType) String() string {
	return FTypeName[ftype]
}

// Ref: https://gist.github.com/lummie/7f5c237a17853c031a57277371528e87

//MarshalJSON Marshal FileType as JSON String
func (ftype FileType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(FTypeName[ftype])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

//UnmarshalJSON UnMarshal FileType JSON String as FileType
func (ftype *FileType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	// lookup value
	*ftype = FTypeIds[s]
	return nil
}

//GetFileType Utility method to get fileType
func GetFileType(fi os.FileInfo) FileType {
	if fi.IsDir() {
		return DIR
	}
	return FILE
}
