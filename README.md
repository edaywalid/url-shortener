# URL SHORTENER

## Description

a URL shortener distributed system that stores and resolves short links using Redis and Zookeeper.

## Requirements

- Go >= 1.22.5
- Docker
- Makefile

## Installation

1. Clone the repository

```bash
   git clone https://github.com/edaywalid/url-shortner.git
```

2. Change directory to the project directory

```bash
   cd url-shortner
```

## Usage

1. If you want to run a single instance u can run the following command

```bash
   make run
```

2. If you want to run a distributed system you can run the following command

```bash
   make runm
```

- &#9658; **_NOTE:_**

  - To stop the system you can run the following command
    ```bash
       make stop
    ```
  - or you can run the following command to stop the distributed system
    ```bash
       make stopm
    ```

## API

there are two endpoints in the system

1. **Shorten URL**

   - **Success Response:**

     - **Code:** 200 <br />
       **Content:** `{ "short_url": "http://localhost/332c" }`

   - **Sample Call:**

   ```bash
   curl -H 'Content-Type: application/json' \
   -d '{ "original_url" : "https://www.youtube.com"}' \
   -X POST \
   http://localhost/shorten
   ```

2. **Resolve URL**

   - **Success Response:**

     - **Code:** 200 <br />
       **Content:** `<a href="https://www.youtube.com">Moved Permanently</a>.`

   - **Sample Call:**

   ```bash
   curl -X GET http://localhost/332c
   ```

- **_NOTE:_**
  - the single instance will run on port 8080 so use `localhost:8080`
  - the distributed system will run on port 80 so use `localhost`
