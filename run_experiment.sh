#!/bin/sh
set -e

function run {
    gcloud compute instances create $1 \
        --image-family=debian-10 \
        --image-project=debian-cloud \
        --machine-type=n1-standard-1 \
        --scopes cloud-platform,compute-rw \
        --metadata-from-file startup-script=startup.sh \
        --metadata experiment_location="gs://go-fuzz-bench/$1/experiment.yaml" \
        --metadata result_location="gs://go-fuzz-bench/$1/$2" \
        --metadata checkout="$3" \
        --zone us-central1 

    echo "started experiment for $3, result will be available at gs://go-fuzz-bench/$1/$2"
}

if [ $# -eq 0 ]; then
    echo "No arguments provided"
    exit 1
fi

exp_name=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 8 | head -n 1)

run "$exp_name" "old.log" "$1"

if [[ ! -z "$2" ]]; then
    run "$exp_name" "new.log" "$2"
fi