# `hpc-cli`

Work with multiple HPCs over SSH.

## Usage

```
NAME:
   hpc-cli - A new cli application

USAGE:
   hpc-cli [global options] command [command options] [arguments...]

COMMANDS:
   exec       Execute the command.
   list-jobs  List all the running jobs.
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE         Load configuration from FILE (default: "config.yml")
   --machine MACHINE, -m MACHINE  Specify the MACHINE to apply on.  (accepts multiple inputs)
   --help, -h                     show help (default: false)
```

### List slurm jobs

```bash
# List running jobs in a table.
hpc-cli list-jobs

# List running jobs in json format then format with `jq`.
hpc-cli list-jobs | jq

# List running jobs on expanse and bridges2.
hpc-cli -m expanse -m bridges2 list-jobs
                                    
```

### Execute command

```bash
# Show the home directory at different machines.
hpc-cli exec pwd

# Show current running jobs on different machines.
hpc-cli exec 'squeue -u $(whoami)'  
```
