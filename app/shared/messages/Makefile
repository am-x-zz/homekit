NANO_PROTOC=~/Downloads/nanopb-0.3.9.3-macosx-x86/generator-bin/protoc
PROTOC=protoc

gen-proto: gen-proto-go gen-proto-nano

gen-proto-go:
	$(PROTOC) -I=. --go_out=. messages.proto

gen-proto-nano:
	rm -f messages.pb.*
	$(NANO_PROTO_CMD) -I=. -I=/Users/alexander/Downloads/nanopb-0.3.9.3-macosx-x86/generator/nanopb --nanopb_out=. messages.proto
	cp messages.pb.* ~/arduino/homekit5/lib/native
	mv messages.pb.* ~/arduino/homekit4/lib/native