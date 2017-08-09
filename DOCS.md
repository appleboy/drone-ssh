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
    username: root
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
    username: root
    password: 1234
    port: 22
    script:
      - echo hello
      - echo world
```

Example configuration for login with user private key:

```diff
pipeline:
  ssh:
    image: appleboy/drone-ssh
    host: foo.com
    username: root
-   password: 1234
+   key: ${DEPLOY_KEY}
    port: 22
    script:
      - echo hello
      - echo world
```

Example configuration for login with file path of user private key:

```diff
pipeline:
  ssh:
    image: appleboy/drone-ssh
    host: foo.com
    username: root
-   password: 1234
+   key_path: ./deploy/key.pem
    port: 22
    script:
      - echo hello
      - echo world
```

Example configuration for command timeout (unit: second), default value is 60 seconds:

```diff
pipeline:
  ssh:
    image: appleboy/drone-ssh
    host: foo.com
    username: root
    password: 1234
    port: 22
+   command_timeout: 10
    script:
      - echo hello
      - echo world
```

Example configuration for execute commands on a remote server using ｀SSHProxyCommand｀:

```diff
pipeline:
  ssh:
    image: appleboy/drone-ssh
    host: foo.com
    username: root
    port: 22
    key: ${DEPLOY_KEY}
    script:
      - echo hello
      - echo world
+   proxy_host: 10.130.33.145
+   proxy_user: ubuntu
+   proxy_port: 22
+   proxy_key: ${PROXY_KEY}
```

Example configuration for success build:

```diff
pipeline:
  ssh:
    image: appleboy/drone-ssh
    host: foo.com
    username: root
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
    username: root
    password: 1234
    port: 22
    script:
      - echo hello
      - echo world
+   when:
+     status: success
+     event: tag
```

Example configuration for using custom secrets:

```diff
pipeline:
  ssh:
    image: appleboy/drone-ssh
    host: foo.com
    username: root
    password: 1234
    port: 22
+    secrets: [ aws_access_key_id ]
+    envs: [ aws_access_key_id ]
    script:
      - export AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID
```

This allows us to use custom secrets (The one's that the plugin doesn't expect) which can be used in the script section.

# Parameter Reference

host
: target hostname or IP

port
: ssh port of target host

username
: account for target host user

password
: password for target host user

key
: plain text of user private key

key_path
: key path of user private key

envs
: custom secrets which are made available in the script section

script
: execute commands on a remote server

timeout
: Timeout is the maximum amount of time for the TCP connection to establish.

command_timeout
: Command timeout is the maximum amount of time for the execute commands, default is 60 secs.

proxy_host
: proxy hostname or IP

proxy_port
: ssh port of proxy host

proxy_username
: account for proxy host user

proxy_password
: password for proxy host user

proxy_key
: plain text of proxy private key

proxy_key_path
: key path of proxy private key
