# hctx

Small project for myself to learn some new things and make one
part of my work day just a little bit easier.

## Usage

### Shell Activation

For this to really work, you will need to have this in your `.bashrc` / `.zshrc`:

```shell
eval "$(hctx activate)"
```

The above assumes that at whatever point you add this, `hctx` is available in your path.

### Use a Stack

```shell
hctx use prod
```

Doing this will only set the environment variables for the config items in the stack. So, for example, if you
only have a Nomad entry in the given stack's config, only the Nomad environment variables will be defined.

For an example config, see [config-example.hcl](./config-example.hcl)

### Stop using a Stack

This will remove the environment variables for the current stack.

```shell
hctx unset
```

Note: There is currently a bug where if you run the `unset` command more than once in a row, it errors out.


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
format   = 'hctx [$env_value]($style)'
```

Then in your 

## TODO Items

- [x] Add `list` command
  - Default should be a simple list of stack names.
  - With something like a `-verbose` flag (w/ an alias of `-v`!), include full values of each stack
    - Probably table format
- [ ] Add self-update
- [ ] Add configuration to indicate an environment is production
  - Could potentially come into play w/ shell prompt updating
- [ ] Add support for stack aliases
  - Let daily usage use shorter names where shell prompt updating uses slightly more verbose naming
- [ ] Add `add` command
- [ ] Add `edit` command
  - I'd want to make sure that a user could modify a single attribute of a stack.


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

Boy, brain really starts going nuts if it is allowed, huh?

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
    it was a bit outside of my area of expertise and I wanted to (for now!)
    focus on just a few languages, so figured learn from `mise` and
    see what I could do on my own in Go.
  - I do use `mise` for some things today.
- [nomctx](https://github.com/mr-karan/nomctx)
  - This is a similar tool, but only handles Nomad. Pretty handy if all you care about is Nomad!
- [target-cli](https://github.com/devops-rob/target-cli)
  - Honestly this project seemed to more fully align with what I needed/wanted,
    but I wanted to explore a different project layout and focus on a subset of products.
    Who knows, might come back to this and decide to contribute instead of writing my own!