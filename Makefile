help:
	@echo "Below are some common options to make:"
	@echo
	@echo "  install          Install dependencies"
	@echo "  run              Run the service"
	@echo "  test             Run tests on the host machine"

# Log into ECRs
install:
	go install

run: 
	go run main.go

test: 
	cd ascii && go test