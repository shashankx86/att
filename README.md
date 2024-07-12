# Arcade Time Tracker (aka att)

```bash
      _    _   
 ___ | |_ | |_ 
| .'||  _||  _|
|__,||_|  |_|                 
```

Arcade Time Tracker (att) is a command-line tool designed to help you manage and track your work sessions efficiently.

## Table of Contents

1. [Installation](#installation)
2. [Usage](#usage)
    - [General Usage](#general-usage)
    - [Commands](#commands)
        - [configure](#configure)
            - [api-token](#api-token)
            - [slack-id](#slack-id)
        - [session](#session)
            - [list](#list)
            - [stats](#stats)
            - [goals](#goals)
            - [history](#history)
            - [start](#start)
            - [pause](#pause)
            - [cancel](#cancel)
        - [ping](#ping)
        - [status](#status)

## Installation

#### Linux/OSX
```
curl -s https://github.com/shashankx86/att/raw/main/install.sh | bash
```

### Manual Build & Install

To manually build and install the CLI tool, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/shashankx86/att.git
   ```

2. Navigate to the project directory:
   ```bash
   cd att
   ```

3. Build the project:
   ```bash
   go build -o att
   ```

4. Move the binary to a directory in your PATH:
   ```bash
   sudo mv att /usr/local/bin/
   ```

## Usage

### General Usage

```bash
att [command] [subcommand] [arguments]
```

### Commands

#### `configure`

The `configure` command is used to set up the CLI tool with necessary configurations such as the API token and Slack ID.

##### `api-token`

Sets the API token required for authentication.

**Usage:**

```bash
att configure api-token [token]
```

**Example:**

```bash
att configure api-token your-api-token
```

##### `slack-id`

Sets the Slack ID for the user.

**Usage:**

```bash
att configure slack-id [id]
```

**Example:**

```bash
att configure slack-id your-slack-id
```

#### `session`

The `session` command group is used to manage work sessions.

##### `list`

Lists the latest session.

**Usage:**

```bash
att session list
```

##### `stats`

Fetches and prints the user's stats.

**Usage:**

```bash
att session stats
```

##### `goals`

Fetches and prints the user's goals.

**Usage:**

```bash
att session goals
```

##### `history`

Fetches and prints the user's session history.

**Usage:**

```bash
att session history
```

##### `start`

Starts a new session. Provide a description of the work to be done.

**Usage:**

```bash
att session start [work description]
```

**Example:**

```bash
att session start "make robot"
```

If no work description is provided, you will be prompted to enter one.

##### `pause`

Pauses or resumes the current session.

**Usage:**

```bash
att session pause
```

##### `cancel`

Cancels the current session.

**Usage:**

```bash
att session cancel
```

#### `ping`

Pings the server to check connectivity.

**Usage:**

```bash
att ping
```

#### `status`

Fetches and prints the status of hack hour.

**Usage:**

```bash
att status
```

---