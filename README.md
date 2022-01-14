<p align="center">
  <h1 align="center">User Voice</h1>
  <div align="center">
    <strong>Platform to collect and organize customer feedback about BCC applications.</strong>
  </div>
  <div align="center">BCC members can share, vote and discuss on suggestions they have to make BCC applications even better.</div>
  <div align="center">It is a slightly modified version of the <a href="https://github.com/getfider/fider">open source project</a> <a href="https://getfider.com">Fider.</a></div>
</p>

<div align="center">
  <sub>Built with ❤︎ by <a href="https://github.com/goenning">Guilherme Oenning</a> and <a href="https://github.com/getfider/fider/graphs/contributors">contributors</a></sub>
</div>

<br />

<img src="etc/homepage.png">

## Modifications

- Allow new registrations when "Private site" is enabled
- Immediately redirect to the only OAuth provider

## Run it locally
To run uservoice on your machine you need to have docker installed. Then you can run one container for the postgres database and another container for user voice.
- open a command line in the root folder of this repository
- run `docker build -t postgres-db ./scripts`
- run `docker run -d --name my-postgresdb-container -p 5555:5432 postgres-db`
- `docker build -t uservoice ./`
- lookup the ip address of your postgres instance using `docker inspect my-postgresdb-container` ([detailed explanation](https://www.tutorialworks.com/container-networking/))
- adjust the ip address in the following command and run it `docker run -it --rm -p 3000:3000 -e HOST_DOMAIN='localhost' -e GO_ENV='development' -e JWT_SECRET='hsjl]W;&ZcHxT&FK;s%bgIQF:#ch=~#Al4:5]N;7V<qPZ3e9lT4a%;go;LIkc%k' -e EMAIL_NOREPLY='noreply@bcc.no' -e EMAIL_SMTP_HOST='localhost' -e EMAIL_SMTP_PORT=1026 -e EMAIL_SMTP_USERNAME='' -e EMAIL_SMTP_PASSWORD='' -e DATABASE_URL='postgres://postgres:docker@172.17.0.2:5432/fider?sslmode=disable' --name uservoice uservoice`
- go to http://localhost:3000 and configure your uservoice/fider instance.
- to be able to use it you need to verify your email address. However the above script does not have a valid smtp server configured and therefore will not send the required email. You can either provide a working smtp server or you need to look up the verification url in the database by querying the database.