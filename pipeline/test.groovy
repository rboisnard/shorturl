def call() {
  stage("test") {
    sh """
      docker run --rm -d -p 5500:5500 -e PORT=5500 --name=shorturl_worker shorturl/worker:staging
      #curl -v localhost:5500
      docker stop shorturl_worker
    """
  }
}

return this
