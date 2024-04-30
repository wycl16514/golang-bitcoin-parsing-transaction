The p2pk (pay to pubkey) script we executed in last section has some problems. First the uncompressed public is too long, it is 65 bytes in binary data, when encode it into string,the length can be doubled, even the compressed SEC format is still too long and difficult to be handled by wallet. 

Second it is no secure enough when machine exachange info usinguncompressed SEC format. And Bitcoin nodes need to indexed all outputs for certain purpose, 
the uncompressed SEC format will requires more disk resources on the part of nodes.

In order to handle problems in p2pk script, there is a updated script format called p2pkh, in this script ,it will use sha256 and ripemd160 hashing to make the 
address much shorter compared with the uncompressed public key. The structure of p2pkh script is like following:


![bitcoin_script](https://github.com/wycl16514/golang-bitcoin-parsing-transaction/assets/7506958/e330c0ee-c047-4fa7-b736-10c1e71611a8)

the first two element of the script commands are data elements therefore they are push to stack directyly:

![bitcoin_script](https://github.com/wycl16514/golang-bitcoin-parsing-transaction/assets/7506958/0a69fe30-96be-4be6-b8b3-df1b59e78b8d)

Notice we push the signature first and then the pubkey, that's why the pubkey is at the top of stack. Then we execute the OP_DUP, this command duplicate the top 
elment on the stack and push it to the stack that cause the stack has two pubkey element at the top:

![bitcoin_script (1)](https://github.com/wycl16514/golang-bitcoin-parsing-transaction/assets/7506958/be34ff79-6081-4152-8c27-06d1c18021c1)

Now take the OP_HASH160 command from the script and take the top element from the stack, do hash160 operation on the elment and push the result to the top of the 
stack:

![bitcoin_script (2)](https://github.com/wycl16514/golang-bitcoin-parsing-transaction/assets/7506958/b15bf706-7c42-43a5-8d52-1d94a0af5fea)

Now there is a hash element on the script, then we push it directly to the stack, this results in two hash elements on the stack:

