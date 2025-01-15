# Using Docker Compose

- Based on `dockerfile` course
- Use `docker compose` instead of plain `docker build` and `docker run`.
  Easier to manage `ENV` (including secrets), multi container, and 
  dependencies.
- `docker-compose.yaml` in parent directory.
- Use `HOST` and `PORT` from `ENV`

## Project Acceptance Criteria:
- Modify your code from the first project and create docker-compose.yaml file, so this application can run on top of docker and run using docker-compose 
- Web server application should run on docker container using docker-compose


## Directory Structure
      - /path/to/base/directory
        - golang  => this repo, cloned from git.
                   rename cloned directory name if needed.
        - docker-compose.yaml  => copied and adjusted from 
                   golang/docker-compose-example.yaml

## Quick Start
- Enter `/path/to/base/directory` (e.g. `~/app/project01`, or
  whatever it is on your setup)
- clone this repo to a directory named `golang`, e.g. 
  (adjust repo URL as needed)
      git clone https://git-server/path/repo-name.git golang
- `cp golang/docker-compose-example.yaml docker-compose.yaml`
  This will be your local configuration, not synced to git. 
  Necessary because env variables in dev (including db connection)
  will be on this file.
- Build and start services: 

      docker compose up --build -t1 -d
- Test: 

      curl http://localhost:10081
- Stop services: 
      
      docker compose down -t1
- Asciinema Preview

[![asciicast](https://asciinema.org/a/567949.png)](https://asciinema.org/a/567949)