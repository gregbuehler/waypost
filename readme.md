# waypost

A simple dns resolver with whitelist/blacklist functionality.


## todo
- [ ] Repository should be a radix-tree
- [ ] Pool should implement weighted members
- [ ] Add a management interface to dynamically communicate configuration changes:
    ```
    $ waypost blacklist add foo.com
    $ waypost blacklist del foo.com
    $ waypost blacklist list
    ```
- [ ] Add a repository based approach
    ```
    $ waypost blacklist fetch http://foo.example.com/blocklist.txt
    ```

- [ ] Add nameservers
    ```
    $ waypost upstream add 8.8.8.8
    ```

- [ ] Add dynamic host mapping
    ```
    $ waypost host add dev.foo.test 127.0.0.1
    ```