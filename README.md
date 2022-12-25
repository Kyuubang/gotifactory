gotifactory
=============
painless artifactory management for GO binary and is highly
customizable. with this, you will make it even easier to implement binary 
self-update like using a normal API

### Repository layout

```
repo
    ├── manifest.json
    └── package
        └── packagebin
```

### Usage 

push binary 

```shell
gotifactory -version 1.0 -pathbin sample/gotifactory -channel latest -commit f2ca1bb6c7 -server http://localhost/scripts/
```

it will be created in your `./repo` directory that contains your binary and
`manifest.json`. So, you can synchronize into your web server directory or s3.

### Manifest guide

manifest used as metadata of binary, that contains `version`, `git commit`,
`download URL`, `checksum`, and also `package name`. Manifest represent as 
JSON data, So you can consume them easily with just as regular API.

==`manifest.json`==

```json
{
  "gotifactory": [
    {
      "channel": "latest",
      "commit": "f2ca1bb6c7",
      "package": "gotifactory",
      "sha256": "01efb3acc22b6e3f2dfee7719c822f69d178c605bc1d5d2208bfced9896ef04f",
      "url": "http://localhost/repo/gotifactory/gotifactory",
      "version": "2.3"
    }
  ]
}
```

#### Self update

if you want to add self update feature to you Go project, I highly
recommend you to use [update](https://github.com/inconshreveable/go-update)




