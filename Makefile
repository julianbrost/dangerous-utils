CMD_SOURCES := $(wildcard cmd/*)
CMD_NAMES := $(notdir $(CMD_SOURCES))
CMD_BINARIES := $(addprefix bin/,$(CMD_NAMES))

.PHONY: all
all: $(CMD_BINARIES)

$(CMD_BINARIES):
	go build -o $@ ./cmd/$(notdir $@)

.PHONY: clean
clean:
	$(RM) $(CMD_BINARIES)
