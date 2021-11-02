def call(String arch) {
  stage("publish") {
    /*
     * Executor nodes run inside containers using docker.
     * They share access to the docker daemon on their host
     * so we can use docker inside containers, but we can't
     * use podman. Sharing the images is done through a
     * local registry.
     */
    // TODO: different tags depending on pull-request, normal branch, release branch ...
    String tag = "staging"
    sh """
      docker tag shorturl/redis:${BRANCH_NAME} raspi01:5000/shorturl/redis:${tag}-${arch}
      docker push raspi01:5000/shorturl/redis:${tag}-${arch}

      docker tag shorturl/worker:${BRANCH_NAME} raspi01:5000/shorturl/worker:${tag}-${arch}
      docker push raspi01:5000/shorturl/worker:${tag}-${arch}

      #docker tag shorturl/server:${BRANCH_NAME} raspi01:5000/shorturl/server:${tag}-${arch}
      #docker push raspi01:5000/shorturl/server:${tag}-${arch}
    """
  }
}

return this
