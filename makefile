EXECUTABLE := $(shell basename $(executable) 2>/dev/null)
LINKER := $(shell basename $(linker) 2>/dev/null)

LDFLAGS=-ldflags "-X=main.Linker=$(LINKER) -X=main.Executable=$(EXECUTABLE)"

ifneq ($(filter build run, $(MAKECMDGOALS)), $())
ifndef linker
$(error Makefile argument 'linker' must be set)
endif
ifndef executable
$(error Makefile argument 'executable' must be set)
endif
ifndef objectfiles-dir
$(error Makefile argument 'objectfiles-dir' must be set)
endif
endif

# Creates a compressed archive containing:
# - the target executable
# - the dynamic linker
# - the object files that the target executable depends on
build/embedded-files.tar.gz:
	@echo Compressing embedded files...
	@mkdir -p build/tmp
	@cp $$(dirname $(executable))/* build/tmp
	@cp $$(dirname $(linker))/* build/tmp
	@cp $(objectfiles-dir)/* build/tmp
	@tar -czf build/embedded-files.tar.gz -C build/tmp .

.PHONY: clean
clean:
	@rm -rf build &>/dev/null
	@rm -rf internal/filebundler/embedded-files.tar.gz &>/dev/null

# Builds the wrapper executable
.PHONY: build
build: build/embedded-files.tar.gz
	@echo Building wrapper binary...
	@echo Linker: $(LINKER)
	@echo Executable: $(EXECUTABLE)

	@# We can only embed files that are within the go package that embeds them
	@# so we need to copy the archive file to that directory.
	@cp build/embedded-files.tar.gz internal/filebundler/embedded-files.tar.gz

	go build \
		$(LDFLAGS) \
		-o $(EXECUTABLE) ./cmd/furl/main.go

	@rm internal/filebundler/embedded-files.tar.gz

.PHONY: run
run: build
	./$(EXECUTABLE)
