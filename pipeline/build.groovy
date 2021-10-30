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
      docker build -t raspi01:5000/shorturl/redis:staging -f src/redis/Dockerfile src/redis/
      docker build -t raspi01:5000/shorturl/worker:staging -f src/worker/Dockerfile src/worker/
      #docker build -t raspi01:5000/shorturl/server:staging -f src/server/Dockerfile src/server/

      docker push raspi01:5000/shorturl/redis:staging
      docker push raspi01:5000/shorturl/worker:staging
      #docker push raspi01:5000/shorturl/server:staging
    """
  }
}

return this
