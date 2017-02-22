package utils

import (
	"encoding/csv"
	"os"
	"sort"
	"strconv"
	"time"
)

const timestampLayout = "02/01/2006 15:04:05.000"

type Record struct {
	UserID    string
	Timestamp time.Time
	HR        int
	BR        float64
	Activity  float64
}

type AverageRecord struct {
	UserID      string
	Timestamp   time.Time
	Count       int
	HRAvg       float64
	BRAvg       float64
	ActivityAvg float64
}

type ByTimestamp [][]string

func (ts ByTimestamp) Len() int           { return len(ts) }
func (ts ByTimestamp) Swap(i, j int)      { ts[i], ts[j] = ts[j], ts[i] }
func (ts ByTimestamp) Less(i, j int) bool { return ts[i][1] < ts[j][1] }

func ReadRecords(filename string) (*[][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	return &records, nil
}

func WriteAverages(filename string, averages *map[string]AverageRecord) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.Write([]string{"UserID", "Timestamp", "Count", "Average HR", "Average BR", "Average Activity"})
	if err != nil {
		return err
	}

	records := make([][]string, 0, len(*averages))
	for _, avg := range *averages {
		records = append(records, marshalAverageRecord(&avg))
	}
	sort.Sort(ByTimestamp(records))
	err = writer.WriteAll(records)
	if err != nil {
		return err
	}

	return nil
}

func CalculateAverages(records *[][]string) (*map[string]AverageRecord, error) {
	averages := make(map[string]AverageRecord)

	for i := 1; i < len(*records); i++ {
		record, err := unmarshalRecord((*records)[i])
		if err != nil {
			return nil, err
		}

		roundedMinute := record.Timestamp.Truncate(time.Minute)
		key := record.UserID + roundedMinute.String()
		avg, ok := averages[key]

		if ok {
			avg.HRAvg = getNewAverage(avg.HRAvg, float64(record.HR), avg.Count)
			avg.BRAvg = getNewAverage(avg.BRAvg, record.BR, avg.Count)
			avg.ActivityAvg = getNewAverage(avg.ActivityAvg, float64(record.Activity), avg.Count)
			avg.Count += 1
			averages[key] = avg
		} else {
			averages[key] = AverageRecord{
				record.UserID,
				roundedMinute,
				1,
				float64(record.HR),
				record.BR,
				record.Activity,
			}
		}
	}

	return &averages, nil
}

func unmarshalRecord(raw []string) (rec *Record, err error) {
	uID := raw[0]

	ts, err := time.Parse(timestampLayout, raw[1])
	if err != nil {
		return nil, err
	}

	hr, err := strconv.Atoi(raw[2])
	if err != nil {
		return nil, err
	}

	br, err := strconv.ParseFloat(raw[3], 64)
	if err != nil {
		return nil, err
	}

	act, err := strconv.ParseFloat(raw[4], 64)
	if err != nil {
		return nil, err
	}

	return &Record{
		uID,
		ts,
		hr,
		br,
		act,
	}, nil
}

// {"UserID", "Timestamp", "Count", "Average HR", "Average BR", "Average Activity"},
func marshalAverageRecord(ar *AverageRecord) []string {
	return []string{
		ar.UserID,
		ar.Timestamp.String(),
		strconv.Itoa(ar.Count),
		strconv.FormatFloat(ar.HRAvg, 'f', 4, 64),
		strconv.FormatFloat(ar.BRAvg, 'f', 4, 64),
		strconv.FormatFloat(ar.ActivityAvg, 'f', 4, 64),
	}
}

func getNewAverage(current, addition float64, count int) float64 {
	return ((current * float64(count)) + addition) / float64((count + 1))
}
