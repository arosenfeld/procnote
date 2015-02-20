# procnote

procnote is a simple program that associates notes with processes.  For example:

    $ procnote add 1234 "This is an important process!"
    
To list all the notes you've added:

    $ procnote list
    PID     STATUS    NOTE
    123     Running   A process calculating the meaning of life
    45678   Running   Don't stop this important process!
    9012    Stopped   Some other process
    
Searching notes is also easy:

    $ procnote search important
    PID     STATUS    NOTE
    45678   Running   Don't stop this important process!
    
To delete a single note use:

    $ procnote del 123
    
You can clear all your notes with:

    $ procnote clear
    
or just the stopped proceses with

    $ procnote clear --stopped
