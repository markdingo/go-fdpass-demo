# go-fdpass-demo

Demo go code showing how to pass a file descriptor between processes
using Unix Domain Sockets.

Since fd-passing is a fairly arcane part of Unix and perhaps an even
more arcane part of go, I figure that others might benefit from having
access to a known working example. While there are plenty of available
examples written in C, at the time of writing I could not readily find
a working example in go. Thus this is that.

### When might you use fd-passing?

There are a number of use-cases, but generally fd-passing is used when
you want to centrally manage and control access to underlying system
resources in a way that is not readily possible with the controls
offered by the operating system.

For example you might want to give clients access to sockets created
on privileged ports but only to some clients and only to some
ports. One way to do this is have a server establish the socket and
fd-pass it back to the client if it passes the access-control rules.

Another example might be if you want to give clients access to some
files in a directory but not others. Such as those under a certain
size or age. The client sends the open request to the server, the
server applies the age/size logic and fd-passes back the opened file
if it's is approved.

Another use-case is to create a server as a container of idle network
connections. If your main server uses a lot of state per connection
and cannot easily be modified then a small modification to the main
server might be to fd-pass idle sockets to the container server which
monitors for activity and then fd-passes active sockets back to the
heavy-state server.

To be fair, the number of use-cases are not large and some use-cases
might be implemented just as easily with fuse or similar. But when you
do have a use-case, now you have a guide to get you up and running.

### How to use

1. Run `make`

2. Run `./server` in one terminal

3. Run `./client` in another terminal

4. Type some lines of text into the `./client` terminal

5. Text should show up on the `./server` terminal

### Runs on?

This demo is know to work on Linux, FreeBSD and macOS using go1.11.6
and beyond.


**--**

**Search terms**: go, golang, fd-passing, CMSG, Control Message, SCM_RIGHTS, recvmsg, sendmsg, Unix.
