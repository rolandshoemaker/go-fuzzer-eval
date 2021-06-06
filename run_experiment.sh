#!/bin/sh
set -ex

function run {
    gcloud compute instances create $1 \
        --image-family=debian-10 \
        --image-project=debian-cloud \
        --machine-type=n1-standard-1 \
        --scopes cloud-platform,compute-rw \
        --metadata-from-file startup-script=startup.sh \
        --metadata runner_location="https://github.com/rolandshoemaker/go-fuzzer-eval" \
        --metadata experiment_location="gs://go-fuzz-eval/$1/experiment.yaml" \
        --metadata result_location="gs://go-fuzz-eval/$1/$2" \
        --metadata checkout="$3" \
        --zone us-central1 

    echo "started experiment for $3, result will be available at gs://go-fuzz-eval/$1/$2"
}

if [ $# -eq 0 ]; then
    echo "No arguments provided"
    exit 1
fi

exp_name=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 8 | head -n 1)

gsutil cp $1 "gs://go-fuzz-eval/$exp_name/experiment.yaml"

run "$exp_name" "old" "$2"

if [[ ! -z "$3" ]]; then
    run "$exp_name" "new" "$3"
fi