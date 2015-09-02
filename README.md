# Swift Resource

Versions objects in an Swift container, by pattern-matching filenames to identify
version numbers.

This resource is heavily inspired by the the [offical S3 resource](https://github.com/concourse/s3-resource)

## Source Configuration

* `username`: *Required.* The username to use when accessing the
  container.

* `api_key`: *Required.* The api secret key to use when accessing
  the container.

* `domain`: *Required.* The domain of the user (V3) 

* `auth_url`: *Required.* The url for authentication (Keystone) 

* `regex`: *.* The pattern to match filenames against. The first
  grouped match is used to extract the version. The extracted version
  needs to be parsable by [go-version](https://github.com/hashicorp/go-version)

## Behavior

### `check`: Extract versions from the bucket.

Objects will be found via the pattern configured by `regex`. The versions
will be used to order them (using [semver](http://semver.org/)). Each
object's filename is the resulting version.

### `in`: Fetch an object from the bucket.

Places the following files in the destination:

* `(filename)`: The file fetched from the bucket.
* `filename`: Contains the filename of the fetched file 
* `version`: Conatins the version identified in the file name.

#### Parameters

*None.*


### `out`: Upload an object to the bucket.

Given a path specified by `from`, upload it to the S3 bucket, optionally to
a directory configured by `to`. The path must identify a single file.

#### Parameters

* `from`: *Required.* A regexp specifying the file to upload.

## Example Configuration

### Resource

``` yaml
- name: release
  type: swift 
    username: user
    api_key: secret
    auth_url: http://somewhere.at.keystone:5000/v3
    domain: comedomain
    container: mycontainer
    regex: file-([.0-9])+\.tgz
```

### Plan

``` yaml
- get: release
```

``` yaml
- put: release
  params:
    from: a/release/path/file-(.*).tgz
```

