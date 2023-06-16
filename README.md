## USEFULL COMMANDS

1. update deps, -from_file need point to go.mod file
```bazel run //:gazelle -- update-repos -from_file=projects/image_processing/go.mod -to_macro=deps.bzl%go_dependencies```
2. run image processing project
```bazel run //projects/image_processing:image_processing```
3. test projects
```bazel test //...```
4. build projects
```bazel build //...```



## BAZEL / GAZELLE

1. add gazelle config to WORKSPACE.bazel: https://github.com/bazelbuild/rules_go#generating-build-files
```
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "6b65cb7917b4d1709f9410ffe00ecf3e160edf674b78c54a894471320862184f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.39.0/rules_go-v0.39.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.39.0/rules_go-v0.39.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "727f3e4edd96ea20c29e8c2ca9e8d2af724d8c7778e7923a854b2c80952bc405",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.30.0/bazel-gazelle-v0.30.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.30.0/bazel-gazelle-v0.30.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.19.3")

gazelle_dependencies()
```
2. create BUILD.bazel in main folder + go project folders, with
```
load("@bazel_gazelle//:def.bzl", "gazelle")

gazelle(name = "gazelle")
```
3. run 
```
bazel run //:gazelle
```
4. generate deps.bzl file with 
```
bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies
```

based on: https://medium.com/@simontoth/golang-with-bazel-2b5310d4ce48
