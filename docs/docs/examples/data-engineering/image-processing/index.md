---
sidebar_label: Simple Image Processing
sidebar_position: 4
description: "How to process images stored in IPFS with Bacalhau"
---
# Simple Image Processing

In this example tutorial, we will show you how to use Bacalhau to process images on a [Landsat dataset](https://ipfs.io/ipfs/QmeZRGhe4PmjctYVSVHuEiA9oSXnqmYa4kQubSHgWbjv72/). 

Bacalhau has the unique capability of operating at a massive scale in a distributed environment. This is made possible because data is naturally sharded across the IPFS network amongst many providers. We can take advantage of this to process images in parallel.

## Prerequisite

To get started, you need to install the Bacalhau client, see more information [here](../../../getting-started/installation.md)


```python
!command -v bacalhau >/dev/null 2>&1 || (export BACALHAU_INSTALL_DIR=.; curl -sL https://get.bacalhau.org/install.sh | bash)
path=!echo $PATH
%env PATH=./:{path[0]}
```

## Running a Bacalhau Job

To submit a workload to Bacalhau, we will use the `bacalhau docker run` command. This command allows to pass input data volume with a `-i ipfs://CID:path` argument just like Docker, except the left-hand side of the argument is a [content identifier (CID)](https://github.com/multiformats/cid). This results in Bacalhau mounting a *data volume* inside the container. By default, Bacalhau mounts the input volume at the path `/inputs` inside the container.

Bacalhau also mounts a data volume to store output data. The `bacalhau docker run` command creates an output data volume mounted at `/outputs`. This is a convenient location to store the results of your job. 


```bash
%%bash --out job_id
bacalhau docker run \
    --wait \
    --wait-timeout-secs 100 \
    --id-only \
    -i ipfs://QmeZRGhe4PmjctYVSVHuEiA9oSXnqmYa4kQubSHgWbjv72:/input_images \
    --entrypoint mogrify \
    dpokidov/imagemagick:7.1.0-47-ubuntu \
    -- -resize 100x100 -quality 100 -path /outputs '/input_images/*.jpg'
```

### Structure of the command

Let's look closely at the command above:

1. `bacalhau docker run`: call to Bacalhau
1. `-i ipfs://QmeZRGhe4PmjctYVSVHuEiA9oSXnqmYa4kQubSHgWbjv72:/input_images`: Specifies the input data, which is stored in IPFS at the given CID.
1. `--entrypoint mogrify`:  Overrides the default ENTRYPOINT of the image, indicating that the mogrify utility from the ImageMagick package will be used instead of the default entry.
1. `dpokidov/imagemagick:7.1.0-47-ubuntu`: The name and the tag of the docker image we are using
1. `-- -resize 100x100 -quality 100 -path /outputs '/input_images/*.jpg'`: These arguments are passed to mogrify and specify operations on the images: resizing to 100x100 pixels, setting quality to 100, and saving the results to the `/outputs` folder.

When a job is submitted, Bacalhau prints out the related `job_id`. We store that in an environment variable so that we can reuse it later on.


```python
%env JOB_ID={job_id}
```
### Declarative job description

The same job can be presented in the [declarative](../../../setting-up/jobs/job-specification/job.md) format. In this case, the description will look like this:

```yaml
name: Simple Image Processing
type: batch
count: 1
tasks:
  - name: My main task
    Engine:
      type: docker
      params:
        Image: dpokidov/imagemagick:7.1.0-47-ubuntu
        Entrypoint:
          - /bin/bash
        Parameters:
          - -c
          - magick mogrify -resize 100x100 -quality 100 -path /outputs '/input_images/*.jpg'
    Publisher:
      Type: ipfs
    ResultPaths:
      - Name: outputs
        Path: /outputs
    InputSources:
    - Target: "/input_images"
      Source:
        Type: "s3"
        Params:
          Bucket: "landsat-image-processing"
          Key: "*"
          Region: "us-east-1"
```

The job description should be saved in `.yaml` format, e.g. `image.yaml`, and then run with the command:
```bash
bacalhau job run image.yaml
```


## Checking the State of your Jobs

**Job status**: You can check the status of the job using `bacalhau list`. 


```bash
%%bash
bacalhau list --id-filter ${JOB_ID} --no-style

Expected Output:
CREATED   ID          JOB                                       STATE      PUBLISHED
11:19:52  4bb743a4    Type:"docker",Params:"map[Entrypoint:[mo  Completed
                       grify] EnvironmentVariables:[] Image:dpo
                       kidov/imagemagick:7.1.0-47-ubuntu Parame
                       ters:[-resize 100x100 -quality 100 -path
                        /outputs /input_images/*.jpg] WorkingDi
                       rectory:]"
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
rm -rf results && mkdir results # Temporary directory to store the results
bacalhau get ${JOB_ID} --output-dir results # Download the results
```

## Viewing your Job Output

To view the file, run the following command:


```bash
%%bash
ls -lah results/outputs
```

### Display the image

To view the images, we will use **glob** to return all file paths that match a specific pattern. 


```python
import glob
from IPython.display import Image, display
for imageName in glob.glob('results/outputs/*.jpg'):
    display(Image(filename=imageName))
```


    
![jpeg](index_files/index_21_0.jpg)
    



    
![jpeg](index_files/index_21_1.jpg)
    



    
![jpeg](index_files/index_21_2.jpg)
    



    
![jpeg](index_files/index_21_3.jpg)
    



    
![jpeg](index_files/index_21_4.jpg)
    



    
![jpeg](index_files/index_21_5.jpg)
    



    
![jpeg](index_files/index_21_6.jpg)
    



    
![jpeg](index_files/index_21_7.jpg)
    



    
![jpeg](index_files/index_21_8.jpg)
    


## Support
If you have questions or need support or guidance, please reach out to the [Bacalhau team via Slack](https://bacalhauproject.slack.com/ssb/redirect) (**#general** channel).
