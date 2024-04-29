The p2pk (pay to pubkey) script we executed in last section has some problems. First the uncompressed public is too long, it is 65 bytes in binary data, when encode it into string,
the length can be doubled, even the compressed SEC format is still too long and difficult to be handled by wallet. Second it is no secure enough when machine exachange info using
uncompressed SEC format. And Bitcoin nodes need to indexed all outputs for certain purpose, the uncompressed SEC format will requires more disk resources on the part of nodes.

