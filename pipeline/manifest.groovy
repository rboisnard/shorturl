def call(String tag) {
  stage("manifest") {
    create_manifest("raspi01:5000/shorturl/redis:${tag}")
    create_manifest("raspi01:5000/shorturl/worker:${tag}")
    //create_manifest("raspi01:5000/shorturl/server:${tag}")
  }
}

def create_manifest(String image_tag) {
  // TODO: review build list
  // using existing arm64 variant while the build is deactivated
  sh """
    export DOCKER_CLI_EXPERIMENTAL=enabled

    docker pull ${image_tag}-amd64
    docker pull ${image_tag}-arm64

    docker manifest create --insecure \
      ${image_tag}                    \
      --amend ${image_tag}-amd64      \
      --amend ${image_tag}-arm64

    docker manifest push --insecure ${image_tag}
    docker manifest inspect --insecure ${image_tag}
  """
}

return this
