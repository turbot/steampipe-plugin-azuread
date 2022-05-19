![image](https://hub.steampipe.io/images/plugins/turbot/azuread-social-graphic.png)

# Azure Active Directory Plugin for Steampipe

Use SQL to query infrastructure including users, groups, applications and more from Azure Active Directory.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/azuread)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/azuread/tables)

- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-azuread/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install azuread
```

Run a query:

```sql
select display_name, user_principal_name, user_type from azuread_user;
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-azuread.git
cd steampipe-plugin-azuread
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/azuread.spc
```

Try it!

```
steampipe query
> .inspect azuread
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-azuread/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Azure Active Directory Plugin](https://github.com/turbot/steampipe-plugin-azuread/labels/help%20wanted)
