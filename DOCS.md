---
date: 2019-08-04T00:00:00+00:00
title: SSH
author: appleboy
tags: [ deploy, publish, ssh ]
repo: appleboy/drone-ssh
logo: term.svg
image: appleboy/drone-ssh
---

Use the SSH plugin to execute commands on a remote server. The below pipeline configuration demonstrates simple usage:

```yaml
- name: ssh commands
  image: appleboy/drone-ssh
  settings:
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
  - name: ssh commands
    image: appleboy/drone-ssh
    settings:
      host:
+       - foo.com
+       - bar.com
      username: root
      password: 1234
      port: 22
      script:
      - echo hello
      - echo world
```

Example configuration for command timeout, default value is 60 seconds:

```diff
  - name: ssh commands
    image: appleboy/drone-ssh
    settings:
      host: foo.com
      username: root
      password: 1234
      port: 22
+       command_timeout: 2m
      script:
        - echo hello
        - echo world
```

Example configuration for execute commands on a remote server using ｀SSHProxyCommand｀:

```diff
  - name: ssh commands
    image: appleboy/drone-ssh
    settings:
      host: foo.com
      username: root
      password: 1234
      port: 22
      script:
        - echo hello
        - echo world
+     proxy_host: 10.130.33.145
+     proxy_user: ubuntu
+     proxy_port: 22
+     proxy_password: 1234
```

Example configuration using password from secrets:

```diff
  - name: ssh commands
    image: appleboy/drone-ssh
    settings:
      host: foo.com
      username: root
+     password:
+       from_secret: ssh_password
      port: 22
      script:
        - echo hello
        - echo world
```

Example configuration using ssh key from secrets:

```diff
  - name: ssh commands
    image: appleboy/drone-ssh
    settings:
      host: foo.com
      username: root
      port: 22
+     key:
+       from_secret: ssh_key
      script:
        - echo hello
        - echo world
```

Example configuration for exporting custom secrets:

```diff
  - name: ssh commands
    image: appleboy/drone-ssh
    settings:
      host: foo.com
      username: root
      password: 1234
      port: 22
+     envs:
        - aws_access_key_id
      script:
        - export AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID
```

Example configuration for stoping script after first failure:

```diff
  - name: ssh commands
    image: appleboy/drone-ssh
    settings:
      host: foo.com
      username: root
      password: 1234
      port: 22
+     script_stop: true
      script:
        - mkdir abc/def/efg
        - echo "you can't see the steps."
```

Example configuration for passphrase which protecting a private key:

```diff
  - name: ssh commands
    image: appleboy/drone-ssh
    settings:
      host: foo.com
      username: root
+     key:
+       from_secret: ssh_key
+     passphrase: 1234
      port: 22
      script:
        - mkdir abc/def/efg
        - echo "you can't see the steps."
```

## Secret Reference

ssh_username
: account for target host user

ssh_password
: password for target host user

ssh_passphrase
: The purpose of the passphrase is usually to encrypt the private key.

ssh_key
: plain text of user private key

proxy_ssh_username
: account for user of proxy server

proxy_ssh_password
: password for user of proxy server

proxy_ssh_passphrase
: The purpose of the passphrase is usually to encrypt the private key.

proxy_ssh_key
: plain text of user private key for proxy server

## Parameter Reference

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

script_stop
: stop script after first failure

timeout
: Timeout is the maximum amount of time for the ssh connection to establish, default is 30 seconds.

command_timeout
: Command timeout is the maximum amount of time for the execute commands, default is 10 minutes.

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
