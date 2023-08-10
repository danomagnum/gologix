If you use a connected message from a logix processor it will send the same message over and over again with the same sequence number.  It sends them at the RPI of the msg command.

You can find the RPI in the msg tag of the tag browser as the ".ConnectionRate" tag.


cipService_GetAttributeAll

ignition handshake:
82 2 32 6 36 1 3 250 6 0 1 2 32 1 36 1 1 0 1 0

kepware handshake:
82        2          32 6 36 1  7 233     6 0                1 2 32 1 36 1       1 0 1 0
service   path size  path       Timeout   msg req size       msg req             path
                     0x20 0x06
                     0x24 0x01
                     connection manager
                     instance 1