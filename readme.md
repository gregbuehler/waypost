# waypost ![status](https://travis-ci.org/gregbuehler/waypost.svg?branch=master "travis status")

A simple dns resolver with whitelist/blacklist functionality.


## todo

- Repository should be a radix-tree
- Pool could/should implement weighted members
- Add mgmt functionality to fetch lists
    ```console
    waypost blacklist fetch http://foo.example.com/blocklist.txt
    ```
