# ss33

*ss33* helps you use a local S3-compatible store in front of real S3 for speedier uploads and downloads.

## Installation

```
$ go get github.com/mark-rushakoff/ss33
$ go install github.com/mark-rushakoff/ss33
```

## Usage

To upload a file that is on neither the permanent nor cache store:

```
$ ss33 put \
    --file /path/to/file.tgz \
    --permanent-secret-access-key $S3_SECRET_ACCESS_KEY \
    --permanent-access-key-id $S3_ACCESS_KEY_ID \
    --permanent-bucket your-bucket \
    --permanent-key "path/within/your/bucket/file.tgz" \
    --cache-endpoint "riak.example.com" \
    --cache-secret-access-key $RIAK_SECRET_ACCESS_KEY \
    --cache-access-key-id $RIAK_ACCESS_KEY_ID \
    --cache-bucket your-cache-bucket \
    --cache-key "path/within/your/cache/bucket/file.tgz"
```

Note that those are the longhand versions of those flags, for clarity's sake.
Check the in-app help (`ss33 --help`) for shorthand versions.

To download a file that is definitely on the permanent store and may or may not be on the cache store:

```
$ ss33 get \
    --file /path/to/file.tgz \
    --permanent-secret-access-key $S3_SECRET_ACCESS_KEY \
    --permanent-access-key-id $S3_ACCESS_KEY_ID \
    --permanent-bucket your-bucket \
    --permanent-key "path/within/your/bucket/file.tgz" \
    --cache-endpoint "riak.example.com" \
    --cache-secret-access-key $RIAK_SECRET_ACCESS_KEY \
    --cache-access-key-id $RIAK_ACCESS_KEY_ID \
    --cache-bucket your-cache-bucket \
    --cache-key "path/within/your/cache/bucket/file.tgz"
```

## License

ss33 is available under the terms of the MIT license.
See [LICENSE.MIT.txt](LICENSE.MIT.txt) for full details.
