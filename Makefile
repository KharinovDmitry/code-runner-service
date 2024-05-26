rebuild_executors:
	docker build -t unpivileged_run cmd/unprivilegedRun
	docker build -t python_executor internal/executor/python

up_python_executor:
	docker run --rm --mount type=bind,source=./internal/executor/python/testFiles,target=/home/tmp --name test.py -d -e FILE_NAME=test.py python_executor sleep infinity

exec_python_executor:
	docker exec -i test.py ./unprivrun 5000 5000 python3 tmp/test.py

run_python_executor:
	docker run --rm --mount type=bind,source=./internal/executor/python/testFiles,target=/home/tmp --name test.py -i -e FILE_NAME=test.py python_executor ./unprivrun 3000 1000 python3 tmp/timeLimitTest.py

gen_mocks:
	mockgen -source internal/domain/executor/executor.go -destination internal/domain/executor/mocks/mock_executor.go -package mockExecutor

	mockgen -source internal/domain/service/testRunner.go -destination internal/domain/service/mocks/mock_testRunner.go -package mockService
	mockgen -source internal/domain/service/executorFactory.go -destination internal/domain/service/mocks/mock_executorFactory.go -package mockService

