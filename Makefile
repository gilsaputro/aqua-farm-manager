.PHONY: deps-init
deps-init:
	echo "INFO: creating dependencies..."
	@docker-compose up -d --build
	@bash ./scripts/init-dep.sh ./scripts/init-vault ./schema/vault ./scripts/init-es ./schema/elasticsearch ./scripts/init-nsq ./schema/nsq
	
.PHONY: deps-tear
deps-tear:
	@docker-compose down --volumes --remove-orphans