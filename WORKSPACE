workspace(name = "llvm-toolchain")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "69de5c704a05ff37862f7e0f5534d4f479418afc21806c887db544a316f3cb6b",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/v0.27.0/rules_go-v0.27.0.tar.gz"],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "de69a09dc70417580aabf20a28619bb3ef60d038470c7cf8442fafcf627c21cb",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz"],
)

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "db427107e4c36a3713c24daf5fcdf7719fe664e0c36e089b5b16d15cf331ce60",
    strip_prefix = "rules_docker-8a4f73fb29a64ba813087220b200f49a1ca10faa",
    urls = ["https://github.com/bazelbuild/rules_docker/archive/8a4f73fb29a64ba813087220b200f49a1ca10faa.tar.gz"],
)

http_archive(
    name = "com_grail_bazel_toolchain",
    sha256 = "fde3d653262a1c7c56dd9bb08ffc4e26b585838cfc5039b02433d986b59a871b",
    patch_args = ["-p1"],
    patches = ["//:com_grail_bazel_toolchain_m1.patch"],
    strip_prefix = "bazel-toolchain-560cf0f5b796d68ba758565e8906ac900a056b5a",
    url = "https://github.com/aspect-build/bazel-toolchain/archive/560cf0f5b796d68ba758565e8906ac900a056b5a.tar.gz",
)

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

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.17.2")

gazelle_dependencies()

load(
    "@bazel_gazelle//:deps.bzl",
    "go_repository",
)

go_repository(
    name = "com_github_confluentinc_confluent_kafka_go",
    build_directives = [
        # The go_rdkafka_generr is an internal tool to generate error
        # constants from librdkafka. Since it doesn't play well with
        # Gazelle, it's safe to exclude it from the Bazel target generation.
        "gazelle:exclude kafka/go_rdkafka_generr/**/*",
        # The pkg-config is not supported by Gazelle, so we disable the
        # auto-generated LibrdkafkaLinkInfo constant.
        "gazelle:exclude kafka/build_*.go",
    ],
    build_file_proto_mode = "disable_global",
    importpath = "github.com/confluentinc/confluent-kafka-go",
    # Patch for statically linking against the vendored build of librdkafka in the repo
    patches = [
        "//:com_github_confluentinc_confluent_kafka_go.patch",
        "//:com_github_confluentinc_confluent_kafka_go_vendor.patch",
    ],
    sum = "h1:YxM/UtMQ2vgJX2gIgeJFUD0ANQYTEvfo4Cs4qKUlmGE=",
    version = "v1.6.1",  # See comments on the build_directives in case this version changes.
)

# Linux sysroot for cross-compilation OSX -> linux
http_archive(
    name = "linux_sysroot",
    build_file_content = """
filegroup(
    name = "sysroot",
    srcs = glob(["*/**"]),
    visibility = ["//visibility:public"],
)
filegroup(
    name = "ensure",
    srcs = [],
    visibility = ["//visibility:public"],
)
""",
    sha256 = "1b2464c308f1c8767f20887354c831c92412ff43b881cf85b121a49337724bb9",
    urls = [
        "https://github.com/aspect-build/debian-sysroot-image-creator/releases/download/87c56f3/debian_stretch_amd64_sysroot.tar.xz",
    ],
)

load("@io_bazel_rules_docker//toolchains/docker:toolchain.bzl",
    docker_toolchain_configure="toolchain_configure"
)

# Hack for rules_docker to be able to find `docker` on M1 mac
docker_toolchain_configure(
  name = "docker_config",
  docker_path="/usr/local/bin/docker",
)

load("@io_bazel_rules_docker//container:container.bzl", "container_pull")

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

container_deps()

container_pull(
    name = "golang_base_image",
    registry = "index.docker.io",
    repository = "golang",
    # 1.17.6-bullseye
    digest = "sha256:ec67c62f48ddfbca1ccaef18f9b3addccd707e1885fa28702a3954340786fcf6"
)
