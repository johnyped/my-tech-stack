.PHONY: help new-service

help: ## Show this help
	@grep -E '^[a-zA-Z0-9_.-]+:.*## ' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

new-service: ## Scaffold a new service. See document/SCAFFOLD_CHECKLIST.md for all flags.
	script/new-service.sh \
		--name "$(NAME)" \
		--module "$(MODULE)" \
		$(if $(PORT),--port "$(PORT)",) \
		$(if $(CACHE),--cache "$(CACHE)",) \
		$(if $(DB),--db "$(DB)",) \
		$(if $(KAFKA),--kafka "$(KAFKA)",) \
		$(if $(STORAGE),--storage "$(STORAGE)",) \
		$(if $(FIREBASE),--firebase "$(FIREBASE)",) \
		$(if $(MOCKOON),--mockoon "$(MOCKOON)",) \
		$(if $(MOCKTCP),--mocktcp "$(MOCKTCP)",) \
		$(if $(OUT),--out "$(OUT)",) \
		$(if $(FORCE),--force,)
