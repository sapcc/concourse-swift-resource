# Swift Resource

Versions objects in an Swift container, by pattern-matching filenames to identify
version numbers.

This resource is inspired by the the [official S3 resource](https://gthub.com/concourse/s3-resource)

## Source Configuration

* `auth_url`: *Required.* The url for authentication (Keystone) 

* `username`: *Required.* The username to use when accessing the
  container.

* `api_key`: *Required.* The password of the given user.

* `domain`: *Required.* The domain of the user 

* `tenant_id`: *Optional.* The project id

* `container`: *Required* The swift container

* `regex`: *.* The pattern to match filenames against. The first
  grouped match is used to extract the version. The extracted version
  needs to be parseable by [go-version](https://github.com/hashicorp/go-version)

## Behaviour

### `check`: Extract versions from the container.

Objects will be found via the pattern configured by `regex`. The versions
will be used to order them (using [semver](http://semver.org/)). Each
object's filename is the resulting version.

### `in`: Fetch an object from the container.

Places the following files in the destination:

* `(filename)`: The file fetched from the container.
* `filename`: Contains the filename of the fetched file 
* `version`: Contains the version identified in the file name.

#### Parameters

*None.*


### `out`: Upload an object to the container.

Given a path specified by `from`, upload it to the Swift container. The path must identify a single file. The filename must conform to the `regex` specified in the resource.

#### Parameters

* `from`: *Required.* A regexp specifying the file to upload.

## Example Configuration

### Resource

``` yaml
- name: myapp
  type: swift 
    username: user
    api_key: secret
    auth_url: http://somewhere.at.keystone:5000/v3
    domain: somedomain
    container: mycontainer
    regex: myapp-([.0-9])+\.tgz
```

### Plan

``` yaml
- get: myapp
```

``` yaml
- put: myapp
  params:
    from: a-task/build/path/file-(.*).tgz
```

