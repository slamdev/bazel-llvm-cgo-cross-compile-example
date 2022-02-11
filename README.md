# Cross-compiling c/c++ toolchain example (Apple M1 support) of a go_binary to CGO enabled

On macos,

```
bazel run //:main
```

Should build the `//:main` go_binary for darwin (Intel or M1) and run it. This binary has CGO enabled as it depends on `github.com/confluentinc/confluent-kafka-go`. We patch the `com_github_confluentinc_confluent_kafka_go` to add the `cc_import` targets to the external repository and set `cdeps`, `clinkopts` and `copts`.

```
bazel run //:main_image
```

Will build the docker image which depends on the `//:main`. In this case `//:main` will be build to linux using the cross-compiling toolchain `@com_grail_bazel_toolchain` (with patches to support M1).

The cross-compiling LLVM toolchain is setup in the WORKSPACE with 

```
load("@com_grail_bazel_toolchain//toolchain:deps.bzl", "bazel_toolchain_dependencies")

bazel_toolchain_dependencies()

load("@com_grail_bazel_toolchain//toolchain:rules.bzl", "llvm_toolchain")

llvm_toolchain(
    name = "llvm_toolchain",
    llvm_version = "13.0.0",
    # https://github.com/aspect-build/llvm-project/releases/download/llvmorg-13.0.0/clang+llvm-13.0.0-arm64-apple-darwin.tar.xz
    llvm_mirror = "https://github.com/aspect-build/llvm-project/releases/download/llvmorg-",
    sysroot = {
        "linux-x86_64": "@linux_sysroot//:sysroot",
    },
)

load("@llvm_toolchain//:toolchains.bzl", "llvm_register_toolchains")

llvm_register_toolchains()
```

We set an `llvm_mirror` for M1 support which points to a M1 release of LLVM 13.0.0 at https://github.com/aspect-build/llvm-project since LLVM does not yes release M1 builds.

We also set a linux `sysroot` for cross-compiling pulled in from

```
# This sysroot is used by github.com/vsco/bazel-toolchains.
http_archive(
    name = "org_chromium_sysroot_linux_x64",
    build_file_content = """
filegroup(
  name = "sysroot",
  srcs = glob(["*/**"]),
  visibility = ["//visibility:public"],
)
""",
    sha256 = "84656a6df544ecef62169cfe3ab6e41bb4346a62d3ba2a045dc5a0a2ecea94a3",
    urls = ["https://commondatastorage.googleapis.com/chrome-linux-sysroot/toolchain/2202c161310ffde63729f29d27fe7bb24a0bc540/debian_stretch_amd64_sysroot.tar.xz"],
)
```
