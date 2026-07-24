package release

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type archiveFile struct {
	Name string
	Mode os.FileMode
	Data []byte
}

func writeArchive(path string, target Target, files []archiveFile, timestamp time.Time) error {
	sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })
	if target.OS == "windows" {
		return writeZip(path, files, timestamp)
	}
	return writeTarGzip(path, files, timestamp)
}

func writeZip(path string, files []archiveFile, timestamp time.Time) error {
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	for _, file := range files {
		header := &zip.FileHeader{Name: file.Name, Method: zip.Deflate}
		header.SetMode(file.Mode)
		header.SetModTime(timestamp.UTC())
		entry, err := writer.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("create zip entry %s: %w", file.Name, err)
		}
		if _, err := entry.Write(file.Data); err != nil {
			return fmt.Errorf("write zip entry %s: %w", file.Name, err)
		}
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("close zip archive: %w", err)
	}
	return writeFileAtomic(path, buffer.Bytes(), 0o644)
}

func writeTarGzip(path string, files []archiveFile, timestamp time.Time) error {
	var buffer bytes.Buffer
	gzipWriter, err := gzip.NewWriterLevel(&buffer, gzip.BestCompression)
	if err != nil {
		return fmt.Errorf("create gzip writer: %w", err)
	}
	gzipWriter.Name = ""
	gzipWriter.Comment = ""
	gzipWriter.ModTime = timestamp.UTC()
	gzipWriter.OS = 255
	tarWriter := tar.NewWriter(gzipWriter)
	for _, file := range files {
		header := &tar.Header{
			Name:       file.Name,
			Mode:       int64(file.Mode.Perm()),
			Size:       int64(len(file.Data)),
			ModTime:    timestamp.UTC(),
			AccessTime: time.Time{},
			ChangeTime: time.Time{},
			Uid:        0,
			Gid:        0,
			Uname:      "",
			Gname:      "",
			Format:     tar.FormatPAX,
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("write tar header %s: %w", file.Name, err)
		}
		if _, err := tarWriter.Write(file.Data); err != nil {
			return fmt.Errorf("write tar entry %s: %w", file.Name, err)
		}
	}
	if err := tarWriter.Close(); err != nil {
		return fmt.Errorf("close tar archive: %w", err)
	}
	if err := gzipWriter.Close(); err != nil {
		return fmt.Errorf("close gzip archive: %w", err)
	}
	return writeFileAtomic(path, buffer.Bytes(), 0o644)
}

func writeFileAtomic(path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".setup-env-release-*")
	if err != nil {
		return err
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Chmod(mode); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Rename(temporaryName, path); err != nil {
		return err
	}
	return nil
}

func readTarGzip(reader io.Reader) (map[string][]byte, error) {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()
	result := map[string][]byte{}
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(tarReader)
		if err != nil {
			return nil, err
		}
		result[header.Name] = data
	}
	return result, nil
}
