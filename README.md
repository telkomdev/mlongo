# Mlongo

MongoDB Command Line Utility.

Current features

- (Manage MongoDB Indexes with easy) With Mlongo, you can `create`, `show indexes`, and `drop index` with simple command line program.

### Usage

Build: this command will produce a binary file insde `bin` folder.
```shell
$ make build
```
Export `bin` folder from your environment file. Eg: `.bash_profile` or `.profile` file. So `mlongo` will available in all place.

```shell
$ vi ~/.bash_profile
```

Add this line (Adjust with yours).
```
export MLONGO_HOME=/Users/wuriyanto/Documents/mlongo
export PATH=$PATH:$MLONGO_HOME/bin
```

```shell
$ source ~/.bash_profile
```

Show available options
```shell
$ ./mlongo -h
```
