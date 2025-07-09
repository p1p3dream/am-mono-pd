include $(ABODEMINE_WORKSPACE)/code/make/core.mk

PYTHON_VENV_DIR := .venv

all: venv-create

# uv docs: https://docs.astral.sh/uv/.
venv-create:
	if [ ! -d $(PYTHON_VENV_DIR) ]; then \
		uv venv $(PYTHON_VENV_DIR); \
	fi

	uv sync --all-extras --link-mode=copy

venv-delete:
	rm -rf $(PYTHON_VENV_DIR)

venv-recreate: venv-delete venv-create
