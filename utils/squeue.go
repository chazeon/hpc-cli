package utils

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/olekukonko/tablewriter"
)

type Job struct {
	Account    string  `json:"account"`
	Name       string  `json:"name"`
	JobID      string  `json:"job_id"`
	WorkDir    string  `json:"work_dir" mapstructure:"WORK_DIR"`
	User       string  `json:"user"`
	Partition  string  `json:"partition"`
	Command    string  `json:"command"`
	State      string  `json:"state"`
	Time       string  `json:"time"`
	TimeLeft   string  `json:"time_left" mapstructure:"TIME_LEFT"`
	SubmitTime string  `json:"submit_time" mapstructure:"SUBMIT_TIME"`
	StartTime  string  `json:"start_time" mapstructure:"START_TIME"`
	EndTime    string  `json:"end_time" mapstructure:"END_TIME"`
	Priority   string  `json:"priority"`
	Machine    Machine `json:"machine"`
}

func ParseJobs(stdout string, machine Machine) (jobs []Job, err error) {

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

		job = Job{
			Machine: machine,
		}

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

func ShowJobs(jobs []Job, fmt string) {

	switch fmt {

	case "json":

		bytes, _ := json.Marshal(jobs)
		println(string(bytes))

	default:
		writer := tablewriter.NewWriter(os.Stdout)

		writer.SetHeader([]string{
			"Machine",
			"JobID",
			"State",
			"Time",
			"WorkDir",
		})

		for _, job := range jobs {
			writer.Append([]string{
				job.Machine.Name,
				job.JobID,
				job.State,
				job.Time,
				job.WorkDir,
			})
		}

		writer.Render()
	}
}
