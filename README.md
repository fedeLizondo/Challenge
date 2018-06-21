# Mercadolibre challenge

## Parte 1 del challenge

Problema tecnico sugerido , conexión en a un servidor en docker , orquestado por docker compose , con servidor en nginex.

__"Al momento de inicializar ​ docker-compose vemos que todo levanta bien, pero no podemos acceder a los métodos de API."__

### Cual fue el problema

El servidor nginx no estaba funcionando como un proxy inverso , la principal causa es que dentro de las configuraciones no se rediregia a la api.
Adicionalmente hay que hacer un cambio en los archivos de configuración ya sea en el archivo `Docker-compose.yml` o en el `default.conf`

### Cual fue la solución

Dentro del archivo de configuración en `default.conf` se pueden observar las siguientes lineas:
```
...
location / {
      proxy_set_header  Host $host;
      proxy_set_header  X-Real-IP $remote_addr;
}
...
```
Se modificaron a
```
...
location / {
    proxy_pass  http://app:8080;
    proxy_set_header  Host $host;
    proxy_set_header  X-Real-IP $remote_addr;
}
...
```
Con esto se soluciona el tema de la redirección a la api , pero tambien se pudo observar que cuando se trataba de accder a http://localhos:8080 , no nos podiamos conectar esto se debia a que dentro de la configuración original del servidor (`default.conf`) esta escuchando el puerto 80 en vez del 8080 con el cual accedemos desde el externamente
```
Server {
    listen       80;
    server_name  localhost;
    ...
```
Si se cambia el parametro listen al puerto 8080 se soluciona, como se muestra en el siguiente ejemplo.
```
Server {
    listen       8080;
    server_name  localhost;
    ...
```
###### Otra solución posible es

en el archivo de configuración `Docker-compose.yml` ,donde originalmente se encontraba asi:
``` yaml
  version: '3'
  services:
    nginx:
      build: nginx
      ports:
        - "8080"
  ...
```
remplazar el ports por `8080:80` donde el container internamente usa el puerto 80 y para acceder externamente se utiliza el puerto 8080
``` yaml
  version: '3'
  services:
    nginx:
      build: nginx
      ports:
        - "8080:80"
  ...
```

### Que aprendí

> Principalmente lo que aprendi en esta parte del challenge fue como es el funcionamiento de docker, como se gestiona y se puede configurar utilizando docker-composite como un archivo yaml
>
> Funcionamiento de nginx, habia trabajado con XAMMP pero nunca tuve que acceder a la configuracion del servidor apache, realmente no es muy complicado realizar una configuración en un servidor nginx, Adicionalmente aprendi que si no corro un `docker-compose build` antes del `docker-compose up` , no actualiza los cambios :sweat_smile:
>
>Realmente me quede asombrado lo sencillo que puede ser poder tener todo el entorno de desarrollo con un solo archivo de configuracion ( ademas de la integracion con github).


## Parte 2 del challenge
### Descripción de la aplicación realizada:
La aplicacion consiste en 5 Endpoints donde dos ya se encontraban por defecto, los cuales permitian crear y buscar un item en la base de datos , incialmente estaba configurado para Sqlite y se paso a Mysql.
Los 3 Endpoint restantes permiten
* Buscar si un archivo almacenado en google drive ,contiene una palabra
* Crear un archivo de texto plano sin contenido en google driver
* Un callback que permite generar el token para la autenticación con google drive , y almacenarlo en una variable de session.

Estructura Basica de la API (Mostrando Imagen Sobre el tema de Autenticacion)

__workflow__
 * Si el usuario quiere crear o buscar un item accede directamente a la URL
 * Si el usuario quiere crear o buscar un archivo
   - Si esta Authenticado accede a la url indicada
   - Si no esta Authenticado , el middleware  primero guardar la url original y los datos para crear un archivo (si los tuviera), lo redirecciona a la URL de google para autenticar, luego de que se autentique es redirreccionado a callback.
 * El callback recupera la url original y guarda el codigo para generar el token y redirecciona a la url original
     - el get Actua normalmenta ya que los paramentros pasan por url
     - el Post verifica si tiene en una variable session un error en los parametros o un json con el titulo y la descripcion, si no los tiene lo toma de el context

* ### la conexion con MYSQL
    ##### ¿Cual fue el problema?
    Habia que realizar el cambio de Sqlite a Mysql server,

    ##### ¿Cual la solución?
    Para realizar el cambio de base de datos no tuve mayores dificultades solamente fue crear el nuevo esquema, cambiar una linea para la conexion a la base de datos y agregar un nuevo repositorio `github.com/go-sql-driver/mysql` dentro del `Dockerfile`.
    Fue investigar sobre los repositorios disponibles y ver documentacion sobre la configuración para la conexion a Mysql

    ``` Go
    ...
    func configDataBase() *sql.DB {
    	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", "user", "userpwd", "db", "db"))
      ...
    ```
    ``` Go
    ...
    func createTable(db *sql.DB) {
    	// create table if not exists
    	sql_table := `
    	CREATE TABLE IF NOT EXISTS items(
    		id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    		name VARCHAR(255),
    		description VARCHAR(255)
    	);`
    	...
    ```
