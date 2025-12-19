package csvzip

import (
	"archive/zip"
	"encoding/csv"
	"io"
)

func WriteCSVToZip(
	w io.Writer,
	filename string,
	header []string,
	rows [][]string,
) error {
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	fileWriter, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(fileWriter)
	defer csvWriter.Flush()

	if err := csvWriter.Write(header); err != nil {
		return err
	}

	for _, row := range rows {
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}
