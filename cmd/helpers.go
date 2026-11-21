package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"nugu.dev/basement/pkg/models"
)

/* Postgres can return duration from AGE() method. This return has []uint8 type.
* This method transforms that to the duration in seconds */

func Fromuint8ToInt(entry []uint8) int {

	entryLen := len(entry)

	total, multiplier := 0, 1

	if entryLen == 0 {
		return total
	}

	for i := entryLen - 1; i >= 0; i-- {

		holder, _ := strconv.Atoi(string(entry[i]))

		switch i {
		case entryLen - 1:
			total += holder * multiplier // 1
			multiplier *= 10
		case entryLen - 2:
			total += holder * multiplier // 10
			multiplier *= 10
		case entryLen - 3:
			multiplier *= 6
			continue
		case entryLen - 4:
			total += holder * multiplier // 60
		case entryLen - 5:
			total += holder * multiplier // 60
		case entryLen - 6:
			multiplier *= 6
			continue
		default:
			total += holder * multiplier
			multiplier *= 10
		}
	}

	return total
}

func From24ClockToTime(date string, hourMinTime string) (time.Time, error) {

	dateTime, err := time.Parse(time.DateOnly, date)

	if err != nil {
		return time.Time{}, err
	}

	holder := strings.Split(hourMinTime, ":")
	entryHours, err := strconv.Atoi(holder[0])

	if err != nil {
		return time.Time{}, err
	}

	entryMinutes, err := strconv.Atoi(holder[1])

	if err != nil {
		return time.Time{}, err
	}

	dateTime = dateTime.Add(time.Hour * time.Duration(entryHours))
	dateTime = dateTime.Add(time.Minute * time.Duration(entryMinutes))

	return dateTime, nil
}

func ReadChangelogFile() (*models.Changelog, error) {
	var c *models.Changelog

	f, err := os.ReadFile("changelog.yaml")

	if err != nil {
		return nil, fmt.Errorf("Error reading changelog file")
	}

	if err := yaml.Unmarshal(f, &c); err != nil {
		return nil, fmt.Errorf("Error reading changelog file")
	}

	return c, nil
}
