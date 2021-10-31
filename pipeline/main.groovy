def call() {
  def build = load "pipeline/build.groovy"
  def test = load "pipeline/test.groovy"
  def publish = load "pipeline/publish.groovy"
  def manifest = load "pipeline/manifest.groovy"

  def actions = [:]
  stash name: "repo"

  archs = [
    "x86": ["node": "x86_slave", "arch": "amd64"]//,
    /*
     * Weird issue with Jenkins remoting, shell steps end up
     * randomly in java.lang.InterruptedException.
     * Possibly because of durable-task plugin.
     *
     * java.lang.InterruptedException
     * at java.base/java.lang.Object.wait(Native Method)
     * at hudson.remoting.Request.call(Request.java:177)
     * at hudson.remoting.Channel.call(Channel.java:1000)
     * at hudson.Launcher$RemoteLauncher.launch(Launcher.java:1122)
     * at hudson.Launcher$ProcStarter.start(Launcher.java:507)
     * at org.jenkinsci.plugins.durabletask.BourneShellScript.launchWithCookie(BourneShellScript.java:176)
     * ...
     *
     * Temporarily deactivating arm build.
     */
    //"arm": ["node": "arm_slave", "arch": "arm64"]
  ]

  archs.each{ arch, config ->
    actions[arch] = {
      node(config.node) {
        unstash "repo"
        build()
        test()
        publish(config.arch)
      }
    }
  }

  parallel(actions)

  node("x86_slave") {
    // TODO: different tags depending on pull-request, normal branch, release branch ...
    // same as in publish()
    manifest("staging")
  }
}

return this
