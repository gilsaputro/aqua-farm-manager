test_dir = $(shell  go list ./... | grep -v /testfile | grep -v mock | grep -v pkg | grep -v scripts | grep -v cmd)

.PHONY: test
test:
	@bash ./scripts/test.sh

.PHONY: deps-init
deps-init:
	echo "INFO: creating dependencies..."
	@docker-compose up -d --build
	@bash ./scripts/init-dep.sh ./scripts/init-vault ./schema/vault ./scripts/init-nsq ./schema/nsq
	
.PHONY: deps-tear
deps-tear:
	@docker-compose down --volumes --remove-orphans