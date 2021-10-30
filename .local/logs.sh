#!/bin/sh

containers=$(podman pod inspect shorturl-pod -f {{.Containers}} | sed -e "s/\[//" -e "s/]//" -e ":again s/{//" -e "t again" -e ": again2 s/}/\n/" -e "t again2" | awk {'print $2'})
for container in ${containers}; do
  echo
  echo "### ${container} ###"
  podman logs ${container}
done