Use the SSH plugin to execute commands on a remote server. You will need to
supply Drone with a private SSH key to being able to connect to a host.

## Config

The following parameters are used to configure the plugin:

* **host** - address or IP of the remote machine
* **port** - port to connect to on the remote machine
* **user** - user to log in as on the remote machine
* **key** - private SSH key for the remote machine
* **sleep** - sleep for seconds between host connections
* **timeout** - timeout for the tcp connection attempt
* **commands** - list of commands to execute

The following secret values can be set to configure the plugin.

* **SSH_HOST** - corresponds to **host**
* **SSH_PORT** - corresponds to **port**
* **SSH_USER** - corresponds to **user**
* **SSH_KEY** - corresponds to **key**
* **SSH_SLEEP** - corresponds to **sleep**
* **SSH_TIMEOUT** - corresponds to **timeout**

It is highly recommended to put the **SSH_KEY** into a secret so it is not
exposed to users. This can be done using the drone-cli.

```bash
drone secret add --image=ssh \
    octocat/hello-world SSH_KEY @path/to/.ssh/id_rsa
```

Then sign the YAML file after all secrets are added.

```bash
drone sign octocat/hello-world
```

See [secrets](http://readme.drone.io/0.5/usage/secrets/) for additional
information on secrets

## Examples

Example configuration in your .drone.yml file for a single host:

```yaml
pipeline:
  ssh:
    host: foo.com
    user: root
    port: 22
    commands:
      - echo hello
      - echo world
```

Example configuration in your .drone.yml file for multiple hosts:

```yaml
pipeline:
  ssh:
    host:
     - foo.com
     - bar.com
    user: root
    port: 22
    sleep: 5
    commands:
      - echo hello
      - echo world
```

In the above example Drone executes the commands on multiple hosts
sequentially. If the commands fail on a single host this plugin exits
immediatly, and will not run your commands on the remaining hosts in the
list.

The above example also uses the `sleep` parameter. The sleep parameter
instructs Drone to sleep for N seconds between host executions.
