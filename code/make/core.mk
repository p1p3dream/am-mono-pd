UNAME_M := $(shell uname -m | tr '[:upper:]' '[:lower:]')
UNAME_S := $(shell uname -s | tr '[:upper:]' '[:lower:]')
UNAME_P := $(shell uname -p | tr '[:upper:]' '[:lower:]')

# Use ":=" to ensure a single UUID is immediately generated
# and will persist for the lifetime of a MAKE.
UUID := $(shell uuidgen | tr '[:upper:]' '[:lower:]')

# Use just "=" to ensure a new UUID is generated
# every time this var is expanded.
UUIDGEN = $(shell uuidgen | tr '[:upper:]' '[:lower:]')
