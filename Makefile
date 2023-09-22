__DOCKER=$(shell which docker)

lint:
	goreturns -w cmd internal configs

clean:
	$(__DOCKER) system prune -f
