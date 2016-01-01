# drone-ssh

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-ssh/status.svg)](http://beta.drone.io/drone-plugins/drone-ssh)
[![](https://badge.imagelayers.io/plugins/drone-ssh:latest.svg)](https://imagelayers.io/?images=plugins/drone-ssh:latest 'Get your own badge on imagelayers.io')

Drone plugin for executing commands on a remote host through SSH

## Usage

```
./drone-ssh <<EOF
{
    "repo" : {
        "owner": "foo",
        "name": "bar",
        "full_name": "foo/bar"
    },
    "build" : {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "commit": "9f2849d5",
        "branch": "master",
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
    },
    "vargs": {
        "host": "foo.com",
        "user": "root",
        "port": 22,
        "commands": [
            "echo hello",
            "echo world"
        ]
    }
}
EOF
```

## Docker

Build the Docker container using `make`:

```
make deps build docker
```

### Example

```sh
docker run -i plugins/drone-ssh <<EOF
{
    "repo" : {
        "owner": "foo",
        "name": "bar",
        "full_name": "foo/bar"
    },
    "build" : {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "commit": "9f2849d5",
        "branch": "master",
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
    },
    "vargs": {
        "host": "foo.com",
        "user": "root",
        "port": 22,
        "commands": [
            "echo hello",
            "echo world"
        ]
    }
}
EOF
```
