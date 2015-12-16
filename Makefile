.PHONY: deploy
deploy:
	goxc -bc="windows linux darwin" -tasks-=go-vet -tasks-=go-test default publish-github
