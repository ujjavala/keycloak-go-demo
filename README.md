# keycloak-go-demo

This repository contains go code, when run, would authenticate a user via keycloak, provided keycloak (configured) is running on your local.

## Prerequisites:
 You need to have docker and go installed

## Configuring keycloak

You need keycloak docker image to start the server. For this execute
`
docker run -d -p 8080:8080 -e KEYCLOAK_USER=keycloak -e KEYCLOAK_PASSWORD=k --name keycloak jboss/keycloak:4.1.0.Final
`
Once you have this sorted, your container would be running on localhost:8080/auth

1. open localhost:8080/auth in your browser, navigate to the administration console and login with username keycloak and k as the corresponding password.
2. create a realm named demo
3. turn off the requirement of ssl for this realm. You can find this under realm settings -> login
4. create a client named demo-client and change the "Access Type" to confidential
5. configure the demo-client to be confidential and use http://localhost:8181/demo/callback as a valid redirect URI.
6. create a user named **demo** with password **demo** . Make sure to activate and impersonate this user.
7. Under client, note the client secret and keep this handy. You will need this in the go code.

**In the go code in demo.go replace the clientSecret value with the above obtained handy secret
**

## Run and Test

1. Navigate to this cloned repo and then run the code by executing `go run demo.go`. Make sure to execute `go mod tidy` just to be sure all dependencies are sorted. 


Navitage to http://localhost:8181. The request should reach your go server, which tries to authenticate you. Since you did not send a token, the go server will redirecty you to keycloak to authenticate. You should see the login screen of keycloak. Login with the demo user you have created for this realm (demo/demo). If you have configured your keycloak correctly, it will authenticate you and redirect you to your go server callback.

Copy your access token and use curl to verify if the server is able to accept your token

You can first export your token by executing:

`export TOKEN="eyJhbG..."
`
. Now curl the request by executing 
`curl -H "Authorization: Bearer $TOKEN" localhost:8181
`.

Your output would look like this:
**Yayy!!!Verified user id_token through keycloak.**

