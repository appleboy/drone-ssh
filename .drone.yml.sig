workspace:
  base: /srv/app
  path: src/github.com/appleboy/drone-scp

pipeline:
  clone:
    image: plugins/git
    tags: true

  # restore the cache from an sftp server
  restore_cache:
    image: appleboy/drone-sftp-cache
    restore: true
    mount: [ .glide, vendor ]
    ignore_branch: true

  test:
    image: appleboy/golang-testing
    pull: true
    environment:
      TAGS: netgo
      GOPATH: /srv/app
    commands:
      - adduser -h /home/drone-scp -s /bin/bash -D -S drone-scp
      - passwd -d drone-scp
      - mkdir -p /home/drone-scp/.ssh
      - chmod 700 /home/drone-scp/.ssh
      - cp tests/.ssh/id_rsa.pub /home/drone-scp/.ssh/authorized_keys
      - chown -R drone-scp /home/drone-scp/.ssh
      # install ssh and start server
      - apk update && apk add openssh openrc
      - rm -rf /etc/ssh/ssh_host_rsa_key /etc/ssh/ssh_host_dsa_key
      - ./tests/entrypoint.sh /usr/sbin/sshd -D &
      - make dep_install
      - make vet
      - make lint
      - make test
      - make coverage
      - make build
      # build binary for docker image
      - make static_build
    when:
      event: [ push, tag, pull_request ]

  release:
    image: appleboy/golang-testing
    pull: true
    environment:
      TAGS: netgo
      GOPATH: /srv/app
    commands:
      - make release
    when:
      event: [ tag ]
      branch: [ refs/tags/* ]

  docker:
    image: plugins/docker
    repo: ${DRONE_REPO}
    tags: [ '${DRONE_TAG}' ]
    when:
      event: [ tag ]
      branch: [ refs/tags/* ]

  docker:
    image: plugins/docker
    repo: ${DRONE_REPO}
    tags: [ 'latest' ]
    when:
      event: [ push ]
      branch: [ master ]

  github:
    image: plugins/github-release
    files:
      - dist/release/*
    when:
      event: [ tag ]
      branch: [ refs/tags/* ]

  # rebuild the cache on the sftp server
  rebuild_cache:
    image: appleboy/drone-sftp-cache
    rebuild: true
    mount: [ .glide, vendor ]
    ignore_branch: true
    when:
      branch: master
