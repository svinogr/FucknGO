- open terminal in «deployments» directory;
- create network `docker network create fuckngo_network`;
- run PostreSQL:
    - Windows:  
        `docker run --name db -e POSTGRES_PASSWORD=postgres --network fuckngo_network -p 5432:5432 -v /${PWD}/postgresql/data:/var/lib/postgresql/data -v /${PWD}/postgresql/init.sql:/docker-entrypoint-initdb.d/init.sql -d postgres`
    - Linux:  
        `docker run --name db -e POSTGRES_PASSWORD=postgres --network fuckngo_network -p 5432:5432 -v ./postgresql/data:/var/lib/postgresql/data -v ./postgresql/init.sql:/docker-entrypoint-initdb.d/init.sql -d postgres`