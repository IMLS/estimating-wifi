# input-initial-configuration

This tool is intended to be used as part of the bootstrap. It reads in information that is critical to the operation of the Pi as a session counter.

For use in the bootstrap, running 

```
sudo input-initial-configuration --all --write
```

is appropriate. This will prompt for all inputs, and write them to `/etc/session-counter/auth.yaml`. Using `--path` will allow for a new location to be chosen for the config.

* `--fcfs-seq` prompts for the FCFS Seq Id (checked via regex)
* `--tag` asks for their hardware id/tag (user chosen, freeform)
* `--word-pairs` prompts for word pairs that are decoded into an api key (checked via wordlist)

## for testing

For developers

```
input-initial-configuration --plain-token
```

will prompt for an API key, and then write out the encrypted version of that key. 

