workspace(name = "dlframework")

git_repository(
    name = "io_bazel_rules_go",
    commit = "570488593c55ad61a18c3d6095344f25da8a84e1",  # master on 2018-01-04
    remote = "https://github.com/bazelbuild/rules_go",
)

git_repository(
    name = "bazel_gazelle",
    commit = "9e43c85089c3247fece397f95dabc1cb63096a59",  # master on 2018-01-09
    remote = "https://github.com/bazelbuild/bazel-gazelle",
)

git_repository(
    name = "com_github_tnarg_rules_go_swagger",
    commit = "f87f37a5e097329f03f3bcb73f942e003cd56164",
    remote = "https://github.com/tnarg/rules_go_swagger.git",
)

git_repository(
    name = "io_bazel_rules_docker",
    commit = "8aeab63328a82fdb8e8eb12f677a4e5ce6b183b1",
    remote = "https://github.com/bazelbuild/rules_docker.git",
)

git_repository(
    name = "build_tools",
    commit = "72c03d74aab4cbe9fe6860f6d5571f2aa292e47c",
    remote = "https://github.com/rai-project/build_tools",
)

git_repository(
    name = "distroless",
    commit = "d061f8989a5fc1c59f2e84bfedb299a5fbb20cb6",
    remote = "https://github.com/GoogleCloudPlatform/distroless",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")
load("@io_bazel_rules_go//proto:def.bzl", "proto_register_toolchains")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("@com_github_tnarg_rules_go_swagger//go/swagger:def.bzl", "go_swagger_deps", "go_swagger_repositories", "go_swagger_repository")
load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
    container_repositories = "repositories",
)
load("@build_tools//:esc.bzl", "esc_repositories")

go_rules_dependencies()

go_register_toolchains()

gazelle_dependencies()

go_swagger_deps()

go_swagger_repositories()

esc_repositories()

container_pull(
    name = "nvidia_cuda_container_amd64",
    registry = "index.docker.io",
    repository = "nvidia/cuda",
    tag = "8.0-cudnn6-devel-ubuntu16.04",
)

container_pull(
    name = "nvidia_cuda_container_ppc64le",
    registry = "index.docker.io",
    repository = "nvidia/cuda-ppc64le",
    tag = "8.0-cudnn6-devel-ubuntu16.04",
)

# This is NOT needed when going through the language lang_image
# "repositories" function(s).
container_repositories()

# container_pull(
#     name = "nvidia_cuda_container",
#     registry = "index.docker.io",
#     repository = select({
#         # see https://github.com/bazelbuild/rules_go/blob/master/go/platform/list.bzl#L31:14
#         "@bazel_tools//platforms:x86_64": "nvidia/cuda",
#         "@bazel_tools//platforms:ppc": "nvidia/cuda-ppc64le",
#     }),
#     tag = "8.0-cudnn6-devel-ubuntu16.04",
# )

go_repository(
    name = "com_github_jteeuwen_go_bindata",
    commit = "a0ff2567cfb70903282db057e799fd826784d41d",
    importpath = "github.com/jteeuwen/go-bindata",
)

go_repository(
    name = "com_github_elazarl_go_bindata_assetfs",
    commit = "30f82fa23fd844bd5bb1e5f216db87fd77b5eb43",
    importpath = "github.com/elazarl/go-bindata-assetfs",
)

go_repository(
    name = "org_golang_google_grpc",
    commit = "6913ad5caedced5d627918609375b057963334a5",
    importpath = "google.golang.org/grpc",
)

go_repository(
    name = "com_github_gogo_protobuf_proto",
    commit = "160de10b2537169b5ae3e7e221d28269ef40d311",
    importpath = "github.com/gogo/protobuf/proto",
)

go_repository(
    name = "com_github_gogo_protobuf_gogoproto",
    commit = "160de10b2537169b5ae3e7e221d28269ef40d311",
    importpath = "github.com/gogo/protobuf/gogoproto",
)

go_repository(
    name = "com_github_golang_protobuf_protoc_gen_go",
    commit = "1e59b77b52bf8e4b449a57e6f79f21226d571845",
    importpath = "github.com/golang/protobuf/protoc-gen-go",
)

go_repository(
    name = "com_github_gogo_protobuf_protoc_gen_gofast",
    commit = "160de10b2537169b5ae3e7e221d28269ef40d311",
    importpath = "github.com/gogo/protobuf/protoc-gen-gofast",
)

go_repository(
    name = "com_github_gogo_protobuf_protoc_gen_gogofaster",
    commit = "160de10b2537169b5ae3e7e221d28269ef40d311",
    importpath = "github.com/gogo/protobuf/protoc-gen-gogofaster",
)

go_repository(
    name = "com_github_gogo_protobuf_protoc_gen_gogoslick",
    commit = "160de10b2537169b5ae3e7e221d28269ef40d311",
    importpath = "github.com/gogo/protobuf/protoc-gen-gogoslick",
)

go_repository(
    name = "com_github_grpc_ecosystem_grpc_gateway_protoc_gen_grpc_gateway",
    commit = "61c34cc7e0c7a0d85e4237d665e622640279ff3d",
    importpath = "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway",
)

go_repository(
    name = "com_github_grpc_ecosystem_grpc_gateway_protoc_gen_swagger",
    commit = "61c34cc7e0c7a0d85e4237d665e622640279ff3d",
    importpath = "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger",
)

go_repository(
    name = "com_github_go_swagger_go_swagger_cmd_swagger",
    commit = "acf3c15f3a1fd86f271220a05558717ec1c61d32",
    importpath = "github.com/go-swagger/go-swagger/cmd/swagger",
)
