.PHONY: deps-init
deps-init:
	echo "INFO: creating dependencies..."
	@docker-compose up -d --build
	@bash ./scripts/deps-check.sh ./scripts/init-vault ./schema/vault ./scripts/init-es ./schema/elasticsearch
	
.PHONY: deps-tear
deps-tear:
	@docker-compose down --volumes --remove-orphans