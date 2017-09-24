[![Build Status](https://travis-ci.org/sapcc/concourse-swift-resource.png?branch=master)](https://travis-ci.org/sapcc/concourse-swift-resource)
# Concourse resource for OpenStack Swift

Versions objects in a Swift container, by pattern-matching filenames to identify
version numbers.

This resource is heavily inspired by the [official S3 resource](https://github.com/concourse/s3-resource)

## Deploying to Concourse

In your bosh deployment manifest, add to the following to the `groundcrew.additional_resource_types`:

```yaml
- image: docker:///databus23/concourse-swift-resource
  type: swift
```

## Source Configuration

* `auth_url`: *Required.* The url for authentication (Keystone)

* `username`: *Required.* The username to use when accessing the
  container.

* `api_key`: *Required.* The password of the given user.

* `domain`: *Required.* The domain of the user

* `tenant_id`: *Required.* The project id

* `container`: *Required* The swift container

* `regex`: *Required* The pattern to match filenames against. The first
  grouped match is used to extract the version. The extracted version
  needs to be parseable by [go-version](https://github.com/hashicorp/go-version)

* `disable_tls_verify`: Disable certificate verification for https connections to keystone and swift. Default: false

## Behaviour

### `check`: Extract versions from the container.

Objects will be found via the pattern configured by `regex`. The versions
will be used to order them (using [semver](http://semver.org/)). Each
object's filename is the resulting version.

### `in`: Fetch an object from the container.

Places the following files in the destination:

* `$objectname`: The object fetched from the container.
* `filename`: Contains the filename of the fetched object ($objectname).
* `version`: Contains the version extracted from the file name.

#### Parameters

*None.*


### `out`: Upload an object to the container.

Given a path specified by `from`, upload it to the Swift container. The path must identify a single file. The filename must conform to the `regex` specified in the resource.

#### Parameters

* `from`: *Required.* A regexp specifying the file to upload.

* `segment_container` : *Optional* Name of the container where segments should be stored. Segments will be created for files larger than 1 GB. (Default: ${container}_segments)
* `segment_size`: *Optional* Size of segments (Default: 1073741824 Byte) (1 GB)
* `delete_after`: *Optional* Instruct Swift to schedule automatic deletion of the Object in `n`seconds (Default: disabled). Useful to define Object retention policies.

## Example Configuration

### Resource

``` yaml
- name: myapp
  type: swift
  source:
    username: user
    api_key: secret
    auth_url: http://somewhere.at.keystone:5000/v3
    domain: somedomain
    container: mycontainer
    regex: myapp-([.0-9])+\.tgz
    tenant_id: tenantid
```

### Plan

``` yaml
- get: myapp
```

``` yaml
- put: myapp
  params:
    from: a-task/build/path/file-(.*).tgz
    segment_container: mysegmentcontainer
    segment_size: 1073741824
    delete_after: 3600

```
