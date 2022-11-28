package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
	"time"
)

type waterService struct{}

func NewWaterService() *waterService { return &waterService{} }

func (s *waterService) writeColdData(text string, parTable [2]string, i int) string {

	var data string
	if i == 4 {
		data = text[10:]
	} else {
		data = text[11:]
	}
	//string([]rune(text)[11:])
	coldData, err := strconv.ParseFloat(data, 32)
	fmt.Println(coldData)
	if err != nil {
		return "Incorrect number"
	} else {
		if !s.checkData(coldData, i) {
			return "Check number: new data < prev data"
		} else {
			dt := time.Now()
			date := dt.Format("02.01.2006")
			file, err := excelize.OpenFile("WaterData.xlsx")
			if err != nil {
				log.Println(err)
				return "Not open file"
			}
			defer func() {
				// Close the spreadsheet.
				if err := file.Close(); err != nil {
					log.Println(err)
				}
			}()
			rows, err := file.GetCols("Water data")
			if err != nil {
				log.Fatal(err)
			}
			partDate := parTable[0] + strconv.FormatInt(int64(len(rows[i])+1), 10)
			fmt.Println("Data:", partDate)
			partValueWater := parTable[1] + strconv.FormatInt(int64(len(rows[i+1])+1), 10)
			fmt.Println("water:", partValueWater)
			if len(rows[i]) != len(rows[i+1]) {
				log.Println("Errors date in table")
				if len(rows[i]) > len(rows[i+1]) {
					partValueWater = parTable[1] + strconv.FormatInt(int64(len(rows[i])+1), 10)
				} else {
					partDate = parTable[0] + strconv.FormatInt(int64(len(rows[i+1])+1), 10)
				}
			}
			cell, err := file.GetCellValue("Water data", partValueWater)
			if err != nil {
				log.Println(err)
			}
			fmt.Println("CELL:", cell, len(cell))
			num, err := strconv.ParseInt(partValueWater[1:], 10, 16)
			if err != nil {
				log.Println(err)
			}
			for len(cell) == 0 {
				num--
				numStr := strconv.FormatInt(num, 10)
				partValueWater = string(partValueWater[0]) + numStr
				cell, err = file.GetCellValue("Water data", partValueWater)
				if err != nil {
					break
				}
				fmt.Println(num)
			}
			num++
			numStr := strconv.FormatInt(num, 10)
			partValueWater = string(partValueWater[0]) + numStr
			partDate = string(partDate[0]) + numStr

			file.SetCellFloat("Water data", partValueWater, coldData, 2, 64)
			file.SetCellValue("Water data", partDate, date)
			if err := file.SaveAs("WaterData.xlsx"); err != nil {
				return "Failed to save data :("
			}
		}
	}
	return "Ok, date saved"
}

func (s *waterService) checkData(d1 float64, i int) bool {
	d2 := s.getPrevData(i + 1)
	if d2 == "No previous data" {
		return true
	}
	newD2, err := strconv.ParseFloat(d2, 32)
	if err != nil {
		log.Println(err)
	}
	if d1 < newD2 {
		return false
	}
	return true
}

func (s *waterService) getPrevData(i int) string {
	file, err := excelize.OpenFile("WaterData.xlsx")
	if err != nil {
		log.Println(err)
		return "Not open file"
	}
	defer func() {
		// Close the spreadsheet.
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()
	// Get all the rows in the Water data.
	rows, err := file.GetCols("Water data")
	if err != nil {
		log.Fatal(err)
	}
	coldData := rows[i][len(rows[i])-1]
	if coldData == "" || coldData == "ХВС" {
		coldData = "No previous data"
	}
	return coldData
}
