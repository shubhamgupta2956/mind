# mind

## Building project

1. Copy the `sample.info` file as `info` and add Slack and Todoist client IDs and secrets.
    ```bash
    $ cp sample.info info
    $ vim info # add client ids and secrets
    ```

    Scopes required for Slack are:
    1. `channels:read`
    1. `channels:history`
    1. `chat:write:user`
    1. `users:read`

    Scopes required for Todoist are:
    1. `data:read_write`

1. Run `make` or `make build` to build the project into `target` directory.
    ```bash
    $ make # builds the project as ./target/mind
    ```

1. To use the CLI as `mind ...`, move the binary somewhere in your `$PATH`.
    ```bash
    $ sudo mv target/mind /usr/local/bin
    ```

## Using CLI

1. To configure ouput type, use `configure` command.
    ```bash
    $ mind configure plain # for plain text output
    $ mind configure json # for json text output
    ```

1. To authorize with Slack or Todoist use their `auth` command.
    ```bash
    $ mind slack auth # authorize slack
    $ mind todoist auth # authorize todoist
    ```
    You'll be prompted with a URL. It's the authorization URL. Visit the URL and authorize Mind for your workspace/account.
    Also the redirectURL is hardcoded as `http://127.0.0.1:12345`.

1. To send messages on Slack, use `send` command.
    ```bash
    $ mind slack send general "What's the progress on assignment?" # sends the message to general channel on Slack
    ```

1. To fetch messages from Slack channel, use `unreads` command.
    ```bash
    $ mind slack unreads general # fetches latest messages from general channel
    $ mind slack unreads --limit 5 general # --limit or -l fetches latest 'n' messages only
    ```
    ***Note:*** Slack API does not provide to fetch unread messages by the user, hence this just fetches "n" latest messages from Slack channel.

1. To add a new task on todoist, use `add` command.
    ```bash
    $ mind todoist add "Complete assignment" # adds the task with due-date today
    ```

1. To fetch today's tasks, use `today` command.
    ```bash
    $ mind todoist today
    ```

1. To print the summary of commands used on any date, use `summary` command.
    ```bash
    $ mind session summary 2020-05-11 # gets summary for May 11, 2020
    ```

1. Use `-h` or `--help` with any command or `help [command]` to print help.
    ```bash
    $ mind slack unreads --help # prints help
    ```
