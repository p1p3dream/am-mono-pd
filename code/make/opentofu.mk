ifdef ABODEMINE_BUILD_ENV
GOMPLATE_BUILD_ENV_DATASOURCE := -d build_env=$(ABODEMINE_BUILD_ENV)?type=application/x-env
GOMPLATE_BUILD_ENV_DEP := abodemine_build_env
endif

%.tf: $(GOMPLATE_BUILD_ENV_DEP) ${ABODEMINE_WORKSPACE}/.env $(ABODEMINE_CONFIG) %.tf.gotmpl
	gomplate \
		$(GOMPLATE_BUILD_ENV_DATASOURCE) \
		-d config=$(ABODEMINE_CONFIG) \
		-d env=${ABODEMINE_WORKSPACE}/.env?type=application/x-env \
		-f $*.tf.gotmpl \
	> $*.tf

	tofu fmt $*.tf
