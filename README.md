# go-scraping-tiktok

# Run in local docker

## build the image
docker build -t go-tiktok-scrapper .

## run with env params below
docker run -p 8080:10000 --cap-add=SYS_ADMIN --shm-size=1g go-tiktok-scrapper

## test
http://localhost:8080/search?query=golang