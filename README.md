# Globus HTTP CLI Client

This program impliments a subset of the [Globus API for HTTP access to
collections](https://docs.globus.org/globus-connect-server/v5/https-access-collections/).
It currently supports upload and download to [Globus Guest
Collections](https://docs.globus.org/guides/tutorials/manage-files/share-files/).

## Usage

A configuration file (config.toml) should be generated before use.  See the Configuration section.

### Globus HTTP URLs

There are a [few ways to determin the Collection's HTTPS base
URL](https://docs.globus.org/globus-connect-server/v5/https-access-collections/#determining_the_collection_https_base_url).
The easiest is to browse to the collection in the file browser, then select the
"Get Link" button.

### Download

```bash
globus-http download <source-url> <destination-filename>
```

For example:

```
globus-http download https://g-123456.12345.1234.data.globus.org/filename.txt filename.txt
```

### Upload

```bash
globus-http upload <source-filename> <destination-url>
```

For example:

```
globus-http upload filename.txt https://g-123456.12345.1234.data.globus.org/filename.txt 
```

## Configuration

This client needs some configuration values that are read from a TOML file. By default, it will read from a config.toml file in the current directory, but you can change the path using the `-config <filename>` option.  The configuration file should have the client ID, secret and needed scopes. Sample file:


```toml
# ClientID and ClientSecret are Globus Auth client credentials (see the
# Client Credentials section in the readme).
ClientID = "f1a33410-8df7-4df4-b112-a3b77f83b6e6"
ClientSecret = "secret"

# Scopes dictate which collection this client should access.
# For guest collections, the format is:
#   https://auth.globus.org/scopes/<collection-uuid>/https
#
Scopes = [
    "https://auth.globus.org/scopes/d5f73dee-4c7e-4d9e-9731-6439a5b82332/https",
    "https://auth.globus.org/scopes/2a18f093-13e6-467d-91b3-0b1b58f3b18c/https",
]
```



## Client Credentials

Currently the `client_credentials` option is supported for authentication. To generate credentials:

1. Go to the [Globus developers page](https://app.globus.org/settings/developers).
1. Select Register an App -> "Advanced registration".
1. Select a Project, or specify create a new project. If creating a new project, enter project name, contact email, and project admins.
1. Provide an App name.
1. Leave Redirects, Required Identity, Pre-selected Identity Provider, Use effective identity, and Prompt for Named Grant, empty and unchecked.
1. Click register app.
1. The ClientID needed in the configuration file (see section above) is the Client UUID.
1. Click Add Client Secret. Record the provided client secret in the configuration file. See the section above.

## Granting the client access to a guest collection

Before the client can access a guest collection, you need to grant it access to
the desired [Guest Collection]((https://docs.globus.org/guides/tutorials/manage-files/share-files/).

Once you've created a guest collection and client credentials, select Add Permissions on the guest collection and share with `<client-id>@clients.auth.globus.org` (e.g. `f1a33410-8df7-4df4-b112-a3b77f83b6e6@clients.auth.globus.org`). You can select read and/or write as appropriate.

