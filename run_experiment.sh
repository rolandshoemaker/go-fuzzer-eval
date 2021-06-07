#!/bin/sh
set -ex

function run {
    gcloud compute instances create "exp-$1-$2" \
        --image-family=debian-10 \
        --image-project=debian-cloud \
        --machine-type=n1-standard-16 \
        --scopes cloud-platform,compute-rw \
        --metadata-from-file startup-script=startup.sh \
        --metadata runner_location="https://github.com/rolandshoemaker/go-fuzzer-eval",experiment_location="gs://go-fuzz-eval/$1/experiment.yaml",result_location="gs://go-fuzz-eval/$1/$2.log",checkout="$3" \
        --zone us-central1-a

    echo "started experiment for $3, result will be available at gs://go-fuzz-eval/$1/$2"
}

if [ $# -eq 0 ]; then
    echo "No arguments provided"
    exit 1
fi

exp_name=$(openssl rand -hex 4)

gsutil cp $1 "gs://go-fuzz-eval/$exp_name/experiment.yaml"

run "$exp_name" "old" "$2"

if [[ ! -z "$3" ]]; then
    run "$exp_name" "new" "$3"
fi
