ABODEMINE_PROJECT_NAME := saas-whitelabel

OUT_DIR ?= ${ABODEMINE_WORKSPACE}/.local/build/clients-react-$(ABODEMINE_PROJECT_NAME)

build:
	pnpm run build --mode ${ABODEMINE_NAMESPACE} --outDir ${OUT_DIR}

dev:
	pnpm run dev --host 0.0.0.0 --port ${ABODEMINE_CLIENTS_APP_PORT}
