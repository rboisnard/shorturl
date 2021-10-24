def call() {
  stage("build") {
    sh """
      #docker build -t shorturl/redis:staging -f src/redis/Dockerfile src/redis/
      docker build -t shorturl/worker:staging -f src/worker/Dockerfile src/worker/
      #docker build -t shorturl/server:staging -f src/server/Dockerfile src/server/
    """
  }
}

return this
