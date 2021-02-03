### Development environment
- Enable experimental on both docker daemon and client
```sh
apt update && apt install -y wget curl docker.io screen
echo $'{\n    "experimental": true\n}' | tee /etc/docker/daemon.json;
mkdir -p ~/.docker/cli-plugins/
echo '{ "experimental": "enabled" }' > ~/.docker/config.json
export DOCKER_CLI_EXPERIMENTAL=enabled
systemctl restart docker
docker version
```
- Download buildx
```sh
wget https://github.com/docker/buildx/releases/download/v0.2.0/buildx-v0.2.0.linux-amd64
mv buildx-v0.2.0.linux-amd64 ~/.docker/cli-plugins/docker-buildx
chmod +x ~/.docker/cli-plugins/docker-buildx
export DOCKER_BUILDKIT=1
docker buildx ls
```
- Download Golang
```sh
apt-get install -y mercurial git-core wget make build-essential rubygems ruby-dev
wget https://storage.googleapis.com/golang/go1.13.8.linux-amd64.tar.gz
tar -C /usr/local -xzf go*-*.tar.gz
```

### Replace
- Change all import library url
- Change gitlab-runner to cicd-runner
- Change company info
- Change logrus message
- Change ci/release_docker_image
- Change dockerfile for both ubuntu&alpine
- Change deb packaging
- Change gitlab-runner-helper to cicd-runner-helper

### Compile
```sh
export PATH=$PATH:/usr/local/go/bin:/root/go/bin/
make runner-and-helper-docker-host
```

### Release
- Push docker images
```sh
docker images|grep YOUR_DOCKERHUB_ID|head -n 5
docker push YOUR_DOCKERHUB_ID/cicd-runner:alpine-latest
docker push YOUR_DOCKERHUB_ID/cicd-runner-helper:x86_64-a1b2c3d4
```

### Run
- Docker executor
  - Register runner
```sh
docker run --rm -it -v $HOME/:/etc/cicd-runner debu99/cicd-runner:alpine-latest register
```
  - Verify privileged setting
  ```sh
  cat $HOME/config.toml
  ...
  privileged = true
  ...
  ```
  - Start runner
  ```sh
docker run -v $HOME/:/etc/cicd-runner -d --name cicd-runner -e DEBUG=true -e LOG_LEVEL=debug -v /var/run/docker.sock:/var/run/docker.sock debu99/cicd-runner:alpine-latest
	```
- Kubernetes executor
  - Generate token in base64
  ```sh
  echo -n 'YOUR_TOKEN_CODE' | base64
  ```
  - Update token in k8s/all.yaml
  - Apply yaml in K8s cicd namespace
  ```sh
  kubectl create ns cicd
  kubectl apply -f k8s/all.yaml -n cicd
  ```