# ==============================================================================
# Makefile helper functions for generate necessary files
#

.PHONY: gen.run
gen.run: gen.clean gen.errcode gen.protobuf

.PHONY: gen.errcode
gen.errcode: gen.errcode.code gen.errcode.doc

.PHONY: gen.errcode.code
gen.errcode.code: tools.verify.codegen
	@echo "===========> Generating bifrost error code go source files"
	@codegen -type=int -fullname=Bifrost ${ROOT_DIR}/internal/pkg/code

.PHONY: gen.errcode.doc
gen.errcode.doc: tools.verify.codegen
	@echo "===========> Generating error code markdown documentation"
	@codegen -type=int -fullname=Bifrost -doc \
		-output ${ROOT_DIR}/docs/guide/zh-CN/api/error_code_generated.md ${ROOT_DIR}/internal/pkg/code

.PHONY: gen.protobuf
gen.protobuf: tools.verify.protoc-gen-go
	@echo "===========> Generating gRPC protobuf go source files"
	@protoc -I=${ROOT_DIR} --go_out=plugins=grpc:${ROOT_DIR} ${ROOT_DIR}/api/protobuf-spec/bifrostpb/v1/bifrost.proto

.PHONY: gen.defaultconfigs
gen.defaultconfigs:
	@${ROOT_DIR}/scripts/gen_default_config.sh

.PHONY: gen.clean
gen.clean:
	@rm -rf ./api/protobuf-spec/bifrostpb/*/*.pb.go
	@$(FIND) -type f -name '*_generated.go' -delete
