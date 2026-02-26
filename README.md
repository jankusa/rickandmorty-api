# Rick & Morty Backend (Go)

A simple Go backend that provides data from the [Rick and Morty API](https://rickandmortyapi.com/).
It supports two main endpoints:

* `/search` – search for characters, locations, and episodes by name
* `/top-pairs` – top character pairs appearing together in episodes

## Run Locally

1. Clone the repository:

```bash
git clone https://github.com/your-username/rickmorty-api.git
cd rickandmorty-backend
```

2. Download dependencies:

```bash
go mod download
```

3. Start the server:

```bash
go run main.go
```

The server will be available at `http://localhost:8080`.

## Docker

### Build the Docker image:

```bash
docker build -t rickmorty-backend .
```

### Run the container:

```bash
docker run -p 8080:8080 rickmorty-backend
```

The backend will now be accessible at `http://localhost:8080`.
