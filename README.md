# hctx

Small project for myself to learn some new things and make one
part of my work day just a little bit easier.

For contributing, please see out [contributing docs](./CONTRIBUTING.md).

## Getting Started

### Installation

Download the latest [release](https://github.com/Shackelford-Arden/hctx/releases/latest) for your OS.

Place it wherever you'd like, so long as it is in your $PATH (or equivalent for your shell).

#### dist

If you have [dist](https://github.com/ekristen/distillery) installed, you can install `hctx` like this:

```shell
dist install github/Shackelford-Arden/hctx
```

### Configure Shell

The following assumes that `hctx` is in your `PATH`. If it is not, you will need to either
make it available in `PATH` or update the path to hctx's binary.

#### Bash/ZSH

Add this in your `.bashrc` / `.zshrc`:

```shell
eval "$(hctx activate)"
```

Then either start a new shell, or import to your current session:

```shell
# Bash
source ~/.bashrc

# zsh
source ~/.zshrc
```

#### Nushell

To use hctx with Nushell, you'll need to add this to your `config.nu`:

```nu
# need to make sure the autoload directory exists
mkdir ~/.config/nushell/autoload
# Nu doesn't have an equivalent to Bash's `eval`, so we have to save the script to disk
# This does mean that if future versions change the script, this will need to be done
hctx activate | save -f ~/.config/nushell/autoload/hctx.nu
```

##### Notes

In future versions of `hctx`, we may make this more of a "managed" experience, but opting for the manual
path for now until an optimal path exists. Thinking about doing something like `hctx init nushell` and have
this be what gets planed in the user's `config.nu`.

### Define Your Configuration

`hctx` assumes the config file will be in `~/.config/hctx/config.hcl`. If it doesn't exist, it will create an empty
one for you when you first run it.

Here is an example configuration file:

```hcl
stack "local" {
  nomad {
    address = "http://localhost:4646"
  }

  consul {
    address = "http://localhost:8500"
  }
}

stack "prd" {
  alias = "PRODUCTION - CAREFUL!"

  nomad {
    address = "https://fancy.cluster:4646"
  }
}
```

Currently, you can have a block for `nomad`, `consul`, and `vault` in a single stack.

Each product supports `address` and `namespace` as configurable items. 

### Listing Configured Stacks

To view a list of the configured stacks:

```shell
hctx list

# shorthand
hctx l
```

If you already have a stack activated, you should see a `*` next to it's name in the list.

You can also add a `-d` or `--detailed` flag to also see the values configured for each stack.

### Use a Stack

Working off the example config above, let's say we want to use the `prd` stack:

```shell
hctx use prd
```

Doing this will only set the environment variables for the config items in the stack. In this case, it would set
`NOMAD_ADDR` and `HCTX_STACK_NAME`.

Notice that the `alias` for this stick is different from the stack name given at the block level. This can be handy
if/when you configure your shell prompt to potentially change colors depending on a particular environment variable.

### Stop using a Stack

This will remove the environment variables for the current stack.

_Note: hctx will only unset the environment variables that are configured in the config._

```shell
hctx unset
```

### Caching Tokens

_Only applies to Nomad and Consul. Vault has built-in caching._

By default, `hctx` does _not_ attempt to cache any credentials/tokens for Nomad or Consul.

To enable it, simply set the global setting in your config:

```hcl
cache_auth = true
```

With this enabled, `hctx` will store credentials when switching between stacks.

This can be helpful when/if you need to quickly switch between two or more stacks, but
don't want to bother with authenticating each time you switch.

_Note: Using `unset` will cache any token currently set in environment variables._

Preferably, Nomad and Consul CLIs would do the caching for you. If
either implement this in the future, `hctx` will be updated to prefer
those methods over itself.

You can find cache file in `~/.config/hctx/cache.json`.

You can also view the current cache by running:

```shell
hctx cache show
```

#### Sharing Tokens

##### Nomad

Sometimes you're working between multiple Nomad regions/datacenters and the tokens are federated.

To maintain the same token between stacks that:

```hcl
share_nomad_tokens = true
```

Setting this will first check if the current token (value of `NOMAD_TOKEN`) is valid against the target
Nomad cluster. If it is, it will once again be set as the current token.

#### Cache Management

`hctx` includes a few commands to interact with your cache:

* `cache show`
* `cache clean` - This simply finds stacks that have expired tokens and cleans them out.
* `cache clear` - This removes all cache items.

### Shell Prompts

This section contains information about how one _might_ configure the
designated shell prompt to update based on the selected stack.

#### Starship

Site: [starship.rs](https://starship.rs)

This example makes use of Starship's built-in support
for getting values from [environment variables](https://starship.rs/config/#environment-variable).

##### config.toml

```toml
format = """
...
${env_var.HCTX_STACK_NAME}\
...
"""

[env_var.HCTX_STACK_NAME]
variable = 'HCTX_STACK_NAME'
format = 'hctx [$env_value]($style)'
```

## TODO Items

- [x] Add `list` command
  - Default should be a simple list of stack names.
  - With something like a `-verbose` flag (w/ an alias of `-v`!), include full values of each stack
    - Probably table format
- [X] Add self-update
- [x] Add configuration to indicate an environment is production
  - This is "available" by letting users use aliases. Users can update their prompts accordingly.
  - Could potentially come into play w/ shell prompt updating
- [x] Add support for stack aliases
  - Let daily usage use shorter names where shell prompt updating uses slightly more verbose naming

## Maybes

These are items that sound interesting, but feel like maybe getting
outside the scope of what I want to do here.

### Contextual Configs

If for example, there is a `.hctx.hcl` file in the current directory,
assume the user wants to use that over the global config file.

Could maybe be overridden by a `-g` flag.

### Passthrough Commands

Could be interesting to do something like:

```shell
hctx use aws-prd
hctx run nomad node status
```

Internally, `hctx` would look for a `nomad` binary and execute it
with the known variables (setting them as environment variables) and
passing in the sub-commands of the given binary.

I wonder if there is a way to have a shell setup to pass raw commands
such as `nomad` to `hctx` without having to include `hctx run` :eyes:

### Binary Version Management

Name is fairly self-explanatory. Main reason for thinking about this is
there may be times there some environments have version sprawl and
ensuring you are using the same version of the binary as what is in
the environment can help w/ ensuring consistency.

## Inspired By

Along with my own workflows, the following projects inspired me to try
my hand at figuring out something like them. These projects
each have their own places and just didn't fully encapsulate
what I wanted/needed.

- [mise](https://github.com/jdx/mise)
    - This project is pretty cool, and covers _alot_ of installs. Written in Rust,
      it was a bit outside my area of expertise and I wanted to (for now!)
      focus on just a few languages, so figured learn from `mise` and
      see what I could do on my own in Go.
    - I do use `mise` for some things today.
- [nomctx](https://github.com/mr-karan/nomctx)
    - This is a similar tool, but only handles Nomad. Pretty handy if all you care about is Nomad!
- [target-cli](https://github.com/devops-rob/target-cli)
    - Honestly this project seemed to more fully align with what I needed/wanted,
      but I wanted to explore a different project layout and focus on a subset of products.
      Who knows, might come back to this and decide to contribute instead of writing my own!