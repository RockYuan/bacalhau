---
sidebar_label: "DuckDB"
sidebar_position: 3
---
# Using Bacalhau with DuckDB

DuckDB is a relational table-oriented database management system that supports SQL queries for producing analytical results. It also comes with various features that are useful for data analytics.

DuckDB is suited for the following use cases:

1. Processing and storing tabular datasets, e.g. from CSV or Parquet files
2. Interactive data analysis, e.g., joining and aggregating multiple large tables
3. Concurrent large changes to multiple large tables, e.g., appending rows, adding/removing/updating columns.
4. Large result set transfer to the client

In this example tutorial, we will show how to use DuckDB with Bacalhau. The advantage of using DuckDB with Bacalhau is that you don’t need to install it, and there is no need to download the datasets since they are already available on IPFS or the web.


## Prerequisites

To get started, you need to install the Bacalhau client, see more information [here](../../../getting-started/installation.md)

## Containerize Script using Docker

:::info
You can skip this entirely and directly go to running on Bacalhau.
:::

If you want any additional dependencies to be installed along with DuckDB, you need to build your own container.

To build your own docker container, create a `Dockerfile`, which contains instructions to build your DuckDB docker container.


```Dockerfile
FROM mcr.microsoft.com/vscode/devcontainers/python:3.9

RUN apt-get update && apt-get install -y nodejs npm g++

# Install dbt
RUN pip3 --disable-pip-version-check --no-cache-dir install duckdb==0.4.0 dbt-duckdb==1.1.4 \
    && rm -rf /tmp/pip-tmp

# Install duckdb cli
RUN wget https://github.com/duckdb/duckdb/releases/download/v0.4.0/duckdb_cli-linux-amd64.zip \
    && unzip duckdb_cli-linux-amd64.zip -d /usr/local/bin \
    && rm duckdb_cli-linux-amd64.zip

# Configure Workspace
ENV DBT_PROFILES_DIR=/workspaces/datadex
WORKDIR /workspaces/datadex

```

:::info
See more information on how to containerize your script/app [here](https://docs.docker.com/get-started/02_our_app/)
:::


### Build the container

We will run the `docker build` command to build the container:

```
docker build -t <hub-user>/<repo-name>:<tag> .
```

Before running the command replace;

**`hub-user`** with your docker hub username, If you don’t have a docker hub account [follow these instructions to create docker account](https://docs.docker.com/docker-id/), and use the username of the account you created

**`repo-name`** with the name of the container, you can name it anything you want

**`tag`** this is not required but you can use the latest tag

In our case

```bash
docker build -t davidgasquez/datadex:v0.2.0
```

### Push the container

Next, upload the image to the registry. This can be done by using the Docker hub username, repo name or tag.

```
docker push <hub-user>/<repo-name>:<tag>
```

In our case

```bash
docker push davidgasquez/datadex:v0.2.0
```

## Running a Bacalhau Job

After the repo image has been pushed to Docker Hub, we can now use the container for running on Bacalhau. To submit a job, run the following Bacalhau command:


```bash
%%bash --out job_id
bacalhau docker run \
    --workdir /inputs/ \
    --wait \
    --id-only \
    davidgasquez/datadex:v0.2.0 \
    -- duckdb -s "select 1"
```

### Structure of the command

Let's look closely at the command above:

1. `bacalhau docker run`: call to Bacalhau
1. `davidgasquez/datadex:v0.2.0`: the name and the tag of the docker image we are using
1. `/inputs/`: path to input dataset
1. `duckdb -s "select 1"`: execute DuckDB


When a job is submitted, Bacalhau prints out the related `job_id`. We store that in an environment variable so that we can reuse it later on.


```python
%env JOB_ID={job_id}
```

### Declarative job description

The same job can be presented in the [declarative](../../../setting-up/jobs/job-specification/job.md) format. In this case, the description will look like this:

```yaml
name: DuckDB Hello World
type: batch
count: 1
tasks:
  - name: My main task
    Engine:
      type: docker
      params:
        Image: davidgasquez/datadex:v0.2.0
        Entrypoint:
          - /bin/bash
        Parameters:
          - -c
          - duckdb -s "select 1"
```

The job description should be saved in `.yaml` format, e.g. `duckdb1.yaml`, and then run with the command:
```bash
bacalhau job run duckdb1.yaml
```

## Checking the State of your Jobs

**Job status**: You can check the status of the job using `bacalhau list`.


```bash
%%bash
bacalhau list --id-filter ${JOB_ID}
```

When it says `Published` or `Completed`, that means the job is done, and we can get the results.

**Job information**: You can find out more information about your job by using `bacalhau describe`.


```bash
%%bash
bacalhau describe ${JOB_ID}
```

 **Job download**: You can download your job results directly by using `bacalhau get`. Alternatively, you can choose to create a directory to store your results. In the command below, we created a directory (`results`) and downloaded our job output to be stored in that directory.


```bash
%%bash
rm -rf results && mkdir -p results
bacalhau get $JOB_ID --output-dir results
```

## Viewing your Job Output

Each job result contains an `outputs` subfolder and `exitCode`, `stderr` and `stdout` files with relevant content. To view the file, run the following command:

```bash
%%bash
cat results/stdout  # displays the contents of the file

Expected Output:
┌───┐
│ 1 │
├───┤
│ 1 │
└───┘
```

## Running Arbitrary SQL commands

Below is the `bacalhau docker run` command to to run arbitrary SQL commands over the yellow taxi trips dataset


```bash
%%bash --out job_id
bacalhau docker run \
    -i ipfs://bafybeiejgmdpwlfgo3dzfxfv3cn55qgnxmghyv7vcarqe3onmtzczohwaq \
    --workdir /inputs \
    --id-only \
    --wait \
    davidgasquez/duckdb:latest \
    -- duckdb -s "select count(*) from '0_yellow_taxi_trips.parquet'"
```

### Structure of the command

Let's look closely at the command above:

`bacalhau docker run`: call to Bacalhau

`-i ipfs://bafybeiejgmdpwlfgo3dzfxfv3cn55qgnxmghyv7vcarqe3onmtzczohwaq \`: CIDs to use on the job. Mounts them at '/inputs' in the execution.

`davidgasquez/duckdb:latest`: the name and the tag of the docker image we are using

`/inputs`: path to input dataset

`duckdb -s`: execute DuckDB


When a job is submitted, Bacalhau prints out the related `job_id`. We store that in an environment variable so that we can reuse it later on.

### Declarative job description

The same job can be presented in the [declarative](../../../setting-up/jobs/job-specification/job.md) format. In this case, the description will look like this:

```yaml
name: DuckDB Parquet Query
type: batch
count: 1
tasks:
  - name: My main task
    Engine:
      type: docker
      params:
        WorkingDirectory: "/inputs"
        Image: davidgasquez/duckdb:latest
        Entrypoint:
          - /bin/bash
        Parameters:
          - -c
          - duckdb -s "select count(*) from '0_yellow_taxi_trips.parquet'"
    InputSources:
    - Target: "/inputs"
      Source:
        Type: "s3"
        Params:
          Bucket: "bacalhau-duckdb"
          Key: "*"
          Region: "us-east-1"
```

The job description should be saved in `.yaml` format, e.g. `duckdb2.yaml`, and then run with the command:
```bash
bacalhau job run duckdb2.yaml
```


**Job status**: You can check the status of the job using `bacalhau list`:


```bash
%%bash
bacalhau list --id-filter ${JOB_ID} --wide
```

**Job information**: You can find out more information about your job by using `bacalhau describe`.



```bash
%%bash
bacalhau describe ${JOB_ID}
```

**Job download**: You can download your job results directly by using `bacalhau get`. Alternatively, you can choose to create a directory to store your results. In the command below, we created a directory (`results`) and downloaded our job output to be stored in that directory.


```bash
%%bash
rm -rf results && mkdir -p results
bacalhau get ${JOB_ID} --output-dir results
```

## Viewing your Job Output

Each job result contains an `outputs` subfolder and `exitCode`, `stderr` and `stdout` files with relevant content. To view the file, run the following command:


```bash
%%bash
cat results/stdout

Expected Output:
┌──────────────┐
│ count_star() │
│    int64     │
├──────────────┤
│     24648499 │
└──────────────┘
```

## Support
If you have questions or need support or guidance, please reach out to the [Bacalhau team via Slack](https://bacalhauproject.slack.com/ssb/redirect) (**#general** channel).