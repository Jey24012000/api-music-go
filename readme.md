# Bienvenido a mi API Music REST

Hola, espero que estés muy bien. Esta API sirve para buscar canciones desde las API's de Itunes y ChartLyric.

Esta API tiene 7 rutas, las cuales explicare como funcionan:

## Rutas

La primer ruta es el home de la api.
 - `/api/` GET: la cual nos muestra un mensaje de bienvenida.

 ### Rutas de busqueda de canciones
 - `/api/search` GET: busca canciones desde la API de Itunes.
 Esta requiere de un parámetro _term_ para poder buscar una canción, como por ejemplo `http://localhost:8080/api/search?term=avicii-the-nights`, cada palabra del term debe de estar separada por un  _-_, entre más palabras coloques en el term, más específica se vuelve la búsqueda. Esta trae de vuelta 10 resultados.
 - `/api/searchlyric` GET: busca canciones desde la API de ChartLyrics. Esta requiere de dos parámetros. El primero es un _artist_ en el cual va el nombre del artista y un _song_ en el cual es colocado el nombre de la canción, como por ejemplo `http://localhost:8080/api/searchlyric?artist=avicii&song=the-nights`, esta trae de vuelta un resultado.

 ### Rutas para guardar canciones
 Para poder usar estas rutas necesitamos una validación, en este caso se usa el api.
 - `/api/jwt`  GET: arroja un token que debemos usar para poder usar las demás rutas. Este token se obtiene usando un header _Access_ con valor de _1234_ , el resultado es es un token con la siguiente estructura `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjcyMjkzNTF9.V_J4QiAhRVXkOATK4NbrF1UKlLBPTNyWEqxh0G1R9hs`.
- `/api/mysongs` GET: Esta ruta muestra las canciones que guardemos en nuestra base de datos. Esta tiene autorización, por lo que necesita un header _Token_ con el valor del token generado en la ruta anterior.
- `/api/mysongs` POST: Esta guarda una canción en la base de datos desde la API de Itunes.
Esta tiene autorización, por lo que necesita un header _Token_ con el valor del token generado en la ruta anterior. Esta requiere de un parámetro _term_ para poder buscar una canción, como por ejemplo `http://localhost:8080/api/search?term=avicii-the-nights`, cada palabra del term debe de estar separada por un  _-_, entre más palabras coloques en el term, más específica se vuelve la búsqueda. Esta trae un solo resultado.
- `/api/songs` POST: Esta guarda una canción en la base de datos desde la API de Chartlyric.Esta tiene autorización, por lo que necesita un header _Token_ con el valor del token generado en la ruta anterior. Esta requiere de dos parámetros. El primero es un _artist_ en el cual va el nombre del artista y un _song_ en el cual es colocado el nombre de la canción, como por ejemplo `http://localhost:8080/api/searchlyric?artist=avicii&song=the-nights`, esta trae de vuelta un resultado y lo guarda en la base de datos.


## Tecnologías
Esta API esta hecho con Golang usando Gorm para la conexión de la base de datos. Gorm crea la base de datos y las tablas correspondientes de los modelos creados para un base de datos Postgres. Esta API usa docker.

## Docker
El proyecto tiene un docker-compose.yml que crea el entorno de desarrollo para Postgres y para Golang.

Solo ejecuta:
`docker-compose up --build`

