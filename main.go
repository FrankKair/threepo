package main

import (
    "archive/zip"
    "encoding/csv"
    "encoding/json"
    "encoding/xml"
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

type sharedStrings struct {
    SI []struct {
        T string `xml:"t"`
    } `xml:"si"`
}

type worksheet struct {
    Rows []struct {
        Cells []struct {
            T string `xml:"t,attr"` 
            V string `xml:"v"`
        } `xml:"c"`
    } `xml:"sheetData>row"`
}

func die(format string, args ...any) {
    fmt.Fprintf(os.Stderr, format+"\n", args...)
    os.Exit(1)
}

func main() {
    if len(os.Args) != 2 {
        die("usage: threepo <file.xlsx|file.csv>")
    }
    path := os.Args[1]

    var rows[][]string
    var err error
    switch strings.ToLower(filepath.Ext(path)) {
    case ".xlsx":
        rows, err = readXLSX(path)
    case ".csv":
        rows, err = readCSV(path)
    default:
        die("unsupported format: %s (expected .xlsx or .csv)", filepath.Ext(path))
    }
    if err != nil {
        die("%v", err)
    }

    out, _ := json.MarshalIndent(buildObject(rows), "", " ")
    os.Stdout.Write(append(out, '\n'))
}

func readCSV(path string) ([][]string, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    return csv.NewReader(f).ReadAll()
}

func readXLSX(path string) ([][]string, error) {
    r, err := zip.OpenReader(path)
    if err != nil {
        return nil, err
    }
    defer r.Close()

    var ss sharedStrings
    readZipXML(r, "xl/sharedStrings.xml", &ss)
    sst := make([]string, len(ss.SI))
    for i, si := range ss.SI {
        sst[i] = si.T
    }

    var ws worksheet
    if err := readZipXML(r, "xl/worksheets/sheet1.xml", &ws); err != nil {
        return nil, err
    }

    var rows [][]string
    for _, row := range ws.Rows {
        var vals []string
        for _, c := range row.Cells {
            v := c.V
            if c.T == "s" {
                if idx, err := strconv.Atoi(v); err == nil && idx < len(sst) {
                    v = sst[idx]
                }
            }
            vals = append(vals, v)
        }
        rows = append(rows, vals)
    }
    return rows, nil
}

func readZipXML(r *zip.ReadCloser, name string, v any) error {
    for _, f := range r.File {
        if f.Name == name {
            rc, err := f.Open()
            if err != nil {
                return err
            }
            defer rc.Close()
            return xml.NewDecoder(rc).Decode(v)
        }
    }
    return fmt.Errorf("%s not found", name)
}

func buildObject(rows [][]string) map[string]any {
    if len(rows) < 2 {
        return map[string]any{}
    }
    headers := rows[0]
    keyIdx := -1
    for i, h := range headers {
        if h == "key" {
            keyIdx = i
            break
        }
    }
    if keyIdx == -1 {
        die("no 'key' column found in header row")
    }
    result := map[string]any{}
    for _, row := range rows[1:] {
        vals := make(map[string]string)
        key := ""
        for i, v := range row {
            if i == keyIdx {
                key = v
            } else if i < len(headers) {
                vals[headers[i]] = v
            }
        }
        if key != "" {
            setNested(result, key, vals)
        }
    }
    return result
}

func setNested(m map[string]any, key string, value any) {
    parts := strings.Split(key, ".")
    for _, p := range parts[:len(parts)-1] {
        child, ok := m[p]
        if !ok {
            child = map[string]any{}
            m[p] = child
        }
        m = child.(map[string]any)
    }
    m[parts[len(parts)-1]] = value
}
