#!/bin/sh
set -ex

function cleanup {
    # Delete the instance
    NAME=$(curl -s http://metadata.google.internal/computeMetadata/v1/instance/name -H 'Metadata-Flavor: Google')
    ZONE=$(curl -s http://metadata.google.internal/computeMetadata/v1/instance/zone -H 'Metadata-Flavor: Google')
    gcloud --quiet compute instances delete $NAME --zone=$ZONE
}
trap cleanup EXIT

curl -s "https://storage.googleapis.com/signals-agents/logging/google-fluentd-install.sh" | bash
service google-fluentd restart &

EXPERIMENT_LOC=$(curl -s http://metadata.google.internal/computeMetadata/v1/instance/attributes/experiment_location -H "Metadata-Flavor: Google")
RESULT_LOC=$(curl -s http://metadata.google.internal/computeMetadata/v1/instance/attributes/result_location -H "Metadata-Flavor: Google")
CHECKOUT=$(curl -s http://metadata.google.internal/computeMetadata/v1/instance/attributes/checkout -H "Metadata-Flavor: Google")
RUNNER_LOCATION=$(curl -s http://metadata.google.internal/computeMetadata/v1/instance/attributes/runner_location -H "Metadata-Flavor: Google")

# Download the experiment description from the GCS bucket
gsutil cp "$EXPERIMENT_LOC" experiment.yaml

# Clone the runner and cehckout the commit we are interested in
git clone https://go.googlesource.com/go
cd go

# Checkout either the commit or CL
git fetch https://go.googlesource.com/go $CHECKOUT && git checkout FETCH_HEAD

# Build go
./src/make.bash
cd -

# Clone the runner
git clone "$RUNNER_LOCATION" runner

# Run the experiment
WORKDIR="$(pwd)"
go/bin/go run ./runner -experiment experiment.yaml -result result.log -go $WORKDIR/go/bin/go

# Copy results to the GCS bucket
gsutil cp result.log "$RESULT_LOC"