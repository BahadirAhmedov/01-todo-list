package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
	"github.com/mergestat/timediff"
)

func main(){

	
	headers := []string{"ID", "Task", "Created", "Done"}
	writingHeadsToCsv(headers)

	// ooo, _ := os.Open("storage.csv")
	// newReader := csv.NewReader(ooo)
	// reder, _ := newReader.ReadAll()
	// fmt.Println(len(reder))
	
	// id := numberOfLines()
	Time := timeCreated()

	if len(os.Args) < 2{   
		fmt.Println("Enter command!")
		return
	}else{

	}
	// THE ADD COMMAND
	if os.Args[1] == "add" {
		if len(os.Args) > 3 {
			panic("Invalid")
		}
		
		id := numberOfLines()

		csvForm := []string{strconv.Itoa(id), os.Args[2], Time, "False"}

		appendToCsv(csvForm)

		accecingRecs := accecingRecsWithHead()
		
		fmt.Println("Len: ", len(accecingRecs))
		fmt.Println("Id: ", accecingRecs[len(accecingRecs)-1][0])

		writer := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)
		// Writing headers to the command line
		writer.Write(
			[]byte(fmt.Sprintf("%s\t%s\t%s\n", accecingRecs[0][0], accecingRecs[0][1], accecingRecs[0][2])),
		)
		parsedTime, _ := time.Parse(time.RFC3339, accecingRecs[id][2])

		humanReadableTime := timediff.TimeDiff(parsedTime)	
		writer.Write(
			[]byte(fmt.Sprintf("%s\t%s\t%s\n", accecingRecs[id][0], accecingRecs[id][1], humanReadableTime)),
		)	
		writer.Flush()
	// THE LIST COMMAND
	}else if os.Args[1] == "list" {
		accecingRecs := accecingRecsWithHead()
		writer := tabwriter.NewWriter(os.Stdout,  0, 2, 4, ' ', 0)
		
		if len(os.Args) == 3 {
			if os.Args[2] == "-a" || os.Args[2] == "--all" {
				writer.Write(
					[]byte(fmt.Sprintf("%s\t%s\t%s\t%s\n", accecingRecs[0][0], accecingRecs[0][1], accecingRecs[0][2], accecingRecs[0][3])),
				)
				for i := 1; i <= len(accecingRecs)-1; i++{
						parsedTime, _ := time.Parse(time.RFC3339, accecingRecs[i][2])
						humanReadableTime := timediff.TimeDiff(parsedTime)
						writer.Write(
							[]byte(fmt.Sprintf("%s\t%s\t%s\t%s\n", accecingRecs[i][0], accecingRecs[i][1], humanReadableTime, accecingRecs[i][3])),
						)
				}
			}else{
				panic("invalid")
			}	
			
		}else{
			writer.Write(
				[]byte(fmt.Sprintf("%s\t%s\t%s\n", accecingRecs[0][0], accecingRecs[0][1], accecingRecs[0][2])),
			)	
			for i := 1; i <= len(accecingRecs)-1; i++{
				if accecingRecs[i][3] == "False"{
					parsedTime, _ := time.Parse(time.RFC3339, accecingRecs[i][2])
					humanReadableTime := timediff.TimeDiff(parsedTime)
					writer.Write(
						[]byte(fmt.Sprintf("%s\t%s\t%s\n", accecingRecs[i][0], accecingRecs[i][1], humanReadableTime)),
					)
				}
			}
		}
		writer.Flush()
	}else if os.Args[1] == "complete" {

		accecingRecs := accecingRecsWithHead()
		index, _ := strconv.Atoi(os.Args[2])

		currTask := accecingRecs[index][0:len(accecingRecs[1])-1]	
		currTaskTrue := append(currTask, "True")
		
		writer := tabwriter.NewWriter(os.Stdout,  0, 2, 4, ' ', 0)

		writer.Write(   
			[]byte(fmt.Sprintf("%s\t%s\t%s\t%s\n", accecingRecs[0][0], accecingRecs[0][1], accecingRecs[0][2], accecingRecs[0][3])),
		)
		parsedTime, _ := time.Parse(time.RFC3339, currTaskTrue[2])
		humanReadableTime := timediff.TimeDiff(parsedTime)	
		
		writer.Write(
			[]byte(fmt.Sprintf("%s\t%s\t%s\t%s\n", currTaskTrue[0], currTaskTrue[1], humanReadableTime, currTaskTrue[3])),
		)

		writer.Flush()

		file, _ := loadFile("storage.csv")
		w := csv.NewWriter(file)
		
		for i := 0; i < len(accecingRecs); i++ {
			w.Write(accecingRecs[i])				
		}

		defer file.Close()
		defer w.Flush()
		

	}else if os.Args[1] == "delete" {
		accecingRecs := accecingRecords()
		removedRecStored, err := removeRecord(accecingRecs)
		
		if err != nil {
			fmt.Println(err)
		}else{
			os.Truncate("storage.csv", 0)		
		
			//WRITING TO THE FILE
			file, _ := loadFile("storage.csv")
			defer file.Close()
			w := csv.NewWriter(file)
			headers := []string{"ID", "Task", "Created", "Done"}
			w.Write(headers)
			w.WriteAll(removedRecStored)
			defer w.Flush()
		}

	}	
}

func removeRecord(records [][]string)([][]string, error){
	
	var index int
	flag := 0
	
	for i := 0; i < len(records); i++ {
		if records[i][0] == os.Args[2] {
			index += i
			flag += 1
			break
		}
	}
	if flag == 1 {
		head := records[:index]		
		tail := records[index+1:]		
		record := append(head, tail...)
		
		var newArr [][]string

		counter := 0
		for i := 0; i < len(record); i++ {
			counter++
			numStr := strconv.Itoa(counter)
			_= numStr
			
			finalRec := replaceId(numStr, record[i][1:])


			newArr = append(newArr, finalRec)
		}

		return  newArr, nil
	}else{
		return nil, errors.New("there is no task with id")
	}		
}

func replaceId(id string, task []string)([]string){
	 head := task
	 tail := id
	slice1 := append(head[3:], tail)
	
	finalSlice := append(slice1, head...)

	return finalSlice
}


func accecingRecsWithHead()([][]string){
	file, _  := loadFile("storage.csv")	
	r := csv.NewReader(file)
	records, _ := r.ReadAll()	
	defer file.Close()
	return records
}

func accecingRecords()([][]string){
	file, _  := loadFile("storage.csv")	
	r := csv.NewReader(file)
	r.Read()
	records, _ := r.ReadAll()	
	defer file.Close()
	return records
}

func timeCreated()(string){
	createdAt := time.Now()
	timeStr := createdAt.Format(time.RFC3339)
	return timeStr
}

func numberOfLines()(int){
	file, _  := loadFile("storage.csv")	
	r := csv.NewReader(file)
	records, _ := r.ReadAll()	
	task_id := len(records)
	defer file.Close()
	return task_id
}

func writingHeadsToCsv(headers []string){
	file, _  := loadFile("storage.csv")	
	w := csv.NewWriter(file)
	
	defer file.Close()
 	
	w.Write(headers)
	defer w.Flush()
}

func appendToCsv(task []string) {
	file, _ := loadFileAppend("storage.csv")
	w := csv.NewWriter(file)
	defer file.Close()
 	
	w.Write(task)
	defer w.Flush()
}

