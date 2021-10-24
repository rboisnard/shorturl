# shorturl

Welcome here

In this repository we explore the theme of an URL shortener application.

**The rules are:**
* ~~use redis~~ - *not implemented yet*
* write the application in go - *ongoing*
* ~~serve multiple workers with NGINX~~ - *not implemented yet*

**Technical choices:**
* The available Jenkins CI on a private network provides x86 and arm nodes but has a dependency on *docker* for containers, my development environement has *podman* available.
* My private servers are not powerful enough to run a Kubernetes cluster, if you are not satisfied, you are invited to help fund better hardware ;)