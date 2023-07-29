docker build -t mochi-broker:v0.0.26 .
docker tag mochi-broker:v0.0.26  us-central1-docker.pkg.dev/freezer-monitor-dev-e7d4c/mochi-broker/mochi-broker:v0.0.26
docker push us-central1-docker.pkg.dev/freezer-monitor-dev-e7d4c/mochi-broker/mochi-broker:v0.0.26
gcloud compute instances update-container instance-1 --zone us-central1-a --container-image=us-central1-docker.pkg.dev/freezer-monitor-dev-e7d4c/mochi-broker/mochi-broker:v0.0.26


// PROD
docker build -t mochi-broker:v0.0.26 .
docker tag mochi-broker:v0.0.26  us-central1-docker.pkg.dev/freezer-monitor-dev-e7d4c/mochi-broker/mochi-broker:v0.0.26
docker push us-central1-docker.pkg.dev/freezer-monitor-dev-e7d4c/mochi-broker/mochi-broker:v0.0.26
gcloud compute instances update-container instance-1 --zone us-central1-a --container-image=us-central1-docker.pkg.dev/freezer-monitor-dev-e7d4c/mochi-broker/mochi-broker:v0.0.26