# go-scraping-tiktok

# Run in local docker

## build the image
docker build -t go-tiktok-scrapper .

## run with env params below
docker run -p 8000:8000 --cap-add=SYS_ADMIN --shm-size=1g go-tiktok-scrapper

## test API
http://localhost:8000/search?query=golang