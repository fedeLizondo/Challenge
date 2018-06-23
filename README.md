# Mercadolibre challenge

## Parte 1 del challenge

Problema técnico sugerido , conexión a un container en docker , orquestado por docker compose , con servidor en nginx.

__"Al momento de inicializar ​docker-compose vemos que todo levanta bien, pero no podemos acceder a los métodos de API."__

### Cuál fue el problema

El servidor nginx no estaba funcionando como un proxy inverso , la razón de ese comportamiento fue que dentro de las configuración no se redireccionaba a la api.
Adicionalmente hay que hacer un cambio en los archivos de configuración ya sea en el archivo `Docker-compose.yml` o en el `default.conf`

### Cuál fue la solución

Dentro del archivo de configuración en `default.conf` se pueden observar las siguientes líneas:
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
Con esto se resolvía el tema de la redirección a la api , pero también se debe tener en cuenta que cuando se trataba de acceder a http://localhost:8080  no nos podíamos conectar .Esto se debía a que dentro de la configuración original del servidor (`default.conf`) está escuchando el puerto 80 en vez del 8080 con el cual accedemos desde el externamente
```
Server {
    listen       80;
    server_name  localhost;
    ...
```
Si se cambia el parámetro `listen` al puerto `8080` se soluciona, como se muestra en el siguiente ejemplo.
```
Server {
    listen       8080;
    server_name  localhost;
    ...
```
###### Otra solución posible es

en el archivo de configuración `Docker-compose.yml` ,donde originalmente se encontraba así:
``` yaml
  version: '3'
  services:
    nginx:
      build: nginx
      ports:
        - "8080"
  ...
```
reemplazar el parámetro ` ports` por `8080:80` donde el container internamente usa el puerto 80 y el puerto 8080 se encuentra expuesto para acceder externamente
``` yaml
  version: '3'
  services:
    nginx:
      build: nginx
      ports:
        - "8080:80"
  ...
```

### Qué aprendí

> Principalmente lo que aprendí en esta parte del challenge fue como es el funcionamiento de docker, cómo se gestiona y se puede configurar utilizando docker-composite como un archivo yaml
>
> Funcionamiento de nginx, había trabajado con XAMMP pero nunca tuve que acceder a la configuración del servidor apache, realmente no es muy complicado realizar una configuración en un servidor nginx, Adicionalmente aprendí luego de unas extensas horas y un duro análisis de por qué no se aplicaban los cambios en el container de nginx que si no corro un `docker-compose build` antes del `docker-compose up` , no actualiza los cambios
>
>Realmente me quede asombrado lo sencillo que puede ser poder tener todo el entorno de desarrollo con un solo archivo de configuración ( ademas de la integración con github).


## Parte 2 del challenge
### Descripción de la aplicación realizada:
La aplicación consiste en 5 Endpoints en los cuales dos ya se encontraban proporcionados , los cuales permiten crear y buscar un item en la base de datos , inicialmente estaba configurado para Sqlite y se pasó a Mysql.
Los 3 Endpoint restantes permiten
* Buscar si un archivo almacenado en google drive ,contiene una palabra
* Crear un archivo de texto plano sin contenido en google driver
* Un callback que permite generar el token para la autenticación con google drive , y almacenarlo en una variable de session.

__workflow__

<svg id="mermaid-1529796882204" width="100%" xmlns="http://www.w3.org/2000/svg" style="max-width: 784px;" viewBox="0 0 784 692.3999938964844">
<style>
#mermaid-1529796882204 .label {
  font-family: 'trebuchet ms', verdana, arial;
  color: #333; }

#mermaid-1529796882204 .node rect,
#mermaid-1529796882204 .node circle,
#mermaid-1529796882204 .node ellipse,
#mermaid-1529796882204 .node polygon {
  fill: #ECECFF;
  stroke: #9370DB;
  stroke-width: 1px; }

#mermaid-1529796882204 .node.clickable {
  cursor: pointer; }

#mermaid-1529796882204 .arrowheadPath {
  fill: #333333; }

#mermaid-1529796882204 .edgePath .path {
  stroke: #333333;
  stroke-width: 1.5px; }

#mermaid-1529796882204 .edgeLabel {
  background-color: #e8e8e8; }

#mermaid-1529796882204 .cluster rect {
  fill: #ffffde !important;
  stroke: #aaaa33 !important;
  stroke-width: 1px !important; }

#mermaid-1529796882204 .cluster text {
  fill: #333; }

#mermaid-1529796882204 div.mermaidTooltip {
  position: absolute;
  text-align: center;
  max-width: 200px;
  padding: 2px;
  font-family: 'trebuchet ms', verdana, arial;
  font-size: 12px;
  background: #ffffde;
  border: 1px solid #aaaa33;
  border-radius: 2px;
  pointer-events: none;
  z-index: 100; }

#mermaid-1529796882204 .actor {
  stroke: #CCCCFF;
  fill: #ECECFF; }

#mermaid-1529796882204 text.actor {
  fill: black;
  stroke: none; }

#mermaid-1529796882204 .actor-line {
  stroke: grey; }

#mermaid-1529796882204 .messageLine0 {
  stroke-width: 1.5;
  stroke-dasharray: '2 2';
  marker-end: 'url(#arrowhead)';
  stroke: #333; }

#mermaid-1529796882204 .messageLine1 {
  stroke-width: 1.5;
  stroke-dasharray: '2 2';
  stroke: #333; }

#mermaid-1529796882204 #arrowhead {
  fill: #333; }

#mermaid-1529796882204 #crosshead path {
  fill: #333 !important;
  stroke: #333 !important; }

#mermaid-1529796882204 .messageText {
  fill: #333;
  stroke: none; }

#mermaid-1529796882204 .labelBox {
  stroke: #CCCCFF;
  fill: #ECECFF; }

#mermaid-1529796882204 .labelText {
  fill: black;
  stroke: none; }

#mermaid-1529796882204 .loopText {
  fill: black;
  stroke: none; }

#mermaid-1529796882204 .loopLine {
  stroke-width: 2;
  stroke-dasharray: '2 2';
  marker-end: 'url(#arrowhead)';
  stroke: #CCCCFF; }

#mermaid-1529796882204 .note {
  stroke: #aaaa33;
  fill: #fff5ad; }

#mermaid-1529796882204 .noteText {
  fill: black;
  stroke: none;
  font-family: 'trebuchet ms', verdana, arial;
  font-size: 14px; }


#mermaid-1529796882204 .section {
  stroke: none;
  opacity: 0.2; }

#mermaid-1529796882204 .section0 {
  fill: rgba(102, 102, 255, 0.49); }

#mermaid-1529796882204 .section2 {
  fill: #fff400; }

#mermaid-1529796882204 .section1,
#mermaid-1529796882204 .section3 {
  fill: white;
  opacity: 0.2; }

#mermaid-1529796882204 .sectionTitle0 {
  fill: #333; }

#mermaid-1529796882204 .sectionTitle1 {
  fill: #333; }

#mermaid-1529796882204 .sectionTitle2 {
  fill: #333; }

#mermaid-1529796882204 .sectionTitle3 {
  fill: #333; }

#mermaid-1529796882204 .sectionTitle {
  text-anchor: start;
  font-size: 11px;
  text-height: 14px; }


#mermaid-1529796882204 .grid .tick {
  stroke: lightgrey;
  opacity: 0.3;
  shape-rendering: crispEdges; }

#mermaid-1529796882204 .grid path {
  stroke-width: 0; }


#mermaid-1529796882204 .today {
  fill: none;
  stroke: red;
  stroke-width: 2px; }



#mermaid-1529796882204 .task {
  stroke-width: 2; }

#mermaid-1529796882204 .taskText {
  text-anchor: middle;
  font-size: 11px; }

#mermaid-1529796882204 .taskTextOutsideRight {
  fill: black;
  text-anchor: start;
  font-size: 11px; }

#mermaid-1529796882204 .taskTextOutsideLeft {
  fill: black;
  text-anchor: end;
  font-size: 11px; }


#mermaid-1529796882204 .taskText0,
#mermaid-1529796882204 .taskText1,
#mermaid-1529796882204 .taskText2,
#mermaid-1529796882204 .taskText3 {
  fill: white; }

#mermaid-1529796882204 .task0,
#mermaid-1529796882204 .task1,
#mermaid-1529796882204 .task2,
#mermaid-1529796882204 .task3 {
  fill: #8a90dd;
  stroke: #534fbc; }

#mermaid-1529796882204 .taskTextOutside0,
#mermaid-1529796882204 .taskTextOutside2 {
  fill: black; }

#mermaid-1529796882204 .taskTextOutside1,
#mermaid-1529796882204 .taskTextOutside3 {
  fill: black; }


#mermaid-1529796882204 .active0,
#mermaid-1529796882204 .active1,
#mermaid-1529796882204 .active2,
#mermaid-1529796882204 .active3 {
  fill: #bfc7ff;
  stroke: #534fbc; }

#mermaid-1529796882204 .activeText0,
#mermaid-1529796882204 .activeText1,
#mermaid-1529796882204 .activeText2,
#mermaid-1529796882204 .activeText3 {
  fill: black !important; }


#mermaid-1529796882204 .done0,
#mermaid-1529796882204 .done1,
#mermaid-1529796882204 .done2,
#mermaid-1529796882204 .done3 {
  stroke: grey;
  fill: lightgrey;
  stroke-width: 2; }

#mermaid-1529796882204 .doneText0,
#mermaid-1529796882204 .doneText1,
#mermaid-1529796882204 .doneText2,
#mermaid-1529796882204 .doneText3 {
  fill: black !important; }


#mermaid-1529796882204 .crit0,
#mermaid-1529796882204 .crit1,
#mermaid-1529796882204 .crit2,
#mermaid-1529796882204 .crit3 {
  stroke: #ff8888;
  fill: red;
  stroke-width: 2; }

#mermaid-1529796882204 .activeCrit0,
#mermaid-1529796882204 .activeCrit1,
#mermaid-1529796882204 .activeCrit2,
#mermaid-1529796882204 .activeCrit3 {
  stroke: #ff8888;
  fill: #bfc7ff;
  stroke-width: 2; }

#mermaid-1529796882204 .doneCrit0,
#mermaid-1529796882204 .doneCrit1,
#mermaid-1529796882204 .doneCrit2,
#mermaid-1529796882204 .doneCrit3 {
  stroke: #ff8888;
  fill: lightgrey;
  stroke-width: 2;
  cursor: pointer;
  shape-rendering: crispEdges; }

#mermaid-1529796882204 .doneCritText0,
#mermaid-1529796882204 .doneCritText1,
#mermaid-1529796882204 .doneCritText2,
#mermaid-1529796882204 .doneCritText3 {
  fill: black !important; }

#mermaid-1529796882204 .activeCritText0,
#mermaid-1529796882204 .activeCritText1,
#mermaid-1529796882204 .activeCritText2,
#mermaid-1529796882204 .activeCritText3 {
  fill: black !important; }

#mermaid-1529796882204 .titleText {
  text-anchor: middle;
  font-size: 18px;
  fill: black; }

#mermaid-1529796882204 g.classGroup text {
  fill: #9370DB;
  stroke: none;
  font-family: 'trebuchet ms', verdana, arial;
  font-size: 10px; }

#mermaid-1529796882204 g.classGroup rect {
  fill: #ECECFF;
  stroke: #9370DB; }

#mermaid-1529796882204 g.classGroup line {
  stroke: #9370DB;
  stroke-width: 1; }

#mermaid-1529796882204 .classLabel .box {
  stroke: none;
  stroke-width: 0;
  fill: #ECECFF;
  opacity: 0.5; }

#mermaid-1529796882204 .classLabel .label {
  fill: #9370DB;
  font-size: 10px; }

#mermaid-1529796882204 .relation {
  stroke: #9370DB;
  stroke-width: 1;
  fill: none; }

#mermaid-1529796882204 #compositionStart {
  fill: #9370DB;
  stroke: #9370DB;
  stroke-width: 1; }

#mermaid-1529796882204 #compositionEnd {
  fill: #9370DB;
  stroke: #9370DB;
  stroke-width: 1; }

#mermaid-1529796882204 #aggregationStart {
  fill: #ECECFF;
  stroke: #9370DB;
  stroke-width: 1; }

#mermaid-1529796882204 #aggregationEnd {
  fill: #ECECFF;
  stroke: #9370DB;
  stroke-width: 1; }

#mermaid-1529796882204 #dependencyStart {
  fill: #9370DB;
  stroke: #9370DB;
  stroke-width: 1; }

#mermaid-1529796882204 #dependencyEnd {
  fill: #9370DB;
  stroke: #9370DB;
  stroke-width: 1; }

#mermaid-1529796882204 #extensionStart {
  fill: #9370DB;
  stroke: #9370DB;
  stroke-width: 1; }

#mermaid-1529796882204 #extensionEnd {
  fill: #9370DB;
  stroke: #9370DB;
  stroke-width: 1; }

#mermaid-1529796882204 .commit-id,
#mermaid-1529796882204 .commit-msg,
#mermaid-1529796882204 .branch-label {
  fill: lightgrey;
  color: lightgrey; }
</style><style>#mermaid-1529796882204 {
    color: rgba(0, 0, 0, 0.65);
    font: ;
  }</style><g transform="translate(-12, -12)"><g class="output"><g class="clusters"/><g class="edgePaths"><g class="edgePath" style="opacity: 1;"><path class="path" d="M301.75,60.82038834951456L259,86L259,111" marker-end="url(#arrowhead25930)" style="fill:none"/><defs><marker id="arrowhead25930" viewBox="0 0 10 10" refX="9" refY="5" markerUnits="strokeWidth" markerWidth="8" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" class="arrowheadPath" style="stroke-width: 1; stroke-dasharray: 1, 0;"/></marker></defs></g><g class="edgePath" style="opacity: 1;"><path class="path" d="M370.75,60.82038834951456L413.5,86L413.5,111" marker-end="url(#arrowhead25931)" style="fill:none"/><defs><marker id="arrowhead25931" viewBox="0 0 10 10" refX="9" refY="5" markerUnits="strokeWidth" markerWidth="8" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" class="arrowheadPath" style="stroke-width: 1; stroke-dasharray: 1, 0;"/></marker></defs></g><g class="edgePath" style="opacity: 1;"><path class="path" d="M362.5,149.14638783269962L282,177L282.5,202.5000030517578" marker-end="url(#arrowhead25932)" style="fill:none"/><defs><marker id="arrowhead25932" viewBox="0 0 10 10" refX="9" refY="5" markerUnits="strokeWidth" markerWidth="8" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" class="arrowheadPath" style="stroke-width: 1; stroke-dasharray: 1, 0;"/></marker></defs></g><g class="edgePath" style="opacity: 1;"><path class="path" d="M237.814589147534,316.21458609577616L133.5,395.8999938964844L133.5,431.3999938964844" marker-end="url(#arrowhead25933)" style="fill:none"/><defs><marker id="arrowhead25933" viewBox="0 0 10 10" refX="9" refY="5" markerUnits="strokeWidth" markerWidth="8" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" class="arrowheadPath" style="stroke-width: 1; stroke-dasharray: 1, 0;"/></marker></defs></g><g class="edgePath" style="opacity: 1;"><path class="path" d="M327.18541429612884,316.2145887556288L430.5,395.8999938964844L430.5,431.3999938964844" marker-end="url(#arrowhead25934)" style="fill:none"/><defs><marker id="arrowhead25934" viewBox="0 0 10 10" refX="9" refY="5" markerUnits="strokeWidth" markerWidth="8" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" class="arrowheadPath" style="stroke-width: 1; stroke-dasharray: 1, 0;"/></marker></defs></g><g class="edgePath" style="opacity: 1;"><path class="path" d="M430.5,472.3999938964844L430.5,507.8999938964844L430.5,543.3999938964844" marker-end="url(#arrowhead25935)" style="fill:none"/><defs><marker id="arrowhead25935" viewBox="0 0 10 10" refX="9" refY="5" markerUnits="strokeWidth" markerWidth="8" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" class="arrowheadPath" style="stroke-width: 1; stroke-dasharray: 1, 0;"/></marker></defs></g><g class="edgePath" style="opacity: 1;"><path class="path" d="M430.5,584.3999938964844L430.5,619.8999938964844L513.8616071428571,655.3999938964844" marker-end="url(#arrowhead25936)" style="fill:none"/><defs><marker id="arrowhead25936" viewBox="0 0 10 10" refX="9" refY="5" markerUnits="strokeWidth" markerWidth="8" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" class="arrowheadPath" style="stroke-width: 1; stroke-dasharray: 1, 0;"/></marker></defs></g><g class="edgePath" style="opacity: 1;"><path class="path" d="M610.1383928571429,655.3999938964844L693.5,619.8999938964844L693.5,563.8999938964844L693.5,507.8999938964844L693.5,451.8999938964844L693.5,395.8999938964844L693.5,281.1999969482422L693.5,177L464.5,139.7875" marker-end="url(#arrowhead25937)" style="fill:none"/><defs><marker id="arrowhead25937" viewBox="0 0 10 10" refX="9" refY="5" markerUnits="strokeWidth" markerWidth="8" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" class="arrowheadPath" style="stroke-width: 1; stroke-dasharray: 1, 0;"/></marker></defs></g></g><g class="edgeLabels"><g class="edgeLabel" style="opacity: 1;" transform=""><g transform="translate(0,0)" class="label"><foreignObject width="0" height="0"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;"><span class="edgeLabel"></span></div></foreignObject></g></g><g class="edgeLabel" style="opacity: 1;" transform=""><g transform="translate(0,0)" class="label"><foreignObject width="0" height="0"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;"><span class="edgeLabel"></span></div></foreignObject></g></g><g class="edgeLabel" style="opacity: 1;" transform=""><g transform="translate(0,0)" class="label"><foreignObject width="0" height="0"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;"><span class="edgeLabel"></span></div></foreignObject></g></g><g class="edgeLabel" style="opacity: 1;" transform="translate(133.5,395.8999938964844)"><g transform="translate(-6,-10.5)" class="label"><foreignObject width="12" height="21"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;"><span class="edgeLabel">SI</span></div></foreignObject></g></g><g class="edgeLabel" style="opacity: 1;" transform="translate(430.5,395.8999938964844)"><g transform="translate(-10,-10.5)" class="label"><foreignObject width="20" height="21"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;"><span class="edgeLabel">NO</span></div></foreignObject></g></g><g class="edgeLabel" style="opacity: 1;" transform="translate(430.5,507.8999938964844)"><g transform="translate(-162.5,-10.5)" class="label"><foreignObject width="325" height="21"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;"><span class="edgeLabel">Redirección de localhost a Google para autenticar</span></div></foreignObject></g></g><g class="edgeLabel" style="opacity: 1;" transform="translate(430.5,619.8999938964844)"><g transform="translate(-108,-10.5)" class="label"><foreignObject width="216" height="21"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;"><span class="edgeLabel">Redirección de Google a callback</span></div></foreignObject></g></g><g class="edgeLabel" style="opacity: 1;" transform="translate(693.5,451.8999938964844)"><g transform="translate(-94.5,-10.5)" class="label"><foreignObject width="189" height="21"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;"><span class="edgeLabel">Redirecciona a la Url Original</span></div></foreignObject></g></g></g><g class="nodes"><g class="node" style="opacity: 1;" id="A" transform="translate(336.25,40.5)"><rect rx="0" ry="0" x="-34.5" y="-20.5" width="69" height="41"/><g class="label" transform="translate(0,0)"><g transform="translate(-24.5,-10.5)"><foreignObject width="49" height="21"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;">Usuario</div></foreignObject></g></g></g><g class="node" style="opacity: 1;" id="B" transform="translate(259,131.5)"><rect rx="5" ry="5" x="-53.5" y="-20.5" width="107" height="41"/><g class="label" transform="translate(0,0)"><g transform="translate(-43.5,-10.5)"><foreignObject width="87" height="21"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;">Get/Post Item</div></foreignObject></g></g></g><g class="node" style="opacity: 1;" id="C" transform="translate(413.5,131.5)"><rect rx="5" ry="5" x="-51" y="-20.5" width="102" height="41"/><g class="label" transform="translate(0,0)"><g transform="translate(-41,-10.5)"><foreignObject width="82" height="21"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;">Get/Post File</div></foreignObject></g></g></g><g class="node" style="opacity: 1;" id="D" transform="translate(282,281.1999969482422)"><polygon points="79.2,0 158.4,-79.2 79.2,-158.4 0,-79.2" rx="5" ry="5" transform="translate(-79.2,79.2)"/><g class="label" transform="translate(0,0)"><g transform="translate(-57.5,-10.5)"><foreignObject width="115" height="21"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;">Está Autenticado?</div></foreignObject></g></g></g><g class="node" style="opacity: 1;" id="E" transform="translate(133.5,451.8999938964844)"><rect rx="5" ry="5" x="-113.5" y="-20.5" width="227" height="41"/><g class="label" transform="translate(0,0)"><g transform="translate(-103.5,-10.5)"><foreignObject width="207" height="21"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;">Continuar Con la URL Ingresada</div></foreignObject></g></g></g><g class="node" style="opacity: 1;" id="F" transform="translate(430.5,451.8999938964844)"><rect rx="5" ry="5" x="-133.5" y="-20.5" width="267" height="41"/><g class="label" transform="translate(0,0)"><g transform="translate(-123.5,-10.5)"><foreignObject width="247" height="21"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;">Guardar Url Ingresada y datos del post</div></foreignObject></g></g></g><g class="node" style="opacity: 1;" id="G" transform="translate(430.5,563.8999938964844)"><rect rx="5" ry="5" x="-82" y="-20.5" width="164" height="41"/><g class="label" transform="translate(0,0)"><g transform="translate(-72,-10.5)"><foreignObject width="144" height="21"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;">El usuario se autentica</div></foreignObject></g></g></g><g class="node" style="opacity: 1;" id="I" transform="translate(562,675.8999938964844)"><rect rx="5" ry="5" x="-142" y="-20.5" width="284" height="41"/><g class="label" transform="translate(0,0)"><g transform="translate(-132,-10.5)"><foreignObject width="264" height="21"><div xmlns="http://www.w3.org/1999/xhtml" style="display: inline-block; white-space: nowrap;">Generar y guardar token a partir del code</div></foreignObject></g></g></g></g></g></g></svg>

* Si el usuario quiere crear o buscar un item accede directamente a la URL
* Si el usuario quiere crear o buscar un archivo:
  - Si se encuentra Autenticado accede a la url indicada
  - Si no se encuentra Autenticado , el middleware  primero guardar la url original y los datos para crear un archivo (si los tuviera), lo redirecciona a la URL de google para autenticar, luego de que se autentique es redireccionado a callback.

* El callback recupera la url original y guarda el código para generar el token y redirecciona a la url original
 * El `GET` actúa normalmente ya que los parámetros pasan por url
 * El `POST` verifica si tiene en una variable session un error en los parámetros o un json con el título y la descripción, si no los tiene lo toma de el context

* ### Conexión con MYSQL
    ##### ¿Cuál fue el problema?
    Había que realizar el cambio de base de datos de Sqlite a Mysql server.

    ##### ¿Cuál fue la solución?
    Para realizar el cambio de base de datos no hubo mayores inconvenientes solamente fue crear el nuevo esquema, cambiar una línea para la conexión a la base de datos y agregar un nuevo repositorio `github.com/go-sql-driver/mysql` dentro del `Dockerfile`.
    El mayor tiempo utilizado fue investigar sobre los repositorios disponibles y ver documentación sobre la configuración para la conexión a Mysql

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

* ### Autenticación con Oauth2
  ##### ¿Cuál fue el problema?
Uno de los principales problemas es como manejar la autenticación,las librerías disponibles para utilizar con go.
  ##### ¿Cuál fue la solución?
La forma de cómo resolver la autenticación fue unas de las tareas que requirió más tiempo del pensado,debido a que para tratar de mantener los Endpoints originales propuestos en el challenge decidí no crear 2 endpoints adicionales los cuales serían para loguearse y para desloguearse , también tomé la decisión de que cuando uno no tiene la autorización para acceder a los servicios del drive , lo redireccione la url que brinda google para autenticar al usuario esto conllevo a una mayor dificultad para probar con Postman.
Como en el caso anterior la mayoría del tiempo para implementar la solución la encontré en primero aprender cómo funciona los modelo de autenticación , luego en la búsqueda de documentación de cómo se encontraba implementado en go.

* ### La Api de Google Drive
  ##### ¿Cuál fue el problema?
  * La utilización de la libreria `google.golang.org/api/drive/v3`,
  * Buscar si un archivo tiene una palabra
  * Crear un archivo en google drive

  ##### ¿Cuál fue la solución?
Al inicio me habia basado en el modelo que propone en  https://developers.google.com/drive/api/v3/quickstart/go , debo reconocer que comprender en su totalidad lo que hacia ese ejemplo me llevo tiempo,no obstante debo mensionar que  no me convencia la forma en la cual recibia el código de autorización ,y posteriormente creaba y almacenaba en un archivo plano un token en el servidor, esto conllevo a que vea la forma para autenticar, aprendiera un poco sobre el funcionamiento del oauth2.
Cuando logre obtener el token necesario para generar un cliente y este me provea del servicio del drive, note que a pesar que en la documentación está claramente explicado cómo es la utilización de rest api de google drive, al utilizar  `google.golang.org/api/drive/v3` tuve que investigar cómo estaban implementados los servicios que provee google drive, y hasta en algunas ocasiones por la falta de documentación para la librería , fue directamente leer el código ,entender cómo funcionaba y cómo implementarlo.
Buscar si un archivo tiene una palabra esto fue algo que me sorprendió , no hay forma directa para saber si un archivo tiene esa palabra en su contenido,título o descripción, la única forma por la documentación de google que encontré fue ,la solución que opte para implementar consiste en primero verificar si existe el id ingresado, si existe  obtener una lista de archivos los cuales tiene la palabra indicada agregando el filtro que tengan el mismo nombre que el archivo ingresado y fecha de creación , y buscar dentro de la lista de archivos si existe alguno con la id ingresada.
Para la creación de dentro del drive, no hubo mayores dificultades , solamente la decisión de no crear un archivo plano y subirlo al drive, debido a que no tenia ningun contenido para agregar al archivo , asi que directamente creo un archivo dentro del drive con los parámetros pasados (Título y Descripción).

* ### El manejo de sessiones
##### ¿Cuál fue el problema?
* Una manera de manera de preservar la url y los posibles datos ingresados
* Librerías disponibles para utilizar en go
* No persistencia de datos en la sesión
##### ¿Cuál fue la solución?
Una dificultad con la cual tuve que lidiar fue que al momento de autenticar con google al redireccionar pierdo el contexto , por lo tanto perdía la url a la cual se había accedido originalmente , además en el caso de que sea un post se pierden los datos ingresados. Para solucionar eso necesitaba una forma de persistir esos datos, una solución pudo ser almacenar esos datos en una cookie , pero como tambien tenia la intención de almacenar el token que obtenía al autenticar en google opte por almacenarlos en sessions, otra alternativa hubiera sido almacenarlos en un servidor estilo Redis pero esta opción cambiaria la arquitectura propuesta.
Por alguna extraña razón tuve la complicación de que en la sesión no tenía persistencia , luego de muchos intentos , muchas horas y borrar la session en el navegador que estaba utilizando , funcionó correctamente.


## Conclusiones
Realmente en la dificultad del proyecto la pude encontrar en la utilización de las tecnologías propuestas con las cuales nunca antes había trabajado como por ejemplo docker , o tuve que configurar un servidor como nginx. Programar en Golang fue volver a pensar en estructurado, tener que investigar mucho más de lo común y hasta incluso leer código fuente de las librerías para poder entender y aplicar las funcionalidades que posee o que la comunidad brinda. pero a su vez me permitió experimentar con un lenguaje con nuevas reglas como, que una función puede retornar 2 valores , o la facilidad con la cual se puede levantar un servidor http, o lo preparado que esta para la concurrencia,a su vez con `gin gonic`  vi lo fácil que es manejar rutas y validaciones, que la documentación que presentan es más que útil.
Fue la primera vez que me tocó programar en una REST API, entender la metodología ,el funcionamiento general, cómo funciona la autenticación, en resumen puedo asegurar que fue reto en varios niveles al enfrentarme a con tantas tecnologías , técnicas y metodologías  pero así mismo me dejó un montón de nuevos conocimientos, preguntas y  ganas de seguir aprendiendo,jugando y experimentando sobre los mismos.
