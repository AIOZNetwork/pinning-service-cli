# pinning-cli

This application provides a Command Line Interface for Pinning APIs.

### Build pinning-cli binary directly from source

Requirements :
- git
- go

Build it:
```
cd pinning
go build
```

## Commands

### Available commands

- login:          Save api key and secret key to configuration file `~/.pinning/credentials`
- pin :           Add a hash to Pinning for asynchronous pinning
- get-pin:        Retrieves a pinned item from the IPFS network by ID
- list-pins :     List of user pins
- info :          Get total information about your content
- unpin :         Unpin a hash from Pinning

### Help

```
pinning --help
```

There is many flags available on each command, so do not forget to run help on each of them for more details on these flags.

### Operation result and logs

Operation result is available on stdout in JSON format.

Logs are available on stderr.

## Store credentials

Credentials to call Pinning APIs are sourced in the following order :
- from command line flags `login` `--key` and `--secret`
- from a configuration file `~/.pinning/credentials` 

### Command line login
```
pinning login --key "your_key" --secret "your_secret"
```

### Configuration file
Configuration file is written in JSON.

```
{"APIKey":"your_key","SecretKey":"your_secret"}
```

## Example usage

You can pin a file
```
pinning pin --file afile.txt
```

You can pin a whole directory
```
pinning pin --file ../some/where
```

You can wrap the directory name in the parent hash
```
pinning pin --file ../some/where -w
```

You can choose a name for display in Pinning
```
pinning pin --file afile.txt --name "a_name"
```

You can choose add metadata for your own usage
```
pinning pin --file afile.txt --keyvalue key1:value1 --keyvalue key2:value2
```

You can add a hash to be pinned
```
pinning pin --hash QmdYTBNig2d4dQd5o1LXM3NHbCYA7168NpN5R9m44vDj88 --keyvalue key1:value1 --keyvalue key2:value2
```

Get pin by id
```
pinning get-pin --id=00000000-0000-0000-0000-000000000000
```

Or list all your pins
```
pinning list-pins
```

Or your custom
```
./pinning list-pins --pinned=true --sortBy=created_at --sortOrder=DESC --limit=10 --offset=0 --keyvalue=key1:value1 --keyvalue=key2:value2
```

And finally unpin a hash
```
pinning unpin --id=00000000-0000-0000-0000-000000000000
```