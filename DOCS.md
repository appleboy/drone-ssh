---
date: 2017-01-29T00:00:00+00:00
title: SSH
author: appleboy
tags: [ publish, ssh ]
repo: appleboy/drone-ssh
logo: term.svg
image: appleboy/drone-ssh
---

Use the SSH plugin to execute commands on a remote server. The below pipeline configuration demonstrates simple usage:

```yaml
pipeline:
  ssh:
    image: appleboy/drone-ssh
    host: foo.com
    user: root
    password: 1234
    port: 22
    script:
      - echo hello
      - echo world
```

Example configuration in your `.drone.yml` file for multiple hosts:

```diff
pipeline:
  ssh:
    image: appleboy/drone-ssh
    host:
+    - foo.com
+    - bar.com
    user: root
    port: 22
    script:
      - echo hello
      - echo world
```

Example configuration for success build:

```diff
pipeline:
  ssh:
    image: appleboy/drone-ssh
    host: foo.com
    user: root
    password: 1234
    port: 22
    script:
      - echo hello
      - echo world
+   when:
+     status: success
```

Example configuration for tag event:

```diff
pipeline:
  ssh:
    image: appleboy/drone-ssh
    host: foo.com
    user: root
    password: 1234
    port: 22
    script:
      - echo hello
      - echo world
+   when:
+     status: success
+     event: tag
```

# Parameter Reference

host
: target hostname or IP

port
: ssh port of target host

user
: account for target host user

password
: password for target host user

key
: plain text of user private key

script
: execute commands on a remote server

timeout
: Timeout is the maximum amount of time for the TCP connection to establish.
