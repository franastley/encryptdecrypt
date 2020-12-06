# Encrypt and Decrupt 
This is a Go webservice that has 2 endpoints that can encrypt and ecrypt a value received in a JSON Value Key.

## Requirements
- An Ubuntu Server
- Docker and Docker Compose
- A domain name for the Nginx Server

## Installation
We navigate to the user folder
```bash
cd ~/
```
Then clone this repository.
```bash
git clone https://github.com/franastley/encryptdecrypt/
```
this will create a folder named encryptdecrypt with all the contents of the repository. We need to navigate to this folder.

```bash
cd ~/encryptdecrypt 
```

We will use Nginx and LetsEncryptIt as a secured proxy. These will run in a seperate container running in the background. For this we will use the following command.

```bash
docker-compose -f nginx-proxy-compose.yaml up -d
```
The nginx-proxy-compose.yaml is already in the Repository. 

After this we will create a new container with the actual go webservice but before you will need to update the go-app-compose.yaml file, for this use
```bash
nano go-app-compose.yaml
```
```bash
version: '2'
services:
  go-web-app:
    restart: always
    build:
      dockerfile: Dockerfile
      context: .
    environment:
      - VIRTUAL_HOST=<host>
      - LETSENCRYPT_HOST=<host>
```
replace <host> for the host you are going to use. In the repository it is set to encrypt.facturatek.site.
Now we can create our container for the Go App using 
```bash 
docker-compose -f go-app-compose.yaml up -d
```
If you navigate to your host via browser you will be able to see with a SSL Certification thanks to LetsEncryptIt 

Homepage. You can enjoy our encrypt and decrypt rest api endpoints.

For the endpoints PostMan was used.
### Encrypt 
Request: https://host//encrypt
Body: {"Value":"test"} 

You should get a encrypted string as a response

### Decrypt 
Request: https://host//decrypt
Body: {"Value":"950e5cabcc9502fe1885a0b5bbce7abfa7e4604e008e1665933b2effcb4447bf"} 

You should get a test as your response
