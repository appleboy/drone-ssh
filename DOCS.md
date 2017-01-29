Use the SSH plugin to execute commands on a remote server. You will need to
supply Drone with a private SSH key to being able to connect to a host.

## Config

The following parameters are used to configure the plugin:

* **host** - address or IP of the remote machine
* **port** - port to connect to on the remote machine
* **user** - user to log in as on the remote machine
* **passsword** - password to log in as on the remote machine
* **key** - private SSH key for the remote machine
* **timeout** - timeout for the tcp connection attempt
* **script** - list of commands to execute

The following secret values can be set to configure the plugin.

* **SSH_HOST** - corresponds to **host**
* **SSH_PORT** - corresponds to **port**
* **SSH_USER** - corresponds to **user**
* **SSH_PASSWORD** - corresponds to **password**
* **SSH_KEY** - corresponds to **key**
* **SSH_TIMEOUT** - corresponds to **timeout**

## Examples

Example configuration in your .drone.yml file for a single host:

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

Example configuration in your .drone.yml file for multiple hosts:

```yaml
pipeline:
  ssh:
    image: appleboy/drone-ssh
    host:
     - foo.com
     - bar.com
    user: root
    port: 22
    script:
      - echo hello
      - echo world
```

In the above example Drone executes the commands on multiple hosts
sequentially. If the commands fail on a single host this plugin exits
immediatly, and will not run your commands on the remaining hosts in the
list.
