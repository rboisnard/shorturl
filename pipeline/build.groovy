def call() {
  stage("build") {
    /*
     * Executor nodes run inside containers using docker.
     * They share access to the docker daemon on their host
     * so we can use docker inside containers, but we can't
     * use podman. Sharing the images is done through a
     * local registry.
     */
    sh """
      docker build -t shorturl/redis:${BRANCH_NAME} -f src/redis/Dockerfile src/redis/
      docker build -t shorturl/worker:${BRANCH_NAME} -f src/worker/Dockerfile src/worker/
      #docker build -t shorturl/server:${BRANCH_NAME} -f src/server/Dockerfile src/server/
    """
  }
}

return this
