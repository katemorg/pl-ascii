# Image to Ascii Art API

This module serves as an HTTP service that accepts an image, then returns an [ASCII version](https://en.wikipedia.org/wiki/ASCII_art) of that image.

**Contents:**

- [Setup](#setup)
- [Endpoints and Parameters](#endpoints-and-parameters)
- [Technical Overview](#technical-overview) 

## Setup

The `Makefile` has a number of commands set up to install dependencies, run, and test the application.

### Prerequisites

Have Golang 1.19 and `make` installed.

### Running Locally

Clone the repository, then run the below commands. 

```
$ cd <path/to>/pl-ascii
$ make install
$ make run
```

If you do not have make installed, you can run the following to install, run, and test the service:

```
$ cd <path/to>/pl-ascii
$ go install
$ go run main.go
$ cd ascii && go test
```

The service is set to run on port 8090. You can change this port in main.go.


## Endpoints and Parameters

The Ascii API contains the following endpoints:

#### GET Healthcheck /ping

Simple healthcheck endpoint that returns PONG! if the service is running.

Example request:

```
curl 'http://localhost:8090/ping' 
```

Example response:

```
PONG!
```


#### POST Image to Ascii /image-to-ascii

This endpoint takes in a file in multipart/form-data and returns ascii art and details about the art in application/json. Note that the endpoint supports multiple images. 

Example request:

```
 curl -F 'media=@path/to/local/file’ 'http://localhost:8090/image-to-ascii'
```

Parameters:

| **Field**       | **Type**     | **Required?** | **Description**                                              |
| --------------- | ------------ | ------------- | ------------------------------------------------------------ |
| media           | file         | true          | A .jpg or .png image file. Can provide more than one.        |



Example response:

```
[
	{
		"Filename": "cat2.jpg",
		"AspectRatio": "12:6",
		"Art": "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@#*#@@@@@@@@@@\n@@@@@@@@@@#===+%@@@@@@@@@@@@@@%+=+=+%@@@@@@@@@\n@@@@@@@@@@#+=====+***#%@@#+=====++=+%@@@@@@@@@\n@@@@@@@@@@%+-==--==-----:::--=====-*@@@@@@@@@@\n@@@@@@@@@@@%=--=====-:...:-+#*=---+%@@@@@@@@@@\n@@@@@@@@@@@@+:-=++*%#=:.:=+*%*==++=#@@@@@@@@@@\n@@@@@@@@@@@%+-==++=:..=++=---=++=+++#@@@@@@@@@\n@@@@@@@@@@@#+===++-:.::+*==+++***++#@@@@@@@@@@\n@@@@@@@@@@@%=-==--=++*%@@@@@%##**++*%@@@@@@@@@\n@@@@@@@@@@@%+-====+*%@%#**#%%#****+#@@@@@@@@@@\n@@@@@@@@@@%%#=-======++====+++**+=+#@@@@@@@@@@\n@@@@@@@@@%%#+-:-===++++*******++==-=%@@@@@@@@@\n#*++====-===:..:=++++**######**+==-:-#@@@@@@@@\n-::--=======: ....:::-=+*##*++==--::.=%@@@@@@@\n====+===-===:........::-=====---::...:*@@@@@@@\n=+++=====+=-.........:::::----::.....:*@@@@@@@\n"
	}
]
```

The response contains the following top-level fields:

| **Field**        | **Type**              | **Description**                                              |
| ---------------- | --------------------- | ------------------------------------------------------------ |
| Filename         | string                | The name of the original image file.                         |
| AspectRatio      | string                | The aspect ratio of the returned ascii art.                  |
| Art              | json                  | Json of the ascii art to parse.                              |


The art response is in JSON format. The client is expected to parse this json to view the art. A [free online json parser](http://json.parser.online.fr/) is helpful for testing. 
<br><br>
![cat2](https://github.com/katemorg/pl-ascii/assets/27252257/d55b63ed-f6e6-4271-a4ad-191a0f42c3dc)
<img width="464" alt="image" src="https://github.com/katemorg/pl-ascii/assets/27252257/159137d8-c803-490d-8e2f-399567f3ad4b">



## Technical Overview

### Libraries

This service uses two libraries: 

[github.com/gorilla/mux](https://pkg.go.dev/github.com/gorilla/mux)https://pkg.go.dev/github.com/gorilla/mux is used as the request router. Right now the service is simple and could rely only on standard net/http functionality, but mux allows us to scale the service with easier url matching, request routing, and adding subrouters. We can easily enforce request types and headers on the endpoints. This package is actively maintained. 

[github.com/stretchr/testify](https://pkg.go.dev/github.com/stretchr/testify/assert)https://pkg.go.dev/github.com/stretchr/testify/assert is a widely used, actively maintained testing library that I have imported to make use of the assert package. This is a personal preference, as I find the assertions to be less verbose and more readable than the standard "testing" package. 
