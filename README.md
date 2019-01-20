# drone-ignore-config

A Drone Configuration Provider plugin that allows build events to ignore chosen file patterns.

Patterns are defined at a `.droneignore` file, similar to `.gitignore` and `.dockerignore`.

## Installation

At the moment it supports only GitHub.

1. Generate a GitHub access token with `repo` permission. This token is used to fetch the `.droneignore` configuration file.

2. Generate a shared secret key. This key is used to secure communication between the server and agents. The secret should be 32 bytes.

        $ openssl rand -hex 16
        558f3eacbfd5928157cbfe34823ab921

3. Run the container somewhere where the drone server can reach it:

        docker run \
          -p ${PLUGIN_PORT}:3000 \
          -e PLUGIN_SECRET=558f3eacbfd5928157cbfe34823ab921 \
          -e GITHUB_TOKEN=GITHUB8168c98304b \
          fbcbarbosa/drone-ignore

4. Update your drone server with information about the plugin:

          -e DRONE_YAML_ENDPOINT=http://${PLUGIN_HOST}:${PLUGIN_PORT}
          -e DRONE_YAML_SECRET=558f3eacbfd5928157cbfe34823ab921


See [the official docs](https://docs.drone.io/extend/config) for extra information on installing a Configuration Provider Plugin.

## Configuration

Add a `.droneignore` file to the root of your repository and add patterns of files to be ignored. For example:

```
*.md    # ignores changs to .md files
docs/** # ignores changes to all files within 'docs' folder
```

This plugin uses [sabhiram/go-gitignore](github.com/sabhiram/go-gitignore), which should behave just like [gitignore](https://git-scm.com/docs/gitignore). See [gitignore docs](https://git-scm.com/docs/gitignore) for more information on usage.

If _all_ files in the commit diff match, then drone-ignore returns the error below:

```
errors.New(".droneignore match, skipping build")
```

This causes the build to be skipped, while also unfortunately logging an error on the Drone server. The alternative is to send an empty dronefile back, however this would still enqueue the build and cause the default clone step to run. With this strategy, the build is skipped entirely, as if [ci skip] was present on the commit message.