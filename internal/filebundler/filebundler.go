package internal

import (
	"archive/tar"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"os"
)

//go:embed embedded-files.tar.gz
var embeddedFS embed.FS

func DebundleFiles(outputPath string) {
	decompressTarGzArchive("embedded-files.tar.gz", outputPath)
}

func decompressTarGzArchive(archiveName string, outputPath string) {

	gzipFile, err := embeddedFS.Open(archiveName)
	if err != nil {
		panic(err)
	}
	defer gzipFile.Close()

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		panic(err)
	}

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // Done reading archive
		} else if err != nil {
			panic(err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// Ignore
		case tar.TypeReg, tar.TypeRegA:
			outputFilePath := outputPath + "/" + header.FileInfo().Name()

			file, err := os.Create(outputFilePath)
			if err != nil {
				panic(err)
			}

			err = os.Chmod(outputFilePath, 0777)
			if err != nil {
				panic(err)
			}

			_, err = io.Copy(file, tarReader)
			if err != nil {
				panic(err)
			}

			file.Close()
		default:
			fmt.Printf("Error extracting %s: found type %v in %v\n",
				archiveName, header.Typeflag, header.Name)
		}
	}
}
