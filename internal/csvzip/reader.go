package csvzip

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"mime/multipart"
)

func ReadCSVFromMultipart(file multipart.File) ([][]string, error) {
	zipData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	r, err := zip.NewReader(
		bytes.NewReader(zipData),
		int64(len(zipData)),
	)
	if err != nil {
		return nil, err
	}

	for _, f := range r.File {
		if f.Name == "data.csv" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			reader := csv.NewReader(rc)
			return reader.ReadAll()
		}
	}

	return nil, errors.New("data.csv not found in archive")
}