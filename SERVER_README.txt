The server modes are listed below.
The order of arguments is irrelevant.
Example invocations are given in the end.

Show help message and exit:
    java -jar server.jar [-h]
Omitting the -h argument shows an abbreviated description.
Providing the -h argument shows this detailed description.

Run a client on a level or a directory of levels, optionally output to GUI and/or log file:
    java -jar server.jar -c <client-cmd> -l <level-file-or-dir-path> [-t <seconds>]
                         [-g [<screen>] [-s <ms-per-action>] [-p] [-f] [-i]]
                         [-o <log-file-path>]
Where the arguments are as follows:
    -c <client-cmd>
        Specifies the command the server will use to start the client process, including all client arguments.
        The <client-cmd> string will be naïvely tokenized by splitting on whitespace, and
        then passed to your OS native process API through Java's ProcessBuilder.
    -l <level-file-or-dir-path>
        Specifies the path to either a single level file, or to a directory containing one or more level files.
    -t <seconds>
        Optional. Specifies a timeout in Seconds for the client.
        The server will terminate the client after this timeout.
        If this argument is not given, the server will never time the client out.
    -g [<screen>]
        Optional. Enables GUI output.
        The optional <screen> argument specifies which screen to start the GUI on.
        See notes on <screen> below for an explanation of the valid values.
    -s <ms-per-action>
        Optional. Set the GUI playback speed in milliseconds per action.
        By default the speed is 250 ms per action.
    -p
        Optional. Start the GUI paused.
        By default the GUI starts playing immediately.
    -f
        Optional. Start the GUI in fullscreen mode.
        By default the GUI starts in windowed mode.
    -i
        Optional. Start the GUI with interface hidden.
        By default the GUI shows interface elements for navigating playback.
    -o <log-file-path>
        Optional. Writes log file(s) of the client run.
        If the -l argument is a level file path, then a single log file is written to the given log file path.
        If the -l argument is a level directory path, then logs for the client run on all levels in the
        level directory are compressed as a zip file written to the given log file path.
        NB: The log file may *not* already exist (the server does not allow overwriting files).

Replay one or more log files, optionally output to synchronized GUIs:
    java -jar server.jar -r <log-file-path> [<log-file-path> ...]
                         [-g [<screen> ...] [-s <ms-per-action>] [-p] [-f] [-i]]
Where the arguments are as follows:
    -r <log-file-path> [<log-file-path> ...]
        Specifies one or more log files to replay.
    -g [<screen> ...]
        Optional. Enables GUI output. The playback of the replays are synchronized.
        The optional <screen> arguments specify which screen to start the GUI on for each log file.
        See notes on <screen> below for an explanation of the valid values.
    -s <ms-per-action>
        Optional. Set the GUI playback speed in milliseconds per action.
        By default the speed is 250 ms per action.
    -p
        Optional. Start the GUI paused.
        By default the GUI starts playing immediately.
    -f
        Optional. Start the GUI in fullscreen mode.
        By default the GUI starts in windowed mode.
    -i
        Optional. Start the GUI with interface hidden.
        By default the GUI shows interface elements for navigating playback.

Notes on the <screen> arguments:
    Values for the <screen> arguments are integers in the range 0..(<num-screens> - 1).
    The server attemps to enumerate screens from left-to-right, breaking ties with top-to-bottom.
    The real underlying screen ordering is system-defined, and the server may fail at enumerating in the above order.
    If no <screen> argument is given, then the 'default' screen is used, which is a system-defined screen.

    E.g. in a grid aligned 2x2 screen setup, the server will attempt to enumerate the screens as:
    0: Top-left screen.
    1: Bottom-left screen.
    2: Top-right screen.
    3: Bottom-right screen.

    E.g. in a 1x3 horizontally aligned screen setup, the server will attempt to enumerate the screens as:
    0: The left-most screen.
    1: The middle screen.
    2: The right-most screen.

Supported domains (case-sensitive):
    hospital

Client example invocations:
    # Client on single level, no output.
    java -jar server.jar -c "java ExampleClient" -l "levels/example.lvl"

    # Client on single level, output to GUI on default screen.
    java -jar server.jar -c "java ExampleClient" -l "levels/example.lvl" -g

    # Client on single level, output to GUI on screen 0.
    java -jar server.jar -c "java ExampleClient" -l "levels/example.lvl" -g 0

    # Client on single level, output to log file.
    java -jar server.jar -c "java ExampleClient" -l "levels/example.lvl" -o "logs/example.log"

    # Client on single level, output to GUI on default screen and to log file.
    java -jar server.jar -c "java ExampleClient" -l "levels/example.lvl" -g -o "logs/example.log"

    # Client on a directory of levels, no output.
    java -jar server.jar -c "java ExampleClient" -l "levels"

    # Client on a directory of levels, output to log archive.
    java -jar server.jar -c "java ExampleClient" -l "levels" -o "logs.zip"

Replay example invocations:
    # Replay a single log file, no output.
    java -jar server.jar -r "logs/example.log"

    # Replay a single log file, output to GUI on default screen.
    java -jar server.jar -r "logs/example.log" -g

    # Replay two log files, output to synchronized GUIs on screen 0 and 1.
    # Start the GUIs paused, in fullscreen mode and with hidden interface elements to avoid spoilers.
    # Play back actions at a speed of one action every 500 milliseconds.
    java -jar server.jar -r "logs/example1.log" "logs/example2.log" -g 0 1 -p -f -i -s 500
