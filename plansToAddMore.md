# When there more WRITE opertions:
    maintain a connection to all the replica and forward write operations
# When there is a more READ ops
    maintain the list of addresses to the repica server and for each write operation create connection, forward write ops and later close that connection

# While working on Accepting commands from master at handshake connection, I forget to pass handshake conn to HandleCmd function and beacuse of which even though master was forwarding WRITE ops but there were no changes being made at replica side.


