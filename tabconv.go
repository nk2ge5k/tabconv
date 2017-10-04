package tabconv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
)

// convert convets either xlsx to csv or csv to xlsx
func Convert(filepath, outdir string, delim rune) error {
	var isCSV bool

	if err := checkXlsxFile(filepath); err != nil {
		if err != ErrFormat {
			return errors.Wrap(err, "convert")
		}

		isCSV = true
	}

	var converter func(string, string, rune) error

	converter = xlsxToCsv
	if isCSV {
		converter = csvToXlsx
	}

	if err := converter(filepath, outdir, delim); err != nil {
		return errors.Wrap(err, "convert")
	}

	return nil
}

// csvToXlsx converts file form csv to xlsx
func csvToXlsx(filepath, outdir string, delim rune) error {
	f, err := os.Open(filepath)
	if err != nil {
		return errors.Wrap(err, "csv to xlsx")
	}
	defer f.Close()

	basename := basename(filepath)

	xlFile, err := xlsx.OpenFile()
	if err != nil {
		return errors.Wrap(err, "csv to xlsx")
	}

	var (
		sheet *xlsx.Sheet
		n     = 1
	)
	// loop till we find previosly not existed sheet name
	for {
		name := fmt.Sprintf("Sheet %d", n)
		if _, ok := xlFile.Sheet[name]; !ok {
			var err error

			sheet, err = xlFile.AddSheet(name)
			if err != nil {
				return errors.Wrap(err, "csv to xlsx")
			}
			break
		}
		n++
	}

	r := csv.NewReader(f)
	r.Comma = delim

	if err := copySheetCsv(sheet, r); err != nil {
		return errors.Wrap(err, "csv to xlsx")
	}

	return xlFile.Save(basename + ".xlsx")
}

// xlsxToCsv converts file form xlsx to csv
func xlsxToCsv(filepath, outdir string, delim rune) error {
	xlFile, err := xlsx.OpenFile(filepath)
	if err != nil {
		return errors.Wrap(err, "xlsx to csv")
	}

	// if no sheets in excel file then we finished
	if len(xlFile.Sheets) == 0 {
		return nil
	}

	if len(xlFile.Sheets) > 1 {
		outdir = path.Join(outdir, basename(filepath))
		if err := os.Mkdir(outdir, 0755); err != nil {
			return errors.Wrap(err, "xlsx to csv")
		}
	}

	for _, sheet := range xlFile.Sheets {
		fname := path.Join(outdir, fix(sheet.Name)) + ".csv"

		f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return errors.Wrap(err, "xlsx to csv")
		}

		w := csv.NewWriter(f)
		w.Comma = delim

		err = copyCsvSheet(w, sheet)
		// check close errors first because we must close either way
		if e := f.Close(); e != nil {
			return errors.Wrap(e, "xlsx to csv")
		}

		if err != nil {
			return errors.Wrap(err, "xlsx to csv")
		}
	}

	return nil
}

// copyCsvSheet copies content form the xlsx sheet to csv writer
func copyCsvSheet(dst *csv.Writer, src *xlsx.Sheet) error {
	// allocate string slice for row values
	rowbuf := make([]string, 0, src.MaxCol)

	for _, row := range src.Rows {
		// skip empty lines
		if row == nil {
			continue
		}
		rowbuf = rowbuf[:0]

		for _, cell := range row.Cells {
			val, err := cell.FormattedValue()
			if err != nil {
				return errors.Wrap(err, "failed to format xlsx value")
			}

			rowbuf = append(rowbuf, val)
		}

		if err := dst.Write(rowbuf); err != nil {
			return errors.Wrap(err, "failed to write in csv")
		}
	}

	dst.Flush()

	if err := dst.Error(); err != nil {
		return errors.Wrap(err, "failed to flush csv content")
	}

	return nil
}

// copySheetCsv copies content form csv file to xlsx sheet
func copySheetCsv(dst *xlsx.Sheet, src *csv.Reader) error {
	for {
		rec, err := src.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.Wrap(err, "failed to copy csv to xlsx")
		}

		row := dst.AddRow()
		for _, val := range rec {
			row.AddCell().Value = val
		}
	}
	return nil
}
