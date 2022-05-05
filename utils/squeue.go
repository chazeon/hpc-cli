package utils

import (
	"encoding/csv"
	"io"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type Job struct {
	Account    string `json:"account"`
	Name       string `json:"name"`
	JobID      string `json:"job_id"`
	WorkDir    string `json:"work_dir" mapstructure:"WORK_DIR"`
	User       string `json:"user"`
	Partition  string `json:"partition"`
	Command    string `json:"command"`
	State      string `json:"state"`
	Time       string `json:"time"`
	TimeLeft   string `json:"time_left" mapstructure:"time_left"`
	SubmitTime string `json:"submit_time" mapstructure:"submit_time"`
	StartTime  string `json:"start_time" mapstructure:"start_time"`
	EndTime    string `json:"end_time" mapstructure:"end_time"`
	Machine    Machine
}

func ParseSqueue(stdout string, machine Machine) (jobs []Job, err error) {

	var (
		keys   []string
		values []string
		record map[string]string
		job    Job
	)

	jobs = make([]Job, 0)
	reader := csv.NewReader(strings.NewReader(stdout))
	reader.Comma = '|'

	keys, err = reader.Read()

	if err == io.EOF {
		err = nil
		return
	}

	if err != nil {
		return
	}

	for {

		values, err = reader.Read()

		if err == io.EOF {
			err = nil
			return
		}

		if err != nil {
			return
		}

		record = make(map[string]string)

		for i, key := range keys {
			record[key] = values[i]
		}

		mapstructure.Decode(record, &job)

		jobs = append(jobs, job)
	}

}
