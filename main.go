package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kristofmic/minute_averages/utils"
)

// const defaultDirInputPath = "./input"
// const defaultDirOutputPath = "./output"

func main() {
	start := time.Now()

	argDirInputPath := os.Args[1]
	argDirOutputPath := os.Args[2]
	if argDirInputPath == "" {
		log.Fatal("Relative input path of CSV files must be defined")
	}
	if argDirOutputPath == "" {
		log.Fatal("Relative output path of CSV files must be defined")
	}

	dirInputAbsPath, err := filepath.Abs(argDirInputPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error reading input directory %s: %s", dirInputAbsPath, err))
	}

	files, err := ioutil.ReadDir(dirInputAbsPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error reading CSV files in directory %s: %s", dirInputAbsPath, err))
	}

	dirOutputAbsPath, err := filepath.Abs(argDirOutputPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error reading output directory %s: %s", dirOutputAbsPath, err))
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".csv" {
			processFile(file.Name(), dirInputAbsPath, dirOutputAbsPath)
		}
	}

	log.Print(fmt.Sprintf("Completed calculating averages (took %s)", time.Since(start).String()))
}

func processFile(fileName, inputPath, outputPath string) {
	records, err := utils.ReadRecords(inputPath + "/" + fileName)
	if err != nil {
		log.Fatal("Error reading CSV: ", err)
	}

	averages, err := utils.CalculateAverages(records)
	if err != nil {
		log.Fatal("Error calculating averages: ", err)
	}

	outputFileName := strings.Replace(fileName, ".csv", ".AVERAGE.csv", 1)
	outputFilePath := outputPath + "/" + outputFileName

	log.Print("Writing averages to: ", outputFilePath)

	err = utils.WriteAverages(outputFilePath, averages)
	if err != nil {
		log.Fatal("Error writing output: ", err)
	}
}
